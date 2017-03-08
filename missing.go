package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	itunes "github.com/rydrman/go-itunes-library"
)

// MissingLog is a log to hold missing entries and where they belong
type MissingLog struct {
	LogFile string
	Entries map[string][]string
}

// InitMissingLog starts a new missing log for outputting at the end of the session
func InitMissingLog(itunesLibraryPath string) *MissingLog {

	ext := path.Ext(itunesLibraryPath)
	baseName := itunesLibraryPath[0 : len(itunesLibraryPath)-len(ext)]
	logFile := fmt.Sprintf("%s.itsp.missing", baseName)

	log := &MissingLog{
		LogFile: logFile,
		Entries: make(map[string][]string),
	}

	return log

}

// SaveLog saves this log to the file system based on the library that
// it was initialized for (overwriting existing log is it exists)
func (ml *MissingLog) SaveLog() error {

	jsonData, err := json.Marshal(ml.Entries)
	if nil != err {
		fmt.Printf("error marshalling log: %s", err)
		return err
	}

	return ioutil.WriteFile(ml.LogFile, jsonData, os.ModeExclusive)

}

// Log logs the given track in this map
func (ml *MissingLog) Log(destination string, track *itunes.Track) {

	if _, ok := ml.Entries[destination]; !ok {
		ml.Entries[destination] = make([]string, 0)
	}

	ml.Entries[destination] = append(ml.Entries[destination], ItunesCacheString(track))

}
