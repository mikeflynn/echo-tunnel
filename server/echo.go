package main

import (
	"github.com/gorilla/mux"
)

func EchoRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/echo/test", EchoHelloWorld)
	return router
}

// Request Types

type EchoRequest struct {
	Version string      `json:"version"`
	Session EchoSession `json:"session"`
	Request EchoReqBody `json:"request"`
}

type EchoSession struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes struct {
		String map[string]interface{} `json:"string"`
	} `json:"attributes"`
	User struct {
		UserID string `json:"string"`
	} `json:"user"`
}

type EchoReqBody struct {
	Type      string     `json:"type"`
	RequestID string     `json:"requestId"`
	Timestamp string     `json:"timestamp"`
	Intent    EchoIntent `json:"intent,omitempty"`
	Reason    string     `json:"reason,omitempty"`
}

type EchoIntent struct {
	Name  string              `json:"name"`
	Slots map[string]EchoSlot `json:"slots"`
}

type EchoSlots struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Response Types

type EchoResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes"`
	Response          EchoRespBody           `json:"response"`
}

type EchoRespBody struct {
	OutputSpeech     EchoRespPayload `json:"outputSpeech"`
	Card             EchoRespPayload `json:"card"`
	Reprompt         EchoRespPayload `json:"reprompt"`
	ShouldEndSession bool            `json:"shouldEndSession"`
}

type EchoRespPayload struct {
	Type    string `json:"type"`
	Title   string `json:"string,omitempty"`
	Content string `json:"string"`
}
