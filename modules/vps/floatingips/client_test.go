package floatingips

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
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
	if client.baseClient == nil {
		t.Error("expected baseClient to be initialized")
	}
}

// TestClient_List tests the List method returning []*FloatingIP
func TestClient_List(t *testing.T) {
	tests := []struct {
		name           string
		opts           *floatingips.ListFloatingIPsOptions
		mockResponse   *floatingips.FloatingIPListResponse
		mockStatusCode int
		wantErr        bool
		expectedCount  int
		checkPath      string
	}{
		{
			name: "list all floating IPs",
			opts: nil,
			mockResponse: &floatingips.FloatingIPListResponse{
				FloatingIPs: []*floatingips.FloatingIP{
					{ID: "fip-1", Address: "203.0.113.1", Status: floatingips.FloatingIPStatusActive, ProjectID: "proj-123", UserID: "user-123", Reserved: false, CreatedAt: "2025-11-11T10:00:00Z"},
					{ID: "fip-2", Address: "203.0.113.2", Status: floatingips.FloatingIPStatusDown, ProjectID: "proj-123", UserID: "user-123", Reserved: false, CreatedAt: "2025-11-11T10:00:00Z"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			expectedCount:  2,
			checkPath:      "/api/v1/project/proj-123/floatingips",
		},
		{
			name: "list with status filter",
			opts: &floatingips.ListFloatingIPsOptions{Status: "ACTIVE"},
			mockResponse: &floatingips.FloatingIPListResponse{
				FloatingIPs: []*floatingips.FloatingIP{
					{ID: "fip-1", Address: "203.0.113.1", Status: floatingips.FloatingIPStatusActive, ProjectID: "proj-123", UserID: "user-123", Reserved: false, CreatedAt: "2025-11-11T10:00:00Z"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			expectedCount:  1,
			checkPath:      "/api/v1/project/proj-123/floatingips?status=ACTIVE",
		},
		{
			name:           "server error",
			opts:           nil,
			mockResponse:   nil,
			mockStatusCode: http.StatusInternalServerError,
			wantErr:        true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkPath != "" {
					actualPath := r.URL.Path
					if r.URL.RawQuery != "" {
						actualPath += "?" + r.URL.RawQuery
					}
					if actualPath != tt.checkPath {
						t.Errorf("expected path %s, got %s", tt.checkPath, actualPath)
					}
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
			client := NewClient(baseClient, "proj-123")

			result, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil && tt.expectedCount > 0 {
				t.Fatal("expected non-nil result")
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d items, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

// TestClient_Create tests the Create method
func TestClient_Create(t *testing.T) {
	tests := []struct {
		name           string
		request        *floatingips.FloatingIPCreateRequest
		mockResponse   *floatingips.FloatingIP
		mockStatusCode int
		wantErr        bool
	}{
		{
			name: "create floating IP successfully",
			request: &floatingips.FloatingIPCreateRequest{
				Description: "Test FIP",
			},
			mockResponse: &floatingips.FloatingIP{
				ID:          "fip-123",
				Address:     "203.0.113.10",
				Status:      "PENDING",
				ProjectID:   "proj-123",
				Description: "Test FIP",
			},
			mockStatusCode: http.StatusCreated,
			wantErr:        false,
		},
		{
			name:           "create without description",
			request:        &floatingips.FloatingIPCreateRequest{},
			mockResponse:   &floatingips.FloatingIP{ID: "fip-456", Address: "203.0.113.20", Status: "ACTIVE"},
			mockStatusCode: http.StatusCreated,
			wantErr:        false,
		},
		{
			name:           "server error",
			request:        &floatingips.FloatingIPCreateRequest{},
			mockStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}

				expectedPath := "/api/v1/project/proj-123/floatingips"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
			client := NewClient(baseClient, "proj-123")

			result, err := client.Create(context.Background(), tt.request)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Fatal("expected result, got nil")
			}

			if result.ID != tt.mockResponse.ID {
				t.Errorf("expected ID %s, got %s", tt.mockResponse.ID, result.ID)
			}
		})
	}
}

// TestClient_Get tests the Get method
func TestClient_Get(t *testing.T) {
	tests := []struct {
		name           string
		fipID          string
		mockResponse   *floatingips.FloatingIP
		mockStatusCode int
		wantErr        bool
	}{
		{
			name:  "get floating IP successfully",
			fipID: "fip-123",
			mockResponse: &floatingips.FloatingIP{
				ID:      "fip-123",
				Address: "203.0.113.10",
				Status:  "ACTIVE",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "floating IP not found",
			fipID:          "fip-999",
			mockStatusCode: http.StatusNotFound,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/api/v1/project/proj-123/floatingips/" + tt.fipID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
			client := NewClient(baseClient, "proj-123")

			result, err := client.Get(context.Background(), tt.fipID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.ID != tt.mockResponse.ID {
				t.Errorf("expected ID %s, got %s", tt.mockResponse.ID, result.ID)
			}
		})
	}
}

// TestClient_Update tests the Update method
func TestClient_Update(t *testing.T) {
	fipID := "fip-123"
	request := &floatingips.FloatingIPUpdateRequest{
		Description: "Updated description",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/floatingips/" + fipID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&floatingips.FloatingIP{
			ID:          fipID,
			Description: "Updated description",
		})
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	result, err := client.Update(context.Background(), fipID, request)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if result.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got %s", result.Description)
	}
}

// TestClient_Delete tests the Delete method
func TestClient_Delete(t *testing.T) {
	fipID := "fip-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/floatingips/" + fipID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	err := client.Delete(context.Background(), fipID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestClient_Approve tests the Approve method
func TestClient_Approve(t *testing.T) {
	fipID := "fip-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/floatingips/" + fipID + "/approve"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	err := client.Approve(context.Background(), fipID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestClient_Reject tests the Reject method
func TestClient_Reject(t *testing.T) {
	fipID := "fip-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/floatingips/" + fipID + "/reject"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	err := client.Reject(context.Background(), fipID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestClient_Disassociate tests the Disassociate method
func TestClient_Disassociate(t *testing.T) {
	fipID := "fip-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/floatingips/" + fipID + "/disassociate"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	err := client.Disassociate(context.Background(), fipID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
