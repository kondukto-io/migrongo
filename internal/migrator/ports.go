package migrator

type Migrator interface {
	Up() error
	Down() error
	LatestVersion() (string, error)
}
