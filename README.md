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
	"log"
	"github.com/kondukto-io/migrongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB URI and script directory
	clientOptions := options.Client()
	scriptDir := "./scripts"

	// Run Up migrations
	if err := migrongo.Up(clientOptions, scriptDir); err != nil {
		log.Fatalf("Failed to run up migrations: %v", err)
	}

	// Run Down migrations
	if err := migrongo.Down(clientOptions, scriptDir); err != nil {
		log.Fatalf("Failed to run down migrations: %v", err)
	}
}
```

### Example

1. **Up Migrations**: To apply all pending "up" migrations:

    ```go
    err := migrongo.Up(opts, "./scripts")
    if err != nil {
        log.Fatalf("Failed to apply up migrations: %v", err)
    }
    ```

2. **Down Migrations**: To rollback the most recent "down" migrations:

    ```go
    err := migrongo.Down(opts, "./scripts")
    if err != nil {
        log.Fatalf("Failed to rollback down migrations: %v", err)
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
