package models

import (
	"unicode/utf8"

	"github.com/flobbe9/go/cmdLineMenu/constants"
	"github.com/flobbe9/go/cmdLineMenu/constants/key"
)

// Represents the raw bytes read from [os.Stdin] produced by one key stroke.
type RawStdin struct {
	Buff []byte
}

// Indicates that [this.Buff] represents a printable ascii / utf-8 char.
//
// [return] [true] if [this.Buff] is printable ascii or [this.Buff] consists of exactly two bytes
// the first of which not beeing regular ascii.
func (this *RawStdin) IsUtf8() bool {
	return this.IsPrintableAscii() || len(this.Buff) == 2 && this.Buff[0] > constants.ASCII_MAX;
}

// [return] the decoded utf-8 char from [this.Buff] or [utf8.RuneError] if is no utf-8. 
// Use [this.IsUtf8] as criteria
func (this *RawStdin) DecodeUtf8Rune() rune {
	if this.IsUtf8() {
		r, _ := utf8.DecodeRune(this.Buff);
		return r;
	}

	return utf8.RuneError;
}

// [return] [true] if there's only one byte in [this.Buff] beeing within range of 0 - 127
func (this *RawStdin) IsRegularAscii() bool {
	return len(this.Buff) == 1 && this.Buff[0] <= constants.ASCII_MAX;
}

// Indicates that [this.Buff] represents a printable, non-utf-8 char.
//
// [return] [true] if there's only one byte in [this.Buff] beeing within range of 32 - 127
func (this *RawStdin) IsPrintableAscii() bool {
	return this.IsRegularAscii() && this.Buff[0] >= constants.ASCII_PRINTABLE_MIN;
}

// Indicates that Ctrl + some other key is pressed
// 
// [return] [true] if there's only one byte in [this.Buff] beeing within range of 0 - 32
func (this *RawStdin) IsControlKeyAscii() bool {
	return this.IsRegularAscii() && !this.IsPrintableAscii();
}

// Indicates that the first bytes in [this.Buff] match an opening ansi sequence. Does not guarantee
// that [this.Buff] is actually ansi!
//
// [return] [true] if [this.Buff] starts with {27} or {27, 91 | 79}
func (this *RawStdin) LooksLikeAnsi() bool {
	if len(this.Buff) == 0 {
		return false
	}

	if this.Buff[0] == byte(constants.ANSI_ESCAPE_SEQ) {
		// case: has more bytes but the wrong ones
		if len(this.Buff) >= 2 && this.Buff[1] != 91 && this.Buff[1] != 79 {
			return false;
		}
		return true;
	}

	return false;
}

// The significant ansi key in this context is the 3 byte in [this.Buff]. It's used to determine the key
// beeing pressed. See also [ansiKey] package
// 
// [return] the 3rd byte or [key.NUL] if [this.Buff] does not look like ansi
func (this *RawStdin) GetSignificantAnsiKey() byte {
	if !this.LooksLikeAnsi() || len(this.Buff) < 3 {
		return key.NUL
	}

	return this.Buff[2];
}