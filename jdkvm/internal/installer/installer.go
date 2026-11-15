// For installing the jdkvm tool.
package intaller;

// install(version)
	// create if not exists: /jdkvm
	// if /jdkvm/{version} already exists
		// prompt exit gracefully "already installed"
	// create /jdkvm/{version} at user roaming dir
	// create or point /jdkvm/jdkvm-current symlink to /jdkvm/{version}
		// add to path if not present yet
	// download matching version zip from github
		// version should be injected at build time
	// unzip to jdkvm/{version}


// TODO
	// utilities
	// as admin utils
		// mkdir
		// create / update symlink
			// accept location and symlinkTarget args
		// write file
		// delete file 
		// unzip (?)
		// updateSysEnvPath(path, doPrepend, mode = 'system' | 'user')
		// updateSysEnv(mode = 'system' | 'user')
			// create or update