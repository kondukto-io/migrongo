# Migrongo

`Migrongo` is a Go package designed to handle MongoDB migrations using JavaScript files and the `mongosh` shell.

## Overview

Migrongo allows you to manage your MongoDB schema evolution by running migration scripts written in JavaScript. The package provides an easy way to run "up" and "down" migrations to apply and rollback changes.

## Installation

First, install the package using `go get`:

```bash
go get github.com/kondukto-io/migrongo
```

## Usage

To use Migrongo, import the package in your Go project and call the `Up` and `Down` functions as needed.

```go
package main

import (
	"fmt"
	"log"

	"github.com/kondukto-io/migrongo/pkg"
	"your_project/infra" // Update with the correct import path for your infra package
)

func main() {
    // Get MongoDB connection options
    opts := infra.GetMongoDBOpts()
  
    // Directory containing the migration scripts
    scriptsDir := "/scripts-path"
  
    config := migrongo.Config{
      FilePathFetcherConfig: &migrongo.FilePathFetcherConfig{Dir: scriptsDir},
      MongoSHMigratorConfig: &migrongo.MongoSHConfig{
        DBName:             infra.DBName,
        MongoClientOptions: opts,
      },
    }
  
    migron, err := migrongo.NewClient(config)
    if err != nil {
      return fmt.Errorf("failed to init migrongo client :%w", err)
    }
  
    // Retrieve the latest migration version
    version, err := migron.LatestVersion()
    if err != nil {
      return fmt.Errorf("error getting the latest version: %w", err)
    }
  
    logger.Log.Infof("current latest version: %s\n", version)
  
    // Run migrations in the 'up' direction
    if err := migron.Up(); err != nil {
      return fmt.Errorf("error running migrations up: %w", err)
    }
  
    // Run migrations in the 'down' direction
    if err := migron.Down(); err != nil {
      return fmt.Errorf("error running migrations down: %w", err)
    }
}
```

## Writing Migrations

Migration scripts should be placed in a directory (e.g., `./scripts`) and should follow a specific naming convention to ensure proper ordering and version control.

### Naming Convention

- Migrations should be named with a version number, direction (`up` or `down`), and a description. For example:
    - `001_up_create_users.js`
    - `001_down_create_users.js`
    - `002_up_add_email_to_users.js`
    - `002_down_add_email_to_users.js`

### Script Content

Each migration script should contain valid JavaScript code that can be executed in the MongoDB shell. For example:

```javascript
// 001_up_create_users.js
db.createCollection("users");

// 002_up_add_email_to_users.js
db.users.updateMany({}, { $set: { email: "" } });
```

## Version Tracking

Migrongo uses a `migrations` collection in your MongoDB database to keep track of which migrations have been applied. This prevents migrations from being run multiple times.

- **Applying Migrations**: When you run `Up`, the applied migrations are recorded in the `migrations` collection.
- **Rolling Back Migrations**: When you run `Down`, the corresponding migration records are removed from the `migrations` collection to keep the state consistent.

## Error Handling

If a migration fails, Migrongo stops the execution and returns an error. It’s recommended to handle these errors in your application logic to ensure consistent state management.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

## License

Migrongo is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contact

For any questions or issues, please open an issue on the GitHub repository.
