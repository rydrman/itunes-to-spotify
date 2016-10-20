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
    termui.Body.Align()
    termui.Render(termui.Body)
}
