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

func TestItUnderstandsItIsReducingShortlist(t *testing.T) {

	var validGuesses = []words.Word{"CHANT"}
	var validSolutions = []words.Word{"SCARE", "SHARE", "SNARE", "SPARE", "STARE"}

	p := NewPlayer(validSolutions, validGuesses)

	_, evaluation := p.GetNextGuess(false)
	got := evaluation.worstCaseShortlistCarryOverRatio
	expected := 0.2

	if got != expected {
		t.Errorf("It reduces the shortlist by a factor of %.2f with a guess, but it seems to think the shortlist reduction is %.2f", expected, got)
	}
}

func TestItUnderstandsBadGuessesDontReduceShortlist(t *testing.T) {

	var validGuesses = []words.Word{"XXXXX"}
	var validSolutions = []words.Word{"SCARE", "SHARE", "SNARE", "SPARE", "STARE"}

	p := NewPlayer(validSolutions, validGuesses)

	_, evaluation := p.GetNextGuess(false)
	got := evaluation.worstCaseShortlistCarryOverRatio
	expected := 1.0

	if got != expected {
		t.Errorf("If you force it to choose a really unhelpful guess, shortlist carry-over ratio is %.2f, but it thinks it is %.2f", expected, got)
	}
}
