# Cloud SDK

A Go SDK for interacting with Cloud VPS APIs.

## Features

- **Type-safe API**: Idiomatic Go interfaces for all VPS operations
- **Project-scoped**: Bind operations to a project once, no repetition
- **Automatic retry**: Exponential backoff with jitter for transient failures (429, 502, 503, 504)
- **Structured errors**: Rich error information with HTTP status codes and metadata
- **Context support**: Timeout control and cancellation via context
- **Tested**: Comprehensive unit and integration tests

## Installation

```bash
go get github.com/Zillaforge/cloud-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk"
)

func main() {
    // Create client with base URL and bearer token
    client, err := cloudsdk.New("https://api.example.com", "your-bearer-token")
    if err != nil {
        log.Fatal(err)
    }
    
    // Get project-scoped VPS client
    vpsClient := client.Project("your-project-id").VPS()
    
    // Use the VPS client for operations
    // Example: vpsClient.Servers().List(ctx, opts)
    
    fmt.Println("SDK initialized successfully")
}
```

## Development

### Prerequisites

- Go 1.21 or later
- Optional: golangci-lint (for linting)
- Optional: goimports (for import formatting)

Run `make install-tools` to install optional tools.

### Quick Commands

```bash
# Show all available commands
make help

# Format code
make fmt

# Run linters
make lint

# Run tests
make test

# Build
make build

# Run all checks (format + lint + test)
make check
```

### Building

```bash
go build ./...
```

### Testing

```bash
# Run all tests
make test

# Run tests with race detector (requires CGO)
make test-race

# Run tests with coverage report
make coverage

# Or use go directly
go test ./...
```

### Available Make Targets

| Target | Description |
|--------|-------------|
| `help` | Show all available commands (default) |
| `fmt` | Format code using gofmt and goimports |
| `lint` | Run golangci-lint or go vet |
| `test` | Run tests with coverage |
| `test-race` | Run tests with race detector |
| `coverage` | Generate HTML coverage report |
| `build` | Build all packages |
| `clean` | Remove build artifacts |
| `deps` | Download and tidy dependencies |
| `install-tools` | Install dev tools (goimports, golangci-lint) |
| `check` | Run fmt + lint + test |

**Note**: The specs directory is excluded from build/test/lint operations as it contains documentation files.

### Linting

```bash
# Run linter
golangci-lint run ./...

# Or use make
make lint
```

### Formatting

```bash
# Format code
make fmt
```

## Project Structure

```
.
├── client.go              # Top-level SDK client
├── errors.go              # Error types and constructors
├── internal/              # Internal packages (not exported)
│   ├── backoff/           # Retry backoff logic
│   ├── http/              # HTTP client wrapper
│   └── types/             # Shared internal types
├── modules/               # Service modules
│   └── vps/               # VPS service client
│       ├── client.go
│       ├── models/        # Request/response types
│       ├── tests/         # Service-specific tests
│       └── waiters/       # Async operation helpers
└── Makefile               # Build and test commands
```

## Architecture

- **Layered design**: Top-level client → Project selector → Service client
- **Internal utilities**: Shared HTTP handling, retry logic, and error mapping
- **TDD approach**: Tests written first, implementation follows
- **Constitution compliance**: Follows Cloud SDK design principles

## License

TBD

## Contributing

TBD
