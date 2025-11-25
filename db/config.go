package db

import (
	"time"

	"gorm.io/gorm"
)

// DatabaseType represents the type of database (MySQL or PostgreSQL)
type DatabaseType string

const (
	// MySQL database type
	MySQL DatabaseType = "mysql"
	// PostgreSQL database type
	PostgreSQL DatabaseType = "postgres"
)

// DBConfig represents the configuration for a single database connection
type DBConfig struct {
	// Type is the database type (mysql or postgres)
	Type DatabaseType
	// DSN is the Data Source Name for the database connection
	DSN string
	// PoolConfig contains connection pool settings
	PoolConfig PoolConfig
}

// ManagerConfig represents the configuration for the database manager
type ManagerConfig struct {
	// Primary is the primary database configuration (required)
	Primary DBConfig
	// Sources is a list of source database configurations for read-write operations
	Sources []DBConfig
	// Replicas is a list of replica database configurations for read-only operations
	Replicas []DBConfig
	// GormConfig contains GORM-specific configuration
	GormConfig *gorm.Config
}

// DefaultPoolConfig returns a PoolConfig with default values
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		ConnMaxIdleTime: 10 * time.Minute,
		ConnMaxLifetime: 60 * time.Minute,
		MaxIdleConns:    5,
		MaxOpenConns:    10,
	}
}
