package widgets

import (
	"fmt"
	"superk/cmd/commands"
	"superk/cmd/utils"

	"github.com/jroimartin/gocui"
)

const (
	// TreeWidgetName is the name of this widget
	TreeWidgetName  string = "tree"
	treeWidgetTitle string = "Commands"
	treeWidgetHelp  string = "Commands \x7c \x1b[7mENTER\x1b[0m Update \x7c \x1b[7m^R\x1b[0m Reuse \x7c \x1b[7m^C\x1b[0m Copy \x7c \x1b[7m^D\x1b[0m Delete \x7c \x1b[7m^X\x1b[0m Exit"
)

// Check interface
var _ IWidget = &TreeWidget{}

// TreeWidget represents a tree of kubectl commands
type TreeWidget struct {
	Widget
	commands  *commands.CTree
	clipboard *utils.Clipboard
	widgets   *Widgets
}

// NewTreeWidget creates a new TreeWidget
func NewTreeWidget(
	commands *commands.CTree,
	clipboard *utils.Clipboard,
	widgets *Widgets) *TreeWidget {
	return &TreeWidget{
		Widget:    Widget{Name: TreeWidgetName, Title: treeWidgetTitle},
		commands:  commands,
		clipboard: clipboard,
		widgets:   widgets,
	}
}

// AddCommand adds a new command to the tree
func (widget *TreeWidget) AddCommand(g *gocui.Gui, command string) error {
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

	// Run command
	return widget.run(g, v, false)
}

// GetName returns the name of the widget
func (widget *TreeWidget) GetName() string { return widget.Name }

// Layout shows the contents of the widget on screen
func (widget *TreeWidget) Layout(g *gocui.Gui, x, y int, w, h int) (*gocui.View, error) {
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
func (widget *TreeWidget) Refresh(g *gocui.Gui) (*gocui.View, error) {
	return widget.Layout(g, widget.X, widget.Y, widget.W, widget.H)
}

// SetAsCurrentView sets the widget as the current view
func (widget *TreeWidget) SetAsCurrentView(g *gocui.Gui) error {
	v, err := g.SetCurrentView(widget.Name)
	if err != nil {
		return err
	}

	if err := widget.run(g, v, true); err != nil {
		return err
	}

	if err := widget.widgets.Status().SetStatus(g, treeWidgetHelp); err != nil {
		return err
	}
	return nil
}

// SetKeyBindings sets keybindings for the widget
func (widget *TreeWidget) SetKeyBindings(g *gocui.Gui) error {
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

	if err := g.SetKeybinding(widget.Name, gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return widget.SetAsCurrentView(g)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return widget.run(g, v, true)
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding(widget.Name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return widget.run(g, v, false)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.KeyCtrlR, gocui.ModNone, widget.reuse); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.KeyCtrlC, gocui.ModNone, widget.copyToClipboard); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.KeyCtrlD, gocui.ModNone, widget.delete); err != nil {
		return err
	}

	return nil
}

func (widget *TreeWidget) moveCursorUp(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, -1, false)
	return widget.run(g, v, true)
}

func (widget *TreeWidget) moveCursorDown(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, 1, false)
	return widget.run(g, v, true)
}

func (widget *TreeWidget) moveCursorLeft(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(-1, 0, false)
	return widget.run(g, v, true)
}

func (widget *TreeWidget) moveCursorRight(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(1, 0, false)
	return widget.run(g, v, true)
}

func (widget *TreeWidget) reuse(g *gocui.Gui, v *gocui.View) error {
	position := getCommandPosition(v)
	if cmd := widget.commands.GetCmd(position); cmd != nil {
		if err := widget.widgets.Command().SetContent(g, cmd.ToString()); err != nil {
			return err
		}
	}
	return nil
}

func (widget *TreeWidget) copyToClipboard(g *gocui.Gui, v *gocui.View) error {
	position := getCommandPosition(v)
	if cmd := widget.commands.GetCmd(position); cmd != nil {
		widget.clipboard.Content = cmd.ToString()
	}
	return nil
}

func (widget *TreeWidget) run(g *gocui.Gui, v *gocui.View, cacheFirst bool) error {
	position := getCommandPosition(v)
	if cmd := widget.commands.GetCmd(position); cmd != nil {
		_ = cmd.Run(cacheFirst)

		if err := widget.widgets.Output().SetCommandOutput(g, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (widget *TreeWidget) delete(g *gocui.Gui, v *gocui.View) error {
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
