//go:build !no_postgres && (all_db || postgres || postgresql || !no_default_db)

package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newPostgreSQLDialector creates a PostgreSQL dialector from the given DSN
func newPostgreSQLDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
