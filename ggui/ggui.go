package ggui

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/leerumler/gengar/ggdb"
)

type ggMenu struct {
	cat        *ggdb.Category
	exp        *ggdb.Expansion
	phrase     ggdb.Phrase
	gooey      *gocui.Gui
	maxX, maxY int
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

// readSel reads the currently selected line and returns a string
// containing its contents, without trailing spaces.
func readSel(curView *gocui.View) string {

	_, posY := curView.Cursor()
	selection, _ := curView.Line(posY)
	selection = strings.TrimSpace(selection)

	return selection
}

// readCat reads the currently selected category name and matches it to a ggdb.Category.
func readCat(catView *gocui.View, cats []ggdb.Category) *ggdb.Category {

	// Read the name of the currently selected category.
	curCatName := readSel(catView)

	// Search for a category that matches that name.
	var curCat ggdb.Category
	for _, cat := range cats {
		if curCatName == cat.Name {
			curCat = cat
		}
	}

	// And return it.
	return &curCat
}

// readExp reads the currently selected expansion name and matches it to a ggdb.Expansion.
func readExp(expView *gocui.View, exps []ggdb.Expansion) *ggdb.Expansion {

	// Read the name of the currently selected expansion.
	curExpName := readSel(expView)

	// Search for an expansion that matches that name.
	var curExp ggdb.Expansion
	for _, exp := range exps {
		if curExpName == exp.Name {
			curExp = exp
		}
	}

	// And return it.
	return &curExp
}

// selUp moves the cursor/selection up one line.
func selUp(gooey *gocui.Gui, view *gocui.View) error {

	if view != nil {
		view.MoveCursor(0, -1, false)
	}

	return nil
}

// selDown moves the selected menu item down one line, without moving past the last line.
func selDown(gooey *gocui.Gui, view *gocui.View) error {

	if view != nil {
		view.MoveCursor(0, 1, false)

		// If the cursor moves to an empty line, move it back. :P
		if readSel(view) == "" {
			view.MoveCursor(0, -1, false)
		}
	}

	return nil
}

// resetHighlights changes the SelBgColor of the categories, expansions,
// and phrases views back to their "default" blue.
func resetHighlights(gooey *gocui.Gui) error {

	if catView, err := gooey.View("categories"); err == nil {
		catView.SelBgColor = gocui.ColorBlue
	} else {
		return err
	}

	if expView, err := gooey.View("expansions"); err == nil {
		expView.SelBgColor = gocui.ColorBlue
	} else {
		return err
	}

	if phraseView, err := gooey.View("phrases"); err == nil {
		phraseView.SelBgColor = gocui.ColorBlue
	} else {
		return err
	}

	return nil
}

// focusCat changes the focus to the categories view.
func focusCat(gooey *gocui.Gui, view *gocui.View) error {

	// Focus on the categories view.
	if err := gooey.SetCurrentView("categories"); err != nil {
		log.Fatal(err)
	}

	// Reset every view's highlight colors back to blue, then set the
	// set the categories view's highlight color to cyan.
	// So everyone's blue but categories.
	resetHighlights(gooey)
	if catView, err := gooey.View("categories"); err == nil {
		catView.SelBgColor = gocui.ColorCyan
	} else {
		log.Fatal(err)
	}

	return nil
}

// focusExp changes the focus to the expansions view.
func focusExp(gooey *gocui.Gui, view *gocui.View) error {

	// Focus on the epxansions view.
	if err := gooey.SetCurrentView("expansions"); err != nil {
		return err
	}

	// Reset every view's highlight colors back to blue, then set the
	// expansions view's highlight color to cyan.
	// So everyone's blue but expansions.
	resetHighlights(gooey)
	if expView, err := gooey.View("expansions"); err == nil {
		expView.SelBgColor = gocui.ColorCyan
	} else {
		return err
	}

	return nil
}

// focusPhrase changes the focus to the phrases view.  It also sets
// the highlight color to Cyan for clarity.
func focusPhrase(gooey *gocui.Gui, view *gocui.View) error {

	// Focus on the phrases view.
	if err := gooey.SetCurrentView("phrases"); err != nil {
		log.Fatal(err)
	}

	// Reset every view's highlight colors back to blue, then set the
	// set the phrases view's highlight color to cyan.
	// So everyone's blue but phrases.
	resetHighlights(gooey)
	if phraseView, err := gooey.View("phrases"); err == nil {
		phraseView.SelBgColor = gocui.ColorCyan
	} else {
		log.Fatal(err)
	}

	return nil
}

// setKeyBinds is a necessary evil.
func setKeyBinds(gooey *gocui.Gui) error {

	// Always have an exit strategy.
	if err := gooey.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// If the categories view is focused and ↑ is pressed, move the selected menu item up one.
	if err := gooey.SetKeybinding("categories", gocui.KeyArrowUp, gocui.ModNone, selUp); err != nil {
		return err
	}

	// If the categories view is focused and ↓ is pressed, move the selected menu item down one.
	if err := gooey.SetKeybinding("categories", gocui.KeyArrowDown, gocui.ModNone, selDown); err != nil {
		return err
	}

	// If the categories view is focused and → is pressed, move the focus to the expansions view.
	if err := gooey.SetKeybinding("categories", gocui.KeyArrowRight, gocui.ModNone, focusExp); err != nil {
		return err
	}

	// If the categories view is focused and Enter is pressed, move the focus to the expansions view.
	if err := gooey.SetKeybinding("categories", gocui.KeyEnter, gocui.ModNone, focusExp); err != nil {
		return err
	}

	// If the expansions view is focused and ↑ is pressed, move the selected menu item up one.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowUp, gocui.ModNone, selUp); err != nil {
		return err
	}

	// If the expansions view is focused and ↓ is pressed, move the selected menu item down one.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowDown, gocui.ModNone, selDown); err != nil {
		return err
	}

	// If the expansions view is focused and ← is pressed, move the focus to the categories menu.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowLeft, gocui.ModNone, focusCat); err != nil {
		return err
	}

	// If the expansions view is focused and → is pressed, move the focus to the phrases view.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowRight, gocui.ModNone, focusPhrase); err != nil {
		return err
	}

	// If the phrases view is focused and ↑ is pressed, move the selected menu item up one.
	if err := gooey.SetKeybinding("phrases", gocui.KeyArrowUp, gocui.ModNone, selUp); err != nil {
		return err
	}

	// If the phrases view is focused and ↓ is pressed, move the selected menu item down one.
	if err := gooey.SetKeybinding("phrases", gocui.KeyArrowDown, gocui.ModNone, selDown); err != nil {
		return err
	}

	// If the phrases view is focused and ← is pressed, move the focus to the expansions menu.
	if err := gooey.SetKeybinding("phrases", gocui.KeyArrowLeft, gocui.ModNone, focusExp); err != nil {
		return err
	}

	return nil
}

// drawHeader adds the "Gengar Configuration Editor" header to the top of the menu.
func drawHeader(menu *ggMenu) error {

	// Place the header view at the top of the menu and extend it down two pixels.
	if header, err := menu.gooey.SetView("header", 0, 0, menu.maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Center the header text and print it to the view.
		head := centerText("Gengar Configuration Editor", menu.maxX)
		fmt.Fprintln(header, head)
	}

	return nil

}

// drawCategories draws the categories menu, which displays all available categories.
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

	// Create a view for the categories if it doesn't already exist.
	if catView, err := menu.gooey.SetView("categories", 0, minY+2, menu.maxX/6, menu.maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Since the categories view will default to being
		// the default view, it defaults to cyan.
		catView.SelBgColor = gocui.ColorCyan
		catView.Highlight = true
	}

	// Redraw the categories menu on the categories view.
	if catView, err := menu.gooey.View("categories"); err == nil {

		// Clear the internal buffer, since this allows us to check
		// for more categories every time the function is executed.
		catView.Clear()

		// Read the categories from the database.
		cats := ggdb.ReadCategories()

		// Print the name of each category in rows on the menu.
		for _, cat := range cats {
			category := padText(cat.Name, menu.maxX/6)
			fmt.Fprintln(catView, category)
		}

		// Check the currently selected row and store its matching category.
		menu.cat = readCat(catView, cats)

	} else {
		log.Fatal(err)
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

	// Create the expansion view, if it doesn't already exist.
	if expView, err := menu.gooey.SetView("expansions", minX, minY+2, (menu.maxX/6)*5, menu.maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// The currently selected expansion will default to Blue.
		expView.SelBgColor = gocui.ColorBlue
		expView.Highlight = true
	}

	// Redraw the expansions menu on the expansions view.
	if expView, err := menu.gooey.View("expansions"); err == nil {

		// Clear the internal buffer.
		expView.Clear()

		// Read the expansions from the database.
		exps := ggdb.ReadExpansions(menu.cat)

		// Print name of each expansion to the view.
		for _, exp := range exps {
			expName := padText(exp.Name, menu.maxX/6*4)
			fmt.Fprintln(expView, expName)
		}

		// Check the currently selected row and store its matching expansion.
		menu.exp = readExp(expView, exps)

	} else {
		return err
	}

	return nil
}

// drawPhrases creates the phrases view, which displays all of the phrases
// mapped to the currently selected expansion.
func drawPhrases(menu *ggMenu) error {

	// Find the lowest coordinate of the header view, which will serve as
	// the top coordinate on the header view.
	_, _, _, minY, err := menu.gooey.ViewPosition("header")
	if err != nil {
		log.Fatal(err)
	}

	// Find minX, which is the right side of the expansions view.
	_, _, minX, _, err := menu.gooey.ViewPosition("expansions")
	if err != nil {
		log.Fatal(err)
	}

	// Create a view for the phrases header and print the title.
	if phraseHead, err := menu.gooey.SetView("phraseHead", minX, minY, menu.maxX-1, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(phraseHead, "Phrases")
	}

	// Create the phrases view, listing the phrases that are mapped to the currently selected expansion.
	if phraseView, err := menu.gooey.SetView("phrases", minX, minY+2, menu.maxX-1, menu.maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Selected phrase will be highlighted in Blue.
		phraseView.SelBgColor = gocui.ColorBlue
		phraseView.Highlight = true
	}

	// Draw the phrase menu on to the phrase view.
	if phraseView, err := menu.gooey.View("phrases"); err == nil {

		// Empty the phrase view of all contents.
		phraseView.Clear()

		// Read the phrases from the database.
		phrases := ggdb.ReadPhrases(menu.exp)

		// Print each of the phrases to the view.
		for _, phrase := range phrases {
			phraseText := padText(phrase.Phrase, menu.maxX/6)
			fmt.Fprintln(phraseView, phraseText)
		}

	} else {
		return err
	}

	return nil
}

func drawHelp(menu *ggMenu) error {

	if help, err := menu.gooey.SetView("help", 0, menu.maxY-3, menu.maxX-1, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(help, "Help text goes here")
	}

	return nil
}

func runMenu(gooey *gocui.Gui) error {

	// Create a ggMenu variable and populate it with some basic info.
	var menu ggMenu
	menu.gooey = gooey
	menu.maxX, menu.maxY = menu.gooey.Size()

	// Draw the views in the menu.
	if err := drawHeader(&menu); err != nil {
		log.Panicln(err)
	}
	if err := drawCategories(&menu); err != nil {
		log.Panicln(err)
	}
	if err := drawExpansions(&menu); err != nil {
		log.Panicln(err)
	}
	if err := drawPhrases(&menu); err != nil {
		log.Panicln(err)
	}

	// If the current view isn't set, set it to categories.
	if gooey.CurrentView() == nil {
		if err := gooey.SetCurrentView("categories"); err != nil {
			log.Fatal(err)
		}
	}

	// And finally, draw the help menu.
	if err := drawHelp(&menu); err != nil {
		log.Panicln(err)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// GengarMenu creates an example GUI.
func GengarMenu() {

	//
	gooey := gocui.NewGui()
	if err := gooey.Init(); err != nil {
		log.Panicln(err)
	}
	defer gooey.Close()

	// Set the layout handler.
	gooey.SetLayout(runMenu)

	//
	if err := setKeyBinds(gooey); err != nil {
		log.Panicln(err)
	}

	// Start the main event loop.
	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
