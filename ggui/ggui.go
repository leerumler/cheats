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

func centerText(text string, maxX int) string {
	numSpaces := maxX/2 - len(text)/2
	for i := 0; i < numSpaces; i++ {
		text = " " + text
	}
	return text
}

func padText(text string, maxX int) string {
	numSpaces := maxX - len(text)
	for i := 0; i < numSpaces; i++ {
		text = text + " "
	}
	return text
}

func drawHeader(gooey *gocui.Gui) error {

	maxX, _ := gooey.Size()

	if header, err := gooey.SetView("header", 0, 0, maxX, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		head := centerText("Gengar Configuration Editor", maxX)

		fmt.Fprintln(header, head)
	}

	return nil

}

func drawPhrases(gooey *gocui.Gui) error {

	// Get window size.
	maxX, maxY := gooey.Size()

	_, _, _, minY, err := gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	if phraseHead, err := gooey.SetView("phraseHead", 0, minY, maxX/6, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(phraseHead, "Phrases")
	}

	//
	if phraseView, err := gooey.SetView("phrases", 0, minY+2, maxX/6, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		phraseView.SelBgColor = gocui.ColorCyan
		phraseView.Highlight = true

		for _, exp := range *exps {
			phrase := padText(exp.Phrase, maxX)
			fmt.Fprintln(phraseView, phrase)
		}
	}

	return nil

}

func drawExpansions(gooey *gocui.Gui) error {

	// Get window size.
	maxX, maxY := gooey.Size()

	_, _, _, minY, err := gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	_, _, minX, _, err := gooey.ViewPosition("phrases")
	if err != nil {
		log.Fatal(err)
	}

	if expHead, err := gooey.SetView("expHead", minX, minY, maxX-1, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(expHead, "Expansions")
	}

	if expView, err := gooey.SetView("galaxy", minX, minY+2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(expView, "Long ago, in a galaxy far, far away, there was no need for overly complicated programming languages designed to make simple tasks easier for users and more difficult for programmers.  In fact, there were no users, programmers, or computers.  There was nothing, in fact, because the galaxy did not contain life.")
		expView.Wrap = true
	}

	return nil
}

func mapPhrases(gooey *gocui.Gui) error {
	drawHeader(gooey)
	drawPhrases(gooey)
	drawExpansions(gooey)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func readExpanders() {
	exps = ggdb.ReadExpanders()
}

// GengarMenu creates an example GUI.
func GengarMenu() {

	readExpanders()

	gooey := gocui.NewGui()
	if err := gooey.Init(); err != nil {
		log.Fatal(err)
	}
	defer gooey.Close()

	gooey.SetLayout(mapPhrases)

	if err := gooey.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
