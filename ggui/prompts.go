package ggui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

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

// newCatPrompt opens a prompt to create a new category.
func newCatPrompt(gooey *gocui.Gui, view *gocui.View) error {

	//
	minX := menu.maxX * 1 / 4
	maxX := menu.maxX * 3 / 4
	midY := menu.maxY / 2

	if _, err := menu.gooey.SetView("promptHead", minX, midY-2, maxX, midY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if promptHead, err := menu.gooey.SetViewOnTop("promptHead"); err == nil {

		promptHead.Clear()
		title := "New Category Name:"
		fmt.Fprintln(promptHead, *centerText(&title, menu.maxX/2))

	} else {
		return err
	}

	if prompt, err := menu.gooey.SetView("newCatPrompt", minX, midY, maxX, midY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := menu.gooey.SetKeybinding("newCatPrompt", gocui.KeyCtrlS, gocui.ModNone, newCat); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("newCatPrompt", gocui.KeyEnter, gocui.ModNone, newCat); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("newCatPrompt", gocui.KeyCtrlX, gocui.ModNone, closeCatPrompt); err != nil {
			return err
		}

		prompt.Editable = true

	}

	if prompt, err := menu.gooey.SetViewOnTop("newCatPrompt"); err == nil {

		menu.gooey.Editor = gocui.EditorFunc(singleLineEditor)
		menu.gooey.Cursor = true
		prompt.SetCursor(0, 0)
		prompt.Clear()

	} else {
		return err
	}

	if err := menu.gooey.SetCurrentView("newCatPrompt"); err != nil {
		return err
	}

	return nil
}

// upCatPrompt opens a prompt to update the currently selected category name.
func upCatPrompt(gooey *gocui.Gui, view *gocui.View) error {

	//
	minX := menu.maxX * 1 / 4
	maxX := menu.maxX * 3 / 4
	midY := menu.maxY / 2

	if _, err := menu.gooey.SetView("promptHead", minX, midY-2, maxX, midY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if promptHead, err := menu.gooey.SetViewOnTop("promptHead"); err == nil {

		promptHead.Clear()
		title := "New Category Name:"
		fmt.Fprintln(promptHead, *centerText(&title, menu.maxX/2))

	} else {
		return err
	}

	if prompt, err := menu.gooey.SetView("upCatPrompt", minX, midY, maxX, midY+2); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}

		if err := menu.gooey.SetKeybinding("upCatPrompt", gocui.KeyCtrlS, gocui.ModNone, upCat); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("upCatPrompt", gocui.KeyEnter, gocui.ModNone, upCat); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("upCatPrompt", gocui.KeyCtrlX, gocui.ModNone, closeCatPrompt); err != nil {
			return err
		}

		prompt.Editable = true

	}

	if prompt, err := menu.gooey.SetViewOnTop("upCatPrompt"); err == nil {

		menu.gooey.Editor = gocui.EditorFunc(singleLineEditor)
		menu.gooey.Cursor = true
		prompt.Clear()
		fmt.Fprintln(prompt, menu.cat.Name)
		prompt.SetCursor(len(menu.cat.Name), 0)

	} else {
		return err
	}

	if err := menu.gooey.SetCurrentView("upCatPrompt"); err != nil {
		return err
	}

	return nil
}

// newExpPrompt opens a prompt to add a new expansion.
func newExpPrompt(gooey *gocui.Gui, view *gocui.View) error {

	//
	minX := menu.maxX * 1 / 4
	maxX := menu.maxX * 3 / 4
	midY := menu.maxY / 2

	if _, err := menu.gooey.SetView("promptHead", minX, midY-2, maxX, midY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if promptHead, err := menu.gooey.SetViewOnTop("promptHead"); err == nil {

		promptHead.Clear()
		title := "New Expansion Name:"
		fmt.Fprintln(promptHead, *centerText(&title, menu.maxX/2))

	} else {
		return err
	}

	if prompt, err := menu.gooey.SetView("newExpPrompt", minX, midY, maxX, midY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := menu.gooey.SetKeybinding("newExpPrompt", gocui.KeyCtrlS, gocui.ModNone, newExp); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("newExpPrompt", gocui.KeyEnter, gocui.ModNone, newExp); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("newExpPrompt", gocui.KeyCtrlX, gocui.ModNone, closeExpPrompt); err != nil {
			return err
		}

		prompt.Editable = true

	}

	if prompt, err := menu.gooey.SetViewOnTop("newExpPrompt"); err == nil {

		menu.gooey.Editor = gocui.EditorFunc(singleLineEditor)
		menu.gooey.Cursor = true
		prompt.SetCursor(0, 0)
		prompt.Clear()

	} else {
		return err
	}

	if err := menu.gooey.SetCurrentView("newExpPrompt"); err != nil {
		return err
	}

	return nil
}

// upExpPrompt opens a prompt to update the currently selected expansion name.
func upExpPrompt(gooey *gocui.Gui, view *gocui.View) error {

	//
	minX := menu.maxX * 1 / 4
	maxX := menu.maxX * 3 / 4
	midY := menu.maxY / 2

	if _, err := menu.gooey.SetView("promptHead", minX, midY-2, maxX, midY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if promptHead, err := menu.gooey.SetViewOnTop("promptHead"); err == nil {

		promptHead.Clear()
		title := "New Expansion Name:"
		fmt.Fprintln(promptHead, *centerText(&title, menu.maxX/2))

	} else {
		return err
	}

	if prompt, err := menu.gooey.SetView("upExpPrompt", minX, midY, maxX, midY+2); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}

		if err := menu.gooey.SetKeybinding("upExpPrompt", gocui.KeyCtrlS, gocui.ModNone, upExp); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("upExpPrompt", gocui.KeyEnter, gocui.ModNone, upExp); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("upExpPrompt", gocui.KeyCtrlX, gocui.ModNone, closeExpPrompt); err != nil {
			return err
		}

		prompt.Editable = true

	}

	if prompt, err := menu.gooey.SetViewOnTop("upExpPrompt"); err == nil {

		menu.gooey.Editor = gocui.EditorFunc(singleLineEditor)
		menu.gooey.Cursor = true
		prompt.Clear()
		fmt.Fprintln(prompt, menu.exp.Name)
		prompt.SetCursor(len(menu.exp.Name), 0)

	} else {
		return err
	}

	if err := menu.gooey.SetCurrentView("upExpPrompt"); err != nil {
		return err
	}

	return nil
}

// newPhrasePrompt opens a prompt to add a new phrase.
func newPhrasePrompt(gooey *gocui.Gui, view *gocui.View) error {

	//
	minX := menu.maxX * 1 / 4
	maxX := menu.maxX * 3 / 4
	midY := menu.maxY / 2

	if _, err := menu.gooey.SetView("promptHead", minX, midY-2, maxX, midY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if promptHead, err := menu.gooey.SetViewOnTop("promptHead"); err == nil {

		promptHead.Clear()
		title := "New Phrase Name:"
		fmt.Fprintln(promptHead, *centerText(&title, menu.maxX/2))

	} else {
		return err
	}

	if prompt, err := menu.gooey.SetView("newPhrasePrompt", minX, midY, maxX, midY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := menu.gooey.SetKeybinding("newPhrasePrompt", gocui.KeyCtrlS, gocui.ModNone, newPhrase); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("newPhrasePrompt", gocui.KeyEnter, gocui.ModNone, newPhrase); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("newPhrasePrompt", gocui.KeyCtrlX, gocui.ModNone, closePhrasePrompt); err != nil {
			return err
		}

		prompt.Editable = true

	}

	if prompt, err := menu.gooey.SetViewOnTop("newPhrasePrompt"); err == nil {

		menu.gooey.Editor = gocui.EditorFunc(singleLineEditor)
		menu.gooey.Cursor = true
		prompt.SetCursor(0, 0)
		prompt.Clear()

	} else {
		return err
	}

	if err := menu.gooey.SetCurrentView("newPhrasePrompt"); err != nil {
		return err
	}

	return nil
}

// upPhrasePrompt opens a prompt to update the currently selected phrase.
func upPhrasePrompt(gooey *gocui.Gui, view *gocui.View) error {

	//
	minX := menu.maxX * 1 / 4
	maxX := menu.maxX * 3 / 4
	midY := menu.maxY / 2

	if _, err := menu.gooey.SetView("promptHead", minX, midY-2, maxX, midY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if promptHead, err := menu.gooey.SetViewOnTop("promptHead"); err == nil {

		promptHead.Clear()
		title := "New Expansion Name:"
		fmt.Fprintln(promptHead, *centerText(&title, menu.maxX/2))

	} else {
		return err
	}

	if prompt, err := menu.gooey.SetView("upPhrasePrompt", minX, midY, maxX, midY+2); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}

		if err := menu.gooey.SetKeybinding("upPhrasePrompt", gocui.KeyCtrlS, gocui.ModNone, upPhrase); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("upPhrasePrompt", gocui.KeyEnter, gocui.ModNone, upPhrase); err != nil {
			return err
		}

		if err := menu.gooey.SetKeybinding("upPhrasePrompt", gocui.KeyCtrlX, gocui.ModNone, closePhrasePrompt); err != nil {
			return err
		}

		prompt.Editable = true

	}

	if prompt, err := menu.gooey.SetViewOnTop("upPhrasePrompt"); err == nil {

		menu.gooey.Editor = gocui.EditorFunc(singleLineEditor)
		menu.gooey.Cursor = true
		prompt.Clear()
		fmt.Fprintln(prompt, menu.phrase.Name)
		prompt.SetCursor(len(menu.phrase.Name), 0)

	} else {
		return err
	}

	if err := menu.gooey.SetCurrentView("upPhrasePrompt"); err != nil {
		return err
	}

	return nil
}
