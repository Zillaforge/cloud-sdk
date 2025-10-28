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

// TestKeypairsGet_Success verifies successful keypair retrieval
func TestKeypairsGet_Success(t *testing.T) {
	mockKeypair := &keypairs.Keypair{
		ID:          "key-123",
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
		if r.URL.Path != "/vps/api/v1/project/proj-123/keypairs/key-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockKeypair)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()

	ctx := context.Background()
	keypair, err := keypairsClient.Get(ctx, "key-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keypair == nil {
		t.Fatal("expected keypair, got nil")
	}
	if keypair.ID != "key-123" {
		t.Errorf("expected ID key-123, got %s", keypair.ID)
	}
	if keypair.Name != "my-keypair" {
		t.Errorf("expected name my-keypair, got %s", keypair.Name)
	}
}

// TestKeypairsGet_Errors verifies error handling
func TestKeypairsGet_Errors(t *testing.T) {
	tests := []struct {
		name         string
		keypairID    string
		mockStatus   int
		mockResponse interface{}
		expectError  bool
	}{
		{
			name:       "not found - 404",
			keypairID:  "key-nonexistent",
			mockStatus: http.StatusNotFound,
			mockResponse: map[string]interface{}{
				"error":   "Not Found",
				"message": "Keypair not found",
			},
			expectError: true,
		},
		{
			name:       "unauthorized - 401",
			keypairID:  "key-123",
			mockStatus: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			},
			expectError: true,
		},
		{
			name:       "forbidden - 403",
			keypairID:  "key-123",
			mockStatus: http.StatusForbidden,
			mockResponse: map[string]interface{}{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
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
			keypair, err := keypairsClient.Get(ctx, tt.keypairID)

			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectError && keypair != nil {
				t.Errorf("expected nil keypair on error, got %v", keypair)
			}
		})
	}
}
