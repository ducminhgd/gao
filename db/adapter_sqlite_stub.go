//go:build no_sqlite || (!all_db && !sqlite && no_default_db)

package db

import (
	"fmt"

	"gorm.io/gorm"
)

// newSQLiteDialector is a stub that returns an error when SQLite support is disabled
func newSQLiteDialector(dsn string) gorm.Dialector {
	panic(fmt.Errorf("SQLite support not compiled in. Build with -tags sqlite to enable"))
}
