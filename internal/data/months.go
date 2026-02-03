package data

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Month struct {
	Key  string
	Days []Day
}

var monthFileRe = regexp.MustCompile(`^\d{4}-\d{2}\.json$`)

func LoadAll(dir string) ([]Month, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var months []Month
	for _, e := range entries {
		if e.IsDir() || !monthFileRe.MatchString(e.Name()) {
			continue
		}
		path := filepath.Join(dir, e.Name())
		days, err := Load(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		key := strings.TrimSuffix(e.Name(), ".json")
		months = append(months, Month{Key: key, Days: days})
	}

	sort.Slice(months, func(i, j int) bool {
		return months[i].Key < months[j].Key
	})

	if len(months) == 0 {
		return nil, fmt.Errorf("no month files found in %s", dir)
	}

	return months, nil
}
