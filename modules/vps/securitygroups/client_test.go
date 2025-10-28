package securitygroups

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

// TestNewClient tests the NewClient constructor
func TestNewClient(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	projectID := "proj-123"

	client := NewClient(baseClient, projectID)

	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.projectID != projectID {
		t.Errorf("expected projectID %s, got %s", projectID, client.projectID)
	}
	if client.basePath != "/api/v1/project/proj-123" {
		t.Errorf("expected basePath /api/v1/project/proj-123, got %s", client.basePath)
	}
}

// TestClient_List tests successful security group listing
func TestClient_List(t *testing.T) {
	mockResponse := &securitygroups.SecurityGroupListResponse{
		SecurityGroups: []securitygroups.SecurityGroup{
			{
				ID:          "sg-1",
				Name:        "default",
				Description: "Default security group",
				ProjectID:   "proj-123",
				UserID:      "user-1",
				Rules:       []securitygroups.SecurityGroupRule{},
			},
			{
				ID:          "sg-2",
				Name:        "web",
				Description: "Web security group",
				ProjectID:   "proj-123",
				UserID:      "user-1",
				Rules:       []securitygroups.SecurityGroupRule{},
			},
		},
		Total: 2,
	}

	expectedPath := "/api/v1/project/proj-123/security_groups"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	response, err := client.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response == nil {
		t.Fatal("expected response, got nil")
	}
	if response.Total != 2 {
		t.Errorf("expected total 2, got %d", response.Total)
	}
	if len(response.SecurityGroups) != 2 {
		t.Errorf("expected 2 security groups, got %d", len(response.SecurityGroups))
	}
}

// TestClient_List_WithFilters tests listing with filters
func TestClient_List_WithFilters(t *testing.T) {
	mockResponse := &securitygroups.SecurityGroupListResponse{
		SecurityGroups: []securitygroups.SecurityGroup{
			{
				ID:          "sg-web",
				Name:        "web",
				Description: "Web security group",
				ProjectID:   "proj-123",
				UserID:      "user-1",
				Rules: []securitygroups.SecurityGroupRule{
					{
						ID:         "rule-1",
						Direction:  "ingress",
						Protocol:   "tcp",
						PortMin:    intPtr(80),
						PortMax:    intPtr(80),
						RemoteCIDR: "0.0.0.0/0",
					},
				},
			},
		},
		Total: 1,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		// Verify query parameters
		nameParam := r.URL.Query().Get("name")
		if nameParam != "web" {
			t.Errorf("expected name query param 'web', got %s", nameParam)
		}

		detailParam := r.URL.Query().Get("detail")
		if detailParam != "true" {
			t.Errorf("expected detail query param 'true', got %s", detailParam)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	opts := &securitygroups.ListSecurityGroupsOptions{
		Name:   "web",
		Detail: true,
	}
	response, err := client.List(ctx, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(response.SecurityGroups) != 1 {
		t.Errorf("expected 1 security group, got %d", len(response.SecurityGroups))
	}
	if len(response.SecurityGroups[0].Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(response.SecurityGroups[0].Rules))
	}
}

// TestClient_List_Error tests error handling in List
func TestClient_List_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	_, err := client.List(ctx, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestClient_Create tests successful security group creation
func TestClient_Create(t *testing.T) {
	mockSG := &securitygroups.SecurityGroup{
		ID:          "sg-new",
		Name:        "my-sg",
		Description: "My security group",
		ProjectID:   "proj-123",
		UserID:      "user-1",
		Rules:       []securitygroups.SecurityGroupRule{},
	}

	expectedPath := "/api/v1/project/proj-123/security_groups"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockSG)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	req := securitygroups.SecurityGroupCreateRequest{
		Name:        "my-sg",
		Description: "My security group",
	}
	response, err := client.Create(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response == nil {
		t.Fatal("expected response, got nil")
	}
	if response.ID != "sg-new" {
		t.Errorf("expected ID sg-new, got %s", response.ID)
	}
	if response.Name != "my-sg" {
		t.Errorf("expected name my-sg, got %s", response.Name)
	}
}

// TestClient_Create_Error tests error handling in Create
func TestClient_Create_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	req := securitygroups.SecurityGroupCreateRequest{
		Name: "", // Invalid
	}
	_, err := client.Create(ctx, req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestClient_Get tests successful security group retrieval
func TestClient_Get(t *testing.T) {
	mockSG := &securitygroups.SecurityGroup{
		ID:          "sg-123",
		Name:        "my-sg",
		Description: "My security group",
		ProjectID:   "proj-123",
		UserID:      "user-1",
		Rules: []securitygroups.SecurityGroupRule{
			{
				ID:         "rule-1",
				Direction:  "ingress",
				Protocol:   "tcp",
				PortMin:    intPtr(22),
				PortMax:    intPtr(22),
				RemoteCIDR: "0.0.0.0/0",
			},
		},
	}

	expectedPath := "/api/v1/project/proj-123/security_groups/sg-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockSG)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	response, err := client.Get(ctx, "sg-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response == nil {
		t.Fatal("expected response, got nil")
	}
	if response.ID != "sg-123" {
		t.Errorf("expected ID sg-123, got %s", response.ID)
	}
	if len(response.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(response.Rules))
	}
}

// TestClient_Get_Error tests error handling in Get
func TestClient_Get_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	_, err := client.Get(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestClient_Update tests successful security group update
func TestClient_Update(t *testing.T) {
	mockSG := &securitygroups.SecurityGroup{
		ID:          "sg-123",
		Name:        "updated-sg",
		Description: "Updated description",
		ProjectID:   "proj-123",
		UserID:      "user-1",
		Rules:       []securitygroups.SecurityGroupRule{},
	}

	expectedPath := "/api/v1/project/proj-123/security_groups/sg-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockSG)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	newName := "updated-sg"
	newDesc := "Updated description"
	req := securitygroups.SecurityGroupUpdateRequest{
		Name:        &newName,
		Description: &newDesc,
	}
	response, err := client.Update(ctx, "sg-123", req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response == nil {
		t.Fatal("expected response, got nil")
	}
	if response.Name != "updated-sg" {
		t.Errorf("expected name updated-sg, got %s", response.Name)
	}
}

// TestClient_Update_Error tests error handling in Update
func TestClient_Update_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	newName := "updated-sg"
	req := securitygroups.SecurityGroupUpdateRequest{
		Name: &newName,
	}
	_, err := client.Update(ctx, "nonexistent", req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestClient_Delete tests successful security group deletion
func TestClient_Delete(t *testing.T) {
	expectedPath := "/api/v1/project/proj-123/security_groups/sg-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	err := client.Delete(ctx, "sg-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestClient_Delete_Error tests error handling in Delete
func TestClient_Delete_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	err := client.Delete(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func intPtr(i int) *int {
	return &i
}
