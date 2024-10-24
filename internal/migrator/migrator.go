package migrator

import (
	"github.com/kondukto-io/migrongo/internal/migrator/mongosh"
	"github.com/kondukto-io/migrongo/internal/script_fetcher"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMigrator(config Config) (Migrator, error) {
	switch config.Type {
	default:
		return mongosh.NewMongoSH(config.MongoSHConfig.MongoClientOptions, config.ScriptFetcher, config.MongoSHConfig.DBName)
	}
}

type (
	Config struct {
		Type          string
		ScriptFetcher script_fetcher.ScriptFetcher
		MongoSHConfig *MongoSHConfig
	}

	MongoSHConfig struct {
		DBName             string
		MongoClientOptions *options.ClientOptions
	}
)
