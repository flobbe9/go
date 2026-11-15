package winutils

import (
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// Restarts this go executable with admin privileges. This will prompt the user to grant those.
//
// [windowMode] see [windows] package "windows.SW_...". Determines how the terminal for the rerun is opened
//
// Returns an error if [windows.ShellExecute] threw one
func RerunSelfAsAdmin(windowMode int32) error {
	verbPtr, _ := syscall.UTF16PtrFromString("runas");

	exe, _ := os.Executable(); // the executable of this go program
	exePtr, _ := syscall.UTF16PtrFromString(exe);

	cwd, _ := os.Getwd();
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd);

	args := strings.Join(os.Args[1:], " ") // args passed to this go program
	argsPtr, _ := syscall.UTF16PtrFromString(args);

	err := windows.ShellExecute(
		0,
		verbPtr,
		exePtr,
		argsPtr,
		cwdPtr,
		windowMode,
	);

	return err;
}

// Indicates whether this go program is beeing run with admin privileges.
func IsSelfRunningAsAdmin() bool {
	// try to open the main hard drive which should exist on any windows machine and be named exactly the same. 
	// Therefore an error can only mean missing admin priviledges
    _, err := os.Open("\\\\.\\PHYSICALDRIVE0"); 
    return err == nil;
}