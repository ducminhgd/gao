package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type gormManager struct {
	db *gorm.DB
}

type GORMManager interface {
	AddSourceDialector(gorm.Dialector, PoolConfig) (*gormManager, error)
	AddReplicaDialector(gorm.Dialector, PoolConfig) (*gormManager, error)
	WithLogger(logger.Interface) GORMManager
}

type PoolConfig struct {
	ConnMaxIdleTime time.Duration `default:"10m"`
	ConnMaxLifetime time.Duration `default:"60m"`
	MaxIdleConns    int           `default:"5"`
	MaxOpenConns    int           `default:"10"`
}

// NewGORMManager creates a new gormManager instance.
//
// Parameters:
// - d: the gorm.Dialector to be used for the gormManager.
// - opts: a variadic parameter of gorm.Option to be applied to the gormManager.
//
// Returns:
// - *gormManager: the newly created gormManager instance.
// - error: an error if the gormManager creation fails.
func NewGORMManager(d gorm.Dialector, opts ...gorm.Option) (*gormManager, error) {
	db, err := gorm.Open(d, opts...)
	if err != nil {
		return nil, err
	}
	return &gormManager{
		db: db,
	}, nil
}

// WithLogger sets the logger for the gormManager instance.
//
// Parameters:
// - lgr: the logger.Interface to be set for the gormManager.
//
// Returns:
// - *gormManager: the updated gormManager instance with the new logger.
func (m *gormManager) WithLogger(lgr logger.Interface) *gormManager {
	m.db.Logger = lgr
	return m
}

// AddSourceDialector adds a source dialector to the gormManager.
//
// Parameters:
// - d: the gorm.Dialector to be added as a source.
// - pc: the PoolConfig containing connection settings.
//
// Returns:
// - *gormManager: the updated gormManager instance.
// - error: an error if the dialector registration fails.
func (m *gormManager) AddSourceDialector(d gorm.Dialector, pc PoolConfig) (*gormManager, error) {
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

// AddReplicaDialector adds a replica dialector to the gormManager.
//
// Parameters:
// - d: the gorm.Dialector to be added as a replica.
// - pc: the PoolConfig containing connection settings.
//
// Returns:
// - *gormManager: the updated gormManager instance.
// - error: an error if the dialector registration fails.
func (m *gormManager) AddReplicaDialector(d gorm.Dialector, pc PoolConfig) (*gormManager, error) {
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
