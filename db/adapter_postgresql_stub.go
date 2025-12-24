//go:build no_postgres || (!all_db && !postgres && !postgresql && no_default_db)

package db

import (
	"fmt"

	"gorm.io/gorm"
)

// newPostgreSQLDialector is a stub that returns an error when PostgreSQL support is disabled
func newPostgreSQLDialector(dsn string) gorm.Dialector {
	panic(fmt.Errorf("PostgreSQL support not compiled in. Build with -tags postgres to enable"))
}
