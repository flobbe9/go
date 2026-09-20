package models

import "github.com/flobbe9/go/utils/stringUtils"

// Wrapper around an arbitrary, non-blank substring of an option including it's index
type OptionSubstring struct {
	// Substring of an option
	Substr string

	// Start index of [substr] in option. May be -1
	Start int
}

func (this *OptionSubstring) End() int {
	if len(this.Substr) == 0 {
		return 0;
	}

	end := stringUtils.Len(this.Substr) - 1 + this.Start;

	if end < 0 {
		return -1;
	}

	return end;
}

// Indicates that either [opt.Substr] is overlapping [this.Substr] or the other way round.
//
// Equal end and start are an overlap as well, e.g. for option "abcde": "abc" overlaps "cde"
func (this *OptionSubstring) IsOverlapping(opt OptionSubstring) bool {
	if stringUtils.IsBlank(this.Substr) && stringUtils.IsBlank(opt.Substr) {
		return false;
	}

	isThisOverlappingOpt2 := this.Start >= opt.Start && this.Start <= opt.End();
	isOpt2OverlappingThis := opt.Start >= this.Start && opt.Start <= this.End();

	return isThisOverlappingOpt2 || isOpt2OverlappingThis;
}