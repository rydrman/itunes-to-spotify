package main

import (
	"os"
	"path/filepath"
	"time"

	itunes "github.com/rydrman/go-itunes-library"
)

func main() {

	program := &SimpleCommandProgram{}
	var err error

	program.Log("Welcome to the iTunes to Spotify utility!")
	program.Log("at any time you can exit by using ctrl+c")

	if "" == clientID || "" == clientSecret {
		program.Error("app identifiers not found (clientID, clientSecret)")
	}

	err = initSession()
	if nil != err {
		program.Error(err.Error())
		os.Exit(1)
	}
	Session.start()

	////////////
	// authenticate with spotify
	////////////

	program.Log("you will need to login to get started")
	program.Log("to open the login page, press enter:")

	_ = program.CaptureInput()

	err = Session.Authenticate()
	if nil != err {
		program.Error(err.Error())
	}

	program.Log("waiting for login response...")
	for Session.IsAuthenticated() == false {
		time.Sleep(time.Millisecond * 250)
	}

	name := "<UNKNOWN>"
	usr, err := Session.Client().CurrentUser()
	if nil != err {
		program.Warningf("error getting user information: %s", err)
	} else {
		name = usr.DisplayName
	}

	program.Logf("Login Successful! Welcome, %s", name)
	program.Log("")

	////////////
	// read the itunes library
	////////////
	var lib *itunes.Library
	for {
		fileName := program.AskStringDefault(
			"enter path to itunes library XML file", "")
		fileName = filepath.Clean(fileName)
		lib, err = itunes.ParseFile(fileName)
		if nil == err {
			break
		}
		program.Error(err.Error())
	}

	program.Log("Library file read successfully!")
	program.Log(lib.String())

	if !Session.IsAuthenticated() {
		program.Warning("You are not logged Spotify, import cannot continue")
		os.Exit(1)
	}

	////////////
	// hand off to importer
	////////////
	importer := NewImporter(program, lib)
	importer.Run()

	Session.Logout()
	os.Exit(0)

}
