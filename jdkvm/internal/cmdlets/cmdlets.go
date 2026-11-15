// For calling cmdlets only. See also [helpers] package
package cmdlets

import (
	"jdkvm/internal/models"
	"os"

	"golang.org/x/sys/windows"
)

func MkdDir(path string) error {

	script := models.New();
	script.Verify = func() error {
		_, err := os.Lstat(path)
		return err
	};
	script.File = "mkdir.exe";
	script.Args = []string{path};

	err := script.Run(windows.SW_NORMAL)
	if err != nil {
		return err;
	}
	
	err = script.AwaitVerification();
	if err != nil {
		return err;
	}

	err = script.LogCmdletLogs();
	if err != nil {
		return err;
	}
	
	// err = script.CleanUp();

	return err
}

