package words

import (
	"fmt"
)

type Word string

func NewWord(val string) (Word, error) {
	if len(val) != 5 {
		return "", fmt.Errorf("all words in a wordle game must be characters long, got '%s'", val)
	}
	var w = Word(val)

	return w, nil
}

func (w *Word) Equals(another Word) bool {
	return *w == another
}
