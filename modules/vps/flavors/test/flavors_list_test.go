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

// TestFlavorsList_Success verifies successful flavor listing with various filters
func TestFlavorsList_Success(t *testing.T) {
	tests := []struct {
		name          string
		opts          *flavors.ListFlavorsOptions
		mockFlavors   []*flavors.Flavor
		expectedQuery map[string]string
	}{
		{
			name: "list all flavors",
			opts: nil,
			mockFlavors: []*flavors.Flavor{
				{
					ID:          "flav-1",
					Name:        "small",
					Description: "Small instance",
					VCPU:        2,
					Memory:      4096,
					Disk:        20,
					Public:      true,
					Tags:        []string{"general"},
				},
				{
					ID:          "flav-2",
					Name:        "medium",
					Description: "Medium instance",
					VCPU:        4,
					Memory:      8192,
					Disk:        40,
					Public:      true,
					Tags:        []string{"general", "balanced"},
				},
			},
			expectedQuery: map[string]string{},
		},
		{
			name: "filter by name",
			opts: &flavors.ListFlavorsOptions{Name: "small"},
			mockFlavors: []*flavors.Flavor{
				{
					ID:          "flav-1",
					Name:        "small",
					Description: "Small instance",
					VCPU:        2,
					Memory:      4096,
					Disk:        20,
					Public:      true,
				},
			},
			expectedQuery: map[string]string{"name": "small"},
		},
		{
			name: "filter by public=true",
			opts: func() *flavors.ListFlavorsOptions {
				publicFlag := true
				return &flavors.ListFlavorsOptions{Public: &publicFlag}
			}(),
			mockFlavors: []*flavors.Flavor{
				{
					ID:     "flav-1",
					Name:   "public-flavor",
					VCPU:   2,
					Memory: 4096,
					Disk:   20,
					Public: true,
				},
			},
			expectedQuery: map[string]string{"public": "true"},
		},
		{
			name: "filter by tag",
			opts: &flavors.ListFlavorsOptions{Tags: []string{"gpu"}},
			mockFlavors: []*flavors.Flavor{
				{
					ID:     "flav-gpu",
					Name:   "gpu-large",
					VCPU:   8,
					Memory: 16384,
					Disk:   100,
					Public: true,
					Tags:   []string{"gpu", "compute"},
				},
			},
			expectedQuery: map[string]string{"tag": "gpu"},
		},
		{
			name:          "empty list",
			opts:          nil,
			mockFlavors:   []*flavors.Flavor{},
			expectedQuery: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				if r.URL.Path != "/vps/api/v1/project/proj-123/flavors" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				// Verify query parameters
				query := r.URL.Query()
				for key, expectedValue := range tt.expectedQuery {
					if got := query.Get(key); got != expectedValue {
						t.Errorf("expected query param %s=%s, got %s", key, expectedValue, got)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := &flavors.FlavorListResponse{Flavors: tt.mockFlavors}
				_ = json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			client := cloudsdk.NewClient(server.URL, "test-token")
			vpsClient := client.Project("proj-123").VPS()
			flavorsClient := vpsClient.Flavors()

			ctx := context.Background()
			resp, err := flavorsClient.List(ctx, tt.opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected response, got nil")
			}
			if len(resp.Flavors) != len(tt.mockFlavors) {
				t.Errorf("expected %d flavors, got %d", len(tt.mockFlavors), len(resp.Flavors))
			}
		})
	}
}

// TestFlavorsList_Errors verifies error handling for flavor listing
func TestFlavorsList_Errors(t *testing.T) {
	tests := []struct {
		name         string
		mockStatus   int
		mockResponse interface{}
		expectError  bool
	}{
		{
			name:       "unauthorized - 401",
			mockStatus: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			},
			expectError: true,
		},
		{
			name:       "forbidden - 403",
			mockStatus: http.StatusForbidden,
			mockResponse: map[string]interface{}{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
			},
			expectError: true,
		},
		{
			name:       "internal server error - 500",
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
			_, err := flavorsClient.List(ctx, nil)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
