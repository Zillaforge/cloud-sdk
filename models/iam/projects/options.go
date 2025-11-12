// Package projects provides data models for IAM project operations.
package projects

// ListProjectsOptions represents optional parameters for listing projects.
// Uses single options struct pattern consistent with VPS/VRM modules.
type ListProjectsOptions struct {
	Offset *int
	Limit  *int
	Order  *string
}
