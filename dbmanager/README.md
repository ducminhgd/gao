# Database Manager

A flexible database manager package for Go that supports multiple database connections with GORM, including MySQL and PostgreSQL with support for source and replica configurations.

## Features

- Support for MySQL and PostgreSQL databases
- Multiple source databases for read-write operations
- Multiple replica databases for read-only operations
- Connection pool configuration
- Built on top of GORM and dbresolver plugin
- Simple and intuitive API

## Installation

```bash
go get github.com/ducminhgd/gao/dbmanager
```

## Usage

### Basic Usage with Single Database

```go
package main

import (
    "log"
    "github.com/ducminhgd/gao/dbmanager"
)

func main() {
    config := dbmanager.ManagerConfig{
        Primary: dbmanager.DBConfig{
            Type: dbmanager.MySQL,
            DSN:  "user:password@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local",
            PoolConfig: dbmanager.DefaultPoolConfig(),
        },
    }

    manager, err := dbmanager.NewManager(config)
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Close()

    // Use the database
    db := manager.DB()
    // ... perform database operations
}
```

### Usage with Sources and Replicas

```go
package main

import (
    "log"
    "time"
    "github.com/ducminhgd/gao/dbmanager"
)

func main() {
    config := dbmanager.ManagerConfig{
        Primary: dbmanager.DBConfig{
            Type: dbmanager.MySQL,
            DSN:  "user:password@tcp(primary.example.com:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local",
            PoolConfig: dbmanager.PoolConfig{
                ConnMaxIdleTime: 10 * time.Minute,
                ConnMaxLifetime: 60 * time.Minute,
                MaxIdleConns:    10,
                MaxOpenConns:    100,
            },
        },
        Sources: []dbmanager.DBConfig{
            {
                Type: dbmanager.MySQL,
                DSN:  "user:password@tcp(source1.example.com:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local",
                PoolConfig: dbmanager.DefaultPoolConfig(),
            },
            {
                Type: dbmanager.MySQL,
                DSN:  "user:password@tcp(source2.example.com:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local",
                PoolConfig: dbmanager.DefaultPoolConfig(),
            },
        },
        Replicas: []dbmanager.DBConfig{
            {
                Type: dbmanager.PostgreSQL,
                DSN:  "host=replica1.example.com user=postgres password=password dbname=mydb port=5432 sslmode=disable",
                PoolConfig: dbmanager.DefaultPoolConfig(),
            },
            {
                Type: dbmanager.PostgreSQL,
                DSN:  "host=replica2.example.com user=postgres password=password dbname=mydb port=5432 sslmode=disable",
                PoolConfig: dbmanager.DefaultPoolConfig(),
            },
        },
    }

    manager, err := dbmanager.NewManager(config)
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Close()

    // Use the database
    db := manager.DB()

    // Write operations go to sources
    db.Create(&User{Name: "John"})

    // Read operations can use replicas
    var users []User
    db.Find(&users)
}
```

### Custom GORM Configuration

```go
config := dbmanager.ManagerConfig{
    Primary: dbmanager.DBConfig{
        Type: dbmanager.MySQL,
        DSN:  "user:password@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local",
        PoolConfig: dbmanager.DefaultPoolConfig(),
    },
    GormConfig: &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    },
}

manager, err := dbmanager.NewManager(config)
```

## Configuration

### DatabaseType

Supported database types:
- `dbmanager.MySQL` - MySQL database
- `dbmanager.PostgreSQL` - PostgreSQL database

### DBConfig

- `Type`: Database type (MySQL or PostgreSQL)
- `DSN`: Data Source Name connection string
- `PoolConfig`: Connection pool configuration

### PoolConfig

- `ConnMaxIdleTime`: Maximum amount of time a connection may be idle (default: 10 minutes)
- `ConnMaxLifetime`: Maximum amount of time a connection may be reused (default: 60 minutes)
- `MaxIdleConns`: Maximum number of idle connections in the pool (default: 5)
- `MaxOpenConns`: Maximum number of open connections to the database (default: 10)

### ManagerConfig

- `Primary`: Primary database configuration (required)
- `Sources`: List of source databases for read-write operations
- `Replicas`: List of replica databases for read-only operations
- `GormConfig`: GORM-specific configuration (optional)

## Connection String Examples

### MySQL DSN Format

```
user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

### PostgreSQL DSN Format

```
host=localhost user=postgres password=password dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Shanghai
```

## License

See the main repository LICENSE file.
