package words

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const wordListValidGuessesFile = "./words/data/wordlist-valid-guesses.csv"
const wordListValidSolutionsFile = "./words/data/wordlist-valid-solutions.csv"

func GetValidGuessesWordList(writer io.Writer) (WordList, error) {
	fmt.Fprintf(writer, "Reading from file: %s\n", wordListValidGuessesFile)
	return makeWordListFromFile(wordListValidGuessesFile)
}

func GetValidSolutionsWordList(writer io.Writer) (WordList, error) {
	fmt.Fprintf(writer, "Reading from file: %s\n", wordListValidSolutionsFile)
	return makeWordListFromFile(wordListValidSolutionsFile)
}

func makeWordListFromFile(filename string) (WordList, error) {

	// To prepare to read the contents, open the file
	file, err := os.Open(filename)
	if err != nil {
		return WordList{}, err
	}
	defer file.Close()

	// To provide a list of Words in the file, scan its contents
	wordsFromFile := make([]Word, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		word := MakeWord(text)
		wordsFromFile = append(wordsFromFile, word)
	}
	return WordList{Words: wordsFromFile}, nil
}
