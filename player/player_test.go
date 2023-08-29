package player

import (
	"testing"
	"wordle/words"
)

func TestChoosingBestWordToFindOutSolution(t *testing.T) {

	var validGuesses = []words.Word{"CHANT", "ZZZZZ"}
	var validSolutions = []words.Word{"SCARE", "SHARE", "SNARE", "STARE"}

	p := NewPlayer(validSolutions, validGuesses)

	got, _ := p.GetNextGuess(false)
	expected := "CHANT"

	if got.String() != expected {
		t.Errorf("Guessing %q would have successfully identified correct soluion, but Player guessed %q", expected, got)
	}
}

func TestGuessingAPossibleSolutionOnLastTurn(t *testing.T) {

	var validGuesses = []words.Word{"CHANT", "SCARE"}
	var validSolutions = []words.Word{"SCARE", "SHARE", "SNARE", "STARE"}

	p := NewPlayer(validSolutions, validGuesses)

	got, _ := p.GetNextGuess(true)
	expected := "SCARE"

	if got.String() != expected {
		t.Errorf("Because it's the last guess, guessing %q might have won the game, but Player guessed %q, which can't win", expected, got)
	}
}
