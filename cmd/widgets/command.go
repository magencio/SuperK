package widgets

import (
	"fmt"
	"strings"
	"superk/cmd/utils"

	"github.com/jroimartin/gocui"
)

// Check interface
var _ IWidget = &CommandWidget{}

// CommandWidget represents a kubectl command to run
type CommandWidget struct {
	Widget
	editor *gocui.Editor
	tree   *CTreeWidget
	status *StatusBarWidget
}

// NewCommandWidget creates a new CommandWidget
func NewCommandWidget(
	name string,
	editor *gocui.Editor,
	tree *CTreeWidget,
	status *StatusBarWidget) *CommandWidget {
	return &CommandWidget{
		Widget: Widget{Name: name, Title: "New Command"},
		editor: editor,
		tree:   tree,
		status: status}
}

// GetName returns the name of the widget
func (widget *CommandWidget) GetName() string { return widget.Name }

// Layout shows the contents of the widget on screen
func (widget *CommandWidget) Layout(g *gocui.Gui, x, y int, w, h int) (*gocui.View, error) {
	widget.X, widget.Y, widget.W, widget.H = x, y, w, h

	v, err := g.SetView(widget.Name, x, y, x+w-1, y+h-1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	v.Title = widget.Title
	v.Editable = true
	v.Editor = *widget.editor

	// If we click with the mouse outside of the command the user wrote, set the cursor at the end
	cx, cy := v.Cursor()
	command, err := v.Line(cy)
	if err != nil {
		cx = 0
	} else {
		cx = utils.Min(cx, len(command))
	}
	if err := v.SetCursor(cx, cy); err != nil {
		return nil, err
	}

	return v, nil
}

// Refresh updates the contents of the widget on screen
func (widget *CommandWidget) Refresh(g *gocui.Gui) (*gocui.View, error) {
	return widget.Layout(g, widget.X, widget.Y, widget.W, widget.H)
}

// SetAsCurrentView sets the widget as the current view
func (widget *CommandWidget) SetAsCurrentView(g *gocui.Gui) error {
	if _, err := g.SetCurrentView(widget.Name); err != nil {
		return err
	}
	if err := widget.status.SetStatus(g, "New Command \x7c \x1b[7mENTER\x1b[0m Add \x7c \x1b[7m^W\x1b[0m Paste \x7c \x1b[7m^D\x1b[0m Delete \x7c \x1b[7m^X\x1b[0m Exit"); err != nil {
		return err
	}
	return nil
}

// SetKeyBindings sets keybindings for the widget
func (widget *CommandWidget) SetKeyBindings(g *gocui.Gui) error {

	if err := g.SetKeybinding(widget.Name, gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return widget.SetAsCurrentView(g)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.KeyEnter, gocui.ModNone, widget.run); err != nil {
		return err
	}

	return nil
}

func (widget *CommandWidget) run(g *gocui.Gui, v *gocui.View) error {
	command, err := v.Line(0)
	if err != nil {
		command = ""
	}
	command = strings.TrimSpace(command)
	if command != "kubectl" && !strings.HasPrefix(command, "kubectl ") {
		command = fmt.Sprintf("kubectl %s", command)
	}

	if err := widget.tree.AddCommand(g, command); err != nil {
		return err
	}

	v.Clear()
	return nil
}
