package mongosh

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/kondukto-io/migrongo/internal/script_fetcher"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoSH struct {
	dbName        string
	clientOpts    *options.ClientOptions
	dbClient      *mongo.Client
	scriptFetcher script_fetcher.ScriptFetcher
}

// NewMongoSH creates a new MongoSH instance
func NewMongoSH(mongoClientOptions *options.ClientOptions, scriptFetcher script_fetcher.ScriptFetcher, dbName string) (*MongoSH, error) {
	if dbName == "" {
		return nil, errors.New("db name cannot be empty")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoClientOptions.GetURI()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection
	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &MongoSH{
		clientOpts:    mongoClientOptions,
		scriptFetcher: scriptFetcher,
		dbClient:      client,
		dbName:        dbName,
	}, nil
}

// runScript executes a given JavaScript file using mongosh
func (m *MongoSH) runScript(scriptPath string) error {
	var dbURI = m.clientOpts.GetURI()

	cmd := exec.Command("mongosh", dbURI, "--file", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run script %s: %w", scriptPath, err)
	}

	return nil
}

// AppliedMigrations retrieves applied migration versions from the migrations collection in the specified database
func (m *MongoSH) appliedMigrations() (map[string]bool, error) {
	collection := m.dbClient.Database(m.dbName).Collection("migrations")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch applied migrations: %w", err)
	}
	defer cursor.Close(context.Background())

	applied := make(map[string]bool)
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode applied migration: %w", err)
		}
		version, ok := result["version"].(string)
		if ok {
			applied[version] = true
		}
	}

	// Check for any error during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error encountered while iterating cursor: %w", err)
	}

	// If no migrations have been applied, return an empty map without an error
	return applied, nil
}

// recordMigration records a migration as applied in the database
func (m *MongoSH) recordMigration(version string) error {
	var db = m.dbName
	collection := m.dbClient.Database(db).Collection("migrations")

	_, err := collection.InsertOne(context.Background(), map[string]interface{}{
		"version":   version,
		"appliedAt": time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return nil
}

// removeMigrationRecord removes a migration record from the database after a rollback
func (m *MongoSH) removeMigrationRecord(version string) error {
	var db = m.dbName
	collection := m.dbClient.Database(db).Collection("migrations")

	_, err := collection.DeleteOne(context.Background(), map[string]interface{}{
		"version": version,
	})
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	return nil
}
