package ggui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/leerumler/gengar/ggconf"
	"github.com/leerumler/gengar/ggdb"
)

var exps *[]ggconf.Expander

// getGooey creates the GUI.
func getGooey() {

	// Create a new GUI.
	gooey := gocui.NewGui()
	if err := gooey.Init(); err != nil {
		log.Fatal(err)
	}
	defer gooey.Close()
}

func drawPhrases(gooey *gocui.Gui) error {

	// Get window size.
	maxX, maxY := gooey.Size()

	//
	if phraseView, err := gooey.SetView("phrases", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		phraseView.SelBgColor = gocui.ColorCyan
		phraseView.Highlight = true

		for _, exp := range *exps {
			fmt.Fprintln(phraseView, exp.Phrase)
		}
	}

	return nil

}

func drawExpansions(gooey *gocui.Gui) error {

	// Get window size.
	maxX, maxY := gooey.Size()

	if expView, err := gooey.SetView("galaxy", 10, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(expView, "Long ago, in a galaxy far, far away, there was no need for overly complicated programming languages designed to make simple tasks easier for users and more difficult for programmers.  In fact, there were no users, programmers, or computers.  There was nothing, in fact, because the galaxy did not contain life.")
		expView.Wrap = true
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func readExpanders() {
	exps = ggdb.ReadExpanders()
}

// Example creates an example GUI.
func Example() {

	readExpanders()

	gooey := gocui.NewGui()
	if err := gooey.Init(); err != nil {
		log.Fatal(err)
	}
	defer gooey.Close()

	gooey.SetLayout(drawPhrases)
	gooey.SetLayout(drawExpansions)

	if err := gooey.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
