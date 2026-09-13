package cmdLineMenu

import (
	"fmt"
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
	rerender(formatMenu(state.Options, state.FocusedOptionIndex), 0);
	if (config.IsOptionSearchEnabled) {
		updateSearchInputPlaceholder();
		// TODO maybe make a input-like border
	}

	// await user input
	err := handleUserInput();
	if err != nil {
		return "", err;
	}
	
	answer := ansi.Strip(state.Options[state.FocusedOptionIndex]);

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

			if (rawStdin.Buff[0] == key.Enter) {
				if config.IsOptionSearchEnabled && !isNoSearchResults() {
					// clear search
					fmt.Print(ansi.DeleteLine(1));
				}
				
			} else if rawStdin.Buff[0] == key.Backspace || rawStdin.Buff[0] == key.CtrlBackspace {
				// TODO repetitive
					// try search input on top again
				if config.IsOptionSearchEnabled {
					prevOptionsLen := len(state.Options);
					updateSearchInputPlaceholder();
					searchOptionsAndUpdateStates();
					rerender(formatMenu(state.Options, state.FocusedOptionIndex), prevOptionsLen);
				}
				
			} else if rawStdin.IsUtf8() {
				if config.IsOptionSearchEnabled {
					prevOptionsLen := len(state.Options);
					updateSearchInputPlaceholder();
					searchOptionsAndUpdateStates();
					rerender(formatMenu(state.Options, state.FocusedOptionIndex), prevOptionsLen);
				}
			}

			switch rawStdin.GetSignificantAnsiKey() {
			case ansiKey.ArrowDown:
				if len(state.Options) > 0 {
					updateFocusedOptionIndex(false);
					rerender(formatMenu(state.Options, state.FocusedOptionIndex), len(state.Options));
				}

			case ansiKey.ArrowUp:
				if len(state.Options) > 0 {
					updateFocusedOptionIndex(true);
					rerender(formatMenu(state.Options, state.FocusedOptionIndex), len(state.Options));
				}

			case ansiKey.Delete:
				if config.IsOptionSearchEnabled {
					prevOptionsLen := len(state.Options);
					updateSearchInputPlaceholder();
					searchOptionsAndUpdateStates();
					rerender(formatMenu(state.Options, state.FocusedOptionIndex), prevOptionsLen);
				}
			}
		})

		// don't exit if nothing focused
		if !isNoSearchResults() || err != nil {
			return err;
		}
	}
}

// Indicates that search has produced no results.
//
// [return] [true] if [state.FocusedOptionIndex == -1]
func isNoSearchResults() bool {
	return state.FocusedOptionIndex == -1;
}

// TODO
func searchOptionsAndUpdateStates() {
	prevSearchResults := slices.Clone(state.Options);

	if stringUtils.IsBlank(state.SearchQuery) {
		didSearchResultsChange := !ansiUtils.EqualsSlicesIgnoreAnsi(prevSearchResults, initialOptions);
		// don't modify menu if no changes
		if !didSearchResultsChange {
			return;
		}
		state.Options = initialOptions;
		state.FocusedOptionIndex = 0;
		return;
	}

	searchResults := search.SearchAndHighlightOptions(initialOptions, state.SearchQuery);

	if len(searchResults) == 0 {
		// TODO should not print '<'
		state.Options = []string{fmt.Sprintf("No results for search '%v'", state.SearchQuery)};
		state.FocusedOptionIndex = -1; 

	} else {
		// TODO will only underline until first highlighted substring
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

// [return] formatted menu line including a line break and possibly underlined if [isSelected == true]
func formatMenuLine(lineContent string, isSelected bool) string {
	return fmt.Sprintf("> %v\n", ansi.NewStyle().Underline(isSelected).Styled(lineContent));
}

// [focusedOptionIndex] the index of the option currently focused
//
// [return] the whole menu consisting of all select options. End on a new line
// TODO remove args?
func formatMenu(options []string, focusedOptionIndex int) []string {
	return sliceUtils.Map(options, func(option string, i int) string {
		return formatMenuLine(option, i == focusedOptionIndex);
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

// Erase the menu assuming the cursor is currently at the last menu option.
//
// Also erase search prompt if enabled.
func clearMenu() {
	numLines := len(state.Options);
	if config.IsOptionSearchEnabled {
		// TODO
		// numLines++; // also erase search line
	}

	for range numLines {
		fmt.Printf("%v%v", ansi.EraseLine(2), ansi.CursorPreviousLine(1));
	}

	// clear last line
	fmt.Printf("%v", ansi.EraseLine(2));
}

// Print [answer] next to [question] and make sure to bring the cursor back to the bottom of the menu afterwards.
func printAnswer(answer string) {
	linesToMoveUp := len(state.Options);
	linesToMoveUp++;
	if (config.IsOptionSearchEnabled) {
		// linesToMoveUp++;
	}

	// move up to question line
	fmt.Printf("%v%v", ansi.CursorBackward(len(state.SearchQuery)), ansi.CursorUp(linesToMoveUp));

	fmt.Printf("%v", ansi.EraseLine(2)); // erase whole question line including hint
	// reprint question line, now with answer
	fmt.Printf("%v - %v\n", question, ansi.NewStyle().ForegroundColor(ansi.RGBColor{R: 100, G: 100, B: 255}).Styled(answer));

	// move back to bottom most line
	fmt.Print(ansi.CursorNextLine(linesToMoveUp));
}