## Purpose

Provides a typed, project-scoped SDK surface for the VPS project APIs defined by the authoritative Swagger contract.

## ADDED Requirements

### Requirement: Project-scoped VPS client
The SDK MUST provide a VPS client bound to a project ID, bearer-token authenticated, and exposing typed context-aware methods without requiring callers to construct raw URLs or headers.

#### Scenario: Initialize project client
- **WHEN** a caller selects a project from a client initialized with a base URL and bearer token
- **THEN** the VPS client uses that project for all project-scoped requests and includes authorization automatically

### Requirement: Project resource operations
The VPS client MUST support the documented project-scoped operations for networks, floating IPs, servers, keypairs, routers, security groups, flavors, and quotas, including each resource's documented subresources and actions.

#### Scenario: Execute documented resource operation
- **WHEN** a caller invokes a supported typed operation with valid input and context
- **THEN** the client sends the corresponding `/api/v1/project/{project-id}/` request and decodes the Swagger-aligned result

### Requirement: Structured errors and safe retries
The SDK MUST return structured errors preserving HTTP status and server details, map client-side failures to status code zero, and retry only safe reads on 429/502/503/504 for at most three attempts with backoff and jitter.

#### Scenario: Transient read failure
- **WHEN** a GET request receives a retryable transient response
- **THEN** the client retries within the configured limit and returns a structured error if all attempts fail

### Requirement: Asynchronous operation support
The SDK MUST return promptly for documented 202 operations and provide optional context-aware waiter helpers for callers that need to poll for a desired state.

#### Scenario: Wait for asynchronous state
- **WHEN** a caller invokes a waiter after an operation returns 202
- **THEN** polling follows the caller's context and configured backoff and ends when the desired state or an error is observed
