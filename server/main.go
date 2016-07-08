package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	//"github.com/mikeflynn/go-alexa/skillserver"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Conn struct {
	ws   *websocket.Conn
	send chan []byte
}

func (c *Conn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *Conn) writePump() {
	defer c.ws.Close()

	for {
		message, ok := <-c.send
		if !ok {
			c.write(websocket.CloseMessage, []byte{})
			return
		}

		c.ws.SetWriteDeadline(time.Now().Add(writeWait))
		w, err := c.ws.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var index map[string]*Conn = map[string]*Conn{}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	mt, message, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Received: %s\n", message)

	conn := &Conn{send: make(chan []byte, 256), ws: ws}
	go conn.writePump()

	index[string(message)] = conn

	fmt.Println("Current clients:")
	for id, _ := range index {
		fmt.Println(id)
	}
	fmt.Println("---------------")

	err = ws.WriteMessage(mt, []byte("Welcome, "+string(message)))
	if err != nil {
		log.Println(err)
	}
}

func pushMessages() {
	for {
		for cid, conn := range index {
			msg := []byte(cid + ": " + RandStringBytes(5))
			conn.send <- msg
			d, _ := time.ParseDuration("5s")
			time.Sleep(d)
		}
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	go pushMessages()

	http.HandleFunc("/join", joinHandler)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
