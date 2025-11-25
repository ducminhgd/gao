package dbmanager

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDefaultPoolConfig(t *testing.T) {
	config := DefaultPoolConfig()

	assert.Equal(t, 10*time.Minute, config.ConnMaxIdleTime)
	assert.Equal(t, 60*time.Minute, config.ConnMaxLifetime)
	assert.Equal(t, 5, config.MaxIdleConns)
	assert.Equal(t, 10, config.MaxOpenConns)
}

func TestCreateDialector(t *testing.T) {
	tests := []struct {
		name      string
		config    DBConfig
		expectErr bool
	}{
		{
			name: "MySQL dialector",
			config: DBConfig{
				Type: MySQL,
				DSN:  "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
			},
			expectErr: false,
		},
		{
			name: "PostgreSQL dialector",
			config: DBConfig{
				Type: PostgreSQL,
				DSN:  "host=localhost user=postgres password=postgres dbname=testdb port=5432 sslmode=disable",
			},
			expectErr: false,
		},
		{
			name: "Invalid database type",
			config: DBConfig{
				Type: "invalid",
				DSN:  "some-dsn",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialector, err := createDialector(tt.config)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, dialector)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dialector)
			}
		})
	}
}

func TestManagerConfig(t *testing.T) {
	config := ManagerConfig{
		Primary: DBConfig{
			Type:       MySQL,
			DSN:        "user:password@tcp(localhost:3306)/primary?charset=utf8mb4&parseTime=True&loc=Local",
			PoolConfig: DefaultPoolConfig(),
		},
		Sources: []DBConfig{
			{
				Type:       MySQL,
				DSN:        "user:password@tcp(localhost:3307)/source1?charset=utf8mb4&parseTime=True&loc=Local",
				PoolConfig: DefaultPoolConfig(),
			},
		},
		Replicas: []DBConfig{
			{
				Type:       PostgreSQL,
				DSN:        "host=localhost user=postgres password=postgres dbname=replica1 port=5433 sslmode=disable",
				PoolConfig: DefaultPoolConfig(),
			},
		},
		GormConfig: &gorm.Config{},
	}

	assert.Equal(t, MySQL, config.Primary.Type)
	assert.Equal(t, 1, len(config.Sources))
	assert.Equal(t, 1, len(config.Replicas))
	assert.Equal(t, PostgreSQL, config.Replicas[0].Type)
}

func TestDatabaseTypes(t *testing.T) {
	assert.Equal(t, DatabaseType("mysql"), MySQL)
	assert.Equal(t, DatabaseType("postgres"), PostgreSQL)
}
