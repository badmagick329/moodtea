package ui

import (
	tea "charm.land/bubbletea/v2"

	"moodtea/internal/data"
)

type Model struct {
	months     []data.Month
	monthIndex int
	cursor     int
	err        error
}

func NewModel(months []data.Month, startKey string, err error) Model {
	m := Model{months: months, err: err}
	if len(months) == 0 {
		return m
	}
	m.monthIndex = 0
	for i, month := range months {
		if month.Key == startKey {
			m.monthIndex = i
			break
		}
	}
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "]", "pgdown":
			m.moveMonth(1)
		case "[", "pgup":
			m.moveMonth(-1)
		case "right", "l":
			if m.cursor < len(m.currentDays())-1 {
				m.cursor++
			}
		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor+7 < len(m.currentDays()) {
				m.cursor += 7
			} else if len(m.currentDays()) > 0 {
				m.cursor = len(m.currentDays()) - 1
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

func (m *Model) currentDays() []data.Day {
	if len(m.months) == 0 {
		return nil
	}
	return m.months[m.monthIndex].Days
}

func (m *Model) moveMonth(delta int) {
	if len(m.months) == 0 {
		return
	}
	if delta == 0 {
		return
	}

	oldDays := m.currentDays()
	selectedDay := 1
	if len(oldDays) > 0 && m.cursor >= 0 && m.cursor < len(oldDays) {
		selectedDay = oldDays[m.cursor].Date.Day()
	}

	next := m.monthIndex + delta
	if next < 0 {
		next = 0
	} else if next >= len(m.months) {
		next = len(m.months) - 1
	}
	if next == m.monthIndex {
		return
	}
	m.monthIndex = next

	newDays := m.currentDays()
	if len(newDays) == 0 {
		m.cursor = 0
		return
	}

	m.cursor = findDayIndex(newDays, selectedDay)
}

func findDayIndex(days []data.Day, day int) int {
	for i, d := range days {
		if d.Date.Day() == day {
			return i
		}
	}
	return len(days) - 1
}
