package models

// Wrapper around an arbitrary, non-blank substring of an option including it's index
type OptionSubstring struct {
	// Substring of an option
	Substr string

	// Index of [substr] in option. May be -1
	Index int
}
