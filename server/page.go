package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func PageRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", homePage)
	router.HandleFunc("/about", aboutPage)
	return router
}

func homePage(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Home Page!")
}

func aboutPage(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "About Page!")
}
