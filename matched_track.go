package main

import (
	itunes "github.com/rydrman/go-itunes-library"
	"github.com/zmb3/spotify"
)

// MatchedTrack represents a itunes / spotify track
// pair that has been successfully matched
type MatchedTrack struct {
	itunes  *itunes.Track
	spotify *spotify.FullTrack
	sAlbum  *spotify.FullAlbum

	score float64
}

// FullAlbum returns the full album for this matches spotify track,
// fetching it if necessary from the server
func (mt *MatchedTrack) FullAlbum() *spotify.FullAlbum {
	if !mt.Valid() {
		return nil
	}
	if nil == mt.sAlbum {
		var err error
		for {
			mt.sAlbum, err = Session.Client().GetAlbum(mt.spotify.Album.ID)
			if Session.ShouldTryAgain(err) {
				continue
			}
			break
		}
	}
	return mt.sAlbum
}

// Valid returns true if this is a complete and valid match
func (mt *MatchedTrack) Valid() bool {
	return nil != mt.spotify
}

type byScoreAndDate []*MatchedTrack

func (s byScoreAndDate) Len() int      { return len(s) }
func (s byScoreAndDate) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byScoreAndDate) Less(i, j int) bool {

	if !s[i].Valid() && s[j].Valid() {
		return false
	}
	if s[i].Valid() && !s[j].Valid() {
		return false
	}

	return s[i].score < s[j].score

}
