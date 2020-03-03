package widgets

import "github.com/jroimartin/gocui"

// Position represents the coordinates of the top left corner of a widget
type Position struct {
	X, Y int
}

// Size represents the width and height of a widget
type Size struct {
	W, H int
}

// Widget represents the basic definition of a widget
type Widget struct {
	Name  string
	Title string
	Position
	Size
}

// IDrawable represents a widget that can be drawn
type IDrawable interface {
	Layout(g *gocui.Gui, x, y int, w, h int) (*gocui.View, error)
	Refresh(g *gocui.Gui) (*gocui.View, error)
	SetAsCurrentView(g *gocui.Gui) error
}

// IBindable represents a widget with key bindings
type IBindable interface {
	SetKeyBindings(g *gocui.Gui) error
}

// IWidget represents the basic behavior of a widget
type IWidget interface {
	GetName() string
	IDrawable
	IBindable
}
