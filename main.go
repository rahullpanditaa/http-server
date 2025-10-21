package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	assetsFileServer := http.FileServer(http.Dir("assets"))

	mux.Handle("/assets/", http.StripPrefix("/assets/", assetsFileServer))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
