package flavors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/flavors"
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

// TestClient_List tests flavor listing with various options
func TestClient_List(t *testing.T) {
	tests := []struct {
		name         string
		opts         *flavors.ListFlavorsOptions
		mockFlavors  []*flavors.Flavor
		expectedPath string
	}{
		{
			name: "list all flavors",
			opts: nil,
			mockFlavors: []*flavors.Flavor{
				{ID: "flav-1", Name: "small", VCPU: 2, Memory: 4096, Disk: 20, Public: true},
				{ID: "flav-2", Name: "medium", VCPU: 4, Memory: 8192, Disk: 40, Public: true},
			},
			expectedPath: "/api/v1/project/proj-123/flavors",
		},
		{
			name: "filter by name",
			opts: &flavors.ListFlavorsOptions{Name: "small"},
			mockFlavors: []*flavors.Flavor{
				{ID: "flav-1", Name: "small", VCPU: 2, Memory: 4096, Disk: 20, Public: true},
			},
			expectedPath: "/api/v1/project/proj-123/flavors",
		},
		{
			name: "filter by public=true",
			opts: func() *flavors.ListFlavorsOptions {
				publicFlag := true
				return &flavors.ListFlavorsOptions{Public: &publicFlag}
			}(),
			mockFlavors: []*flavors.Flavor{
				{ID: "flav-1", Name: "public-flavor", VCPU: 2, Memory: 4096, Disk: 20, Public: true},
			},
			expectedPath: "/api/v1/project/proj-123/flavors",
		},
		{
			name: "filter by tag",
			opts: &flavors.ListFlavorsOptions{Tags: []string{"gpu"}},
			mockFlavors: []*flavors.Flavor{
				{ID: "flav-gpu", Name: "gpu-large", VCPU: 8, Memory: 16384, Disk: 100, Public: true, Tags: []string{"gpu"}},
			},
			expectedPath: "/api/v1/project/proj-123/flavors",
		},
		{
			name:         "empty list",
			opts:         nil,
			mockFlavors:  []*flavors.Flavor{},
			expectedPath: "/api/v1/project/proj-123/flavors",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				if r.URL.Path != tt.expectedPath {
					t.Errorf("expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				// Verify query parameters if opts provided
				if tt.opts != nil {
					query := r.URL.Query()
					if tt.opts.Name != "" && query.Get("name") != tt.opts.Name {
						t.Errorf("expected name=%s, got %s", tt.opts.Name, query.Get("name"))
					}
					if tt.opts.Public != nil && query.Get("public") != "" {
						expected := "true"
						if !*tt.opts.Public {
							expected = "false"
						}
						if query.Get("public") != expected {
							t.Errorf("expected public=%s, got %s", expected, query.Get("public"))
						}
					}
					if len(tt.opts.Tags) > 0 {
						tags := query["tag"]
						if len(tags) != len(tt.opts.Tags) {
							t.Errorf("expected %d tags, got %d", len(tt.opts.Tags), len(tags))
						}
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := &flavors.FlavorListResponse{Flavors: tt.mockFlavors}
				_ = json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			resp, err := client.List(ctx, tt.opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected response, got nil")
			}
			if len(resp) != len(tt.mockFlavors) {
				t.Errorf("expected %d flavors, got %d", len(tt.mockFlavors), len(resp))
			}
		})
	}
}

// TestClient_List_Errors tests error handling for flavor listing
func TestClient_List_Errors(t *testing.T) {
	tests := []struct {
		name         string
		mockStatus   int
		mockResponse map[string]interface{}
	}{
		{
			name:       "unauthorized",
			mockStatus: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"error":   "Unauthorized",
				"message": "Invalid token",
			},
		},
		{
			name:       "internal server error",
			mockStatus: http.StatusInternalServerError,
			mockResponse: map[string]interface{}{
				"error":   "Internal Server Error",
				"message": "Something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{}, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			_, err := client.List(ctx, nil)

			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

// TestClient_Get tests retrieving a specific flavor
func TestClient_Get(t *testing.T) {
	mockFlavor := &flavors.Flavor{
		ID:          "flav-123",
		Name:        "large",
		Description: "Large compute instance",
		VCPU:        8,
		Memory:      16384,
		Disk:        80,
		Public:      true,
		Tags:        []string{"compute", "balanced"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/project/proj-123/flavors/flav-123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockFlavor)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	flavor, err := client.Get(ctx, "flav-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flavor == nil {
		t.Fatal("expected flavor, got nil")
	}
	if flavor.ID != mockFlavor.ID {
		t.Errorf("expected flavor ID %s, got %s", mockFlavor.ID, flavor.ID)
	}
	if flavor.Name != mockFlavor.Name {
		t.Errorf("expected flavor name %s, got %s", mockFlavor.Name, flavor.Name)
	}
	if flavor.VCPU != mockFlavor.VCPU {
		t.Errorf("expected %d VCPU, got %d", mockFlavor.VCPU, flavor.VCPU)
	}
}

// TestClient_Get_Errors tests error handling for flavor retrieval
func TestClient_Get_Errors(t *testing.T) {
	tests := []struct {
		name         string
		flavorID     string
		mockStatus   int
		mockResponse map[string]interface{}
	}{
		{
			name:       "not found",
			flavorID:   "nonexistent",
			mockStatus: http.StatusNotFound,
			mockResponse: map[string]interface{}{
				"error":   "Not Found",
				"message": "Flavor not found",
			},
		},
		{
			name:       "unauthorized",
			flavorID:   "flav-123",
			mockStatus: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"error":   "Unauthorized",
				"message": "Invalid token",
			},
		},
		{
			name:       "internal server error",
			flavorID:   "flav-123",
			mockStatus: http.StatusInternalServerError,
			mockResponse: map[string]interface{}{
				"error":   "Internal Server Error",
				"message": "Something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{}, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			_, err := client.Get(ctx, tt.flavorID)

			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

// TestContractFlavorAPIResponse tests that API responses match the expected contract.
func TestContractFlavorAPIResponse(t *testing.T) {
	// This test validates that the API response structure matches our Flavor model
	// as defined in the vps.yaml swagger spec (pb.FlavorInfo)

	tests := []struct {
		name         string
		responseJSON string
		validate     func(*testing.T, *flavors.Flavor)
	}{
		{
			name: "complete flavor with all fields",
			responseJSON: `{
				"id": "flavor-123",
				"name": "m1.xlarge",
				"description": "Extra large flavor with GPU",
				"vcpu": 16,
				"memory": 32768,
				"disk": 200,
				"gpu": {
					"count": 2,
					"is_vgpu": false,
					"model": "NVIDIA Tesla V100"
				},
				"public": true,
				"tags": ["gpu", "compute", "large"],
				"project_ids": ["proj-1", "proj-2"],
				"az": "az-east-1",
				"createdAt": "2025-11-08T10:00:00Z",
				"updatedAt": "2025-11-08T12:00:00Z"
			}`,
			validate: func(t *testing.T, f *flavors.Flavor) {
				if f.ID != "flavor-123" {
					t.Errorf("ID: expected 'flavor-123', got '%s'", f.ID)
				}
				if f.Name != "m1.xlarge" {
					t.Errorf("Name: expected 'm1.xlarge', got '%s'", f.Name)
				}
				if f.VCPU != 16 {
					t.Errorf("VCPU: expected 16, got %d", f.VCPU)
				}
				if f.Memory != 32768 {
					t.Errorf("Memory: expected 32768, got %d", f.Memory)
				}
				if f.Disk != 200 {
					t.Errorf("Disk: expected 200, got %d", f.Disk)
				}
				if f.GPU == nil {
					t.Fatal("GPU: expected non-nil, got nil")
				}
				if f.GPU.Count != 2 {
					t.Errorf("GPU.Count: expected 2, got %d", f.GPU.Count)
				}
				if f.GPU.Model != "NVIDIA Tesla V100" {
					t.Errorf("GPU.Model: expected 'NVIDIA Tesla V100', got '%s'", f.GPU.Model)
				}
				if len(f.Tags) != 3 {
					t.Errorf("Tags: expected 3 tags, got %d", len(f.Tags))
				}
				if len(f.ProjectIDs) != 2 {
					t.Errorf("ProjectIDs: expected 2 project IDs, got %d", len(f.ProjectIDs))
				}
				if f.AZ != "az-east-1" {
					t.Errorf("AZ: expected 'az-east-1', got '%s'", f.AZ)
				}
				if f.CreatedAt == nil {
					t.Error("CreatedAt: expected non-nil, got nil")
				}
				if f.UpdatedAt == nil {
					t.Error("UpdatedAt: expected non-nil, got nil")
				}
			},
		},
		{
			name: "minimal flavor without optional fields",
			responseJSON: `{
				"id": "flavor-456",
				"name": "m1.small",
				"vcpu": 2,
				"memory": 4096,
				"disk": 20,
				"public": true
			}`,
			validate: func(t *testing.T, f *flavors.Flavor) {
				if f.ID != "flavor-456" {
					t.Errorf("ID: expected 'flavor-456', got '%s'", f.ID)
				}
				if f.Name != "m1.small" {
					t.Errorf("Name: expected 'm1.small', got '%s'", f.Name)
				}
				if f.VCPU != 2 {
					t.Errorf("VCPU: expected 2, got %d", f.VCPU)
				}
				if f.Memory != 4096 {
					t.Errorf("Memory: expected 4096, got %d", f.Memory)
				}
				if f.Disk != 20 {
					t.Errorf("Disk: expected 20, got %d", f.Disk)
				}
				if !f.Public {
					t.Error("Public: expected true, got false")
				}
				if f.GPU != nil {
					t.Errorf("GPU: expected nil, got %+v", f.GPU)
				}
				if f.Description != "" {
					t.Errorf("Description: expected empty, got '%s'", f.Description)
				}
				if len(f.Tags) != 0 {
					t.Errorf("Tags: expected 0 tags, got %d", len(f.Tags))
				}
			},
		},
		{
			name: "flavor with vGPU",
			responseJSON: `{
				"id": "flavor-vgpu",
				"name": "m1.vgpu",
				"vcpu": 8,
				"memory": 16384,
				"disk": 100,
				"gpu": {
					"count": 1,
					"is_vgpu": true,
					"model": "NVIDIA GRID K2"
				},
				"public": true
			}`,
			validate: func(t *testing.T, f *flavors.Flavor) {
				if f.GPU == nil {
					t.Fatal("GPU: expected non-nil, got nil")
				}
				if f.GPU.Count != 1 {
					t.Errorf("GPU.Count: expected 1, got %d", f.GPU.Count)
				}
				if !f.GPU.IsVGPU {
					t.Error("GPU.IsVGPU: expected true, got false")
				}
				if f.GPU.Model != "NVIDIA GRID K2" {
					t.Errorf("GPU.Model: expected 'NVIDIA GRID K2', got '%s'", f.GPU.Model)
				}
			},
		},
		{
			name: "flavor with deletion timestamp",
			responseJSON: `{
				"id": "flavor-deleted",
				"name": "m1.deprecated",
				"vcpu": 4,
				"memory": 8192,
				"disk": 40,
				"public": false,
				"createdAt": "2025-10-01T10:00:00Z",
				"updatedAt": "2025-11-01T10:00:00Z",
				"deletedAt": "2025-11-05T10:00:00Z"
			}`,
			validate: func(t *testing.T, f *flavors.Flavor) {
				if f.DeletedAt == nil {
					t.Error("DeletedAt: expected non-nil, got nil")
				}
				if f.CreatedAt == nil {
					t.Error("CreatedAt: expected non-nil, got nil")
				}
				if f.UpdatedAt == nil {
					t.Error("UpdatedAt: expected non-nil, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.responseJSON))
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			flavor, err := client.Get(ctx, "test-flavor")

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if flavor == nil {
				t.Fatal("expected flavor, got nil")
			}

			tt.validate(t, flavor)
		})
	}
}

// TestContractFlavorListResponse tests that list API responses match the expected contract.
func TestContractFlavorListResponse(t *testing.T) {
	responseJSON := `{
		"flavors": [
			{
				"id": "flavor-1",
				"name": "m1.small",
				"vcpu": 2,
				"memory": 4096,
				"disk": 20,
				"public": true
			},
			{
				"id": "flavor-2",
				"name": "m1.large",
				"description": "Large with GPU",
				"vcpu": 8,
				"memory": 16384,
				"disk": 100,
				"gpu": {
					"count": 1,
					"is_vgpu": false,
					"model": "NVIDIA T4"
				},
				"public": true,
				"tags": ["gpu"]
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseJSON))
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
	if len(response) != 2 {
		t.Fatalf("expected 2 flavors, got %d", len(response))
	}

	// Validate first flavor (minimal)
	f1 := response[0]
	if f1.ID != "flavor-1" {
		t.Errorf("Flavor 1 ID: expected 'flavor-1', got '%s'", f1.ID)
	}
	if f1.VCPU != 2 {
		t.Errorf("Flavor 1 VCPU: expected 2, got %d", f1.VCPU)
	}
	if f1.Memory != 4096 {
		t.Errorf("Flavor 1 Memory: expected 4096, got %d", f1.Memory)
	}

	// Validate second flavor (with GPU)
	f2 := response[1]
	if f2.ID != "flavor-2" {
		t.Errorf("Flavor 2 ID: expected 'flavor-2', got '%s'", f2.ID)
	}
	if f2.GPU == nil {
		t.Fatal("Flavor 2 GPU: expected non-nil, got nil")
	}
	if f2.GPU.Count != 1 {
		t.Errorf("Flavor 2 GPU.Count: expected 1, got %d", f2.GPU.Count)
	}
	if f2.GPU.Model != "NVIDIA T4" {
		t.Errorf("Flavor 2 GPU.Model: expected 'NVIDIA T4', got '%s'", f2.GPU.Model)
	}
}

// TestContractMultipleTagsFiltering tests that multiple tags can be passed as query parameters.
func TestContractMultipleTagsFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify multiple tag parameters are sent correctly
		tags := r.URL.Query()["tag"]
		if len(tags) != 2 {
			t.Errorf("expected 2 tag parameters, got %d", len(tags))
		}
		if tags[0] != "gpu" || tags[1] != "large" {
			t.Errorf("expected tags [gpu, large], got %v", tags)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := &flavors.FlavorListResponse{
			Flavors: []*flavors.Flavor{
				{
					ID:     "flavor-gpu-large",
					Name:   "m1.gpu-large",
					VCPU:   16,
					Memory: 32768,
					Disk:   200,
					Public: true,
					Tags:   []string{"gpu", "large"},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	opts := &flavors.ListFlavorsOptions{
		Tags: []string{"gpu", "large"},
	}
	_, err := client.List(ctx, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestContractResizeServerIDFilter tests the resize_server_id query parameter.
func TestContractResizeServerIDFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resizeServerID := r.URL.Query().Get("resize_server_id")
		if resizeServerID != "server-123" {
			t.Errorf("expected resize_server_id=server-123, got %s", resizeServerID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := &flavors.FlavorListResponse{
			Flavors: []*flavors.Flavor{
				{ID: "flavor-1", Name: "compatible", VCPU: 4, Memory: 8192, Disk: 40, Public: true},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	opts := &flavors.ListFlavorsOptions{
		ResizeServerID: "server-123",
	}
	_, err := client.List(ctx, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
