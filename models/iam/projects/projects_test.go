package projects_test

import (
	"encoding/json"
	"testing"

	"github.com/Zillaforge/cloud-sdk/models/iam/common"
	"github.com/Zillaforge/cloud-sdk/models/iam/projects"
)

func TestProject_JSONParsing(t *testing.T) {
	jsonData := `{
		"projectId": "91457b61-0b92-4aa8-b136-b03d88f04946",
		"displayName": "prj1762875055667",
		"description": "Test project",
		"extra": {
			"iservice": {
				"projectSysCode": "TCI111222"
			}
		},
		"namespace": "ci.asus.com",
		"frozen": false,
		"createdAt": "2025-11-11T15:31:02Z",
		"updatedAt": "2025-11-11T15:31:06Z"
	}`

	var project projects.Project
	err := json.Unmarshal([]byte(jsonData), &project)
	if err != nil {
		t.Fatalf("Failed to unmarshal project: %v", err)
	}

	if project.ProjectID != "91457b61-0b92-4aa8-b136-b03d88f04946" {
		t.Errorf("ProjectID = %v, want %v", project.ProjectID, "91457b61-0b92-4aa8-b136-b03d88f04946")
	}
	if project.DisplayName != "prj1762875055667" {
		t.Errorf("DisplayName = %v, want %v", project.DisplayName, "prj1762875055667")
	}
	if project.Namespace != "ci.asus.com" {
		t.Errorf("Namespace = %v, want %v", project.Namespace, "ci.asus.com")
	}
	if project.Frozen != false {
		t.Errorf("Frozen = %v, want %v", project.Frozen, false)
	}
}

func TestProjectMembership_JSONParsing(t *testing.T) {
	jsonData := `{
		"project": {
			"projectId": "test-project-id",
			"displayName": "Test Project",
			"description": "",
			"extra": {},
			"namespace": "test.com",
			"frozen": false,
			"createdAt": "2025-01-01T00:00:00Z",
			"updatedAt": "2025-01-01T00:00:00Z"
		},
		"globalPermissionId": "perm-global-id",
		"globalPermission": {
			"id": "perm-global-id",
			"label": "DEFAULT"
		},
		"userPermissionId": "perm-user-id",
		"userPermission": {
			"id": "perm-user-id",
			"label": "ADMIN"
		},
		"frozen": false,
		"tenantRole": "TENANT_MEMBER",
		"extra": {}
	}`

	var membership projects.ProjectMembership
	err := json.Unmarshal([]byte(jsonData), &membership)
	if err != nil {
		t.Fatalf("Failed to unmarshal project membership: %v", err)
	}

	if membership.Project == nil {
		t.Fatal("Project should not be nil")
	}
	if membership.Project.ProjectID != "test-project-id" {
		t.Errorf("Project.ProjectID = %v, want %v", membership.Project.ProjectID, "test-project-id")
	}
	if membership.GlobalPermissionID != "perm-global-id" {
		t.Errorf("GlobalPermissionID = %v, want %v", membership.GlobalPermissionID, "perm-global-id")
	}
	if membership.TenantRole != common.TenantRoleMember {
		t.Errorf("TenantRole = %v, want %v", membership.TenantRole, common.TenantRoleMember)
	}
}

func TestListProjectsResponse_Parsing(t *testing.T) {
	jsonData := `{
		"projects": [
			{
				"project": {
					"projectId": "project-1",
					"displayName": "Project 1",
					"description": "",
					"extra": {},
					"namespace": "test.com",
					"frozen": false,
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z"
				},
				"globalPermissionId": "perm-1",
				"globalPermission": {
					"id": "perm-1",
					"label": "DEFAULT"
				},
				"userPermissionId": "perm-1",
				"userPermission": {
					"id": "perm-1",
					"label": "DEFAULT"
				},
				"frozen": false,
				"tenantRole": "TENANT_OWNER",
				"extra": {}
			},
			{
				"project": {
					"projectId": "project-2",
					"displayName": "Project 2",
					"description": "",
					"extra": {},
					"namespace": "test.com",
					"frozen": false,
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z"
				},
				"globalPermissionId": "perm-2",
				"globalPermission": {
					"id": "perm-2",
					"label": "DEFAULT"
				},
				"userPermissionId": "perm-2",
				"userPermission": {
					"id": "perm-2",
					"label": "DEFAULT"
				},
				"frozen": false,
				"tenantRole": "TENANT_MEMBER",
				"extra": {}
			}
		],
		"total": 2
	}`

	var response projects.ListProjectsResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal list projects response: %v", err)
	}

	if len(response.Projects) != 2 {
		t.Errorf("len(Projects) = %v, want %v", len(response.Projects), 2)
	}
	if response.Total != 2 {
		t.Errorf("Total = %v, want %v", response.Total, 2)
	}
	if response.Projects[0].Project.ProjectID != "project-1" {
		t.Errorf("Projects[0].Project.ProjectID = %v, want %v", response.Projects[0].Project.ProjectID, "project-1")
	}
}

func TestListProjectsOptions_StructFields(t *testing.T) {
	// Test nil options
	var opts *projects.ListProjectsOptions
	if opts != nil {
		t.Error("nil options should be nil")
	}

	// Test with values
	offset := 10
	limit := 20
	order := "displayName"

	opts = &projects.ListProjectsOptions{
		Offset: &offset,
		Limit:  &limit,
		Order:  &order,
	}

	if opts.Offset == nil || *opts.Offset != 10 {
		t.Errorf("Offset = %v, want %v", opts.Offset, 10)
	}
	if opts.Limit == nil || *opts.Limit != 20 {
		t.Errorf("Limit = %v, want %v", opts.Limit, 20)
	}
	if opts.Order == nil || *opts.Order != "displayName" {
		t.Errorf("Order = %v, want %v", opts.Order, "displayName")
	}
}

func TestGetProjectResponse_Parsing(t *testing.T) {
	jsonData := `{
		"projectId": "test-project-id",
		"displayName": "Test Project",
		"description": "Test description",
		"extra": {
			"metadata": "value"
		},
		"namespace": "test.com",
		"frozen": false,
		"globalPermission": {
			"id": "global-perm-id",
			"label": "DEFAULT"
		},
		"userPermission": {
			"id": "user-perm-id",
			"label": "ADMIN"
		},
		"createdAt": "2025-01-01T00:00:00Z",
		"updatedAt": "2025-01-02T00:00:00Z"
	}`

	var response projects.GetProjectResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal GetProjectResponse: %v", err)
	}

	if response.ProjectID != "test-project-id" {
		t.Errorf("ProjectID = %v, want %v", response.ProjectID, "test-project-id")
	}
	if response.DisplayName != "Test Project" {
		t.Errorf("DisplayName = %v, want %v", response.DisplayName, "Test Project")
	}
	if response.GlobalPermission == nil || response.GlobalPermission.Label != "DEFAULT" {
		t.Error("GlobalPermission should have DEFAULT label")
	}
	if response.UserPermission == nil || response.UserPermission.Label != "ADMIN" {
		t.Error("UserPermission should have ADMIN label")
	}
	if response.Frozen != false {
		t.Errorf("Frozen = %v, want %v", response.Frozen, false)
	}
}
