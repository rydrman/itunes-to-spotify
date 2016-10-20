package main

import (
    "strings"

    "github.com/gizak/termui"
)

var quitCommands = []string{
    "q", "quit", "exit",
}

// RunCommand runs the given command string
func RunCommand(query string) {

    // check quit commands first and foremost
    if StringInSlice(query, quitCommands) {
        termui.StopLoop()
    }

    parts := strings.Split(query, " ")
    cmd := parts[0]

    Console.Log("command received: " + cmd)

}
