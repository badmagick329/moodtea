package config

import (
	"encoding/json"
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
	return cfg, nil
}
