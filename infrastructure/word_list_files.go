package infrastructure

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"wordle/domain/words"
)

const wordListValidGuessesFile = "./data/wordlist-valid-guesses.csv"
const wordListValidSolutionsFile = "./data/wordlist-valid-solutions.csv"

func GetValidGuessesWordList(writer io.Writer) (words.WordList, error) {
	fmt.Fprintf(writer, "Reading from file: %s\n", wordListValidGuessesFile)
	return makeWordListFromFile(wordListValidGuessesFile)
}

func GetValidSolutionsWordList(writer io.Writer) (words.WordList, error) {
	fmt.Fprintf(writer, "Reading from file: %s\n", wordListValidSolutionsFile)
	return makeWordListFromFile(wordListValidSolutionsFile)
}

func makeWordListFromFile(filename string) (words.WordList, error) {

	// To prepare to read the contents, open the file
	file, err := os.Open(filename)
	if err != nil {
		return words.WordList{}, err
	}
	defer file.Close()

	// To provide a list of Words in the file, scan its contents
	wordsFromFile := make([]words.Word, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		word := words.MakeWord(text)
		wordsFromFile = append(wordsFromFile, word)
	}
	return words.WordList{Words: wordsFromFile}, nil
}
