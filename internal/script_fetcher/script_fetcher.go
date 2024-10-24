package script_fetcher

import (
	"github.com/kondukto-io/migrongo/internal/script_fetcher/file_path"
)

func NewScriptFetcher(config Config) (ScriptFetcher, error) {
	switch config.Type {
	default:
		return file_path.NewFilePath(config.FilePathConfig.Dir), nil
	}
}

type (
	Config struct {
		Type           string
		FilePathConfig *FilePathConfig
	}

	FilePathConfig struct {
		Dir string
	}
)
