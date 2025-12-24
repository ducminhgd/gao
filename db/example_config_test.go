package db_test

import (
	"fmt"

	"github.com/ducminhgd/gao/db"
)

// Example demonstrating DSN auto-building for MySQL
func ExampleDBConfig_GetDSN_mysql() {
	config := db.DBConfig{
		Type:     db.MySQL,
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "secret",
		Database: "myapp",
		Params:   "charset=utf8mb4&parseTime=True&loc=Local",
	}

	fmt.Println(config.GetDSN())
	// Output: root:secret@tcp(localhost:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local
}

// Example demonstrating DSN auto-building for PostgreSQL
func ExampleDBConfig_GetDSN_postgresql() {
	config := db.DBConfig{
		Type:     db.PostgreSQL,
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "secret",
		Database: "myapp",
		Params:   "sslmode=disable",
	}

	fmt.Println(config.GetDSN())
	// Output: host=localhost port=5432 user=postgres password=secret dbname=myapp sslmode=disable
}

// Example demonstrating that explicitly set DSN takes precedence
func ExampleDBConfig_GetDSN_explicit() {
	config := db.DBConfig{
		Type:     db.MySQL,
		DSN:      "custom-dsn-string",
		Host:     "localhost",
		Database: "myapp",
	}

	fmt.Println(config.GetDSN())
	// Output: custom-dsn-string
}

// Example demonstrating simple single database configuration
func ExampleNewSimpleManagerConfig() {
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
	fmt.Printf("Primary: %s\n", config.Primary.Host)
	fmt.Printf("Sources: %d\n", len(config.Sources))
	fmt.Printf("Replicas: %d\n", len(config.Replicas))
	// Output:
	// Primary: localhost
	// Sources: 0
	// Replicas: 0
}

// Example demonstrating read-write split configuration
func ExampleNewReadWriteSplitConfig() {
	primary := db.DBConfig{
		Type:       db.MySQL,
		Host:       "primary.db.local",
		Database:   "myapp",
		Username:   "root",
		Password:   "secret",
		PoolConfig: db.DefaultPoolConfig(),
	}

	// Additional write databases (optional)
	sources := []db.DBConfig{
		{
			Type:       db.MySQL,
			Host:       "source1.db.local",
			Database:   "myapp",
			Username:   "root",
			Password:   "secret",
			PoolConfig: db.DefaultPoolConfig(),
		},
	}

	// Read-only replicas
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

	config := db.NewReadWriteSplitConfig(primary, sources, replicas, nil)
	fmt.Printf("Primary: %s\n", config.Primary.Host)
	fmt.Printf("Sources: %d\n", len(config.Sources))
	fmt.Printf("Replicas: %d\n", len(config.Replicas))
	// Output:
	// Primary: primary.db.local
	// Sources: 1
	// Replicas: 2
}

// Example demonstrating read-write split with only replicas (no additional sources)
func ExampleNewReadWriteSplitConfig_onlyReplicas() {
	primary := db.DBConfig{
		Type:       db.PostgreSQL,
		Host:       "primary.db.local",
		Database:   "myapp",
		Username:   "postgres",
		Password:   "secret",
		PoolConfig: db.DefaultPoolConfig(),
	}

	// Read-only replicas (primary handles all writes)
	replicas := []db.DBConfig{
		{
			Type:       db.PostgreSQL,
			Host:       "replica1.db.local",
			Database:   "myapp",
			Username:   "readonly",
			Password:   "secret",
			PoolConfig: db.DefaultPoolConfig(),
		},
	}

	config := db.NewReadWriteSplitConfig(primary, nil, replicas, nil)
	fmt.Printf("Primary handles writes: %t\n", len(config.Sources) == 0)
	fmt.Printf("Read replicas: %d\n", len(config.Replicas))
	// Output:
	// Primary handles writes: true
	// Read replicas: 1
}
