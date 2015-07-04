package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func EchoRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/echo/", echoTest)
	router.HandleFunc("/echo/test", echoTest)
	return router
}

func echoTest(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Echo page")

	//foo := req.URL.Query().Get("foo")
}
