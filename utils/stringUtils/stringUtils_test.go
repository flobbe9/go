package stringUtils_test

import (
	"testing"

	"github.com/flobbe9/go/utils/stringUtils"
)


func TestStartsWith_true(t *testing.T) {
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

func TestStartsWith_false(t *testing.T) {
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

func TestEndsWith_true(t *testing.T) {
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

func TestEndsWith_false(t *testing.T) {
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

func TestUnslash_shouldNotModify(t *testing.T) {
	testCases := []string{
		"",
		" ",
		"asdf",
		"asdf/asdf",
		"asdf/ ",
		" /asdf",
		" / ",
	}

	for _, testCase := range testCases {
		if unslashedStr := stringUtils.Unslash(testCase); unslashedStr != testCase {
			t.Errorf("Expected '%v' to equal '%v'", testCase, unslashedStr);
		}
	}
}

func TestUnslash_shouldModify(t *testing.T) {
	testCases := []string{
		"/",
		"/ ",
		"/asdf",
		"asdf/",
		" /",
		"///asdfasdf///",
		"////",
	}

	for _, testCase := range testCases {
		if unslashedStr := stringUtils.Unslash(testCase); unslashedStr == testCase {
			t.Errorf("Expected '%v' not to equal '%v'", testCase, unslashedStr);
		}
	}
}

func TestIsBlank_shouldBeTrue(t *testing.T) {
	testCases := []string{
		"",
		" ",
		"     ",
	}

	for _, testCase := range testCases {
		if !stringUtils.IsBlank(testCase) {
			t.Errorf("Expected '%v' to be blank", testCase);
		};
	}
}

func TestIsBlank_shouldBeFalse(t *testing.T) {
	testCases := []string{
		"d",
		"3",
		"\n",
		"\r\n",
		"asdfasfas",
		" d",
	}

	for _, testCase := range testCases {
		if stringUtils.IsBlank(testCase) {
			t.Errorf("Expected '%v' not to be blank", testCase);
		}
	}
}
