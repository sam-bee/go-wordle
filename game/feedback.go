package game

import (
	"wordle/words"
)

const (
	green int = iota
	yellow
	grey
)

type Feedback struct {
	colours []int
}

func GetFeedback(solution words.Word, guess words.Word) Feedback {
	colours := []int{}
	for i := range solution.String() {
		colours = append(colours, getFeedbackColour(solution, guess, i))
	}
	return Feedback{colours: colours}
}

func getFeedbackColour(solution words.Word, guess words.Word, index int) int {
	if solution.String()[index] == guess.String()[index] {
		return green
	}

	for j := 0; j < len(solution.String()); j++ {
		if solution.String()[j] == guess.String()[index] && j != index{
			return yellow
		}
	}

	return grey
}

func (f *Feedback) String() string {
	feedbackString := ""
	for _, colour := range f.colours {
		switch colour {
			case grey:
				feedbackString += "-"
			case yellow:
				feedbackString += "Y"
			case green:
				feedbackString += "G"
		}
	}
	return feedbackString
}

func (f *Feedback) Equals(another Feedback) bool {
	return f.String() == another.String()
}
