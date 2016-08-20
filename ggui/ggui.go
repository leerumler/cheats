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
func drawHeader(ggmenu *curView) error {

	// The header will be dynamically placed at the top of the menu and will extend down two pixels.
	if header, err := ggmenu.gooey.SetView("header", 0, 0, ggmenu.maxX, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// The text will be centered at the top of the menu.
		head := centerText("Gengar Configuration Editor", ggmenu.maxX)
		fmt.Fprintln(header, head)
	}

	return nil

}

func drawCategories(ggmenu *curView) error {

	// Get window size.
	// maxX, maxY := gooey.Size()

	_, _, _, minY, err := ggmenu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	if catHead, err := ggmenu.gooey.SetView("catHead", 0, minY, ggmenu.maxX/6, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(catHead, "Categories")
	}

	if catView, err := ggmenu.gooey.SetView("categories", 0, minY+2, ggmenu.maxX/6, ggmenu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		catView.SelBgColor = gocui.ColorCyan
		catView.Highlight = true

		cats := ggdb.ReadCategories()

		for _, cat := range *cats {
			category := padText(cat.Name, ggmenu.maxX)
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

		ggmenu.cat = curCat

	}
	return nil
}

func drawExpansions(ggmenu *curView) error {

	// Get window size.
	// maxX, maxY := gooey.Size()

	_, _, _, minY, err := ggmenu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	_, _, minX, _, err := ggmenu.gooey.ViewPosition("categories")
	if err != nil {
		log.Fatal(err)
	}

	if expHead, err := ggmenu.gooey.SetView("expHead", minX, minY, (ggmenu.maxX/6)*5, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(expHead, "Expansions")
	}

	if expView, err := ggmenu.gooey.SetView("expansions", minX, minY+2, (ggmenu.maxX/6)*5, ggmenu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		exps := ggdb.ReadExpansionsInCategory(ggmenu.cat)

		for _, exp := range *exps {
			expName := padText(exp.Name, ggmenu.maxX)
			fmt.Fprintln(expView, expName)
		}

	}

	return nil
}

func drawPhrases(ggmenu *curView) error {

	_, _, _, minY, err := ggmenu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	_, _, minX, _, err := ggmenu.gooey.ViewPosition("expansions")
	if err != nil {
		log.Fatal(err)
	}

	if phraseHead, err := ggmenu.gooey.SetView("phraseHead", minX, minY, ggmenu.maxX, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(phraseHead, "Phrases")
	}

	//
	if phraseView, err := ggmenu.gooey.SetView("phrases", minX, minY+2, ggmenu.maxX, ggmenu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		phraseView.SelBgColor = gocui.ColorCyan
		phraseView.Highlight = true

		// for _, exp := range *exps {
		// 	phrase := padText(exp.Phrase, maxX)
		// 	fmt.Fprintln(phraseView, phrase)
		// }
	}

	return nil

}

func catView(gooey *gocui.Gui) error {

	var ggmenu curView
	ggmenu.gooey = gooey
	ggmenu.maxX, ggmenu.maxY = ggmenu.gooey.Size()

	drawHeader(&ggmenu)
	drawCategories(&ggmenu)
	drawExpansions(&ggmenu)
	drawPhrases(&ggmenu)
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

	gooey.SetLayout(catView)

	if err := gooey.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
