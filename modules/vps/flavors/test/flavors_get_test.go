package flavors_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/flavors"
)

// TestFlavorsGet_Success verifies successful flavor retrieval
func TestFlavorsGet_Success(t *testing.T) {
	mockFlavor := &flavors.Flavor{
		ID:          "flav-123",
		Name:        "large",
		Description: "Large compute instance",
		VCPUs:       8,
		RAM:         16384,
		Disk:        80,
		Public:      true,
		Tags:        []string{"compute", "balanced"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/flavors/flav-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockFlavor)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	flavorsClient := vpsClient.Flavors()

	ctx := context.Background()
	flavor, err := flavorsClient.Get(ctx, "flav-123")

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
	if flavor.VCPUs != mockFlavor.VCPUs {
		t.Errorf("expected %d VCPUs, got %d", mockFlavor.VCPUs, flavor.VCPUs)
	}
	if flavor.RAM != mockFlavor.RAM {
		t.Errorf("expected %d RAM, got %d", mockFlavor.RAM, flavor.RAM)
	}
	if flavor.Disk != mockFlavor.Disk {
		t.Errorf("expected %d Disk, got %d", mockFlavor.Disk, flavor.Disk)
	}
}

// TestFlavorsGet_Errors verifies error handling for flavor retrieval
func TestFlavorsGet_Errors(t *testing.T) {
	tests := []struct {
		name         string
		flavorID     string
		mockStatus   int
		mockResponse interface{}
		expectError  bool
	}{
		{
			name:       "unauthorized - 401",
			flavorID:   "flav-123",
			mockStatus: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			},
			expectError: true,
		},
		{
			name:       "forbidden - 403",
			flavorID:   "flav-123",
			mockStatus: http.StatusForbidden,
			mockResponse: map[string]interface{}{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
			},
			expectError: true,
		},
		{
			name:       "not found - 404",
			flavorID:   "nonexistent",
			mockStatus: http.StatusNotFound,
			mockResponse: map[string]interface{}{
				"error":   "Not Found",
				"message": "Flavor not found",
			},
			expectError: true,
		},
		{
			name:       "internal server error - 500",
			flavorID:   "flav-123",
			mockStatus: http.StatusInternalServerError,
			mockResponse: map[string]interface{}{
				"error":   "Internal Server Error",
				"message": "Something went wrong",
			},
			expectError: true,
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

			client := cloudsdk.NewClient(server.URL, "test-token")
			vpsClient := client.Project("proj-123").VPS()
			flavorsClient := vpsClient.Flavors()

			ctx := context.Background()
			_, err := flavorsClient.Get(ctx, tt.flavorID)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
