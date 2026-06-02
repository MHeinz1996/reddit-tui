package posts

import "github.com/charmbracelet/bubbles/key"

type postsKeyMap struct {
	Home   key.Binding
	Search key.Binding
	Sort   key.Binding
	SortHot key.Binding
	SortNew key.Binding
	SortRising key.Binding
	SortTop key.Binding
	SortControversial key.Binding
	Back   key.Binding
	Load   key.Binding
}

var postsKeys = postsKeyMap{
	Home: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "home")),
	Search: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "subreddit search")),
	Sort: key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5"),
		key.WithHelp("1-5", "sort")),
	SortHot: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "hot")),
	SortNew: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "new")),
	SortRising: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "rising")),
	SortTop: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "top")),
	SortControversial: key.NewBinding(
		key.WithKeys("5"),
		key.WithHelp("5", "controversial")),
	Back: key.NewBinding(
		key.WithKeys("bs"),
		key.WithHelp("bs", "back")),
	Load: key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "load more posts")),
}

func (k postsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Home, k.Search, k.Sort, k.Load}
}

func (k postsKeyMap) FullHelp() []key.Binding {
	return []key.Binding{k.Home, k.Search, k.SortHot, k.SortNew, k.SortRising, k.SortTop, k.SortControversial, k.Back, k.Load}
}
