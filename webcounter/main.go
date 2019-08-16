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

	rediscounter "github.com/slavrd/go-redis-counter"
)

// command line flags
var bindAddr = flag.String("bind-addr", "0.0.0.0:8000", "server bind address")
var tplPath = flag.String("html-tpl", "html/index.gohtml", "path to the golang html template to use")
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
var redisConnInfo string // string to display on the web page
var htmlCounterTpl *template.Template

var counter *rediscounter.RedisCounter

func main() {
	flag.Parse()

	// set up global variables
	var err error
	htmlCounterTpl, err = loadTemplate(*tplPath)
	if err != nil {
		log.Fatalf("error loading html template: %v", err)
	}

	// generate redis connection info string to display on the webpage
	redisConnInfo = fmt.Sprintf("redis @ %s:%v, db: %v, key: %q", *redisHost, *redisPort, *redisDB, *redisKey)

	// initialize the server's RedisCounter instance
	counter, err = rediscounter.NewCounter(fmt.Sprintf("%s:%v", *redisHost, *redisPort), *redisPass, *redisKey, *redisDB)
	if err != nil {
		log.Fatal(err)
	}

	// setup server handlers
	http.Handle("/incr", newHandler(newIncrCtx, htmlCounterTpl))
	http.Handle("/get", newHandler(newGetCtx, htmlCounterTpl))

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

// newHandler returns a http.Handler func which renders the tpl template
// with a counterCtx incstace created using cf func
func newHandler(ctxf func(*rediscounter.RedisCounter) (*counterCtx, error), tpl *template.Template) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {

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
