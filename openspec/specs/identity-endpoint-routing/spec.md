## Purpose

Enables Cloud SDK callers to direct IAM traffic to a dedicated identity endpoint while other service traffic continues using the general API endpoint.

## Requirements

### Requirement: Dedicated identity endpoint configuration
The SDK SHALL provide a public constructor that accepts a base URL, an identity URL, a bearer token, and the existing client options. The constructor MUST reject an identity URL that is not a valid URL with a scheme, using validation behavior equivalent to the base URL validation.

#### Scenario: Construct client with valid dedicated identity URL
- **WHEN** a caller provides valid base and identity URLs, a non-empty token, and client options
- **THEN** the SDK returns a configured client that preserves the supplied options

#### Scenario: Reject identity URL without a scheme
- **WHEN** a caller provides an identity URL without a URL scheme
- **THEN** the SDK returns an error and does not return a client

### Requirement: Service-specific endpoint routing
When a client is constructed with a dedicated identity URL, the SDK SHALL send IAM operations to that identity URL with the IAM service path appended. The SDK MUST continue sending VPS and VRM operations to paths derived from the base URL.

#### Scenario: IAM request uses dedicated identity endpoint
- **WHEN** a caller invokes an IAM operation on a client constructed with a dedicated identity URL
- **THEN** the request is sent to the identity URL's IAM endpoint rather than the base URL

#### Scenario: Project lookup uses dedicated identity endpoint
- **WHEN** a caller resolves a project by ID or project code on a client constructed with a dedicated identity URL
- **THEN** every IAM request made during the lookup is sent to the identity URL's IAM endpoint

#### Scenario: Project-scoped non-IAM request uses base endpoint
- **WHEN** a caller obtains a VPS or VRM client from a client constructed with a dedicated identity URL
- **THEN** that service client is configured from the base URL and not the identity URL
