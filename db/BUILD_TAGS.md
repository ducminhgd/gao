# Database Adapter Build Tags

This package supports conditional compilation of database adapters using Go build tags. This allows you to reduce binary size by only including the database drivers you actually need.

## Default Behavior

By default, **all database adapters** (MySQL, PostgreSQL, SQLite) are included in the build.

## Build Tags

### Including Specific Databases Only

To build with **only specific databases**, use the `no_default_db` tag combined with the databases you want:

```bash
# SQLite only
go build -tags "no_default_db,sqlite" ./...

# MySQL only
go build -tags "no_default_db,mysql" ./...

# PostgreSQL only
go build -tags "no_default_db,postgres" ./...

# MySQL and PostgreSQL (no SQLite)
go build -tags "no_default_db,mysql,postgres" ./...
```

### Excluding Specific Databases

To exclude specific databases from the default build:

```bash
# Exclude MySQL (keep PostgreSQL and SQLite)
go build -tags "no_mysql" ./...

# Exclude PostgreSQL (keep MySQL and SQLite)
go build -tags "no_postgres" ./...

# Exclude SQLite (keep MySQL and PostgreSQL)
go build -tags "no_sqlite" ./...

# Exclude multiple databases
go build -tags "no_mysql,no_postgres" ./...
```

### Including All Databases Explicitly

```bash
# Explicitly include all databases
go build -tags "all_db" ./...
```

## Tag Reference

| Tag | Description |
|-----|-------------|
| `no_default_db` | Disable all databases by default (use with specific database tags) |
| `all_db` | Explicitly include all database adapters |
| `mysql` | Include MySQL support |
| `postgres` / `postgresql` | Include PostgreSQL support |
| `sqlite` | Include SQLite support |
| `no_mysql` | Exclude MySQL support |
| `no_postgres` | Exclude PostgreSQL support |
| `no_sqlite` | Exclude SQLite support |

## Examples

### Example 1: SQLite-only Application

```bash
# Build with SQLite only
go build -tags "no_default_db,sqlite" -o myapp ./cmd/myapp

# Or in your code, only use SQLite
config := db.DBConfig{
    Type: db.SQLite,
    Database: "app.db",
}
```

### Example 2: PostgreSQL-only Application

```bash
# Build with PostgreSQL only
go build -tags "no_default_db,postgres" -o myapp ./cmd/myapp
```

### Example 3: MySQL and PostgreSQL (No SQLite)

```bash
# Build without SQLite
go build -tags "no_default_db,mysql,postgres" -o myapp ./cmd/myapp
```

## Binary Size Comparison

Here's an approximate binary size reduction when excluding databases:

- **All databases**: ~15-20MB (baseline)
- **SQLite only**: ~8-10MB (saves ~7-10MB)
- **PostgreSQL only**: ~10-12MB (saves ~5-8MB)
- **MySQL only**: ~10-12MB (saves ~5-8MB)

*Actual sizes depend on your application code and other dependencies.*

## Testing with Build Tags

When running tests with specific build tags:

```bash
# Test with SQLite only
go test -tags "no_default_db,sqlite" ./db/...

# Test with all databases
go test -tags "all_db" ./db/...
```

## Using in go.mod and Projects

If you're building a library or application that always uses specific databases, you can document this in your README and build scripts:

```bash
# In your Makefile or build script
build-sqlite:
	go build -tags "no_default_db,sqlite" -ldflags="-s -w" -o bin/app ./cmd/app

build-postgres:
	go build -tags "no_default_db,postgres" -ldflags="-s -w" -o bin/app ./cmd/app
```

## Error Handling

If you try to use a database adapter that wasn't compiled in, you'll get a panic at runtime:

```
panic: MySQL support not compiled in. Build with -tags mysql to enable
```

Make sure to build with the appropriate tags for the databases your application uses.

## Additional Optimization

For even smaller binaries, combine with linker flags:

```bash
# SQLite only with stripped binary
go build -tags "no_default_db,sqlite" -ldflags="-s -w" -o myapp ./cmd/myapp

# -s: Omit symbol table
# -w: Omit DWARF debug information
```
