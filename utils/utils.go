package utils

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/eiannone/keyboard"
)

// Log [err] and exit if [err] is not [nil]. Do nothing otherwise
func ErrorLogExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Log [err] if not [nil] and prompt user to press Enter before exiting with code 1
func ErrorLogPromptExit(err error) {
	if err == nil {
		return;
	}
	
	slog.Error(err.Error());
	fmt.Print("Press Enter to exit...");
	fmt.Scanln();
	os.Exit(1);
}


// Execute [callback] repeatedly every [interval] until the ticker is closed or the program ends. Does not block.
func GoInterval(interval time.Duration, callback func(ticker *time.Ticker, time time.Time)) *time.Ticker {
	ticker := time.NewTicker(interval)

	go func() {
		for time := range ticker.C {
			if callback != nil {
				callback(ticker, time)
			}
		}
	}()

	return ticker
}

// Handles terminal key press calling [handler] or exiting on error.
// 
// Need to call [keyboard.Open()] followed by [defer keyboard.Close()] first.
//
// Notice that paste (ctrl + v) will trigger as many keystrokes as there are chars in the pasted text.
func HandleKeyPress(handler func (char rune, key keyboard.Key)) error {
	char, key, err := keyboard.GetKey();
	if err != nil {
		return err;
	}
	
	handler(char, key);
	
	return err;
}

// Behaves almost like [fmt.Scanln] but not imitating all modern ansi terminal behaviour. In return 
// the [callback] can be used to react to every keystroke.
//
// Terminates on "Enter" press.
//
// [callback]
// 
// Arg [char] the typed character, see [HandleKeyPress]
//
// Arg [key] the pressed key, see [HandleKeyPress]
//
// Arg [currentLine] the accumulated user input as printed to the terminal
//
// Arg [cursorIndex] the current cursorIndex (0-based)
//
// [return] error if [keyboard] has already been opened or the keypress handler returned an error
// func ScanlnRaw(callback func (char rune, key keyboard.Key, currentLine string, cursorIndex int)) error {
// 	err := keyboard.Open();
// 	if err != nil {
// 		return err;
// 	}
// 	defer keyboard.Close();

// 	cursorIndex := 0;
// 	currentLine := "";
	
// 	decrementCursorIndex := func() {
// 		if cursorIndex > 0 {
// 			cursorIndex--;
// 		}
// 	}
// 	incrementCursorIndex := func() {
// 		if cursorIndex < len(currentLine) {
// 			cursorIndex++;
// 		}
// 	}

// 	running := true;
// 	for running {
// 		err := HandleKeyPress(func(char rune, key keyboard.Key) {
// 			switch key {
// 			case keyboard.KeyBackspace:
// 				if (cursorIndex > 0) {
// 					currentLine = currentLine[0:cursorIndex - 1] + currentLine[cursorIndex:];

// 					fmt.Printf("%v%v", ansi.CursorBackward(1), ansi.DeleteCharacter(1));
// 					decrementCursorIndex();
// 				}

// 			case keyboard.KeyEnter:
// 				running = false;
				
// 			case keyboard.KeySpace:
// 				currentLine = currentLine[0:cursorIndex] + " " + currentLine[cursorIndex:];

// 				fmt.Printf("%v%v", ansi.InsertCharacter(1), " ");
// 				incrementCursorIndex();

// 			case keyboard.KeyArrowLeft:
// 				fmt.Print(ansi.CursorBackward(1));
// 				decrementCursorIndex();
				
// 			case keyboard.KeyArrowRight:
// 				if cursorIndex < len(currentLine) {
// 					fmt.Print(ansi.CursorForward(1));
// 					incrementCursorIndex();
// 				}
			
// 			case keyboard.KeyDelete:
// 				if len(currentLine) > 0 {
// 					currentLine = currentLine[0:cursorIndex] + currentLine[cursorIndex + 1:];
// 				}
// 				fmt.Print(ansi.DeleteCharacter(1));

// 			// will naturally only print keys with an utf-8 representation. Other keys simply get ignored.
// 			// Paste will trigger one keystroke per pasted char
// 			default:
// 				// splice char into inputC
// 				currentLine = currentLine[0:cursorIndex] + string(char) + currentLine[cursorIndex:];

// 				fmt.Printf("%v%v", ansi.InsertCharacter(1), string(char))
// 				incrementCursorIndex();
// 			}

// 			if callback != nil {
// 				callback(char, key, currentLine, cursorIndex);
// 			}
// 		});
// 		if err != nil {
// 			return err;
// 		}
// 	}

// 	return nil;
// }

// Asynchronously wait until interrupt signal is received, then run [handler] if not [nil].
// See also [os.Interrupt]
func GoHandleUserInterrupt(handler func()) {
	if handler == nil {
		return;
	}

	go func() {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt);
		defer cancel();

		// block until interrupt signal is received
		<- ctx.Done();

		handler();
	}();
}