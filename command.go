package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"

    "github.com/fatih/color"
)

// SimpleCommandProgram manages a simple interactive command
// program for go processes running in the command line
type SimpleCommandProgram struct{}

// ClearInput clears any current input in the stdin buffer
// so that a new anser can be read
func (p *SimpleCommandProgram) ClearInput() {

    stats, err := os.Stdin.Stat()
    for err != nil && stats != nil && stats.Size() > 0 {
        _ = p.CaptureInput()
        stats, err = os.Stdin.Stat()
    }

}

// CaptureInput gets the next input string from the user
func (p *SimpleCommandProgram) CaptureInput() string {

    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    return text[:len(text)-2]

}

// AskString simply promts the user for a string answer
func (p *SimpleCommandProgram) AskString(question string) string {

    p.ClearInput()
    fmt.Printf("%s:", question)
    text := p.CaptureInput()
    return text

}

// AskStringDefault simply promts the user for a string answer but
// also provides a default value if they simply hit return
func (p *SimpleCommandProgram) AskStringDefault(question, def string) string {

    p.ClearInput()
    fmt.Printf("%s [%s]:", question, def)
    text := p.CaptureInput()
    if text == "" {
        return def
    }
    return text

}

// AskYesNo asks the given question as a yes or no option and waits for
// a valid answer from the user. def is the default answer if the user simply
// hits the return key
func (p *SimpleCommandProgram) AskYesNo(question string, def bool) bool {

    p.ClearInput()
    defStr := "[y/N]"
    if def == true {
        defStr = "[Y/n]"
    }

    for {
        fmt.Printf("%s %s:", question, defStr)
        text := p.CaptureInput()
        switch text {
        case "":
            return def

        case "yes":
            fallthrough
        case "y":
            fallthrough
        case "Yes":
            fallthrough
        case "Y":
            return true

        case "no":
            fallthrough
        case "n":
            fallthrough
        case "No":
            fallthrough
        case "N":
            return false

        default:
            fmt.Println("invalid answer")

        }

    }

}

// AskOption gives the user a selection betweek the given
// options and returns an integer representing their selection
func (p *SimpleCommandProgram) AskOption(question string, options []string) int {

    p.ClearInput()
    fmt.Printf("%s:\n", question)
    for i, option := range options {
        fmt.Printf("[%d] %s\n", i, option)
    }
    fmt.Printf("please select one of the above options:")
    for {
        text := p.CaptureInput()
        if text == "" {
            continue
        }
        sel, err := strconv.Atoi(text)
        if err != nil || sel < 0 || sel >= len(options) {
            fmt.Printf("%s is not a valid option, select again:", text)
            continue
        }
        return sel
    }

}

// AskOptionDefault gives the user a selection betweek the given
// options and returns an integer representing their selection
func (p *SimpleCommandProgram) AskOptionDefault(question string, options []string, def int) int {

    p.ClearInput()
    fmt.Printf("%s:\n", question)
    for i, option := range options {
        if i == def {
            fmt.Printf("[*%d] %s\n", i, option)
        } else {
            fmt.Printf("[%d] %s\n", i, option)
        }
    }
    for {
        fmt.Printf("please select one of the above options [%d]:", def)
        text := p.CaptureInput()
        if text == "" {
            return def
        }
        sel, err := strconv.Atoi(text)
        if err != nil || sel < 0 || sel >= len(options) {
            fmt.Printf("%s is not a valid option, select again:", text)
            continue
        }
        return sel
    }

}

// AskOptionCustom gives the user a selection betweek the given
// options or creating a custom string value returns an integer and
// string representing their selection. custom is a string that should
// fin into the phrase select one option, or enter <custom>
func (p *SimpleCommandProgram) AskOptionCustom(question string, options []string, custom string) (int, string) {

    p.ClearInput()
    fmt.Printf("%s:\n", question)
    cursor := 0
    cursorLimit := 15
    for {
        for cursor < cursorLimit && cursor < len(options) {
            fmt.Printf("[%d] %s\n", cursor, options[cursor])
            cursor++
        }
        if cursor >= len(options)-1 {
            fmt.Printf("select one option, or enter %s:", custom)
        } else {
            fmt.Printf("select one option, enter %s, or enter for next page:", custom)
        }
        text := p.CaptureInput()
        if text == "" {
            cursorLimit += 15
            continue
        }
        sel, err := strconv.Atoi(text)
        if err == nil && (sel < 0 || sel >= len(options)) {
            fmt.Printf("%s is not a valid option, select again:", text)
            continue
        }
        if err != nil {
            return -1, text
        }
        return sel, options[sel]
    }

}

// Log is used to log a message to the console
func (p *SimpleCommandProgram) Log(msg string) {
    color.White(msg)
}

// Logf uses fmt.Sprintf to format the given message and vars
func (p *SimpleCommandProgram) Logf(msg string, vars ...interface{}) {
    p.Log(fmt.Sprintf(msg, vars...))
}

// Debug prints a debug message to the program
func (p *SimpleCommandProgram) Debug(msg string) {
    msg = fmt.Sprintf("DEBUG: %s", msg)
    color.Cyan(msg)
}

// Debugf uses fmt.Sprintf to format the given message and vars
func (p *SimpleCommandProgram) Debugf(msg string, vars ...interface{}) {
    p.Debug(fmt.Sprintf(msg, vars...))
}

// Warning prints a warning to this program
func (p *SimpleCommandProgram) Warning(msg string) {
    msg = fmt.Sprintf("WARNING: %s", msg)
    color.Yellow(msg)
}

// Warningf uses fmt.Sprintf to format the given message and vars
func (p *SimpleCommandProgram) Warningf(msg string, vars ...interface{}) {
    p.Warning(fmt.Sprintf(msg, vars...))
}

func (p *SimpleCommandProgram) Error(msg string) {
    msg = fmt.Sprintf("ERROR: %s", msg)
    color.Red(msg)
}

// Errorf uses fmt.Sprintf to format the given message and vars
func (p *SimpleCommandProgram) Errorf(msg string, vars ...interface{}) {
    p.Error(fmt.Sprintf(msg, vars...))
}
