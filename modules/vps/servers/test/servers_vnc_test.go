package servers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
)

// TestServersVNCURL_Success verifies successful VNC URL retrieval
func TestServersVNCURL_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"url": "https://console.example.com/vnc?token=abc123",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/servers/svr-1/vnc_url" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client, err := cloudsdk.New(mockServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	vpsClient := client.Project("proj-123").VPS()

	vncResp, err := vpsClient.Servers().GetVNCConsoleURL(context.Background(), "svr-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedURL := "https://console.example.com/vnc?token=abc123"
	if vncResp.URL != expectedURL {
		t.Errorf("expected URL '%s', got '%s'", expectedURL, vncResp.URL)
	}
}

// TestServersVNCURL_Errors verifies error handling for VNC URL retrieval
func TestServersVNCURL_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"unauthorized - 401", http.StatusUnauthorized},
		{"forbidden - 403", http.StatusForbidden},
		{"not found - 404", http.StatusNotFound},
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

			_, err = vpsClient.Servers().GetVNCConsoleURL(context.Background(), "svr-1")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
