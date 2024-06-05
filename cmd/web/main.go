package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("POST /snippet/create", snippetCreate)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Println(err)
}
