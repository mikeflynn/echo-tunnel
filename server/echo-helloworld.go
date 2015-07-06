package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
)

func EchoHelloWorld(w http.ResponseWriter, r *http.Request) {
	var req EchoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HTTPError(w, err.Error(), "Bad Request", 400)
	}

	if req.Request.Type == "IntentRequest" || req.Request.Type == "LaunchRequest" {
		resp := &EchoResponse{
			Version: "1.0",
			Response: EchoRespBody{
				OutputSpeech: EchoRespPayload{
					Type: "PlainText",
					Text: "Hello world from my Echo test app!",
				},
				Card: EchoRespPayload{
					Type:    "Simple",
					Title:   "Hello World!",
					Content: "This is a test card.",
				},
				//				Reprompt: &EchoReprompt{
				//					OutputSpeech: EchoRespPayload{
				//						Type: "PlainText",
				//						Text: "blah blah blah?",
				//					},
				//				},
				ShouldEndSession: true,
			},
		}

		jsonStr, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(jsonStr)
	}
}
