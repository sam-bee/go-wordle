package words

type WordList struct {
	Words []Word
}

func (wl WordList) Count() int {
	return len(wl.Words)
}

func MakeWordList(words []Word) WordList {
	return WordList{Words: words}
}

func Count(wl WordList) int {
	return len(wl.Words)
}
