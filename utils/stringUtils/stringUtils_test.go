package stringUtils_test

import (
	"testing"

	"github.com/flobbe9/go/utils/stringUtils"
)


func TestStartsWith_True(t *testing.T) {
	assertTrue := func(str string, subStrs []string) {
		for _, subStr := range subStrs {
			if (!stringUtils.StartsWith(str, subStr)) {
				t.Errorf("Expected '%s' to start with '%v'", str, subStr);
			}
		}
	}

	assertTrue("test", []string{
		"",
		"t",
		"te",
		"test",
	});

	assertTrue("", []string{
		"",
	});

	assertTrue(" ", []string{
		"",
		" ",
	});
}

func TestStartsWith_False(t *testing.T) {
	assertFalse := func(str string, subStrs []string) {
		for _, subStr := range subStrs {
			if (stringUtils.StartsWith(str, subStr)) {
				t.Errorf("Expected '%s' to start with '%v'", str, subStr);
			}
		}
	}

	assertFalse("test", []string{
		"est",
		" test",
		"testing",
		"T",
		"xyz",
		" ",
	});
}

func TestEndsWith_True(t *testing.T) {
	assertTrue := func(str string, subStrs []string) {
		for _, subStr := range subStrs {
			if (!stringUtils.EndsWith(str, subStr)) {
				t.Errorf("Expected '%s' to end with '%v'", str, subStr);
			}
		}
	}

	assertTrue("test", []string{
		"",
		"t",
		"st",
		"test",
	});

	assertTrue("", []string{
		"",
	});

	assertTrue(" ", []string{
		"",
		" ",
	});
}

func TestEndsWith_False(t *testing.T) {
	assertFalse := func(str string, subStrs []string) {
		for _, subStr := range subStrs {
			if (stringUtils.EndsWith(str, subStr)) {
				t.Errorf("Expected '%s' to end with '%v'", str, subStr);
			}
		}
	}

	assertFalse("test", []string{
		"tes",
		"tes ",
		"testing",
		"T",
		"xyz",
		" ",
	});
}