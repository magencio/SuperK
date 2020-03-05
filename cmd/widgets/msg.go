package widgets

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	// MsgWidgetName is the name of this widget
	MsgWidgetName string = "msg"
)

// Check interface
var _ IWidget = &MsgWidget{}

// MsgWidget represents a popup message
type MsgWidget struct {
	Widget
	message string
}

// NewMsgWidget creates a new NewMsgWidget
func NewMsgWidget() *MsgWidget {
	return &MsgWidget{Widget: Widget{Name: MsgWidgetName}}
}

// ShowMsg shows a popup message to user
func (widget *MsgWidget) ShowMsg(g *gocui.Gui, title, message string) error {
	widget.Title, widget.message = title, message

	maxX, maxY := g.Size()
	if _, err := widget.Layout(g, 0, 0, maxX, maxY); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(widget.Name); err != nil {
		return err
	}

	return nil
}

// HideMsg hides the popup message
func (widget *MsgWidget) HideMsg(g *gocui.Gui) error {
	if err := g.DeleteView(widget.Name); err != nil {
		return err
	}
	return nil
}

// GetName returns the name of the widget
func (widget *MsgWidget) GetName() string { return widget.Name }

// Layout shows the contents of the widget on screen
func (widget *MsgWidget) Layout(g *gocui.Gui, x, y int, w, h int) (*gocui.View, error) {
	widget.X, widget.Y, widget.W, widget.H = x, y, w, h

	len := len(widget.message)
	x0, y0 := w/2-len/2-2, h/2-1
	x1, y1 := x0+len+3, y0+2

	v, err := g.SetView(widget.Name, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	v.Title = widget.Title
	v.Clear()
	fmt.Fprintf(v, " %s", widget.message)
	if err := v.SetCursor(len, 0); err != nil {
		return nil, err
	}

	return v, nil
}

// Refresh updates the contents of the widget on screen
func (widget *MsgWidget) Refresh(g *gocui.Gui) (*gocui.View, error) {
	return widget.Layout(g, widget.X, widget.Y, widget.W, widget.H)
}

// SetAsCurrentView sets the widget as the current view
func (widget *MsgWidget) SetAsCurrentView(g *gocui.Gui) error {
	if _, err := g.SetCurrentView(widget.Name); err != nil {
		return err
	}
	return nil
}

// SetKeyBindings sets keybindings for the widget
func (widget *MsgWidget) SetKeyBindings(g *gocui.Gui) error {
	hide := func(g *gocui.Gui, v *gocui.View) error {
		return widget.HideMsg(g)
	}
	if err := g.SetKeybinding(widget.Name, gocui.MouseLeft, gocui.ModNone, hide); err != nil {
		return err
	}
	if err := g.SetKeybinding(widget.Name, gocui.KeyEnter, gocui.ModNone, hide); err != nil {
		return err
	}
	return nil
}
