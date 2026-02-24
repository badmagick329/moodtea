package ui

import "charm.land/bubbles/v2/key"

type keyMap struct {
	MoveDay   key.Binding
	MoveWeek  key.Binding
	MoveMonth key.Binding
	Scroll    key.Binding
	Today     key.Binding
	GotoMonth key.Binding
	Quit      key.Binding
	Confirm   key.Binding
	Cancel    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.MoveDay, k.MoveWeek, k.MoveMonth, k.Scroll, k.Today, k.GotoMonth, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.MoveDay, k.MoveWeek, k.MoveMonth},
		{k.Scroll, k.Today, k.GotoMonth, k.Quit},
	}
}

func (k keyMap) GotoShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Cancel, k.Quit}
}

func (k keyMap) GotoFullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Confirm, k.Cancel, k.Quit}}
}

type gotoKeyMap struct {
	base keyMap
}

func (k gotoKeyMap) ShortHelp() []key.Binding {
	return k.base.GotoShortHelp()
}

func (k gotoKeyMap) FullHelp() [][]key.Binding {
	return k.base.GotoFullHelp()
}

var keys = keyMap{
	MoveDay: key.NewBinding(
		key.WithKeys("left", "right", "h", "l"),
		key.WithHelp("←/→ h/l", "move day"),
	),
	MoveWeek: key.NewBinding(
		key.WithKeys("up", "down", "k", "j"),
		key.WithHelp("↑/↓ k/j", "move week"),
	),
	MoveMonth: key.NewBinding(
		key.WithKeys("[", "]", "pgup", "pgdown"),
		key.WithHelp("[/] PgUp/PgDn", "move month"),
	),
	Scroll: key.NewBinding(
		key.WithKeys("J", "K"),
		key.WithHelp("J/K", "scroll"),
	),
	Today: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "today"),
	),
	GotoMonth: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "go to month"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "jump"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}
