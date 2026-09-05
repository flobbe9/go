package constants;

import "github.com/charmbracelet/x/ansi"

var ANSWER_TEXT_COLOR = ansi.RGBColor{R: 100, G: 100, B: 255}
// Every ansi opening and closing sequence starts with this string.
var ANSI_ESCAPE_SEQ = "\x1b";