package cmdLineMenu

import (
	"fmt"
	"strings"

	internalModels "github.com/flobbe9/go/cmdLineMenu/internal/models"
	models "github.com/flobbe9/go/cmdLineMenu/models"
	utils "github.com/flobbe9/go/utils"

	"github.com/charmbracelet/x/ansi"
	"github.com/eiannone/keyboard"
)

// Store the current state of the menu depending on user input.
var state = &internalModels.State{};

// options to choose from
var SelectOptions []string;

// Passed to [Prompt()]
var opts *models.Options;


// Prompt user to answer [question] by selecting from [selectOptions] using Arrow-up / -down keys and submitting
// with Enter.
// 
// [return] the selected option from [selectOptions].
func Prompt(question string, selectOptions []string, optsArg models.Options) (string, error) {
	if (len(selectOptions) <= 0) {
		return "", fmt.Errorf("Specify at least one option");
	}

	// init global vars
	opts = &optsArg;
	SelectOptions = selectOptions;
	state.CurrentSelectionIndex = 0;

	// question
	if (opts.IsShowSubmitHint) {
		fmt.Printf("%v (Submit with Enter)\n", question);

	} else {
		fmt.Println(question);
	}

	// first menu-print
	fmt.Print(formatMenu(SelectOptions, state.CurrentSelectionIndex));

	// user input
	err := handleUserInput();
	if err != nil {
		return "", err;
	}
	
	// clean
	if (opts.IsClearMenuOnSubmit) {
		clearMenu();
	}

	answer := SelectOptions[state.CurrentSelectionIndex];
	
	// print answer
	if (opts.IsDisplayAnswer) {
		fmt.Printf("%v", ansi.CursorUp(1));
		fmt.Printf("%v", ansi.EraseLine(2)); // erase hint too
		fmt.Printf("%v - %v\n", question, ansi.NewStyle().ForegroundColor(ansi.RGBColor{R: 100, G: 100, B: 255}).Styled(answer));
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
	for {
		if (didSelect) {
			break;
		}

		err := utils.HandleKeyPress(func(char rune, key keyboard.Key) {
			switch key {
			case keyboard.KeyEnter:
				didSelect = true;

			case keyboard.KeyArrowDown:
				updateCurrentSelection(false);
				rerender(formatMenu(SelectOptions, state.CurrentSelectionIndex), len(SelectOptions));

			case keyboard.KeyArrowUp:
				updateCurrentSelection(true);
				rerender(formatMenu(SelectOptions, state.CurrentSelectionIndex), len(SelectOptions));
			}
		});

		if err != nil {
			return err;
		}
	}

	return nil;
}

// Override existing stdout by moving the current terminal's cursor up by [numLines] and all the way backwards, then
// print [content] (not printing a line break at the end).
//
// [numLines] use 0 or 1 to stay at current line?
func rerender(content string, numLines int) {
	cursorUpAnsi := ansi.CursorPreviousLine(numLines);
	// ansi would move up a line even with 0 for some reason
	if (numLines == 0) {
		cursorUpAnsi = "";
	}

	fmt.Printf("%v%v%v",
		cursorUpAnsi, 
		ansi.EraseLine(2), 
		content,
	);
}

// [return] formatted menu line including a line break and possibly underlined if [isSelected == true]
func formatMenuLine(lineContent string, isSelected bool) string {
	style := ansi.NewStyle().Underline(isSelected);

	return fmt.Sprintf("> %v\n", style.Styled(lineContent));
}

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
		state.CurrentSelectionIndex = len(SelectOptions) - 1;

	} else if (state.CurrentSelectionIndex >= len(SelectOptions)) {
		state.CurrentSelectionIndex = 0;
	}
}

// Erase the menu assuming the cursor is currently at the last menu option.
func clearMenu() {
	numLines := len(SelectOptions);

	for range numLines {
		fmt.Printf("%v%v", ansi.EraseLine(2), ansi.CursorPreviousLine(1));
	}

	// clear last line
	fmt.Printf("%v", ansi.EraseLine(2));
}