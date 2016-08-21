package ghostie

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
)

// Xinfos holds information about the current X connection state.
type Xinfos struct {
	XUtil        *xgbutil.XUtil
	Root, Active *xproto.Window
}

// SendKeys lets ghostie send simulated keystrokes to type messages in to the active window.  If it doesn't
// understand the keystroke (which it may not), it will do nothing.
func SendKeys(xinfo Xinfos, expansion string) {

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
func Backspace(xinfo Xinfos, numKeys int) {
	for i := 0; i < numKeys; i++ {
		backspace := nilKey
		backCodes := keybind.StrToKeycodes(xinfo.XUtil, "BackSpace")
		backspace.Detail = backCodes[0]
		backspace.Root = *xinfo.Root
		backspace.Event = *xinfo.Active
		xproto.SendEvent(xinfo.XUtil.Conn(), false, *xinfo.Active, xproto.EventMaskKeyRelease, string(backspace.Bytes()))
	}
}
