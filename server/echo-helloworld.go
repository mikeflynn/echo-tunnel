package main

import (
	"fmt"
	"net/http"
)

func EchoHelloWorld(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Echo page")

	//foo := req.URL.Query().Get("foo")
}
