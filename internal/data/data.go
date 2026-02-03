package data

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Day struct {
	Date   time.Time
	Mood   int
	Energy int
}

type rawDay struct {
	Date   string `json:"date"`
	Mood   int    `json:"mood"`
	Energy int    `json:"energy"`
}

func Load(path string) ([]Day, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw []rawDay
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("no entries found in %s", path)
	}

	out := make([]Day, 0, len(raw))
	for i, r := range raw {
		dt, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid date %q (want YYYY-MM-DD): %w", i, r.Date, err)
		}
		if r.Mood < 1 || r.Mood > 5 {
			return nil, fmt.Errorf("row %d (%s): mood out of range (1-5): %d", i, r.Date, r.Mood)
		}
		if r.Energy < 1 || r.Energy > 5 {
			return nil, fmt.Errorf("row %d (%s): energy out of range (1-5): %d", i, r.Date, r.Energy)
		}
		out = append(out, Day{Date: dt, Mood: r.Mood, Energy: r.Energy})
	}

	return out, nil
}
