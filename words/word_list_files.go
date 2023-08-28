package words

import (
	_ "embed"
	"strings"
)

//go:embed data/wordlist-valid-guesses.csv
var guessesFile string

//go:embed data/wordlist-valid-solutions.csv
var solutionsFile string

func GetValidGuessesWordList() ([]Word, error) {
	return makeWordListFromString(guessesFile)
}

func GetValidSolutionsWordList() ([]Word, error) {
	return makeWordListFromString(solutionsFile)
}

func makeWordListFromString(s string) ([]Word, error) {
	lines := strings.Split(s, "\n")
	wl := make([]Word, 0, len(lines))
	for _, line := range lines {
		w, err := NewWord(line)
		if (err != nil) {
			return []Word{}, err
		}
		wl = append(wl, w)
	}
	return wl, nil
}
