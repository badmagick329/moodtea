package ui

import (
	tea "charm.land/bubbletea/v2"

	"moodtea/internal/data"
)

type Model struct {
	days   []data.Day
	cursor int
	err    error
}

func NewModel(days []data.Day, err error) Model {
	return Model{days: days, err: err}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "right", "l":
			if m.cursor < len(m.days)-1 {
				m.cursor++
			}
		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor+7 < len(m.days) {
				m.cursor += 7
			} else if len(m.days) > 0 {
				m.cursor = len(m.days) - 1
			}
		case "up", "k":
			if m.cursor-7 >= 0 {
				m.cursor -= 7
			} else {
				m.cursor = 0
			}
		}
	}
	return m, nil
}
