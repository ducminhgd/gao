package dbmanager

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// Manager manages database connections with support for sources and replicas
type Manager struct {
	db     *gorm.DB
	config ManagerConfig
}

// NewManager creates a new database manager with the provided configuration
//
// Parameters:
// - config: the ManagerConfig containing database configuration
//
// Returns:
// - *Manager: the newly created Manager instance
// - error: an error if the manager creation fails
func NewManager(config ManagerConfig) (*Manager, error) {
	// Create the primary database connection
	dialector, err := createDialector(config.Primary)
	if err != nil {
		return nil, fmt.Errorf("failed to create primary dialector: %w", err)
	}

	// Use provided GORM config or default
	gormConfig := config.GormConfig
	if gormConfig == nil {
		gormConfig = &gorm.Config{}
	}

	// Open the primary database connection
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open primary database: %w", err)
	}

	// Apply pool configuration to primary database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from primary: %w", err)
	}
	applyPoolConfig(sqlDB, config.Primary.PoolConfig)

	manager := &Manager{
		db:     db,
		config: config,
	}

	// Register sources and replicas if provided
	if len(config.Sources) > 0 || len(config.Replicas) > 0 {
		if err := manager.registerSourcesAndReplicas(); err != nil {
			return nil, fmt.Errorf("failed to register sources/replicas: %w", err)
		}
	}

	return manager, nil
}

// DB returns the underlying *gorm.DB instance
//
// Returns:
// - *gorm.DB: the underlying GORM database instance
func (m *Manager) DB() *gorm.DB {
	return m.db
}

// Close closes all database connections
//
// Returns:
// - error: an error if closing fails
func (m *Manager) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// registerSourcesAndReplicas registers source and replica databases with dbresolver
func (m *Manager) registerSourcesAndReplicas() error {
	var sources []gorm.Dialector
	var replicas []gorm.Dialector

	// Create dialectors for sources
	for _, sourceConfig := range m.config.Sources {
		dialector, err := createDialector(sourceConfig)
		if err != nil {
			return fmt.Errorf("failed to create source dialector: %w", err)
		}
		sources = append(sources, dialector)
	}

	// Create dialectors for replicas
	for _, replicaConfig := range m.config.Replicas {
		dialector, err := createDialector(replicaConfig)
		if err != nil {
			return fmt.Errorf("failed to create replica dialector: %w", err)
		}
		replicas = append(replicas, dialector)
	}

	// Register with dbresolver
	resolverConfig := dbresolver.Config{
		Sources:  sources,
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}

	resolver := dbresolver.Register(resolverConfig)

	// Apply pool configuration for sources
	for _, sourceConfig := range m.config.Sources {
		resolver = resolver.SetConnMaxIdleTime(sourceConfig.PoolConfig.ConnMaxIdleTime).
			SetConnMaxLifetime(sourceConfig.PoolConfig.ConnMaxLifetime).
			SetMaxIdleConns(sourceConfig.PoolConfig.MaxIdleConns).
			SetMaxOpenConns(sourceConfig.PoolConfig.MaxOpenConns)
	}

	// Apply pool configuration for replicas
	for _, replicaConfig := range m.config.Replicas {
		resolver = resolver.SetConnMaxIdleTime(replicaConfig.PoolConfig.ConnMaxIdleTime).
			SetConnMaxLifetime(replicaConfig.PoolConfig.ConnMaxLifetime).
			SetMaxIdleConns(replicaConfig.PoolConfig.MaxIdleConns).
			SetMaxOpenConns(replicaConfig.PoolConfig.MaxOpenConns)
	}

	return m.db.Use(resolver)
}

// createDialector creates a GORM dialector based on the database type and DSN
func createDialector(config DBConfig) (gorm.Dialector, error) {
	switch config.Type {
	case MySQL:
		return mysql.Open(config.DSN), nil
	case PostgreSQL:
		return postgres.Open(config.DSN), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}

// applyPoolConfig applies connection pool configuration to sql.DB
func applyPoolConfig(sqlDB *sql.DB, config PoolConfig) {
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
}
