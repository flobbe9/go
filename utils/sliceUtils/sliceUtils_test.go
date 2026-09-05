package sliceUtils_test

import (
	"slices"
	"testing"

	"github.com/flobbe9/go/utils/sliceUtils"
)

func TestMap_shouldPanicIfSliceArgNil(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic for nil args")
        }
    }();

	sliceUtils.Map(nil, func(el string, i int) string {return "";});
}

func TestMap_shouldPanicIfCallbackArgNil(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic for nil args")
        }
    }();

	var callback func(el string, i int) string = nil;
	sliceUtils.Map([]string{}, callback);
}

func TestMap_shouldReturnEmptySliceIfSliceArgEmpty(t *testing.T) {
	result := sliceUtils.Map([]string{}, func(str string, i int) string {return str})
	if len(result) > 0 {
		t.Errorf("Expected result to be empty for empty slice arg");
	}
}

func TestMap_shouldReturnExpectedSlice(t *testing.T) {
	s := []rune{'a', 'b', 'c', 'd'};
	expectedResult := []int{int('a'), int('b'), int('c'), int('d')};

	result := sliceUtils.Map(s, func(char rune, i int) int {return int(char)})

	if len(result) != len(s) {
		t.Errorf("Expected result length to equal slice arg length");
	}
	
	for i, expectedEl := range expectedResult {
		if expectedEl != result[i] {
			t.Errorf("Expected element at index %v to equal '%v' but got '%v'", i, expectedEl, result[i]);
		}
	}
}


func TestFilter_shouldPanicIfSliceArgNil(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic for nil args")
        }
    }();

	sliceUtils.Filter(nil, func(el string, i int) bool {return true});
}

func TestFilter_shouldPanicIfCallbackArgNil(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic for nil args")
        }
    }();

	var callback func(el string, i int) bool = nil;
	sliceUtils.Filter([]string{}, callback);
}

func TestFilter_shouldReturnEmptySliceIfEmptyArg(t *testing.T) {
	s := []string{};
	var callback func(el string, i int) bool = func(el string, i int) bool {return true};
	result := sliceUtils.Filter(s, callback);

	if len(result) > 0 {
		t.Errorf("Expected empty result for empty slice arg but got %v", result);
	}
}

func TestFilter_shouldReturnEmptySliceIfNoMatch(t *testing.T) {
	s := []int{ 0, 1, 2 };
	var callback func(el int, i int) bool = func(el int, i int) bool {
		return el > 100;
	};
	result := sliceUtils.Filter(s, callback);

	if len(result) > 0 {
		t.Errorf("Expected empty result for arg s %v arg but got %v", s, result);
	}
}

func TestFilter_shouldFilterByPredicateTrue(t *testing.T) {
	s := []int{};
	expectedResult := []int {};
	var callback func(el int, i int) bool = nil;
	
	assertFiltered := func() {
		result := sliceUtils.Filter(s, callback);
		if !slices.Equal(expectedResult, result) {
			t.Errorf("Expected result %v arg but got %v (arg %v)", expectedResult, result, s);
		}
	}
	
	s = []int{0, 1, 2};
	expectedResult = []int {1, 2};
	callback = func(el int, i int) bool {
		return el > 0;
	};
	assertFiltered();

	s = []int{0, 1, 2};
	expectedResult = []int {0, 1, 2};
	callback = func(el int, i int) bool {
		return el > -1;
	};
	assertFiltered();
	
	s = []int{0, 1, 2, 3};
	expectedResult = []int {0, 2};
	callback = func(el int, i int) bool {
		return el % 2 == 0;
	};
	assertFiltered();
}