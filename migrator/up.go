package migrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Up applies all up migrations in the scripts directory that haven't been applied yet
func (m *Migrator) Up() error {
	files, err := os.ReadDir(m.ScriptDir)
	if err != nil {
		return fmt.Errorf("failed to read script directory: %w", err)
	}

	appliedMigrations, err := m.appliedMigrations()
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".js") && strings.Contains(file.Name(), "up") {
			version := extractVersionFromFileName(file.Name())
			if appliedMigrations[version] {
				fmt.Printf("Migration %s already applied, skipping.\n", version)
				continue
			}

			scriptPath := filepath.Join(m.ScriptDir, file.Name())
			if err := m.runScript(scriptPath); err != nil {
				return err
			}

			if err := m.recordMigration(version); err != nil {
				return err
			}
		}
	}

	return nil
}
