package main

import (
    "fmt"

    "github.com/rydrman/termui"
)

type console struct {
    messages []string
    element  *termui.List
}

type message struct {
    content string
    expires int
    expired bool
}

// Console is a public singleton instance of the
// console class used for logging and output in Spotr
var Console *console

func initConsole() error {

    if Console != nil {
        return fmt.Errorf("console already initialized")
    }

    c := &console{
        messages: []string{},
        element:  termui.NewList(),
    }

    Console = c
    return nil

}

// Start is called at the beginning of the program
// after the ui is initialized
func (c *console) start() {
    elem := c.element
    elem.Overflow = "wrap"
    elem.Items = c.messages
    elem.Height = termui.TermHeight() - Command.element.Height
    elem.BorderFg = termui.ColorWhite

    termui.Handle("/update", c.Update)
    termui.Handle("/sys/wnd/resize", c.Resize)
}

// Element returns the termui element for this console
func (c console) Element() termui.GridBufferer {
    return c.element
}

// Update is used once per "frame"
func (c *console) Update(e termui.Event) {

    //now := e.Data.(int64)

    /*
       validMessages := []string{}
       for _, msg := range c.messages {
           if msg.expires != 0 && msg.expires < now {
               msg.expired = true
           }
       }*/

}

func (c *console) Resize(e termui.Event) {
    w := e.Data.(termui.EvtWnd)
    c.element.Height = w.Height
}

// Clear clears all messages from the console
func (c *console) Clear() {
    c.messages = make([]string, 0)
    c.element.Items = c.messages
}

// Log is used to log a message to the console
func (c *console) Log(msg string) {
    c.messages = append(c.messages, msg)
    c.element.Items = c.messages
    c.element.Y = len(c.messages)
}

// Logf uses fmt.Sprintf to format the given message and vars
func (c *console) Logf(msg string, vars ...interface{}) {
    c.Log(fmt.Sprintf(msg, vars...))
}

func (c *console) Debug(msg string) {
    msg = fmt.Sprintf("[DEBUG: %s](fg-cyan)", msg)
    c.Log(msg)
}

// Debugf uses fmt.Sprintf to format the given message and vars
func (c *console) Debugf(msg string, vars ...interface{}) {
    c.Debug(fmt.Sprintf(msg, vars...))
}

func (c *console) Warning(msg string) {
    msg = fmt.Sprintf("[WARNING: %s](fg-yellow)", msg)
    c.Log(msg)
}

// Warningf uses fmt.Sprintf to format the given message and vars
func (c *console) Warningf(msg string, vars ...interface{}) {
    c.Warning(fmt.Sprintf(msg, vars...))
}

func (c *console) Error(msg string) {
    msg = fmt.Sprintf("[ERROR: %s](fg-red)", msg)
    c.Log(msg)
}

// Errorf uses fmt.Sprintf to format the given message and vars
func (c *console) Errorf(msg string, vars ...interface{}) {
    c.Error(fmt.Sprintf(msg, vars...))
}
