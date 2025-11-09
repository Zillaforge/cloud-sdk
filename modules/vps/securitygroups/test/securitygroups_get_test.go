package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

func TestContract_GetSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}

		response := securitygroups.SecurityGroup{
			ID:          "sg-123",
			Name:        "web-sg",
			Description: "Web server security group",
			ProjectID:   "proj-1",
			UserID:      "user-1",
			Namespace:   "default",
			Rules: []securitygroups.SecurityGroupRule{
				{
					ID:         "rule-1",
					Direction:  securitygroups.DirectionIngress,
					Protocol:   securitygroups.ProtocolTCP,
					PortMin:    80,
					PortMax:    80,
					RemoteCIDR: "0.0.0.0/0",
				},
			},
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	result, err := sgClient.Get(context.Background(), "sg-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "sg-123" {
		t.Errorf("expected ID sg-123, got %s", result.ID)
	}

	if result.Name != "web-sg" {
		t.Errorf("expected name 'web-sg', got %s", result.Name)
	}

	if len(result.SecurityGroup.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(result.SecurityGroup.Rules))
	}

	// Verify Rules() method returns a rules client
	if result.Rules() == nil {
		t.Error("expected Rules() to return a RulesClient, got nil")
	}
}

func TestContract_GetSecurityGroup_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "security group not found"})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	_, err := sgClient.Get(context.Background(), "nonexistent-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
