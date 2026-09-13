package scanner

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"unicode/utf8"

	"github.com/charmbracelet/x/ansi"
	"github.com/flobbe9/go/cmdLineMenu/constants/ansiKey"
	"github.com/flobbe9/go/cmdLineMenu/constants/key"
	"github.com/flobbe9/go/cmdLineMenu/models"
	"golang.org/x/term"
)

// Put the current terminal into raw mode and read up to 3 bytes (but at least 1) from [os.Stdin]. Tries to guess the format of the input
// using the first byte, e.g. utf-8 (needing 2 bytes) or ansi (needing at least 3 bytes) sothat reading the next byte
// should never require a second keystroke. This restriction means that ansi sequences of more than 3 bytes cannot be read by this method.
//
// [switchOnRawMode] pass [false] if the terminal is already in raw mode, e.g. when using this function inside a loop.
//
// [return] the read bytes or an error if the terminal could not be switched into raw mode, or Ctrl + C was pressed
func ScanRaw(switchOnRawMode bool) ([]byte, error) {
	if switchOnRawMode {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()));
		if err != nil {
			return nil, err;
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState);
	}
	
	rawStdin := models.RawStdin{Buff: []byte{}};
	reader := bufio.NewReader(os.Stdin);
	
	b, err := reader.ReadByte();
	if err != nil {
		return nil, err;
	}
	
	// raw mode requires manual exit handling 
	if b == key.CtrlC {
		return nil, fmt.Errorf("User interrupt");
	} 

	rawStdin.Buff = append(rawStdin.Buff, b);

	// might be utf8
	if !rawStdin.IsRegularAscii() {
		b2, err := reader.ReadByte();
		if err == nil {
			rawStdin.Buff = append(rawStdin.Buff, b2);
		}
	}
	
	// might be ansi
	if rawStdin.LooksLikeAnsi() {
		b2, err := reader.ReadByte();
		if err != nil {
			return nil, err;
		}
		rawStdin.Buff = append(rawStdin.Buff, b2);
		
		b3, err := reader.ReadByte();
		if err != nil {
			return nil, err;
		}
		rawStdin.Buff = append(rawStdin.Buff, b3);
	}

	return rawStdin.Buff, nil;
}

// Uses [ScanRaw] for repeatedly reading bytes from [os.Stdin], stopping when Enter key is pressed. This scan function
// tries to replace a regular scanner, reducing some functionality but allowing to react to every keystroke. 
// Notice that keystrokes are limited to 3 bytes, see [ScanRaw].
//
// [callback]
//
// Arg [rawStdin] wrapper around the read bytes providing functionality to interpret the last keystroke
// 
// Arg [line] the accumulated printable chars typed so far, as visible to the user
//
// Arg [cursorIndex] 0-based index of the current terminal cursor position
//
// [return] the read bytes or an error if the terminal could not be switched into raw mode, or Ctrl + C was pressed
func ScanlnRaw(callback func (rawStdin models.RawStdin, line string, cursorIndex int)) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()));
    if err != nil {
		return err;
    }
    defer term.Restore(int(os.Stdin.Fd()), oldState);

	cursorIndex := 0;
	line := []rune{}; // make sure that utf-8 chars are represented as 1 element only
	
	decrementCursorIndex := func() {
		if cursorIndex > 0 {
			cursorIndex--;
		}
	}
	incrementCursorIndex := func() {
		if cursorIndex < len(line) {
			cursorIndex++;
		}
	}

	for {
		buff, err := ScanRaw(false);
		if err != nil {
			return err;
		}
		rawStdin := models.RawStdin{Buff: buff};

		// UTF-8
		if rawStdin.Buff[0] == key.Backspace {
			if (cursorIndex > 0) {
				line = slices.Concat(line[0:cursorIndex - 1], line[cursorIndex:]);
	
				fmt.Printf("%v%v", ansi.CursorBackward(1), ansi.DeleteCharacter(1));
				decrementCursorIndex();
			}

		} else if rawStdin.Buff[0] == key.Enter {
			if callback != nil {
				callback(rawStdin, string(line), cursorIndex);
			}
			break;

		} else if rawStdin.Buff[0] == key.CtrlBackspace {
			// TODO have this behave like normal ansi terminal
			line = []rune{};
			fmt.Print(ansi.DeleteLine(1));
			cursorIndex = 0;

		} else if rawStdin.IsUtf8() {
			char := rawStdin.DecodeUtf8Rune();
			if char != utf8.RuneError {
				// splice char into current position
				line = slices.Concat(line[0:cursorIndex], []rune{char}, line[cursorIndex:]);

				fmt.Printf("%v%v", ansi.InsertCharacter(1), string(char));
				incrementCursorIndex();
			}
		}

		// ANSI
		switch rawStdin.GetSignificantAnsiKey() {
		case ansiKey.Delete:
			if cursorIndex < len(line) {
				line = slices.Concat(line[0:cursorIndex], line[cursorIndex + 1:]);
				fmt.Print(ansi.DeleteCharacter(1));
			}

		case ansiKey.ArrowRight:
			if cursorIndex < len(line) {
				fmt.Print(ansi.CursorForward(1));
				incrementCursorIndex();
			}
			
		case ansiKey.ArrowLeft: 
			fmt.Print(ansi.CursorBackward(1));
			decrementCursorIndex();
		}

		if callback != nil {
			callback(rawStdin, string(line), cursorIndex);
		}
	}

	return nil;
}