package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/keypairs"
)

// TestKeypairsList_Success verifies successful keypair listing
func TestKeypairsList_Success(t *testing.T) {
	mockKeypairs := []keypairs.Keypair{
		{
			ID:          "key-1",
			Name:        "my-keypair",
			Description: "Primary SSH key",
			PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
			Fingerprint: "SHA256:abc123...",
			UserID:      "user-1",
		},
		{
			ID:          "key-2",
			Name:        "backup-key",
			Description: "Backup SSH key",
			PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
			Fingerprint: "SHA256:def456...",
			UserID:      "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/keypairs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := &keypairs.KeypairListResponse{Keypairs: mockKeypairs}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()

	ctx := context.Background()
	result, err := keypairsClient.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keypairs, got %d", len(result))
	}
	if result[0].Name != "my-keypair" {
		t.Errorf("expected name my-keypair, got %s", result[0].Name)
	}
}

// TestKeypairsList_WithFilter verifies listing with name filter
func TestKeypairsList_WithFilter(t *testing.T) {
	mockKeypair := keypairs.Keypair{
		ID:          "key-1",
		Name:        "my-keypair",
		Description: "Primary SSH key",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
		Fingerprint: "SHA256:abc123...",
		UserID:      "user-1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		// Verify query parameter
		if name := r.URL.Query().Get("name"); name != "my-keypair" {
			t.Errorf("expected name=my-keypair, got %s", name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := &keypairs.KeypairListResponse{Keypairs: []keypairs.Keypair{mockKeypair}}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()

	ctx := context.Background()
	opts := &keypairs.ListKeypairsOptions{Name: "my-keypair"}
	result, err := keypairsClient.List(ctx, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 keypair, got %d", len(result))
	}
}

// TestKeypairsList_Empty verifies empty list response
func TestKeypairsList_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := &keypairs.KeypairListResponse{Keypairs: []keypairs.Keypair{}}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()

	ctx := context.Background()
	result, err := keypairsClient.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 keypairs, got %d", len(result))
	}
}

// TestKeypairsList_Errors verifies error handling
func TestKeypairsList_Errors(t *testing.T) {
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
				"message": "An unexpected error occurred",
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
			keypairsClient := vpsClient.Keypairs()

			ctx := context.Background()
			_, err := keypairsClient.List(ctx, nil)

			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
