package main

import "testing"

func TestCompareTitle(t *testing.T) {

	if 1 > TitleCompare("first song", "second song live") {
		t.Error("'live' should automatically make titles unequal")
	}

	if 1 > TitleCompare("first song", "second song karaoke") {
		t.Error("'karaoke' should automatically make titles unequal")
	}

	if TitleCompare("first song", "first song - single version") > simpleEffect {
		t.Error("'single version' should not effect title matching")
	}

}

func TestCompareArtist(t *testing.T) {

	var score float64
	if ArtistCompare("david & goliath", "david and goliath") > cleanEffect {
		t.Error("and usage should not affect artist comparisons")
	}
	if ArtistCompare("david & goliath / me", "david, goliath, and me") > cleanEffect {
		t.Error("listing should not affect artist comparisons")
	}

	if 0 != ArtistCompare("me too & david & goliath", "goliath & david & me too") {
		t.Error("order should not affect artist comparisons")
	}

	if ArtistCompare("The Eagles", "eagles") > simpleEffect {
		t.Error("'the' should be an ignorable word")
	}

	score = ArtistCompare("Harry \"Haywire Mac\" McClintock", "Harry McClintock")
	if score > complexEffect {
		t.Errorf("nick names should not be necessary %f", score)
	}

	score = ArtistCompare("Pink Floyd", "The Pink Boyz")
	if score <= thresholdMatched {
		t.Error("'Pink Floyd' should not match 'The Pink Boyz'")
	}

}

func TestCompareAlbum(t *testing.T) {

	var score float64

	if 1 > AlbumCompare("first song", "second song live", false) {
		t.Error("'live' should automatically make albums unequal")
	}

	if 1 > AlbumCompare("first song", "second song karaoke", false) {
		t.Error("'karaoke' should automatically make albums unequal")
	}

	score = AlbumCompare("Hotel California", "Hotel California (Remastered)", false)
	if score > simpleEffect {
		t.Error("remastered and brackets should have no negative effect")
	}

	score = AlbumCompare(
		"Snakes On A Plane: The Album",
		"Snakes On A Plane [Bring It] (1-track DMD) [0.1221]", false)
	if score > complexEffect {
		t.Error("'Snakes On A Plane: The Album' should match 'Snakes On A Plane [Bring It] (1-track DMD) [0.1221]'")
	}

}

func TestSCompareScore(t *testing.T) {

	var score float64

	if 0 != SCompareScore("option", "option") {
		t.Error("expected same string to return score or 0")
	}

	if SCompareScore("name", "name (with extras)") > simpleEffect {
		t.Error("should ignore data in brackets")
	}

	if SCompareScore("name", "name [with extras]") > simpleEffect {
		t.Error("should ignore data in square brackets")
	}

	score = SCompareScore("tenlettera", "abcdfghijx")
	if 1 != score {
		t.Errorf("all different strings should be a 1, not %f", score)
	}

}
