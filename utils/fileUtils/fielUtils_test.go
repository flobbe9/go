package fileUtils_test

import (
	"path/filepath"
	"testing"

	"github.com/flobbe9/go/utils/fileUtils"
	"github.com/flobbe9/go/utils/internal/test/fileUtils/constants"
	"github.com/flobbe9/go/utils/internal/test/testUtils"
)

// false
// dir not exists
// assert error nil
// some other error
// also assert that error is not nil
// true
// dir exists
// assert error nil

func TestDirExists_shouldBeTrue(t *testing.T) {

	root, err := testUtils.GetRootDir();
	testUtils.AssertNoError(err, t);

	testResourcesDir := root + constants.TEST_RESOURCES_DIR;
	exists, err := fileUtils.DirExists(testResourcesDir);
	testUtils.AssertNoError(err, t);
	
	if !exists {
		t.Errorf("Expected '%v' to exist", testResourcesDir);
	}
}

func TestDirExists_shouldBeFalse(t *testing.T) {

	root, err := testUtils.GetRootDir();
	testUtils.AssertNoError(err, t);

	testResourcesDir := root + constants.TEST_RESOURCES_DIR;
	nonExistentDir := testResourcesDir + string(filepath.Separator) + "anyasdf";

	exists, err := fileUtils.DirExists(nonExistentDir);
	testUtils.AssertNoError(err, t);
	if exists {
		t.Errorf("Expected '%v' not to exist", testResourcesDir);
	}
}
