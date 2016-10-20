package main

import "github.com/gizak/termui"

type command struct {
    element *termui.Par
}

// Command is a singleton instance class for managing
// the Spotr command input
var Command = newCommand()

// newCommand is not exported because we are forcing a singleton
func newCommand() *command {

    c := &command{}

    cmd := termui.NewPar("")
    cmd.Height = 3
    cmd.BorderFg = termui.ColorWhite

    c.element = cmd
    return c

}

func (c command) Element() termui.GridBufferer {
    return c.element
}

func (c command) HandleKeyboard(e termui.EvtKbd) {

    switch e.KeyStr {

    case "<space>":
        c.element.Text += " "

    case "C-8": //backspace
        if len(c.element.Text) > 0 {
            c.element.Text = c.element.Text[:len(c.element.Text)-1]
        }

    case "<enter>":
        RunCommand(c.element.Text)
        c.element.Text = ""

    default:
        c.element.Text += e.KeyStr

    }

    Refresh()

}
