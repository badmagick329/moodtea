package data

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsInvalidDateFormat(t *testing.T) {
	path := writeMonthFile(t, `[{"date":"2026/02/01","mood":3,"energy":4}]`)
	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "invalid date") {
		t.Fatalf("expected invalid date error, got %v", err)
	}
}

func TestLoadRejectsOutOfRangeValues(t *testing.T) {
	path := writeMonthFile(t, `[{"date":"2026-02-01","mood":6,"energy":4}]`)
	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "mood out of range") {
		t.Fatalf("expected mood range error, got %v", err)
	}
}

func TestLoadRejectsEmptyMonth(t *testing.T) {
	path := writeMonthFile(t, `[]`)
	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "no entries found") {
		t.Fatalf("expected no entries error, got %v", err)
	}
}

func writeMonthFile(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "2026-02.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	return path
}
