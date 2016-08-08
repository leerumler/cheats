package gengar

import (
	"log"
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/leerumler/gengar/ggconf"
	"github.com/leerumler/gengar/ggdb"
	"github.com/leerumler/gengar/ghostie"
)

// comm defines a struct that holds communication channels between
// our various goroutines.
type comm struct {
	input   chan string
	refresh chan bool
}

// parseMatch checks if there is a match for the input and returns either
// an empty string for no match or the expansion for a match.
func parseMatch(input string, exps *[]ggconf.Expander) string {
	var expansion string
	for _, exp := range *exps {
		if comp := strings.Compare(input, exp.Phrase); comp == 0 {
			expansion = string(exp.Expansion)
		}
	}
	return expansion
}

// getActiveWindow returns a pointer to the active window.
func getActive(xinfo ggconf.Xinfos) *xproto.Window {

	// Check the active window, or die if we can't get it.
	active, err := ewmh.ActiveWindowGet(xinfo.XUtil)
	if err != nil {
		log.Fatal(err)
	}

	// Return a pointer to that value.
	return &active
}

func conX() ggconf.Xinfos {

	// Create a space in memory to hold information
	// about the current X connection state.
	var xinfo ggconf.Xinfos

	// Connect to X, or die.  Initialize keybind, so we can
	// use some of its functions.
	xu, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	keybind.Initialize(xu)

	// Store that connection in xinfo.
	xinfo.XUtil = xu

	// Check the root window and store that in xinfo.
	root := xinfo.XUtil.RootWin()
	xinfo.Root = &root

	// Check the active window, or die.  Store that in xinfo.
	xinfo.Active = getActive(xinfo)

	// Finally, return that info.
	return xinfo
}

// BaitAndSwitch listens for input, which in this case is any word typed within the X server,
func BaitAndSwitch(com comm) {

	// Establish a connection to X, find the root and active windows.
	xinfo := conX()

	// Load up possible expansions from testInfo().
	expanders := ggdb.ReadExpanders()

	// Listen and wait for instructions in com channels.
	for {
		select {

		// If we get input, check and replace.
		case keys := <-com.input:

			// Log the received input.  This is for debugging purposes and should
			// be disabled on production systems.  Gengar is not a keylogger.
			log.Println("Received Input:", keys)

			// Check the input, and if it matches, return an expansion.
			expansion := parseMatch(keys, expanders)

			// If we got an expansion back, send a series of backspaces to wipe out
			// the input and replace it with the expansion.
			if expansion != "" {
				ghostie.Backspace(xinfo, len(keys)+1)
				ghostie.SendKeys(xinfo, expansion)
			}

		case <-com.refresh:

			// If we receive a signal to refresh, reset the current active window
			xinfo.Active = getActive(xinfo)

			// And re-read the expansion database.
			expanders = ggdb.ReadExpanders()
		}
	}
}

// WatchKeys connects to the active window and sends input whenever it reaches a terminator.
func WatchKeys(xinfo ggconf.Xinfos, com comm) {

	// Listen for KeyPress events on the active window.
	err := xwindow.New(xinfo.XUtil, *xinfo.Active).Listen(xproto.EventMaskKeyPress)
	if err != nil {
		log.Fatal(err)
	}

	// This is where we'll be storing the keystrokes.
	var keyBytes []byte

	// Attach a callback function that listens for keypress events.
	xevent.KeyPressFun(
		func(xu *xgbutil.XUtil, keyPress xevent.KeyPressEvent) {

			// Whenever we see a keypress event, look up what key was pressed
			keyStr := keybind.LookupString(xu, keyPress.State, keyPress.Detail)

			// Check to see if the key is in skipKeys.
			keep := true
			for _, skip := range skipKeys {
				if keyStr == skip {
					log.Println("Skipped Key:", keyStr)
					keep = false
				}
			}

			// As long as it's not in skipKeys, add it to keyBytes.
			if keep {

				// and append that to our byte slice.
				keyBytes = append(keyBytes, keyStr...)

				// Log the key that was pressed.  This should really be disabled
				// whenever it isn't necessary, as it's a bit of a security risk.
				log.Println("Key logged:", keyStr)
			}

			// If we get a sendKey, send off the collected byte slice to BaitAndSwitch
			// and empty out keyBytes for the next word.
			for _, send := range sendKeys {
				if keyStr == send {

					// Trim off the last keyStr before sending.
					keys := strings.TrimSuffix(string(keyBytes), send)
					com.input <- keys
					keyBytes = nil
				}
			}

			// If we get a stopKey, just empty out keyBytes.
			for _, stop := range stopKeys {
				if keyStr == stop {
					log.Println("Not sending:", string(keyBytes))
					keyBytes = nil
				}
			}

		}).Connect(xinfo.XUtil, *xinfo.Active)
}

// KeepFocus watches for changes in the _NET_ACTIVE_WINDOW property of the root window.
// If a change is detected, it sends a com.refresh signal and quits the X event loop.
func KeepFocus(xinfo ggconf.Xinfos, com comm) {

	// Listen for property changes on the root window or die.
	err := xwindow.New(xinfo.XUtil, *xinfo.Root).Listen(xproto.EventMaskPropertyChange)
	if err != nil {
		log.Fatal(err)
	}

	// Whenever a change is detected, check if it was to the _NET_ACTIVE_WINDOW property.
	// If it was, print to the log, send a signal to com.refresh, and quit the X event loop.
	xevent.PropertyNotifyFun(
		func(xu *xgbutil.XUtil, propEve xevent.PropertyNotifyEvent) {
			if xinfo.XUtil.AtomNames[propEve.Atom] == "_NET_ACTIVE_WINDOW" {
				log.Println("Focus changed, restarting event loop.")
				com.refresh <- true
				xevent.Quit(xinfo.XUtil)
			}
		}).Connect(xinfo.XUtil, *xinfo.Root)
}

// ListenClosely does some things.
func ListenClosely() {

	// Establish some communication channels.
	com := comm{
		input:   make(chan string),
		refresh: make(chan bool),
	}

	// Start Bait and Switch.
	go BaitAndSwitch(com)

	// Establish an X connection and get some information.
	xinfo := conX()

	// gengar will listen to keyboard input until killed.
	for {

		// Initialize some values.  These are actually the default values, but
		// we want to reset them in case they got unset, which they will be.
		keybind.Initialize(xinfo.XUtil)
		xinfo.XUtil.Quit = false

		// Establish a connection to X, find the root and active windows.
		xinfo.Active = getActive(xinfo)
		log.Println("Starting listen loop.")

		// Hook on the callback functions.
		KeepFocus(xinfo, com)

		// Don't try and start WatchKeys if we don't have an active window.
		if *xinfo.Active != 0 {
			WatchKeys(xinfo, com)
		}

		// Start the main event loop.
		xevent.Main(xinfo.XUtil)

		// Detach the callback functions.
		keybind.Detach(xinfo.XUtil, *xinfo.Root)
		xevent.Detach(xinfo.XUtil, *xinfo.Root)
		keybind.Detach(xinfo.XUtil, *xinfo.Active)
		xevent.Detach(xinfo.XUtil, *xinfo.Active)

		// Just in case.
		xevent.Quit(xinfo.XUtil)

		// Tell me when it's done.
		log.Println("Listen loop stopped.")
	}
}
