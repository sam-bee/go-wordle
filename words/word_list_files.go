package words

import (
	_ "embed"
	"strings"
)

//go:embed data/wordlist-valid-guesses.csv
var guessesFile string

//go:embed data/wordlist-valid-solutions.csv
var solutionsFile string

func GetValidGuessesWordList() (WordList) {
	return makeWordListFromString(guessesFile)
}

func GetValidSolutionsWordList() (WordList) {
	return makeWordListFromString(solutionsFile)
}

func makeWordListFromString(s string) WordList {
	lines := strings.Split(s, "\n")
	words := make([]Word, 0)
	for _, line := range lines {
		words = append(words, MakeWord(line))
	}
	return WordList{Words: words}
}
