package words

type Word struct {
	Characters []rune
	String     string
}

func MakeWord(word string) Word {
	return Word{Characters: []rune(word), String: word}
}

func (w Word) Equals(another Word) bool {
	return w.String == another.String
}
