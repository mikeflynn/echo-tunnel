package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

var origin = "http://localhost/"
var url = "ws://localhost:8888/join"

func main() {
	clientID := flag.String("id", "", "Client ID")
	flag.Parse()

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte(*clientID)
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent: %s\n", message)

	for {
		var msg = make([]byte, 512)
		_, err = ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received: %s\n", msg)
	}
}
