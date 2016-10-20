package main

import (
    "fmt"
    "strings"

    "github.com/gizak/termui"
)

type console struct {
    messages []string
    element  *termui.Par
}

// Console is a public singleton instance of the
// console class used for logging and output in Spotr
var Console = newConsole()

// newConsole not public because we are forcing a singleton here
func newConsole() *console {

    c := &console{
        messages: []string{"Welcome to Spotr!"},
    }

    elem := termui.NewPar(strings.Join(c.messages, "\n"))
    elem.Height = 8
    elem.BorderFg = termui.ColorWhite

    c.element = elem
    return c

}

// Element returns the termui element for this console
func (c console) Element() termui.GridBufferer {
    return c.element
}

// Log is used to log a message to the console
func (c *console) Log(msg string) {
    c.messages = append(c.messages, msg)
    c.element.Text = strings.Join(c.messages, "\n")
    Refresh()
}

func (c *console) Error(msg string) {
    msg = fmt.Sprintf("[%s](fg-red)", msg)
    c.Log(msg)
}
