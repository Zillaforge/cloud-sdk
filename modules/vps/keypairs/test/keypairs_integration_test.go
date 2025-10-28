package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/keypairs"
)

// TestKeypairsIntegration_FullLifecycle tests complete keypair lifecycle
func TestKeypairsIntegration_FullLifecycle(t *testing.T) {
	var createdKeypairID string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// CREATE
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/keypairs") {
			var req keypairs.KeypairCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode create request: %v", err)
			}

			createdKeypairID = "key-integration-test"
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(&keypairs.Keypair{
				ID:          createdKeypairID,
				Name:        req.Name,
				Description: req.Description,
				PublicKey:   req.PublicKey,
				Fingerprint: "SHA256:integration-test-fingerprint",
				UserID:      "user-1",
			})
			return
		}

		// GET
		if r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/keypairs/key-integration-test") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(&keypairs.Keypair{
				ID:          "key-integration-test",
				Name:        "integration-test-keypair",
				Description: "Integration test SSH key",
				PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
				Fingerprint: "SHA256:integration-test-fingerprint",
				UserID:      "user-1",
			})
			return
		}

		// UPDATE
		if (r.Method == http.MethodPut || r.Method == http.MethodPatch) && strings.Contains(r.URL.Path, "/keypairs/key-integration-test") {
			var req keypairs.KeypairUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode update request: %v", err)
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(&keypairs.Keypair{
				ID:          "key-integration-test",
				Name:        "integration-test-keypair",
				Description: req.Description,
				PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
				Fingerprint: "SHA256:integration-test-fingerprint",
				UserID:      "user-1",
			})
			return
		}

		// DELETE
		if r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/keypairs/key-integration-test") {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Unexpected request
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()
	ctx := context.Background()

	// Step 1: Create keypair
	t.Run("Step1_Create", func(t *testing.T) {
		createReq := &keypairs.KeypairCreateRequest{
			Name:        "integration-test-keypair",
			Description: "Integration test SSH key",
			PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
		}
		keypair, err := keypairsClient.Create(ctx, createReq)
		if err != nil {
			t.Fatalf("failed to create keypair: %v", err)
		}
		if keypair.ID != "key-integration-test" {
			t.Errorf("expected ID key-integration-test, got %s", keypair.ID)
		}
		if keypair.Name != "integration-test-keypair" {
			t.Errorf("expected name integration-test-keypair, got %s", keypair.Name)
		}
		createdKeypairID = keypair.ID
	})

	// Step 2: Get keypair
	t.Run("Step2_Get", func(t *testing.T) {
		keypair, err := keypairsClient.Get(ctx, createdKeypairID)
		if err != nil {
			t.Fatalf("failed to get keypair: %v", err)
		}
		if keypair.ID != createdKeypairID {
			t.Errorf("expected ID %s, got %s", createdKeypairID, keypair.ID)
		}
		if keypair.Fingerprint != "SHA256:integration-test-fingerprint" {
			t.Errorf("expected fingerprint SHA256:integration-test-fingerprint, got %s", keypair.Fingerprint)
		}
	})

	// Step 3: Update keypair
	t.Run("Step3_Update", func(t *testing.T) {
		updateReq := &keypairs.KeypairUpdateRequest{
			Description: "Updated integration test description",
		}
		keypair, err := keypairsClient.Update(ctx, createdKeypairID, updateReq)
		if err != nil {
			t.Fatalf("failed to update keypair: %v", err)
		}
		if keypair.Description != "Updated integration test description" {
			t.Errorf("expected updated description, got %s", keypair.Description)
		}
	})

	// Step 4: Delete keypair
	t.Run("Step4_Delete", func(t *testing.T) {
		err := keypairsClient.Delete(ctx, createdKeypairID)
		if err != nil {
			t.Fatalf("failed to delete keypair: %v", err)
		}
	})
}

// TestKeypairsIntegration_ListAndFilter tests listing and filtering keypairs
func TestKeypairsIntegration_ListAndFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || !strings.HasSuffix(r.URL.Path, "/keypairs") {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check for filter parameter
		nameFilter := r.URL.Query().Get("name")

		var response keypairs.KeypairListResponse
		if nameFilter == "" {
			// Return all keypairs
			response = keypairs.KeypairListResponse{
				Keypairs: []keypairs.Keypair{
					{ID: "key-1", Name: "keypair-one", Description: "First keypair"},
					{ID: "key-2", Name: "keypair-two", Description: "Second keypair"},
					{ID: "key-3", Name: "special-keypair", Description: "Special keypair"},
				},
				Total: 3,
			}
		} else if nameFilter == "special" {
			// Return filtered keypairs
			response = keypairs.KeypairListResponse{
				Keypairs: []keypairs.Keypair{
					{ID: "key-3", Name: "special-keypair", Description: "Special keypair"},
				},
				Total: 1,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()
	ctx := context.Background()

	// Test 1: List all keypairs
	t.Run("ListAll", func(t *testing.T) {
		response, err := keypairsClient.List(ctx, nil)
		if err != nil {
			t.Fatalf("failed to list keypairs: %v", err)
		}
		if response.Total != 3 {
			t.Errorf("expected 3 keypairs, got %d", response.Total)
		}
		if len(response.Keypairs) != 3 {
			t.Errorf("expected 3 keypairs in list, got %d", len(response.Keypairs))
		}
	})

	// Test 2: List with filter
	t.Run("ListWithFilter", func(t *testing.T) {
		opts := &keypairs.ListKeypairsOptions{
			Name: "special",
		}
		response, err := keypairsClient.List(ctx, opts)
		if err != nil {
			t.Fatalf("failed to list filtered keypairs: %v", err)
		}
		if response.Total != 1 {
			t.Errorf("expected 1 keypair, got %d", response.Total)
		}
		if len(response.Keypairs) != 1 {
			t.Errorf("expected 1 keypair in list, got %d", len(response.Keypairs))
		}
		if response.Keypairs[0].Name != "special-keypair" {
			t.Errorf("expected name special-keypair, got %s", response.Keypairs[0].Name)
		}
	})
}

// TestKeypairsIntegration_CreateWithGeneration tests creating keypair with auto-generation
func TestKeypairsIntegration_CreateWithGeneration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/keypairs") {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var req keypairs.KeypairCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode create request: %v", err)
		}

		// When no PublicKey provided, server generates one
		generatedKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADA...GENERATED"
		privateKey := ""
		if req.PublicKey == "" {
			privateKey = "-----BEGIN RSA PRIVATE KEY-----\nGenerated private key...\n-----END RSA PRIVATE KEY-----"
		}

		response := &keypairs.Keypair{
			ID:          "key-generated",
			Name:        req.Name,
			Description: req.Description,
			PublicKey:   generatedKey,
			Fingerprint: "SHA256:generated-fingerprint",
			UserID:      "user-1",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// In real API, private key would be returned in a special field
		// For this test, we'll add it to the response
		responseMap := map[string]interface{}{
			"id":          response.ID,
			"name":        response.Name,
			"description": response.Description,
			"public_key":  response.PublicKey,
			"fingerprint": response.Fingerprint,
			"user_id":     response.UserID,
		}
		if privateKey != "" {
			responseMap["private_key"] = privateKey
		}

		_ = json.NewEncoder(w).Encode(responseMap)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()
	ctx := context.Background()

	// Create keypair without providing public key
	createReq := &keypairs.KeypairCreateRequest{
		Name:        "auto-generated-keypair",
		Description: "Keypair with auto-generated keys",
		// PublicKey is intentionally omitted
	}

	keypair, err := keypairsClient.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("failed to create keypair: %v", err)
	}
	if keypair.ID != "key-generated" {
		t.Errorf("expected ID key-generated, got %s", keypair.ID)
	}
	if keypair.PublicKey == "" {
		t.Errorf("expected generated public key, got empty string")
	}
	if keypair.Fingerprint != "SHA256:generated-fingerprint" {
		t.Errorf("expected fingerprint SHA256:generated-fingerprint, got %s", keypair.Fingerprint)
	}
}

// TestKeypairsIntegration_ErrorScenarios tests various error scenarios
func TestKeypairsIntegration_ErrorScenarios(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Simulate duplicate name error on create
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/keypairs") {
			var req keypairs.KeypairCreateRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req.Name == "existing-keypair" {
				w.WriteHeader(http.StatusConflict)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error":   "Conflict",
					"message": "Keypair with this name already exists",
				})
				return
			}
		}

		// Simulate not found error
		if strings.Contains(r.URL.Path, "nonexistent") {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Not Found",
				"message": "Keypair not found",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	keypairsClient := vpsClient.Keypairs()
	ctx := context.Background()

	// Test 1: Duplicate name on create
	t.Run("DuplicateName", func(t *testing.T) {
		createReq := &keypairs.KeypairCreateRequest{
			Name:      "existing-keypair",
			PublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
		}
		_, err := keypairsClient.Create(ctx, createReq)
		if err == nil {
			t.Errorf("expected error for duplicate name, got nil")
		}
	})

	// Test 2: Get nonexistent keypair
	t.Run("GetNonexistent", func(t *testing.T) {
		_, err := keypairsClient.Get(ctx, "nonexistent")
		if err == nil {
			t.Errorf("expected error for nonexistent keypair, got nil")
		}
	})

	// Test 3: Update nonexistent keypair
	t.Run("UpdateNonexistent", func(t *testing.T) {
		updateReq := &keypairs.KeypairUpdateRequest{
			Description: "Test",
		}
		_, err := keypairsClient.Update(ctx, "nonexistent", updateReq)
		if err == nil {
			t.Errorf("expected error for nonexistent keypair, got nil")
		}
	})

	// Test 4: Delete nonexistent keypair
	t.Run("DeleteNonexistent", func(t *testing.T) {
		err := keypairsClient.Delete(ctx, "nonexistent")
		if err == nil {
			t.Errorf("expected error for nonexistent keypair, got nil")
		}
	})
}
