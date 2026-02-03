package ui

import (
	tea "charm.land/bubbletea/v2"

	"moodtea/internal/data"
)

type Model struct {
	months []data.Month
	state  State
	err    error
}

const (
	keyQuit       = "ctrl+c"
	keyQuitAlt    = "q"
	keyQuitAlt2   = "esc"
	keyNextMonth  = "]"
	keyPrevMonth  = "["
	keyNextMonth2 = "pgdown"
	keyPrevMonth2 = "pgup"
	keyRight      = "right"
	keyRightAlt   = "l"
	keyLeft       = "left"
	keyLeftAlt    = "h"
	keyDown       = "down"
	keyDownAlt    = "j"
	keyUp         = "up"
	keyUpAlt      = "k"
)

func NewModel(months []data.Month, startKey string, err error) Model {
	m := Model{months: months, err: err}
	if len(months) == 0 {
		return m
	}
	m.state.MonthIndex = 0
	for i, month := range months {
		if month.Key == startKey {
			m.state.MonthIndex = i
			break
		}
	}
	if months[m.state.MonthIndex].Key != startKey {
		m.state.MonthIndex = lastMonthIndexBefore(months, startKey)
	}
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyQuit, keyQuitAlt, keyQuitAlt2:
			return m, tea.Quit
		case keyNextMonth, keyNextMonth2:
			m.moveMonth(1)
		case keyPrevMonth, keyPrevMonth2:
			m.moveMonth(-1)
		case keyRight, keyRightAlt:
			if m.state.Cursor < len(m.currentDays())-1 {
				m.state.Cursor++
			}
		case keyLeft, keyLeftAlt:
			if m.state.Cursor > 0 {
				m.state.Cursor--
			}
		case keyDown, keyDownAlt:
			if m.state.Cursor+7 < len(m.currentDays()) {
				m.state.Cursor += 7
			} else if len(m.currentDays()) > 0 {
				m.state.Cursor = len(m.currentDays()) - 1
			}
		case keyUp, keyUpAlt:
			if m.state.Cursor-7 >= 0 {
				m.state.Cursor -= 7
			} else {
				m.state.Cursor = 0
			}
		}
	}
	return m, nil
}

func (m *Model) currentDays() []data.Day {
	if len(m.months) == 0 {
		return nil
	}
	return m.months[m.state.MonthIndex].Days
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
	if len(oldDays) > 0 && m.state.Cursor >= 0 && m.state.Cursor < len(oldDays) {
		selectedDay = oldDays[m.state.Cursor].Date.Day()
	}

	next := clamp(m.state.MonthIndex+delta, 0, len(m.months)-1)
	if next == m.state.MonthIndex {
		return
	}
	m.state.MonthIndex = next

	newDays := m.currentDays()
	if len(newDays) == 0 {
		m.state.Cursor = 0
		return
	}

	m.state.Cursor = findDayIndex(newDays, selectedDay)
}

func findDayIndex(days []data.Day, day int) int {
	for i, d := range days {
		if d.Date.Day() == day {
			return i
		}
	}
	return len(days) - 1
}

func clamp(v int, min int, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func lastMonthIndexBefore(months []data.Month, key string) int {
	if len(months) == 0 {
		return 0
	}
	idx := 0
	for i, month := range months {
		if month.Key <= key {
			idx = i
		}
	}
	return idx
}
