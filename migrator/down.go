package migrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Down applies all down migrations in the scripts directory that have been applied
func (m *Migrator) Down() error {
	files, err := os.ReadDir(m.ScriptDir)
	if err != nil {
		return fmt.Errorf("failed to read script directory: %w", err)
	}

	appliedMigrations, err := m.appliedMigrations()
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".js") && strings.Contains(file.Name(), "down") {
			version := extractVersionFromFileName(file.Name())
			if !appliedMigrations[version] {
				fmt.Printf("Migration %s not applied, skipping.\n", version)
				continue
			}

			scriptPath := filepath.Join(m.ScriptDir, file.Name())
			if err := m.runScript(scriptPath); err != nil {
				return err
			}

			// Remove the migration record after a successful rollback
			if err := m.removeMigrationRecord(version); err != nil {
				return err
			}
		}
	}

	return nil
}
