package ui

import (
	"regexp"
	"time"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"

	"moodtea/internal/data"
)

type Model struct {
	months    []data.Month
	state     State
	helpModel help.Model
	err       error
}

const (
	keyQuit       = "ctrl+c"
	keyQuitAlt    = "q"
	keyQuitAlt2   = "esc"
	keyQuitAlt3   = "escape"
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
	keyToday      = "t"
	keyGoto       = "g"
	keyEnter      = "enter"
	keyBackspace  = "backspace"
)

var monthKeyRe = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])$`)

func NewModel(months []data.Month, startKey string, err error) Model {
	m := Model{months: months, err: err, helpModel: help.New()}
	m.state.Width = 100
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
	m.state.CursorDay = clamp(time.Now().Day(), 1, m.currentMonthDayCount())
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.state.Width = msg.Width
	case tea.KeyMsg:
		if m.state.Mode == InputModeGotoMonth {
			switch msg.String() {
			case keyQuit, keyQuitAlt:
				return m, tea.Quit
			default:
				m.updateGotoMode(msg)
				return m, nil
			}
		}

		switch msg.String() {
		case keyQuit, keyQuitAlt, keyQuitAlt2, keyQuitAlt3:
			return m, tea.Quit
		}

		switch msg.String() {
		case keyNextMonth, keyNextMonth2:
			m.moveMonth(1)
		case keyPrevMonth, keyPrevMonth2:
			m.moveMonth(-1)
		case keyRight, keyRightAlt:
			m.setCursorDay(m.state.CursorDay + 1)
		case keyLeft, keyLeftAlt:
			m.setCursorDay(m.state.CursorDay - 1)
		case keyDown, keyDownAlt:
			m.setCursorDay(m.state.CursorDay + 7)
		case keyUp, keyUpAlt:
			m.setCursorDay(m.state.CursorDay - 7)
		case keyToday:
			now := time.Now()
			m.jumpToKeyDay(now.Format("2006-01"), now.Day())
		case keyGoto:
			m.state.Mode = InputModeGotoMonth
			m.state.GotoBuffer = ""
			m.state.GotoError = ""
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

	selectedDay := m.state.CursorDay
	if selectedDay < 1 {
		selectedDay = 1
	}

	next := clamp(m.state.MonthIndex+delta, 0, len(m.months)-1)
	if next == m.state.MonthIndex {
		return
	}
	m.state.MonthIndex = next
	m.setCursorDay(selectedDay)
}

func (m *Model) setCursorDay(day int) {
	maxDay := m.currentMonthDayCount()
	if maxDay <= 0 {
		m.state.CursorDay = 1
		return
	}
	m.state.CursorDay = clamp(day, 1, maxDay)
}

func (m *Model) currentMonthDayCount() int {
	if len(m.months) == 0 {
		return 0
	}
	days := m.currentDays()
	if len(days) > 0 {
		first := days[0].Date
		return time.Date(first.Year(), first.Month()+1, 0, 0, 0, 0, 0, first.Location()).Day()
	}
	parsed, err := time.Parse("2006-01", m.months[m.state.MonthIndex].Key)
	if err != nil {
		return 31
	}
	return time.Date(parsed.Year(), parsed.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
}

func (m *Model) currentDateForCursor() time.Time {
	if len(m.months) == 0 {
		return time.Time{}
	}
	days := m.currentDays()
	if len(days) > 0 {
		first := days[0].Date
		return time.Date(first.Year(), first.Month(), m.state.CursorDay, 0, 0, 0, 0, first.Location())
	}
	parsed, err := time.Parse("2006-01", m.months[m.state.MonthIndex].Key)
	if err != nil {
		return time.Time{}
	}
	return time.Date(parsed.Year(), parsed.Month(), m.state.CursorDay, 0, 0, 0, 0, time.Local)
}

func (m *Model) currentSelectedEntry() (data.Day, bool) {
	for _, day := range m.currentDays() {
		if day.Date.Day() == m.state.CursorDay {
			return day, true
		}
	}
	return data.Day{}, false
}

func (m *Model) jumpToKeyDay(key string, day int) {
	if len(m.months) == 0 {
		return
	}
	idx := -1
	for i, month := range m.months {
		if month.Key == key {
			idx = i
			break
		}
	}
	if idx == -1 {
		idx = lastMonthIndexBefore(m.months, key)
	}
	m.state.MonthIndex = idx
	m.setCursorDay(day)
}

func (m *Model) updateGotoMode(msg tea.KeyMsg) {
	switch msg.String() {
	case keyQuitAlt2, keyQuitAlt3:
		m.state.Mode = InputModeNormal
		m.state.GotoBuffer = ""
		m.state.GotoError = ""
		return
	case keyEnter:
		if !monthKeyRe.MatchString(m.state.GotoBuffer) {
			m.state.GotoError = "Invalid format (use YYYY-MM)"
			return
		}
		m.jumpToKeyDay(m.state.GotoBuffer, m.state.CursorDay)
		m.state.Mode = InputModeNormal
		m.state.GotoBuffer = ""
		m.state.GotoError = ""
		return
	case keyBackspace:
		if len(m.state.GotoBuffer) > 0 {
			m.state.GotoBuffer = m.state.GotoBuffer[:len(m.state.GotoBuffer)-1]
			m.state.GotoError = ""
		}
		return
	}

	r := msg.String()
	if len(r) != 1 {
		return
	}
	ch := r[0]
	if (ch >= '0' && ch <= '9') || ch == '-' {
		if len(m.state.GotoBuffer) < 7 {
			m.state.GotoBuffer += r
			m.state.GotoError = ""
		}
	}
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
