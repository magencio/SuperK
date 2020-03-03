package editors

import (
	"superk/cmd/utils"

	"github.com/jroimartin/gocui"
)

// NewCustomEditor creates a new simple editor
func NewCustomEditor(clipboard *utils.Clipboard) *gocui.Editor {
	var editor gocui.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case ch != 0 && mod == 0:
			v.EditWrite(ch)
		case key == gocui.KeyCtrlW:
			// KeyCtrlV cannot be intercepted. It will paste contents from system clipboard.
			// I wanted to use KeyCtrlP, but VS Code intercepts it.
			// Using another key combination instead.
			for _, ch := range clipboard.Content {
				v.EditWrite(ch)
			}
		case key == gocui.KeyCtrlD:
			v.Clear()
		case key == gocui.KeySpace:
			v.EditWrite(' ')
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
			v.EditDelete(true)
		case key == gocui.KeyDelete:
			v.EditDelete(false)
		case key == gocui.KeyInsert:
			v.Overwrite = !v.Overwrite
		case key == gocui.KeyArrowLeft:
			v.MoveCursor(-1, 0, false)
		case key == gocui.KeyArrowRight:
			v.MoveCursor(1, 0, false)
		}
	})
	return &editor
}
