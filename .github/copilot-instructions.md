# cloud-sdk Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-10-26

## Active Technologies
- Go 1.21 + Go standard library (`encoding/json`, `context`, `net/http`) (002-fix-network-model)
- Go 1.21+ + Standard library (net/http, encoding/json, context, time) (001-fix-flavors-model)
- N/A (API client SDK) (001-fix-flavors-model)
- Go 1.21+ + Go standard library (encoding/json, context, net/http, time) (006-fix-keypair-model)
- N/A (SDK client library) (006-fix-keypair-model)
- Go 1.21+ + Go standard library (`encoding/json`, `context`, `net/http`, `fmt`) (007-fix-floating-ip-model)
- N/A (SDK client library, no persistence layer) (008-vrm-tag-repository)

- Go 1.21+ + Standard library (`net/http`, `encoding/json`, `context`, `time`); consider code-gen from Swagger (e.g., `oapi-codegen` for models only, justified for type safety and maintainability) (1-vps-project-api)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.21+

## Code Style

Go 1.21+: Follow standard conventions

## Recent Changes
- 008-vrm-tag-repository: Added Go 1.21+
- 007-fix-floating-ip-model: Added Go 1.21+ + Go standard library (`encoding/json`, `context`, `net/http`, `fmt`)
- 006-fix-keypair-model: Added Go 1.21+ + Go standard library (encoding/json, context, net/http, time)


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
