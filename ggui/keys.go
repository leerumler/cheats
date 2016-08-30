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

	// Refresh the categories, expansions, and phrases views.
	if err := runMenu(gooey); err != nil {
		return err
	}

	// Refresh text view.
	if err := upText(); err != nil {
		return err
	}

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

	// Refresh the categories, expansions, and phrases views.
	if err := runMenu(gooey); err != nil {
		return err
	}

	// Refresh text view.
	if err := upText(); err != nil {
		return err
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

	// Refresh the categories, expansions, and phrases views.
	if err := runMenu(gooey); err != nil {
		return err
	}

	// Refresh text view.
	if err := upText(); err != nil {
		return err
	}

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

	// Refresh the categories, expansions, and phrases views.
	if err := runMenu(gooey); err != nil {
		return err
	}

	// Refresh text view.
	if err := upText(); err != nil {
		return err
	}

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

	// Refresh the categories, expansions, and phrases views.
	if err := runMenu(gooey); err != nil {
		return err
	}

	return nil
}

// focusText changes the focus to the text view.
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

// saveText saves the text inside of the text view to Gengar's database.
func saveText(gooey *gocui.Gui, textView *gocui.View) error {

	// Update menu.exp.Text to the text in the view.
	menu.exp.Text = strings.TrimSpace(textView.ViewBuffer())

	// Update the database with the new value.
	ggdb.UpdateExpansionText(menu.exp)

	// Refresh the text view.
	// if err := upText(); err != nil {
	// 	return err
	// }

	// Change focus to the expansions view.
	if err := focusExp(gooey, nil); err != nil {
		return err
	}

	return nil
}

// closePrompt doesn't actually close anything.  It just pretends
// to close a prompt by putting the expansions view on top of it.
func closePrompt(gooey *gocui.Gui, prompt *gocui.View) error {

	// Set the expansions view on top of any other views.
	if _, err := gooey.SetViewOnTop("expansions"); err != nil {
		return nil
	}

	return nil
}

// closeCatPrompt pretends to close a category prompt and focuses
// back on the categories view.
func closeCatPrompt(gooey *gocui.Gui, prompt *gocui.View) error {

	// Pretend to close the prompt.
	if err := closePrompt(gooey, nil); err != nil {
		return err
	}

	// Focus on the categories view.
	if err := focusCat(gooey, nil); err != nil {
		return err
	}

	return nil
}

// closeExpPrompt pretends to close an expansion prompt and focuses
// back on the expansions view.
func closeExpPrompt(gooey *gocui.Gui, prompt *gocui.View) error {

	// Pretend to close the prompt.
	if err := closePrompt(gooey, nil); err != nil {
		return err
	}

	// Focus on the expansions view.
	if err := focusExp(gooey, nil); err != nil {
		return err
	}

	return nil
}

// closePhrasePrompt pretends to close a phrase prompt and focuses
// back on the phrase view.
func closePhrasePrompt(gooey *gocui.Gui, prompt *gocui.View) error {

	// Pretend to close the prompt.
	if err := closePrompt(gooey, nil); err != nil {
		return err
	}

	// Focus on the phrases view.
	if err := focusPhrase(gooey, nil); err != nil {
		return err
	}

	return nil
}

// newCat reads a category name from a prompt and attempts to
// add that category to Gengar's database.
func newCat(gooey *gocui.Gui, prompt *gocui.View) error {

	// Create a new category and read the name from the prompt.
	var cat ggdb.Category
	cat.Name = strings.TrimSpace(prompt.ViewBuffer())

	// Insert the new category in to the database.
	ggdb.AddCategory(&cat)

	// Pretend to close the prompt.
	if err := closeCatPrompt(gooey, prompt); err != nil {
		return err
	}

	return nil
}

// upCat reads the currently selected category's new name from a prompt
// and attempts to update that name in Gengar's database.
func upCat(gooey *gocui.Gui, prompt *gocui.View) error {

	// Update the currently selected category's name from the
	// prompt and update that in the database.
	menu.cat.Name = strings.TrimSpace(prompt.ViewBuffer())
	ggdb.UpdateCategory(menu.cat)

	// Pretend to close the prompt.
	if err := closeCatPrompt(gooey, prompt); err != nil {
		return err
	}

	return nil
}

// delCat deletes the currently selected category, along with
// all of its expansions and their associated phrases.
func delCat(gooey *gocui.Gui, view *gocui.View) error {

	// Double check the currently selected category.
	menu.cat = readCat()

	// Delete it from the database.
	ggdb.DeleteCategory(menu.cat)

	// Move the cursor up one to account for the deletion.
	if err := selUp(gooey, view); err != nil {
		return err
	}

	return nil
}

// newExp reads an expansion's name from a prompt and attempts to
// add that expansion to Gengar's database.
func newExp(gooey *gocui.Gui, prompt *gocui.View) error {

	// Create a new expansion, read the name from the prompt, and
	// set the category ID to the currently selected category.
	var exp ggdb.Expansion
	exp.Name = strings.TrimSpace(prompt.ViewBuffer())
	exp.CatID = menu.cat.ID

	// Attempt to add that expansion to the database.
	ggdb.AddExpansion(&exp)

	// Pretend to close the prompt.
	if err := closeExpPrompt(gooey, nil); err != nil {
		return err
	}

	return nil
}

// upExp reads the currently selected expansion's new name from a prompt
// and attempts to update that name in Gengar's database.
func upExp(gooey *gocui.Gui, prompt *gocui.View) error {

	// Read the new expansion name from the prompt.
	menu.exp.Name = strings.TrimSpace(prompt.ViewBuffer())

	// Attempt to update that in the database.
	ggdb.UpdateExpansionName(menu.exp)

	// Pretend to close the prompt.
	if err := closeExpPrompt(gooey, nil); err != nil {
		return err
	}

	return nil
}

// delExp deletes the currently selected expansion along with its phrases.
func delExp(gooey *gocui.Gui, view *gocui.View) error {

	// Double check the currently selected expansion and delete it.
	menu.exp = readExp()
	ggdb.DeleteExpansion(menu.exp)

	// Move the cursor up one to account for the deletion.
	if err := selUp(gooey, view); err != nil {
		return err
	}

	return nil
}

// newPhrase reads a phrase's name from a prompt and attempts to
// add that phrase to Gengar's database.
func newPhrase(gooey *gocui.Gui, prompt *gocui.View) error {

	// Create a new phrase, read the name from the prompt, and set
	// the expansion ID to the currently selected expansion.
	var phrase ggdb.Phrase
	phrase.Name = strings.TrimSpace(prompt.ViewBuffer())
	phrase.ExpID = menu.exp.ID

	// Attempt to add that phrase to the database.
	ggdb.AddPhrase(&phrase)

	// Pretend to close the prompt.
	if err := closePhrasePrompt(gooey, nil); err != nil {
		return err
	}

	return nil
}

// upPhrase reads the currently selected phrases's new name from a prompt
// and attempts to update that name in Gengar's database.
func upPhrase(gooey *gocui.Gui, prompt *gocui.View) error {

	// Read the new phrase name from the prompt.
	menu.phrase.Name = strings.TrimSpace(prompt.ViewBuffer())

	// Attempt to update that in the database.
	ggdb.UpdatePhrase(menu.phrase)

	// Pretend to close the prompt.
	if err := closePhrasePrompt(gooey, nil); err != nil {
		return err
	}

	return nil
}

// delPhrase deletes the currently selected phrase from Gengar's database.
func delPhrase(gooey *gocui.Gui, view *gocui.View) error {

	// Double check the currently selected phrase and delete it.
	menu.phrase = readPhrase()
	ggdb.DeletePhrase(menu.phrase)

	// Move the cursor up one to account for the deletion.
	if err := selUp(gooey, view); err != nil {
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

	// If the categories view is focused and ctrl+d is pressed, delete the currently selected category.
	if err := gooey.SetKeybinding("categories", gocui.KeyCtrlD, gocui.ModNone, delCat); err != nil {
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

	// If the expansions view is focused and ctrl+e is pressed, edit the current expansion name.
	if err := gooey.SetKeybinding("expansions", gocui.KeyCtrlE, gocui.ModNone, upExpPrompt); err != nil {
		return err
	}

	// If the expansions view is focused and ctrl+d is pressed, delete the currently selected expansion.
	if err := gooey.SetKeybinding("expansions", gocui.KeyCtrlD, gocui.ModNone, delExp); err != nil {
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

	// If the phrases view is focused and ctrl+n is pressed, add a new phrase.
	if err := gooey.SetKeybinding("phrases", gocui.KeyCtrlN, gocui.ModNone, newPhrasePrompt); err != nil {
		return err
	}

	// If the phrases view is focused and ctrl+e is pressed, edit the current phrase.
	if err := gooey.SetKeybinding("phrases", gocui.KeyCtrlE, gocui.ModNone, upPhrasePrompt); err != nil {
		return err
	}

	// If the phrases view is focused and ctrl+d is pressed, delete the currently selected phrase.
	if err := gooey.SetKeybinding("phrases", gocui.KeyCtrlD, gocui.ModNone, delPhrase); err != nil {
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
