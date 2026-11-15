package main

import (
	cmdletHelpers "jdkvm/internal/cmdlets/helpers"
	"log/slog"
	"os"

	"github.com/flobbe9/go/utils"
	"github.com/flobbe9/go/utils/winutils"
	"golang.org/x/sys/windows"
)

// TODO handle exists
func main() {

	// TODO validate args
		// or consider using cobra
	slog.SetLogLoggerLevel(slog.LevelDebug);

	fullExeName := os.Args[0];
	scriptUuid := os.Args[1];
	targetDir := os.Args[2];


	// make sure logs can be found by script caller
	err := cmdletHelpers.SetupLog(scriptUuid, fullExeName);
	utils.ErrorLogExit(err);

	// possibly rerun as admin
	// TODO test this
	isRerunAsAdmin, err := cmdletHelpers.IsRerunAsAdminToModifyPath(targetDir);
	utils.ErrorLogExit(err);
	if isRerunAsAdmin {
		err = winutils.RerunSelfAsAdmin(windows.SW_NORMAL);
		utils.ErrorLogExit(err);
	}

	err = os.Mkdir(targetDir, os.ModeDir); // TODO permissions?
	utils.ErrorLogExit(err);
}
