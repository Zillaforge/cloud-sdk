## Purpose

Represents VPS flavor data using the API-defined names, types, enum values, and JSON serialization contract.

## ADDED Requirements

### Requirement: Swagger-aligned flavor model
The Flavor model MUST expose all fields defined by the Swagger specification, including correctly named CPU and memory fields and documented enum-like values.

#### Scenario: Decode flavor response
- **WHEN** a valid flavor response is unmarshaled
- **THEN** all documented fields are populated in the Flavor model

### Requirement: Flavor JSON compatibility
Flavor models MUST marshal using the API's documented JSON names and structure.

#### Scenario: Marshal flavor
- **WHEN** a Flavor model containing all supported fields is marshaled
- **THEN** the JSON matches the API contract without legacy field aliases
