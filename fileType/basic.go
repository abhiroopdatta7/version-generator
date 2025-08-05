package filetype

import (
	"os"
	"path/filepath"
)

type BasicFile struct {
}

func (b *BasicFile) WriteVersion(filePath string, version string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	// Write file (this will overwrite existing file)
	return os.WriteFile(filePath, []byte(version+"\n"), 0644)
}
