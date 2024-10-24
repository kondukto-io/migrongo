package script_fetcher

type ScriptFetcher interface {
	GetScripts() (*ScriptMetadata, error)
}
