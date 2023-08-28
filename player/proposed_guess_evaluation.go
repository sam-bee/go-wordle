package player

import (
	"wordle/game"
	"wordle/words"
)

type ProposedGuessEvaluation struct {
	Guess                            words.Word
	shortlistSize                    int
	potentialFeedbackCounts          map[string]int
	worstCaseScenarioFeedbackString  string
	worstCaseShortlistCarryOverRatio float64
	isPotentialSolution              bool
}

func MakeProposedGuessEvaluation(
	proposedGuess words.Word,
	sizeOfCurrentShortlist int,
	currrentShortlist []words.Word,
) ProposedGuessEvaluation {

	isPotentialSolution := false

	for _, wordInCurrentShortlist := range currrentShortlist {
		if wordInCurrentShortlist.String() == proposedGuess.String() {
			isPotentialSolution = true
		}
	}

	return ProposedGuessEvaluation{
		proposedGuess,
		sizeOfCurrentShortlist,
		make(map[string]int),
		"",
		0.0,
		isPotentialSolution,
	}
}

func (proposedGuessEvaluation *ProposedGuessEvaluation) AddPossibleOutcome(possibleSolution words.Word, feedback game.Feedback) {
	proposedGuessEvaluation.potentialFeedbackCounts[feedback.String()] += 1
}

func (proposedGuessEvaluation *ProposedGuessEvaluation) GetWorstCaseScenarioFeedbackString() string {
	if proposedGuessEvaluation.worstCaseScenarioFeedbackString == "" {
		proposedGuessEvaluation.calculate()
	}
	return proposedGuessEvaluation.worstCaseScenarioFeedbackString
}

func (proposedGuessEvaluation *ProposedGuessEvaluation) GetWorstCaseShortlistCarryOverRatio() float64 {
	if proposedGuessEvaluation.worstCaseShortlistCarryOverRatio == 0.0 {
		proposedGuessEvaluation.calculate()
	}
	return proposedGuessEvaluation.worstCaseShortlistCarryOverRatio
}

func (proposedGuessEvaluation *ProposedGuessEvaluation) isBetterThan(another ProposedGuessEvaluation) bool {
	if proposedGuessEvaluation.GetWorstCaseShortlistCarryOverRatio() < another.GetWorstCaseShortlistCarryOverRatio() {
		return true
	}
	if proposedGuessEvaluation.GetWorstCaseShortlistCarryOverRatio() > another.GetWorstCaseShortlistCarryOverRatio() {
		return false
	}
	if proposedGuessEvaluation.isPotentialSolution && !another.isPotentialSolution {
		return true
	}
	return false
}

func (proposedGuessEvaluation *ProposedGuessEvaluation) calculate() {
	worstCaseScenario := struct {
		feedbackString string
		count          int
	}{
		"",
		0,
	}

	for potentialFeedbackString, potentialFeedbackCount := range proposedGuessEvaluation.potentialFeedbackCounts {

		if potentialFeedbackCount > worstCaseScenario.count {
			worstCaseScenario = struct {
				feedbackString string
				count          int
			}{
				potentialFeedbackString,
				potentialFeedbackCount,
			}
		}

	}

	proposedGuessEvaluation.worstCaseShortlistCarryOverRatio = float64(worstCaseScenario.count) / float64(proposedGuessEvaluation.shortlistSize)
	proposedGuessEvaluation.worstCaseScenarioFeedbackString = worstCaseScenario.feedbackString
}
