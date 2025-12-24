package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type GORMManager struct {
	db     *gorm.DB
	config ManagerConfig
}

type PoolConfig struct {
	ConnMaxIdleTime time.Duration `default:"10m"`
	ConnMaxLifetime time.Duration `default:"60m"`
	MaxIdleConns    int           `default:"5"`
	MaxOpenConns    int           `default:"10"`
}

// NewGORMManager creates a new GORMManager instance.
//
// Parameters:
// - d: the gorm.Dialector to be used for the GORMManager.
// - opts: a variadic parameter of gorm.Option to be applied to the GORMManager.
//
// Returns:
// - *GORMManager: the newly created GORMManager instance.
// - error: an error if the GORMManager creation fails.
func NewGORMManager(d gorm.Dialector, opts ...gorm.Option) (*GORMManager, error) {
	db, err := gorm.Open(d, opts...)
	if err != nil {
		return nil, err
	}
	return &GORMManager{
		db: db,
	}, nil
}

// WithLogger sets the logger for the GORMManager instance.
//
// Parameters:
// - lgr: the logger.Interface to be set for the GORMManager.
//
// Returns:
// - *GORMManager: the updated GORMManager instance with the new logger.
func (m *GORMManager) WithLogger(lgr logger.Interface) *GORMManager {
	m.db.Logger = lgr
	return m
}

// AddSourceDialector adds a source dialector to the GORMManager.
//
// Parameters:
// - d: the gorm.Dialector to be added as a source.
// - pc: the PoolConfig containing connection settings.
//
// Returns:
// - *GORMManager: the updated GORMManager instance.
// - error: an error if the dialector registration fails.
func (m *GORMManager) AddSourceDialector(d gorm.Dialector, pc PoolConfig) (*GORMManager, error) {
	err := m.db.Use(dbresolver.Register(
		dbresolver.Config{
			Sources:           []gorm.Dialector{d},
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}).
		SetConnMaxIdleTime(pc.ConnMaxIdleTime).
		SetConnMaxLifetime(pc.ConnMaxLifetime).
		SetMaxOpenConns(pc.MaxOpenConns).
		SetMaxIdleConns(pc.MaxIdleConns),
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// AddReplicaDialector adds a replica dialector to the GORMManager.
//
// Parameters:
// - d: the gorm.Dialector to be added as a replica.
// - pc: the PoolConfig containing connection settings.
//
// Returns:
// - *GORMManager: the updated GORMManager instance.
// - error: an error if the dialector registration fails.
func (m *GORMManager) AddReplicaDialector(d gorm.Dialector, pc PoolConfig) (*GORMManager, error) {
	err := m.db.Use(dbresolver.Register(
		dbresolver.Config{
			Replicas:          []gorm.Dialector{d},
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}).
		SetConnMaxIdleTime(pc.ConnMaxIdleTime).
		SetConnMaxLifetime(pc.ConnMaxLifetime).
		SetMaxOpenConns(pc.MaxOpenConns).
		SetMaxIdleConns(pc.MaxIdleConns),
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// DB returns the underlying *gorm.DB instance of the GORMManager.
//
// Parameters:
// - None
//
// Returns:
// - *gorm.DB: the underlying *gorm.DB instance.
func (m *GORMManager) DB() *gorm.DB {
	return m.db
}

// Close closes all database connections
//
// Returns:
// - error: an error if closing fails
func (m *GORMManager) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// NewGORMManagerFromConfig creates a new GORMManager with the provided configuration.
// This function supports MySQL and PostgreSQL databases with multiple sources and replicas.
//
// Parameters:
// - config: the ManagerConfig containing database configuration
//
// Returns:
// - *GORMManager: the newly created GORMManager instance
// - error: an error if the manager creation fails
func NewGORMManagerFromConfig(config ManagerConfig) (*GORMManager, error) {
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

	manager := &GORMManager{
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

// registerSourcesAndReplicas registers source and replica databases with dbresolver
func (m *GORMManager) registerSourcesAndReplicas() error {
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
	dsn := config.GetDSN()
	if dsn == "" {
		return nil, fmt.Errorf("DSN is empty for database type: %s", config.Type)
	}

	switch config.Type {
	case MySQL:
		return mysql.Open(dsn), nil
	case PostgreSQL:
		return postgres.Open(dsn), nil
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
