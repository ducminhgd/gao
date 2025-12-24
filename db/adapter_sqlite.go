//go:build !no_sqlite && (all_db || sqlite || !no_default_db)

package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newSQLiteDialector creates a SQLite dialector from the given DSN
func newSQLiteDialector(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
