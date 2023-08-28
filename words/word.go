package words

import (
	"fmt"
)

type Word struct {
	value string
}

func NewWord(val string) (w Word, err error) {
	if len(val) != 5 {
		return Word{}, fmt.Errorf("all words in a wordle game must be characters long, got '%s'", val)
	}

	return Word{value: val}, nil
}

func (w *Word) String() string {
	return w.value
}
