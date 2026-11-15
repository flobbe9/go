// Helpers for cmdlets.
package helpers

import (
	"fmt"
	"jdkvm/internal/constants"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/flobbe9/go/utils/stringUtils"
	"github.com/flobbe9/go/utils/winutils"
)

// Indicates whether the current exe needs to be rerun as admin in order to modify the resource at [path].
// Takes into acccount whether the exe is already running as admin, which is why this function does ont definitively
// indicate whether [path] needs admin priviliges in general (e.g. located at /).
func IsRerunAsAdminToModifyPath(path string) (bool, error) {
	if stringUtils.IsBlank(path) {
		return false, fmt.Errorf("Arg 'path' cannot be blank");
	}

	if !winutils.IsSelfRunningAsAdmin() {
		return false, nil;
	}

	usrHome, err := os.UserHomeDir();
	if err != nil {
		return false, err;
	}

	return !stringUtils.StartsWith(path, usrHome), nil;
}

// Only call by cmdlet!
func SetupLog(cmdletUuid string, fullExeName string) error {
	if stringUtils.IsBlank(fullExeName) {
		return fmt.Errorf("Failed to setup script. Arg 'fullExeName' cannot be blank");
	}

	exeName := GetExeName(fullExeName);
	logFileName, err := GetCmdletLogFileName(cmdletUuid, exeName);
	if err != nil {
		return err;
	}	

	scriptLogFile, err := os.Create(logFileName);
	if err != nil {
		return err;
	}

	// have cmdlet logs write to dedicated file
	os.Stdout = scriptLogFile;
	os.Stderr = scriptLogFile;
	log.SetOutput(scriptLogFile);

	return nil;
}

func GetCmdletLogFileName(cmdletUuid string, exeName string) (string, error) {
	if stringUtils.IsBlank(cmdletUuid) {
		return "", fmt.Errorf("Failed to setup script. Arg 'cmdletUuid' cannot be blank");
	}
	if stringUtils.IsBlank(exeName) {
		return "", fmt.Errorf("Failed to setup script. Arg 'exeName' cannot be blank");
	}

	// create dedicated log file for cmdlet caller to read
	logFileName := fmt.Sprintf("%v/%v.%v.log", constants.TMP_DIR, exeName, cmdletUuid);
	slog.Debug("Creating script log file: " + logFileName);

	return logFileName, nil;
}

// TODO
func GetExeName(fullExeName string) string {
	var splitPaths []string;
	if strings.Contains(fullExeName, "\\") {
		splitPaths = strings.Split(fullExeName, "\\");

	} else {
		splitPaths = strings.Split(fullExeName, "/");
	}

	exeFileName := splitPaths[len(splitPaths) - 1];

	return exeFileName;
}