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

// TestKeypairsUpdate_Success verifies successful keypair update
func TestKeypairsUpdate_Success(t *testing.T) {
	mockKeypair := &keypairs.Keypair{
		ID:          "key-123",
		Name:        "my-keypair",
		Description: "Updated description",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
		Fingerprint: "SHA256:abc123...",
		UserID:      "user-1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			t.Errorf("expected PUT or PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/keypairs/key-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Verify request body
		var req keypairs.KeypairUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Description != "Updated description" {
			t.Errorf("expected description 'Updated description', got %s", req.Description)
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
	req := &keypairs.KeypairUpdateRequest{
		Description: "Updated description",
	}
	keypair, err := keypairsClient.Update(ctx, "key-123", req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keypair == nil {
		t.Fatal("expected keypair, got nil")
	}
	if keypair.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got %s", keypair.Description)
	}
}

// TestKeypairsUpdate_Errors verifies error handling
func TestKeypairsUpdate_Errors(t *testing.T) {
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
			name:       "validation error - 400",
			keypairID:  "key-123",
			mockStatus: http.StatusBadRequest,
			mockResponse: map[string]interface{}{
				"error":   "Bad Request",
				"message": "Invalid update request",
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
			req := &keypairs.KeypairUpdateRequest{Description: "test"}
			keypair, err := keypairsClient.Update(ctx, tt.keypairID, req)

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
