package ghostie

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/leerumler/gengar/ggconf"
)

var weirdSyms = map[rune]string{
	' ':  "space",
	'!':  "exclam",
	'@':  "at",
	'#':  "numbersign",
	'$':  "dollar",
	'%':  "percent",
	'^':  "asciicircum",
	'&':  "ampersand",
	'*':  "asterisk",
	'(':  "parenleft",
	')':  "parenright",
	'[':  "bracketleft",
	']':  "bracketright",
	'{':  "braceleft",
	'}':  "braceright",
	'-':  "minus",
	'_':  "underscore",
	'=':  "equal",
	'+':  "plus",
	'\\': "backslash",
	'|':  "bar",
	';':  "semicolon",
	':':  "colon",
	'\'': "quoteright",
	'"':  "quotedbl",
	'<':  "less",
	'>':  "greater",
	',':  "comma",
	'.':  "period",
	'/':  "slash",
	'?':  "question",
	'`':  "quoteleft",
	'~':  "asciitilde",
}

var shiftySyms = []rune{
	'~',
	'!',
	'@',
	'#',
	'$',
	'%',
	'^',
	'&',
	'*',
	'(',
	')',
	'_',
	'+',
	'{',
	'}',
	'|',
	':',
	'"',
	'<',
	'>',
	'?',
	'A',
	'B',
	'C',
	'D',
	'E',
	'F',
	'G',
	'H',
	'I',
	'J',
	'K',
	'L',
	'M',
	'N',
	'O',
	'P',
	'Q',
	'R',
	'S',
	'T',
	'U',
	'V',
	'W',
	'X',
	'Y',
	'Z',
}

var nilKey = xproto.KeyPressEvent{
	// Detail:     nil,
	// Root:       *root,
	// Event:      *active,
	Sequence:   6,
	Time:       xproto.TimeCurrentTime,
	Child:      0,
	RootX:      1,
	RootY:      1,
	EventX:     1,
	EventY:     1,
	State:      0,
	SameScreen: true,
}

// SendKeys lets ghostie send simulated keystrokes to type messages in to the active window.  If it doesn't
// understand the keystroke (which it may not), it will do nothing.
func SendKeys(xinfo ggconf.Xinfos, expansion string) {

	keybind.Initialize(xinfo.XUtil)

	for _, charByte := range expansion {

		// var keycodes []xproto.Keycode

		charStr := string(charByte)
		if sym, okay := weirdSyms[charByte]; okay {
			charStr = sym
		}
		keycodes := keybind.StrToKeycodes(xinfo.XUtil, charStr)

		// fmt.Println(keycodes)

		var needShift bool
		for _, match := range shiftySyms {
			if match == charByte {
				needShift = true
			}
		}

		for _, keycode := range keycodes {
			key := nilKey
			key.Detail = keycode
			key.Root = *xinfo.Root
			key.Event = *xinfo.Active
			if needShift {
				key.State = xproto.ModMaskShift
			}
			xproto.SendEvent(xinfo.XUtil.Conn(), false, *xinfo.Active, xproto.EventMaskKeyRelease, string(key.Bytes()))
		}
	}
}

// Backspace inserts as many backspaces as its told to the active window.
func Backspace(xinfo ggconf.Xinfos, numKeys int) {
	for i := 0; i < numKeys; i++ {
		backspace := nilKey
		backCodes := keybind.StrToKeycodes(xinfo.XUtil, "BackSpace")
		backspace.Detail = backCodes[0]
		backspace.Root = *xinfo.Root
		backspace.Event = *xinfo.Active
		xproto.SendEvent(xinfo.XUtil.Conn(), false, *xinfo.Active, xproto.EventMaskKeyRelease, string(backspace.Bytes()))
	}
}
