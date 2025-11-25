package db

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

func TestGetDSN_WithExistingDSN(t *testing.T) {
	config := DBConfig{
		Type: MySQL,
		DSN:  "existing-dsn-string",
	}

	assert.Equal(t, "existing-dsn-string", config.GetDSN())
}

func TestGetDSN_MySQL(t *testing.T) {
	tests := []struct {
		name     string
		config   DBConfig
		expected string
	}{
		{
			name: "MySQL with all parameters",
			config: DBConfig{
				Type:     MySQL,
				Host:     "localhost",
				Port:     3306,
				Username: "root",
				Password: "password",
				Database: "testdb",
				Params:   "charset=utf8mb4&parseTime=True&loc=Local",
			},
			expected: "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name: "MySQL without port (default)",
			config: DBConfig{
				Type:     MySQL,
				Host:     "localhost",
				Username: "root",
				Password: "password",
				Database: "testdb",
			},
			expected: "root:password@tcp(localhost:3306)/testdb",
		},
		{
			name: "MySQL without params",
			config: DBConfig{
				Type:     MySQL,
				Host:     "localhost",
				Port:     3307,
				Username: "user",
				Password: "pass",
				Database: "mydb",
			},
			expected: "user:pass@tcp(localhost:3307)/mydb",
		},
		{
			name: "MySQL with empty username and password",
			config: DBConfig{
				Type:     MySQL,
				Host:     "localhost",
				Database: "testdb",
			},
			expected: ":@tcp(localhost:3306)/testdb",
		},
		{
			name: "MySQL without host (empty DSN)",
			config: DBConfig{
				Type:     MySQL,
				Username: "root",
				Password: "password",
				Database: "testdb",
			},
			expected: "",
		},
		{
			name: "MySQL without database (empty DSN)",
			config: DBConfig{
				Type:     MySQL,
				Host:     "localhost",
				Username: "root",
				Password: "password",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.GetDSN())
		})
	}
}

func TestGetDSN_PostgreSQL(t *testing.T) {
	tests := []struct {
		name     string
		config   DBConfig
		expected string
	}{
		{
			name: "PostgreSQL with all parameters",
			config: DBConfig{
				Type:     PostgreSQL,
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "password",
				Database: "testdb",
				Params:   "sslmode=disable&TimeZone=UTC",
			},
			expected: "host=localhost port=5432 user=postgres password=password dbname=testdb sslmode=disable&TimeZone=UTC",
		},
		{
			name: "PostgreSQL without port (default)",
			config: DBConfig{
				Type:     PostgreSQL,
				Host:     "localhost",
				Username: "postgres",
				Password: "password",
				Database: "testdb",
			},
			expected: "host=localhost port=5432 user=postgres password=password dbname=testdb",
		},
		{
			name: "PostgreSQL without params",
			config: DBConfig{
				Type:     PostgreSQL,
				Host:     "localhost",
				Port:     5433,
				Username: "user",
				Password: "pass",
				Database: "mydb",
			},
			expected: "host=localhost port=5433 user=user password=pass dbname=mydb",
		},
		{
			name: "PostgreSQL with special characters in password",
			config: DBConfig{
				Type:     PostgreSQL,
				Host:     "localhost",
				Username: "postgres",
				Password: "p@ss word!",
				Database: "testdb",
			},
			expected: "host=localhost port=5432 user=postgres password=p%40ss+word%21 dbname=testdb",
		},
		{
			name: "PostgreSQL without host (empty DSN)",
			config: DBConfig{
				Type:     PostgreSQL,
				Username: "postgres",
				Password: "password",
				Database: "testdb",
			},
			expected: "",
		},
		{
			name: "PostgreSQL without database (empty DSN)",
			config: DBConfig{
				Type:     PostgreSQL,
				Host:     "localhost",
				Username: "postgres",
				Password: "password",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.GetDSN())
		})
	}
}

func TestGetDSN_UnsupportedType(t *testing.T) {
	config := DBConfig{
		Type:     "unsupported",
		Host:     "localhost",
		Database: "testdb",
	}

	assert.Equal(t, "", config.GetDSN())
}

func TestNewSimpleManagerConfig(t *testing.T) {
	primary := DBConfig{
		Type:     MySQL,
		Host:     "localhost",
		Database: "testdb",
	}
	gormConfig := &gorm.Config{}

	config := NewSimpleManagerConfig(primary, gormConfig)

	assert.Equal(t, MySQL, config.Primary.Type)
	assert.Equal(t, "localhost", config.Primary.Host)
	assert.Equal(t, "testdb", config.Primary.Database)
	assert.Equal(t, gormConfig, config.GormConfig)
	assert.Empty(t, config.Sources)
	assert.Empty(t, config.Replicas)
}

func TestNewSimpleManagerConfig_NilGormConfig(t *testing.T) {
	primary := DBConfig{
		Type:     PostgreSQL,
		Host:     "localhost",
		Database: "testdb",
	}

	config := NewSimpleManagerConfig(primary, nil)

	assert.Equal(t, PostgreSQL, config.Primary.Type)
	assert.Nil(t, config.GormConfig)
	assert.Empty(t, config.Sources)
	assert.Empty(t, config.Replicas)
}

func TestNewReadWriteSplitConfig(t *testing.T) {
	primary := DBConfig{
		Type:     MySQL,
		Host:     "primary.local",
		Database: "maindb",
	}

	sources := []DBConfig{
		{
			Type:     MySQL,
			Host:     "source1.local",
			Database: "maindb",
		},
		{
			Type:     MySQL,
			Host:     "source2.local",
			Database: "maindb",
		},
	}

	replicas := []DBConfig{
		{
			Type:     MySQL,
			Host:     "replica1.local",
			Database: "maindb",
		},
		{
			Type:     MySQL,
			Host:     "replica2.local",
			Database: "maindb",
		},
	}

	gormConfig := &gorm.Config{}

	config := NewReadWriteSplitConfig(primary, sources, replicas, gormConfig)

	assert.Equal(t, "primary.local", config.Primary.Host)
	assert.Equal(t, 2, len(config.Sources))
	assert.Equal(t, "source1.local", config.Sources[0].Host)
	assert.Equal(t, "source2.local", config.Sources[1].Host)
	assert.Equal(t, 2, len(config.Replicas))
	assert.Equal(t, "replica1.local", config.Replicas[0].Host)
	assert.Equal(t, "replica2.local", config.Replicas[1].Host)
	assert.Equal(t, gormConfig, config.GormConfig)
}

func TestNewReadWriteSplitConfig_EmptySources(t *testing.T) {
	primary := DBConfig{
		Type:     PostgreSQL,
		Host:     "primary.local",
		Database: "maindb",
	}

	replicas := []DBConfig{
		{
			Type:     PostgreSQL,
			Host:     "replica1.local",
			Database: "maindb",
		},
	}

	config := NewReadWriteSplitConfig(primary, nil, replicas, nil)

	assert.Equal(t, "primary.local", config.Primary.Host)
	assert.Empty(t, config.Sources)
	assert.Equal(t, 1, len(config.Replicas))
	assert.Nil(t, config.GormConfig)
}
