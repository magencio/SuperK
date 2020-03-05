package widgets

import (
	"superk/cmd/commands"
	"superk/cmd/editors"
	"superk/cmd/utils"
)

// Widgets represents all the widgets in the app
type Widgets struct {
	widgets map[string]IWidget
}

// NewWidgets creates a new Widgets
func NewWidgets(commands *commands.CTree) *Widgets {
	clipboard := utils.NewClipboard()
	editor := editors.NewCustomEditor(clipboard)

	all := Widgets{widgets: map[string]IWidget{}}
	all.widgets[MsgWidgetName] = NewMsgWidget()
	all.widgets[StatusWidgetName] = NewStatusWidget()
	all.widgets[OutputWidgetName] = NewOutputWidget(clipboard, &all)
	all.widgets[TreeWidgetName] = NewTreeWidget(commands, clipboard, &all)
	all.widgets[CommandWidgetName] = NewCommandWidget(editor, &all)
	all.widgets[MainScreenWidgetName] = NewMainScreenWidget(&all)

	return &all
}

// Msg returns the message widget
func (all *Widgets) Msg() *MsgWidget { return all.widgets[MsgWidgetName].(*MsgWidget) }

// Status returns the status widget
func (all *Widgets) Status() *StatusWidget { return all.widgets[StatusWidgetName].(*StatusWidget) }

// Output returns the output widget
func (all *Widgets) Output() *OutputWidget { return all.widgets[OutputWidgetName].(*OutputWidget) }

// Tree returns the command tree widget
func (all *Widgets) Tree() *TreeWidget { return all.widgets[TreeWidgetName].(*TreeWidget) }

// Command returns the command widget
func (all *Widgets) Command() *CommandWidget { return all.widgets[CommandWidgetName].(*CommandWidget) }

// MainScreen returns the main screen widget
func (all *Widgets) MainScreen() *MainScreenWidget {
	return all.widgets[MainScreenWidgetName].(*MainScreenWidget)
}

// Widgets returns all widgets
func (all *Widgets) Widgets() []IWidget {
	values := make([]IWidget, 0, len(all.widgets))
	for _, v := range all.widgets {
		values = append(values, v)
	}
	return values
}
