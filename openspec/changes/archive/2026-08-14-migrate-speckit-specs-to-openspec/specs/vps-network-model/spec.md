## Purpose

Represents VPS network responses and requests with the complete Swagger-defined field set and serialization behavior.

## ADDED Requirements

### Requirement: Complete network representation
The Network model MUST expose all documented network fields, including gateway, nameservers, router, project and user references, status, status reason, bonding, and shared state.

#### Scenario: Decode complete network
- **WHEN** a Swagger-shaped network response is decoded
- **THEN** every documented field is populated and accessible through the typed model

### Requirement: Network serialization contract
Network request and response JSON MUST use the field names and types defined by the API contract.

#### Scenario: Round-trip network data
- **WHEN** a fully populated Network model is marshaled and unmarshaled
- **THEN** documented values are retained without field-name or type loss
