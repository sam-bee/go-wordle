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

func (player *Player) GetNextGuess(lastTurn bool) (words.Word, ProposedGuessEvaluation) {

	bestGuessEvaluation := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}

	if len(player.PossibleSolutions) == 1 {
		return player.PossibleSolutions[0], bestGuessEvaluation
	}

	if lastTurn {
		return player.PossibleSolutions[0], bestGuessEvaluation
	}

	bestGuess := player.identifyBestPossibleGuess(player.ValidGuesses)

	return bestGuess.Guess, bestGuess
}

func (player *Player) EvaluatePossibleGuess(guess words.Word) ProposedGuessEvaluation {

	proposedGuessEvaluation := MakeProposedGuessEvaluation(guess, len(player.PossibleSolutions), player.PossibleSolutions)

	for _, possibleSolution := range player.PossibleSolutions {
		feedback := game.GetFeedback(possibleSolution, guess)
		proposedGuessEvaluation.AddPossibleOutcome(possibleSolution, feedback)
	}

	return proposedGuessEvaluation
}

func (player *Player) TakeFeedbackFromGuess(word words.Word, feedback game.Feedback) {

	var shortlist []words.Word

	for _, solutionStillOnShortlist := range player.PossibleSolutions {
		feedbackIfThisWordWereSolution := game.GetFeedback(solutionStillOnShortlist, word)
		if feedbackIfThisWordWereSolution.Equals(feedback) {
			shortlist = append(shortlist, solutionStillOnShortlist)
		}
	}

	player.PossibleSolutions = shortlist
}

func (player *Player) GetNoOfPossibleSolutions() int {
	return len(player.PossibleSolutions)
}

func (player *Player) GetPossibleSolutions() string {

	var words []string
	for _, word := range player.PossibleSolutions {
		words = append(words, word.String())
	}

	return strings.Join(words, ", ")
}

func fanoutGuessEvaluation(potentialGuesses []words.Word) <-chan words.Word {
	fanoutChannel := make(chan words.Word)
	go func() {
		for _, g := range potentialGuesses {
			fanoutChannel <- g
		}
		close(fanoutChannel)
	}()
	return fanoutChannel
}

func (player *Player) evaluatePotentialGuesses(fanoutChannel <-chan words.Word) <-chan ProposedGuessEvaluation {
	faninChannel := make(chan ProposedGuessEvaluation)
	go func() {

		bestGuess := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}

		for word := range fanoutChannel {
			evaluation := player.EvaluatePossibleGuess(word)
			if evaluation.isBetterThan(bestGuess) {
				bestGuess = evaluation
			}
		}

		faninChannel <- bestGuess
		close(faninChannel)
	}()
	return faninChannel
}

func (player *Player) identifyBestPossibleGuess(validGuesses []words.Word) ProposedGuessEvaluation {

	// To fan out the guesses to the workers, create a fan out channel
	fanoutChannel := fanoutGuessEvaluation(validGuesses)

	// To collate the results from the workers, create one fan in channel per worker
	noOfWorkers := max(runtime.NumCPU()-1, 1)
	fanInChannels := make([]<-chan ProposedGuessEvaluation, noOfWorkers)

	for i := 0; i < noOfWorkers; i++ {
		fanInChannels[i] = player.evaluatePotentialGuesses(fanoutChannel)
	}

	// To identify the best guess from any of the workers, loop through their channels
	best := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}
	for i := range fanInChannels {
		for guess := range fanInChannels[i] {
			if guess.isBetterThan(best) {
				best = guess
			}
		}
	}

	return best
}
