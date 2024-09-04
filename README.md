# **Migrongo**

**Migrongo** is a Go package designed to handle MongoDB migrations using JavaScript files. It executes each migration in order and keeps track of the versions that have already been applied, similar to tools like `go-migrate`.

## **Features**

- Execute MongoDB migrations using `mongosh`.
- Automatic version tracking to prevent re-running migrations.

## **Installation**

To install the package, you can use `go get`:

```bash
go get github.com/kondukto-io/migrongo
```

## **Usage**

### **1. Create Migration Files**

Create your migration scripts in a designated directory (e.g., `migrations/`). Each file should be named with a version prefix followed by an underscore and a description. Example:

```
migrations/
│
├── 001_initial_setup.js
├── 002_add_indexes.js
└── 003_update_schema.js
```

### **2. Write Your Migrations**

Each migration script should contain valid MongoDB JavaScript commands. For example:

```javascript
// 001_initial_setup.js
db.createCollection("users");

// 002_add_indexes.js
db.users.createIndex({ "email": 1 }, { unique: true });
```

### **3. Integrate Migrongo in Your Go Application**

In your main application, initialize and run the migrator:

```go
package main

import (
	"log"
	"github.com/kondukto-io/migrongo"
)

func main() {
	clientURI := "mongodb://localhost:27017"
	migrationsDir := "./migrations"

	migrator := migrongo.NewMigrator(clientURI, migrationsDir)
	err := migrator.RunMigrations()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}
```

### **4. Run Your Application**

Run your Go application, and Migrongo will automatically execute the migration files in order.

```bash
go run main.go
```

### **5. Migration Versioning**

MongoMigrate automatically tracks applied migrations in a special collection (`migrations`) within your MongoDB. This ensures each migration runs only once.

## **License**

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more details.

## **Contributing**

Contributions are welcome! Please feel free to submit a pull request or open an issue if you encounter any bugs or have suggestions for new features.
