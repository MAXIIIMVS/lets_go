package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello from snippetbox"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", home)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Println(err)
}
