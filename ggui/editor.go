package ggui

import "github.com/jroimartin/gocui"

// textEditor is used as the default gocui editor.
func multiLineEditor(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	switch {
	case char != 0 && mod == 0:
		view.EditWrite(char)
	case key == gocui.KeySpace:
		view.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		view.EditDelete(true)
	case key == gocui.KeyDelete:
		view.EditDelete(false)
	case key == gocui.KeyInsert:
		view.Overwrite = !view.Overwrite
	case key == gocui.KeyEnter:
		view.EditNewLine()
	case key == gocui.KeyArrowDown:
		view.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		view.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		view.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		view.MoveCursor(1, 0, false)
	}
}

// textEditor is used as the default gocui editor.
func singleLineEditor(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	switch {
	case char != 0 && mod == 0:
		view.EditWrite(char)
	case key == gocui.KeySpace:
		view.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		view.EditDelete(true)
	case key == gocui.KeyDelete:
		view.EditDelete(false)
	case key == gocui.KeyInsert:
		view.Overwrite = !view.Overwrite
	case key == gocui.KeyArrowDown:
		view.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		view.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		view.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		view.MoveCursor(1, 0, false)
	}
}
