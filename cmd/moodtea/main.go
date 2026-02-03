package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"moodtea/internal/data"
	"moodtea/internal/ui"
)

func main() {
	var path string
	flag.StringVar(&path, "file", "data/january_2026.json", "path to JSON data file")
	flag.Parse()

	days, err := data.Load(path)
	m := ui.NewModel(days, err)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
