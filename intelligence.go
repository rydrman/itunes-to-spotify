package main

import (
	"fmt"
	"math"
	re "regexp"
	"sort"
	"strings"

	itunes "github.com/rydrman/go-itunes-library"
	"github.com/xrash/smetrics"
	"github.com/zmb3/spotify"
)

const (
	thresholdMatched float64 = 0.15
	thresholdLikely          = 0.5
	thresholdSimilar         = 1.0

	cleanEffect   = 0.025
	simpleEffect  = 0.05
	complexEffect = 0.95

	titleWeight      = 0.5
	artistWeight     = 0.3
	albumWeight      = 0.1
	popularityWeight = 0.1
)

// cleanReplacements are regexs that attempt
// to clean names so they are more similar,
// replacing common string permutations
var cleanReplacements = map[string][]*re.Regexp{
	"${1}": {
		re.MustCompile(`(^|\s+)&\s+`),
		re.MustCompile(`,(\s+)`),
		re.MustCompile(`(^|\s+)and\s+`),
		re.MustCompile(`(^|\s+)(/|\\)\s+`),
	},
	"${1}in${2}": {
		re.MustCompile(`(\w+)ing(\s|$)`),
		re.MustCompile(`(\w+)in'(\s|$)`),
	},
}

// simpleReplacements are regexs that attempt to
// simplify names so that they don't include common
// extra information that is not relevant
var simpleReplacements = map[string][]*re.Regexp{
	" ": {
		re.MustCompile(`\s+the\s+`),
		re.MustCompile(`\s*(vs\.?|versus)\s*$`),
	},
	"": {
		re.MustCompile(`[!]`),
		re.MustCompile(`(^|\s+)the(\s+|$)`),
		re.MustCompile(`\s*\(.*\)`),
		re.MustCompile(`\s*\[.*\]`),
		re.MustCompile(` - (\w+ )?from .*$`),
		re.MustCompile(` - single version.*$`),
		re.MustCompile(` - radio edit.*$`),
	},
}

// complexReplacements are regexs that attempt to
// coerce names into being similar by removing potentially
// important but maybe not important information
var complexReplacements = map[string][]*re.Regexp{
	"": {
		re.MustCompile(`:\s+.*$`),
		re.MustCompile(`["'].*["']\s*`),
	},
}

// PreprocessTrackArtists attempts to make itunes track names more friendly to
// spotify by removing featured artists and adding them to the main aritst list
func PreprocessTrackArtists(goal *itunes.Track) *itunes.Track {

	newTrack := itunes.Track(*goal)

	featureRe := re.MustCompile(`[\s\(\[](feat\.?|ft\.?|featuring)\s([\s\w,&]*)`)

	for groups := featureRe.FindStringSubmatch(newTrack.Name); len(groups) > 0; groups = featureRe.FindStringSubmatch(newTrack.Name) {

		newTrack.Name = strings.Replace(newTrack.Name, groups[0], "", 1)
		newTrack.Artist += " & " + groups[2]

	}

	for groups := featureRe.FindStringSubmatch(newTrack.Artist); len(groups) > 0; groups = featureRe.FindStringSubmatch(newTrack.Artist) {

		newTrack.Artist = strings.Replace(newTrack.Artist, groups[0], "", 1)
		newTrack.Artist += " & " + groups[2]

	}

	return &newTrack

}

// TrackCompare intelligently compares the itunes track to the spotify track and
// returns a number from 0 to 1 0 being exaclty the same to 1 being totally different
func TrackCompare(goal *itunes.Track, test *spotify.FullTrack, preferOriginal, ignoreAlbum bool) float64 {

	score := 0.0

	// first compare title
	score += titleWeight * TitleCompare(goal.Name, test.Name)

	// then artist
	score += artistWeight * ArtistCompare(goal.Artist, artist(test))

	// then album
	score += albumWeight * AlbumCompare(goal.Album, test.Album.Name, ignoreAlbum)
	if preferOriginal && test.Album.AlbumType == "consolidation" {
		score += 0.1
	}

	// account for popularity
	score += (1.0 - float64(test.Popularity)/100.0) * popularityWeight

	return score

}

// TitleCompare compares two track titles to estimate the likelyhood
// of a match, returns a probability float (can be greater than 1, but that
// means the match is even less likely)
func TitleCompare(a, b string) float64 {

	a = strings.Trim(strings.ToLower(a), " ")
	b = strings.Trim(strings.ToLower(b), " ")

	specialStrings := []*re.Regexp{
		re.MustCompile(`(^|\s+)live(\s+|$)`),
		re.MustCompile(`karaoke`),
		re.MustCompile(`instrumental`),
		re.MustCompile(`cover`),
	}

	for _, regex := range specialStrings {

		if regex.MatchString(a) != regex.MatchString(b) {
			return 1.0 / titleWeight
		}

	}

	return SCompareScore(a, b)

}

// AlbumCompare compares two album titles to estimate the likelyhood
// of a match, returns a probability float (can be greater than 1, but that
// means the match is even less likely)
//
// simpleCompare will foregoe string comparisons in an attempt to
// only look for albums that are not obviously problematic
func AlbumCompare(a, b string, simpleCompare bool) float64 {

	a = strings.Trim(strings.ToLower(a), " ")
	b = strings.Trim(strings.ToLower(b), " ")

	// empty string makes us unsure but not devastatingly
	if a == "" || b == "" {
		return 0.5
	}

	specialStrings := []*re.Regexp{
		re.MustCompile(`(^|\s+)live(\s+|$)`),
		re.MustCompile(`karaoke`),
		re.MustCompile(`(^|\s+)cast(\s+|$)`),
		re.MustCompile(`soundtrack`),
		re.MustCompile(`cover`),
	}

	for _, regex := range specialStrings {

		if regex.MatchString(a) != regex.MatchString(b) {
			return 1.0 / albumWeight
		}

	}

	// these regular expressions denote album titles that
	// should never be allwed to pass through a simple comparison
	// as they are likely very important (is soundtracks)
	noSimpleStrings := []*re.Regexp{
		re.MustCompile(`(^|\s+)orignal (%w+\s)?cast(\s+|$)`),
		re.MustCompile(`(^|\s+)cast(\s+|$)`),
		re.MustCompile(`(^|\s+)broadway cast(\s+|$)`),
		re.MustCompile(`(^|\s+)london cast(\s+|$)`),
		re.MustCompile(`(^|\s+)soundtrack(\s+|$)`),
		re.MustCompile(`(^|\s+)motion picture(\s+|$)`),
		re.MustCompile(`(^|\s+)musical(\s+|$)`),
	}

	for _, regex := range noSimpleStrings {

		if regex.MatchString(a) || regex.MatchString(b) {
			simpleCompare = false
			break
		}

	}

	if simpleCompare {
		return 0.0
	}

	return SCompareScore(a, b)

}

// ArtistCompare compares two track artists to estimate the likelyhood
// of a match, returns a probability float (can be greater than 1, but that
// means the match is even less likely)
func ArtistCompare(a, b string) float64 {

	a = strings.Trim(strings.ToLower(a), " ")
	b = strings.Trim(strings.ToLower(b), " ")

	specialStrings := []*re.Regexp{
		re.MustCompile(`karaoke`),
		re.MustCompile(`(^|\s+)cast(\s+|$)`),
		re.MustCompile(`soundtrack`),
	}

	for _, regex := range specialStrings {

		if regex.MatchString(a) != regex.MatchString(b) {
			return 1.0 / artistWeight
		}

	}

	// sort all strings alphabetically for better chance
	// of matching one another
	partsA := strings.Split(a, " ")
	sort.Strings(partsA)
	partsB := strings.Split(b, " ")
	sort.Strings(partsB)

	a = strings.Join(partsA, " ")
	b = strings.Join(partsB, " ")

	return SCompareScore(a, b)

}

// SCompareScore returns a score to compare the given strings
// based on basic string compare methods and common title deviants
func SCompareScore(a, b string) float64 {

	score := wagnerFischerRelative(a, b, 1, 1, 1)

	if score < cleanEffect {
		return score
	}

	for r, options := range cleanReplacements {

		for _, option := range options {

			a = option.ReplaceAllString(a, r)
			b = option.ReplaceAllString(b, r)

			score = math.Min(
				score,
				cleanEffect+wagnerFischerRelative(a, b, 1, 1, 1),
			)

			if score == cleanEffect {
				return cleanEffect
			}

		}

	}

	if score < simpleEffect {
		return score
	}

	for r, options := range simpleReplacements {

		for _, option := range options {

			a = option.ReplaceAllString(a, r)
			b = option.ReplaceAllString(b, r)

			score = math.Min(
				score,
				simpleEffect+wagnerFischerRelative(a, b, 1, 1, 1),
			)

			if score == simpleEffect {
				return simpleEffect
			}

		}

	}

	if score < complexEffect {
		return score
	}

	for r, options := range complexReplacements {

		for _, option := range options {

			a = option.ReplaceAllString(a, r)
			b = option.ReplaceAllString(b, r)

			score = math.Min(
				score,
				complexEffect+wagnerFischerRelative(a, b, 1, 1, 1),
			)

			if score == complexEffect {
				return complexEffect
			}

		}

	}

	return score

}

func wagnerFischerRelative(a, b string, icost, dcost, scost int) float64 {

	score := float64(smetrics.WagnerFischer(a, b, icost, dcost, scost))
	longest := math.Max(float64(len(a)), float64(len(b)))
	return score / longest

}

// SearchAttempts returns a list of strings
// to try in searching for this track
func SearchAttempts(goal *itunes.Track) []string {

	normName := strings.ToLower(goal.Name)
	normArtist := strings.ToLower(goal.Artist)
	normAlbum := strings.ToLower(goal.Album)

	queries := []string{
		fmt.Sprintf(`"%s"`, normName),
		fmt.Sprintf(`"%s" "%s"`, normName, normArtist),
		fmt.Sprintf(`track:"%s"`, normName),
		fmt.Sprintf(`track:"%s" artist:"%s"`, normName, normArtist),
	}

	cleanName := normName
	cleanArtist := normArtist
	cleanAlbum := normAlbum
	for r, options := range cleanReplacements {

		for _, option := range options {

			cleanName = option.ReplaceAllString(cleanName, r)
			cleanArtist = option.ReplaceAllString(cleanArtist, r)
			cleanAlbum = option.ReplaceAllString(cleanAlbum, r)

		}

	}

	queries = append(
		queries, fmt.Sprintf(`track:"%s"`, cleanName))
	queries = append(
		queries, fmt.Sprintf(`track:"%s" artist:"%s"`, cleanName, cleanArtist))

	simpleName := cleanName
	simpleArtist := cleanArtist
	simpleAlbum := cleanAlbum
	for r, options := range simpleReplacements {

		for _, option := range options {

			simpleName = option.ReplaceAllString(simpleName, r)
			simpleArtist = option.ReplaceAllString(simpleArtist, r)
			simpleAlbum = option.ReplaceAllString(simpleAlbum, r)

		}

	}

	queries = append(
		queries, fmt.Sprintf(`track:"%s"`, simpleName))
	queries = append(
		queries, fmt.Sprintf(`track:"%s" artist:"%s"`, simpleName, simpleArtist))

	queries = append(
		queries, fmt.Sprintf(`"%s"`, cleanName))
	queries = append(
		queries, fmt.Sprintf(`"%s" "%s"`, cleanName, cleanArtist))
	queries = append(
		queries, fmt.Sprintf(`"%s"`, simpleName))
	queries = append(
		queries, fmt.Sprintf(`"%s" "%s"`, simpleName, simpleArtist))

	queries = append(
		queries, fmt.Sprintf(`track:"%s" artist:"%s" album:"%s"`, normName, normArtist, normAlbum))
	queries = append(
		queries, fmt.Sprintf(`track:"%s" artist:"%s" album:"%s"`, cleanName, cleanArtist, cleanAlbum))
	queries = append(
		queries, fmt.Sprintf(`track:"%s" artist:"%s" album:"%s"`, simpleName, simpleArtist, simpleAlbum))

	queries = append(
		queries, fmt.Sprintf(`"%s" "%s" "%s"`, normName, normArtist, normAlbum))
	queries = append(
		queries, fmt.Sprintf(`"%s" "%s" "%s"`, cleanName, cleanArtist, cleanAlbum))
	queries = append(
		queries, fmt.Sprintf(`"%s" "%s" "%s"`, simpleName, simpleArtist, simpleAlbum))

	// sort out duplicates
	var ret []string
	for _, q := range queries {
		if !StringInSlice(q, ret) {
			ret = append(ret, q)
		}
	}

	return ret

}
