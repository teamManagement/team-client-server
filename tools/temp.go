package tools

import (
	"os"
	"path/filepath"
)

func CreateTempFile() (*os.File, error) {
	tmpDir := filepath.Join(os.TempDir(), "team")
	_ = os.MkdirAll(tmpDir, 0755)
	return os.CreateTemp(tmpDir, "cache*")
}
