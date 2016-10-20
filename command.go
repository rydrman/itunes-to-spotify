package main

import (
    "fmt"

    "github.com/gizak/termui"
)

type command struct {
    element  *termui.Par
    focused  bool
    prompt   string
    value    string
    callback func(string, error)
}

// Command is a singleton instance class for managing
// the Spotr command input
var Command = newCommand()

// newCommand is not exported because we are forcing a singleton
func newCommand() *command {

    c := &command{
        prompt:  ">",
        focused: true,
        value:   "",
    }

    cmd := termui.NewPar("")
    cmd.Height = 3
    cmd.BorderFg = termui.ColorWhite

    c.element = cmd
    return c

}

// Element returns the termui element for this command box
func (c command) Element() termui.GridBufferer {
    c.Paint()
    return c.element
}

// Clear clears the current value of the command input
func (c *command) Clear() {

    c.value = ""
    c.Paint()
}

// GetInput registers a message for the user to respond to
// and a callback function when complete.
func (c *command) Prompt(msg string, cb func(string, error)) {
    c.prompt = msg + ">"
    c.callback = cb
    Refresh()
}

// RemovePrompt removes any current prompt and callback in the command window
func (c *command) RemovePrompt() {
    c.prompt = ">"
    c.callback = nil
    c.Paint()
}

// Focus is called when the user is focused on this command window
func (c *command) Focus() {
    c.focused = true
    c.Paint()
}

// Blur is called when the user navigates away from this command window
func (c *command) Blur() {
    c.focused = false
    if c.callback != nil {
        c.callback("", fmt.Errorf("user exited the prompt"))
    }
    c.RemovePrompt()
    c.Paint()
}

// Paint re-paints the command ui box with any new changes
func (c *command) Paint() {

    promptColor := "fg-cyan"
    if !c.focused {
        promptColor = "fg-white"
    }

    c.element.Text = fmt.Sprintf("[%s](%s)%s", c.prompt, promptColor, c.value)
    Refresh()
}

// HandleKeyboard is used to have this command box handle the
// given keyboard input event from termui
func (c *command) HandleKeyboard(e termui.EvtKbd) {

    switch e.KeyStr {

    case "<space>":
        c.value += " "
        c.Paint()

    case "C-8": //backspace
        if len(c.value) > 0 {
            c.value = c.value[:len(c.value)-1]
            c.Paint()
        }

    case "<escape>":
        Console.Warning("command did not expect to receive escape!")

    case "<enter>":
        if c.callback != nil {
            c.callback(c.value, nil)
            c.RemovePrompt()
        } else {
            RunCommand(c.value)
        }
        c.Clear()

    default:
        c.value += e.KeyStr
        c.Paint()

    }

}
