package install

import (
	"fmt"

	"github.com/spf13/cobra"
)

// [args] starts counting after "install":
//  [args[0] == jdkVersion]
// This command is configured to panic if no arg is specified
func Run(cmd *cobra.Command, args []string) {

	jdkVersionInput := args[0];
	fmt.Println("install", jdkVersionInput)

}

// parse version
	// should return a fully qualified release_name
	// possible inputs
		// unsigned integer for major version, e.g. 25, 17, 15
		// strings for whole version, e.g. jdk-21.8.0.1+1 
	// if is int
		// return latest release for major: /v3/assets/latest/{feature_version}/{jvm_impl}
		// if not exists
			// error
// check if version already installed
// check if version exists
	// fetch from /v3/assets/latest/{feature_version}/{jvm_impl}
	// if not exists
		// err
// retrieve os
	// runtime.GOOS
// retrieve cpu arch
// if is major
	// download latest: /v3/binary/latest/{feature_version}/{release_type}/{os}/{arch}/{image_type}/{jvm_impl}/{heap_size}/{vendor}
// else
	// download specific: /v3/binary/version/{release_name}/{os}/{arch}/{image_type}/{jvm_impl}/{heap_size}/{vendor}
// extract to /jdkvm/{jdk-version}/...
// prompt "run 'jdkvm use {version}' in order to use the installed version"