package script_fetcher

import "github.com/kondukto-io/migrongo/internal/script_fetcher/types"

type ScriptFetcher interface {
	GetScripts() (*types.ScriptMetadata, error)
}
