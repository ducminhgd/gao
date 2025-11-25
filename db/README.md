# Database Configuration

This package provides a flexible configuration system for GORM database connections with support for MySQL and PostgreSQL.

## Features

- Auto-build DSN from individual connection parameters
- Support for MySQL and PostgreSQL
- Connection pool configuration
- Read-write separation with dbresolver
- Multiple database sources and replicas

## Configuration Modes

### 1. Simple Mode (Single Database)

Use `NewSimpleManagerConfig` for a single database connection:

```go
primary := db.DBConfig{
    Type:       db.MySQL,
    Host:       "localhost",
    Port:       3306,
    Username:   "root",
    Password:   "secret",
    Database:   "myapp",
    Params:     "charset=utf8mb4&parseTime=True",
    PoolConfig: db.DefaultPoolConfig(),
}

config := db.NewSimpleManagerConfig(primary, nil)
manager, err := db.NewGORMManagerFromConfig(config)
```

### 2. Read-Write Split Mode

Use `NewReadWriteSplitConfig` for read-write separation:

```go
primary := db.DBConfig{
    Type:     db.MySQL,
    Host:     "primary.db.local",
    Database: "myapp",
    Username: "root",
    Password: "secret",
}

// Optional: Additional write databases
sources := []db.DBConfig{
    {
        Type:     db.MySQL,
        Host:     "source1.db.local",
        Database: "myapp",
        Username: "root",
        Password: "secret",
    },
}

// Read-only replicas
replicas := []db.DBConfig{
    {
        Type:     db.MySQL,
        Host:     "replica1.db.local",
        Database: "myapp",
        Username: "readonly",
        Password: "secret",
    },
}

config := db.NewReadWriteSplitConfig(primary, sources, replicas, nil)
manager, err := db.NewGORMManagerFromConfig(config)
```

## DSN Configuration

You have two options for specifying database connections:

### Option 1: Explicit DSN

Set the DSN string directly:

```go
config := db.DBConfig{
    Type: db.MySQL,
    DSN:  "root:password@tcp(localhost:3306)/myapp?charset=utf8mb4",
}
```

### Option 2: Auto-build from Parameters

Specify individual connection parameters, and the DSN will be built automatically:

```go
config := db.DBConfig{
    Type:     db.MySQL,
    Host:     "localhost",
    Port:     3306,  // Optional: defaults to 3306 for MySQL, 5432 for PostgreSQL
    Username: "root",
    Password: "password",
    Database: "myapp",
    Params:   "charset=utf8mb4&parseTime=True&loc=Local",  // Optional
}

// DSN is automatically built when GetDSN() is called
dsn := config.GetDSN()
// Result: root:password@tcp(localhost:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local
```

### MySQL DSN Format

When auto-building, MySQL DSN follows this format:
```
username:password@tcp(host:port)/database?params
```

### PostgreSQL DSN Format

When auto-building, PostgreSQL DSN follows this format:
```
host=host port=port user=username password=password dbname=database params
```

Note: PostgreSQL passwords with special characters are automatically URL-encoded.

## Understanding Primary, Sources, and Replicas

The configuration structure supports three types of database connections:

- **Primary**: The initial database connection (required)
  - Always the first connection established
  - Handles all operations when Sources and Replicas are empty

- **Sources**: Additional write databases (optional)
  - Used for write operations (INSERT, UPDATE, DELETE)
  - When specified, writes are distributed among Sources using the configured policy
  - If empty, Primary handles all writes

- **Replicas**: Read-only databases (optional)
  - Used for read operations (SELECT)
  - Reads are distributed among Replicas using the configured policy
  - If empty, reads go to Primary (or Sources if specified)

### Common Scenarios

**Scenario 1: Single Database**
```go
// Primary only - handles all reads and writes
config := db.NewSimpleManagerConfig(primary, nil)
```

**Scenario 2: Read Replicas Only**
```go
// Primary handles writes, replicas handle reads
config := db.NewReadWriteSplitConfig(primary, nil, replicas, nil)
```

**Scenario 3: Multiple Write Sources and Read Replicas**
```go
// Primary + sources handle writes, replicas handle reads
config := db.NewReadWriteSplitConfig(primary, sources, replicas, nil)
```

## Connection Pool Configuration

Configure connection pooling for each database:

```go
poolConfig := db.PoolConfig{
    ConnMaxIdleTime: 10 * time.Minute,
    ConnMaxLifetime: 60 * time.Minute,
    MaxIdleConns:    5,
    MaxOpenConns:    10,
}

config := db.DBConfig{
    // ... other fields
    PoolConfig: poolConfig,
}
```

Or use the defaults:

```go
config := db.DBConfig{
    // ... other fields
    PoolConfig: db.DefaultPoolConfig(),
}
```

## Complete Example

```go
package main

import (
    "log"
    "github.com/ducminhgd/gao/db"
)

func main() {
    // Configure primary database
    primary := db.DBConfig{
        Type:       db.MySQL,
        Host:       "primary.db.local",
        Port:       3306,
        Username:   "root",
        Password:   "secret",
        Database:   "myapp",
        Params:     "charset=utf8mb4&parseTime=True&loc=Local",
        PoolConfig: db.DefaultPoolConfig(),
    }

    // Configure read replicas
    replicas := []db.DBConfig{
        {
            Type:       db.MySQL,
            Host:       "replica1.db.local",
            Database:   "myapp",
            Username:   "readonly",
            Password:   "secret",
            PoolConfig: db.DefaultPoolConfig(),
        },
        {
            Type:       db.MySQL,
            Host:       "replica2.db.local",
            Database:   "myapp",
            Username:   "readonly",
            Password:   "secret",
            PoolConfig: db.DefaultPoolConfig(),
        },
    }

    // Create configuration
    config := db.NewReadWriteSplitConfig(primary, nil, replicas, nil)

    // Create manager
    manager, err := db.NewGORMManagerFromConfig(config)
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Close()

    // Use the database
    db := manager.DB()
    // ... perform database operations
}
```

## Notes

- DSN strings take precedence over individual parameters
- Default ports: MySQL (3306), PostgreSQL (5432)
- Special characters in PostgreSQL passwords are automatically URL-encoded
- The `GetDSN()` method returns an empty string if required parameters (Host, Database) are missing
