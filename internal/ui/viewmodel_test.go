package ui

import (
	"testing"
	"time"

	"moodtea/internal/data"
)

func TestMonthAverages(t *testing.T) {
	t.Run("normal month", func(t *testing.T) {
		days := []data.Day{
			{Date: mustDate(t, "2026-01-01"), Mood: 5, Energy: 4},
			{Date: mustDate(t, "2026-01-02"), Mood: 3, Energy: 2},
			{Date: mustDate(t, "2026-01-03"), Mood: 4, Energy: 5},
		}

		avgMood, avgEnergy, recorded := monthAverages(days)

		if avgMood != 4.0 {
			t.Fatalf("avg mood = %v, want 4.0", avgMood)
		}
		if avgEnergy != 11.0/3.0 {
			t.Fatalf("avg energy = %v, want %v", avgEnergy, 11.0/3.0)
		}
		if recorded != 3 {
			t.Fatalf("recorded = %d, want 3", recorded)
		}
	})

	t.Run("sparse month", func(t *testing.T) {
		days := []data.Day{
			{Date: mustDate(t, "2026-02-01"), Mood: 2, Energy: 1},
			{Date: mustDate(t, "2026-02-15"), Mood: 4, Energy: 5},
		}

		avgMood, avgEnergy, recorded := monthAverages(days)

		if avgMood != 3.0 {
			t.Fatalf("avg mood = %v, want 3.0", avgMood)
		}
		if avgEnergy != 3.0 {
			t.Fatalf("avg energy = %v, want 3.0", avgEnergy)
		}
		if recorded != 2 {
			t.Fatalf("recorded = %d, want 2", recorded)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		avgMood, avgEnergy, recorded := monthAverages(nil)

		if avgMood != 0 {
			t.Fatalf("avg mood = %v, want 0", avgMood)
		}
		if avgEnergy != 0 {
			t.Fatalf("avg energy = %v, want 0", avgEnergy)
		}
		if recorded != 0 {
			t.Fatalf("recorded = %d, want 0", recorded)
		}
	})
}

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatalf("parse date %q: %v", s, err)
	}
	return d
}
