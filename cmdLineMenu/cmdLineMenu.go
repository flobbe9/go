package cmdLineMenu

import (
	"fmt"
	"regexp"
	"slices"
	"unicode/utf8"

	"github.com/flobbe9/go/cmdLineMenu/constants/ansiKey"
	"github.com/flobbe9/go/cmdLineMenu/constants/key"
	internalModels "github.com/flobbe9/go/cmdLineMenu/internal/models"
	"github.com/flobbe9/go/cmdLineMenu/internal/scanner"
	"github.com/flobbe9/go/cmdLineMenu/internal/search"
	"github.com/flobbe9/go/cmdLineMenu/internal/utils/ansiUtils"
	"github.com/flobbe9/go/cmdLineMenu/models"
	"github.com/flobbe9/go/utils/sliceUtils"
	"github.com/flobbe9/go/utils/stringUtils"

	"github.com/charmbracelet/x/ansi"
)

// Store the current state of the menu depending on user input.
var state = &internalModels.State{};

// options to choose from
var initialOptions []string; // initialOptions constant
var question string;

// Global settings passed to [Prompt()]
var config *models.Config;

// The prefix for every menu option
const OPTION_PREFIX = "> ";
const NO_SEARCH_RESULTS_MSG_START = "No results for search";


// Prompt user to answer [question] by selecting from [options] using Arrow-up / -down keys and submitting
// with Enter.
// 
// [return] the selected option from [options].
func Prompt(_question string, _initialOptions []string, _config models.Config) (string, error) {
	if (len(_initialOptions) <= 0) {
		return "", fmt.Errorf("Specify at least one option");
	}

	// init global vars
	config = &_config;
	state.Options, initialOptions = _initialOptions, _initialOptions;
	question = _question;
	state.FocusedOptionIndex = 0;
	// clear state on exit
	defer func() {
		state = &internalModels.State{};
	}();

	// print question
	if (config.IsShowSubmitHint) {
		fmt.Printf("%v (Submit with Enter)\n", question);

	} else {
		fmt.Println(question);
	}

	// first menu-print
	rerender(formatMenu(), 0);
	if (config.IsOptionSearchEnabled) {
		updateSearchInputPlaceholder();
	}

	// await user input
	err := handleUserInput();
	if err != nil {
		return "", err;
	}
	
	answer, err := getCurrentlyFocusedOption();
	if err != nil {
		return "", err;
	}

	if (config.IsDisplayAnswer) {
		printAnswer(answer);
	}

	if (config.IsClearMenuOnSubmit) {
		clearMenu();
	}

	return answer, nil;
}

// Blocks until user has pressed Enter key. Increment / Decrement [state.FocusedOptionIndex] depending on arrow-down/-up keys.
func handleUserInput() error {
	for {
		err := scanner.ScanlnRaw(func(rawStdin models.RawStdin, line string, cursorIndex int) {
			if config.IsOptionSearchEnabled {
				state.SearchQuery = line;
				state.SearchQueryCursorIndex = cursorIndex;
			}

			handleStdinChange := func() {
				if config.IsOptionSearchEnabled {
					prevOptionsLen := len(state.Options);
					updateSearchInputPlaceholder();
					searchOptionsAndUpdateStates();
					rerender(formatMenu(), prevOptionsLen);
				}
			}

			switch rawStdin.GetSignificantAsciiKey() {
			case key.Backspace, key.CtrlBackspace:
				handleStdinChange();
			}

			// case: did print something
			if rawStdin.IsUtf8() {
				handleStdinChange();
			}

			switch rawStdin.GetSignificantAnsiKey() {
			case ansiKey.ArrowDown:
				if len(state.Options) > 0 {
					updateFocusedOptionIndex(false);
					rerender(formatMenu(), len(state.Options));
				}

			case ansiKey.ArrowUp:
				if len(state.Options) > 0 {
					updateFocusedOptionIndex(true);
					rerender(formatMenu(), len(state.Options));
				}

			case ansiKey.Delete:
				handleStdinChange();
			}
		})

		// don't exit if nothing focused
		if !areSearchResultsEmpty() || err != nil {
			return err;
		}
	}
}

// Indicates that search has produced no results.
//
// [return] [true] if [state.FocusedOptionIndex == -1]
func areSearchResultsEmpty() bool {
	return state.FocusedOptionIndex == -1;
}

// [return] the [state.Option] currently focused using [state.FocusedOptionIndex]. Return error if 
// search results are empty.
func getCurrentlyFocusedOption() (string, error) {
	if areSearchResultsEmpty() {
		return "", fmt.Errorf("Failed to get focused option. No search results");
	}
	return ansi.Strip(state.Options[state.FocusedOptionIndex]), nil;
}

// Search [state.Options] using [state.SearchQuery] and update options and [state.FocusedOptionIndex] accordingly.
func searchOptionsAndUpdateStates() {
	prevSearchResults := slices.Clone(state.Options);

	if stringUtils.IsBlank(state.SearchQuery) {
		didSearchResultsChange := !ansiUtils.EqualsSlicesIgnoreAnsi(prevSearchResults, initialOptions);
		state.Options = initialOptions;

		if didSearchResultsChange {
			state.FocusedOptionIndex = 0;
		}
		return;
	}

	searchResults := search.SearchAndHighlightOptions(initialOptions, state.SearchQuery);

	if len(searchResults) == 0 {
		state.Options = []string{fmt.Sprintf("%v '%v'", NO_SEARCH_RESULTS_MSG_START, state.SearchQuery)};
		state.FocusedOptionIndex = -1; 

	} else {
		state.Options = searchResults;

		didSearchResultsChange := !ansiUtils.EqualsSlicesIgnoreAnsi(prevSearchResults, searchResults);
		if didSearchResultsChange {
			state.FocusedOptionIndex = 0;
		}
	}
}

// Override existing stdout by moving the current terminal's cursor up by [backwardLines] and then
// printing [content] (not printing a line break at the end).
//
// [content] to render after moving back
//
// [backwardLines] number of lines to move the cursor up before printing [content]. Use 0 or 1 to stay at current line
func rerender(content []string, backwardLines int) {
	// move up
	if (backwardLines > 0) {
		fmt.Printf("%v", ansi.CursorPreviousLine(backwardLines));
	}

	// delete old menu
	fmt.Print(ansi.DeleteLine(backwardLines));

	// print new menu
	for _, line := range content {
		fmt.Printf("%v%v", ansi.InsertLine(1), line);
	}

	// move search query cursor back
	if config.IsOptionSearchEnabled && state.SearchQueryCursorIndex > 0 {
		fmt.Print(ansi.CursorForward(state.SearchQueryCursorIndex));
	}
}

// [return] formatted menu line including a line break and possibly underlined if [isFocused == true]
func formatMenuLine(lineContent string, isFocused bool) string {
	if isFocused {
		var anyClosingSeq = regexp.MustCompile(`\x1b\[0?m`); 
		// Make sure that underline sequence is not closed by a color sequence. 
		// Replace all closing sequences with underline opening seq. Ansi still works if there're more opening 
		// sequences than closing ones, as long as the whole string ends with a closing sequence.
		lineContent = anyClosingSeq.ReplaceAllString(lineContent, "$0\x1b[4m");
	}

	line := fmt.Sprintf("%v\n", ansi.NewStyle().Underline(isFocused).Styled(lineContent));

	if !stringUtils.StartsWith(lineContent, NO_SEARCH_RESULTS_MSG_START) {
		line = fmt.Sprintf("%v%v", OPTION_PREFIX, line);
	}

	return line;
}

// [return] the whole menu consisting of all select options. End on a new line
func formatMenu() []string {
	return sliceUtils.Map(state.Options, func(option string, i int) string {
		return formatMenuLine(option, i == state.FocusedOptionIndex);
	});
}

// Increment / Decrement [state.FocusedOptionIndex] by one making sure it loops to the start / end if out of bounds
func updateFocusedOptionIndex(isDecrease bool) {
	if (isDecrease) {
		state.FocusedOptionIndex--;
	} else {
		state.FocusedOptionIndex++;
	}
	
	// loop user selection
	if (state.FocusedOptionIndex < 0) {
		state.FocusedOptionIndex = len(state.Options) - 1;

	} else if (state.FocusedOptionIndex >= len(state.Options)) {
		state.FocusedOptionIndex = 0;
	}
}

// Either print a greyish placeholder-like "Search..." text at the current cursor pos or, if [state.SearchQuery]
// is not empty, erase the placeholder.
// No line breaks
func updateSearchInputPlaceholder() {
	placeholder := getSearchPlaceholder();
	searchQueryLength := utf8.RuneCountInString(state.SearchQuery);

	// print placeholder
	switch searchQueryLength {
	case 0:
		fmt.Print(placeholder);
		fmt.Print(ansi.CursorBackward(len(placeholder)));

	// delete placeholder
	case 1:
		// move cursor to index 1
		if state.SearchQueryCursorIndex == 0 {
			fmt.Printf("%v", ansi.CursorForward(1));
		}

		fmt.Printf("%v", ansi.DeleteCharacter(len(placeholder)));

		// move cursor back
		if state.SearchQueryCursorIndex == 0 {
			fmt.Print(ansi.CursorBackward(1));
		}
	}
}

func getSearchPlaceholder() string {
	return ansi.NewStyle().ForegroundColor(ansi.BrightBlack).Styled("Search...");
}

// Erase the menu assuming the cursor is currently at the search input below the menu.
//
// Also erase search prompt if enabled.
//
// Cursor will end up at search input line.
func clearMenu() {
	numLines := len(state.Options);

	// move up to first option
	fmt.Print(ansi.CursorPreviousLine(numLines));
	
	// delete menu
	fmt.Print(ansi.DeleteLine(numLines));
	// delete search input line as well
	fmt.Print(ansi.DeleteLine(1));
}

// Print [answer] next to [question] and make sure to bring the cursor back to the bottom of the menu afterwards.
func printAnswer(answer string) {
	linesToMoveUp := len(state.Options) + 1;

	// move up to question line
	fmt.Print(ansi.CursorPreviousLine(linesToMoveUp))

	// reprint question line, now with answer
	fmt.Printf("%v%v - %v", 
		ansi.EraseLine(2), 
		question,
		ansi.NewStyle().ForegroundColor(ansi.RGBColor{R: 100, G: 100, B: 255}).Styled(answer),
	);

	// move back to bottom most line
	fmt.Print(ansi.CursorNextLine(linesToMoveUp));
}