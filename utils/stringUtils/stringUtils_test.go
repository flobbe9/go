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
	assertNotModified := func(str string) {
		if unslashedStr := stringUtils.Unslash(str); unslashedStr != str {
			t.Errorf("Expected '%v' to equal '%v'", str, unslashedStr);
		}
	}

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
		assertNotModified(testCase);
	}
}

func TestUnslash_shouldModify(t *testing.T) {
	assertModified := func(str string) {
		if unslashedStr := stringUtils.Unslash(str); unslashedStr == str {
			t.Errorf("Expected '%v' not to equal '%v'", str, unslashedStr);
		}
	}

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
		assertModified(testCase);
	}
}

func TestIsBlank_shouldBeTrue(t *testing.T) {
	assertBlank := func(str string) {
		if !stringUtils.IsBlank(str) {
			t.Errorf("Expected '%v' to be blank", str);
		}
	}

	testCases := []string{
		"",
		" ",
		"     ",
	}

	for _, testCase := range testCases {
		assertBlank(testCase);
	}
}

func TestIsBlank_shouldBeFalse(t *testing.T) {
	assertNotBlank := func(str string) {
		if stringUtils.IsBlank(str) {
			t.Errorf("Expected '%v' not to be blank", str);
		}
	}

	testCases := []string{
		"d",
		"3",
		"\n",
		"\r\n",
		"asdfasfas",
		" d",
	}

	for _, testCase := range testCases {
		assertNotBlank(testCase);
	}
}
