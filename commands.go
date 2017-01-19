package main

import (
    "fmt"
    "strings"
)

var quitCommands = []string{
    "q", "quit", "exit",
}

// CmdDef defines a command that can be called in
// Spotr console
type CmdDef struct {
    Handler func([]string) error
    Args    string
    Help    string
}

var commands = map[string]*CmdDef{
    "login": {
        doLogin, "",
        "login to spotify"},
    "import": {
        doImport, "[itunes library file]",
        "begin the playlist import process for the given itunes library"},
}

func showHelp() {
    Console.Clear()
    for name, def := range commands {
        padded := fmt.Sprintf("%-10s", name)
        if len(padded) > 10 {
            Console.Log(padded)
            Console.Logf("           %s", def.Help)
        } else {
            Console.Logf("%s %s", name, def.Help)
        }
    }
}

// RunCommand runs the given command string
func RunCommand(query string) {

    // check quit commands first and foremost
    if StringInSlice(query, quitCommands) {
        Quit()
    }

    parts := strings.Split(query, " ")
    cmd := parts[0]

    def := commands[cmd]
    if nil == def {
        showHelp()
        return
    }
    err := def.Handler(parts[1:])

    if err != nil {
        Console.Error(err.Error())
    }

}

func doLogin(q []string) error {
    return Session.Authenticate()
}

func doImport(q []string) error {
    Console.Logf("%s", q)

    return nil
}
