package constants

import (
	"github.com/charmbracelet/x/ansi"
)

var ANSWER_TEXT_COLOR = ansi.RGBColor{R: 100, G: 100, B: 255}
// Every ansi opening and closing sequence starts with this string (decimal 27)
var ANSI_ESCAPE_SEQ = '\x1b';

// The highest char code for regular ascii
const ASCII_MAX = 127;
// The lowest printable char code for regular ascii
const ASCII_PRINTABLE_MIN = 32;