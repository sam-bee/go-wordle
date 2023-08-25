package player

import (
	"fmt"
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

func (player Player) identifyBestPossibleGuess(validGuesses []words.Word) ProposedGuessEvaluation {

	guessEvaluations := proposedGuessEvaluationContainer{value: make([]ProposedGuessEvaluation, 0, 1000)}

	var wg sync.WaitGroup

	for startOfBatch := 0; startOfBatch < len(validGuesses) - 1; startOfBatch = startOfBatch+1000 {
		lengthOfBatch := min(1000, len(validGuesses) - startOfBatch -1) + startOfBatch // 1000, or the remainder of the list
		wg.Add(1)
		go evaluatePossibleGuesses(validGuesses[startOfBatch:lengthOfBatch], &guessEvaluations, &wg, player)
	}

	wg.Wait()

	bestGuessEvaluation := ProposedGuessEvaluation{worstCaseShortlistCarryOverRatio: 1.0}

	for _, guessEvaluation := range guessEvaluations.value {
		if guessEvaluation.isBetterThan(bestGuessEvaluation) {
			bestGuessEvaluation = guessEvaluation
		}
	}

	return bestGuessEvaluation
}

func evaluatePossibleGuesses(proposedGuesses []words.Word, guessEvaluations *proposedGuessEvaluationContainer, wg *sync.WaitGroup, player Player) {

	evaluations := make([]ProposedGuessEvaluation, 0, len(proposedGuesses))

	for _, proposedGuess := range proposedGuesses {
		proposedGuessEvaluation := player.EvaluatePossibleGuess(proposedGuess)
		evaluations = append(evaluations, proposedGuessEvaluation)
	}

	guessEvaluations.mu.Lock()
	guessEvaluations.value = append(guessEvaluations.value, evaluations...)
	guessEvaluations.mu.Unlock()

	wg.Done()
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
