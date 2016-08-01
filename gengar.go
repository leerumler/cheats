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

type comm struct {
	input   chan []byte
	refresh chan bool
	// wait  sync.WaitGroup
}

type xinfos struct {
	xu           *xgbutil.XUtil
	root, active *xproto.Window
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
func BaitAndSwitch(com comm) {

	// Establish a connection to X, find the root and active windows.
	xinfo := conX()

	expansions := testInfo()

	// Listen and wait for input.
	for {
		select {
		case keys := <-com.input:
			//
			keyCheck := bytes.TrimSpace(keys)
			expansion := parseMatch(keyCheck, expansions)
			// fmt.Println(exp)
			if expansion != "" {
				Backspace(xinfo, len(keys))
				SendKeys(xinfo, expansion)
			}

		case <-com.refresh:

			var err error
			active, err := ewmh.ActiveWindowGet(xinfo.xu)
			if err != nil {
				log.Fatal(err)
			}
			xinfo.active = &active

		}
	}
}

// WatchKeys connects to the active window and sends input whenever it reaches a terminator.
func WatchKeys(xinfo xinfos, com comm) {

	// Listen for KeyPress events on the active window.
	xwindow.New(xinfo.xu, *xinfo.active).Listen(xproto.EventMaskKeyPress)

	var inputBytes []byte

	listenForKeys := func(xu *xgbutil.XUtil, keyPress xevent.KeyPressEvent) {

		keyStr := keybind.LookupString(xu, keyPress.State, keyPress.Detail)
		inputBytes = append(inputBytes, keyStr...)

		// fmt.Println(keyStr)
		if keyStr == " " {

			//
			com.input <- inputBytes
			xevent.Quit(xu)
		}
	}
	xevent.KeyPressFun(listenForKeys).Connect(xinfo.xu, *xinfo.active)
}

func conX() xinfos {
	// Create a space in memory to hold information about the current
	// X connection state.
	var xinfo xinfos

	xu, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	keybind.Initialize(xu)

	xinfo.xu = xu
	root := xinfo.xu.RootWin()
	xinfo.root = &root

	// Check the active window.
	active, err := ewmh.ActiveWindowGet(xinfo.xu)
	if err != nil {
		log.Fatal(err)
	}

	xinfo.active = &active

	return xinfo
}

// KeepFocus watches for changes in the _NET_ACTIVE_WINDOW property of the root window,
// sends a com.refresh signal, and quits the X event loop.
func KeepFocus(xinfo xinfos, com comm) {

	//
	err := xwindow.New(xinfo.xu, *xinfo.root).Listen(xproto.EventMaskPropertyChange)
	if err != nil {
		log.Fatal(err)
	}

	xevent.PropertyNotifyFun(
		func(xu *xgbutil.XUtil, propEve xevent.PropertyNotifyEvent) {
			if xinfo.xu.AtomNames[propEve.Atom] == "_NET_ACTIVE_WINDOW" {
				// log.Println(propEve.Atom)
				log.Println("Focus changed, restarting event loop.")
				com.refresh <- true
				xevent.Quit(xinfo.xu)
			}
		}).Connect(xinfo.xu, *xinfo.root)

}

func main() {

	// Establish some communication channels.
	com := comm{
		input:   make(chan []byte),
		refresh: make(chan bool),
	}

	go BaitAndSwitch(com)

	// gengar will listen to keyboard input until killed.
	for {

		// Establish a connection to X, find the root and active windows.
		xinfo := conX()

		//
		KeepFocus(xinfo, com)
		WatchKeys(xinfo, com)
		xevent.Main(xinfo.xu)

		// It's called the double tap.
		keybind.Detach(xinfo.xu, *xinfo.active)
		xevent.Detach(xinfo.xu, *xinfo.active)
		xevent.Quit(xinfo.xu)
	}
}
