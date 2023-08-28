package player

import (
	"runtime"
	"strings"
	"wordle/domain/game"
	"wordle/domain/words"
)

type Player struct {
	PossibleSolutions words.WordList
	ValidGuesses      words.WordList
}

func (player Player) GetNextGuess(isSixthTurn bool) (words.Word, ProposedGuessEvaluation) {

	bestGuessEvaluation := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}

	if player.PossibleSolutions.Count() == 1 {
		return player.PossibleSolutions.Words[0], bestGuessEvaluation
	}

	if isSixthTurn {
		return player.PossibleSolutions.Words[0], bestGuessEvaluation
	}

	bestGuess := player.identifyBestPossibleGuess(player.ValidGuesses.Words)

	return bestGuess.ProposedGuess, bestGuess
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
func fanoutGuessEvaluation(potentialGuesses []words.Word) <-chan words.Word {
	fanoutChannel := make(chan words.Word)
	go func() {
		for _, potentialGuess := range potentialGuesses {
			select {
				case fanoutChannel <- potentialGuess:
			}
		}
		close(fanoutChannel)
	}()
	return fanoutChannel
}

func (player Player) evaluatePotentialGuesses(fanoutChannel <-chan words.Word) <-chan ProposedGuessEvaluation {
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

func (player Player) identifyBestPossibleGuess(validGuesses []words.Word) ProposedGuessEvaluation {

	// To fan out the guesses to the workers, create a fan out channel
	fanoutChannel := fanoutGuessEvaluation(validGuesses)

	// To collate the results from the workers, create one fan in channel per worker
	noOfWorkers := runtime.NumCPU() - 1
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
