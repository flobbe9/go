package main

import (
	"fmt"
	"jdkvm/cmd/jdkvm/install"

	"github.com/flobbe9/go/utils"
	"github.com/spf13/cobra"
)

// -------- Injected vars --------
// the current jdkvm installation version
var version string = "SNAPSHOT";

var jdkvmCmd = &cobra.Command{Use: "jdkvm", Run: RunJdkvm};
var installCmd = &cobra.Command{
	Use:   "install",
	Aliases: []string{"i"},
	Args: cobra.ExactArgs(1), // expect jdk-version arg
	Short: "TODO",
	Run: install.Run,
}

const versionFlagName = "version";

func main() {
	jdkvmCmd.AddCommand(installCmd);
	jdkvmCmd.Flags().Bool(versionFlagName, false, "Display the version of the current jdkvm installation.");

	jdkvmCmd.Execute();
}

// Only runs if no other command is specified. Wont run any logic, just display help if no other flag is present
func RunJdkvm(cmd *cobra.Command, args []string) {
	isShowVersion, err := cmd.Flags().GetBool(versionFlagName);
	utils.ErrorLogPromptExit(err);

	if isShowVersion {
		fmt.Println(version);

	// case: no flags
	} else {
		cmd.Help();
	}
}

func GetJdkvmVersion() string {
	return version;
}