package testUtils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func GetRootDir() (string, error) {
	info, err := os.Stat(fmt.Sprintf("..%v..%v..%v", string(filepath.Separator), string(filepath.Separator), string(filepath.Separator)));

	return info.Name(), err;
}

func AssertNoError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("Expected 'err' to be nil but was: %v", err);
	}
}