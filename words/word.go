package words

import (
	"fmt"
)

type Word struct {
	value string
}

func NewWord(value string) (w Word, err error) {
	if len(value) != 5 {
		return Word{}, fmt.Errorf("all words in a wordle game must be characters long, got '%s'", value)
	}

	return Word{value: value}, nil
}

func (w *Word) String() string {
	return w.value
}
