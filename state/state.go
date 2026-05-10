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
	Version int    `json:"version"`
	Query   string `json:"query"`
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

func Save(s AppState) {
	// TODO: handle errors
	s.Version = currentVersion
	data, _ := json.MarshalIndent(s, "", "  ")
	_ = os.WriteFile(filepath.Join(".", filename), data, 0o644)
}
