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

// TestKeypairsCreate_Success verifies successful keypair creation
func TestKeypairsCreate_Success(t *testing.T) {
	mockKeypair := &keypairs.Keypair{
		ID:          "key-new",
		Name:        "new-keypair",
		Description: "Newly created keypair",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
		Fingerprint: "SHA256:xyz789...",
		UserID:      "user-1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/keypairs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Verify request body
		var req keypairs.KeypairCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Name != "new-keypair" {
			t.Errorf("expected name new-keypair, got %s", req.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockKeypair)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()

	ctx := context.Background()
	req := &keypairs.KeypairCreateRequest{
		Name:        "new-keypair",
		Description: "Newly created keypair",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
	}
	keypair, err := keypairsClient.Create(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keypair == nil {
		t.Fatal("expected keypair, got nil")
	}
	if keypair.ID != "key-new" {
		t.Errorf("expected ID key-new, got %s", keypair.ID)
	}
	if keypair.Name != "new-keypair" {
		t.Errorf("expected name new-keypair, got %s", keypair.Name)
	}
}

// TestKeypairsCreate_Generate verifies keypair generation (no public key provided)
func TestKeypairsCreate_Generate(t *testing.T) {
	mockKeypair := &keypairs.Keypair{
		ID:          "key-generated",
		Name:        "generated-key",
		Description: "Auto-generated keypair",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E... (generated)",
		Fingerprint: "SHA256:generated123...",
		UserID:      "user-1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req keypairs.KeypairCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify PublicKey is empty (generation request)
		if req.PublicKey != "" {
			t.Errorf("expected empty PublicKey for generation, got %s", req.PublicKey)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockKeypair)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()

	ctx := context.Background()
	req := &keypairs.KeypairCreateRequest{
		Name:        "generated-key",
		Description: "Auto-generated keypair",
		// PublicKey omitted - should generate
	}
	keypair, err := keypairsClient.Create(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keypair.PublicKey == "" {
		t.Error("expected generated public key, got empty")
	}
}

// TestKeypairsCreate_Errors verifies error handling
func TestKeypairsCreate_Errors(t *testing.T) {
	tests := []struct {
		name         string
		mockStatus   int
		mockResponse interface{}
		expectError  bool
	}{
		{
			name:       "validation error - 400",
			mockStatus: http.StatusBadRequest,
			mockResponse: map[string]interface{}{
				"error":   "Bad Request",
				"message": "Invalid keypair name",
			},
			expectError: true,
		},
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
			name:       "duplicate name - 409",
			mockStatus: http.StatusConflict,
			mockResponse: map[string]interface{}{
				"error":   "Conflict",
				"message": "Keypair with this name already exists",
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
			req := &keypairs.KeypairCreateRequest{Name: "test-key"}
			keypair, err := keypairsClient.Create(ctx, req)

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
