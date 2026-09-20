package search

import (
	"fmt"
	"log/slog"
	"math"
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
// TODO
// should consider the second occurrence of a letter for the rest of the string
func SearchAndHighlightOptions(options []string, searchQuery string) []string {
	if (len(options) == 0 || stringUtils.IsBlank(searchQuery)) {
		return options;
	}

	rankedHighlightedOptions := []models.RankedOption{};
	for _, option := range options {
		relevantSubstrs := findMatchingOptionSubstrings(option, searchQuery);
		points := calculateSearchRankingPoints(option, relevantSubstrs);
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

// Find substrings and their indices in [option] sothat a substr can also be found in [searchQuery].
//
// Each substring in [searchQuery] will only match once at most.
//
// Option substrings will never be matched multiple times.
//
// Whitespace will not be considered as substring.
//
// Make sure that matches can only appear in the order from [searchQuery].
//
// Filter out overlapping matches.
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
			searchQuerySubstrNotInOrder := len(matchingOptionSubstrings) > 0 && matchingOptionSubstrings[len(matchingOptionSubstrings) - 1].Start > index;
			if searchQuerySubstrDoesNotMatchOption || searchQuerySubstrNotInOrder {
				break;
			}

			matchingOptionSubstrings = append(matchingOptionSubstrings, models.OptionSubstring{Substr: searchQuerySubstr, Start: index});
		}
	}

	// extra iterations, shame on me
	matchingOptionSubstrings = filterOverlappingOptionSubstrings(matchingOptionSubstrings);

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

	// rebuild option highlighting substrs 
	var highlightedOption strings.Builder;
	for i := 0; i < len(option); i++ {
		nextOptionSubstr := optionSubstrs[0];

		// case: not a substr to highlight
		if (nextOptionSubstr.Start != i) {
			_, err := highlightedOption.WriteString(string(option[i]));
			utils.ErrorLogExit(err);

		} else {
			// append highlighted substr
			_, err := highlightedOption.WriteString(ansi.NewStyle().ForegroundColor(constants.ANSWER_TEXT_COLOR).Styled(nextOptionSubstr.Substr));
			utils.ErrorLogExit(err);

			// skip ahead, we don't want to append chars multiple times
			i += len(nextOptionSubstr.Substr) - 1;

			// pop substr from "queue"
			optionSubstrs = optionSubstrs[1:];
		}

		// case: nothing more to highlight
		if (len(optionSubstrs) == 0) {
			// write the remaining unhighlighted chars
			_, err := highlightedOption.WriteString(option[i + 1:]);
			utils.ErrorLogExit(err);
			break;
		}
	}

	return highlightedOption.String();
}

// Remove overlapping substrings keeping only the longest ones but always the one with the lowest index
// Assume that [optionSubstrs] is sorted by [Index].
//
// [optionSubstrs] matching substrings for an option
//
// [return] non-overlapping substrings
func filterOverlappingOptionSubstrings(optionSubstrs []models.OptionSubstring) []models.OptionSubstring {
	slog.Debug(fmt.Sprintf("Substrings with indices: %v", optionSubstrs))

	// case: no overlap possible
	if len(optionSubstrs) <= 1 {
		return optionSubstrs;
	}

	// squash options with same index sothat the last one stays
	squashedOptionSubstrs := []models.OptionSubstring{};
	var prev models.OptionSubstring;
	for i := 1; i < len(optionSubstrs); i++ {
		cur := optionSubstrs[i];
		prev = optionSubstrs[i - 1];
		isCurIndexDifferent := cur.Start != prev.Start;
		
		if isCurIndexDifferent {
			squashedOptionSubstrs = append(squashedOptionSubstrs, prev);
		}
			
		// add last element
		if i == len(optionSubstrs) - 1 && (isCurIndexDifferent || len(squashedOptionSubstrs) == 0) { 
			squashedOptionSubstrs = append(squashedOptionSubstrs, cur);
		}
	}
	
	// filter out overlapping substrs
	optionSubstrsWithoutOverlaps := []models.OptionSubstring{squashedOptionSubstrs[0]};
	for _, squashedOptionSubstr := range squashedOptionSubstrs {
		lastOptionWithSubstrWithoutOverlap := optionSubstrsWithoutOverlaps[len(optionSubstrsWithoutOverlaps) - 1];
		if !squashedOptionSubstr.IsOverlapping(lastOptionWithSubstrWithoutOverlap) {
			optionSubstrsWithoutOverlaps = append(optionSubstrsWithoutOverlaps, squashedOptionSubstr);
		}
	}
	slog.Debug(fmt.Sprintf("Substrings with distinct indices: %v", optionSubstrsWithoutOverlaps))

	return optionSubstrsWithoutOverlaps;
}

// Determine ranking points for [option] depending on [relevantSubstrs].
//
// [relevantSubstrs] derrived from [filterOverlappingOptionSubstrings], expected to be ordered
// 
// [return] ranking points, never negative
func calculateSearchRankingPoints(option string, relevantSubstrs []models.OptionSubstring) int {
	if (len(relevantSubstrs) == 0) {
		return 0;
	}

	var points int;
	var totalMatchingChars int;
	var numSubstrsLongerThan1 int;

	for _, relevantOptionSubstr := range relevantSubstrs {
		totalMatchingChars += stringUtils.Len(relevantOptionSubstr.Substr);
		if len(relevantOptionSubstr.Substr) > 1 {
			numSubstrsLongerThan1++;
		}
	}

	points += totalMatchingChars;
	slog.Debug(fmt.Sprintf("Search points for total num chars: %v", points));

	points += numSubstrsLongerThan1;
	slog.Debug(fmt.Sprintf("Plus search points for long substrs: %v", points));

	// always prioritise starts with
	if stringUtils.StartsWith(option, relevantSubstrs[0].Substr) {
		points += math.MaxInt / 2;
		slog.Debug(fmt.Sprintf("Plus search points for startswith: %v", points));
	}

	return points;
}