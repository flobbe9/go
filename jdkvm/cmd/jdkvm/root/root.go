package root

import (
	"github.com/spf13/cobra"
)

// expect absolute path arg
// allow no arg, display root in that case

// validate arg
	// is abs path
		// files.IsAbs
	// is dir
		// filepath.Lstat.isdir
// create if not exists
// store config


func Run(cmd *cobra.Command, args []string) {
}