package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
			for _, user := range UserIndex.Users {
				for _, client := range user.Clients {
					if client.isActive() {
						fmt.Println(fmt.Sprintf("* %s:%s", user.ID, client.Name))
					}
				}
			}
			fmt.Println("================")
		case strings.TrimSpace(commands[0]) == "send":
			parts := strings.Split(commands[0], ":")

			uid, err := strconv.ParseInt(parts[0], 10, 64)
			user, err := UserIndex.getUserByID(uid)
			if err != nil {
				fmt.Println("Invalid client.")
			}

			if client, err := user.getClient(parts[1]); err == nil {
				client.Connection.send <- []byte(strings.TrimSpace(strings.Join(commands[2:], " ")))
			} else {
				fmt.Println("Invalid client.")
			}
		default:
			fmt.Println("Invalid command.")
		}
	}
}

func Debug(msg string) {
	if *verbose {
		log.Println(msg)
	}
}

func EchoLaunchHandler(req *alexa.EchoRequest, res *alexa.EchoResponse) {
	mongo, db := MongoConnect()
	GetSession(db, req.GetSessionID())
	defer mongo.Close()

	_, err := UserIndex.getUserByAmazonID(req.GetUserID())
	if err != nil {
		UserIndex.addUser(req.GetUserID())
		res.OutputSpeech("Welcome new user. You need some clients to connect to us first!").EndSession(true)
		return
	}

	res.OutputSpeech("Which computer do you want to connect to?").EndSession(false)
}

func EchoIntentHandler(req *alexa.EchoRequest, res *alexa.EchoResponse) {
	mongo, db := MongoConnect()
	defer mongo.Close()

	session := GetSession(db, req.GetSessionID())

	user, err := UserIndex.getUserByAmazonID(req.GetUserID())
	if err != nil {
		res.OutputSpeech("Sorry, we don't have an account setup for this user.").EndSession(true)
		return
	}

	switch req.GetIntentName() {
	case "SelectBox":
		target, _ := req.GetSlotValue("target")
		if target == "" {
			res.OutputSpeech("I didn't get that. Can you tell me the computer you want again?").EndSession(false)
			return
		}

		if !user.clientExists(target) {
			res.OutputSpeech("The computer you requested isn't online.").EndSession(true)
			return
		}

		session.Target = target
		session.Update(db)

		runCommand(session, user, res)
	case "RunCommand":
		cmd, err := req.GetSlotValue("cmd")
		if err != nil {
			res.OutputSpeech("I'm sorry, but what should I tell " + session.Target + "to do?").EndSession(false)
			return
		}

		session.Cmd = cmd

		payload, _ := req.GetSlotValue("payload")
		session.Payload = payload
		session.Update(db)

		runCommand(session, user, res)
	case "ListCommand":
		names := []string{}
		for _, client := range user.Clients {
			if client.isActive() {
				names = append(names, client.Name)
			}
		}

		res.OutputSpeech("Your available computers are: " + strings.Join(names, ", ")).EndSession(true)
	default:
		res.OutputSpeech("I'm sorry, I didn't understand your request.").EndSession(false)
	}
}

func runCommand(session *TunnelSession, user *User, res *alexa.EchoResponse) {
	if session.Target == "" {
		res.OutputSpeech("Which computer do you want to run this on?").EndSession(false)
		return
	}

	if session.Cmd == "" {
		res.OutputSpeech("What should I tell " + session.Target + "to do?").EndSession(false)
		return
	}

	client, err := user.getClient(session.Target)
	if err != nil {
		res.OutputSpeech("The computer you requested isn't online.").EndSession(true)
		return
	}

	client.Connection.send <- []byte(session.Cmd + " " + session.Payload)
	res.OutputSpeech("Done!").EndSession(true)
}

var Applications = map[string]interface{}{
	"/echo/tunnel": alexa.EchoApplication{ // Route
		AppID:    os.Getenv("ECHOTUNNEL_APP_ID"), // Echo App ID from Amazon Dashboard
		OnIntent: EchoIntentHandler,
		OnLaunch: EchoLaunchHandler,
	},
	"/client/join": alexa.StdApplication{
		Methods: "GET",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			ws, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				Debug("Upgrade Error:" + err.Error())
				return
			}

			mt, message, err := ws.ReadMessage()
			if err != nil {
				Debug(err.Error())
			}

			Debug(fmt.Sprintf("Client Connection: %s\n", message))

			conn := &Conn{send: make(chan []byte, 256), ws: ws}
			go conn.writePump()

			retMessage := ""

			parts := strings.Split(string(message), ":")

			uid, err := strconv.ParseInt(parts[0], 10, 64)
			user, err := UserIndex.getUserByID(uid)
			if err != nil {
				retMessage = "ERROR: User ID not found."

				defer close(conn.send)
			} else {
				user.addClient(parts[1], conn)

				retMessage = fmt.Sprintf("Welcome, %s for user %s", parts[1], parts[0])
			}

			if err = conn.write(mt, []byte(retMessage)); err != nil {
				Debug(err.Error())
			}
		},
	},
}

var verbose = flag.Bool("v", false, "Verbose logging.")
var startRepl = flag.Bool("repl", false, "Start with a repl.")
var port = flag.String("port", "8888", "Port number.")
var database = flag.String("db", "", "File path for the persistent JSON DB")

func main() {
	flag.Parse()

	if *database != "" {
		if err := UserIndex.fileLoad(*database); err != nil {
			log.Println(err.Error())
		}

		go func() {
			for {
				dur, _ := time.ParseDuration("5m")
				time.Sleep(dur)

				if err := UserIndex.fileSave(*database); err != nil {
					log.Println(err.Error())
				}
			}
		}()
	}

	if *startRepl {
		go alexa.Run(Applications, *port)
		repl()
	} else {
		alexa.Run(Applications, *port)
	}
}
