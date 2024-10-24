package pkg

import (
	"fmt"

	"github.com/kondukto-io/migrongo/internal/migrator"
	"github.com/kondukto-io/migrongo/internal/script_fetcher"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(config Config) (*Client, error) {
	fetcherConfig := script_fetcher.Config{
		FilePathConfig: &script_fetcher.FilePathConfig{Dir: config.FilePathFetcherConfig.Dir},
	}

	fetcher, err := script_fetcher.NewScriptFetcher(fetcherConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize new script fetcher: %w", err)
	}

	migratorConfig := migrator.Config{
		ScriptFetcher: fetcher,
		MongoSHConfig: &migrator.MongoSHConfig{
			DBName:             config.MongoSHMigratorConfig.DBName,
			MongoClientOptions: config.MongoSHMigratorConfig.MongoClientOptions,
		},
	}

	migratorClient, err := migrator.NewMigrator(migratorConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize new migrator: %w", err)
	}

	return &Client{
		migrator:      migratorClient,
		scriptFetcher: fetcher,
	}, nil
}

func (c Client) Up() error {
	return c.migrator.Up()
}

func (c Client) Down() error {
	return c.migrator.Down()
}

func (c Client) LatestVersion() (string, error) {
	return c.migrator.LatestVersion()
}

type (
	Client struct {
		migrator      migrator.Migrator
		scriptFetcher script_fetcher.ScriptFetcher
	}

	Config struct {
		FilePathFetcherConfig *FilePathFetcherConfig
		MongoSHMigratorConfig *MongoSHConfig
	}

	FilePathFetcherConfig struct {
		Dir string
	}

	MongoSHConfig struct {
		DBName             string
		MongoClientOptions *options.ClientOptions
	}
)
