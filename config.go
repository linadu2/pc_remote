package main

import (
	"encoding/json"
	"os"
)

// Config holds the application configuration
type Config struct {
	Port  int    `json:"port"`
	Token string `json:"token"`
}

const ConfigFile = "config.json"

// LoadConfig attempts to read the config file.
// Returns the config and true if successful and valid.
// Returns empty config and false if file missing or invalid.
func LoadConfig() (Config, bool) {
	file, err := os.Open(ConfigFile)
	if err != nil {
		return Config{}, false
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, false
	}

	// Basic validation
	if cfg.Port < 1024 || cfg.Port > 65535 || len(cfg.Token) < 8 {
		return Config{}, false
	}

	return cfg, true
}

// SaveConfig writes the config to disk
func SaveConfig(cfg Config) error {
	file, err := os.Create(ConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(cfg)
}
