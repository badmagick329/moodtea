package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	NotesPath string `json:"notes_path"`
	DataPath  string `json:"data_path"`
}

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}
	cfg.NotesPath = strings.TrimSpace(cfg.NotesPath)
	cfg.DataPath = strings.TrimSpace(cfg.DataPath)
	if cfg.NotesPath == "" && cfg.DataPath == "" {
		return Config{}, fmt.Errorf("config has no notes_path or data_path in %s", path)
	}
	return cfg, nil
}
