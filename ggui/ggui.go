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
	phrase     *ggdb.Phrase
	cats       []ggdb.Category
	exps       []ggdb.Expansion
	phrases    []ggdb.Phrase
	gooey      *gocui.Gui
	maxX, maxY int
}

var menu ggMenu

// quit quits the main event loop.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// centerText takes a string of text and a length and pads the beginning
// of the string with spaces to center that text in the available space.
func centerText(text *string, maxX int) *string {

	numSpaces := maxX/2 - len(*text)/2
	for i := 1; i < numSpaces; i++ {
		*text = " " + *text
	}

	return text
}

// padText takes a string of text and pads the end of it with spaces to
// fill the available space in a cell.
func padText(text *string, maxX int) *string {

	numSpaces := maxX - len(*text)
	for i := 0; i < numSpaces; i++ {
		*text += " "
	}

	return text
}

// emptyLine returns a string of spaces of the specified length.
func emptyLine(length int) *string {

	var blank string
	for i := 0; i < length; i++ {
		blank += " "
	}

	return &blank
}

// readSel reads the currently selected line and returns a string
// containing its contents, without trailing spaces.
func readSel(view *gocui.View) string {

	_, posY := view.Cursor()
	selection, _ := view.Line(posY)
	selection = strings.TrimSpace(selection)

	return selection
}

// readCat reads the currently selected category name and matches it to a ggdb.Category.
func readCat() *ggdb.Category {

	var curCat ggdb.Category

	catView, err := menu.gooey.View("categories")
	if err != nil {
		return &curCat
	}

	// Read the name of the currently selected category.
	curCatName := readSel(catView)

	// Search for a category that matches that name.
	for _, cat := range menu.cats {
		if curCatName == cat.Name {
			curCat = cat
		}
	}

	// And return it.
	return &curCat
}

// readExp reads the currently selected expansion name and matches it to a ggdb.Expansion.
func readExp() *ggdb.Expansion {

	var curExp ggdb.Expansion
	expView, err := menu.gooey.View("expansions")
	if err != nil {
		return &curExp
	}

	// Read the name of the currently selected expansion.
	curExpName := readSel(expView)

	// Search for an expansion that matches that name.
	for _, exp := range menu.exps {
		if curExpName == exp.Name {
			curExp = exp
		}
	}

	// And return it.
	return &curExp
}

// readPhrase reads the currently selected phrase name and matches it to a ggdb.Phrase.
func readPhrase() *ggdb.Phrase {

	var curPhrase ggdb.Phrase
	phraseView, err := menu.gooey.View("phrases")
	if err != nil {
		return &curPhrase
	}

	// Read the name of the currently selected phrase.
	curPhraseName := readSel(phraseView)

	// Search for a phrase that matches that name.
	for _, phrase := range menu.phrases {
		if curPhraseName == phrase.Name {
			curPhrase = phrase
		}
	}

	// And return it.
	return &curPhrase
}

// drawHeader adds the "Gengar Configuration Editor" header to the top of the menu.
func drawHeader() error {

	// Place the header view at the top of the menu and extend it down two pixels.
	if header, err := menu.gooey.SetView("header", 0, 0, menu.maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Center the header text and print it to the view.
		ggHead := "Gengar Configuration Editor"
		fmt.Fprintln(header, *centerText(&ggHead, menu.maxX))
	}

	return nil

}

// drawCategories draws the categories menu, which displays all available categories.
func drawCategories() error {

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
	if catView, err := menu.gooey.SetView("categories", 0, minY+2, menu.maxX/6, menu.maxY*2/3); err != nil {
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
		menu.cats = ggdb.ReadCategories()

		// If the category list is empty, print a blank line.
		if len(menu.cats) == 0 {
			fmt.Fprintln(catView, *emptyLine(menu.maxX / 6))
		}

		// Print the name of each category in rows on the menu.
		for _, cat := range menu.cats {
			fmt.Fprintln(catView, *padText(&cat.Name, menu.maxX/6))
		}

		// Check the currently selected row and store its matching category.
		menu.cat = readCat()

	} else {
		return err
	}

	return nil
}

// drawExpansions creates the expansions view, which displays all of the expansions
// in the currently selected category.
func drawExpansions() error {

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
	if expHead, err := menu.gooey.SetView("expHead", minX, minY, menu.maxX*5/6, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(expHead, "Expansions")
	}

	// Create the expansion view, if it doesn't already exist.
	if expView, err := menu.gooey.SetView("expansions", minX, minY+2, menu.maxX*5/6, menu.maxY*2/3); err != nil {
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
		menu.exps = ggdb.ReadExpansions(menu.cat)

		// If the expansions list is empty, print a blank line.
		if len(menu.exps) == 0 {
			fmt.Fprintln(expView, *emptyLine(menu.maxX * 4 / 6))
		}

		// Print name of each expansion to the view.
		for _, exp := range menu.exps {
			fmt.Fprintln(expView, *padText(&exp.Name, menu.maxX*4/6))
		}

		// Check the currently selected row and store its matching expansion.
		menu.exp = readExp()

	} else {
		return err
	}

	return nil
}

// drawPhrases creates the phrases view, which displays all of the phrases
// mapped to the currently selected expansion.
func drawPhrases() error {

	// Find the lowest coordinate of the header view, which will serve as
	// the top coordinate on the header view.
	_, _, _, minY, err := menu.gooey.ViewPosition("header")
	if err != nil {
		return err
	}

	// Find minX, which is the right side of the expansions view.
	_, _, minX, _, err := menu.gooey.ViewPosition("expansions")
	if err != nil {
		return err
	}

	// Create a view for the phrases header and print the title.
	if phraseHead, err := menu.gooey.SetView("phraseHead", minX, minY, menu.maxX-1, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(phraseHead, "Phrases")
	}

	// Create the phrases view, listing the phrases that are mapped to the currently selected expansion.
	if phraseView, err := menu.gooey.SetView("phrases", minX, minY+2, menu.maxX-1, menu.maxY*2/3); err != nil {
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
		menu.phrases = ggdb.ReadPhrases(menu.exp)

		// If the expansions list is empty, print a blank line.
		if len(menu.phrases) == 0 {
			fmt.Fprintln(phraseView, *emptyLine(menu.maxX / 6))
		}

		// Print each of the phrases to the view.
		for _, phrase := range menu.phrases {
			fmt.Fprintln(phraseView, *padText(&phrase.Name, menu.maxX/6))
		}

		menu.phrase = readPhrase()

	} else {
		return err
	}

	return nil
}

// drawHelp creates and updates the help view, which displays the keybindings/controls
// for the currently selected view.
func drawHelp() error {

	// Find minY, which will be the bottom of the expansions view.
	_, _, _, minY, err := menu.gooey.ViewPosition("expansions")
	if err != nil {
		return err
	}

	// Create a view to hold the help menu.
	if _, err := menu.gooey.SetView("help", 0, minY, menu.maxX-1, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	// Redraw the help menu on the help view.
	if help, err := menu.gooey.View("help"); err == nil {

		// Check if the current view is set.
		if curView := menu.gooey.CurrentView(); curView != nil {

			// Clear the view.
			help.Clear()

			// Print a help message that relates to the currently selected view.
			var helpText string
			switch curView.Name() {
			case "categories":
				helpText = "select: ⏎ | new: ctrl+n | edit: ctrl+e | delete: ctrl+d | quit: ctrl+q"
			case "expansions":
				helpText = "select: ⏎ | new: ctrl+n | edit: ctrl+e | delete: ctrl+d | quit: ctrl+q"
			case "phrases":
				helpText = "new: ctrl+n | edit: ctrl+e | delete: ctrl+d | quit: ctrl+q"
			case "text":
				helpText = "save: ctrl+s | reload: ctrl+r | exit: ctrl+x | quit: ctrl+q "
			case "newCatPrompt":
				helpText = "save: ctrl+s | exit: ctrl+x | quit: ctrl+q "
			case "upCatPrompt":
				helpText = "save: ctrl+s | exit: ctrl+x | quit: ctrl+q "
			case "newExpPrompt":
				helpText = "save: ctrl+s | exit: ctrl+x | quit: ctrl+q "
			case "upExpPrompt":
				helpText = "save: ctrl+s | exit: ctrl+x | quit: ctrl+q "
			case "newPhrasePrompt":
				helpText = "save: ctrl+s | exit: ctrl+x | quit: ctrl+q "
			case "upPhrasePrompt":
				helpText = "save: ctrl+s | exit: ctrl+x | quit: ctrl+q "
			}
			helpText = *centerText(&helpText, menu.maxX)
			fmt.Fprintln(help, helpText)
		}

	} else {
		return err
	}

	return nil
}

// drawText creates the text view, which displays the text of the currently selected expansion.
// The text view allows users to edit and update the expansion text.
func drawText() error {

	// Find minY, which will be the bottom of the help view.
	_, _, _, minY, err := menu.gooey.ViewPosition("help")
	if err != nil {
		log.Fatal(err)
	}

	// Create the text view if it doesn't already exist.
	if textView, err := menu.gooey.SetView("text", 0, minY, menu.maxX-1, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Enable editing and line wrapping.
		textView.Editable = true
		textView.Wrap = true

	} else {
		return err
	}

	return nil
}

// upText updates the text view with the text from the currently selected expansion.
// This is separated from drawText to prevent unwanted updates.
func upText() error {

	// Ensure the currently selected expansion is up-to-date.
	// menu.exps = ggdb.ReadExpansions(menu.cat)
	// menu.exp = readExp()

	// Clear the text view and re-fill it with the current expansion text.
	if textView, err := menu.gooey.View("text"); err == nil {
		textView.Clear()
		textView.SetCursor(0, 0)
		fmt.Fprintln(textView, menu.exp.Text)
	} else {
		return err
	}

	return nil
}

// runMenu is ggui's main layout handler.  It draws each of ggui's main views.
func runMenu(gooey *gocui.Gui) error {

	// Create a ggMenu variable and populate it with some basic info.
	menu.gooey = gooey
	menu.maxX, menu.maxY = menu.gooey.Size()

	// Draw the views in the menu.
	if err := drawHeader(); err != nil {
		return err
	}
	if err := drawCategories(); err != nil {
		return err
	}
	if err := drawExpansions(); err != nil {
		return err
	}
	if err := drawPhrases(); err != nil {
		return err
	}
	if err := drawHelp(); err != nil {
		return err
	}
	if err := drawText(); err != nil {
		return err
	}

	// If the current view isn't set, set it to categories.
	if gooey.CurrentView() == nil {
		if _, err := gooey.SetCurrentView("categories"); err != nil {
			return err
		}
		if err := upText(); err != nil {
			return err
		}
	}

	return nil
}

// GengarMenu creates a text-based menu to control Gengar's settings.
func GengarMenu() {

	// Create a new gocui.Gui and initialize it.
	gooey, err := gocui.NewGui()
	if err != nil {
		log.Panicln(err)
	}
	defer gooey.Close()

	// Set the layout handler.
	gooey.SetManagerFunc(runMenu)

	// Always have an exit strategy.
	if err := gooey.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// Set the rest of the keybindinds.
	if err := setKeyBinds(gooey); err != nil {
		log.Panicln(err)
	}

	// Start the main event loop.
	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
