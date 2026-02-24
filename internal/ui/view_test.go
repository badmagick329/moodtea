package ui

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"moodtea/internal/data"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestViewIncludesMonthlyAverages(t *testing.T) {
	months := []data.Month{
		{
			Key: "2026-01",
			Days: []data.Day{
				{Date: mustDate(t, "2026-01-01"), Mood: 5, Energy: 4},
				{Date: mustDate(t, "2026-01-02"), Mood: 3, Energy: 2},
			},
		},
	}
	m := NewModel(months, "2026-01", nil)
	m.state.CursorDay = 1

	out := stripANSI(fmt.Sprintf("%v", m.View()))
	if !strings.Contains(out, "MoodTea — January 2026") {
		t.Fatalf("view missing title line, got:\n%s", out)
	}
	if !strings.Contains(out, "2 days recorded") {
		t.Fatalf("view missing recorded-days badge, got:\n%s", out)
	}
	if !strings.Contains(out, "Avg M 4.0") || !strings.Contains(out, "Med E 3.0") {
		t.Fatalf("view missing stat chips, got:\n%s", out)
	}
	if !strings.Contains(out, "M 3-5") || !strings.Contains(out, "7d E 3.0") {
		t.Fatalf("view missing min/max or rolling chips, got:\n%s", out)
	}
	if !strings.Contains(out, "move day") || !strings.Contains(out, "scroll") {
		t.Fatalf("view missing help hints, got:\n%s", out)
	}
	if !strings.Contains(out, "2026-01-01") || !strings.Contains(out, "Mood █████ (5)") {
		t.Fatalf("view missing footer card selection header, got:\n%s", out)
	}
}

func TestViewAveragesUpdateWhenMonthChanges(t *testing.T) {
	months := []data.Month{
		{
			Key: "2026-01",
			Days: []data.Day{
				{Date: mustDate(t, "2026-01-01"), Mood: 5, Energy: 4},
				{Date: mustDate(t, "2026-01-02"), Mood: 3, Energy: 2},
			},
		},
		{
			Key: "2026-02",
			Days: []data.Day{
				{Date: mustDate(t, "2026-02-01"), Mood: 1, Energy: 5},
				{Date: mustDate(t, "2026-02-02"), Mood: 2, Energy: 4},
				{Date: mustDate(t, "2026-02-03"), Mood: 3, Energy: 3},
			},
		},
	}
	m := NewModel(months, "2026-01", nil)

	first := stripANSI(fmt.Sprintf("%v", m.View()))
	if !strings.Contains(first, "Avg M 4.0") {
		t.Fatalf("month 1 average chip mismatch, got:\n%s", first)
	}

	m.moveMonth(1)
	second := stripANSI(fmt.Sprintf("%v", m.View()))
	if !strings.Contains(second, "Avg M 2.0") || !strings.Contains(second, "Avg E 4.0") {
		t.Fatalf("month 2 average chip mismatch, got:\n%s", second)
	}
	if !strings.Contains(second, "Med M 2.0") || !strings.Contains(second, "M 1-3") {
		t.Fatalf("month 2 trend chip mismatch, got:\n%s", second)
	}
}

func TestViewGotoModeShowsPromptAndContextHelp(t *testing.T) {
	months := []data.Month{
		{
			Key: "2026-01",
			Days: []data.Day{
				{Date: mustDate(t, "2026-01-01"), Mood: 5, Energy: 4},
			},
		},
	}
	m := NewModel(months, "2026-01", nil)
	m.state.Mode = InputModeGotoMonth
	m.state.CursorDay = 1
	m.state.GotoBuffer = "2026-"
	m.state.GotoError = "Invalid format (use YYYY-MM)"

	out := stripANSI(fmt.Sprintf("%v", m.View()))
	if !strings.Contains(out, "Go to month (YYYY-MM):") || !strings.Contains(out, "2026-") {
		t.Fatalf("view missing goto prompt, got:\n%s", out)
	}
	if !strings.Contains(out, "jump") || !strings.Contains(out, "cancel") {
		t.Fatalf("view missing goto help hints, got:\n%s", out)
	}
	if !strings.Contains(out, "Invalid format (use YYYY-MM)") {
		t.Fatalf("view missing goto error, got:\n%s", out)
	}
}

func TestRenderStatChipsWrapsWhenWidthIsNarrow(t *testing.T) {
	vm := ViewModel{
		AvgMood:        3.4,
		AvgEnergy:      3.2,
		MedianMood:     4.0,
		MedianEnergy:   3.0,
		MinMood:        1,
		MaxMood:        5,
		MinEnergy:      1,
		MaxEnergy:      5,
		Rolling7Mood:   2.9,
		Rolling7Energy: 2.9,
	}

	out := renderStatChips(vm, 24)
	if !strings.Contains(out, "\n") {
		t.Fatalf("expected wrapped chips output, got:\n%s", out)
	}
}

func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}
