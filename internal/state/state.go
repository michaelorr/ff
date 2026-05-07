// Package state handles serialization / deserialization of application state
package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	filename       = ".ff"
	currentVersion = 1
)

type AppState struct {
	Version int      `json:"version"`
	Query   string   `json:"query"`
	Results []string `json:"results"`
}

func Load() (AppState, bool) {
	data, err := os.ReadFile(filepath.Join(".", filename))
	if err != nil {
		return AppState{}, false
	}
	var s AppState
	if err := json.Unmarshal(data, &s); err != nil {
		return AppState{}, false
	}
	return s, true
}

func Save(s AppState) error {
	s.Version = currentVersion
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(".", filename), data, 0o644)
}
