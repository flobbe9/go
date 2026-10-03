package errorUtils

import (
	"fmt"
)

// [error] implementation with an exit code. See also [ExitCode].
type ErrorExitCode struct {
	code ExitCode;

	msg string;
}

// Unmodified error message
func (this *ErrorExitCode) Msg() string {
	return this.msg;
}

func (this ErrorExitCode) Code() ExitCode {
	return this.code;
}

// Indicates that [this] contains an error, that is [this.Code > 0]
func (this ErrorExitCode) IsError() bool {
	return this.code > 0;
}

func (this ErrorExitCode) String() string {
	return fmt.Sprintf("%v (Exit code %v)", this.msg, this.code);
}

func (this ErrorExitCode) Error() string {
	return this.String();
}

// Create new instance where [args] is inserted into [msg]
func NewErrorExitCodef(code ExitCode, msg string, args ...any) *ErrorExitCode {
	return &ErrorExitCode{
		code: code,
		msg: fmt.Sprintf(msg, args...),
	}
}

// Create new instance
func NewErrorExitCode(code ExitCode, msg string) *ErrorExitCode {
	return NewErrorExitCodef(code, "%s", msg);
}
