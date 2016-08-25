package ggui

import (
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/leerumler/gengar/ggdb"
)

// nothibg does nothing.  It should always do nothing, no matter what.
func nothing(gooey *gocui.Gui, view *gocui.View) error {
	return nil
}

// selUp moves the cursor/selection up one line.
func selUp(gooey *gocui.Gui, view *gocui.View) error {

	// Move the cursor up one line.
	if view != nil {
		view.MoveCursor(0, -1, false)
	}

	if err := runMenu(gooey); err != nil {
		return err
	}

	// Refresh text view.
	upText()

	return nil
}

// selDown moves the selected menu item down one line, without moving past the last line.
func selDown(gooey *gocui.Gui, view *gocui.View) error {

	// Move the cursor down one line.
	if view != nil {
		view.MoveCursor(0, 1, false)

		// If the cursor moves to an empty line, move it back. :P
		if readSel(view) == "" {
			view.MoveCursor(0, -1, false)
		}
	}

	if err := runMenu(gooey); err != nil {
		return err
	}

	// Refresh text view.
	upText()

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

	gooey.Cursor = false

	// Focus on the categories view.
	if err := gooey.SetCurrentView("categories"); err != nil {
		return err
	}

	// Reset every view's highlight colors back to blue, then set the
	// set the categories view's highlight color to cyan.
	// So everyone's blue but categories.
	resetHighlights(gooey)
	if catView, err := gooey.View("categories"); err == nil {
		catView.SelBgColor = gocui.ColorCyan
	} else {
		return err
	}

	// Refresh text view.
	upText()

	return nil
}

// focusExp changes the focus to the expansions view.
func focusExp(gooey *gocui.Gui, view *gocui.View) error {

	gooey.Cursor = false

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

	// Refresh text view.
	upText()

	return nil
}

// focusPhrase changes the focus to the phrases view.  It also sets
// the highlight color to Cyan for clarity.
func focusPhrase(gooey *gocui.Gui, view *gocui.View) error {

	gooey.Cursor = false

	// Focus on the phrases view.
	if err := gooey.SetCurrentView("phrases"); err != nil {
		return err
	}

	// Reset every view's highlight colors back to blue, then set the
	// set the phrases view's highlight color to cyan.
	// So everyone's blue but phrases.
	resetHighlights(gooey)
	if phraseView, err := gooey.View("phrases"); err == nil {
		phraseView.SelBgColor = gocui.ColorCyan
	} else {
		return err
	}

	return nil
}

func focusText(gooey *gocui.Gui, view *gocui.View) error {

	// Focus on the text view.
	if err := gooey.SetCurrentView("text"); err != nil {
		return err
	}

	// Set the editor function to textEditor and enable the cursor.
	gooey.Editor = gocui.EditorFunc(multiLineEditor)
	gooey.Cursor = true

	return nil
}

func saveText(gooey *gocui.Gui, textView *gocui.View) error {

	// // Update menu.exp.Text to the text in the view.
	menu.exp.Text = strings.TrimSpace(textView.ViewBuffer())

	// Update the database with the new value.
	ggdb.UpdateExpansionText(menu.exp)

	if err := upText(); err != nil {
		return err
	}

	// Change focus to the expansions view.
	if err := focusExp(gooey, nil); err != nil {
		return err
	}

	return nil
}

func closePrompt(gooey *gocui.Gui, prompt *gocui.View) error {

	if _, err := gooey.SetViewOnTop("expansions"); err != nil {
		return nil
	}
	return nil
}

func closeCatPrompt(gooey *gocui.Gui, prompt *gocui.View) error {
	if err := closePrompt(gooey, nil); err != nil {
		return err
	}
	if err := focusCat(gooey, nil); err != nil {
		return err
	}
	return nil
}

func newCat(gooey *gocui.Gui, promptView *gocui.View) error {

	var cat ggdb.Category
	cat.Name = strings.TrimSpace(promptView.ViewBuffer())
	ggdb.AddCategory(&cat)

	if err := closeCatPrompt(gooey, promptView); err != nil {
		return err
	}

	return nil
}

func upCat(gooey *gocui.Gui, promptView *gocui.View) error {

	menu.cat.Name = strings.TrimSpace(promptView.ViewBuffer())
	ggdb.UpdateCategory(menu.cat)

	if err := closeCatPrompt(gooey, promptView); err != nil {
		return err
	}

	return nil
}

func closeExpPrompt(gooey *gocui.Gui, prompt *gocui.View) error {
	if err := closePrompt(gooey, nil); err != nil {
		return err
	}
	if err := focusExp(gooey, nil); err != nil {
		return err
	}
	return nil
}

func newExp(gooey *gocui.Gui, prompt *gocui.View) error {

	var exp ggdb.Expansion

	exp.Name = strings.TrimSpace(prompt.ViewBuffer())
	exp.CatID = menu.cat.ID
	ggdb.AddExpansion(&exp)

	if err := closeExpPrompt(gooey, nil); err != nil {
		return err
	}

	return nil
}

func upExp(gooey *gocui.Gui, prompt *gocui.View) error {

	menu.exp.Name = strings.TrimSpace(prompt.ViewBuffer())
	ggdb.UpdateExpansionName(menu.exp)

	if err := closeExpPrompt(gooey, nil); err != nil {
		return err
	}

	return nil
}

// setKeyBinds is a necessary evil.
func setKeyBinds(gooey *gocui.Gui) error {

	// If the categories view is focused and ↑ is pressed, move the selected menu item up one.
	if err := gooey.SetKeybinding("categories", gocui.KeyArrowUp, gocui.ModNone, selUp); err != nil {
		return err
	}

	// If the categories view is focused and ↓ is pressed, move the selected menu item down one.
	if err := gooey.SetKeybinding("categories", gocui.KeyArrowDown, gocui.ModNone, selDown); err != nil {
		return err
	}

	// If the categories view is focused and → is pressed, move focus to the expansions view.
	if err := gooey.SetKeybinding("categories", gocui.KeyArrowRight, gocui.ModNone, focusExp); err != nil {
		return err
	}

	// If the categories view is focused and ctrl+n is pressed, add a new category.
	if err := gooey.SetKeybinding("categories", gocui.KeyCtrlN, gocui.ModNone, newCatPrompt); err != nil {
		return err
	}

	// If the categories view is focused and ctrl+e is pressed, edit the current category.
	if err := gooey.SetKeybinding("categories", gocui.KeyCtrlE, gocui.ModNone, upCatPrompt); err != nil {
		return err
	}

	// If the categories view is focused and Enter is pressed, move focus to the expansions view.
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

	// If the expansions view is focused and ← is pressed, move focus to the categories menu.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowLeft, gocui.ModNone, focusCat); err != nil {
		return err
	}

	// If the expansions view is focused and → is pressed, move focus to the phrases view.
	if err := gooey.SetKeybinding("expansions", gocui.KeyArrowRight, gocui.ModNone, focusPhrase); err != nil {
		return err
	}

	// If the expansions view is focused and Enter is pressed, move focus to the phrases view.
	if err := gooey.SetKeybinding("expansions", gocui.KeyEnter, gocui.ModNone, focusText); err != nil {
		return err
	}

	// If the expansions view is focused and ctrl+n is pressed, add a new expansion.
	if err := gooey.SetKeybinding("expansions", gocui.KeyCtrlN, gocui.ModNone, newExpPrompt); err != nil {
		return err
	}

	// If the expansions view is focused and ctrl+e is pressed, edit the current category.
	if err := gooey.SetKeybinding("expansions", gocui.KeyCtrlE, gocui.ModNone, upExpPrompt); err != nil {
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

	// If the phrases view is focused and ← is pressed, move focus to the expansions menu.
	if err := gooey.SetKeybinding("phrases", gocui.KeyArrowLeft, gocui.ModNone, focusExp); err != nil {
		return err
	}

	// If the text view is focused and Escape is pressed, move focus to the expansions menu.
	if err := gooey.SetKeybinding("text", gocui.KeyCtrlX, gocui.ModNone, focusExp); err != nil {
		return err
	}

	// If the text view is focused and Escape is pressed, move focus to the expansions menu.
	if err := gooey.SetKeybinding("text", gocui.KeyCtrlS, gocui.ModNone, saveText); err != nil {
		return err
	}

	return nil
}
