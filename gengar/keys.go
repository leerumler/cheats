package gengar

// When gengar hits a stopKey, it will reset its
// collection of logged keystrokes.
var stopKeys = []string{

	// Navigation.
	"Up",
	"Down",
	"Left",
	"Right",
	"Page_Up",
	"Next",
	"Home",
	"End",

	// Deletions.
	"BackSpace",
	"Delete",

	// Modifiers.
	"Control_L",
	"Control_R",
	"Alt_L",
	"Alt_R",
	"Super_L",
	"Super_R",

	// Function Keys.
	"Escape",
	"Print",
	"Pause",
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"F8",
	"F9",
	"F10",
	"L1",
	"L2",
}

// When gengar hits a sendKey, it will check the keystrokes its
// collected and replace them if they match an expansion.
var sendKeys = []string{
	" ",
	"Tab",
	"Return",
}

// gengar ignores skipKeys.
var skipKeys = []string{

	// Shift keys.
	"Shift_R",
	"Shift_L",

	// Lock Keys.
	"Scroll_Lock",
	"Caps_Lock",
	"Num_Lock",
}

// Gengar ignores keystrokes that accompany skipMods.
var skipMods = []string{
	"control",
	"mod1",
	"mod4",
}
