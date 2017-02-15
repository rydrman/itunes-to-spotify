package main

import (
	"math/rand"

	"github.com/zmb3/spotify"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// StringInSlice returns whether a string is contained
// in a slice of strings
func StringInSlice(s string, list []string) bool {
	for _, m := range list {
		if m == s {
			return true
		}
	}
	return false
}

// RandomToken returns a 64 character string of random characters
func RandomToken() string {
	b := make([]byte, 64)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func artist(track *spotify.FullTrack) string {
	artistStr := ""
	for i, artist := range track.Artists {
		if i != 0 {
			artistStr += " & "
		}
		artistStr += artist.Name
	}
	return artistStr
}

func albumTracks(album *spotify.FullAlbum) []spotify.FullTrack {

	for {

		var ids []spotify.ID
		for _, t := range album.Tracks.Tracks {
			ids = append(ids, t.ID)
		}
		tracks, err := Session.Client().GetTracks(ids...)
		if Session.ShouldTryAgain(err) {
			continue
		}
		ret := make([]spotify.FullTrack, len(tracks))
		for i := 0; i < len(tracks); i++ {
			ret[i] = *tracks[i]
		}
		return ret

	}

}
