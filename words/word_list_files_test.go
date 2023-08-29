package words

import (
	"testing"
)

func TestValidGuessesReadCorrectly(t *testing.T) {
	wl, _ := GetValidGuesses()
	testWordList(t, wl, 12947, "AAHED")
}

func TestValidSolutionsReadCorrectly(t *testing.T) {
	wl, _ := GetValidSolutions()
	testWordList(t, wl, 2309, "ABACK")
}

func testWordList(t *testing.T, wl []Word, expectedLength int, expectedFirst string) {
	gotLength := len(wl)
	gotFirst := wl[0].String()

	if gotLength != expectedLength {
		t.Errorf("Expected %d words, got %d", expectedLength, gotLength)
	}

	if gotFirst != expectedFirst {
		t.Errorf("Expected %q, got %q", expectedFirst, gotFirst)
	}
}
