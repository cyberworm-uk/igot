package ert_test

import (
	"igot/runes"
	"testing"
)

func TestCost(t *testing.T) {
	if runes.Calc(355173) != 713 {
		t.Fatalf("bad runes")
	}
	if runes.Calc(1) != 1 {
		t.Fatalf("bad runes")
	}
	if runes.Calc(6051) != 150 {
		t.Fatalf("bad runes")
	}
	if runes.Calc(0) != 0 {
		t.Fatalf("bad runes")
	}
	if runes.Calc(355174) != 0 {
		t.Fatalf("bad runes")
	}
}
