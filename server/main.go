package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

	fmt.Printf("NEW CLIENT: %s\n", message)

	conn := &Conn{send: make(chan []byte, 256), ws: ws}
	go conn.writePump()

	index[string(message)] = conn

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

func repl() {
	var line string
	var err error

	fmt.Printf("Terminal\n" + "==========================\n\n")

	for {
		fmt.Printf("> ")

		in := bufio.NewReader(os.Stdin)
		if line, err = in.ReadString('\n'); err != nil {
			log.Fatal(err)
		}

		commands := strings.Split(line, " ")
		switch {
		case strings.TrimSpace(commands[0]) == "list":
			fmt.Println("Current clients:")
			for name, _ := range index {
				fmt.Println("* " + name)
			}
			fmt.Println("================")
		case strings.TrimSpace(commands[0]) == "send":
			if ws, ok := index[commands[1]]; ok {
				ws.send <- []byte(strings.TrimSpace(strings.Join(commands[2:], " ")))
			} else {
				fmt.Println("Invalid client.")
			}
		default:
			fmt.Println("Invalid command. Try \"commands\" for a command list.")
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
	go func() {
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	repl()
}
