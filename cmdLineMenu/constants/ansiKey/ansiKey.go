package ansiKey;

// Default escape sequence is "27 91"
const (
	// Remember that most terminals have also default behaviour for Escape. Using this with 
	// RawScanln will propably lead to defered callback execution
	Escape = 27;
	Ctrl_ArrowRight = 49;
	Delete = 51;
	ArrowUp = 65;
	ArrowDown = 66;
	ArrowRight = 67;
	ArrowLeft = 68;
	// With ansi opening sequence 91 79
	F1 = 80;
	// With ansi opening sequence 91 79
	F2 = 81;
	// With ansi opening sequence 91 79
	F3 = 82;
	// With ansi opening sequence 91 79
	F4 = 83;
)