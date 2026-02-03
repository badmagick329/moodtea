package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"moodtea/internal/data"
	"moodtea/internal/ui"
)

func main() {
	var path string
	flag.StringVar(&path, "file", "data/2026-01.json", "path to JSON data file")
	flag.Parse()

	dir := filepath.Dir(path)
	startKey := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	months, err := data.LoadAll(dir)
	m := ui.NewModel(months, startKey, err)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
