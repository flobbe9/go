package utils

import (
	"fmt"
	"log"
	"log/slog"
	"os"
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
	fmt.Scanln("Press Enter to exit...");
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
func HandleKeyPress(handler func (char rune, key keyboard.Key)) error {
	char, key, err := keyboard.GetKey();
	if err != nil {
		return err;
	}
	
	handler(char, key);
	
	return err;
}