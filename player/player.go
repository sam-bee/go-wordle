package player

import (
	"runtime"
	"strings"
	"wordle/game"
	"wordle/words"
)

type Player struct {
	PossibleSolutions []words.Word
	ValidGuesses      []words.Word
}

func (player *Player) GetNextGuess(isSixthTurn bool) (words.Word, ProposedGuessEvaluation) {

	bestGuessEvaluation := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}

	if len(player.PossibleSolutions) == 1 {
		return player.PossibleSolutions[0], bestGuessEvaluation
	}

	if isSixthTurn {
		return player.PossibleSolutions[0], bestGuessEvaluation
	}

	bestGuess := player.identifyBestPossibleGuess(player.ValidGuesses)

	return bestGuess.ProposedGuess, bestGuess
}

func (player *Player) EvaluatePossibleGuess(possibleGuess words.Word) ProposedGuessEvaluation {

	proposedGuessEvaluation := MakeProposedGuessEvaluation(possibleGuess, len(player.PossibleSolutions), player.PossibleSolutions)

	for _, possibleSolution := range player.PossibleSolutions {
		feedback := game.GetFeedback(possibleSolution, possibleGuess)
		proposedGuessEvaluation.AddPossibleOutcome(possibleSolution, feedback)
	}

	return proposedGuessEvaluation
}

func (player *Player) TakeFeedbackFromGuess(word words.Word, feedback game.Feedback) {

	var newShortlist []words.Word

	for _, solutionStillOnShortlist := range player.PossibleSolutions {
		feedbackIfThisWordWereSolution := game.GetFeedback(solutionStillOnShortlist, word)
		if feedbackIfThisWordWereSolution.Equals(feedback) {
			newShortlist = append(newShortlist, solutionStillOnShortlist)
		}
	}

	player.PossibleSolutions = newShortlist
}

func (player *Player) GetNoOfPossibleSolutions() int {
	return len(player.PossibleSolutions)
}

func (player *Player) GetPossibleSolutions() string {

	var wordsAsStrings []string
	for _, word := range player.PossibleSolutions {
		wordsAsStrings = append(wordsAsStrings, word.String())
	}

	return strings.Join(wordsAsStrings, ", ")
}

func fanoutGuessEvaluation(potentialGuesses []words.Word) <-chan words.Word {
	fanoutChannel := make(chan words.Word)
	go func() {
		for _, potentialGuess := range potentialGuesses {
			fanoutChannel <- potentialGuess
		}
		close(fanoutChannel)
	}()
	return fanoutChannel
}

func (player *Player) evaluatePotentialGuesses(fanoutChannel <-chan words.Word) <-chan ProposedGuessEvaluation {
	faninChannel := make(chan ProposedGuessEvaluation)
	go func() {

		bestGuessEvaluation := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}

		for proposedGuess := range fanoutChannel {
			proposedGuessEvaluation := player.EvaluatePossibleGuess(proposedGuess)
			if proposedGuessEvaluation.isBetterThan(bestGuessEvaluation) {
				bestGuessEvaluation = proposedGuessEvaluation
			}
		}

		faninChannel <- bestGuessEvaluation
		close(faninChannel)
	}()
	return faninChannel
}

func (player *Player) identifyBestPossibleGuess(validGuesses []words.Word) ProposedGuessEvaluation {

	// To fan out the guesses to the workers, create a fan out channel
	fanoutChannel := fanoutGuessEvaluation(validGuesses)

	// To collate the results from the workers, create one fan in channel per worker
	noOfWorkers := max(runtime.NumCPU() - 1, 1)
	fanInChannels := make([]<-chan ProposedGuessEvaluation, noOfWorkers)

	for i := 0; i < noOfWorkers; i++ {
		fanInChannels[i] = player.evaluatePotentialGuesses(fanoutChannel)
	}

	// To identify the best guess from any of the workers, loop through their channels
	bestGuessEvaluation := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}
	for i := range fanInChannels {
		for proposedGuessEvaluation := range fanInChannels[i] {
			if proposedGuessEvaluation.isBetterThan(bestGuessEvaluation) {
				bestGuessEvaluation = proposedGuessEvaluation
			}
		}
	}

	return bestGuessEvaluation
}
