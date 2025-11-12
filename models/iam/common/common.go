// Package common provides shared types used across IAM models.
package common

// Permission represents access control permissions.
type Permission struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// TenantRole represents a user's role within a project/tenant.
type TenantRole string

const (
	// TenantRoleMember represents a standard project member
	TenantRoleMember TenantRole = "TENANT_MEMBER"
	// TenantRoleAdmin represents a project administrator
	TenantRoleAdmin TenantRole = "TENANT_ADMIN"
	// TenantRoleOwner represents a project owner (highest authority)
	TenantRoleOwner TenantRole = "TENANT_OWNER"
)

// String returns the string representation of the TenantRole.
func (r TenantRole) String() string {
	return string(r)
}

// IsValid checks if the TenantRole is a valid value.
func (r TenantRole) IsValid() bool {
	switch r {
	case TenantRoleMember, TenantRoleAdmin, TenantRoleOwner:
		return true
	}
	return false
}
