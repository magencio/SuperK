package widgets

import (
	"fmt"
	"strings"
	"superk/cmd/commands"
	"superk/cmd/utils"
	"time"
	"unicode"

	"github.com/jroimartin/gocui"
)

const (
	// OutputWidgetName is the name of this widget
	OutputWidgetName  string = "output"
	outputWidgetTitle string = "Output"
	outputWidgetHelp  string = "Output \x7c \x1b[7m^C\x1b[0m Copy word \x7c \x1b[7m^L\x1b[0m Copy line \x7c \x1b[7m^X\x1b[0m Exit"
)

// Check interface
var _ IWidget = &OutputWidget{}

// OutputWidget represents the output of a kubectl command
type OutputWidget struct {
	Widget
	output    *string
	clipboard *utils.Clipboard
	widgets   *Widgets
}

// NewOutputWidget creates a new OutputWidget
func NewOutputWidget(
	clipboard *utils.Clipboard,
	widgets *Widgets) *OutputWidget {
	return &OutputWidget{
		Widget:    Widget{Name: OutputWidgetName, Title: outputWidgetTitle},
		clipboard: clipboard,
		widgets:   widgets}
}

// SetCommandOutput sets the command and its output that this widget will show to user
func (widget *OutputWidget) SetCommandOutput(g *gocui.Gui, cmd *commands.Cmd) error {
	// Refresh widget
	widget.Title = fmt.Sprintf("Output [%s] [%s]", cmd.ToString(), cmd.RunTime.Format(time.UnixDate))
	widget.output = cmd.CmdOutput.Output
	v, err := widget.Refresh(g)
	if err != nil {
		return err
	}

	// Scroll view and set cursor at the end of the output
	// TODO: This is not taking wrapped lines into account!
	_, height := v.Size()
	lineCount := len(v.BufferLines())
	originY := utils.Max(0, lineCount-height)
	cursorY := lineCount - originY - 1
	if err := v.SetOrigin(0, originY); err != nil {
		return err
	}
	if err := v.SetCursor(0, cursorY); err != nil {
		return err
	}

	return nil
}

// GetName returns the name of the widget
func (widget *OutputWidget) GetName() string { return widget.Name }

// Layout shows the contents of the widget on screen
func (widget *OutputWidget) Layout(g *gocui.Gui, x, y int, w, h int) (*gocui.View, error) {
	widget.X, widget.Y, widget.W, widget.H = x, y, w, h

	v, err := g.SetView(widget.Name, x, y, x+w-1, y+h-1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	v.Title = widget.Title
	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	v.Wrap = true
	v.Clear()
	if widget.output != nil {
		fmt.Fprintln(v, *widget.output)
	}

	return v, nil
}

// Refresh updates the contents of the widget on screen
func (widget *OutputWidget) Refresh(g *gocui.Gui) (*gocui.View, error) {
	return widget.Layout(g, widget.X, widget.Y, widget.W, widget.H)
}

// SetAsCurrentView sets the widget as the current view
func (widget *OutputWidget) SetAsCurrentView(g *gocui.Gui) error {
	if _, err := g.SetCurrentView(widget.Name); err != nil {
		return err
	}
	if err := widget.widgets.Status().SetStatus(g, outputWidgetHelp); err != nil {
		return err
	}
	return nil
}

// SetKeyBindings sets keybindings for the widget
func (widget *OutputWidget) SetKeyBindings(g *gocui.Gui) error {
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

	if err := g.SetKeybinding(widget.Name, gocui.KeyCtrlC, gocui.ModNone, widget.copyWordToClipboard); err != nil {
		return err
	}
	if err := g.SetKeybinding(widget.Name, gocui.KeyCtrlL, gocui.ModNone, widget.copyLineToClipboard); err != nil {
		return err
	}

	if err := g.SetKeybinding(widget.Name, gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return widget.SetAsCurrentView(g)
	}); err != nil {
		return err
	}

	return nil
}

func (widget *OutputWidget) moveCursorUp(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, -1, false)
	return nil
}

func (widget *OutputWidget) moveCursorDown(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, 1, false)
	return nil
}

func (widget *OutputWidget) moveCursorLeft(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(-1, 0, false)
	return nil
}

func (widget *OutputWidget) moveCursorRight(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(1, 0, false)
	return nil
}

func (widget *OutputWidget) copyWordToClipboard(g *gocui.Gui, v *gocui.View) error {
	xc, yc := v.Cursor()
	if line, err := v.Line(yc); err == nil && xc <= len(line) {
		lastSpaceBeforeWord := -1
		firstSpaceAfterWord := len(line)
		for index, ch := range line {
			if unicode.IsSpace(ch) && index < xc {
				lastSpaceBeforeWord = index
			}
			if unicode.IsSpace(ch) && index == xc {
				firstSpaceAfterWord = -1
				break
			}
			if unicode.IsSpace(ch) && index > xc {
				firstSpaceAfterWord = index
				break
			}
		}
		if firstSpaceAfterWord > lastSpaceBeforeWord {
			word := line[lastSpaceBeforeWord+1 : firstSpaceAfterWord]
			widget.clipboard.Content = word
		} else {
			widget.clipboard.Content = ""
		}
	}
	return nil
}

func (widget *OutputWidget) copyLineToClipboard(g *gocui.Gui, v *gocui.View) error {
	_, yc := v.Cursor()
	if line, err := v.Line(yc); err == nil {
		fields := strings.FieldsFunc(line, func(ch rune) bool {
			return ch == '\n'
		})
		if len(fields) > 0 {
			widget.clipboard.Content = fields[0]
		}
	}
	return nil
}
