package filetype

import (
	"os"
	"path/filepath"
)

type GoType struct {
}

func (g *GoType) WriteVersion(filePath string, version string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	// Write file (this will overwrite existing file)
	data := "package main\n\nconst Version = \"" + version + "\"\n"
	return os.WriteFile(filePath, []byte(data), 0644)
}
