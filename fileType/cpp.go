package filetype

import (
	"os"
	"path/filepath"
)

type CPPType struct {
}

func (c *CPPType) WriteVersion(filePath string, version string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	// Write file (this will overwrite existing file)
	data := "#define VERSION \"" + version + "\"\n"
	return os.WriteFile(filePath, []byte(data), 0644)
}
