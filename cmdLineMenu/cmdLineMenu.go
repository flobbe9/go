package cmdLineMenu

import (
	"fmt"
	"strings"

	internalModels "github.com/flobbe9/go/cmdLineMenu/internal/models"
	"github.com/flobbe9/go/cmdLineMenu/models"
	"github.com/flobbe9/go/utils"

	"github.com/charmbracelet/x/ansi"
	"github.com/eiannone/keyboard"
)

// Store the current state of the menu depending on user input.
var state = &internalModels.State{};

// options to choose from
var selectOptions []string;

// Passed to [Prompt()]
var opts *models.Options;


// Prompt user to answer [question] by selecting from [selectOptions] using Arrow-up / -down keys and submitting
// with Enter.
// 
// [return] the selected option from [selectOptions].
func Prompt(question string, _selectOptions []string, optsArg models.Options) (string, error) {
	if (len(_selectOptions) <= 0) {
		return "", fmt.Errorf("Specify at least one option");
	}

	// init global vars
	opts = &optsArg;
	selectOptions = _selectOptions;
	state.CurrentSelectionIndex = 0;

	// print question
	if (opts.IsShowSubmitHint) {
		fmt.Printf("%v (Submit with Enter)\n", question);

	} else {
		fmt.Println(question);
	}

	// first menu-print
	fmt.Print(formatMenu(selectOptions, state.CurrentSelectionIndex));

	// await user input
	err := handleUserInput();
	if err != nil {
		return "", err;
	}
	
	answer := selectOptions[state.CurrentSelectionIndex];
	
	// print answer
	if (opts.IsDisplayAnswer) {
		// move up to question line
		fmt.Printf("%v", ansi.CursorUp(len(selectOptions) + 1));
		// print answer next to question
		fmt.Printf("%v", ansi.EraseLine(2)); // erase hint too
		fmt.Printf("%v - %v\n", question, ansi.NewStyle().ForegroundColor(ansi.RGBColor{R: 100, G: 100, B: 255}).Styled(answer));
		// move back to bottom most line
		fmt.Print(ansi.CursorNextLine(len(selectOptions) + 1))
	}

	// clear select options
	if (opts.IsClearMenuOnSubmit) {
		clearMenu();
	}

	return answer, nil;
}

// Blocks until user has pressed Enter key. Increment / Decrement [state.CurrentSelectionIndex] depending on arrow-down/-up keys.
func handleUserInput() (error) {
	err := keyboard.Open();
	if err != nil {
		return err;
	}
	defer keyboard.Close();
	
	var didSelect bool;
	for !didSelect {
		err := utils.HandleKeyPress(func(char rune, key keyboard.Key) {
			switch key {
			case keyboard.KeyEnter:
				didSelect = true;

			case keyboard.KeyArrowDown:
				updateCurrentSelection(false);
				rerender(formatMenu(selectOptions, state.CurrentSelectionIndex), len(selectOptions));

			case keyboard.KeyArrowUp:
				updateCurrentSelection(true);
				rerender(formatMenu(selectOptions, state.CurrentSelectionIndex), len(selectOptions));
			}
		});

		if err != nil {
			return err;
		}
	}

	return nil;
}

// Override existing stdout by moving the current terminal's cursor up by [numLines] and then
// printing [content] (not printing a line break at the end).
//
// [numLines] use 0 or 1 to stay at current line
func rerender(content string, numLines int) {
	cursorUpAnsi := ansi.CursorPreviousLine(numLines);
	if (numLines == 0) {
		// 0 would move up a line for some reason
		cursorUpAnsi = "";
	}

	// move up
	fmt.Printf("%v", cursorUpAnsi);
	// delete following lines in case new content is shorter than previous content
	fmt.Print(ansi.DeleteLine(numLines));

	fmt.Print(content);
}

// [return] formatted menu line including a line break and possibly underlined if [isSelected == true]
func formatMenuLine(lineContent string, isSelected bool) string {
	style := ansi.NewStyle().Underline(isSelected);

	return fmt.Sprintf("> %v\n", style.Styled(lineContent));
}

// [return] the whole menu consisting of all select options. End on a new line
func formatMenu(options []string, selectionIndex int) string {
	var menuStr strings.Builder;

	for i, option := range options {
		menuStr.WriteString(formatMenuLine(option, i == selectionIndex)); 
	}

	return menuStr.String();
}

// Increment / Decrement [state.CurrentSelectionIndex] by one making sure it loops to the start / end if out of bounds
func updateCurrentSelection(isDecrease bool) {
	if (isDecrease) {
		state.CurrentSelectionIndex--;
	} else {
		state.CurrentSelectionIndex++;
	}
	
	// loop user selection
	if (state.CurrentSelectionIndex < 0) {
		state.CurrentSelectionIndex = len(selectOptions) - 1;

	} else if (state.CurrentSelectionIndex >= len(selectOptions)) {
		state.CurrentSelectionIndex = 0;
	}
}

// Erase the menu assuming the cursor is currently at the last menu option.
func clearMenu() {
	numLines := len(selectOptions);

	for range numLines {
		fmt.Printf("%v%v", ansi.EraseLine(2), ansi.CursorPreviousLine(1));
	}

	// clear last line
	fmt.Printf("%v", ansi.EraseLine(2));
}