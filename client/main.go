package main

import (
	"flag"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

// Config
var origin = "http://localhost/"
var url = "ws://localhost:8888/join"

// Flags
var clientID = flag.String("id", "", "Client ID")
var verbose = flag.Bool("v", false, "Verbose logging")

func main() {
	flag.Parse()

	for {
		connect(*clientID)
		Debug("Reconnecting in 10 seconds...")

		d, _ := time.ParseDuration("10s")
		time.Sleep(d)
	}
}

func connect(clientID string) {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		Debug(err.Error())
		return
	}

	message := []byte(clientID)
	_, err = ws.Write(message)
	if err != nil {
		Debug(err.Error())
		return
	}
	Debug("Sent: %s\n", message)

	for {
		var msg = make([]byte, 512)
		_, err = ws.Read(msg)
		if err != nil {
			Debug(err.Error())
			return
		}
		Debug("Received: %s\n", msg)
	}
}

func Debug(msg) {
	if *verbose {
		log.Println(msg)
	}
}
