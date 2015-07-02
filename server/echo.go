package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func EchoRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", echoTest)
	router.HandleFunc("/test", echoTest)
	return router
}

func echoTest(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Echo page")
}
