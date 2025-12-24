package db

import (
	"fmt"
	"net/url"
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
	// If empty, it will be built from Host, Port, Username, Password, Database, and Params
	DSN string
	// Host is the database host (used if DSN is empty)
	Host string
	// Port is the database port (used if DSN is empty)
	Port int
	// Username is the database username (used if DSN is empty)
	Username string
	// Password is the database password (used if DSN is empty)
	Password string
	// Database is the database name (used if DSN is empty)
	Database string
	// Params contains additional connection parameters (used if DSN is empty)
	// For MySQL: e.g., "charset=utf8mb4&parseTime=True&loc=Local"
	// For PostgreSQL: e.g., "sslmode=disable&TimeZone=UTC"
	Params string
	// PoolConfig contains connection pool settings
	PoolConfig PoolConfig
}

// ManagerConfig represents the configuration for the database manager
//
// Two configuration modes are supported:
//
// 1. Simple mode (single database):
//    Set only Primary for a single database connection
//
// 2. Multi-database mode (read-write separation):
//    - Sources: List of write databases (if empty, Primary is used as the sole write DB)
//    - Replicas: List of read-only databases
//    - Primary: The initial connection (always required)
//
// Note: If you specify Sources, the Primary connection becomes the default connection,
// and Sources + Replicas are used by dbresolver for read-write splitting.
// If Sources is empty, Primary handles both reads and writes.
type ManagerConfig struct {
	// Primary is the primary database configuration (required)
	// This is always the initial connection to the database
	Primary DBConfig

	// Sources is a list of write database configurations (optional)
	// When specified, these databases (along with Primary) handle write operations
	// If empty, Primary handles all write operations
	Sources []DBConfig

	// Replicas is a list of read-only database configurations (optional)
	// These databases handle read operations when using dbresolver
	// If empty, reads go to Primary (or Sources if specified)
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

// GetDSN returns the DSN string. If DSN is already set, it returns it directly.
// Otherwise, it builds the DSN from Host, Port, Username, Password, Database, and Params.
func (c *DBConfig) GetDSN() string {
	if c.DSN != "" {
		return c.DSN
	}

	switch c.Type {
	case MySQL:
		return c.buildMySQLDSN()
	case PostgreSQL:
		return c.buildPostgreSQLDSN()
	default:
		return ""
	}
}

// buildMySQLDSN builds a MySQL DSN from individual config parameters
// Format: username:password@tcp(host:port)/database?params
func (c *DBConfig) buildMySQLDSN() string {
	if c.Host == "" || c.Database == "" {
		return ""
	}

	port := c.Port
	if port == 0 {
		port = 3306 // Default MySQL port
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.Username,
		c.Password,
		c.Host,
		port,
		c.Database,
	)

	if c.Params != "" {
		dsn += "?" + c.Params
	}

	return dsn
}

// buildPostgreSQLDSN builds a PostgreSQL DSN from individual config parameters
// Format: host=host port=port user=user password=password dbname=dbname params
func (c *DBConfig) buildPostgreSQLDSN() string {
	if c.Host == "" || c.Database == "" {
		return ""
	}

	port := c.Port
	if port == 0 {
		port = 5432 // Default PostgreSQL port
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		c.Host,
		port,
		c.Username,
		url.QueryEscape(c.Password),
		c.Database,
	)

	if c.Params != "" {
		dsn += " " + c.Params
	}

	return dsn
}

// NewSimpleManagerConfig creates a ManagerConfig for a single database connection
// This is a convenience function for the most common use case.
//
// Parameters:
//   - primary: The primary database configuration
//   - gormConfig: Optional GORM configuration (can be nil)
//
// Returns:
//   - ManagerConfig: A manager config with only the primary database
func NewSimpleManagerConfig(primary DBConfig, gormConfig *gorm.Config) ManagerConfig {
	return ManagerConfig{
		Primary:    primary,
		GormConfig: gormConfig,
	}
}

// NewReadWriteSplitConfig creates a ManagerConfig with read-write separation
// This sets up a configuration where writes go to source databases and reads go to replicas.
//
// Parameters:
//   - primary: The primary database configuration (serves as the initial connection)
//   - sources: Additional write databases (optional, can be empty)
//   - replicas: Read-only databases
//   - gormConfig: Optional GORM configuration (can be nil)
//
// Returns:
//   - ManagerConfig: A manager config with read-write separation
func NewReadWriteSplitConfig(primary DBConfig, sources []DBConfig, replicas []DBConfig, gormConfig *gorm.Config) ManagerConfig {
	return ManagerConfig{
		Primary:    primary,
		Sources:    sources,
		Replicas:   replicas,
		GormConfig: gormConfig,
	}
}
