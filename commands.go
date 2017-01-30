package main

import (
    "fmt"
    "strings"

    "github.com/rydrman/go-itunes-library"
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
    "logout": {
        doLogout, "",
        "logout of spotify"},
    "import": {
        doImport, "[itunes library file]",
        "begin the playlist import process for the given itunes library"},
}

func showHelp() {
    Console.Clear()
    for name, def := range commands {
        fullName := fmt.Sprintf("%s %s", name, def.Args)
        padded := fmt.Sprintf("%-15s", fullName)
        if len(padded) > 15 {
            Console.Log(padded)
            Console.Logf("                %s", def.Help)
        } else {
            Console.Logf("%s %s", padded, def.Help)
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
        Console.Errorf(
            "'%s' command not found, use 'help' to see available commands", cmd)
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

func doLogout(q []string) error {
    return Session.Logout()
}

func doImport(q []string) error {

    lib, err := itunes.ParseFile(q[0])
    if nil != err {
        return err
    }

    Console.Log("Library file opened successfully!")
    Console.Log(lib.String())

    if !Session.IsAuthenticated() {
        Console.Warning("You are not logged Spotify, import cannot continue")
        return nil
    }

    return nil
}
