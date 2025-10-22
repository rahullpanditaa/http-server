package main

import (
	"log"
	"net/http"
)

func main() {
	// create http multiplexer - router that maps
	// url patterns to handlers
	mux := http.NewServeMux()

	// custom http handler function
	readinessHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	mux.HandleFunc("/healthz", readinessHandler)

	// fileserver handler that serves static files
	// from the CURRENT DIR
	fileServerHandler := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", fileServerHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
