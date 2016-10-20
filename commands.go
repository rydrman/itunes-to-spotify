package main

import "github.com/gizak/termui"

var quitCommands = []string{
    "q", "quit", "exit",
}

// RunCommand runs the given command string
func RunCommand(cmd string) {
    if stringInSlice(cmd, quitCommands) {
        termui.StopLoop()
    }
}

func stringInSlice(s string, list []string) bool {
    for _, m := range list {
        if m == s {
            return true
        }
    }
    return false
}
