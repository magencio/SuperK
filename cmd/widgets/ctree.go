package widgets

import (
	"fmt"
	"superk/cmd/commands"
	"superk/cmd/utils"

	"github.com/jroimartin/gocui"
)

// Check interface
var _ IWidget = &CTreeWidget{}

// CTreeWidget represents a tree of kubectl commands
type CTreeWidget struct {
	Widget
	commands  *commands.CTree
	clipboard *utils.Clipboard
	output    *OutputWidget
	status    *StatusBarWidget
}

// NewCTreeWidget creates a new CTreeWidget
func NewCTreeWidget(
	name string,
	commands *commands.CTree,
	clipboard *utils.Clipboard,
	output *OutputWidget,
	status *StatusBarWidget) *CTreeWidget {
	return &CTreeWidget{
		Widget:    Widget{Name: name, Title: "Commands"},
		commands:  commands,
		clipboard: clipboard,
		output:    output,
		status:    status,
	}
}

// AddCommand adds a new command to the tree
func (widget *CTreeWidget) AddCommand(g *gocui.Gui, command string) error {
	// Update tree
	if err := widget.commands.MergeCommand(command); err != nil {
		return err
	}
	v, err := widget.Refresh(g)
	if err != nil {
		return err
	}

	// Scroll content and set cursor to new command
	_, height := v.Size()
	position := *widget.commands.GetPosition(command)
	originY := utils.Max(0, position-height)
	cursorY := position - originY - 1

	if err := v.SetOrigin(0, originY); err != nil {
		return err
	}
	if err := v.SetCursor(0, cursorY); err != nil {
		return err
	}

	// Set focus to this widget
	if err := widget.SetAsCurrentView(g); err != nil {
		return err
	}
	return nil
}

// GetName returns the name of the widget
func (widget *CTreeWidget) GetName() string { return widget.Name }

// Layout shows the contents of the widget on screen
func (widget *CTreeWidget) Layout(g *gocui.Gui, x, y int, w, h int) (*gocui.View, error) {
	widget.X, widget.Y, widget.W, widget.H = x, y, w, h

	v, err := g.SetView(widget.Name, x, y, x+w-1, y+h-1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	v.Title = widget.Title
	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	v.Clear()
	for _, item := range widget.commands.ToStrings(2) {
		fmt.Fprintln(v, item)
	}

	return v, nil
}

// Refresh updates the contents of the widget on screen
func (widget *CTreeWidget) Refresh(g *gocui.Gui) (*gocui.View, error) {
	return widget.Layout(g, widget.X, widget.Y, widget.W, widget.H)
}

// SetAsCurrentView sets the widget as the current view
func (widget *CTreeWidget) SetAsCurrentView(g *gocui.Gui) error {
	if _, err := g.SetCurrentView(widget.Name); err != nil {
		return err
	}
	if err := widget.status.SetStatus(g, "Commands \x7c \x1b[7mENTER\x1b[0m Execute \x7c \x1b[7m^C\x1b[0m Copy \x7c \x1b[7m^D\x1b[0m Delete \x7c \x1b[7m^X\x1b[0m Exit"); err != nil {
		return err
	}
	return nil
}

// SetKeyBindings sets keybindings for the widget
func (widget *CTreeWidget) SetKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(widget.Name, gocui.KeyArrowUp, gocui.ModNone, widget.moveCursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(widget.Name, gocui.KeyArrowDown, gocui.ModNone, widget.moveCursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(widget.Name, gocui.KeyArrowLeft, gocui.ModNone, widget.moveCursorLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding(widget.Name, gocui.KeyArrowRight, gocui.ModNone, widget.moveCursorRight); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.KeyCtrlC, gocui.ModNone, widget.copyToClipboard); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return widget.SetAsCurrentView(g)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.MouseLeft, gocui.ModNone, widget.run); err != nil {
		return err
	}
	if err := g.SetKeybinding(widget.Name, gocui.KeyEnter, gocui.ModNone, widget.run); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.KeyCtrlD, gocui.ModNone, widget.delete); err != nil {
		return err
	}

	return nil
}

func (widget *CTreeWidget) moveCursorUp(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, -1, false)
	return nil
}

func (widget *CTreeWidget) moveCursorDown(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, 1, false)
	return nil
}

func (widget *CTreeWidget) moveCursorLeft(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(-1, 0, false)
	return nil
}

func (widget *CTreeWidget) moveCursorRight(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(1, 0, false)
	return nil
}

func (widget *CTreeWidget) copyToClipboard(g *gocui.Gui, v *gocui.View) error {
	position := getCommandPosition(v)
	if cmd := widget.commands.GetCmd(position); cmd != nil {
		widget.clipboard.Content = cmd.ToString()
	}
	return nil
}

func (widget *CTreeWidget) run(g *gocui.Gui, v *gocui.View) error {
	position := getCommandPosition(v)
	if cmd := widget.commands.GetCmd(position); cmd != nil {
		output := cmd.Run()

		if err := widget.output.SetCommandOutput(g, cmd, output); err != nil {
			return err
		}
	}
	return nil
}

func (widget *CTreeWidget) delete(g *gocui.Gui, v *gocui.View) error {
	position := getCommandPosition(v)
	if err := widget.commands.RemoveCommand(position); err != nil {
		return nil
	}
	return nil
}

func getCommandPosition(v *gocui.View) int {
	_, yc := v.Cursor()
	_, yo := v.Origin()
	return yc + yo + 1
}
