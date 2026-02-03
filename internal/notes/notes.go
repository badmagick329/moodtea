package notes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"moodtea/internal/data"
)

var (
	yearRe    = regexp.MustCompile(`^\d{4}$`)
	monthRe   = regexp.MustCompile(`^(0[1-9]|1[0-2])$`)
	dayFileRe = regexp.MustCompile(`^(0[1-9]|[12]\d|3[01])\.txt$`)
	moodRe    = regexp.MustCompile(`@mood\s*[:=]?\s*(\d+)`)
	energyRe  = regexp.MustCompile(`@energy\s*[:=]?\s*(\d+)`)
)

func BuildMonthlyData(root string) (map[string][]data.Day, error) {
	months := make(map[string][]data.Day)

	years, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, y := range years {
		if !y.IsDir() || !yearRe.MatchString(y.Name()) {
			continue
		}
		yearPath := filepath.Join(root, y.Name())
		monthDirs, err := os.ReadDir(yearPath)
		if err != nil {
			return nil, err
		}

		yearInt, err := strconv.Atoi(y.Name())
		if err != nil {
			return nil, err
		}

		for _, m := range monthDirs {
			if !m.IsDir() || !monthRe.MatchString(m.Name()) {
				continue
			}
			monthPath := filepath.Join(yearPath, m.Name())
			files, err := os.ReadDir(monthPath)
			if err != nil {
				return nil, err
			}

			monthInt, err := strconv.Atoi(m.Name())
			if err != nil {
				return nil, err
			}
			key := fmt.Sprintf("%04d-%02d", yearInt, monthInt)

			for _, f := range files {
				if f.IsDir() || !dayFileRe.MatchString(f.Name()) {
					continue
				}
				dayStr := strings.TrimSuffix(f.Name(), ".txt")
				dayInt, err := strconv.Atoi(dayStr)
				if err != nil {
					return nil, err
				}

				date := time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.Local)
				if date.Year() != yearInt || int(date.Month()) != monthInt || date.Day() != dayInt {
					continue
				}

				path := filepath.Join(monthPath, f.Name())
				b, err := os.ReadFile(path)
				if err != nil {
					return nil, err
				}

				mood, energy, ok := parseMoodEnergy(string(b))
				if !ok {
					continue
				}

				months[key] = append(months[key], data.Day{
					Date:   date,
					Mood:   mood,
					Energy: energy,
				})
			}
		}
	}

	return months, nil
}

func WriteMonthlyJSON(months map[string][]data.Day, outDir string) (int, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return 0, err
	}

	keys := make([]string, 0, len(months))
	for k := range months {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	written := 0
	for _, key := range keys {
		days := months[key]
		if len(days) == 0 {
			continue
		}
		sort.Slice(days, func(i, j int) bool {
			return days[i].Date.Before(days[j].Date)
		})

		type rawDay struct {
			Date   string `json:"date"`
			Mood   int    `json:"mood"`
			Energy int    `json:"energy"`
		}
		out := make([]rawDay, 0, len(days))
		for _, d := range days {
			out = append(out, rawDay{
				Date:   d.Date.Format("2006-01-02"),
				Mood:   d.Mood,
				Energy: d.Energy,
			})
		}

		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return written, err
		}
		outPath := filepath.Join(outDir, key+".json")
		if err := os.WriteFile(outPath, b, 0o644); err != nil {
			return written, err
		}
		written++
	}

	return written, nil
}

func parseMoodEnergy(s string) (int, int, bool) {
	mm := moodRe.FindStringSubmatch(s)
	em := energyRe.FindStringSubmatch(s)
	if mm == nil || em == nil {
		return 0, 0, false
	}
	mood, err := strconv.Atoi(mm[1])
	if err != nil {
		return 0, 0, false
	}
	energy, err := strconv.Atoi(em[1])
	if err != nil {
		return 0, 0, false
	}
	if mood < 1 || mood > 5 || energy < 1 || energy > 5 {
		return 0, 0, false
	}
	return mood, energy, true
}
