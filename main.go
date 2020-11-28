package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Register() {
	r := mux.NewRouter()

	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Hello from Brain!\n</h1>")
	})

	http.ListenAndServe(":80", r)
}

func main() {
	Register()
}
