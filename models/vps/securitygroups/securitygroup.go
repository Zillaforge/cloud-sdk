// Package securitygroups provides data structures for VPS security group resources.
package securitygroups

// SecurityGroup represents a security group in the VPS service.
type SecurityGroup struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	ProjectID   string              `json:"project_id"`
	UserID      string              `json:"user_id"`
	Namespace   string              `json:"namespace"`
	Rules       []SecurityGroupRule `json:"rules"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
	Project     *IDName             `json:"project,omitempty"`
	User        *IDName             `json:"user,omitempty"`
}

// IDName represents a simple ID and name pair.
type IDName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SecurityGroupCreateRequest represents the request body for creating a security group.
type SecurityGroupCreateRequest struct {
	Name        string                           `json:"name"`
	Description string                           `json:"description,omitempty"`
	Rules       []SecurityGroupRuleCreateRequest `json:"rules,omitempty"`
}

// SecurityGroupUpdateRequest represents the request body for updating a security group.
type SecurityGroupUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListSecurityGroupsOptions represents query parameters for listing security groups.
type ListSecurityGroupsOptions struct {
	Name   string
	UserID string
	Detail bool
}

// SecurityGroupListResponse represents the response from listing security groups.
// Note: The Total field has been removed as it is not present in the API specification (pb.SgListOutput).
// Use len(resp.SecurityGroups) to count results.
type SecurityGroupListResponse struct {
	SecurityGroups []SecurityGroup `json:"security_groups"`
}
