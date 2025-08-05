package filetype

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type YAMLFile struct {
}

func (y *YAMLFile) WriteVersion(filePath string, version string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	// Write file (this will overwrite existing file)
	data := map[string]string{"version": version}
	out, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, out, 0644)
}
