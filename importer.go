package main

import (
	"fmt"
	"sort"
	"strings"

	itunes "github.com/rydrman/go-itunes-library"
	"github.com/zmb3/spotify"
)

// Importer is used to import an itunes library
type Importer struct {

	// import settings
	AddToLibrary      bool
	LibraryAsPlaylist bool
	ImportPlaylists   bool
	GroupPlaylists    bool
	ImportDisabled    bool
	PlaylistGroup     string

	// match settings
	PreferOriginal bool
	GuessMatching  bool

	// match processing
	matchNum   int
	matchTotal int

	// cache
	missingLog []string
	matchCache *MatchCache
	trackCache map[int]*MatchedTrack
	albumCache map[string][]spotify.FullTrack

	// runtime
	lib     *itunes.Library
	program *SimpleCommandProgram
}

// NewImporter creates a new importer for the command program and
// itunes library that are supplied
func NewImporter(program *SimpleCommandProgram, lib *itunes.Library) *Importer {

	i := &Importer{
		AddToLibrary:    program.AskYesNo("Add all tracks to your library?", false),
		ImportPlaylists: program.AskYesNo("Import playlists?", true),
		GroupPlaylists:  false, //program.AskYesNo("Group all itunes playlists?", true),
		PlaylistGroup:   "iTunes Playlists",
		PreferOriginal:  program.AskYesNo("Prefer non-consolidation albums?", true),
		GuessMatching:   program.AskYesNo("Guess when there are mutliple excellent matches?", true),
		ImportDisabled:  program.AskYesNo("Import unchecked songs?", false),

		matchTotal: len(lib.Tracks),

		matchCache: InitMatchCache(lib.LibraryFile),
		trackCache: make(map[int]*MatchedTrack),
		albumCache: make(map[string][]spotify.FullTrack),
		lib:        lib,
		program:    program,
	}

	return i

}

// Run this importer with the current configuration
func (i *Importer) Run() {

	var user *spotify.PrivateUser
	var err error

	i.program.Log("gathering necessary data...")

	user, err = Session.Client().CurrentUser()
	if nil != err {
		i.program.Errorf("Error getting current user: %s", err)
		return
	}

	if i.AddToLibrary {

		i.program.Log("adding tracks to library...")

		var chunk []spotify.ID
		i.matchNum = 0
		i.matchTotal = len(i.lib.Tracks)
		for _, track := range i.lib.Tracks {

			if i.shouldSkipTrack(track) {
				i.program.Logf("            \nskipping:  %s\n",
					ItunesCacheString(track))
				i.matchTotal--
				continue
			}

			mt := i.getMappedTrack(track.TrackID)
			if mt != nil && mt.Valid() {
				chunk = append(chunk, mt.spotify.ID)
			}
			if len(chunk) == 100 {
				for {
					err = Session.Client().AddTracksToLibrary(chunk...)
					if Session.ShouldTryAgain(err) {
						continue
					} else if err != nil {
						i.program.Errorf("Error adding tracks to library: %s", err)
						return
					}
					break
				}
				chunk = make([]spotify.ID, 0)
			}
		}
		if len(chunk) > 0 {
			for {
				err = Session.Client().AddTracksToLibrary(chunk...)
				if Session.ShouldTryAgain(err) {
					continue
				} else if err != nil {
					i.program.Errorf("Error adding tracks to library: %s", err)
					return
				}
				break
			}
		}
	}

	i.program.Log("creating itunes library playlist...")

	var libraryPlaylist *spotify.FullPlaylist
	for {
		libraryPlaylist, err = Session.Client().CreatePlaylistForUser(
			user.ID, "iTunes Library", false)
		if Session.ShouldTryAgain(err) {
			continue
		} else if err != nil {
			i.program.Errorf("Error creating library playlist: %s", err)
			return
		}
		break
	}

	var chunk []spotify.ID
	i.matchNum = 0
	i.matchTotal = len(i.lib.Tracks)
	for _, track := range i.lib.Tracks {

		if i.shouldSkipTrack(track) {
			i.program.Logf("            \nskipping:  %s\n",
				ItunesCacheString(track))
			i.matchTotal--
			continue
		}

		mt := i.getMappedTrack(track.TrackID)
		if mt != nil && mt.Valid() {
			chunk = append(chunk, mt.spotify.ID)
		}
		if len(chunk) == 100 {
			for {
				_, err = Session.Client().AddTracksToPlaylist(
					user.ID, libraryPlaylist.ID, chunk...)
				if Session.ShouldTryAgain(err) {
					continue
				} else if err != nil {
					i.program.Errorf("Error adding tracks to playlist: %s", err)
					return
				}
				break
			}
			chunk = make([]spotify.ID, 0)
		}
	}
	if len(chunk) > 0 {
		for {
			_, err = Session.Client().AddTracksToPlaylist(
				user.ID, libraryPlaylist.ID, chunk...)
			if Session.ShouldTryAgain(err) {
				continue
			} else if err != nil {
				i.program.Errorf("Error adding tracks to playlist: %s", err)
				return
			}
			break
		}
	}

	if i.ImportPlaylists {

		i.program.Log("importing itunes playlists...")

		for _, iList := range i.lib.Playlists {

			if iList.Master ||
				iList.TVShows ||
				iList.Movies {
				continue
			}

			i.program.Logf("importing %s...", iList.Name)

			var sList *spotify.FullPlaylist
			for {
				sList, err = Session.Client().CreatePlaylistForUser(
					user.ID, iList.Name, false)
				if Session.ShouldTryAgain(err) {
					continue
				} else if err != nil {
					i.program.Errorf("Error creating playlist: %s", err)
					return
				}
				break
			}

			var chunk []spotify.ID
			i.matchNum = 0
			i.matchTotal = len(iList.PlaylistItems)
			for _, track := range iList.PlaylistItems {

				if i.shouldSkipTrack(track) {
					i.program.Logf("            \nskipping:  %s\n",
						ItunesCacheString(track))
					i.matchTotal--
					continue
				}

				mt := i.getMappedTrack(track.TrackID)
				if mt != nil && mt.Valid() {
					chunk = append(chunk, mt.spotify.ID)
				}
				if len(chunk) == 100 {
					for {
						_, err = Session.Client().AddTracksToPlaylist(
							user.ID, sList.ID, chunk...)
						if Session.ShouldTryAgain(err) {
							continue
						} else if err != nil {
							i.program.Errorf("Error adding tracks to playlist: %s", err)
							return
						}
						break
					}
					chunk = make([]spotify.ID, 0)
				}
			}
			if len(chunk) > 0 {
				for {
					_, err = Session.Client().AddTracksToPlaylist(
						user.ID, sList.ID, chunk...)
					if Session.ShouldTryAgain(err) {
						continue
					} else if err != nil {
						i.program.Errorf("Error adding tracks to playlist: %s", err)
						return
					}
					break
				}
			}

		}

	}

}

func (i *Importer) shouldSkipTrack(track *itunes.Track) bool {

	if track.Podcast || track.Movie || track.ITunesU || track.TVShow {
		return true
	}

	if track.Disabled && !i.ImportDisabled {
		return true
	}

	return false

}

func (i *Importer) cacheTrack(mt *MatchedTrack) *MatchedTrack {

	i.program.Logf("  @%1.4f  %s", mt.score, SpotifyCacheString(mt.spotify))

	i.trackCache[mt.itunes.TrackID] = mt
	i.matchCache.TrackMap.Store(mt)

	if mt.itunes.Album != "" && mt.Valid() {
		i.albumCache[mt.itunes.Album] = albumTracks(mt.FullAlbum())
		i.matchCache.AlbumMap.Store(mt)
	}

	i.matchCache.SaveCache()

	return mt

}

func (i *Importer) getMappedTrack(itunesTrackID int) *MatchedTrack {

	i.matchNum++

	// see if it has been cached in this session
	if t, ok := i.trackCache[itunesTrackID]; ok {
		return t
	}

	goal := i.lib.TracksByID[itunesTrackID]

	i.program.Logf("            \n%04d/%04d: %s\n",
		i.matchNum, i.matchTotal, ItunesCacheString(goal))

	// see if it exists in a previous cache
	cached := i.matchCache.TrackMap.GetMatch(goal)
	if nil != cached {
		return i.cacheTrack(cached)
	}

	// TODO special case
	if strings.ToLower(goal.Artist) == "taylor swift" {
		return nil
	}

	goal = PreprocessTrackArtists(goal)

	// see if the album was already mapped
	aTracks, ok := i.albumCache[goal.Album]
	if !ok {
		// see if it exists in a previous cache
		a := i.matchCache.AlbumMap.GetMatch(goal.Album)
		if a != nil {
			aTracks = albumTracks(a)
		}
	}

	var scored []*MatchedTrack

	// use the album if available to look for this track
	if nil != aTracks && len(aTracks) > 0 {
		scored := i.scoreTracks(aTracks, goal)

		sort.Sort(byScoreAndDate(scored))

		if scored[0].score <= thresholdMatched {
			return i.cacheTrack(scored[0])
		}

	}

	// move on to querying spotify
	queryOptions := SearchAttempts(goal)

	for _, query := range queryOptions {

		results := Session.SearchTracks(query, 2) // fetch 2 pages - max 40 tracks
		//fmt.Printf("%s %d\n", query, len(results))

		if 0 == len(results) {
			continue
		}

		var newTracks []spotify.FullTrack
		for _, res := range results {
			found := false
			for _, sd := range scored {
				if SpotifyCacheString(sd.spotify) == SpotifyCacheString(&res) {
					found = true
					break
				}
			}
			if !found {
				newTracks = append(newTracks, res)
			}
		}

		if 0 == len(newTracks) {
			continue
		}

		scored = append(scored, i.scoreTracks(newTracks, goal)...)

		sort.Sort(byScoreAndDate(scored))

		var matched []*MatchedTrack
		for j := 0; j < len(scored) && scored[j].score <= thresholdMatched; j++ {
			matched = append(matched, scored[j])
		}

		if len(matched) > 0 {

			if len(matched) == 1 || i.GuessMatching {
				return i.cacheTrack(matched[0])
			}

		}

	}

	// try re-scoring them without the album name
	for _, mt := range scored {
		mt.score = TrackCompare(mt.itunes, mt.spotify, i.PreferOriginal, true)
	}

	sort.Sort(byScoreAndDate(scored))

	var matched []*MatchedTrack
	for j := 0; j < len(scored) && scored[j].score <= thresholdMatched; j++ {
		matched = append(matched, scored[j])
	}

	if len(matched) > 0 {

		if len(matched) == 1 || i.GuessMatching {
			return i.cacheTrack(matched[0])
		}

	}

	match := i.askMappedTrackSelection(goal, scored)
	if match == nil {
		match = &MatchedTrack{
			itunes:  goal,
			spotify: nil,
			score:   -1,
		}
	}
	return i.cacheTrack(match)

}

func (i *Importer) scoreTracks(tracks []spotify.FullTrack, goal *itunes.Track) []*MatchedTrack {

	mapped := make([]*MatchedTrack, len(tracks))

	for j := 0; j < len(tracks); j++ {

		vers := float64(j) * 0.025

		mapped[j] = &MatchedTrack{
			itunes:  goal,
			spotify: &tracks[j],
			sAlbum:  nil,
			score:   vers + TrackCompare(goal, &tracks[j], i.PreferOriginal, false),
		}

	}

	return mapped

}

func (i *Importer) askMappedTrackSelection(goal *itunes.Track, tracks []*MatchedTrack) *MatchedTrack {

	var sel int
	var text string
	if len(tracks) == 0 {

		fmt.Printf("no matches found for: %s (%s) %s:\n",
			goal.Name, goal.Artist, goal.Album)
		fmt.Printf("try a custom search or press enter to skip:")
		sel = -1
		text = i.program.CaptureInput()
		if text == "" {
			return nil
		}

	} else {

		var options []string
		for _, mt := range tracks {
			options = append(options, fmt.Sprintf("%s [%1.4f]",
				SpotifyCacheString(mt.spotify),
				mt.score))
		}
		options = append(options, "none of the above")

		sel, text = i.program.AskOptionCustom(
			fmt.Sprintf("which of the following best matches: %s (%s) %s?",
				goal.Name, goal.Artist, goal.Album),
			options,
			fmt.Sprintf("custom search (nomatch=%d)", len(options)-1))

	}

	if sel == -1 {

		// search for the term and ask again
		var results *spotify.SearchResult
		var err error
		for {
			results, err = Session.Client().Search(text, spotify.SearchTypeTrack)
			if Session.ShouldTryAgain(err) {
				continue
			} else if err != nil {
				return nil
			}
			break
		}
		var options []*MatchedTrack
		for j := 0; j < len(results.Tracks.Tracks); j++ {
			options = append(options, &MatchedTrack{
				itunes:  goal,
				spotify: &results.Tracks.Tracks[j],
			})
		}
		return i.askMappedTrackSelection(goal, options)
	}

	if sel == len(tracks) {
		return nil
	}

	return tracks[sel]

}
