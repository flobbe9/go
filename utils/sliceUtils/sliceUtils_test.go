package sliceUtils_test

import (
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