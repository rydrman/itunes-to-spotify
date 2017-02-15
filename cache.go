package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	itunes "github.com/rydrman/go-itunes-library"
	"github.com/zmb3/spotify"
)

// MatchCache is a chache to hold previously mapped track and album data
type MatchCache struct {
	LibraryFile string
	CacheFile   string
	TrackMap    TrackMap
	AlbumMap    AlbumMap
}

// InitMatchCache attemps to load the cache for the given itunes library file
// but will return an empty cache if not found
func InitMatchCache(itunesLibraryPath string) *MatchCache {

	ext := path.Ext(itunesLibraryPath)
	baseName := itunesLibraryPath[0 : len(itunesLibraryPath)-len(ext)]
	cacheFile := fmt.Sprintf("%s.itsp.cache", baseName)

	cache := &MatchCache{
		LibraryFile: itunesLibraryPath,
		CacheFile:   cacheFile,
		TrackMap:    make(TrackMap),
		AlbumMap:    make(AlbumMap),
	}

	if _, err := os.Open(cacheFile); os.IsNotExist(err) {
		fmt.Printf("no cache found for: %s\n", cache.LibraryFile)
		return cache
	}

	jsonData, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		fmt.Printf("error reading cache: %s\n", cache.CacheFile)
		return cache
	}
	err = json.Unmarshal(jsonData, cache)
	if err != nil {
		fmt.Printf("error unmarshalling cache: %s\n", cache.CacheFile)
	}

	fmt.Printf("cache file found: %d tracks, %d albums\n", len(cache.TrackMap), len(cache.AlbumMap))

	return cache

}

// SaveCache saves this cache to the file system based on the library that
// it was initialized for (overwriting existing cache is it exists)
func (mc *MatchCache) SaveCache() error {

	jsonData, err := json.Marshal(mc)
	if nil != err {
		fmt.Printf("error marshalling cache: %s", err)
		return err
	}

	return ioutil.WriteFile(mc.CacheFile, jsonData, os.ModeExclusive)

}

// CachedTrackMatch represents a cached MappedTrack, which
// has enough info to be useful for debugging and as a cache
// but not runtime
type CachedTrackMatch struct {
	ItunesTrack  string
	SpotifyTrack string
	SpotifyID    string
	ItunesID     int
	Score        float64
}

// TrackMap stores simple id mapping data for itunes to spotify mappings
type TrackMap map[int]CachedTrackMatch

// GetMatch tried to build a mapped track instance for the given
// itunes track if it is available in this map
func (tm *TrackMap) GetMatch(goal *itunes.Track) *MatchedTrack {

	if cached, ok := (*tm)[goal.TrackID]; ok {

		mt := &MatchedTrack{
			itunes: goal,
			score:  cached.Score,
		}

		if cached.SpotifyID == "" {
			return mt
		}

		for {
			res, err := Session.Client().GetTrack(spotify.ID(cached.SpotifyID))
			if Session.ShouldTryAgain(err) {
				continue
			}
			if err != nil {
				fmt.Printf("error getting spotify track: %s", err)
				return nil
			}
			mt.spotify = res
			return mt
		}

	}

	return nil
}

// Store stores the given mapping in this map
func (tm *TrackMap) Store(mt *MatchedTrack) {

	spotifyID := ""
	if nil != mt.spotify {
		spotifyID = mt.spotify.ID.String()
	}

	(*tm)[mt.itunes.TrackID] = CachedTrackMatch{
		ItunesTrack:  ItunesCacheString(mt.itunes),
		SpotifyTrack: SpotifyCacheString(mt.spotify),
		SpotifyID:    spotifyID,
		ItunesID:     mt.itunes.TrackID,
		Score:        mt.score,
	}

}

// CachedAlbumMatch represents an album match, which
// has enough info to be useful for debugging and as a cache
// but not runtime
type CachedAlbumMatch struct {
	SpotifyAlbum string
	SpotifyID    string
	Score        float64
}

// AlbumMap stores simple id mapping data for itunes to spotify mappings
type AlbumMap map[string]CachedAlbumMatch

// GetMatch tried to build a mapped track instance for the given
// itunes track if it is available in this map
func (am *AlbumMap) GetMatch(name string) *spotify.FullAlbum {

	if cached, ok := (*am)[name]; ok {

		for {
			res, err := Session.Client().GetAlbum(spotify.ID(cached.SpotifyID))
			if Session.ShouldTryAgain(err) {
				continue
			}
			if err != nil {
				fmt.Printf("error getting cahed album: %s\n", err)
				return nil
			}
			return res
		}

	}

	return nil
}

// Store stores the given mapping in this map
func (am *AlbumMap) Store(mt *MatchedTrack) {

	(*am)[mt.itunes.Album] = CachedAlbumMatch{
		SpotifyAlbum: mt.spotify.Album.Name,
		SpotifyID:    mt.spotify.Album.ID.String(),
		Score:        mt.score,
	}

}

// ItunesCacheString prints out a nice string representation of a itunes
// for debugging use in the cache
func ItunesCacheString(track *itunes.Track) string {
	return fmt.Sprintf("%s (%s)[%s]", track.Name, track.Artist, track.Album)
}

// SpotifyCacheString prints out a nice string representation of a spotify
// for debugging use in the cache
func SpotifyCacheString(track *spotify.FullTrack) string {
	if track == nil {
		return "<NO MATCH>"
	}
	return fmt.Sprintf("%s (%s)[%s]", track.Name, artist(track), track.Album.Name)
}
