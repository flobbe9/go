package models

type State struct {
	// Index of the select option the user is currently focusing (using arrow keys)
	FocusedOptionIndex int
	
	// User search input if enabled
	SearchQuery string;

	// The 0-based cursor index for the search input line
	SearchQueryCursorIndex int;

	// List of menu options. May be updated by search results produced from [SearchQuery] (if search is enabled)
	Options []string;
}
