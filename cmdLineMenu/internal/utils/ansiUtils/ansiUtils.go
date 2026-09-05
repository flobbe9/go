package ansiUtils

import (
	"github.com/charmbracelet/x/ansi"
	"github.com/flobbe9/go/cmdLineMenu/constants"
	"github.com/flobbe9/go/utils/stringUtils"
)

// Iterate each char in [str] treating ansi sequences as one char.
//
// [callback] Executed for every char:
//
// Arg [char] which is either an utf-8 rune of str or an ansi sequence.
//
// Arg [isAnsi] indicates whether [char] is an ansi sequence
//
// Arg [index] the 0-based iteration index counting any type of [char]. Initialized with -1
//
// Arg [ansiCharIndex] the 0-based iteration index for ansi [char]s. Initialized with -1
//
// Arg [nonAnsiCharIndex] the 0-based iteration index for non-ansi [char]s. Initialized with -1
//
// Returns [error] which, if not [nil], will break the loop and return immediatly
//
// Wont handle panics of callback.
//
// [return] the error returned from [callback] or [nil]
// TODO private
func IterateCharsConsiderAnsi(str string, callback func (char string, isAnsi bool, index, ansiCharIndex, nonAnsiCharIndex int) error) error {
	if (len(str) == 0 || callback == nil) {
		return nil;
	}

	var index, ansiCharIndex, nonAnsiCharIndex int = -1, -1, -1;
	var state byte
	p := ansi.NewParser();
	for len(str) > 0 {
		seq, _, n, newState := ansi.DecodeSequence(str, state, p);

		isAnsi := stringUtils.StartsWith(seq, constants.ANSI_ESCAPE_SEQ);
		index++;
		if !isAnsi {
			nonAnsiCharIndex++;

		} else {
			ansiCharIndex++;
		}

		err := callback(seq, isAnsi, index, ansiCharIndex, nonAnsiCharIndex);
		if err != nil {
			return err;
		}

		state = newState;
		str = str[n:];
	}

	return nil
}

// Indicates that there's at least one ansi sequence in [str]
func HasAnsi(str string) bool {
	return ansi.Strip(str) != str;
}

// [return] [true] if [a == b] after stripping all ansi chars from both args
func EqualsIgnoreAnsi(a, b string) bool {
	return ansi.Strip(a) == ansi.Strip(b);
}

// [return] [str] with characters colored with [color], even if [str] is blank
func Colored(str string, color ansi.Color) string {
	return ansi.NewStyle().ForegroundColor(color).Styled(str);
}