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

type SkipReason string

const (
	SkipInvalidYearDir  SkipReason = "invalid_year_dir"
	SkipInvalidMonthDir SkipReason = "invalid_month_dir"
	SkipInvalidDayFile  SkipReason = "invalid_day_filename"
	SkipInvalidDate     SkipReason = "invalid_date"
	SkipReadError       SkipReason = "read_error"
	SkipMissingTag      SkipReason = "missing_tag"
	SkipParseError      SkipReason = "parse_error"
	SkipOutOfRange      SkipReason = "out_of_range"
)

type SkipDetail struct {
	Path   string
	Reason SkipReason
}

type ScanDiagnostics struct {
	YearsScanned  int
	MonthsScanned int
	FilesScanned  int
	Imported      int
	Skipped       int
	ReasonCounts  map[SkipReason]int
	Details       []SkipDetail
}

func (d *ScanDiagnostics) addSkip(path string, reason SkipReason) {
	d.Skipped++
	if d.ReasonCounts == nil {
		d.ReasonCounts = make(map[SkipReason]int)
	}
	d.ReasonCounts[reason]++
	d.Details = append(d.Details, SkipDetail{Path: path, Reason: reason})
}

func BuildMonthlyData(root string) (map[string][]data.Day, error) {
	months, _, err := BuildMonthlyDataWithDiagnostics(root)
	return months, err
}

func BuildMonthlyDataWithDiagnostics(root string) (map[string][]data.Day, ScanDiagnostics, error) {
	months := make(map[string][]data.Day)
	diag := ScanDiagnostics{ReasonCounts: make(map[SkipReason]int)}

	years, err := os.ReadDir(root)
	if err != nil {
		return nil, diag, err
	}

	for _, y := range years {
		if !y.IsDir() || !yearRe.MatchString(y.Name()) {
			diag.addSkip(filepath.Join(root, y.Name()), SkipInvalidYearDir)
			continue
		}
		diag.YearsScanned++
		yearPath := filepath.Join(root, y.Name())
		monthDirs, err := os.ReadDir(yearPath)
		if err != nil {
			return nil, diag, err
		}

		yearInt, err := strconv.Atoi(y.Name())
		if err != nil {
			return nil, diag, err
		}

		for _, m := range monthDirs {
			if !m.IsDir() || !monthRe.MatchString(m.Name()) {
				diag.addSkip(filepath.Join(yearPath, m.Name()), SkipInvalidMonthDir)
				continue
			}
			diag.MonthsScanned++
			monthPath := filepath.Join(yearPath, m.Name())
			files, err := os.ReadDir(monthPath)
			if err != nil {
				return nil, diag, err
			}

			monthInt, err := strconv.Atoi(m.Name())
			if err != nil {
				return nil, diag, err
			}
			key := fmt.Sprintf("%04d-%02d", yearInt, monthInt)

			for _, f := range files {
				path := filepath.Join(monthPath, f.Name())
				if f.IsDir() || !dayFileRe.MatchString(f.Name()) {
					diag.addSkip(path, SkipInvalidDayFile)
					continue
				}
				diag.FilesScanned++
				dayStr := strings.TrimSuffix(f.Name(), ".txt")
				dayInt, err := strconv.Atoi(dayStr)
				if err != nil {
					return nil, diag, err
				}

				date := time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.Local)
				if date.Year() != yearInt || int(date.Month()) != monthInt || date.Day() != dayInt {
					diag.addSkip(path, SkipInvalidDate)
					continue
				}

				b, err := os.ReadFile(path)
				if err != nil {
					diag.addSkip(path, SkipReadError)
					continue
				}

				mood, energy, reason, ok := parseMoodEnergy(string(b))
				if !ok {
					diag.addSkip(path, reason)
					continue
				}

				months[key] = append(months[key], data.Day{
					Date:   date,
					Mood:   mood,
					Energy: energy,
				})
				diag.Imported++
			}
		}
	}

	return months, diag, nil
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

func parseMoodEnergy(s string) (int, int, SkipReason, bool) {
	mm := moodRe.FindStringSubmatch(s)
	em := energyRe.FindStringSubmatch(s)
	if mm == nil || em == nil {
		return 0, 0, SkipMissingTag, false
	}
	mood, err := strconv.Atoi(mm[1])
	if err != nil {
		return 0, 0, SkipParseError, false
	}
	energy, err := strconv.Atoi(em[1])
	if err != nil {
		return 0, 0, SkipParseError, false
	}
	if mood < 1 || mood > 5 || energy < 1 || energy > 5 {
		return 0, 0, SkipOutOfRange, false
	}
	return mood, energy, "", true
}
