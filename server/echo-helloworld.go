package main

import (
	"net/http"

	"github.com/gorilla/context"
)

func EchoHelloWorld(w http.ResponseWriter, r *http.Request) {
	echoReq := context.Get(r, "echoRequest").(*EchoRequest)

	if echoReq.GetRequestType() == "IntentRequest" || echoReq.GetRequestType() == "LaunchRequest" {
		echoResp := NewEchoResponse().OutputSpeech("Hello world from my new Echo test app!").Card("Hello World", "This is a test card.")

		json, _ := echoResp.String()
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(json)
	}
}
