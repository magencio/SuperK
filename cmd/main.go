package main

import (
	"log"
	"superk/cmd/commands"
	"superk/cmd/widgets"

	"github.com/jroimartin/gocui"
)

const (
	backupName string = "superk_backup"
)

func main() {

	commands, err := loadCommandsFromBackup()
	if err != nil {
		log.Panicln(err)
	}
	defer backupCommands(commands)

	g, err := createNewGui()
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	widgets := createWidgets(commands)

	setGuiManager(g, widgets.MainScreen())

	if err := setGlobalKeybindings(g, widgets); err != nil {
		log.Panicln(err)
	}

	if err := setWidgetKeybindings(g, widgets); err != nil {
		log.Panicln(err)
	}

	if err := mainLoop(g); err != nil {
		log.Panicln(err)
	}
}

func loadCommandsFromBackup() (*commands.CTree, error) {
	return commands.NewBackup(backupName).Commands()
}

func backupCommands(commandTree *commands.CTree) {
	if err := commands.NewBackup(backupName).SetCommands(commandTree); err != nil {
		log.Panicln(err)
	}
}

func createNewGui() (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	// Enable mouse support
	g.Cursor = true
	g.Mouse = true

	return g, nil
}

func createWidgets(commands *commands.CTree) *widgets.Widgets {
	return widgets.NewWidgets(commands)
}

func setGuiManager(g *gocui.Gui, widget widgets.IWidget) {
	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		if _, err := widget.Layout(g, 0, 0, maxX, maxY); err != nil {
			return err
		}
		return nil
	})
}

func setGlobalKeybindings(g *gocui.Gui, allWidgets *widgets.Widgets) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlX, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit // Exit the app
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		mainScreen := allWidgets.MainScreen()
		return mainScreen.OnTab(g)
	}); err != nil {
		return err
	}

	return nil
}

func setWidgetKeybindings(g *gocui.Gui, allWidgets *widgets.Widgets) error {
	for _, widget := range allWidgets.Widgets() {
		if err := widget.SetKeyBindings(g); err != nil {
			return err
		}
	}
	return nil
}

func mainLoop(g *gocui.Gui) error {
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}
