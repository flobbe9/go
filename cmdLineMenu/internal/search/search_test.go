package search

import (
	"slices"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/flobbe9/go/cmdLineMenu/constants"
	"github.com/flobbe9/go/cmdLineMenu/internal/models"
	"github.com/flobbe9/go/cmdLineMenu/internal/utils/ansiUtils"
)

func TestFindMatchingOptionSubstrings_shouldReturnEmptySlice(t *testing.T) {
	options := []string{
		"something",
		"",
		" ",
	};
	searchQueries := []string {
		"",
		" ",
		"xyz",
		";$\n",
	}

	for _, option := range options {
		for _, searchQuery := range searchQueries {
			substrs := findMatchingOptionSubstrings(option, searchQuery);
			if (len(substrs) != 0) {
				t.Errorf("Expected substrs not to appear in option '%v' but got %v", option, substrs);
			}
		}
	}
}

func TestFindMatchingOptionSubstrings_shouldReturnSubstrs(t *testing.T) {
	var option string;
	var searchQuery string;
	expectedSubstrs := []models.OptionSubstring{};

	assertEquals := func() {
		actualSubstrs := findMatchingOptionSubstrings(option, searchQuery);
		if (!slices.Equal(expectedSubstrs, actualSubstrs)) {
			t.Errorf("Expected search '%v' over option '%v' to find substrings %v but found %v", searchQuery, option, expectedSubstrs, actualSubstrs);
		}
	}

	option = "s";
	searchQuery = "s";
	expectedSubstrs = []models.OptionSubstring{{Substr: "s", Index: 0}};
	assertEquals();
	
	option = "something";
	searchQuery = "s";
	expectedSubstrs = []models.OptionSubstring{{Substr: "s", Index: 0}};
	assertEquals();
		
	option = "something";
	searchQuery = "so";
	expectedSubstrs = []models.OptionSubstring{{Substr: "s", Index: 0}, {Substr: "so", Index: 0}, {Substr: "o", Index: 1}};
	assertEquals();
	
	option = "something";
	searchQuery = "m";
	expectedSubstrs = []models.OptionSubstring{{Substr: "m", Index: 2}};
	assertEquals();
	
	option = "something";
	searchQuery = "g";
	expectedSubstrs = []models.OptionSubstring{{Substr: "g", Index: 8}};
	assertEquals();
	
	option = "something";
	searchQuery = "om";
	expectedSubstrs = []models.OptionSubstring{{Substr: "o", Index: 1}, {Substr: "om", Index: 1}, {Substr: "m", Index: 2}};
	assertEquals();
		
	// only first occurrence
	option = "test";
	searchQuery = "t";
	expectedSubstrs = []models.OptionSubstring{{Substr: "t", Index: 0}};
	assertEquals();
	// duplicates are fine though
	option = "test";
	searchQuery = "tt";
	expectedSubstrs = []models.OptionSubstring{{Substr: "t", Index: 0}, {Substr: "t", Index: 0}};
	assertEquals();

	// only matches in order
	option = "test";
	searchQuery = "se";
	expectedSubstrs = []models.OptionSubstring{{Substr: "s", Index: 2}};
	assertEquals();
			
	// should trim
	option = " test ";
	searchQuery = " t ";
	expectedSubstrs = []models.OptionSubstring{{Substr: "t", Index: 0}};
	assertEquals();
				
	// should not consider whitespace in between a result
	option = "test";
	searchQuery = "t e";
	expectedSubstrs = []models.OptionSubstring{{Substr: "t", Index: 0}, {Substr: "e", Index: 1}};
	assertEquals();

	option = "t est";
	searchQuery = "t e";
	expectedSubstrs = []models.OptionSubstring{{Substr: "t", Index: 0}, {Substr: "e", Index: 2}};
	assertEquals();
					
	option = "t est";
	searchQuery = "te";
	expectedSubstrs = []models.OptionSubstring{{Substr: "t", Index: 0}, {Substr: "e", Index: 2}};
	assertEquals();
}

func TestHighlightOption_shouldReturnUnmodifiedOptionIfBlankArgs(t *testing.T) {
	var option string
	var searchQuery string

	assertNotHighlighted := func() {
		highlightedOption := highlightOption(option, findMatchingOptionSubstrings(option, searchQuery));
		if option != highlightedOption {
			t.Errorf("Expected highlighted option to be '%v' but was '%v'", option, highlightedOption);
		}
	}

	option = "option";
	searchQuery = "";
	assertNotHighlighted();

	option = "";
	searchQuery = "search";
	assertNotHighlighted();

	option = "";
	searchQuery = "";
	assertNotHighlighted();
}

func TestHighlightOption_shouldReturnUnmodifiedOptionIfSearchQueryDoesNotMatchOption(t *testing.T) {
	var option string;
	var searchQuery string;

	assertNotHighlighted := func() {
		highlightedOption := highlightOption(option, findMatchingOptionSubstrings(option, searchQuery));
		if option != highlightedOption {
			t.Errorf("Expected highlighted option to be '%v' but was '%v'", option, highlightedOption);
		}
	}

	option = "option";
	searchQuery = "xyz"; // search does not contain any option chars
	assertNotHighlighted();
}

func TestHighlightOption_shouldReturnHighlightedOptionIfSearchQueryDoesMatchOption(t *testing.T) {
	var option string;
	var searchQuery string;
	var highlightedOption string;
	var expectedHighlightedOption string;

	assertNonAnsiCharsNotModified := func() {
		if !ansiUtils.EqualsIgnoreAnsi(highlightedOption, option) {
			t.Errorf("Expected highlighted option to be '%v' but was '%v'", ansi.Strip(option), ansi.Strip(highlightedOption));
		}
	}

	assertHighlighted := func() {
		if highlightedOption != expectedHighlightedOption {
			t.Errorf("Expected highlighted option '%v' to be '%v' but was '%v'", option, expectedHighlightedOption, highlightedOption);
		}
	}

	option = "option";
	searchQuery = "o"; 
	expectedHighlightedOption = colored(searchQuery) + "ption";
	highlightedOption = highlightOption(option, findMatchingOptionSubstrings(option, searchQuery));
	assertNonAnsiCharsNotModified();
	assertHighlighted();
	
	option = "option";
	searchQuery = "pt"; 
	expectedHighlightedOption = "o" + colored(searchQuery) + "ion";
	highlightedOption = highlightOption(option, findMatchingOptionSubstrings(option, searchQuery));
	assertNonAnsiCharsNotModified();
	assertHighlighted();
		
	option = "option";
	searchQuery = "oion"; 
	expectedHighlightedOption = colored("o") + "pt" + colored("ion");
	highlightedOption = highlightOption(option, findMatchingOptionSubstrings(option, searchQuery));
	assertNonAnsiCharsNotModified();
	assertHighlighted();

	// first occurrence only
	option = "option";
	searchQuery = "oo"; 
	expectedHighlightedOption = colored("o") + "ption";
	highlightedOption = highlightOption(option, findMatchingOptionSubstrings(option, searchQuery));
	assertNonAnsiCharsNotModified();
	assertHighlighted();

	// only in order
	option = "option";
	searchQuery = "po"; 
	expectedHighlightedOption = "o" + colored("p") + "tion";
	highlightedOption = highlightOption(option, findMatchingOptionSubstrings(option, searchQuery));
	assertNonAnsiCharsNotModified();
	assertHighlighted();
}

func TestSearchAndHighlightOptions_shouldReturnOptionsIfEmptyArgs(t *testing.T) {
	options := []string{};
	searchQuery := "";

	assertResultsEqualOptions := func() {
		results := SearchAndHighlightOptions(options, searchQuery);
		if !slices.Equal(options, results) {
			t.Errorf("Expected results to equal %v but was %v", options, results);
		}
	}

	assertResultsEqualOptions();

	searchQuery = "no-empty";
	assertResultsEqualOptions();

	searchQuery = "";
	options = []string{"a"};
	assertResultsEqualOptions();
}

func TestSearchAndHighlightOptions_shouldReturnEmptyResultIfNoMatch(t *testing.T) {
	options := []string{};
	searchQuery := "";
	expectedResults := []string{};

	assertResultsEqualOptions := func() {
		results := SearchAndHighlightOptions(options, searchQuery);
		if !slices.Equal(expectedResults, results) {
			t.Errorf("Expected results to equal %v but was %v", expectedResults, results);
		}
	}

	options = []string{"a", "b", "c"};
	searchQuery = "d";
	assertResultsEqualOptions();
}

func TestSearchAndHighlightOptions_shouldReturnResults(t *testing.T) {
	options := []string{};
	searchQuery := "";
	expectedResults := []string{};

	assertResultsEqualOptions := func() {
		results := SearchAndHighlightOptions(options, searchQuery);
		if !slices.Equal(expectedResults, results) {
			t.Errorf("Expected results to equal %v but was %v", expectedResults, results);
		}
	}

	
	options = []string{"banana", "apple", "orange"};
	searchQuery = "p";
	expectedResults = []string{"a" + colored("p") + "ple"};
	assertResultsEqualOptions();

	// maintain order
	options = []string{"banana", "apple", "orange"};
	searchQuery = "a";
	expectedResults = []string{"b" + colored("a") + "nana", colored("a") + "pple", "or" + colored("a") + "nge"};
	assertResultsEqualOptions();
		
	// reorder
	options = []string{"banana", "apple", "orange"};
	searchQuery = "app";
	expectedResults = []string{colored("app") + "le", "b" + colored("a") + "nana", "or" + colored("a") + "nge"};
	assertResultsEqualOptions();
			
	// prefer consecutive match
	options = []string{"acb", "abc"};
	searchQuery = "ab";
	expectedResults = []string{colored("ab") + "c", colored("a") + "c" + colored("b")}; // both have 2 matches but one is consecutive
	assertResultsEqualOptions();
}

// [return] colored [str] using [constants.ANSWER_TEXT_COLOR]
func colored(str string) string {
	return ansiUtils.Colored(str, constants.ANSWER_TEXT_COLOR);
}
