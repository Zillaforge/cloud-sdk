package contracts

import (
	"context"

	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

// Client provides operations for managing security groups.
//
// This interface defines the contract for security group CRUD operations
// following the VPS API specification (swagger/vps.yaml). All methods accept
// context.Context for cancellation/timeout and return typed responses with
// wrapped errors.
//
// Usage:
//
//	client := sdk.Project("proj-123").VPS().SecurityGroups()
//	sg, err := client.Create(ctx, &securitygroups.SecurityGroupCreateRequest{...})
type Client interface {
	// List retrieves all security groups for the project with optional filters.
	//
	// API Endpoint: GET /api/v1/project/{project-id}/security_groups
	//
	// Query Parameters:
	//   - name: Filter by exact security group name
	//   - user_id: Filter by owner user ID (admin only)
	//   - detail: Include rules array in response (default: false)
	//
	// Returns:
	//   - *SecurityGroupListResponse containing array of security groups
	//   - Error (nil on success, sdkError with HTTP status on failure)
	//
	// Example:
	//   // List all security groups (no rules)
	//   resp, err := client.List(ctx, nil)
	//
	//   // List security groups with rules
	//   detail := true
	//   resp, err := client.List(ctx, &securitygroups.ListSecurityGroupsOptions{Detail: &detail})
	List(ctx context.Context, opts *securitygroups.ListSecurityGroupsOptions) (*securitygroups.SecurityGroupListResponse, error)

	// Create provisions a new security group with optional initial rules.
	//
	// API Endpoint: POST /api/v1/project/{project-id}/security_groups
	//
	// Request Body: SecurityGroupCreateRequest
	//   - name: Required, non-empty security group name
	//   - description: Optional description
	//   - rules: Optional array of initial rules (can add later via Rules sub-resource)
	//
	// Returns:
	//   - *SecurityGroup containing the created resource
	//   - Error (nil on success, sdkError with HTTP status on failure)
	//
	// HTTP Status Codes:
	//   - 201 Created: Security group created successfully
	//   - 400 Bad Request: Invalid request body or missing required fields
	//   - 409 Conflict: Security group with same name already exists
	//   - 500 Internal Server Error: Server-side error
	//
	// Example:
	//   req := &securitygroups.SecurityGroupCreateRequest{
	//       Name:        "web-servers",
	//       Description: "Security group for web servers",
	//   }
	//   sg, err := client.Create(ctx, req)
	Create(ctx context.Context, req *securitygroups.SecurityGroupCreateRequest) (*securitygroups.SecurityGroup, error)

	// Get retrieves a single security group by ID and returns a resource wrapper
	// that provides access to sub-resource operations (rules).
	//
	// API Endpoint: GET /api/v1/project/{project-id}/security_groups/{sg-id}
	//
	// Returns:
	//   - *SecurityGroupResource wrapping the security group with Rules() accessor
	//   - Error (nil on success, sdkError with HTTP status on failure)
	//
	// HTTP Status Codes:
	//   - 200 OK: Security group retrieved successfully
	//   - 404 Not Found: Security group does not exist
	//   - 500 Internal Server Error: Server-side error
	//
	// Example:
	//   sg, err := client.Get(ctx, "sg-abc123")
	//   if err != nil {
	//       return err
	//   }
	//   // Access sub-resource operations
	//   rules := sg.Rules()
	//   rule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{...})
	Get(ctx context.Context, id string) (*SecurityGroupResource, error)

	// Update modifies a security group's name and/or description.
	//
	// API Endpoint: PUT /api/v1/project/{project-id}/security_groups/{sg-id}
	//
	// Request Body: SecurityGroupUpdateRequest
	//   - name: Optional, new security group name (nil = no change)
	//   - description: Optional, new description (nil = no change)
	//
	// Note: Rules cannot be updated via this endpoint. Use sg.Rules().Create() and
	// sg.Rules().Delete() to manage rules.
	//
	// Returns:
	//   - *SecurityGroup containing the updated resource
	//   - Error (nil on success, sdkError with HTTP status on failure)
	//
	// HTTP Status Codes:
	//   - 200 OK: Security group updated successfully
	//   - 400 Bad Request: Invalid request body
	//   - 404 Not Found: Security group does not exist
	//   - 409 Conflict: Name conflict with existing security group
	//   - 500 Internal Server Error: Server-side error
	//
	// Example:
	//   name := "web-servers-v2"
	//   req := &securitygroups.SecurityGroupUpdateRequest{Name: &name}
	//   sg, err := client.Update(ctx, "sg-abc123", req)
	Update(ctx context.Context, id string, req *securitygroups.SecurityGroupUpdateRequest) (*securitygroups.SecurityGroup, error)

	// Delete removes a security group and all its rules.
	//
	// API Endpoint: DELETE /api/v1/project/{project-id}/security_groups/{sg-id}
	//
	// Returns:
	//   - Error (nil on success, sdkError with HTTP status on failure)
	//
	// HTTP Status Codes:
	//   - 204 No Content: Security group deleted successfully
	//   - 404 Not Found: Security group does not exist
	//   - 409 Conflict: Security group still in use (attached to instances)
	//   - 500 Internal Server Error: Server-side error
	//
	// Example:
	//   err := client.Delete(ctx, "sg-abc123")
	Delete(ctx context.Context, id string) error
}

// SecurityGroupResource wraps a SecurityGroup with sub-resource operations.
//
// This struct is returned by Client.Get() and provides access to rule management
// operations via the Rules() method. It embeds the full SecurityGroup model for
// direct field access.
//
// Example:
//
//	sg, err := client.SecurityGroups().Get(ctx, "sg-abc123")
//	fmt.Println(sg.Name)  // Direct access to security group fields
//
//	// Access rule sub-resource operations
//	rules := sg.Rules()
//	rule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{
//	    Direction:  "ingress",
//	    Protocol:   "tcp",
//	    PortMin:    intPtr(22),
//	    PortMax:    intPtr(22),
//	    RemoteCIDR: "0.0.0.0/0",
//	})
type SecurityGroupResource struct {
	*securitygroups.SecurityGroup
	rulesOps RuleOperations
}

// Rules returns the rule operations for this security group.
//
// This method provides access to the sub-resource operations for managing
// security group rules (create, delete). Rules are scoped to the parent
// security group and automatically include the security group ID in API paths.
//
// Returns:
//   - RuleOperations interface for creating and deleting rules
//
// Example:
//
//	sg, err := client.Get(ctx, "sg-abc123")
//	rules := sg.Rules()
//	rule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{...})
func (sgr *SecurityGroupResource) Rules() RuleOperations {
	return sgr.rulesOps
}

// RuleOperations defines operations on security group rules (sub-resource).
//
// This interface is returned by SecurityGroupResource.Rules() and provides
// scoped operations for managing rules within a specific security group.
// All operations automatically include the parent security group ID in API paths.
//
// Note: Rules can only be created or deleted. To modify a rule, delete the
// existing rule and create a new one with updated parameters.
type RuleOperations interface {
	// Create adds a new rule to the security group.
	//
	// API Endpoint: POST /api/v1/project/{project-id}/security_groups/{sg-id}/rules
	//
	// Request Body: SecurityGroupRuleCreateRequest
	//   - direction: Required ("ingress" or "egress")
	//   - protocol: Required ("tcp", "udp", "icmp")
	//   - port_min: Optional, required for TCP/UDP (0-65535)
	//   - port_max: Optional, required for TCP/UDP (0-65535, must be >= port_min)
	//   - remote_cidr: Required, valid CIDR notation (e.g., "0.0.0.0/0")
	//
	// Returns:
	//   - *SecurityGroupRule containing the created rule
	//   - Error (nil on success, sdkError with HTTP status on failure)
	//
	// HTTP Status Codes:
	//   - 201 Created: Rule created successfully
	//   - 400 Bad Request: Invalid request body, validation errors
	//   - 404 Not Found: Parent security group does not exist
	//   - 409 Conflict: Duplicate rule (same direction, protocol, ports, CIDR)
	//   - 500 Internal Server Error: Server-side error
	//
	// Example (TCP Rule):
	//   rule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{
	//       Direction:  "ingress",
	//       Protocol:   "tcp",
	//       PortMin:    intPtr(80),
	//       PortMax:    intPtr(80),
	//       RemoteCIDR: "0.0.0.0/0",
	//   })
	//
	// Example (ICMP Rule):
	//   rule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{
	//       Direction:  "ingress",
	//       Protocol:   "icmp",
	//       RemoteCIDR: "0.0.0.0/0",
	//   })
	Create(ctx context.Context, req *securitygroups.SecurityGroupRuleCreateRequest) (*securitygroups.SecurityGroupRule, error)

	// Delete removes a rule from the security group.
	//
	// API Endpoint: DELETE /api/v1/project/{project-id}/security_groups/{sg-id}/rules/{sg-rule-id}
	//
	// Returns:
	//   - Error (nil on success, sdkError with HTTP status on failure)
	//
	// HTTP Status Codes:
	//   - 204 No Content: Rule deleted successfully
	//   - 404 Not Found: Rule or parent security group does not exist
	//   - 500 Internal Server Error: Server-side error
	//
	// Example:
	//   err := rules.Delete(ctx, "rule-abc123")
	Delete(ctx context.Context, ruleID string) error
}
