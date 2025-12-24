//go:build !no_mysql && (all_db || mysql || !no_default_db)

package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// newMySQLDialector creates a MySQL dialector from the given DSN
func newMySQLDialector(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}
