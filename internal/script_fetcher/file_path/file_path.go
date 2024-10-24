package file_path

import (
	"fmt"
	"os"

	"github.com/kondukto-io/migrongo/internal/script_fetcher"
)

type FilePath struct {
	Dir string
}

func NewFilePath(dir string) *FilePath {
	return &FilePath{Dir: dir}
}

func (f *FilePath) GetScripts() (*script_fetcher.ScriptMetadata, error) {
	files, err := os.ReadDir(f.Dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read script directory: %w", err)
	}

	return &script_fetcher.ScriptMetadata{DirFiles: files}, nil
}
