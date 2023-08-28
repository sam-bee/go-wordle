package cmd

import (
	"fmt"
	"io"
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
	fmt.Fprintf(writer, "Found %d words\n\n", len(validSolutionsWordList))

	player := player.Player{PossibleSolutions: validSolutionsWordList, ValidGuesses: validGuessesWordList}

	turn := 1
	won := false

	for turn <= 6 && !won {
		printPreAnalysis(writer, player)
		guessWasSolution, guess, feedback, evaluation := takeGuess(turn, &player, solution)
		printEvaluation(writer, evaluation, player)
		won = guessWasSolution
		printTurn(writer, guess, feedback, turn)
		turn += 1
	}

	printOutcome(writer, won, turn-1)
	return nil
}

func takeGuess(guessNo int, player *player.Player, solution words.Word) (won bool, guess words.Word, feedback game.Feedback, evaluation player.GuessEvaluation) {
	guess, evaluation = player.GetNextGuess(guessNo == 6)
	won = guess.String() == solution.String()
	feedback = game.GetFeedback(solution, guess)
	player.TakeFeedbackFromGuess(guess, feedback)
	return
}

func printTurn(writer io.Writer, guess words.Word, feedback game.Feedback, guessNo int) {
	fmt.Fprintf(writer, "Guess number %d: %q\n", guessNo, guess.String())
	fmt.Fprintf(writer, "Feedback from guess was: %q\n", feedback.String())
	fmt.Fprintln(writer)
}

func printPreAnalysis(writer io.Writer, player player.Player) {
	noOfPossibleSolutions := player.GetNoOfPossibleSolutions()
	fmt.Fprintf(writer, "There are currently %d possible solutions\n", player.GetNoOfPossibleSolutions())
	if noOfPossibleSolutions <= 10 {
		fmt.Fprintf(writer, "The remaining possible solutions are: [%s]\n", player.GetPossibleSolutions())
	}
}

func printEvaluation(writer io.Writer, evaluation player.GuessEvaluation, player player.Player) {
	fmt.Fprintf(writer, "The next guess should be %q\n", evaluation.Guess.String())

	if player.GetNoOfPossibleSolutions() > 1 {
		fmt.Fprintf(
			writer, "Worst-case scenario for guess is the feedback %q. Carry-over ratio for possible solutions list would be %.2f%%\n",
			evaluation.GetWorstCaseScenarioFeedbackString(),
			100*evaluation.GetWorstCaseShortlistCarryOverRatio(),
		)
	}
}

func printOutcome(writer io.Writer, won bool, turnNumber int) {
	if won {
		fmt.Fprintf(writer, "Won the Wordle in %d turns", turnNumber)
	} else {
		fmt.Fprintln(writer, "Lost the Wordle after 6 turns :-(")
	}
	fmt.Fprintln(writer)
}
