package ui

import (
	"fmt"
	"strings"
	"testing"

	"moodtea/internal/data"
)

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

	out := fmt.Sprintf("%v", m.View())
	if !strings.Contains(out, "Avg mood 4.0 | Avg energy 3.0 (2 days)") {
		t.Fatalf("view missing average line, got:\n%s", out)
	}
	if !strings.Contains(out, "Median M/E: 4.0 / 3.0 | Min-Max M: 3-5 E: 2-4 | 7d Avg M/E: 4.0 / 3.0") {
		t.Fatalf("view missing trend line, got:\n%s", out)
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

	first := fmt.Sprintf("%v", m.View())
	if !strings.Contains(first, "Avg mood 4.0 | Avg energy 3.0 (2 days)") {
		t.Fatalf("month 1 average line mismatch, got:\n%s", first)
	}

	m.moveMonth(1)
	second := fmt.Sprintf("%v", m.View())
	if !strings.Contains(second, "Avg mood 2.0 | Avg energy 4.0 (3 days)") {
		t.Fatalf("month 2 average line mismatch, got:\n%s", second)
	}
	if !strings.Contains(second, "Median M/E: 2.0 / 4.0 | Min-Max M: 1-3 E: 3-5 | 7d Avg M/E: 2.0 / 4.0") {
		t.Fatalf("month 2 trend line mismatch, got:\n%s", second)
	}
}
