package migrator

import (
	"github.com/kondukto-io/migrongo/internal/migrator/mongosh"
	"github.com/kondukto-io/migrongo/internal/script_fetcher"
)

func NewMigrator(config Config) (Migrator, error) {
	switch config.Type {
	default:
		return mongosh.NewMongoSH(config.MongoSHConfig.MongoURI, config.ScriptFetcher, config.MongoSHConfig.DBName)
	}
}

type (
	Config struct {
		Type          string
		ScriptFetcher script_fetcher.ScriptFetcher
		MongoSHConfig *MongoSHConfig
	}

	MongoSHConfig struct {
		DBName   string
		MongoURI string
	}
)
