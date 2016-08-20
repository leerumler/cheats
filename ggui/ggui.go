package ggui

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/leerumler/gengar/ggdb"
)

type curView struct {
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
func drawHeader(menu *curView) error {

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

func drawCategories(menu *curView) error {

	// Get window size.
	// maxX, maxY := gooey.Size()

	_, _, _, minY, err := menu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	if catHead, err := menu.gooey.SetView("catHead", 0, minY, menu.maxX/6, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(catHead, "Categories")
	}

	if catView, err := menu.gooey.SetView("categories", 0, minY+2, menu.maxX/6, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		catView.SelBgColor = gocui.ColorCyan
		catView.Highlight = true

		cats := ggdb.ReadCategories()

		for _, cat := range *cats {
			category := padText(cat.Name, menu.maxX)
			fmt.Fprintln(catView, category)
		}

		_, posY := catView.Cursor()
		curCatName, _ := catView.Line(posY)
		curCatName = strings.TrimSpace(curCatName)
		var curCat ggdb.Category
		for _, cat := range *cats {
			if curCatName == cat.Name {
				curCat = cat
			}
		}
		// fmt.Println(curCat.ID, curCat.Name)

		menu.cat = curCat

	}
	return nil
}

func drawExpansions(menu *curView) error {

	// Get window size.
	// maxX, maxY := gooey.Size()

	_, _, _, minY, err := menu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	_, _, minX, _, err := menu.gooey.ViewPosition("categories")
	if err != nil {
		log.Fatal(err)
	}

	if expHead, err := menu.gooey.SetView("expHead", minX, minY, (menu.maxX/6)*5, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(expHead, "Expansions")
	}

	if expView, err := menu.gooey.SetView("expansions", minX, minY+2, (menu.maxX/6)*5, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		exps := ggdb.ReadExpansionsInCategory(menu.cat)

		for _, exp := range *exps {
			expName := padText(exp.Name, menu.maxX)
			fmt.Fprintln(expView, expName)
		}

	}

	return nil
}

func drawPhrases(menu *curView) error {

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

func ggMenu(gooey *gocui.Gui) error {

	var menu curView
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

	gooey.SetLayout(ggMenu)

	if err := gooey.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
