package fileUtils

import (
	"errors"
	"io/fs"
	"os"
	"time"
)

func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}


// Generates a filename-safe timestamp for [time.Now()].
// TODO tests
func Timestamp() string {
	return time.Now().Local().Format("2006-01-02T15-04-05.000");
}