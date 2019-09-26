package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	rediscounter "github.com/slavrd/go-redis-counter"
)

// newHandler returns a http.Handler func which renders the tpl template
// with a counterCtx incstace created using cf func
func newHandler(ctxf func(*rediscounter.RedisCounter) (*counterCtx, error), tpl *template.Template) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {

		usageData.add(r.URL.Path)
		ctrMu.Lock()
		ctx, err := ctxf(counter)
		ctrMu.Unlock()
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
			w.Write([]byte(fmt.Sprintf("Cannot connect to redis: %v", err)))
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

// newCrashHandler returns a http.Handler func which will call f(msg).
// Intended for use with log.Fatal to crash the server.
func newCrashHandler(f func(...interface{}), msg string) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {
		f(msg)
	}

	return http.HandlerFunc(h)

}
