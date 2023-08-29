package words

import (
	"testing"
)

func TestValidGuessesReadCorrectly(t *testing.T) {
	wl, err := GetValidGuesses()
	wordListShouldBeAsExpected(t, wl, 12947, "AAHED", err)
}

func TestValidSolutionsReadCorrectly(t *testing.T) {
	wl, err := GetValidSolutions()
	wordListShouldBeAsExpected(t, wl, 2309, "ABACK", err)
}

func wordListShouldBeAsExpected(t *testing.T, wl []Word, expectedLength int, expectedFirst string, err error) {
	gotLength := len(wl)
	gotFirst := wl[0].String()

	if err != nil {
		t.Errorf("Error reading word list: %s", err)
	}

	if gotLength != expectedLength {
		t.Errorf("Expected %d words, got %d", expectedLength, gotLength)
	}

	if gotFirst != expectedFirst {
		t.Errorf("Expected %q, got %q", expectedFirst, gotFirst)
	}
}
