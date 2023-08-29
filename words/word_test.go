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

func TestEquals(t *testing.T) {
	w, _ := NewWord("SPARE")
	same, _ := NewWord("SPARE")
	different, _ := NewWord("MONEY")

	if !w.Equals(same) {
		t.Errorf("String equality check failed with %q and %q", w, same)
	}

	if w.Equals(different) {
		t.Errorf("String equality check failed with %q and %q", w, different)
	}
}
