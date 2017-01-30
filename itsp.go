package main

import (
    "fmt"
    "os"
    "time"

    "github.com/rydrman/termui"
)

var quit = false

func main() {

    termui.Handle("/sys", func(event termui.Event) {
        switch e := event.Data.(type) {
        case termui.EvtKbd:
            Command.HandleKeyboard(e)
        }
    })

    err := initConsole()
    if nil != err {
        fmt.Print(err.Error())
        quit = true
    }
    err = initCommand()
    if nil != err {
        fmt.Print(err.Error())
        quit = true
    }
    err = initSession()
    if nil != err {
        fmt.Print(err.Error())
        quit = true
    }

    err = initUI()
    if nil != err {
        fmt.Print(err.Error())
        quit = true
    }
    defer termui.Close()

    Console.start()
    Command.start()
    Session.start()

    Console.Log("Welcome to the iTunes to Spotify utility!")
    Console.Log("use 'login' or 'help' to get started :)")
    if "" == clientID || "" == clientSecret {
        Console.Error("app identifiers not found (clientID, clientSecret)")
    }

    // enter a forever update loop for the application
    for quit == false {

        start := time.Now().UnixNano()
        end := start

        termui.SendCustomEvt("/update", start)

        for time.Duration(end-start) < 160*time.Millisecond {
            termui.HandleEvents()
            end = time.Now().UnixNano()
        }

        Refresh()

    }

    Command.SaveHistory()
    os.Exit(0)

}

// Quit stops all processes and exits Spotr
func Quit() {
    quit = true
    termui.Close()
}
