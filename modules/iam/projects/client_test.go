package projects_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/iam/common"
	"github.com/Zillaforge/cloud-sdk/models/iam/projects"
	projectsClient "github.com/Zillaforge/cloud-sdk/modules/iam/projects"
)

func TestClient_List_Success_NilOptions(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/projects" {
			t.Errorf("Expected path /api/v1/projects, got %s", r.URL.Path)
		}

		// Verify no query parameters when opts is nil
		if r.URL.RawQuery != "" {
			t.Errorf("Expected no query parameters, got %s", r.URL.RawQuery)
		}

		// Return mock response
		response := projects.ListProjectsResponse{
			Projects: []*projects.ProjectMembership{
				{
					Project: &projects.Project{
						ProjectID:   "project-1",
						DisplayName: "Project 1",
						Description: "",
						Extra:       map[string]interface{}{},
						Namespace:   "test.com",
						Frozen:      false,
						CreatedAt:   "2025-01-01T00:00:00Z",
						UpdatedAt:   "2025-01-01T00:00:00Z",
					},
					GlobalPermissionID: "perm-1",
					GlobalPermission: &common.Permission{
						ID:    "perm-1",
						Label: "DEFAULT",
					},
					UserPermissionID: "perm-1",
					UserPermission: &common.Permission{
						ID:    "perm-1",
						Label: "DEFAULT",
					},
					Frozen:     false,
					TenantRole: common.TenantRoleOwner,
					Extra:      map[string]interface{}{},
				},
			},
			Total: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := projectsClient.NewClient(baseClient, "/api/v1/")

	// Call List with nil options
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := client.List(ctx, nil)

	// Verify result
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}
	if result == nil {
		t.Fatal("List() returned nil result")
	}
	if len(result) != 1 {
		t.Errorf("len(result) = %v, want %v", len(result), 1)
	}
}

func TestClient_List_WithPaginationOptions(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		if query.Get("offset") != "10" {
			t.Errorf("Expected offset=10, got %s", query.Get("offset"))
		}
		if query.Get("limit") != "20" {
			t.Errorf("Expected limit=20, got %s", query.Get("limit"))
		}
		if query.Get("order") != "displayName" {
			t.Errorf("Expected order=displayName, got %s", query.Get("order"))
		}

		// Return mock response
		response := projects.ListProjectsResponse{
			Projects: []*projects.ProjectMembership{},
			Total:    100,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := projectsClient.NewClient(baseClient, "/api/v1/")

	// Create pagination options
	offset := 10
	limit := 20
	order := "displayName"
	opts := &projects.ListProjectsOptions{
		Offset: &offset,
		Limit:  &limit,
		Order:  &order,
	}

	// Call List with options
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := client.List(ctx, opts)

	// Verify result
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}
	if result == nil {
		t.Fatal("List() returned nil result")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d items", len(result))
	}
}

func TestClient_List_EmptyResults(t *testing.T) {
	// Mock HTTP server returning empty list
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := projects.ListProjectsResponse{
			Projects: []*projects.ProjectMembership{},
			Total:    0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := projectsClient.NewClient(baseClient, "/api/v1/")

	// Call List
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := client.List(ctx, nil)

	// Verify result
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("len(result) = %v, want %v", len(result), 0)
	}
}

func TestClient_List_QueryParameterEncoding(t *testing.T) {
	// Mock HTTP server to verify query encoding
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL encoding is correct
		query := r.URL.Query()

		// Check that offset is present
		if !query.Has("offset") {
			t.Error("Expected offset parameter to be present")
		}

		// Check that limit is present
		if !query.Has("limit") {
			t.Error("Expected limit parameter to be present")
		}

		response := projects.ListProjectsResponse{
			Projects: []*projects.ProjectMembership{},
			Total:    0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := projectsClient.NewClient(baseClient, "/api/v1/")

	// Create options with special characters
	offset := 0
	limit := 50
	opts := &projects.ListProjectsOptions{
		Offset: &offset,
		Limit:  &limit,
	}

	// Call List
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := client.List(ctx, opts)

	// Verify no error
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}
}

// ===== Tests for Get() method =====

func TestClient_Get_Success(t *testing.T) {
	projectID := "test-project-123"

	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/project/" + projectID
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Return mock response
		response := projects.GetProjectResponse{
			ProjectID:   projectID,
			DisplayName: "Test Project",
			Description: "Project description",
			Extra:       map[string]interface{}{},
			Namespace:   "test.com",
			Frozen:      false,
			GlobalPermission: &common.Permission{
				ID:    "global-perm-id",
				Label: "DEFAULT",
			},
			UserPermission: &common.Permission{
				ID:    "user-perm-id",
				Label: "ADMIN",
			},
			CreatedAt: "2025-01-01T00:00:00Z",
			UpdatedAt: "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := projectsClient.NewClient(baseClient, "/api/v1/")

	// Call Get
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := client.Get(ctx, projectID)

	// Verify result
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if result == nil {
		t.Fatal("Get() returned nil result")
	}
	if result.ProjectID != projectID {
		t.Errorf("ProjectID = %v, want %v", result.ProjectID, projectID)
	}
	if result.DisplayName != "Test Project" {
		t.Errorf("DisplayName = %v, want %v", result.DisplayName, "Test Project")
	}
	if result.UserPermission == nil || result.UserPermission.Label != "ADMIN" {
		t.Error("UserPermission should have ADMIN label")
	}
}

func TestClient_Get_InvalidProjectID(t *testing.T) {
	// Mock HTTP server returning 400
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"errorCode": 400,
			"message":   "Invalid project ID format",
		}); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := projectsClient.NewClient(baseClient, "/api/v1/")

	// Call Get with invalid ID
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := client.Get(ctx, "invalid-id")

	// Verify error
	if err == nil {
		t.Fatal("Get() should return error for invalid project ID")
	}
	if result != nil {
		t.Error("Get() should return nil result on error")
	}
}

func TestClient_Get_UnauthorizedProject(t *testing.T) {
	// Mock HTTP server returning 403/404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"errorCode": 403,
			"message":   "Access denied to project",
		}); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := projectsClient.NewClient(baseClient, "/api/v1/")

	// Call Get
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := client.Get(ctx, "unauthorized-project")

	// Verify error
	if err == nil {
		t.Fatal("Get() should return error for unauthorized project")
	}
	if result != nil {
		t.Error("Get() should return nil result on error")
	}
}

func TestClient_Get_ContextTimeout(t *testing.T) {
	// Mock HTTP server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client
	httpClient := &http.Client{Timeout: 10 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := projectsClient.NewClient(baseClient, "/api/v1/")

	// Call Get with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result, err := client.Get(ctx, "project-id")

	// Verify timeout error
	if err == nil {
		t.Fatal("Get() should return error for context timeout")
	}
	if result != nil {
		t.Error("Get() should return nil result on timeout")
	}
}
