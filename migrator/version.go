package migrator

import (
	"fmt"
	"os"
	"strings"
)

// LatestVersion retrieves the latest migration version based on the script filenames
func (m *Migrator) LatestVersion() (string, error) {
	files, err := os.ReadDir(m.ScriptDir)
	if err != nil {
		return "", fmt.Errorf("failed to read script directory: %w", err)
	}

	latestVersion := ""
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".js") {
			version := extractVersionFromFileName(file.Name())
			if version > latestVersion {
				latestVersion = version
			}
		}
	}

	return latestVersion, nil
}

// extractVersionFromFileName extracts the version from the migration script filename
func extractVersionFromFileName(fileName string) string {
	parts := strings.Split(fileName, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
