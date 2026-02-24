package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"moodtea/internal/config"
	"moodtea/internal/notes"
)

func main() {
	var configPath string
	var notesPathFlag string
	var outDir string
	var verbose bool
	flag.StringVar(&configPath, "config", "config.json", "path to config.json")
	flag.StringVar(&notesPathFlag, "notes-path", "", "notes root path (overrides config notes_path)")
	flag.StringVar(&outDir, "out", "data", "output directory for month JSON files")
	flag.BoolVar(&verbose, "verbose", false, "print per-file skip diagnostics")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(1)
	}
	notesPath := cfg.NotesPath
	if notesPathFlag != "" {
		notesPath = notesPathFlag
	}
	if notesPath == "" {
		fmt.Fprintln(os.Stderr, "config error: notes_path is required")
		os.Exit(1)
	}

	months, diag, err := notes.BuildMonthlyDataWithDiagnostics(notesPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scan error:", err)
		os.Exit(1)
	}

	written, err := notes.WriteMonthlyJSON(months, outDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "write error:", err)
		os.Exit(1)
	}

	fmt.Printf(
		"Scan summary: years=%d months=%d files=%d imported=%d skipped=%d\n",
		diag.YearsScanned,
		diag.MonthsScanned,
		diag.FilesScanned,
		diag.Imported,
		diag.Skipped,
	)

	if len(diag.ReasonCounts) > 0 {
		reasons := make([]string, 0, len(diag.ReasonCounts))
		for reason := range diag.ReasonCounts {
			reasons = append(reasons, string(reason))
		}
		sort.Strings(reasons)
		fmt.Println("Skip reasons:")
		for _, reason := range reasons {
			fmt.Printf("  %s: %d\n", reason, diag.ReasonCounts[notes.SkipReason(reason)])
		}
	}

	if verbose && len(diag.Details) > 0 {
		grouped := make(map[notes.SkipReason][]string)
		for _, detail := range diag.Details {
			grouped[detail.Reason] = append(grouped[detail.Reason], detail.Path)
		}
		reasons := make([]string, 0, len(grouped))
		for reason := range grouped {
			reasons = append(reasons, string(reason))
		}
		sort.Strings(reasons)
		fmt.Println("Detailed skips:")
		for _, reason := range reasons {
			fmt.Printf("  %s:\n", reason)
			paths := grouped[notes.SkipReason(reason)]
			sort.Strings(paths)
			for _, path := range paths {
				fmt.Printf("    - %s\n", path)
			}
		}
	}

	fmt.Printf("Wrote %d month file(s) to %s\n", written, outDir)
}
