package main

import (
	"log"
	"net/http"
	"time"
)

// handleGet generates a response without changing the counter value
// but updateing the ime only
func handleGet(w http.ResponseWriter, r *http.Request) {

	htmlCounterCtx.Time = time.Now()

	err := htmlCounterTpl.Execute(w, htmlCounterCtx)
	if err != nil {
		log.Printf("error writing response: %v", err)
	}
}

// handleGet increases the counter value by 1 and generates a response
// with the resulting value and updated time
func handleIncr(w http.ResponseWriter, r *http.Request) {

	err := htmlCounterCtx.IncrBy(1)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		log.Printf("error generating response: %v", err)
	}

	err = htmlCounterTpl.Execute(w, htmlCounterCtx)
	if err != nil {
		log.Printf("error writing response: %v", err)
	}
}
