package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	port  string = "8080"
	token string
)

func init() {

	// overwrite port
	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}

	// set token
	if value, ok := os.LookupEnv("TOKEN"); ok {
		token = value
	} else {
		log.Fatal("request token undefined")
	}

}

func authenticate(handler http.Handler) http.Handler {

	// wrap handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// log incoming request
		log.Println(r.Method, r.URL, r.ContentLength)

		// crude authentication mechanism to match tokens
		if r.Header.Get("token") != token {
			http.Error(w, "invalid request token", http.StatusUnauthorized)
			return
		}

		// call original handler
		handler.ServeHTTP(w, r)

	})

}

func main() {

	// build router
	router := mux.NewRouter()

	// file routes
	router.Path("/file").Methods(http.MethodPost).HandlerFunc(CreateFile)
	router.Path("/file").Methods(http.MethodDelete).HandlerFunc(DestroyFile)

	// package routes
	router.Path("/package").Methods(http.MethodPost).HandlerFunc(InstallPackage)
	router.Path("/package").Methods(http.MethodDelete).HandlerFunc(RemovePackage)

	// service routes
	router.Path("/service").Methods(http.MethodPost).HandlerFunc(StartService)
	router.Path("/service").Methods(http.MethodDelete).HandlerFunc(StopService)

	// start server with authenticated routes
	http.ListenAndServe(net.JoinHostPort("", port), authenticate(router))

}
