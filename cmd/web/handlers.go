package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello from snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d\n", id)
}

func snippetCreate(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}
