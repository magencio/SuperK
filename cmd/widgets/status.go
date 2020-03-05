package widgets

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	// StatusWidgetName is the name of this widget
	StatusWidgetName string = "status"
)

// Check interface
var _ IWidget = &StatusWidget{}

// StatusWidget represents the status bar of the app
type StatusWidget struct {
	Widget
	status string
}

// NewStatusWidget creates a new StatusWidget
func NewStatusWidget() *StatusWidget {
	// Get PID so we can show it on screen and attach the debugger to the process
	return &StatusWidget{Widget: Widget{Name: StatusWidgetName}}
}

// GetName returns the name of the widget
func (widget *StatusWidget) GetName() string { return widget.Name }

// Layout shows the contents of the widget on screen
func (widget *StatusWidget) Layout(g *gocui.Gui, x, y int, w, h int) (*gocui.View, error) {
	widget.X, widget.Y, widget.W, widget.H = x, y, w, h

	v, err := g.SetView(widget.Name, x, y, x+w-1, y+h-1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	v.Frame = false
	v.Clear()
	fmt.Fprint(v, widget.status)
	return v, nil
}

// Refresh updates the contents of the widget on screen
func (widget *StatusWidget) Refresh(g *gocui.Gui) (*gocui.View, error) {
	return widget.Layout(g, widget.X, widget.Y, widget.W, widget.H)
}

// SetAsCurrentView sets the widget as the current view
func (widget *StatusWidget) SetAsCurrentView(g *gocui.Gui) error {
	if _, err := g.SetCurrentView(widget.Name); err != nil {
		return err
	}
	return nil
}

// SetKeyBindings sets keybindings for the widget
func (widget *StatusWidget) SetKeyBindings(g *gocui.Gui) error {
	return nil
}

// SetStatus allows us to set the status shown by the status bar
func (widget *StatusWidget) SetStatus(g *gocui.Gui, status string) error {
	widget.status = status
	if _, err := widget.Refresh(g); err != nil {
		return err
	}
	return nil
}
