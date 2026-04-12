package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const dir = "checkpoints"

var nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9]`)

// Checkpoint persists the IMAP UID state for incremental mail processing.
type Checkpoint struct {
	UIDValidity uint32 `json:"uid_validity"`
	LastUID     uint32 `json:"last_uid"`
}

func filePath(host, username, inbox string) string {
	sanitize := func(s string) string {
		return nonAlphanumeric.ReplaceAllString(s, "_")
	}
	name := fmt.Sprintf("%s__%s__%s.json", sanitize(host), sanitize(username), sanitize(inbox))
	return filepath.Join(dir, name)
}

// Load returns the stored checkpoint for the given mailbox, or nil if none exists yet.
func Load(host, username, inbox string) (*Checkpoint, error) {
	path := filePath(host, username, inbox)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint: %w", err)
	}
	var cp Checkpoint
	if err = json.Unmarshal(data, &cp); err != nil {
		return nil, fmt.Errorf("failed to parse checkpoint: %w", err)
	}
	return &cp, nil
}

// Save writes the checkpoint for the given mailbox to persistent storage.
func Save(host, username, inbox string, cp *Checkpoint) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create checkpoint directory: %w", err)
	}
	path := filePath(host, username, inbox)
	data, err := json.Marshal(cp)
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %w", err)
	}
	if err = os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write checkpoint: %w", err)
	}
	return nil
}
