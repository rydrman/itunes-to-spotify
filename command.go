package main

import (
    "bufio"
    "fmt"
    "os"
    "path"

    "github.com/rydrman/termui"
)

type command struct {
    element  *termui.Par
    focused  bool
    prompt   string
    value    string
    callback func(string, error)

    history         []string
    historyLocation int

    Password bool
}

// Command is a singleton instance class for managing
// the Spotr command input
var Command *command

// newCommand is not exported because we are forcing a singleton
func initCommand() error {

    c := &command{
        prompt:          ">",
        focused:         true,
        value:           "",
        historyLocation: 0,
    }

    cmd := termui.NewPar("")
    cmd.Height = 3
    cmd.BorderFg = termui.ColorWhite

    c.element = cmd
    Command = c

    Command.LoadHistory()

    return nil

}

// start is run one at the beginning of the program
// after the ui is initialized
func (c *command) start() {
}

// Element returns the termui element for this command box
func (c command) Element() termui.GridBufferer {
    c.Paint()
    return c.element
}

// Clear clears the current value of the command input
func (c *command) Clear() {

    c.value = ""
    c.historyLocation = len(c.history)
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

func (c *command) LoadHistory() error {

    tempDir := os.TempDir()
    historyFile := path.Join(tempDir, ".itspHistory")

    if _, err := os.Stat(historyFile); os.IsNotExist(err) {
        Console.Debug("no history file to load")
        return nil
    }

    file, err := os.Open(historyFile)
    if nil != err {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    for scanner.Scan() {
        c.history = append(c.history, scanner.Text())
    }

    Console.Debugf("loaded %d history commands", len(c.history))

    c.Clear()

    return nil

}

func (c *command) SaveHistory() {

    tempDir := os.TempDir()
    historyFile := path.Join(tempDir, ".itspHistory")

    file, err := os.Create(historyFile)
    if nil != err {
        Console.Error(err.Error())
        return
    }
    defer file.Close()

    for _, cmd := range c.history {
        file.WriteString(fmt.Sprintf("%s\n", cmd))
    }

}

// Paint re-paints the command ui box with any new changes
func (c *command) Paint() {

    promptColor := "fg-cyan"
    if !c.focused {
        promptColor = "fg-white"
    }

    valueRunes := []rune(c.value)
    if c.Password == true {
        for i := range valueRunes {
            valueRunes[i] = '•'
        }
    }

    c.element.Text = fmt.Sprintf("[%s](%s)%s", c.prompt, promptColor, string(valueRunes))
    Refresh()
}

func (c *command) ShowHistory(record int) {

    if record < 0 {

        return

    }

    if record >= len(c.history) {

        record = len(c.history)
        c.value = ""

    } else {

        c.value = c.history[record]

    }

    c.historyLocation = record
    c.Paint()

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

            if c.value == "" {
                return
            }

            RunCommand(c.value)
            c.history = append(c.history, c.value)

        }
        c.Clear()

    case "<up>":
        c.ShowHistory(c.historyLocation - 1)

    case "<down>":
        c.ShowHistory(c.historyLocation + 1)

    case ":":
        if len(c.value) == 0 {
            return
        }
        fallthrough

    default:
        c.value += e.KeyStr
        c.Paint()

    }

}
