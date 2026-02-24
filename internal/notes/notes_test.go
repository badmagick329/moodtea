package notes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildMonthlyDataWithDiagnostics(t *testing.T) {
	root := t.TempDir()

	writeNote(t, filepath.Join(root, "2026", "02", "01.txt"), "@mood 3\n@energy 4\n")
	writeNote(t, filepath.Join(root, "2026", "02", "02.txt"), "@mood 3\n")
	writeNote(t, filepath.Join(root, "2026", "02", "03.txt"), "@mood 9\n@energy 1\n")
	writeNote(t, filepath.Join(root, "2026", "02", "31.txt"), "@mood 3\n@energy 2\n")
	writeNote(t, filepath.Join(root, "2026", "02", "04.txt"), "@mood 999999999999999999999999999\n@energy 2\n")
	writeNote(t, filepath.Join(root, "2026", "02", "note.md"), "@mood 5\n@energy 5\n")

	if err := os.MkdirAll(filepath.Join(root, "2026", "02", "subdir"), 0o755); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "junk"), 0o755); err != nil {
		t.Fatalf("mkdir junk: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "2026", "13"), 0o755); err != nil {
		t.Fatalf("mkdir invalid month: %v", err)
	}

	months, diag, err := BuildMonthlyDataWithDiagnostics(root)
	if err != nil {
		t.Fatalf("BuildMonthlyDataWithDiagnostics: %v", err)
	}

	if diag.YearsScanned != 1 || diag.MonthsScanned != 1 {
		t.Fatalf("unexpected year/month counts: %+v", diag)
	}
	if diag.FilesScanned != 5 {
		t.Fatalf("FilesScanned = %d, want 5", diag.FilesScanned)
	}
	if diag.Imported != 1 {
		t.Fatalf("Imported = %d, want 1", diag.Imported)
	}
	if diag.Skipped != 8 {
		t.Fatalf("Skipped = %d, want 8", diag.Skipped)
	}

	days := months["2026-02"]
	if len(days) != 1 {
		t.Fatalf("len(months[2026-02]) = %d, want 1", len(days))
	}
	if days[0].Mood != 3 || days[0].Energy != 4 {
		t.Fatalf("unexpected imported day: %+v", days[0])
	}

	assertReasonCount(t, diag, SkipInvalidYearDir, 1)
	assertReasonCount(t, diag, SkipInvalidMonthDir, 1)
	assertReasonCount(t, diag, SkipInvalidDayFile, 2)
	assertReasonCount(t, diag, SkipMissingTag, 1)
	assertReasonCount(t, diag, SkipOutOfRange, 1)
	assertReasonCount(t, diag, SkipInvalidDate, 1)
	assertReasonCount(t, diag, SkipParseError, 1)

	if len(diag.Details) != diag.Skipped {
		t.Fatalf("details=%d skipped=%d", len(diag.Details), diag.Skipped)
	}
}

func assertReasonCount(t *testing.T, diag ScanDiagnostics, reason SkipReason, want int) {
	t.Helper()
	if got := diag.ReasonCounts[reason]; got != want {
		t.Fatalf("%s count = %d, want %d", reason, got, want)
	}
}

func writeNote(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
