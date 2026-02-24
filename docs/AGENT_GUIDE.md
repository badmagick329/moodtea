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
2. Resolve notes root with precedence `-notes-path` flag, then `notes_path`.
3. Scan notes directory with structure `ROOT/YYYY/MM/DD.txt`.
4. Parse each day file for `@mood` and `@energy` tags.
5. Keep only valid scores in range 1-5.
6. Write JSON files to output dir (`-out`, default `data`) as `YYYY-MM.json`.
7. Print scan diagnostics summary; `-verbose` prints per-file skips grouped by reason.

### Viewer flow (`cmd/moodtea`)
1. Load `config.json` (or `-config` override).
2. Resolve data directory with precedence `-data-path` flag, then `data_path`.
3. Load all month files in `data_path` that match `^\d{4}-\d{2}\.json$`.
4. Pick initial month key based on current date (`time.Now().Format("2006-01")`), fallback to closest previous available month.
5. Start Bubble Tea program and render two calendars (Mood and Energy), legends, monthly averages, and selected-day details.

## Data Contracts
### Config schema
`config.json` fields:
- `notes_path` (string): used by importer.
- `data_path` (string): used by viewer.

Validation behavior:
- Config load trims values and parses JSON.
- Importer enforces resolved notes path non-empty.
- Viewer enforces resolved data path non-empty.

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
- `CursorDay`: selected calendar day within active month (`1..daysInMonth`).
- `Mode`: normal navigation vs go-to-month input mode.

Navigation:
- Quit: `q`, `esc`, `ctrl+c`
- Month: `[`/`]` or `pgup`/`pgdown`
- Day: arrows or `h/j/k/l` (calendar-day movement; clamps within month)
- Today: `t` (jump to today key/day with closest-previous month fallback)
- Month jump: `g` (input `YYYY-MM`, Enter to jump, Esc to cancel)

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

Additional trend stats (implemented):
- Median mood/energy from recorded entries.
- Min/max mood and energy from recorded entries.
- Rolling 7-entry average using latest up to 7 recorded entries in month.
- Display line:
  - `Median M/E: X.X / Y.Y | Min-Max M: A-B E: C-D | 7d Avg M/E: P.P / Q.Q`

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
- Importer records skip reasons and counts:
  - `invalid_year_dir`, `invalid_month_dir`, `invalid_day_filename`, `invalid_date`,
    `read_error`, `missing_tag`, `parse_error`, `out_of_range`.

## Build, Run, Test
Prereq:
- Go `1.25.6` (from `go.mod`).

Common commands:
```powershell
# Run importer
go run ./cmd/importer -config config.json -out data
go run ./cmd/importer -config config.json -notes-path C:/path/to/notes -verbose

# Run viewer
go run ./cmd/moodtea -config config.json
go run ./cmd/moodtea -config config.json -data-path C:/code/golang/moodtea/data

# Run all tests
go test ./...
```

Current tests:
- `internal/notes/notes_test.go`: importer scan diagnostics, skip reasons, and valid import handling.
- `internal/data/data_test.go`: strict date/range/empty validations.
- `internal/data/months_test.go`: strict month filename filtering and key sort behavior.
- `internal/ui/model_test.go`: calendar-day cursor behavior, no-entry rendering, `t`/`g` shortcuts.
- `internal/ui/viewmodel_test.go`: monthly average calculations.
- `internal/ui/view_test.go`: average/trend lines rendered + updates on month switch.

## Dependencies
Core runtime dependencies are Charm ecosystem packages:
- `charm.land/bubbletea/v2`
- `charm.land/lipgloss/v2`

Most dependencies in `go.mod` are currently marked indirect.

## Known Constraints and Gotchas
- The viewer only loads month files matching strict `YYYY-MM.json` naming.
- Date parsing in data loader is strict and fail-fast per file.
- Importer skip behavior is non-fatal at file level and surfaced through diagnostics output.
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
7. Run `golangci-lint run` and `go build ./cmd/importer ./cmd/moodtea` for quality gates.
