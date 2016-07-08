package main

import (
	"flag"
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

var origin = "http://localhost/"
var url = "ws://localhost:8888/join"

func connect(clientID string) {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	message := []byte(clientID)
	_, err = ws.Write(message)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Sent: %s\n", message)

	for {
		var msg = make([]byte, 512)
		_, err = ws.Read(msg)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Received: %s\n", msg)
	}
}

func main() {
	clientID := flag.String("id", "", "Client ID")
	flag.Parse()

	for {
		connect(*clientID)
		fmt.Println("Reconnecting in 10 seconds...")

		d, _ := time.ParseDuration("10s")
		time.Sleep(d)
	}
}
