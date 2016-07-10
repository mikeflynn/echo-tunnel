package main

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
)

var (
	newline                  = []byte{'\n'}
	space                    = []byte{' '}
	connIdx map[string]*Conn = map[string]*Conn{}
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
