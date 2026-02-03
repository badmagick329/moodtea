package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"moodtea/internal/config"
	"moodtea/internal/data"
	"moodtea/internal/ui"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "path to config.json")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(1)
	}
	if cfg.DataPath == "" {
		fmt.Fprintln(os.Stderr, "config error: data_path is required")
		os.Exit(1)
	}

	months, err := data.LoadAll(cfg.DataPath)
	startKey := time.Now().Format("2006-01")
	m := ui.NewModel(months, startKey, err)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
