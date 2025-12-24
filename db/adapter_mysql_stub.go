//go:build no_mysql || (!all_db && !mysql && no_default_db)

package db

import (
	"fmt"

	"gorm.io/gorm"
)

// newMySQLDialector is a stub that returns an error when MySQL support is disabled
func newMySQLDialector(dsn string) gorm.Dialector {
	panic(fmt.Errorf("MySQL support not compiled in. Build with -tags mysql to enable"))
}
