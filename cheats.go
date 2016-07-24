package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type expander struct {
	text, expansion string
}

func testInfo() *[]expander {
	exps := make([]expander, 3)
	exps = append(exps, expander{"test1", "this is test 1"})
	exps = append(exps, expander{"test2", "this is test 2"})
	exps = append(exps, expander{"test3", "this is test 3"})
	return &exps
}

func testUsage() {
	exps := testInfo()
	input := prompt("Input Test Statement: ")
	result := parseMatch(input, exps)
	fmt.Println(result)
}

func prompt(question string) *string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	input, _ := reader.ReadString('\n')
	answer := strings.TrimSuffix(input, "\n")
	return &answer
}

func parseMatch(input *string, exps *[]expander) string {
	var expansion string
	for _, exp := range *exps {
		if *input == exp.text {
			expansion = exp.expansion
		}
	}
	return expansion
}

var flagRoot = false

func init() {
	log.SetFlags(0)
	flag.BoolVar(&flagRoot, "root", flagRoot,
		"When set, the keyboard will be grabbed on the root window. "+
			"Make sure you have a way to kill the window created with "+
			"the mouse.")
	flag.Parse()
}

func main() {

	// Connect to the X server using the DISPLAY environment variable.
	xu, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Anytime the keybind (mousebind) package is used, keybind.Initialize
	// *should* be called once. It isn't strictly necessary, but allows your
	// keybindings to persist even if the keyboard mapping is changed during
	// run-time. (Assuming you're using the xevent package's event loop.)
	// It also handles the case when your modifier map is changed.
	keybind.Initialize(xu)

	// Create a new window. We will listen for key presses and translate them
	// only when this window is in focus. (Similar to how `xev` works.)
	// win, err := xwindow.Generate(X)
	// if err != nil {
	// 	log.Fatalf("Could not generate a new window X id: %s", err)
	// }
	// win.Create(X.RootWin(), 0, 0, 500, 500, xproto.CwBackPixel, 0xffffffff)

	// Get the window id of the root window.
	setup := xproto.Setup(xu.Conn())
	root := setup.DefaultScreen(xu.Conn()).Root

	// Get the atom id (i.e., intern an atom) of "_NET_ACTIVE_WINDOW".
	reply, err := xprop.GetProperty(xu, root, "_NET_ACTIVE_WINDOW")
	if err != nil {
		log.Fatal(err)
	}

	//

	winproto := xproto.Window(xgb.Get32(reply.Value))
	win := xwindow.New(xu, winproto)
	//
	// // Listen for Key{Press,Release} events.
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease)
	//
	// // Map the window.
	// win.Map()
	//
	// // Notice that we use xevent.KeyPressFun instead of keybind.KeyPressFun,
	// // because we aren't trying to make a grab *and* because we want to listen
	// // to *all* key press events, rather than just a particular key sequence
	// // that has been pressed.
	wid := win.Id
	if flagRoot {
		wid = xu.RootWin()
	}
	xevent.KeyPressFun(
		func(xu *xgbutil.XUtil, e xevent.KeyPressEvent) {
			// keybind.LookupString does the magic of implementing parts of
			// the X Keyboard Encoding to determine an english representation
			// of the modifiers/keycode tuple.
			// N.B. It's working for me, but probably isn't 100% correct in
			// all environments yet.
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(xu, e.State, e.Detail)
			if len(modStr) > 0 {
				log.Printf("Key: %s-%s\n", modStr, keyStr)
			} else {
				log.Println("Key:", keyStr)
			}

			if keybind.KeyMatch(xu, "Escape", e.State, e.Detail) {
				if e.State&xproto.ModMaskControl > 0 {
					log.Println("Control-Escape detected. Quitting...")
					xevent.Quit(xu)
				}
			}
		}).Connect(xu, wid)

	// If we want root, then we take over the entire keyboard.
	if flagRoot {
		if err := keybind.GrabKeyboard(xu, xu.RootWin()); err != nil {
			log.Fatalf("Could not grab keyboard: %s", err)
		}
		log.Println("WARNING: We are taking *complete* control of the root " +
			"window. The only way out is to press 'Control + Escape' or to " +
			"close the window with the mouse.")
	}

	// Finally, start the main event loop. This will route any appropriate
	// KeyPressEvents to your callback function.
	log.Println("Program initialized. Start pressing keys!")
	xevent.Main(xu)
}
