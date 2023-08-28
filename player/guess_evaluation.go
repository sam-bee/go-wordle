package player

import (
	"wordle/game"
	"wordle/words"
)

type GuessEvaluation struct {
	Guess                            words.Word
	shortlistSize                    int
	potentialFeedbackCounts          map[string]int
	worstCaseScenarioFeedbackString  string
	worstCaseShortlistCarryOverRatio float64
	isPotentialSolution              bool
}

func NewGuessEvaluation(guess words.Word, currrentShortlist []words.Word) GuessEvaluation {

	isPotentialSolution := false

	for _, wordInCurrentShortlist := range currrentShortlist {
		if wordInCurrentShortlist.String() == guess.String() {
			isPotentialSolution = true
		}
	}

	return GuessEvaluation{
		Guess:                   guess,
		shortlistSize:           len(currrentShortlist),
		potentialFeedbackCounts: make(map[string]int),
		isPotentialSolution:     isPotentialSolution,
	}
}

func (ge *GuessEvaluation) AddPossibleOutcome(possibleSolution words.Word, feedback game.Feedback) {
	ge.potentialFeedbackCounts[feedback.String()] += 1
}

func (ge *GuessEvaluation) GetWorstCaseScenarioFeedbackString() string {
	if ge.worstCaseScenarioFeedbackString == "" {
		ge.calculate()
	}
	return ge.worstCaseScenarioFeedbackString
}

func (ge *GuessEvaluation) GetWorstCaseShortlistCarryOverRatio() float64 {
	if ge.worstCaseShortlistCarryOverRatio == 0.0 {
		ge.calculate()
	}
	return ge.worstCaseShortlistCarryOverRatio
}

func (ge *GuessEvaluation) isBetterThan(another GuessEvaluation) bool {
	if ge.GetWorstCaseShortlistCarryOverRatio() < another.GetWorstCaseShortlistCarryOverRatio() {
		return true
	}
	if ge.GetWorstCaseShortlistCarryOverRatio() > another.GetWorstCaseShortlistCarryOverRatio() {
		return false
	}
	if ge.isPotentialSolution && !another.isPotentialSolution {
		return true
	}
	return false
}

func (ge *GuessEvaluation) calculate() {

	type worstCaseScenario struct {
		feedbackString string
		count          int
	}
	worst := worstCaseScenario{}

	for potentialFeedbackString, potentialFeedbackCount := range ge.potentialFeedbackCounts {

		if potentialFeedbackCount > worst.count {
			worst = worstCaseScenario{
				feedbackString: potentialFeedbackString,
				count:          potentialFeedbackCount,
			}
		}

	}

	ge.worstCaseShortlistCarryOverRatio = float64(worst.count) / float64(ge.shortlistSize)
	ge.worstCaseScenarioFeedbackString = worst.feedbackString
}
