# Data Model: Fix Network Model Definition

## Network
- **Description**: Represents a virtual network returned by the VPS API (`pb.NetworkInfo`).
- **Fields**:
  - `id` (`string`): Unique identifier.
  - `name` (`string`)
  - `description` (`string`, optional)
  - `cidr` (`string`)
  - `bonding` (`bool`, optional)
  - `gateway` (`string`, optional)
  - `gw_state` (`bool`, optional, deprecated)
  - `is_default` (`bool`, optional)
  - `nameservers` (`[]string`, optional)
  - `namespace` (`string`, optional)
  - `project` (`*IDName`, optional)
  - `project_id` (`string`, optional)
  - `router` (`*RouterInfo`, optional)
  - `router_id` (`string`, optional)
  - `shared` (`bool`, optional)
  - `status` (`string`, optional)
  - `status_reason` (`string`, optional)
  - `subnet_id` (`string`, optional)
  - `user` (`*IDName`, optional)
  - `user_id` (`string`, optional)
  - `createdAt` (`string`)
  - `updatedAt` (`string`, optional)

## IDName
- **Description**: Shared reference object for `project` and `user` embeddings.
- **Fields**:
  - `id` (`string`)
  - `name` (`string`)

## RouterInfo
- **Description**: Lightweight representation of a router attached to a network, mirroring `pb.RouterInfo`.
- **Fields**:
  - `id` (`string`)
  - `name` (`string`)
  - `description` (`string`, optional)
  - `bonding` (`bool`, optional)
  - `is_default` (`bool`, optional)
  - `shared` (`bool`, optional)
  - `state` (`bool`, optional)
  - `status` (`string`, optional)
  - `status_reason` (`string`, optional)
  - `namespace` (`string`, optional)
  - `project` (`*IDName`, optional)
  - `project_id` (`string`, optional)
  - `user` (`*IDName`, optional)
  - `user_id` (`string`, optional)
  - `extnetwork` (`*ExtNetworkInfo`, optional)
  - `extnetwork_id` (`string`, optional)
  - `gw_addrs` (`[]string`, optional)
  - `createdAt` (`string`, optional)
  - `updatedAt` (`string`, optional)

## ExtNetworkInfo
- **Description**: Embedded external network summary referenced by routers (`pb.ExtNetInfo`).
- **Fields**:
  - `id` (`string`)
  - `name` (`string`)
  - `description` (`string`, optional)
  - `cidr` (`string`, optional)
  - `namespace` (`string`, optional)
  - `segment_id` (`string`, optional)
  - `type` (`string`, optional)
  - `is_default` (`bool`, optional)

## NetworkCreateRequest
- **Description**: Matches `NetCreateInput` for network creation.
- **Fields**:
  - `name` (`string`)
  - `description` (`string`, optional)
  - `cidr` (`string`)
  - `gateway` (`string`, optional)
  - `router_id` (`string`, optional)

## NetworkUpdateRequest
- **Description**: Matches `NetUpdateInput` (unchanged per spec).
- **Fields**:
  - `name` (`string`, optional)
  - `description` (`string`, optional)

## NetworkListResponse
- **Description**: Response wrapper for list operations.
- **Fields**:
  - `networks` (`[]*Network`)

## ListNetworksOptions
- **Description**: Filter parameters for list operations (existing behavior retained).
- **Fields**:
  - `Name` (`string`, optional)
  - *(Future pagination fields remain unchanged.)*

## NetworkPort
- **Description**: Represents an individual port returned with network associations (`NetPort`).
- **Fields**:
  - `id` (`string`)
  - `addresses` (`[]string`, optional)
  - `server` (`*ServerSummary`, optional)

## ServerSummary (port context)
- **Description**: Minimal server summary embedded within `NetPort` when a port is attached to a server, mapping to `pb.ServerInfo` subset.
- **Fields**:
  - `id` (`string`)
  - `name` (`string`, optional)
  - `status` (`string`, optional)
  - `project_id` (`string`, optional)
  - `user_id` (`string`, optional)
