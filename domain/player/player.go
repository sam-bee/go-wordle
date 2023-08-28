package player

import (
	"runtime"
	"strings"
	"sync"
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

type proposedGuessEvaluationContainer struct {
	mu    sync.Mutex
	value []ProposedGuessEvaluation
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
func fanoutGuessEvaluation(signalChannel <-chan struct{}, potentialGuesses []words.Word) <-chan words.Word {
	fanoutChannel := make(chan words.Word)
	go func() {
		for _, potentialGuess := range potentialGuesses {
			select {
			case <-signalChannel:
				return
			case fanoutChannel <- potentialGuess:
			}
		}
		close(fanoutChannel)
	}()
	return fanoutChannel
}

func (player Player) evaluatePotentialGuesses(signalChannel <-chan struct{}, fanoutChannel <-chan words.Word) <-chan ProposedGuessEvaluation {
	faninChannel := make(chan ProposedGuessEvaluation)
	go func() {
		for potentialGuess := range fanoutChannel {
			select {
			case <-signalChannel:
				return
			case faninChannel <- player.EvaluatePossibleGuess(potentialGuess):
			}
		}
		close(faninChannel)
	}()
	return faninChannel
}

func mergeChannelsToMultiplex(signalChannel <-chan struct{}, faninChannels ...<-chan ProposedGuessEvaluation) <-chan ProposedGuessEvaluation {
	var wg sync.WaitGroup

	wg.Add(len(faninChannels))
	multiplexChannel := make(chan ProposedGuessEvaluation)
	multiplex := func(c <-chan ProposedGuessEvaluation) {
		defer wg.Done()
		for i := range c {
			select {
			case <-signalChannel:
				return
			case multiplexChannel <- i:
			}
		}
	}
	for _, c := range faninChannels {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(multiplexChannel)
	}()
	return multiplexChannel
}

func (player Player) identifyBestPossibleGuess(validGuesses []words.Word) ProposedGuessEvaluation {

	// To enable the workers to be shut down, create a signal channel to tell them when to stop
	signalChannel := make(chan struct{})
	defer close(signalChannel)

	// To fan out the guesses to the workers, create a fan out channel
	fanoutChannel := fanoutGuessEvaluation(signalChannel, validGuesses)

	// To collate the results from the workers, create one fan in channel per worker
	noOfWorkers := runtime.NumCPU() - 1
	fanInChannels := make([]<-chan ProposedGuessEvaluation, noOfWorkers)

	for i := 0; i < noOfWorkers; i++ {
		fanInChannels[i] = player.evaluatePotentialGuesses(signalChannel, fanoutChannel)
	}

	// To multiplex the results from the workers, create a multiplex channel
	multiplexChannel := mergeChannelsToMultiplex(signalChannel, fanInChannels...)

	// To identify the best guess, iterate over the multiplex channel
	bestGuessEvaluation := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}

	for proposedGuessEvaluation := range multiplexChannel {
		if proposedGuessEvaluation.isBetterThan(bestGuessEvaluation) {
			bestGuessEvaluation = proposedGuessEvaluation
		}
	}
	return bestGuessEvaluation
}
