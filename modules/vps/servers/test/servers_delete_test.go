package servers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
)

// TestServersDelete_Success verifies successful server deletion
func TestServersDelete_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/servers/svr-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer mockServer.Close()

	client, err := cloudsdk.New(mockServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	vpsClient := client.Project("proj-123").VPS()

	err = vpsClient.Servers().Delete(context.Background(), "svr-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestServersDelete_Errors verifies error handling for server deletion
func TestServersDelete_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"not found - 404", http.StatusNotFound},
		{"conflict - 409", http.StatusConflict},
		{"unauthorized - 401", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"message":    "Error",
					"error_code": tt.statusCode,
				})
			}))
			defer mockServer.Close()

			client, err := cloudsdk.New(mockServer.URL, "test-token")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			vpsClient := client.Project("proj-123").VPS()

			err = vpsClient.Servers().Delete(context.Background(), "svr-1")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
