package cmd

import (
	"fmt"
	"io"
	"github.com/spf13/cobra"
	"strconv"
	"wordle/domain/game"
	"wordle/domain/player"
	"wordle/domain/words"
	"wordle/infrastructure"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play a wordle game",
	Long: `e.g. wordle play SPARE`,
	RunE: func(cmd *cobra.Command, args []string) error {
		solutionArgument := args[0]
		writer := cmd.OutOrStdout()
		err := playWordleGame(cmd, solutionArgument, writer)
		if (err != nil) {
			return err
		}
		fmt.Fprintln(writer, "Finished!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}

func playWordleGame(cmd *cobra.Command, solutionArgument string, writer io.Writer) error {

	// To find out what our guesses might be, read guesses word list from file
	fmt.Fprintln(writer, "Reading valid guesses from file...")
	validGuessesWordList, err := infrastructure.GetValidGuessesWordList(writer)
	if (err != nil) {
		return err
	}
	fmt.Fprintf(writer, "Found %v words\n", words.Count(validGuessesWordList))

	// To find out what the solution might be, read guesses word list from file
	fmt.Fprintln(writer, "Reading valid solutions from file...")
	validSolutionsWordList, err := infrastructure.GetValidSolutionsWordList(writer)
	if (err != nil) {
		return err
	}
	fmt.Fprintf(writer, "Found %v words\n", words.Count(validSolutionsWordList))

	solution := words.MakeWord(solutionArgument)

	player := player.Player{PossibleSolutions: validSolutionsWordList, ValidGuesses: validGuessesWordList}

	turnNumber := 1
	won := false


	for turnNumber <= 6 && !won {
		printPreAnalysis(player)
		guessWasSolution, guess, feedback, evaluation := takeGuess(turnNumber, &player, solution)
		printEvaluation(evaluation)
		won = guessWasSolution
		printTurn(guess, feedback, turnNumber)
		turnNumber += 1
	}

	printOutcome(won, turnNumber-1)
	return nil
}

func takeGuess(guessNo int, player *player.Player, solution words.Word) (won bool, guess words.Word, feedback game.Feedback, evaluation player.ProposedGuessEvaluation) {
	guess, evaluation = player.GetNextGuess(guessNo == 6)
	won = guess.Equals(solution)
	feedback = game.GetFeedback(solution, guess)
	player.TakeFeedbackFromGuess(guess, feedback)
	return
}

func printTurn(guess words.Word, feedback game.Feedback, guessNo int) {
	fmt.Println("Guess number " + strconv.Itoa(guessNo) + ": " + guess.String)
	fmt.Println("Feedback from guess was: " + feedback.String())
	fmt.Println()
}

func printPreAnalysis(player player.Player) {
	noOfPossibleSolutions := player.GetNoOfPossibleSolutions()
	fmt.Print("There are currently " + strconv.Itoa(noOfPossibleSolutions) + " possible solutions")
	if noOfPossibleSolutions <= 10 {
		fmt.Println(" [" + player.GetPossibleSolutions() + "]")
	} else {
		fmt.Println()
	}
}

func printEvaluation(evaluation player.ProposedGuessEvaluation) {
	if !evaluation.IsNullEvaluation() {
		fmt.Println("The next guess should be " + evaluation.ProposedGuess.String)
		fmt.Println("Worst-case scenario for proposed guess is the feedback " + evaluation.GetWorstCaseScenarioFeedbackString() + ". Carry-over ratio for possible solutions list would be " + evaluation.GetWorstCaseShortlistCarryOverRatioString())
	}
}

func printOutcome(won bool, turnNumber int) {
	if won {
		fmt.Println("Won the Wordle in " + strconv.Itoa(turnNumber) + " turns")
	} else {
		fmt.Println("Lost the Wordle after 6 turns :-(")
	}
	fmt.Println()
}
