package cmd

import (
	"fmt"
	"io"
	"strconv"
	"wordle/game"
	"wordle/player"
	"wordle/words"

	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play a wordle game",
	Long:  `e.g. wordle play SPARE`,
	RunE: func(cmd *cobra.Command, args []string) error {
		solutionArgument := args[0]
		writer := cmd.OutOrStdout()
		err := playWordleGame(cmd, solutionArgument, writer)
		if err != nil {
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

	// To initialise the solution, parse command line argument
	solution, err := words.NewWord(solutionArgument)
	if err != nil {
		return err
	}

	// To find out what our guesses might be, read guesses word list from file
	fmt.Fprintln(writer, "Reading valid guesses from file...")
	validGuessesWordList, err := words.GetValidGuessesWordList()
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, "Found %d words\n", len(validGuessesWordList))

	// To find out what the solution might be, read guesses word list from file
	fmt.Fprintln(writer, "Reading valid solutions from file...")
	validSolutionsWordList, err := words.GetValidSolutionsWordList()
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, "Found %d words\n", len(validSolutionsWordList))

	player := player.Player{PossibleSolutions: validSolutionsWordList, ValidGuesses: validGuessesWordList}

	turn := 1
	won := false

	for turn <= 6 && !won {
		printPreAnalysis(player)
		guessWasSolution, guess, feedback, evaluation := takeGuess(turn, &player, solution)
		printEvaluation(evaluation)
		won = guessWasSolution
		printTurn(guess, feedback, turn)
		turn += 1
	}

	printOutcome(won, turn-1)
	return nil
}

func takeGuess(guessNo int, player *player.Player, solution words.Word) (won bool, guess words.Word, feedback game.Feedback, evaluation player.ProposedGuessEvaluation) {
	guess, evaluation = player.GetNextGuess(guessNo == 6)
	won = guess.String() == solution.String()
	feedback = game.GetFeedback(solution, guess)
	player.TakeFeedbackFromGuess(guess, feedback)
	return
}

func printTurn(guess words.Word, feedback game.Feedback, guessNo int) {
	fmt.Println("Guess number " + strconv.Itoa(guessNo) + ": " + guess.String())
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
	fmt.Println("The next guess should be " + evaluation.Guess.String())
	fmt.Println("Worst-case scenario for proposed guess is the feedback " + evaluation.GetWorstCaseScenarioFeedbackString() + ". Carry-over ratio for possible solutions list would be " + evaluation.GetWorstCaseShortlistCarryOverRatioString())
}

func printOutcome(won bool, turnNumber int) {
	if won {
		fmt.Println("Won the Wordle in " + strconv.Itoa(turnNumber) + " turns")
	} else {
		fmt.Println("Lost the Wordle after 6 turns :-(")
	}
	fmt.Println()
}
