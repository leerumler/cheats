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

func getActiveWindow(xu *xgbutil.XUtil, root xproto.Window) xproto.Window {
	reply, err := xprop.GetProperty(xu, root, "_NET_ACTIVE_WINDOW")
	if err != nil {
		log.Fatal(err)
	}

	//
	active := xproto.Window(xgb.Get32(reply.Value))
	return active
}

func replaceInput(xu *xgbutil.XUtil, root, active *xproto.Window, input []byte) []byte {
	// fmt.Println("replacing...")
	nilKey := xproto.KeyPressEvent{
		Sequence: 6,
		// Detail:     37,
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

	// keys := <-input
	// for count := 0; count < len(keys); count++ {
	// 	backspace := nilKey
	// 	backspace.Detail = 22
	// 	xproto.SendEvent(xu.Conn(), false, *active, xproto.EventMaskKeyPress, string(backspace.Bytes()))
	// }
	// keys := <-input:
	for count := 0; count < len(input); count++ {
		backspace := nilKey
		backspace.Detail = 22
		xproto.SendEvent(xu.Conn(), false, *active, xproto.EventMaskKeyPress, string(backspace.Bytes()))
	}
	return []byte{}
}

func listenClosely(xu *xgbutil.XUtil, root, active *xproto.Window, input []byte) []byte {
	// fmt.Println("Listening...")

	// Listen for KeyPress events on the active window.
	xwindow.New(xu, *active).Listen(xproto.EventMaskKeyPress)

	// keys := make([]byte, 10)

	listenForKeys := func(xu *xgbutil.XUtil, keyPress xevent.KeyPressEvent) {
		// Always have a way out.  Press ctrl+Escape to exit.
		if keybind.KeyMatch(xu, "Escape", keyPress.State, keyPress.Detail) {
			if keyPress.State&xproto.ModMaskControl > 0 {
				log.Println("Control-Escape detected. Quitting...")
				xevent.Quit(xu)
			}
		}

		keyStr := keybind.LookupString(xu, keyPress.State, keyPress.Detail)
		input = append(input, keyStr...)
		if keyStr == " " {
			// fmt.Print(input)
			// input <- keys
			xevent.Quit(xu)
		}
	}
	xevent.KeyPressFun(listenForKeys).Connect(xu, *active)

	// Finally, start the main event loop. This will route any appropriate
	// KeyPressEvents to your callback function.
	// log.Println("Program initialized. Start pressing keys!")
	xevent.Main(xu)
	return input
}

func main() {
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
	// win := xwindow.New(xu, active)
	//

	input := []byte{}
	for {
		input = listenClosely(xu, &root, &active, input)
		input = replaceInput(xu, &root, &active, input)
		// <-input
	}

	// replaceInput(xu, root, active, input)
	// input = nil
}
