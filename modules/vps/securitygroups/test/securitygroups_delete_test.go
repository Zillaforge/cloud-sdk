package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
)

func TestContract_DeleteSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	err := sgClient.Delete(context.Background(), "sg-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestContract_DeleteSecurityGroup_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "security group not found"})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	err := sgClient.Delete(context.Background(), "nonexistent-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
