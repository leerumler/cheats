package main

import (
	"bytes"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type expander struct {
	orig, expansion []byte
}

func testInfo() *[]expander {
	exps := make([]expander, 3)
	exps = append(exps, expander{[]byte("test1"), []byte("this is test 1")})
	exps = append(exps, expander{[]byte("test2"), []byte("this is test 2")})
	exps = append(exps, expander{[]byte("test3"), []byte("this is test 3")})
	return &exps
}

func parseMatch(input []byte, exps *[]expander) []byte {
	var expansion []byte
	for _, exp := range *exps {
		if comp := bytes.Compare(input, exp.orig); comp == 0 {
			expansion = exp.expansion
		}
	}
	return expansion
}

func getActiveWindow(xu *xgbutil.XUtil, root xproto.Window) xproto.Window {
	reply, err := xprop.GetProperty(xu, root, "_NET_ACTIVE_WINDOW")
	if err != nil {
		log.Fatal(err)
	}

	//
	active := xproto.Window(xgb.Get32(reply.Value))
	return active
}

func sendKeys(xu *xgbutil.XUtil, root, active *xproto.Window, exp []byte) {
	nilKey := xproto.KeyPressEvent{
		// Detail:     nil,
		Sequence:   6,
		Time:       xproto.TimeCurrentTime,
		Root:       *root,
		Event:      *active,
		Child:      0,
		RootX:      1,
		RootY:      1,
		EventX:     1,
		EventY:     1,
		State:      0,
		SameScreen: true,
	}

	keybind.Initialize(xu)

	for _, charByte := range exp {
		// var keycodes []xproto.Keycode
		charStr := string(charByte)
		if charStr == " " {
			charStr = "space"
		}
		keycodes := keybind.StrToKeycodes(xu, charStr)
		// fmt.Println(keycodes)

		for _, keycode := range keycodes {
			key := nilKey
			key.Detail = keycode
			xproto.SendEvent(xu.Conn(), false, *active, xproto.EventMaskKeyPress, string(key.Bytes()))
		}
	}

	// xproto.SendEvent(c, Propagate, Destination, EventMask, Event)
}

func checkInput(xu *xgbutil.XUtil, root, active *xproto.Window, input chan []byte, listen chan bool) {
	// fmt.Println("replacing...")
	nilKey := xproto.KeyPressEvent{
		// Detail:     nil,
		Sequence:   6,
		Time:       xproto.TimeCurrentTime,
		Root:       *root,
		Event:      *active,
		Child:      0,
		RootX:      1,
		RootY:      1,
		EventX:     1,
		EventY:     1,
		State:      0,
		SameScreen: true,
	}

	exps := testInfo()
	for {
		select {
		case keys := <-input:
			keyCheck := bytes.TrimSpace(keys)
			exp := parseMatch(keyCheck, exps)
			// fmt.Println(exp)
			if exp != nil {
				for range keys {
					backspace := nilKey
					backspace.Detail = 22
					xproto.SendEvent(xu.Conn(), false, *active, xproto.EventMaskKeyPress, string(backspace.Bytes()))
					// log.Println("backspace")
				}
				sendKeys(xu, root, active, exp)
			}
			listen <- true
		}
	}
}

func listenClosely(xu *xgbutil.XUtil, root, active *xproto.Window, input chan []byte) {

	// Listen for KeyPress events on the active window.
	xwindow.New(xu, *active).Listen(xproto.EventMaskKeyPress)

	var inputBytes []byte

	listenForKeys := func(xu *xgbutil.XUtil, keyPress xevent.KeyPressEvent) {
		// Always have a way out.  Press ctrl+Escape to exit.
		if keybind.KeyMatch(xu, "Escape", keyPress.State, keyPress.Detail) {
			if keyPress.State&xproto.ModMaskControl > 0 {
				log.Println("Control-Escape detected. Quitting...")
				xevent.Quit(xu)
			}
		}

		keyStr := keybind.LookupString(xu, keyPress.State, keyPress.Detail)
		inputBytes = append(inputBytes, keyStr...)

		// fmt.Println(keyStr)
		if keyStr == " " {

			//
			input <- inputBytes
			// keybind.Detach(xu, *active)
			// xevent.Detach(xu, *active)
			xevent.Quit(xu)
		}
	}

	// Finally, start the main event loop. This will route any appropriate
	// KeyPressEvents to your callback function.
	// log.Println("Program initialized. Start pressing keys!")
	xevent.KeyPressFun(listenForKeys).Connect(xu, *active)
	xevent.Main(xu)
}

func main() {
	input := make(chan []byte)
	listen := make(chan bool)

	for {
		// fmt.Println("Listening...")
		// Connect to the X server using the DISPLAY environment variable
		// and initialize keybind.
		xu, err := xgbutil.NewConn()
		if err != nil {
			log.Fatal(err)
		}
		keybind.Initialize(xu)

		// Get the window id of the root and active windows.
		root := xproto.Setup(xu.Conn()).DefaultScreen(xu.Conn()).Root
		active := getActiveWindow(xu, root)

		go checkInput(xu, &root, &active, input, listen)
		listenClosely(xu, &root, &active, input)
		<-listen
	}
}
