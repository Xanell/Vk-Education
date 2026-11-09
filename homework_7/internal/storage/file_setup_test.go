package storage

import (
	"path/filepath"
	"testing"
)

func testSetup(t *testing.T) (*FileStorage, string) {
	t.Helper()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_inventory.txt")
	fs := NewFileStorage(filePath)
	return fs, filePath
}
