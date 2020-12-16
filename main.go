package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ice1n36/brain/handlers"
)

func Register() {
	r := mux.NewRouter()

	r.HandleFunc("/hello", handlers.HelloHandler)
	r.HandleFunc("/api/v1/writeNetworkTraffic", handlers.WriteNetworkTrafficHandler).Methods("POST")

	http.ListenAndServe(":80", r)
	// TODO: implement graceful shutdown: https://github.com/gorilla/mux#graceful-shutdown
}

func main() {
	Register()
}
