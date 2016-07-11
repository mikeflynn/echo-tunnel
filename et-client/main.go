package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// Flags
var server = flag.String("server", "localhost", "The hostname for your Echo Tunnel server.")
var clientID = flag.String("id", "", "Client ID")
var verbose = flag.Bool("v", false, "Verbose logging")

var CmdChan chan string = make(chan string)

func main() {
	flag.Parse()

	go cmdPipe()

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

		CmdChan <- string(msg)
	}
}

func cmdPipe() {
	for {
		message, ok := <-CmdChan
		if !ok {
			continue
		}

		args := strings.Split(message, " ")

		if args[0] == "Welcome," {
			continue
		}

		Debug("Command: " + message)

		if event, ok := EventList[args[0]]; ok {
			event.Run(args[1:]...)
		} else {
			EventList["notify"].Fn("Invalid command from Echo Tunnel.")
		}
	}
}

/*
func words2Int(cmd string) string {
	numbers := map[string]int{
		"zero":      0,
		"one":       1,
		"two":       2,
		"three":     3,
		"four":      4,
		"five":      5,
		"six":       6,
		"seven":     7,
		"eight":     8,
		"nine":      9,
		"ten":       10,
		"eleven":    11,
		"twelve":    12,
		"thirteen":  13,
		"fourteen":  14,
		"fifteen":   15,
		"sixteen":   16,
		"seventeen": 17,
		"eighteen":  18,
		"nineteen":  19,
		"twenty":    20,
		"thirty":    30,
		"forty":     40,
		"fifty":     50,
		"sixty":     60,
		"seventy":   70,
		"eighty":    80,
		"ninety":    90,
		"hundred":   100,
	}

	for old, new := range stopwords {
		cmd = strings.Replace(cmd, old, new, -1)
	}

	return cmd
}

func cmdStopWords(cmd) string {
	stopwords := map[string]string{
		"set":     "",
		"the":     "",
		"to":      "",
		"please":  "",
		"percent": "",
	}

	for old, new := range stopwords {
		cmd = strings.Replace(cmd, old, new, -1)
	}

	return cmd
}
*/
func Debug(msg string) {
	if *verbose {
		log.Println(msg)
	}
}
