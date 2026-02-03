package main

import (
	"flag"
	"fmt"
	"os"

	"moodtea/internal/config"
	"moodtea/internal/notes"
)

func main() {
	var configPath string
	var outDir string
	flag.StringVar(&configPath, "config", "config.json", "path to config.json")
	flag.StringVar(&outDir, "out", "data", "output directory for month JSON files")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(1)
	}
	if cfg.NotesPath == "" {
		fmt.Fprintln(os.Stderr, "config error: notes_path is required")
		os.Exit(1)
	}

	months, err := notes.BuildMonthlyData(cfg.NotesPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scan error:", err)
		os.Exit(1)
	}

	written, err := notes.WriteMonthlyJSON(months, outDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "write error:", err)
		os.Exit(1)
	}

	fmt.Printf("Wrote %d month file(s) to %s\n", written, outDir)
}
