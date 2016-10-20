package main

import "github.com/gizak/termui"

func launchUI() error {

    err := termui.Init()
    if err != nil {
        return err
    }
    defer termui.Close()

    err = buildUI()
    if err != nil {
        return err
    }

    Refresh()

    termui.Loop()

    return nil

}

func buildUI() error {

    cmdBox, err := setupCommandBox()
    if err != nil {
        return err
    }

    termui.Body.AddRows(
        termui.NewRow(
            termui.NewCol(12, 0, cmdBox),
        ),
    )

    return nil

}

func setupCommandBox() (termui.GridBufferer, error) {

    cmdBox := termui.NewPar("")
    cmdBox.Height = 3
    cmdBox.BorderFg = termui.ColorWhite

    termui.Handle("/sys/kbd", func(event termui.Event) {

        e := event.Data.(termui.EvtKbd)

        switch e.KeyStr {
        case "<space>":
            cmdBox.Text += " "
        case "<enter>":
            RunCommand(cmdBox.Text)
            cmdBox.Text = ""
        default:
            cmdBox.Text += e.KeyStr
        }

        Refresh()

    })

    return cmdBox, nil
}

// Refresh updates the ui, and should be used whenever something changes
func Refresh() {
    termui.Body.Align()
    termui.Render(termui.Body)
}
