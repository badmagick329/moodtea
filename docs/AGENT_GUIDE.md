# MoodTea Agent Guide

## Purpose
MoodTea is a terminal app for viewing daily mood and energy scores (1-5) on month calendars.
It has two executables:
- `cmd/importer`: scans note files and generates normalized month JSON files.
- `cmd/moodtea`: interactive CLI viewer for those month JSON files.

This guide is the baseline context a future agent should read before changing code.

## Repository Layout
- `cmd/moodtea/main.go`: app entrypoint for the TUI viewer.
- `cmd/importer/main.go`: entrypoint for note-to-JSON import pipeline.
- `internal/config/config.go`: JSON config loading and validation.
- `internal/data/data.go`: parses month JSON files into `[]Day` with validation.
- `internal/data/months.go`: loads all `YYYY-MM.json` files and sorts by month key.
- `internal/notes/notes.go`: scans note hierarchy, extracts `@mood`/`@energy`, writes month JSON.
- `internal/ui/*.go`: Bubble Tea model/update/view and calendar rendering.
- `data/`: generated month JSON files used by the viewer.

## Runtime Flow
### Import flow (`cmd/importer`)
1. Load `config.json` (or `-config` override).
2. Require `notes_path`.
3. Scan notes directory with structure `ROOT/YYYY/MM/DD.txt`.
4. Parse each day file for `@mood` and `@energy` tags.
5. Keep only valid scores in range 1-5.
6. Write JSON files to output dir (`-out`, default `data`) as `YYYY-MM.json`.

### Viewer flow (`cmd/moodtea`)
1. Load `config.json` (or `-config` override).
2. Require `data_path`.
3. Load all month files in `data_path` that match `^\d{4}-\d{2}\.json$`.
4. Pick initial month key based on current date (`time.Now().Format("2006-01")`), fallback to closest previous available month.
5. Start Bubble Tea program and render two calendars (Mood and Energy), legends, monthly averages, and selected-day details.

## Data Contracts
### Config schema
`config.json` fields:
- `notes_path` (string): used by importer.
- `data_path` (string): used by viewer.

Validation behavior:
- If both empty, config load fails.
- Importer separately enforces `notes_path` non-empty.
- Viewer separately enforces `data_path` non-empty.

Example (`config.example.json`):
```json
{
  "notes_path": "C:/path/to/your/notes/root",
  "data_path": "C:/code/golang/moodtea/data"
}
```

### Month data schema
Each `data/YYYY-MM.json` is an array of objects:
```json
[
  {
    "date": "2026-02-24",
    "mood": 2,
    "energy": 3
  }
]
```
Rules enforced by loader:
- `date` format must be `YYYY-MM-DD`.
- `mood` and `energy` must be integers 1-5.
- Empty month files are rejected.

File selection rule:
- Only files named exactly `YYYY-MM.json` are loaded.
- Files like `january_2026.json` are ignored.

## UI Architecture
Key files:
- `internal/ui/model.go`: app state and key handling.
- `internal/ui/viewmodel.go`: transforms model state into render-ready `ViewModel`.
- `internal/ui/calendar.go`: day map, month bounds, and grid rendering.
- `internal/ui/view.go`: composes final output.

State model:
- `MonthIndex`: active month in loaded months slice.
- `Cursor`: selected day index within the month's existing entries.

Navigation:
- Quit: `q`, `esc`, `ctrl+c`
- Month: `[`/`]` or `pgup`/`pgdown`
- Day: arrows or `h/j/k/l`

Rendering details:
- Two 6x7 calendars are always rendered: Mood and Energy.
- Cell background color comes from palette by score.
- Missing days are rendered uncolored.
- Selected day is bold + underlined.

Monthly averages (implemented):
- Computed from recorded entries only.
- `avgMood = sum(mood) / recordedDays`
- `avgEnergy = sum(energy) / recordedDays`
- Displayed with one decimal under month header:
  - `Avg mood X.X | Avg energy Y.Y (N days)`

## Importer Parsing Rules
Notes scanner expects:
- Year directory: `YYYY`
- Month directory: `MM` (`01`-`12`)
- Day file: `DD.txt` (`01`-`31`)

Tag regexes (case-sensitive):
- Mood: `@mood\s*[:=]?\s*(\d+)`
- Energy: `@energy\s*[:=]?\s*(\d+)`

Behavior:
- If either tag missing, the note is skipped.
- If parse fails or values out of 1-5, the note is skipped.
- Invalid calendar dates are skipped (e.g., 2026-02-31).

## Build, Run, Test
Prereq:
- Go `1.25.6` (from `go.mod`).

Common commands:
```powershell
# Run importer
go run ./cmd/importer -config config.json -out data

# Run viewer
go run ./cmd/moodtea -config config.json

# Run all tests
go test ./...
```

Current tests:
- `internal/ui/viewmodel_test.go`: monthly average calculations.
- `internal/ui/view_test.go`: average line rendered + updates on month switch.

## Dependencies
Core runtime dependencies are Charm ecosystem packages:
- `charm.land/bubbletea/v2`
- `charm.land/lipgloss/v2`

Most dependencies in `go.mod` are currently marked indirect.

## Known Constraints and Gotchas
- The viewer only loads month files matching strict `YYYY-MM.json` naming.
- Date parsing in data loader is strict and fail-fast per file.
- Importer skips invalid or incomplete note files silently (no per-file warning output).
- UI output includes ANSI styling; tests should assert with substring matching rather than full snapshot equality.
- `.gitignore` excludes `data/*`, `config.json`, and `*.exe`.

## Safe Extension Points
- Add derived month stats in `internal/ui/viewmodel.go` and surface via `ViewModel`.
- Add new navigation behaviors in `internal/ui/model.go` key handlers.
- Add new note tags in `internal/notes/notes.go` with explicit regex + validation.
- Keep data contract strict in `internal/data/data.go` to avoid downstream UI surprises.

## Suggested Agent Workflow Before Editing
1. Read this file.
2. Read target subsystem files (`internal/ui`, `internal/notes`, or `internal/data`).
3. Confirm data contract impact (config keys, JSON schema, file naming).
4. Implement minimal changes in one subsystem at a time.
5. Add/adjust tests in `internal/ui` (or new package tests as needed).
6. Run `go test ./...` before finishing.
