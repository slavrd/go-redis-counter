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
	"sync"
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

// global variables
var redisAddr string     // host:port address for the redis server
var redisConnInfo string // string to display on the web page

var ctrMu sync.Mutex // guards the global counter
var counter *rediscounter.RedisCounter

var htmlCounterTpl *template.Template // html template to render counter
var htmlMetricsTpl *template.Template // html template to render metrics data
var usageData *metrics

// wcInit performs initialization of the global variables based on passed command line flags
// cannot use init() as of go1.13 packages that call flag.Parse() there will cause tests to fail.
func wcInit() {
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

	wcInit()

	// initialize the server's RedisCounter instance.
	// not done in init() as we need slightly different process for testing initialization
	setGlobalCounter()

	// if VAULT_TOKEN is set up start a goroutine to retrieve the redis password from Vault
	// and update the global counter.
	if os.Getenv("VAULT_TOKEN") != "" {
		spath := os.Getenv("VAULT_RP_PATH")
		if spath == "" {
			spath = "secret/redispassword"
		}
		skey := os.Getenv("VAULT_RP_KEY")
		if skey == "" {
			skey = "pass"
		}
		go setVaultRedisPass(spath, skey)
	}

	// setup server handlers
	http.Handle("/incr", newHandler(newIncrCtx, htmlCounterTpl))
	http.Handle("/decr", newHandler(newDecrCtx, htmlCounterTpl))
	http.Handle("/get", newHandler(newGetCtx, htmlCounterTpl))
	http.Handle("/reset", newHandler(newResetCtx, htmlCounterTpl))
	http.Handle("/health", newHealthHandler())
	http.Handle("/metrics", newMetricsHandler(usageData, htmlMetricsTpl))
	http.Handle("/crash", newCrashHandler(log.Fatal, "/crash called, stopping server!"))

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

// setGlobalCounter sets the value for the counter global var
func setGlobalCounter() {
	var err error
	ctrMu.Lock()
	counter, err = rediscounter.NewCounter(redisAddr, *redisPass, *redisKey, *redisDB)
	ctrMu.Unlock()
	if err != nil {
		log.Printf("error intializing global RedisCounter: %v", err)
	}
}

// setVaultRedisPass sets the global counter with password retrieved from Vault.
// In case an error is received from Vault the func will auto retry in increasing intervals.
func setVaultRedisPass(spath, skey string) {
	vc, err := newVaultClient()
	if err != nil {
		log.Printf("error creating vault client: %v", err)
		return
	}

	var p string
	var pause = 1 * time.Second
	for p == "" {
		p, err = getRedisPass(vc, spath, skey)
		if err != nil {
			log.Printf("error retrieveing password from Vault: %v", err)
		}
		time.Sleep(pause)
		if pause < 60*time.Second {
			pause *= 2
		}
	}

	*redisPass = p
	setGlobalCounter()
	log.Printf("password successfully retrieved from Vault.")
}
