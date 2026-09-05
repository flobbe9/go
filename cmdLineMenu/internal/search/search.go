package search

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/flobbe9/go/cmdLineMenu/constants"
	"github.com/flobbe9/go/cmdLineMenu/internal/models"
	"github.com/flobbe9/go/utils"
	"github.com/flobbe9/go/utils/sliceUtils"
	"github.com/flobbe9/go/utils/stringUtils"
)

// Assume strs are unique
// TODO test this
func SearchAndHighlightOptions(options []string, serachQuery string) []string {
	if (len(options) == 0 || stringUtils.IsBlank(serachQuery)) {
		return options;
	}

	rankedHighlightedOptions := []models.RankedOption{};
	for _, option := range options {
		relevantSubstrs := findMatchingOptionSubstrings(option, serachQuery);
		// TODO startswith should way more
		points := len(relevantSubstrs);
		highlightedOption := highlightOption(option, relevantSubstrs);
		rankedHighlightedOptions = append(rankedHighlightedOptions, models.RankedOption{Option: highlightedOption, Points: points})
	}

	// only include matching results
	rankedHighlightedOptions = sliceUtils.Filter(rankedHighlightedOptions, func(el models.RankedOption, index int) bool {
		return el.Points > 0;
	})

	// sort by points desc
	slices.SortFunc(rankedHighlightedOptions, func(a, b models.RankedOption) int {
		return b.Points - a.Points;
	});

	// map to []string
	return sliceUtils.Map(rankedHighlightedOptions, func(rankedOption models.RankedOption, index int) string {
		return rankedOption.Option;
	});
}

// something
// substrings for 'oths'
// =>
// o: 1
// t: 4
// th: 4
// h: 5

// Find substrings and their indices in [option] sothat a substr can also be found in [searchQuery].
//
// Each substring in [serachQuery] will only match once at most.
//
// Option substrings will never be matched multiple times.
//
// Whitespace will not be considered as substring.
//
// Also make sure that matches can only appear in the order from [searchQuery].
//
// [option] the option that is searched for substrings
//
// [searchQuery] the user input used to match [option] substrings against
//
// [return] list of relevant substrings in [option] that should way in on [option]'s search ranking.
// Ordered by occuring index asc.
func findMatchingOptionSubstrings(option, searchQuery string) []models.OptionSubstring {
	if (stringUtils.IsBlank(option) || stringUtils.IsBlank(searchQuery)) {
		return []models.OptionSubstring{};
	}

	option = strings.Trim(option, " ");
	searchQuery = strings.Trim(searchQuery, " ");

	matchingOptionSubstrings := []models.OptionSubstring{};

	// go through all possible searchQuery substrings
	for i := 0; i < len(searchQuery); i++ {
		if searchQuery[i] == ' ' {
			continue;
		}

		for j := i + 1; j <= len(searchQuery); j++ {
			searchQuerySubstr := searchQuery[i:j];
			if strings.Contains(searchQuerySubstr, " ") {
				break;
			}
			index := strings.Index(option, searchQuerySubstr);

			searchQuerySubstrDoesNotMatchOption := index == -1;
			searchQuerySubstrNotInOrder := len(matchingOptionSubstrings) > 0 && matchingOptionSubstrings[len(matchingOptionSubstrings) - 1].Index > index;
			if searchQuerySubstrDoesNotMatchOption || searchQuerySubstrNotInOrder {
				break;
			}

			matchingOptionSubstrings = append(matchingOptionSubstrings, models.OptionSubstring{Substr: searchQuerySubstr, Index: index});
		}
	}

	return matchingOptionSubstrings;
}

// Highlight all [optionSubstrs] in [option] using [constants.ANSWER_TEXT_COLOR].
//
// [option] the option to highlight
// 
// [optionSubstrs] substrings to highlight in [option]. 
//
// [return] a copy of [option] containing ansi chars to highlight the "foreground color" of [optionSubstrs] when printed.
func highlightOption(option string, optionSubstrs []models.OptionSubstring) string {
	if (len(optionSubstrs) == 0) {
		return option;
	}
	slog.Debug(fmt.Sprintf("Substrings with indices: %v", optionSubstrs))

	// filter out overlapping substrings. Assuming that [optionSubstrs] is sorted by [index] asc
	optionSubstrsWithoutOverlaps := []models.OptionSubstring{};
	prev := optionSubstrs[0];
	for i := 1; i < len(optionSubstrs); i++ {
		if (optionSubstrs[i].Index != prev.Index) {
			optionSubstrsWithoutOverlaps = append(optionSubstrsWithoutOverlaps, prev);
		}
		prev = optionSubstrs[i];
	}
	optionSubstrsWithoutOverlaps = append(optionSubstrsWithoutOverlaps, prev);
	slog.Debug(fmt.Sprintf("Substrings with distinct indices: %v", optionSubstrsWithoutOverlaps))

	// rebuild option highlighting substrs 
	var highlightedOption strings.Builder;
	for i := 0; i < len(option); i++ {
		nextOptionSubstr := optionSubstrsWithoutOverlaps[0];

		// case: not a substr to highlight
		if (nextOptionSubstr.Index != i) {
			_, err := highlightedOption.WriteString(string(option[i]));
			utils.ErrorLogExit(err);

		} else {
			// append highlighted substr
			_, err := highlightedOption.WriteString(ansi.NewStyle().ForegroundColor(constants.ANSWER_TEXT_COLOR).Styled(nextOptionSubstr.Substr));
			utils.ErrorLogExit(err);

			// skip ahead, we don't want to append chars multiple times
			i += len(nextOptionSubstr.Substr) - 1;

			// pop substr from "queue"
			optionSubstrsWithoutOverlaps = optionSubstrsWithoutOverlaps[1:];
		}

		// case: nothing more to highlight
		if (len(optionSubstrsWithoutOverlaps) == 0) {
			// write the remaining unhighlighted chars
			_, err := highlightedOption.WriteString(option[i + 1:]);
			utils.ErrorLogExit(err);
			break;
		}
	}

	return highlightedOption.String();
}
