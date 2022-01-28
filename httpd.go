package main

import (
	"log"
	"net/http"
	"time"
)

func WithLogging(h http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        start := time.Now()

        h.ServeHTTP(rw, r) // serve the original request

		log.Printf("ts: %v, method: %s, url: %s, duration: %v\n",
			start, r.Method, r.RequestURI, time.Since(start))
    })
}

func main() {
	log.Fatal(http.ListenAndServe(":8080",
		WithLogging(http.FileServer(http.Dir("htdocs")))))
}
