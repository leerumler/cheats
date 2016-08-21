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

// drawCategories creates the categories view, with a header.
func drawCategories(menu *ggMenu) error {

	// Find minY, which will be the bottom of the header view.
	_, _, _, minY, err := menu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	// Create a view holding the category header.
	if catHead, err := menu.gooey.SetView("catHead", 0, minY, menu.maxX/6, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(catHead, "Categories")
	}

	// Create a new for the categories themselves.
	if catView, err := menu.gooey.SetView("categories", 0, minY+2, menu.maxX/6, menu.maxY-1); err != nil {
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
			category := padText(cat.Name, menu.maxX)
			fmt.Fprintln(catView, category)
		}

		// var curCat ggdb.Category
		// for _, cat := range cats {
		// 	if curCatName == cat.Name {
		// 		curCat = cat
		// 	}
		// }

		// fmt.Println(curCat.ID, curCat.Name)

	}
	return nil
}

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

	// Create the expansions view, which will list the expansions within the currently selected category.
	if expView, err := menu.gooey.SetView("expansions", minX, minY+2, (menu.maxX/6)*5, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Get the expansions from the database and set the currently
		// selected expansion to the first item in the list.
		exps := ggdb.ReadExpansionsInCategory(menu.cat)
		menu.exp = exps[0]

		// Print the expansions to the view.
		for _, exp := range exps {
			expName := padText(exp.Name, menu.maxX)
			fmt.Fprintln(expView, expName)
		}
	}

	return nil
}

func drawPhrases(menu *ggMenu) error {

	_, _, _, minY, err := menu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	_, _, minX, _, err := menu.gooey.ViewPosition("expansions")
	if err != nil {
		log.Fatal(err)
	}

	if phraseHead, err := menu.gooey.SetView("phraseHead", minX, minY, menu.maxX, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(phraseHead, "Phrases")
	}

	//
	if phraseView, err := menu.gooey.SetView("phrases", minX, minY+2, menu.maxX, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		phraseView.SelBgColor = gocui.ColorCyan
		phraseView.Highlight = true

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

	var menu ggMenu
	menu.gooey = gooey
	menu.maxX, menu.maxY = menu.gooey.Size()

	drawHeader(&menu)
	drawCategories(&menu)
	drawExpansions(&menu)
	drawPhrases(&menu)

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// GengarMenu creates an example GUI.
func GengarMenu() {

	// readExpanders()

	gooey := gocui.NewGui()
	if err := gooey.Init(); err != nil {
		log.Fatal(err)
	}
	defer gooey.Close()

	gooey.SetLayout(runMenu)

	if err := gooey.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
