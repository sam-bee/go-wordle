package words

import (
	"testing"
)

func TestConstructorPreventsInvalidState(t *testing.T) {
	_, err := NewWord("")

	if err == nil {
		t.Errorf("Expected error when constructing word with empty string")
	}
}

func TestStringConversion(t *testing.T) {
	w, _ := NewWord("SPARE")
	res := w.String()

	if res != "SPARE" {
		t.Errorf("String conversion failed, expected 'SPARE', got '%s'", res)
	}
}
