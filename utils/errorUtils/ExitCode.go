package errorUtils

// General exit code trying to follow linux standards.
type ExitCode int;

const (
	// Exit code 0: Exit without error
	EXIT_OK ExitCode = 0;
	// Exit code 1: Exit with any general, miscellaneous error
	EXIT_ERR ExitCode = 1;
	// Exit code 130: Ctrl + C
	EXIT_USER_INTERRUPT ExitCode = 130;
)
