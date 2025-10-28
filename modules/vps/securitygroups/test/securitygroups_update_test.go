package test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

func TestContract_UpdateSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT method, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var req securitygroups.SecurityGroupUpdateRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		if req.Name == nil || *req.Name != "updated-sg" {
			t.Errorf("expected name 'updated-sg', got %v", req.Name)
		}

		response := securitygroups.SecurityGroup{
			ID:          "sg-123",
			Name:        "updated-sg",
			Description: "Updated description",
			ProjectID:   "proj-1",
			UserID:      "user-1",
			Rules:       []securitygroups.SecurityGroupRule{},
			CreatedAt:   "2024-01-01T00:00:00Z",
			UpdatedAt:   "2024-01-02T00:00:00Z",
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	newName := "updated-sg"
	newDesc := "Updated description"
	req := securitygroups.SecurityGroupUpdateRequest{
		Name:        &newName,
		Description: &newDesc,
	}

	result, err := sgClient.Update(context.Background(), "sg-123", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Name != "updated-sg" {
		t.Errorf("expected name 'updated-sg', got %s", result.Name)
	}

	if result.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got %s", result.Description)
	}
}

func TestContract_UpdateSecurityGroup_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "security group not found"})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	newName := "updated-sg"
	req := securitygroups.SecurityGroupUpdateRequest{
		Name: &newName,
	}

	_, err := sgClient.Update(context.Background(), "nonexistent-id", req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
