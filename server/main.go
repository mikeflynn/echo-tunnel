package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	echoRouter := EchoRouter()
	router.PathPrefix("/echo/").Handler(negroni.New(
		negroni.HandlerFunc(ValidateRequest),
		negroni.Wrap(echoRouter),
	))

	pageRouter := PageRouter()
	router.PathPrefix("/").Handler(negroni.New(
		negroni.Wrap(pageRouter),
	))

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}

func ValidateRequest(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Checking request signature...")

	http.Error(rw, "Not Authorized", 401)
}

func verifyCertURL(path string) bool {
	return true
}
