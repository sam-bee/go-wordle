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
	guess words.Word,
	currrentShortlist []words.Word,
) ProposedGuessEvaluation {

	isPotentialSolution := false

	for _, wordInCurrentShortlist := range currrentShortlist {
		if wordInCurrentShortlist.String() == guess.String() {
			isPotentialSolution = true
		}
	}

	return ProposedGuessEvaluation{
		Guess: guess,
		shortlistSize: len(currrentShortlist),
		potentialFeedbackCounts: make(map[string]int),
		isPotentialSolution: isPotentialSolution,
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

	type worstCaseScenario struct {
		feedbackString string
		count          int
	}
	worst := worstCaseScenario{}

	for potentialFeedbackString, potentialFeedbackCount := range proposedGuessEvaluation.potentialFeedbackCounts {

		if potentialFeedbackCount > worst.count {
			worst = worstCaseScenario {
				feedbackString: potentialFeedbackString,
				count: potentialFeedbackCount,
			}
		}

	}

	proposedGuessEvaluation.worstCaseShortlistCarryOverRatio = float64(worst.count) / float64(proposedGuessEvaluation.shortlistSize)
	proposedGuessEvaluation.worstCaseScenarioFeedbackString = worst.feedbackString
}
