package models

import (
	"bufio"
	"fmt"
	cmdletHelpers "jdkvm/internal/cmdlets/helpers"
	"log/slog"
	"os"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/flobbe9/go/utils/fileUtils"
	"github.com/flobbe9/go/utils/stringUtils"
	"golang.org/x/sys/windows"
)

// Represents a runnable go script containing a main function.
type Cmdlet struct {
	// Optional. Verify the successful outcome of the script, e.g. by checking if a file was created, and return [nil] in that case.
	Verify func() error;

	// E.g. "test.exe", relative to "go.mod". It may be assumed that all exe files will be located inside the same directory.
	File string;
	// TODO max tries and / or interval

	// Optional args passed to [File] and separated with " ". These don't include the script uuid as it's 
	// appended by [this.Run()].
	Args []string;

	// TODO
	Uuid string;
}

// TODO use 
func New() Cmdlet {

	return Cmdlet{Uuid: fileUtils.Timestamp()};
} 

func (this *Cmdlet) Run(windowMode int32) error {
	if stringUtils.IsBlank(this.File) {
		return fmt.Errorf("Failed to run script. 'this.File' cannot be blank");
	}

	verbPtr, _ := syscall.UTF16PtrFromString("");

	exePtr, _ := syscall.UTF16PtrFromString(this.File);

	cwd, _ := os.Getwd();
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd);

	args := slices.Concat([]string{this.Uuid}, this.Args);
	argsPtr, _ := syscall.UTF16PtrFromString(strings.Join(args, " "));

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

// TODO
func (this *Cmdlet) AwaitVerification() error {
	if this.Verify == nil {
		return nil;
	}

	var tries int;
	for {
		tries++;

		slog.Debug(fmt.Sprintf("Verifying go script '%v'...", this.File))

		err := this.Verify();
		if err == nil {
			break;
		}

		time.Sleep(500 * time.Millisecond);

		if tries > 10 {
			return err;
		}
	}

	return nil;
}

func (this *Cmdlet) LogCmdletLogs() error {
	cmdletLogFileName, err := this.getCmdletLogFileName();
	if err != nil {
		return err;
	}

	// TODO continue here, does not find log file
	cmdletLogFile, err := os.OpenFile(cmdletLogFileName, os.O_RDONLY, os.ModeAppend); // TODO permission?
	if err != nil {
		return err;
	}

	reader := bufio.NewReader(cmdletLogFile);
	for line, _, err := reader.ReadLine(); err != nil; {
		fmt.Println(string(line));
	}

	return nil;
}

func (this *Cmdlet) CleanUp() error {
	cmdletLogFileName, err := this.getCmdletLogFileName();
	if err != nil {
		return err;
	}

	err = os.Remove(cmdletLogFileName);
	// TODO dont return err if did not exist
	return err;
}

func (this *Cmdlet) getCmdletLogFileName() (string, error) {
	return cmdletHelpers.GetCmdletLogFileName(this.Uuid, cmdletHelpers.GetExeName(this.File))
}

