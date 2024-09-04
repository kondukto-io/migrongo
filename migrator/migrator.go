package migrator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Migrator struct {
	ScriptDir          string
	DBName             string
	MongoClientOptions *options.ClientOptions
	dbClient           *mongo.Client
}

// NewMigrator creates a new Migrator instance
func NewMigrator(mongoClientOptions *options.ClientOptions, dbName, scriptDir string) (*Migrator, error) {
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

	return &Migrator{
		ScriptDir:          scriptDir,
		MongoClientOptions: mongoClientOptions,
		dbClient:           client,
	}, nil
}

// runScript executes a given JavaScript file using mongosh
func (m *Migrator) runScript(scriptPath string) error {
	var dbURI = m.MongoClientOptions.GetURI()

	cmd := exec.Command("mongosh", dbURI, "--file", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run script %s: %w", scriptPath, err)
	}

	return nil
}

// appliedMigrations retrieves applied migration versions from the database
func (m *Migrator) appliedMigrations() (map[string]bool, error) {
	collection := m.dbClient.Database(m.DBName).Collection("migrations")

	cursor, err := collection.Find(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch applied migrations: %w", err)
	}
	defer cursor.Close(context.Background())

	applied := make(map[string]bool)
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode applied migration: %w", err)
		}
		version, ok := result["version"].(string)
		if ok {
			applied[version] = true
		}
	}

	return applied, nil
}

// recordMigration records a migration as applied in the database
func (m *Migrator) recordMigration(version string) error {
	var db = m.DBName
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
func (m *Migrator) removeMigrationRecord(version string) error {
	collection := m.dbClient.Database("migrations").Collection("applied")

	_, err := collection.DeleteOne(context.Background(), map[string]interface{}{
		"version": version,
	})
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	return nil
}
