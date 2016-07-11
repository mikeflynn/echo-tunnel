package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/yannk/gosx-notifier"
)

var EventList map[string]*Event = map[string]*Event{
	"notify": {
		Description:    "Fires a notification to the screen.",
		ArgDescription: "<body> <title>",
		Fn: func(args ...string) string {
			if len(args) < 2 {
				return "Not enough arguments."
			} else {
				title := args[1]
				if title == "" {
					title = "Echo Tunnel"
				}

				n := &Notification{
					Title: title,
					Body:  args[0],
				}

				err := n.notify()
				if err != nil {
					return err.Error()
				}

				return "Notification sent."
			}
		},
	},
	"lock": {
		Description:    "Activates the lock screen.",
		ArgDescription: "",
		Fn: func(args ...string) string {
			result, err := termCommand("/System/Library/CoreServices/Menu Extras/User.menu/Contents/Resources/CGSession", "-suspend")
			if err != nil {
				return err.Error()
			}

			return result
		},
	},
	"open": {
		Description:    "Opens an app, could be in the background.",
		ArgDescription: "<app name> <background flag>",
		Fn: func(args ...string) string {
			var inBackground bool

			if len(args) > 1 && args[1] != "0" {
				inBackground = true
			} else if len(args) > 0 {
				inBackground = false
			} else {
				return "Not enough arguments."
			}

			var cmd *exec.Cmd
			if inBackground {
				cmd = exec.Command("open", "-a", args[0], "-g")
			} else {
				cmd = exec.Command("open", "-a", args[0])
			}

			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				return err.Error()
			}

			if list := out.String(); err != nil {
				return list
			}

			return "App opened."
		},
	},
	"close": {
		Description:    "Closes an app.",
		ArgDescription: "<app name>",
		Fn: func(args ...string) string {
			if len(args) == 0 {
				return "Not enough arguments."
			} else {
				actionScript(fmt.Sprintf("quit app \"%s\"", args[0]))
			}

			return "App closed."
		},
	},
	"brightness": {
		Description:    "Adjusts brightness level.",
		ArgDescription: "<brightness level 0 - 100; ex: 30>",
		Fn: func(args ...string) string {
			if len(args) == 0 {
				return "Not enough arguments."
			} else {
				i, _ := strconv.ParseFloat(args[0], 64)
				percent := strconv.FormatFloat(i/100, 'f', 2, 64)
				if output, err := storedActionScript("brightness.applescript", percent); err != nil {
					return err.Error()
				} else {
					return output
				}
			}
		},
	},
	"volume": {
		Description:    "Adjusts volume without UI",
		ArgDescription: "<volume % 0 - 100>",
		Fn: func(args ...string) string {
			if len(args) == 0 {
				return "Not enough arguments."
			} else {
				if output, err := actionScript(fmt.Sprintf("set volume output volume %s --100%", args[0])); err != nil {
					return err.Error()
				} else {
					return output
				}
			}
		},
	},
}

type Event struct {
	Description     string                 // Name of event for logging.
	ArgDescription  string                 // A description of what arguments it takes.
	Options         map[string]string      // Various options.
	Fn              func(...string) string // Run method
	FollowedBy      []*Event               // Links to events that may come directly after this one.
	LastOccured     uint64                 // Time stamp of last occurrence.
	DownTime        uint64                 // Min number of seconds between occurrences.
	TotalOccurances int                    // Total number of occurrences in this run.
}

func (this *Event) Run(args ...string) {
	fmt.Println(this.Fn(args...))
}

func storedActionScript(scriptName string, params ...string) (string, error) {
	var data []byte
	var err error

	tempDir := "/tmp/"

	// Pull from asset store
	if data, err = Asset("scripts/" + scriptName); err != nil {
		return "", err
	}

	if err = ioutil.WriteFile(tempDir+scriptName, []byte(data), 0644); err != nil {
		return "", err
	}

	if err = os.Chmod(tempDir+scriptName, 0777); err != nil {
		return "", err
	}

	cmd := exec.Command(tempDir+scriptName, params...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	list := out.String()
	return list, nil
}

func actionScript(command string) (string, error) {
	cmd := exec.Command("osascript", "-e", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	list := out.String()
	return list, nil
}

func termCommand(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	list := out.String()
	return list, nil
}

type Notification struct {
	Body     string
	Title    string
	Subtitle string
	Image    string
	Icon     string
}

func (this *Notification) notify() error {
	note := gosxnotifier.NewNotification(this.Body)
	note.Title = this.Title
	note.Subtitle = this.Subtitle
	note.Group = "com.echotunnel"

	if this.Icon != "" {
		note.AppIcon = this.Icon
	}

	if this.Image != "" {
		note.ContentImage = this.Image
	}

	note.Sound = gosxnotifier.Basso

	return note.Push()
}
