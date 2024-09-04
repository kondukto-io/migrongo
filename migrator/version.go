package migrator

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LatestVersion retrieves the latest applied migration version from the database
func (m *Migrator) LatestVersion() (string, error) {
	collection := m.dbClient.Database(m.DBName).Collection("migrations")

	// Find the document with the highest version
	opts := options.FindOne().SetSort(bson.D{{"version", -1}})
	var result bson.M
	err := collection.FindOne(context.Background(), bson.M{}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No migrations have been applied yet
			return "", nil
		}
		return "", fmt.Errorf("failed to fetch latest version: %w", err)
	}

	version, ok := result["version"].(string)
	if !ok {
		return "", fmt.Errorf("version format is not correct in the database")
	}

	return version, nil
}

// extractVersionFromFileName extracts the version from the migration script filename
func extractVersionFromFileName(fileName string) string {
	parts := strings.Split(fileName, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
