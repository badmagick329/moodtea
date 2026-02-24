package data

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAllFiltersAndSortsMonthFiles(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, filepath.Join(dir, "2026-02.json"), `[{"date":"2026-02-01","mood":2,"energy":3}]`)
	writeJSON(t, filepath.Join(dir, "2026-01.json"), `[{"date":"2026-01-01","mood":5,"energy":4}]`)
	writeJSON(t, filepath.Join(dir, "january_2026.json"), `[{"date":"2026-01-01","mood":5,"energy":4}]`)

	months, err := LoadAll(dir)
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	if len(months) != 2 {
		t.Fatalf("len(months) = %d, want 2", len(months))
	}
	if months[0].Key != "2026-01" || months[1].Key != "2026-02" {
		t.Fatalf("unexpected sort order: %+v", []string{months[0].Key, months[1].Key})
	}
}

func TestLoadAllErrorsWithNoValidMonthFiles(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, filepath.Join(dir, "january_2026.json"), `[{"date":"2026-01-01","mood":5,"energy":4}]`)
	_, err := LoadAll(dir)
	if err == nil {
		t.Fatal("expected error with no valid month files")
	}
}

func writeJSON(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
