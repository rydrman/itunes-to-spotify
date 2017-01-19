package main

import "github.com/rydrman/termui"

var initialized = false

func initUI() error {

    err := termui.Init()
    if err != nil {
        return err
    }

    err = buildUI()
    if err != nil {
        return err
    }

    initialized = true

    return nil

}

func buildUI() error {

    termui.Body.AddRows(
        termui.NewRow(
            termui.NewCol(12, 0, Console.Element()),
        ),
        termui.NewRow(
            termui.NewCol(12, 0, Command.Element()),
        ),
    )

    return nil

}

// Refresh updates the ui, and should be used whenever something changes
func Refresh() {

    if !initialized {
        return
    }
    termui.Body.Align()
    termui.Render(termui.Body)

}
