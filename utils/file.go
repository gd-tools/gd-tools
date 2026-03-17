package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadJSON reads a JSON file into the provided struct or slice.
func LoadJSON(name string, v any) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("decode %s: %w", name, err)
	}

	return nil
}

// SaveFile writes data to name using atomic tmp+rename.
// If the existing file already contains identical data,
// the write is skipped.
func SaveFile(name string, data []byte) error {
	old, err := os.ReadFile(name)
	if err == nil && bytes.Equal(old, data) {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", name, err)
	}

	dir := filepath.Dir(name)
	base := filepath.Base(name)

	tmp, err := os.CreateTemp(dir, base+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	tmpName := tmp.Name()
	ok := false
	defer func() {
		if !ok {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, name); err != nil {
		return err
	}

	ok = true
	return nil
}

// SaveJSON writes a struct or slice as formatted JSON to a file.
func SaveJSON(name string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", name, err)
	}

	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}

	return SaveFile(name, data)
}
