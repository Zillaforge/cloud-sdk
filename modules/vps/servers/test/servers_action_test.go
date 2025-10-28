package servers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

// TestServersAction_Success verifies successful server actions
func TestServersAction_Success(t *testing.T) {
	tests := []struct {
		name   string
		action string
	}{
		{"start action", "start"},
		{"stop action", "stop"},
		{"reboot action", "reboot"},
		{"resize action", "resize"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/vps/api/v1/project/proj-123/servers/svr-1/action" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				w.WriteHeader(http.StatusAccepted)
			}))
			defer mockServer.Close()

			client, err := cloudsdk.New(mockServer.URL, "test-token")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			vpsClient := client.Project("proj-123").VPS()

			req := &servers.ServerActionRequest{
				Action: tt.action,
			}

			err = vpsClient.Servers().Action(context.Background(), "svr-1", req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestServersAction_Errors verifies error handling for server actions
func TestServersAction_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"invalid action - 400", http.StatusBadRequest},
		{"unauthorized - 401", http.StatusUnauthorized},
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

			req := &servers.ServerActionRequest{Action: "start"}
			err = vpsClient.Servers().Action(context.Background(), "svr-1", req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
