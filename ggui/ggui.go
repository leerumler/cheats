package ggui

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/leerumler/gengar/ggdb"
)

type ggMenu struct {
	cat        ggdb.Category
	exp        ggdb.Expansion
	phrase     ggdb.Phrase
	gooey      *gocui.Gui
	maxX, maxY int
}

// getGooey returns a new GUI struct.
func getGooey() {

	// Create a new GUI.
	gooey := gocui.NewGui()
	if err := gooey.Init(); err != nil {
		log.Fatal(err)
	}
	defer gooey.Close()
}

// centerText takes a string of text and a length and pads the beginning
// of the string with spaces to center that text in the available space.
func centerText(text string, maxX int) string {
	numSpaces := maxX/2 - len(text)/2
	for i := 0; i < numSpaces; i++ {
		text = " " + text
	}
	return text
}

// padText takes a string of text and pads the end of it with spaces to
// fill the available space in a cell.
func padText(text string, maxX int) string {
	numSpaces := maxX - len(text)
	for i := 0; i < numSpaces; i++ {
		text = text + " "
	}
	return text
}

// drawHeader adds the "Gengar Configuration Editor" header to the top of the menu.
func drawHeader(menu *ggMenu) error {

	// The header will be dynamically placed at the top of the menu and will extend down two pixels.
	if header, err := menu.gooey.SetView("header", 0, 0, menu.maxX, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// The text will be centered at the top of the menu.
		head := centerText("Gengar Configuration Editor", menu.maxX)
		fmt.Fprintln(header, head)
	}

	return nil

}

// drawCategories creates the categories view, which displays all available categories.
func drawCategories(menu *ggMenu) error {

	// Find minY, which will be the bottom of the header view.
	_, _, _, minY, err := menu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	// Set maxX coordinate as 1/6 of the menu size.
	maxX := menu.maxX / 6

	// Create a view holding the category header.
	if catHead, err := menu.gooey.SetView("catHead", 0, minY, maxX, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(catHead, "Categories")
	}

	// Create a new for the categories themselves.
	if catView, err := menu.gooey.SetView("categories", 0, minY+2, maxX, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Selected category will be highlighted in Cyan.
		catView.SelBgColor = gocui.ColorCyan
		catView.Highlight = true

		// Read the categories from the database and set the currently
		// selected category to the first category in the list.
		cats := ggdb.ReadCategories()
		menu.cat = cats[0]

		// Print each of the categories to the view.
		for _, cat := range cats {
			category := padText(cat.Name, maxX)
			fmt.Fprintln(catView, category)
		}

	}
	return nil
}

// drawExpansions creates the expansions view, which displays all of the expansions
// in the currently selected category.
func drawExpansions(menu *ggMenu) error {

	// Find minY, which will be the bottom of the header view.
	_, _, _, minY, err := menu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	// Find minX, which will be the right side of categories view.
	_, _, minX, _, err := menu.gooey.ViewPosition("categories")
	if err != nil {
		log.Fatal(err)
	}

	// Create a view for the expansions header.
	if expHead, err := menu.gooey.SetView("expHead", minX, minY, (menu.maxX/6)*5, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(expHead, "Expansions")
	}

	// Create the expansions view, listing the expansions within the currently selected category.
	if expView, err := menu.gooey.SetView("expansions", minX, minY+2, (menu.maxX/6)*5, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Selected expansion will be highlighted in Cyan.
		expView.SelBgColor = gocui.ColorCyan
		expView.Highlight = true

		// Read the expansions from the database and set the currently
		// selected expansion to the first item in the list.
		exps := ggdb.ReadExpansions(menu.cat)
		menu.exp = exps[0]

		// Print the expansions to the view.
		for _, exp := range exps {
			expName := padText(exp.Name, menu.maxX/6*4)
			fmt.Fprintln(expView, expName)
		}
	}

	return nil
}

// drawPhrases creates the phrases view, which displays all of the phrases
// mapped to the currently selected expansion.
func drawPhrases(menu *ggMenu) error {

	// Find minY, which is the bottom of the header view.
	_, _, _, minY, err := menu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	// Find minX, which is the right side of the expansions view.
	_, _, minX, _, err := menu.gooey.ViewPosition("expansions")
	if err != nil {
		log.Fatal(err)
	}

	// Create a view for the phrases header.
	if phraseHead, err := menu.gooey.SetView("phraseHead", minX, minY, menu.maxX, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(phraseHead, "Phrases")
	}

	// Create the phrases view, listing the phrases mapped to the currently selected expansion.
	if phraseView, err := menu.gooey.SetView("phrases", minX, minY+2, menu.maxX, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Selected phrase will be highlighted in Cyan.
		phraseView.SelBgColor = gocui.ColorCyan
		phraseView.Highlight = true

		// Read the phrases from the database and set the currently selected
		// phrase to the first item in the list.
		phrases := ggdb.ReadPhrases(menu.exp)
		menu.phrase = phrases[0]

		// Print the phrases to the view.
		for _, phrase := range phrases {
			phraseText := padText(phrase.Phrase, menu.maxX/6)
			fmt.Fprintln(phraseView, phraseText)
		}

	}

	return nil

}

// readSel reads the currently selected line and returns a string
// containing its contents, without trailing spaces.
func readSel(curView *gocui.View) string {
	_, posY := curView.Cursor()
	selection, _ := curView.Line(posY)
	selection = strings.TrimSpace(selection)
	return selection
}

func runMenu(gooey *gocui.Gui) error {

	// Create a ggMenu variable and populate it with some basic info.
	var menu ggMenu
	menu.gooey = gooey
	menu.maxX, menu.maxY = menu.gooey.Size()

	// Draw the views in the menu.
	drawHeader(&menu)
	drawCategories(&menu)
	drawExpansions(&menu)
	drawPhrases(&menu)

	// Something like this could be used in the keybindings:
	// var curCat ggdb.Category
	// for _, cat := range cats {
	// 	if curCatName == cat.Name {
	// 		curCat = cat
	// 	}
	// }

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// GengarMenu creates an example GUI.
func GengarMenu() {

	// Create a new gocui.Gui object.
	gooey := gocui.NewGui()
	if err := gooey.Init(); err != nil {
		log.Fatal(err)
	}
	defer gooey.Close()

	// Set the layout handler.
	gooey.SetLayout(runMenu)

	// Make Ctrl+C quit the menu.
	if err := gooey.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// Start the main loop.
	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
