package main

import (
	"bytes"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
)

// expander defines a struct that holds text epansion information.
type expander struct {
	orig, expansion []byte
}

// testInfo populates a slice of expanders with some simple testing data.
func testInfo() *[]expander {
	exps := make([]expander, 3)
	exps = append(exps, expander{[]byte("test1"), []byte("This is test 1!")})
	exps = append(exps, expander{[]byte("test2"), []byte("this is test 2")})
	exps = append(exps, expander{[]byte("test3"), []byte("this is test 3")})
	return &exps
}

// parseMatch checks if there is a match for the input and returns either
// an empty string for no match or the expansion for a match.
func parseMatch(input []byte, exps *[]expander) string {
	var expansion string
	for _, exp := range *exps {
		if comp := bytes.Compare(input, exp.orig); comp == 0 {
			expansion = string(exp.expansion)
		}
	}
	return expansion
}

// getActiveWindow returns the active window in xproto.Window form.  This can probably
// be replaced with ewmh.ActiveWindowGet(), though that will need to be investigated.
func getActiveWindow(xu *xgbutil.XUtil, root xproto.Window) xproto.Window {
	reply, err := xprop.GetProperty(xu, root, "_NET_ACTIVE_WINDOW")
	if err != nil {
		log.Fatal(err)
	}

	//
	active := xproto.Window(xgb.Get32(reply.Value))
	return active
}

// SendKeys lets gengar send simulated keystrokes to type messages in to the active window.  If it doesn't
// understand the keystroke (which it may not), it will do nothing.
func SendKeys(xinfo xinfos, expansion string) {

	keybind.Initialize(xinfo.xu)

	for _, charByte := range expansion {
		// var keycodes []xproto.Keycode
		charStr := string(charByte)
		if sym, okay := weirdSyms[charByte]; okay {
			charStr = sym
		}
		keycodes := keybind.StrToKeycodes(xinfo.xu, charStr)
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
			key.Root = *xinfo.root
			key.Event = *xinfo.active
			if needShift {
				key.State = xproto.ModMaskShift
			}
			xproto.SendEvent(xinfo.xu.Conn(), false, *xinfo.active, xproto.EventMaskKeyPress, string(key.Bytes()))
		}
	}
}

type comm struct {
	input  chan []byte
	listen chan bool
}

type xinfos struct {
	xu           *xgbutil.XUtil
	root, active *xproto.Window
}

// Backspace inserts as many backspaces as its told to the active window.
func Backspace(xinfo xinfos, numKeys int) {
	for i := 0; i < numKeys; i++ {
		backspace := nilKey
		backCodes := keybind.StrToKeycodes(xinfo.xu, "BackSpace")
		backspace.Detail = backCodes[0]
		backspace.Root = *xinfo.root
		backspace.Event = *xinfo.active
		xproto.SendEvent(xinfo.xu.Conn(), false, *xinfo.active, xproto.EventMaskKeyPress, string(backspace.Bytes()))
		// log.Println("backspace")
	}
}

// BaitAndSwitch listens for input, which in this case is any word typed within the X server,
func BaitAndSwitch(xinfo xinfos, com comm) {

	expansions := testInfo()

	// Listen and wait for input.
	for {
		select {
		case keys := <-com.input:

			// Whenever BaitAndSwitch sees input, it immediately stops the main xevent loop.
			xevent.Quit(xinfo.xu)

			//
			keyCheck := bytes.TrimSpace(keys)
			expansion := parseMatch(keyCheck, expansions)
			// fmt.Println(exp)
			if expansion != "" {
				numKeys := len(keys)
				Backspace(xinfo, numKeys)
				SendKeys(xinfo, expansion)
			}
			keybind.Detach(xinfo.xu, *xinfo.active)
			xevent.Detach(xinfo.xu, *xinfo.active)
			com.listen <- true
		}
	}
}

// WatchKeys connects to the active window and sends input whenever it reaches a terminator.
func WatchKeys(xinfo xinfos, com comm) {

	// Listen for KeyPress events on the active window.
	xwindow.New(xinfo.xu, *xinfo.active).Listen(xproto.EventMaskKeyPress)

	var inputBytes []byte

	listenForKeys := func(xu *xgbutil.XUtil, keyPress xevent.KeyPressEvent) {

		// Always have a way out.  Press ctrl+Escape at any time to exit.
		// if keybind.KeyMatch(xu, "Escape", keyPress.State, keyPress.Detail) {
		// 	if keyPress.State&xproto.ModMaskControl > 0 {
		// 		log.Println("Control-Escape detected. Quitting...")
		// 		xevent.Quit(xu)
		// 	}
		// }

		keyStr := keybind.LookupString(xu, keyPress.State, keyPress.Detail)
		inputBytes = append(inputBytes, keyStr...)

		// fmt.Println(keyStr)
		if keyStr == " " {

			//
			com.input <- inputBytes
			xevent.Quit(xu)
		}
	}

	// Finally, start the main event loop. This will route any appropriate
	// KeyPressEvents to your callback function.
	// log.Println("Program initialized. Start pressing keys!")
	xevent.KeyPressFun(listenForKeys).Connect(xinfo.xu, *xinfo.active)
	xevent.Main(xinfo.xu)
}

func conX() *xgbutil.XUtil {
	xu, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	keybind.Initialize(xu)
	return xu
}

// KeepFocus follows the _NET_ACTIVE_WINDOW and starts WatchKeys
func KeepFocus(com comm) {

	// Establish a connection to X, find the root and active windows.
	xu := conX()
	root := xu.RootWin()
	active, err := ewmh.ActiveWindowGet(xu)
	if err != nil {
		log.Println(err)
	}

	// Place that info in an xinfos struct.
	xinfo := xinfos{
		xu:     xu,
		root:   &root,
		active: &active,
	}

	// TODO: watchFocus() and

	go BaitAndSwitch(xinfo, com)
	WatchKeys(xinfo, com)

	<-com.listen

	// It's called the double tap.
	keybind.Detach(xinfo.xu, *xinfo.active)
	xevent.Detach(xinfo.xu, *xinfo.active)
	xevent.Quit(xu)

}

func main() {

	// Establish some communication channels.
	input := make(chan []byte)
	listen := make(chan bool)

	// One X event loop will be listening until it sends input,
	// the loop will then hold while gengar responds to the input.
	// normally, it should just trash everything and move on.
	com := comm{
		input:  input,
		listen: listen,
	}

	// gengar will watch in a loop until killed.
	for {
		KeepFocus(com)
	}
}
