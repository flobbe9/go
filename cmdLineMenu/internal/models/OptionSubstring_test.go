package models

import "testing"

func TestIsOverlap_shouldNotOverlap(t *testing.T) {
	// baseStr := "abcdefghi"
	var opt1 OptionSubstring
	var opt2 OptionSubstring

	assertDoesNotOverlap := func() {
		if opt1.IsOverlapping(opt2) {
			t.Errorf("Expected '%v' not to overlap '%v'", opt2.Substr, opt1.Substr)
		}
		if opt2.IsOverlapping(opt1) {
			t.Errorf("Expected '%v' not to overlap '%v'", opt1.Substr, opt2.Substr)
		}
	}

	// all blank
	assertDoesNotOverlap()

	opt1 = OptionSubstring{Substr: "", Start: 0}
	opt2 = OptionSubstring{Substr: "ghi", Start: 6}
	assertDoesNotOverlap()

	opt1 = OptionSubstring{Substr: "cdefg", Start: 2}
	opt2 = OptionSubstring{Substr: "hi", Start: 7}
	assertDoesNotOverlap()

	opt1 = OptionSubstring{Substr: "cdefg", Start: 2}
	opt2 = OptionSubstring{Substr: "ab", Start: 0}
	assertDoesNotOverlap()
}

func TestIsOverlap_shouldOverlap(t *testing.T) {
	// baseStr := "abcdefghi"
	var opt1 OptionSubstring
	var opt2 OptionSubstring

	assertOverlap := func() {
		if !opt1.IsOverlapping(opt2) {
			t.Errorf("Expected '%v' to overlap '%v'", opt2.Substr, opt1.Substr)
		}

		if !opt2.IsOverlapping(opt1) {
			t.Errorf("Expected '%v' to overlap '%v'", opt1.Substr, opt2.Substr)
		}
	}

	opt1 = OptionSubstring{Substr: "bcdefg", Start: 1}
	opt2 = OptionSubstring{Substr: "ghi", Start: 6}
	assertOverlap()

	opt1 = OptionSubstring{Substr: "bcdefg", Start: 1}
	opt2 = OptionSubstring{Substr: "fghi", Start: 5}
	assertOverlap()

	opt1 = OptionSubstring{Substr: "bcdefg", Start: 1}
	opt2 = OptionSubstring{Substr: "bcdefghi", Start: 1}
	assertOverlap()

	opt1 = OptionSubstring{Substr: "bcdefg", Start: 1}
	opt2 = OptionSubstring{Substr: "abcdefghi", Start: 0}
	assertOverlap()

	opt1 = OptionSubstring{Substr: "bcdefg", Start: 1}
	opt2 = OptionSubstring{Substr: "abcdefg", Start: 0}
	assertOverlap()

	opt1 = OptionSubstring{Substr: "bcdefg", Start: 1}
	opt2 = OptionSubstring{Substr: "abcd", Start: 0}
	assertOverlap()

	opt1 = OptionSubstring{Substr: "bcdefg", Start: 1}
	opt2 = OptionSubstring{Substr: "ab", Start: 0}
	assertOverlap()

	opt1 = OptionSubstring{Substr: "bcdefg", Start: 1}
	opt2 = OptionSubstring{Substr: "bcdefg", Start: 1}
	assertOverlap()
}