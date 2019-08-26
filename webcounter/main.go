// Command webcounter starts a web server that implements the rediscounter package.
// Upon receiving a request the server will increment the counter's value and
// return a web page that displays the resulting value.
package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	rediscounter "github.com/slavrd/go-redis-counter"
)

// command line flags
var bindAddr = flag.String("bind-addr", "0.0.0.0:8000", "server bind address")
var tplPath = flag.String("html-tpl", "html/index.gohtml", "path to the golang html template to use for main page")
var mtplPath = flag.String("html-mtpl", "html/metrics.gohtml", "path to the golang html template to use for metrics page")
var redisHost = flag.String("redis-host", "127.0.0.1", "redis server host address")
var redisPort = flag.Int("redis-port", 6379, "redis server port")
var redisPass = flag.String("redis-pass", "", "redis server password")
var redisDB = flag.Int("redis-db", 0, "redis database index to use")
var redisKey = flag.String("redis-key", "count", "redis key to use")

// TODO: vault implementation
var useVault = flag.Bool("vault", false, "use vault server to retrieve password")
var vSecretPath = flag.String("vault-secret-path", "kv/redispassword", "vault path for redis password kv secret")
var vSecretKey = flag.String("vault-secret-key", "pass", "vault secret key for redis password")

// global variables
var redisAddr string     // host:port address for the redis server
var redisConnInfo string // string to display on the web page
var counter *rediscounter.RedisCounter
var htmlCounterTpl *template.Template // html template to render counter
var htmlMetricsTpl *template.Template // html template to render metrics data
var usageData *metrics

func init() {
	flag.Parse()

	// check which flags have been set
	sf := make(map[string]struct{})
	flag.Visit(func(f *flag.Flag) { sf[f.Name] = struct{}{} })

	// handle environment variable configuration
	// Priority is: passed flag > env var > default flag

	// for redisAddr we need to join redis-host and redis-port
	// or if using REDIS_ADDR check if it contains port
	if _, ok := sf["redis-host"]; !ok {
		if os.Getenv("REDIS_ADDR") != "" {
			redisAddr = os.Getenv("REDIS_ADDR")
		} else {
			redisAddr = *redisHost
		}
	} else {
		redisAddr = *redisHost
	}
	if !strings.ContainsRune(redisAddr, ':') {
		redisAddr = strings.Join([]string{redisAddr, strconv.Itoa(*redisPort)}, ":")
	}

	// set up redis-pass
	if _, ok := sf["redis-pass"]; !ok && os.Getenv("REDIS_PASS") != "" {
		*redisPass = os.Getenv("REDIS_PASS")
	}

	// set global variables

	var err error
	htmlCounterTpl, err = loadTemplate(*tplPath)
	if err != nil {
		log.Fatalf("error loading index html template: %v", err)
	}

	htmlMetricsTpl, err = loadMetricsTpl(*mtplPath)
	if err != nil {
		log.Fatalf("error loading metrics html template: %v", err)
	}

	redisConnInfo = fmt.Sprintf("redis @ %s, db: %v, key: %q", redisAddr, *redisDB, *redisKey)

	usageData = newMetrics(redisConnInfo)
}

func main() {

	// initialize the server's RedisCounter instance.
	// not done in init() as we need slightly different process for testing initialization
	var err error
	counter, err = rediscounter.NewCounter(redisAddr, *redisPass, *redisKey, *redisDB)
	if err != nil {
		log.Fatalf("error intializing global RedisCounter: %v", err)
	}

	// setup server handlers
	http.Handle("/incr", newHandler(newIncrCtx, htmlCounterTpl))
	http.Handle("/get", newHandler(newGetCtx, htmlCounterTpl))
	http.Handle("/health", newHealthHandler(counter.RedisHealth))
	http.Handle("/metrics", newMetricsHandler(usageData, htmlMetricsTpl))

	// start server
	log.Fatal(http.ListenAndServe(*bindAddr, nil))
}

// loadTemplate reads the file from tplPath and parses it as a go html template
func loadTemplate(tplPath string) (*template.Template, error) {

	content, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("indexCounter").Parse(string(content))
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

// loadMetricsTpl reads the file from mtplPath and parses it as a go html template
func loadMetricsTpl(mtplPath string) (*template.Template, error) {
	content, err := ioutil.ReadFile(mtplPath)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("metrics").
		Funcs(template.FuncMap{"timenow": time.Now}).
		Parse(string(content))
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

// newHandler returns a http.Handler func which renders the tpl template
// with a counterCtx incstace created using cf func
func newHandler(ctxf func(*rediscounter.RedisCounter) (*counterCtx, error), tpl *template.Template) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {

		usageData.add(r.URL.Path)

		ctx, err := ctxf(counter)
		if err != nil {
			log.Printf("error generating counter ctx: %v", err)
			w.WriteHeader(500)
			w.Write([]byte("Internal server error!"))
			return
		}

		err = tpl.Execute(w, ctx)
		if err != nil {
			log.Printf("error writing response: %v", err)
		}
	}

	return http.HandlerFunc(h)
}

// newHealthHandler returns a http.Handler func which will call the hcf func and
// will return 200 OK if result is nil or 500 Internal Server Error if result is an error
func newHealthHandler(hcf func() error) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {

		usageData.add(r.URL.Path)

		err := hcf()
		if err == nil {
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(500)
			log.Printf("error healthcheck: %v", err)
			w.Write([]byte("Redis server is down!"))
		}
	}

	return http.HandlerFunc(h)
}

// newMetricsHandler returns a http.Handler func which will render the provided tpl
// with the m context. It will respect the m.Mutex so will be concurency safe.
func newMetricsHandler(m *metrics, tpl *template.Template) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {

		m.add(r.URL.Path)

		m.Mutex <- struct{}{}
		err := tpl.Execute(w, m)
		<-m.Mutex
		if err != nil {
			log.Printf("error writing response: %v", err)
		}

	}

	return http.HandlerFunc(h)
}
