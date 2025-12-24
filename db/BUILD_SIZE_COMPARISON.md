# Binary Size Comparison

This document shows the actual binary size reduction achieved by using build tags to exclude unused database drivers.

## Test Setup

Simple application that initializes a SQLite database:

```go
package main

import (
	"fmt"
	"github.com/ducminhgd/gao/db"
)

func main() {
	config := db.DBConfig{
		Type:     db.SQLite,
		Database: ":memory:",
	}

	manager, err := db.NewGORMManagerFromConfig(db.NewSimpleManagerConfig(config, nil))
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	fmt.Println("Database initialized successfully")
}
```

## Build Results

| Build Configuration | Binary Size | Reduction | Command |
|---------------------|-------------|-----------|---------|
| All databases (default) | 18 MB | Baseline | `go build` |
| SQLite only | 9.2 MB | **49% smaller** | `go build -tags "no_default_db,sqlite"` |
| SQLite only (stripped) | 5.4 MB | **70% smaller** | `go build -tags "no_default_db,sqlite" -ldflags="-s -w"` |

## Key Takeaways

1. **Using build tags saves ~9 MB** (49% reduction) when using only SQLite
2. **Adding stripped builds saves an additional ~4 MB** (total 70% reduction)
3. The more databases you exclude, the smaller your binary

## Recommendations

### For Production Builds

If your application uses only one database type:

```bash
# SQLite production build
go build -tags "no_default_db,sqlite" -ldflags="-s -w" -o app

# MySQL production build
go build -tags "no_default_db,mysql" -ldflags="-s -w" -o app

# PostgreSQL production build
go build -tags "no_default_db,postgres" -ldflags="-s -w" -o app
```

### For Development Builds

Keep all databases for flexibility:

```bash
go build -o app
```

Or use specific database for faster builds:

```bash
go build -tags "no_default_db,sqlite" -o app
```

## Additional Size Optimization Tips

1. **Use UPX compression**: Can reduce binary size by an additional 50-70%
   ```bash
   upx --best app-sqlite-stripped
   ```

2. **Enable CGO_ENABLED=0** for pure Go builds (if applicable):
   ```bash
   CGO_ENABLED=0 go build -tags "no_default_db,sqlite" -ldflags="-s -w" -o app
   ```
   Note: SQLite driver uses CGO, so this won't work for SQLite builds

3. **Use module vendoring**: Can slightly reduce build time
   ```bash
   go mod vendor
   go build -mod=vendor -tags "no_default_db,sqlite" -ldflags="-s -w" -o app
   ```

## Impact on Container Images

When deploying in Docker/containers:

```dockerfile
# Multi-stage build example
FROM golang:1.21 as builder

WORKDIR /app
COPY . .

# Build with only required database
RUN go build -tags "no_default_db,sqlite" -ldflags="-s -w" -o app ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .

CMD ["./app"]
```

**Image size comparison:**
- With all databases: ~25 MB (Alpine + 18 MB binary)
- With SQLite only (stripped): ~12 MB (Alpine + 5.4 MB binary)
- **Savings: ~13 MB per image** (52% reduction)

This is especially important for:
- Container registries (reduced storage and bandwidth)
- Kubernetes deployments (faster pod startup)
- CI/CD pipelines (faster image pulls)
