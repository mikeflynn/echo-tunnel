package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline                    = []byte{'\n'}
	space                      = []byte{' '}
	connIdx   map[string]*Conn = map[string]*Conn{}
	UserIndex *Index           = &Index{}
)

func genUserID() int64 {
	now := time.Now().Unix()
	rand.Seed(now)
	return int64(math.Ceil(float64(time.Now().Unix() * rand.Int63n(999999) / 100000000)))
}

type Index struct {
	Users []*User `json:"users"`
}

func (this *Index) fileLoad(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, this); err != nil {
		return err
	}

	return nil
}

func (this *Index) fileSave(filename string) error {
	data, err := json.Marshal(this)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}

func (this *Index) getUserByAmazonID(amazonID string) (*User, error) {
	for _, user := range this.Users {
		if user.AmazonID == amazonID {
			return user, nil
		}
	}

	return &User{}, errors.New("User not found.")
}

func (this *Index) getUserByID(id int64) (*User, error) {
	for _, user := range this.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return &User{}, errors.New("User not found.")
}

func (this *Index) addUser(amazonID string) *User {
	user := &User{
		ID:       genUserID(),
		AmazonID: amazonID,
	}

	this.Users = append(this.Users, user)

	return user
}

type User struct {
	Name     string    `json:"name"`
	ID       int64     `json:"id"`
	AmazonID string    `json:"amazon_id"`
	Clients  []*Client `json:"clients"`
}

func (this *User) getClient(clientName string) (*Client, error) {
	for _, client := range this.Clients {
		if client.Name == clientName {
			return client, nil
		}
	}

	return &Client{}, errors.New("Client not found.")
}

func (this *User) clientExists(clientName string) bool {
	for _, client := range this.Clients {
		if client.Name == clientName {
			return true
		}
	}

	return false
}

func (this *User) addClient(clientName string, conn *Conn) error {
	c, err := this.getClient(clientName)
	if err == nil {
		c.Connection = conn
		c.LastConnection = time.Now().Unix()
		return nil
	}

	this.Clients = append(this.Clients, &Client{
		Connection:     conn,
		LastConnection: time.Now().Unix(),
		Name:           clientName,
	})

	return nil
}

type Client struct {
	Name           string `json:"name"`
	LastConnection int64  `json:"last_connection"`
	Connection     *Conn
}

func (this *Client) isActive() bool {
	timeout, _ := time.ParseDuration("-5m")
	if this.LastConnection > time.Now().Add(timeout).Unix() {
		return true
	}

	return false
}

type Conn struct {
	ws           *websocket.Conn
	send         chan []byte
	readCallback func(*Conn, string)
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

func (c *Conn) readPump() {
	defer func() {
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				Debug(fmt.Sprintf("ERROR: %v", err))
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		Debug(string(message))

		if c.readCallback != nil {
			c.readCallback(c, string(message))
		}
	}
}
