# Research: VPS Project APIs SDK

**Phase 0 Output** | **Date**: 2025-10-26  
**Purpose**: Resolve all NEEDS CLARIFICATION from Technical Context; document technology choices and patterns

## Research Tasks

### 1. Go Module Structure for Multi-Service SDK

**Decision**: Mono-repo with service-specific packages under root module  
**Rationale**:
- Single `go.mod` at root (`github.com/Zillaforge/cloud-sdk`) simplifies dependency management and versioning
- Each service (VPS, future compute/storage) gets a dedicated package (e.g., `vps/`) with independent `.go` files
- Consumers import: `import "github.com/Zillaforge/cloud-sdk/vps"`
- Easier to maintain shared code (`internal/http`, `errors.go`) and cross-service consistency
- Aligns with Go best practices for SDK structure (see AWS SDK for Go v2, Google Cloud Go SDK)

**Alternatives considered**:
- **Separate modules per service** (e.g., `cloud-sdk-vps`, `cloud-sdk-compute`): More complex versioning; overkill for initial scope; harder to share internal helpers
- **Single flat package**: Naming conflicts inevitable; poor organization for 50+ API operations

**References**:
- [AWS SDK for Go v2 structure](https://github.com/aws/aws-sdk-go-v2) (service/* pattern)
- [Effective Go: Package names](https://go.dev/doc/effective_go#names)

---

### 2. Code Generation from Swagger vs Hand-Written Models

**Decision**: Hand-write request/response structs with Swagger as authoritative reference; defer code-gen unless maintenance burden exceeds threshold  
**Rationale**:
- **Pros of code-gen (oapi-codegen)**:
  - Type safety guaranteed
  - Auto-sync with Swagger updates
  - Reduces boilerplate for ~50 structs
- **Cons of code-gen**:
  - Adds build-time dependency
  - Generated code may include unnecessary fields (admin endpoints)
  - Learning curve for contributors
  - Needs custom templates to match idiomatic naming (e.g., `ID` not `Id`)
- **Hand-written approach**:
  - Full control over struct tags, naming, validation
  - Zero runtime dependencies
  - Easier to test and document
  - Maintenance cost: ~2-3 hours to write initial models; incremental updates when Swagger changes
- **Decision**: Start hand-written; re-evaluate if Swagger changes frequently or struct count exceeds 100

**Alternatives considered**:
- **go-swagger**: Heavier tooling; generates entire server + client (overkill for client-only SDK)
- **OpenAPI Generator**: Java-based; less idiomatic Go output

**Action**: Create `vps/models.go` with hand-crafted structs; add comment linking to Swagger line numbers for traceability

---

### 3. Error Wrapping Pattern for Structured Errors

**Decision**: Custom `SDKError` type implementing `error` interface; wrap with context using `fmt.Errorf` + `%w`  
**Rationale**:
- Spec requires: `statusCode`, `errorCode`, `message`, `meta` fields
- Go 1.13+ error wrapping allows `errors.Is/As` for type assertions
- Pattern:
  ```go
  type SDKError struct {
      StatusCode int
      ErrorCode  int
      Message    string
      Meta       map[string]interface{}
      Cause      error // underlying error (e.g., network timeout)
  }
  func (e *SDKError) Error() string { return e.Message }
  func (e *SDKError) Unwrap() error { return e.Cause }
  ```
- Client-side failures (timeout, canceled, DNS): `StatusCode=0`, `Message="context deadline exceeded"`, `Meta={"category":"timeout", "cause":"..."}`
- Server errors: parse `vpserr.ErrResponse` from JSON body; populate `StatusCode` from HTTP, `ErrorCode`/`Message`/`Meta` from body

**Alternatives considered**:
- **Plain `error` with custom messages**: Loses structured fields; hard to programmatically handle specific errors
- **Multiple error types** (NetworkError, TimeoutError, APIError): Over-engineering; caller needs type switches

**References**:
- [Go blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- Spec FR-011, Error Handling Contract section

---

### 4. Retry Logic: Exponential Backoff + Jitter Implementation

**Decision**: Implement in `internal/backoff` package; use `time.Duration` with exponential multiplier + random jitter  
**Rationale**:
- Spec: Max 3 attempts (1 initial + 2 retries) for GET/HEAD on 429/502/503/504
- Algorithm:
  ```
  baseDelay := 200ms
  for attempt := 0; attempt < maxRetries; attempt++ {
      delay := baseDelay * (2 ^ attempt) + rand.Intn(100ms)
      time.Sleep(delay)
  }
  ```
  - Attempt 1: ~200ms + jitter
  - Attempt 2: ~400ms + jitter
  - Total retry window: ~600ms-1s (acceptable for 30s timeout)
- Jitter prevents thundering herd if many clients retry simultaneously
- Encapsulate in `shouldRetry(statusCode, method) bool` + `backoffDuration(attempt int) time.Duration`

**Alternatives considered**:
- **Fixed delay**: No jitter; collision risk
- **Third-party library** (e.g., `cenkalti/backoff`): Adds dependency; simple enough to implement in ~30 lines

**References**:
- [AWS Architecture Blog: Exponential Backoff and Jitter](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)
- Spec FR-015, Resiliency & Retry Policy

---

### 5. Waiter Pattern for Async Operations (202 Accepted)

**Decision**: Provide `vps.WaitFor*` functions (e.g., `WaitForServerStatus`) with polling loop + context support  
**Rationale**:
- Spec: Optional waiters; non-blocking by default; configurable interval/backoff; honor context deadlines
- Pattern:
  ```go
  func (c *Client) WaitForServerStatus(ctx context.Context, serverID string, targetStatus string, opts ...WaiterOption) error {
      cfg := defaultWaiterConfig().apply(opts) // interval, maxWait, backoff
      ticker := time.NewTicker(cfg.interval)
      defer ticker.Stop()
      deadline := time.Now().Add(cfg.maxWait)
      
      for {
          select {
          case <-ctx.Done():
              return ctx.Err() // canceled or timeout
          case <-ticker.C:
              if time.Now().After(deadline):
                  return ErrWaiterTimeout
              server, err := c.GetServer(ctx, serverID)
              if err != nil { return err }
              if server.Status == targetStatus { return nil }
              // optionally apply backoff multiplier to ticker interval
          }
      }
  }
  ```
- Encapsulate in `vps/waiters.go`; tests use mock clock for determinism

**Alternatives considered**:
- **Callbacks**: Less idiomatic in Go; harder to test
- **Channels**: Over-engineered for simple polling; context is sufficient

**References**:
- [AWS SDK Go Waiters](https://aws.github.io/aws-sdk-go-v2/docs/making-requests/#using-waiters)
- Spec FR-016, Asynchronous Operations & Waiters

---

### 6. HTTP Client Configuration: Timeout and Transport Tuning

**Decision**: Wrap `http.Client` with custom `Transport` for TLS defaults; enforce 30s default timeout; allow override via `ClientOption`  
**Rationale**:
- Spec: 30s default per-request timeout; TLS verification enabled; context-based cancellation
- Implementation:
  ```go
  defaultHTTPClient := &http.Client{
      Timeout: 30 * time.Second,
      Transport: &http.Transport{
          TLSClientConfig: &tls.Config{InsecureSkipVerify: false}, // enforce TLS
          MaxIdleConns:    10,
          IdleConnTimeout: 90 * time.Second,
      },
  }
  ```
- Allow caller to override: `cloudsdk.New(baseURL, token, cloudsdk.WithHTTPClient(custom), cloudsdk.WithTimeout(60*time.Second))`
- Per-request timeout via context: `ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)`; SDK respects the earlier deadline

**Alternatives considered**:
- **No default timeout**: Risk of hung requests
- **60s default**: Too long for most API calls; 30s balances usability and safety

**References**:
- [Go net/http Client docs](https://pkg.go.dev/net/http#Client)
- Spec FR-018, Timeouts section

---

### 7. Testing Strategy: Unit + Contract Tests

**Decision**: Table-driven unit tests + `httptest`-based contract tests per Swagger endpoint  
**Rationale**:
- **Unit tests**: Isolated logic (retry, backoff, error mapping, waiters); mock HTTP using `httptest.Server`
- **Contract tests**: For each API operation, verify:
  - Request URL, method, headers (Authorization, Content-Type)
  - Request body matches Swagger schema
  - Response body parsed correctly
  - Status codes handled per spec
  - Example:
    ```go
    func TestListServers(t *testing.T) {
        tests := []struct{ name, mockResponse string; wantErr bool }{
            {"success", `{"items":[{"id":"1"}]}`, false},
            {"401", `{"errorCode":401,"message":"Unauthorized"}`, true},
        }
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    assert.Equal(t, "GET", r.Method)
                    assert.Equal(t, "/vps/api/v1/project/test-proj/servers", r.URL.Path)
                    w.WriteHeader(200)
                    w.Write([]byte(tt.mockResponse))
                }))
                defer srv.Close()
                client := vps.New(srv.URL, "token", "test-proj")
                _, err := client.ListServers(context.Background(), nil)
                if tt.wantErr { assert.Error(t, err) } else { assert.NoError(t, err) }
            })
        }
    }
    ```
- Contract tests derived directly from Swagger ensure SDK matches API exactly

**Alternatives considered**:
- **Integration tests against live API**: Requires staging environment; slow; reserved for E2E validation (out of scope for TDD)
- **Snapshot tests**: Brittle; manual fixture maintenance

**References**:
- [Effective Go: Testing](https://go.dev/doc/effective_go#testing)
- Spec SC-001, SC-002, SC-003

---

### 8. Import Cycle Prevention: The `internal/types` Pattern

**Decision**: Create `internal/types` package for shared types (SDKError, Logger) used by both root and service packages; re-export from root for public API  
**Rationale**:
- **Problem**: Go prohibits circular imports. If `cloudsdk` package imports `modules/vps`, and both need `SDKError`, we get:
  ```
  cloudsdk → modules/vps → cloudsdk (CYCLE!)
  ```
- **Solution**: Three-layer architecture:
  1. `internal/types/types.go` - defines SDKError, Logger (no external imports)
  2. `errors.go` (root) - re-exports types for public API: `type SDKError = types.SDKError`
  3. `modules/vps/` - imports `internal/types` directly (not root package)
- **Re-export pattern** maintains clean public API while avoiding cycles:
  ```go
  // errors.go (root) - Public API surface
  package cloudsdk
  
  import "github.com/Zillaforge/cloud-sdk/internal/types"
  
  type SDKError = types.SDKError  // Type alias, zero cost
  
  // Re-export constructor helpers for ergonomic error creation
  func NewSDKError(...) *SDKError { return types.NewSDKError(...) }
  func NewNetworkError(...) *SDKError { return types.NewNetworkError(...) }
  func NewTimeoutError(...) *SDKError { return types.NewTimeoutError(...) }
  func NewCanceledError(...) *SDKError { return types.NewCanceledError(...) }
  func NewHTTPError(...) *SDKError { return types.NewHTTPError(...) }
  ```
- **Logger interface** also lives in `internal/types` and is re-exported as `type Logger = types.Logger` at root

**Package visibility rules**:
- `internal/types` can be imported by any package within `cloud-sdk` module
- External consumers use `cloudsdk.SDKError` and `cloudsdk.Logger` (never import internal directly)
- Service packages (`modules/vps`) import `internal/types` for implementation, root package for composition

**Alternatives considered**:
- **Duplicate error types**: Violates DRY; causes type mismatch issues
- **Interface-only in root, concrete in services**: Leaks implementation details; still causes cycles
- **Vendor pattern**: Over-engineered for SDK use case

**References**:
- [Go Command: Internal directories](https://go.dev/doc/go1.4#internalpackages)
- [Go Wiki: Import cycle resolution patterns](https://github.com/golang/go/wiki/CodeReviewComments#import-dot)

---

### 9. Logger Interface Design for Pluggable Observability

**Decision**: Structured logging interface with key-value pairs; optional (nil-safe) injection via ClientOption  
**Rationale**:
- **Interface signature** (in `internal/types/types.go`):
  ```go
  type Logger interface {
      Debug(msg string, keysAndValues ...interface{})
      Info(msg string, keysAndValues ...interface{})
      Error(msg string, keysAndValues ...interface{})
  }
  ```
- **Design choices**:
  - **Structured logging**: Key-value pairs (`logger.Debug("retry attempt", "attempt", 2, "status", 503)`) over printf-style for machine-readability
  - **Three levels only**: Debug (retry/waiter details), Info (major operations), Error (failures) - sufficient for SDK telemetry
  - **Compatible with popular loggers**: Works with zap, logrus, slog (Go 1.21+), zerolog with thin adapters
  - **Optional**: Logger field in Client is `Logger` interface (can be nil); all internal code checks `if c.logger != nil` before logging
- **What gets logged**:
  - Retry attempts: `logger.Debug("retrying request", "method", "GET", "url", "/servers", "attempt", 2, "backoff", "400ms")`
  - Waiter polls: `logger.Debug("polling resource state", "resource", "server", "id", "abc", "target", "ACTIVE", "current", "BUILD")`
  - Request initiation: `logger.Info("HTTP request", "method", "POST", "path", "/servers")` (no token/secrets)
- **Security**: Token values never logged; only presence indicated: `"auth", "present"`

**Alternatives considered**:
- **Printf-style logger** (`func(format string, args ...interface{})`): Less structured; harder to parse programmatically
- **Context-based logging**: Over-engineered; SDK doesn't need request-scoped loggers
- **Required logger**: Friction for simple use cases; optional is more Go-idiomatic

**References**:
- [slog package (Go 1.21+)](https://pkg.go.dev/log/slog)
- Constitution: "expose hooks for logging / observability (pluggable logger interfaces)"

---

### 10. ClientOption Functional Options Pattern

**Decision**: Implement ClientOption as functional option in `client.go` (inline, no separate options.go file); apply after struct creation  
**Rationale**:
- **Pattern** (Functional Options by Dave Cheney):
  ```go
  type ClientOption func(*Client)
  
  func WithLogger(logger Logger) ClientOption {
      return func(c *Client) { c.logger = logger }
  }
  
  func WithTimeout(timeout time.Duration) ClientOption {
      return func(c *Client) { c.httpClient.Timeout = timeout }
  }
  
  func WithHTTPClient(httpClient *http.Client) ClientOption {
      return func(c *Client) { c.httpClient = httpClient }
  }
  
  func New(baseURL, token string, opts ...ClientOption) (*Client, error) {
      // 1. Create with defaults
      client := &Client{
          baseURL: baseURL,
          token:   token,
          httpClient: &http.Client{Timeout: 30 * time.Second},
          logger: nil, // optional
      }
      // 2. Apply options (order matters if they conflict)
      for _, opt := range opts {
          opt(client)
      }
      return client, nil
  }
  ```
- **Why inline in client.go**:
  - Small number of options (3-4 total)
  - Options tightly coupled to Client struct fields
  - Keeps all Client initialization logic in one file
  - Avoids over-engineering with separate `options.go` for <50 lines
- **Type re-export**: `Logger` type must be re-exported at root level so callers don't import internal: `type Logger = types.Logger`

**Alternatives considered**:
- **Separate options.go file**: Overkill for 3 options; would split related initialization code
- **Builder pattern**: More verbose; functional options more idiomatic in Go
- **Config struct**: Loses composability; all-or-nothing configuration

**References**:
- [Dave Cheney: Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

---

### 11. Multi-Service Access Pattern: ProjectClient Intermediate Type

**Decision**: Two-step access pattern `Client.Project(id) → ProjectClient → service.VPS()` with intermediate selector type  
**Rationale**:
- **Architecture**:
  ```go
  // Step 1: Create root client (auth, base URL, HTTP client)
  client := cloudsdk.New("https://api.example.com", "token-123")
  
  // Step 2: Select project (binds project-id, returns project-scoped handle)
  projectClient := client.Project("proj-abc")
  
  // Step 3: Get service client (VPS, Compute, Storage, etc.)
  vpsClient := projectClient.VPS()
  
  // Step 4: Call operations (no project-id needed in params)
  servers, err := vpsClient.Servers().List(ctx, nil)
  ```
- **Why intermediate ProjectClient type**:
  - **Extensibility**: Future project-level operations (get quotas, list services, configure region) without polluting root Client
  - **Multiple projects**: Users can create multiple project-scoped clients from same root: `proj1.VPS()`, `proj2.VPS()`
  - **Clean separation**: Root client = auth/transport; ProjectClient = project context; service client = operations
- **ProjectClient struct** (in `client.go`):
  ```go
  type ProjectClient struct {
      client    *Client  // reference to root (for auth, HTTP)
      projectID string   // bound project
  }
  
  func (c *Client) Project(projectID string) *ProjectClient {
      return &ProjectClient{client: c, projectID: projectID}
  }
  
  func (p *ProjectClient) VPS() *vps.Client {
      return vps.NewClient(
          p.client.baseURL,
          p.client.token,
          p.projectID,
          p.client.httpClient,
          p.client.logger,
      )
  }
  ```
- **Parameter flow**: Root client fields (baseURL, token, httpClient, logger) + project-id flow into service client constructor

**Alternatives considered**:
- **Direct factory**: `client.VPS(projectID)` simpler but loses extensibility; hard to add project-level methods later
- **Service client holds multiple projects**: Violates single responsibility; complicates state management

**References**:
- Similar pattern: AWS SDK v2 uses `cfg → client` with service-specific clients
- Spec: "SDK provides a way to select a project that returns a project-scoped handle"

---

### 12. Service Client Composition and basePath Construction

**Decision**: Service clients compose `internal/http.Client` (not inherit from root Client); construct basePath with project-id  
**Rationale**:
- **VPS Client structure** (in `modules/vps/client.go`):
  ```go
  type Client struct {
      baseClient *internalhttp.Client  // HTTP wrapper with retry/timeout
      projectID  string                // bound project
      basePath   string                // "/vps/api/v1/project/{project-id}"
  }
  
  func NewClient(baseURL, token, projectID string, httpClient *http.Client, logger types.Logger) *Client {
      basePath := "/vps/api/v1/project/" + projectID
      return &Client{
          baseClient: internalhttp.NewClient(baseURL, token, httpClient, logger),
          projectID:  projectID,
          basePath:   basePath,
      }
  }
  ```
- **Why 5 parameters**: baseURL, token, projectID (identity); httpClient, logger (config) - all needed for complete service client
- **basePath construction**: VPS client owns path template; ProjectClient doesn't need to know service-specific URL structure
- **Composition over inheritance**: `internalhttp.Client` provides HTTP primitives (Do, retry, error mapping); VPS Client adds resource operations

**Path construction in operations**:
```go
func (c *Client) ListServers(ctx context.Context, opts *ListOptions) ([]Server, error) {
    path := c.basePath + "/servers"
    // Use c.baseClient.Do(ctx, "GET", path, nil, &result)
}
```

**Alternatives considered**:
- **Root client in service client**: Tight coupling; violates encapsulation; harder to test
- **Pass ProjectClient to service constructor**: Couples service to project abstraction unnecessarily

**References**:
- Plan.md VPS Client structure section

---

## Summary

All technical unknowns resolved:
1. ✅ Go module structure: Mono-repo, service packages (`vps/`)
2. ✅ Models: Hand-written with Swagger reference
3. ✅ Errors: Custom `SDKError` with wrapping
4. ✅ Retry: Exponential backoff + jitter in `internal/backoff`
5. ✅ Waiters: Polling loop + context in `vps/waiters.go`
6. ✅ HTTP: Custom `http.Client` with 30s timeout, TLS enforced
7. ✅ Testing: Table-driven unit + `httptest` contract tests
8. ✅ Import cycles: `internal/types` pattern with re-exports
9. ✅ Logger: Structured interface, optional, nil-safe
10. ✅ Options: Functional options inline in `client.go`
11. ✅ Access pattern: ProjectClient intermediate type
12. ✅ Service composition: VPS Client owns basePath, composes internal HTTP client

**Next**: Phase 1 - Generate data-model.md and contracts/
