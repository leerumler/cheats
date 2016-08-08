package ggui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// GetGooey creates the GUI.
func GetGooey() {
	gooey := gocui.NewGui()
	if err := gooey.Init(); err != nil {
		// handle error
	}
	defer gooey.Close()
}

func layout(gooey *gocui.Gui) error {
	maxX, maxY := gooey.Size()
	if v, err := gooey.SetView("hello", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Example creates an example GUI.
func Example() {
	gooey := gocui.NewGui()
	if err := gooey.Init(); err != nil {
		log.Fatal(err)
	}
	defer gooey.Close()

	gooey.SetLayout(layout)

	if err := gooey.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
