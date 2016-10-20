package main

import "github.com/gizak/termui"

// these variables are set during the build process
// using the specified envrionment varaibles
var clientID string     //SPOTIFY_CLIENT_ID
var clientSecret string //SPOTIFY_CLIENT_SECRET

func main() {

    termui.Handle("/sys", func(event termui.Event) {
        switch e := event.Data.(type) {
        case termui.EvtKbd:
            Command.HandleKeyboard(e)
        }
    })

    launchUI()

}
