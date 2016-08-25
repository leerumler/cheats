package ggui

import (
	"fmt"
	"log"

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

var menu ggMenu

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
		*text = *text + " "
	}

	return text
}

func popUp(message *string) error {

	minX := menu.maxX * 1 / 4
	maxX := menu.maxX * 3 / 4
	minY := menu.maxY/2 - 1
	maxY := menu.maxY/2 + 1

	if popUp, err := menu.gooey.SetView("popUp", minX, minY, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(popUp, *centerText(message, maxX-minX))
	}

	return nil
}

func promptInput(title *string, prev *string) error {

	//
	minX := menu.maxX * 1 / 4
	maxX := menu.maxX * 3 / 4
	midY := menu.maxY / 2

	if promptHead, err := menu.gooey.SetView("promptHead", minX, midY-2, maxX, midY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(promptHead, *centerText(title, menu.maxX/2))
	}

	if prompt, err := menu.gooey.SetView("prompt", minX, midY, maxX, midY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		menu.gooey.Editor = gocui.EditorFunc(singleLineEditor)
		menu.gooey.Cursor = true
		prompt.Editable = true

	}

	// Focus on the prompt view.
	if err := menu.gooey.SetCurrentView("prompt"); err != nil {
		return err
	}

	if prompt, err := menu.gooey.View("prompt"); err == nil {

		prompt.Clear()

		if prev != nil {
			fmt.Fprintln(prompt, *prev)
		}

	} else {
		return err
	}

	return nil
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
		cats := ggdb.ReadCategories()

		// Print the name of each category in rows on the menu.
		for _, cat := range cats {
			fmt.Fprintln(catView, *padText(&cat.Name, menu.maxX/6))
		}

		// Check the currently selected row and store its matching category.
		menu.cat = readCat(catView, cats)

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
		exps := ggdb.ReadExpansions(menu.cat)

		// Print name of each expansion to the view.
		for _, exp := range exps {
			fmt.Fprintln(expView, *padText(&exp.Name, menu.maxX*4/6))
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
		phrases := ggdb.ReadPhrases(menu.exp)

		// Print each of the phrases to the view.
		for _, phrase := range phrases {
			fmt.Fprintln(phraseView, *padText(&phrase.Phrase, menu.maxX/6))
		}

	} else {
		return err
	}

	return nil
}

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

			// var helpText string
			help.Clear()
			helpText := "up: ↑ | down: ↓ | left: ← | right: → | "

			switch curView.Name() {
			case "categories":
				helpText += "new: ctrl+n | edit: ctrl+e | select: ⏎"
			case "expansions":
				helpText += "new: ctrl+n | edit: ctrl+e | select: ⏎"
			case "phrases":
				helpText += "new: ctrl+n | edit: ctrl+e"
			case "text":
				helpText += "exit: ctrl+x | save: ctrl+s | reload: ctrl+r"
			case "prompt":
				helpText = "exit: ctrl+x | save: ctrl+s"
			}
			helpText = *centerText(&helpText, menu.maxX)
			fmt.Fprintln(help, helpText)
		}

	} else {
		return err
	}

	return nil
}

//
func drawText() error {

	// Find minY, which will be the bottom of the expansions view.
	_, _, _, minY, err := menu.gooey.ViewPosition("help")
	if err != nil {
		log.Fatal(err)
	}

	// Create the expansion text view, if it doesn't already exist.
	if textView, err := menu.gooey.SetView("text", 0, minY, menu.maxX-1, menu.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//
		textView.Editable = true
		textView.Wrap = true

	}

	return nil
}

func upText(gooey *gocui.Gui) error {

	if expView, err := gooey.View("expansions"); err == nil {
		exps := ggdb.ReadAllExpansions()
		menu.exp = readExp(expView, exps)
	} else {
		return err
	}

	if textView, err := gooey.View("text"); err == nil {
		textView.Clear()
		fmt.Fprintln(textView, menu.exp.Text)
	} else {
		return err
	}

	return nil
}

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
		if err := gooey.SetCurrentView("categories"); err != nil {
			return err
		}
		if err := upText(menu.gooey); err != nil {
			return err
		}
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

	// Always have an exit strategy.
	if err := gooey.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	//
	if err := setKeyBinds(gooey); err != nil {
		log.Panicln(err)
	}

	// Start the main event loop.
	if err := gooey.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
