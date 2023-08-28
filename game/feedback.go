package game

import (
	"wordle/words"
)

type FeedbackColour int

const (
	Green FeedbackColour = iota
	Yellow
	Grey
)

type Feedback struct {
	FeedbackColours [5]FeedbackColour
}

func GetFeedback(solution words.Word, guess words.Word) Feedback {
	feedbackColours := [5]FeedbackColour{}

	for i := 0; i < 5; i++ {
		feedbackColours[i] = getFeedbackColour(solution, guess, i)
	}
	return Feedback{FeedbackColours: feedbackColours}
}

func getFeedbackColour(solution words.Word, guess words.Word, index int) FeedbackColour {
	if solution.String()[index] == guess.String()[index] {
		return Green
	}

	for j := 0; j < len(solution.String()); j++ {
		if solution.String()[j] == guess.String()[index] && j != index{
			return Yellow
		}
	}

	return Grey
}

func (f *Feedback) String() string {
	feedbackString := ""
	for _, colour := range f.FeedbackColours {
		switch colour {
		case Grey:
			feedbackString += "-"
		case Yellow:
			feedbackString += "Y"
		case Green:
			feedbackString += "G"
		}
	}
	return feedbackString
}

func (f *Feedback) Equals(another Feedback) bool {
	return f.String() == another.String()
}
