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

func homePage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Home Page!")
}

func aboutPage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "About Page!")
}
