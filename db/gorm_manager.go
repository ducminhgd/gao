package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type GORMManager struct {
	db *gorm.DB
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
