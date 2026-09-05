package models

type State struct {
	// Index of the select option the user is currently focusing (using arrow keys)
	FocusedOptionIndex int

	currentRowIndex int
}

func (this *State) GetCurrentRowIndex() int {
	return this.currentRowIndex
}