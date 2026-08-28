package constants

import (
	"fmt"
	"path/filepath"
)

// Relative to "utils" dir. Use [testUtils.getRootDir() + constants.TEST_RESOURCES_DIR] to get a fully qualified path.
var TEST_RESOURCES_DIR = fmt.Sprintf("%vinternal%vtest%vfileUtils%vresources", string(filepath.Separator), string(filepath.Separator), string(filepath.Separator), string(filepath.Separator));