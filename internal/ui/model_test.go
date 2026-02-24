package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"moodtea/internal/data"
)

func TestCalendarDayNavigationClampsWithinMonth(t *testing.T) {
	months := []data.Month{
		{
			Key: "2026-01",
			Days: []data.Day{
				{Date: mustDate(t, "2026-01-01"), Mood: 5, Energy: 4},
				{Date: mustDate(t, "2026-01-15"), Mood: 3, Energy: 2},
			},
		},
	}
	m := NewModel(months, "2026-01", nil)
	m.state.CursorDay = 1

	m = sendKey(m, keyCode(tea.KeyRight))
	if m.state.CursorDay != 2 {
		t.Fatalf("cursor day = %d, want 2", m.state.CursorDay)
	}

	m = sendKey(m, keyCode(tea.KeyUp))
	if m.state.CursorDay != 1 {
		t.Fatalf("cursor day = %d, want 1", m.state.CursorDay)
	}

	m = sendKey(m, keyCode(tea.KeyDown))
	if m.state.CursorDay != 8 {
		t.Fatalf("cursor day = %d, want 8", m.state.CursorDay)
	}
}

func TestSelectedMissingDayShowsNoEntry(t *testing.T) {
	months := []data.Month{
		{
			Key: "2026-01",
			Days: []data.Day{
				{Date: mustDate(t, "2026-01-01"), Mood: 5, Energy: 4},
				{Date: mustDate(t, "2026-01-15"), Mood: 3, Energy: 2},
			},
		},
	}
	m := NewModel(months, "2026-01", nil)
	m.state.CursorDay = 2

	out := fmt.Sprintf("%v", m.View())
	if !strings.Contains(out, "No entry for selected day") {
		t.Fatalf("expected no-entry message, got:\n%s", out)
	}
}

func TestTodayShortcutJumpsToCurrentMonthAndDay(t *testing.T) {
	now := time.Now()
	nowKey := now.Format("2006-01")
	months := []data.Month{
		{
			Key: "2025-12",
			Days: []data.Day{
				{Date: mustDate(t, "2025-12-01"), Mood: 2, Energy: 2},
			},
		},
		{
			Key: nowKey,
			Days: []data.Day{
				{Date: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local), Mood: 4, Energy: 4},
			},
		},
	}
	m := NewModel(months, "2025-12", nil)
	m = sendKey(m, keyText("t"))

	if m.months[m.state.MonthIndex].Key != nowKey {
		t.Fatalf("month key = %s, want %s", m.months[m.state.MonthIndex].Key, nowKey)
	}
	if m.state.CursorDay != now.Day() {
		t.Fatalf("cursor day = %d, want %d", m.state.CursorDay, now.Day())
	}
}

func TestGotoMonthInputValidationAndFallback(t *testing.T) {
	months := []data.Month{
		{
			Key: "2026-01",
			Days: []data.Day{
				{Date: mustDate(t, "2026-01-01"), Mood: 5, Energy: 4},
			},
		},
		{
			Key: "2026-03",
			Days: []data.Day{
				{Date: mustDate(t, "2026-03-01"), Mood: 3, Energy: 3},
			},
		},
	}
	m := NewModel(months, "2026-03", nil)
	m.state.CursorDay = 10

	m = sendKey(m, keyText("g"))
	for _, ch := range []string{"2", "0", "2", "6", "-", "1", "3"} {
		m = sendKey(m, keyText(ch))
	}
	m = sendKey(m, keyCode(tea.KeyEnter))
	if m.state.Mode != InputModeGotoMonth {
		t.Fatalf("mode = %v, want goto mode on invalid input", m.state.Mode)
	}
	if !strings.Contains(m.state.GotoError, "Invalid format") {
		t.Fatalf("unexpected goto error: %q", m.state.GotoError)
	}

	m = sendKey(m, keyCode(tea.KeyEsc))
	if m.state.Mode != InputModeNormal {
		t.Fatalf("mode = %v, want normal after escape", m.state.Mode)
	}

	m = sendKey(m, keyText("g"))
	for _, ch := range []string{"2", "0", "2", "6", "-", "0", "2"} {
		m = sendKey(m, keyText(ch))
	}
	m = sendKey(m, keyCode(tea.KeyEnter))
	if m.months[m.state.MonthIndex].Key != "2026-01" {
		t.Fatalf("month key = %s, want 2026-01 fallback", m.months[m.state.MonthIndex].Key)
	}
	if m.state.Mode != InputModeNormal {
		t.Fatalf("mode = %v, want normal after successful jump", m.state.Mode)
	}
}

func sendKey(m Model, msg tea.KeyMsg) Model {
	next, _ := m.Update(msg)
	return next.(Model)
}

func keyText(s string) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Text: s, Code: rune(s[0])})
}

func keyCode(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}
