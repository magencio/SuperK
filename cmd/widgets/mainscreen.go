package widgets

import (
	"fmt"
	"superk/cmd/utils"

	"github.com/jroimartin/gocui"
)

const (
	// MainScreenWidgetName is the name of this widget
	MainScreenWidgetName string = "mainScreen"
)

// Check interface
var _ IWidget = &MainScreenWidget{}

// MainScreenWidget represents the main screen of the app
type MainScreenWidget struct {
	Widget
	widgets  *Widgets
	tabOrder []IWidget
}

// NewMainScreenWidget creates a new MainScreenWidget
func NewMainScreenWidget(
	widgets *Widgets) *MainScreenWidget {
	return &MainScreenWidget{
		Widget:   Widget{Name: MainScreenWidgetName},
		widgets:  widgets,
		tabOrder: []IWidget{widgets.Command(), widgets.Tree(), widgets.Output()},
	}
}

// GetName returns the name of the widget
func (widget *MainScreenWidget) GetName() string { return widget.Name }

// Layout shows the contents of the widget on screen
func (widget *MainScreenWidget) Layout(g *gocui.Gui, x, y int, w, h int) (*gocui.View, error) {
	widget.X, widget.Y, widget.W, widget.H = x, y, w, h

	v, err := g.SetView(widget.Name, x, y, x+w-1, y+h-1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	minWidth, minHeight := 60, 10
	if width, height := v.Size(); width < minWidth || height < minHeight {
		if err := widget.deleteAllOtherViews(g); err != nil {
			return nil, err
		}
		v.Frame = true
		v.Clear()
		v.Wrap = true
		fmt.Fprintln(v, fmt.Sprintf("Window size not supported. Minimum: width %d x height %d", minWidth, minHeight))
		return v, nil
	}

	v.Title = widget.Title
	v.Frame = false

	if _, err := widget.widgets.Command().Layout(g, x, y, x+w-1, 3); err != nil {
		return nil, err
	}

	if _, err := widget.widgets.Tree().Layout(g, x, y+3, x+w/3-1, h-4); err != nil {
		return nil, err
	}

	if _, err := widget.widgets.Output().Layout(g, x+w/3-1, y+3, x+w-w/3, h-4); err != nil {
		return nil, err
	}

	if _, err := widget.widgets.Status().Layout(g, x, y+h-2, x+w-1, 3); err != nil {
		return nil, err
	}

	if g.CurrentView() == nil {
		if err := widget.widgets.Command().SetAsCurrentView(g); err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (widget *MainScreenWidget) deleteAllOtherViews(g *gocui.Gui) error {
	var viewsToDelete []string
	for _, v := range g.Views() {
		if v.Name() != widget.Name {
			viewsToDelete = append(viewsToDelete, v.Name())
		}
	}
	for _, viewName := range viewsToDelete {
		if err := g.DeleteView(viewName); err != nil {
			return err
		}
	}
	return nil
}

// Refresh updates the contents of the widget on screen
func (widget *MainScreenWidget) Refresh(g *gocui.Gui) (*gocui.View, error) {
	return widget.Layout(g, widget.X, widget.Y, widget.W, widget.H)
}

// SetAsCurrentView sets the widget as the current view
func (widget *MainScreenWidget) SetAsCurrentView(g *gocui.Gui) error {
	if _, err := g.SetCurrentView(widget.Name); err != nil {
		return err
	}
	return nil
}

// SetKeyBindings sets keybindings for the widget
func (widget *MainScreenWidget) SetKeyBindings(g *gocui.Gui) error {
	return nil
}

// OnTab handles the event of user pressing Tab key
func (widget *MainScreenWidget) OnTab(g *gocui.Gui) error {
	name := g.CurrentView().Name()
	for index, current := range widget.tabOrder {
		if current.GetName() == name {
			nextIndex := utils.Mod(index+1, len(widget.tabOrder))
			return widget.tabOrder[nextIndex].SetAsCurrentView(g)
		}
	}
	return nil
}
