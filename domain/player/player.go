package player

import (
	"strings"
	"wordle/domain/game"
	"wordle/domain/words"
)

type Player struct {
	PossibleSolutions words.WordList
	ValidGuesses      words.WordList
}

func (player Player) GetNextGuess(isSixthTurn bool) (guess words.Word, evaluation ProposedGuessEvaluation) {

	bestGuessEvaluation := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}

	if player.PossibleSolutions.Count() == 1 {
		return player.PossibleSolutions.Words[0], bestGuessEvaluation
	}

	if isSixthTurn {
		return player.PossibleSolutions.Words[0], bestGuessEvaluation
	}

	for _, proposedGuess := range player.ValidGuesses.Words {
		proposedGuessEvaluation := player.EvaluatePossibleGuess(proposedGuess)

		if proposedGuessEvaluation.isBetterThan(bestGuessEvaluation) {
			bestGuessEvaluation = proposedGuessEvaluation
		}
	}

	return bestGuessEvaluation.ProposedGuess, bestGuessEvaluation
}

func (player Player) EvaluatePossibleGuess(possibleGuess words.Word) ProposedGuessEvaluation {

	proposedGuessEvaluation := MakeProposedGuessEvaluation(possibleGuess, player.PossibleSolutions.Count(), player.PossibleSolutions)

	for _, possibleSolution := range player.PossibleSolutions.Words {
		feedback := game.GetFeedback(possibleSolution, possibleGuess)
		proposedGuessEvaluation.AddPossibleOutcome(possibleSolution, feedback)
	}

	return proposedGuessEvaluation
}

func (player *Player) TakeFeedbackFromGuess(word words.Word, feedback game.Feedback) {

	var newShortlist []words.Word

	for _, solutionStillOnShortlist := range player.PossibleSolutions.Words {
		feedbackIfThisWordWereSolution := game.GetFeedback(solutionStillOnShortlist, word)
		if feedbackIfThisWordWereSolution.Equals(feedback) {
			newShortlist = append(newShortlist, solutionStillOnShortlist)
		}
	}

	player.PossibleSolutions = words.WordList{Words: newShortlist}
}

func (player Player) GetNoOfPossibleSolutions() int {
	return player.PossibleSolutions.Count()
}

func (player Player) GetPossibleSolutions() string {

	var wordsAsStrings []string
	for _, word := range player.PossibleSolutions.Words {
		wordsAsStrings = append(wordsAsStrings, word.String)
	}

	return strings.Join(wordsAsStrings, ", ")
}
