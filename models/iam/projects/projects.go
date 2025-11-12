// Package projects provides data models for IAM project operations.
package projects

import "github.com/Zillaforge/cloud-sdk/models/iam/common"

// Project represents a project/tenant in the multi-tenant system.
type Project struct {
	ProjectID   string                 `json:"projectId"`
	DisplayName string                 `json:"displayName"`
	Description string                 `json:"description"`
	Extra       map[string]interface{} `json:"extra"`
	Namespace   string                 `json:"namespace"`
	Frozen      bool                   `json:"frozen"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
}

// ProjectMembership wraps a Project with user-specific membership context.
type ProjectMembership struct {
	Project            *Project               `json:"project"`
	GlobalPermissionID string                 `json:"globalPermissionId"`
	GlobalPermission   *common.Permission     `json:"globalPermission"`
	UserPermissionID   string                 `json:"userPermissionId"`
	UserPermission     *common.Permission     `json:"userPermission"`
	Frozen             bool                   `json:"frozen"`
	TenantRole         common.TenantRole      `json:"tenantRole"`
	Extra              map[string]interface{} `json:"extra"`
}

// ListProjectsResponse represents the response from the GET /projects endpoint.
type ListProjectsResponse struct {
	Projects []*ProjectMembership `json:"projects"`
	Total    int                  `json:"total"`
}

// GetProjectResponse represents the response from the GET /project/{id} endpoint.
type GetProjectResponse struct {
	ProjectID        string                 `json:"projectId"`
	DisplayName      string                 `json:"displayName"`
	Description      string                 `json:"description"`
	Extra            map[string]interface{} `json:"extra"`
	Namespace        string                 `json:"namespace"`
	Frozen           bool                   `json:"frozen"`
	GlobalPermission *common.Permission     `json:"globalPermission"`
	UserPermission   *common.Permission     `json:"userPermission"`
	CreatedAt        string                 `json:"createdAt"`
	UpdatedAt        string                 `json:"updatedAt"`
}
