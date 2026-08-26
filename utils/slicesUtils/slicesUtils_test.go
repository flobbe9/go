package slicesUtils_test

import (
	"testing"

	"github.com/flobbe9/go/utils/slicesUtils"
)

func TestSlicesMap_shouldPanicIfSliceArgNil(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic for nil args")
        }
    }();

	slicesUtils.SlicesMap(nil, func(el string, i int) string {return "";});
}

func TestSlicesMap_shouldPanicIfCallbackArgNil(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic for nil args")
        }
    }();

	var callback func(el string, i int) string = nil;
	slicesUtils.SlicesMap([]string{}, callback);
}

func TestSlicesMap_shouldReturnEmptySliceIfSliceArgEmpty(t *testing.T) {
	result := slicesUtils.SlicesMap([]string{}, func(str string, i int) string {return str})
	if len(result) > 0 {
		t.Errorf("Expected result to be empty for empty slice arg");
	}
}

func TestSlicesMap_shouldReturnExpectedSlice(t *testing.T) {
	s := []rune{'a', 'b', 'c', 'd'};
	expectedResult := []int{int('a'), int('b'), int('c'), int('d')};

	result := slicesUtils.SlicesMap(s, func(char rune, i int) int {return int(char)})

	if len(result) != len(s) {
		t.Errorf("Expected result length to equal slice arg length");
	}
	
	for i, expectedEl := range expectedResult {
		if expectedEl != result[i] {
			t.Errorf("Expected element at index %v to equal '%v' but got '%v'", i, expectedEl, result[i]);
		}
	}
}