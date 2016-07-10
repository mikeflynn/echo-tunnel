package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

// Flags
var server = flag.String("server", "localhost", "The hostname for your Echo Tunnel server.")
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

func getServer() string {
	return "ws://" + *server + ":80/client/join"
}

func connect(clientID string) {
	ws, err := websocket.Dial(getServer(), "", "http://localhost/")
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
	Debug(fmt.Sprintf("Sent: %s\n", message))

	for {
		var msg = make([]byte, 512)
		_, err = ws.Read(msg)
		if err != nil {
			Debug(err.Error())
			return
		}
		Debug(fmt.Sprintf("Received: %s\n", msg))
	}
}

func Debug(msg string) {
	if *verbose {
		log.Println(msg)
	}
}
