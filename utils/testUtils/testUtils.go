package testUtils

import "testing"

// Fail the test with [t.Error] if [err] is not [nil]. Do nothing otherwise
func ErrorTestFail(err error, t *testing.T) {
	if err != nil {
		t.Error(err);
	}
}