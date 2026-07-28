package comments

import "github.com/charmbracelet/bubbles/key"

type viewportKeyMap struct {
	CursorUp          key.Binding
	CursorDown        key.Binding
	GoToStart         key.Binding
	GoToEnd           key.Binding
	OpenPost          key.Binding
	Refresh           key.Binding
	GoHome            key.Binding
	Sort              key.Binding
	SortBest          key.Binding
	SortNew           key.Binding
	SortTop           key.Binding
	SortControversial key.Binding
	SortOld           key.Binding
	CollapseComments  key.Binding
	ShowFullHelp      key.Binding
	CloseFullHelp     key.Binding
	Quit              key.Binding
	ForceQuit         key.Binding
}

var commentsKeys = viewportKeyMap{
	CursorUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "prev comment"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "next comment"),
	),
	GoToStart: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("g/home", "first comment"),
	),
	GoToEnd: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("G/end", "last comment"),
	),
	OpenPost: key.NewBinding(
		key.WithKeys("o", "O"),
		key.WithHelp("o", "open post"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	GoHome: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "go home"),
	),
	Sort: key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5"),
		key.WithHelp("1-5", "sort")),
	SortBest: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "best")),
	SortNew: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "new")),
	SortTop: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "top")),
	SortControversial: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "controversial")),
	SortOld: key.NewBinding(
		key.WithKeys("5"),
		key.WithHelp("5", "old")),
	CollapseComments: key.NewBinding(
		key.WithKeys(" ", "space"),
		key.WithHelp("space", "collapse/expand"),
	),
	ShowFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "more"),
	),
	CloseFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "close help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "quit"),
	),
	ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
}

func (k viewportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.CursorUp, k.CursorDown, k.Sort, k.OpenPost, k.Refresh, k.GoHome, k.ShowFullHelp}
}

func (k viewportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown, k.GoToStart, k.GoToEnd, k.OpenPost},
		{k.SortBest, k.SortNew, k.SortTop, k.SortControversial, k.SortOld},
		{k.GoHome, k.Refresh, k.CollapseComments, k.Quit, k.CloseFullHelp},
	}
}
