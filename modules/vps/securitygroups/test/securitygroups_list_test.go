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

func TestContract_ListSecurityGroups_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := securitygroups.SecurityGroupListResponse{
			SecurityGroups: []securitygroups.SecurityGroup{
				{
					ID:          "sg-123",
					Name:        "default",
					Description: "Default security group",
					ProjectID:   "proj-1",
					UserID:      "user-1",
					Namespace:   "default",
					Rules: []securitygroups.SecurityGroupRule{
						{
							ID:         "rule-1",
							Direction:  securitygroups.DirectionIngress,
							Protocol:   securitygroups.ProtocolTCP,
							PortMin:    22,
							PortMax:    22,
							RemoteCIDR: "0.0.0.0/0",
						},
					},
					CreatedAt: "2024-01-01T00:00:00Z",
					UpdatedAt: "2024-01-01T00:00:00Z",
				},
				{
					ID:          "sg-456",
					Name:        "web",
					Description: "Web server security group",
					ProjectID:   "proj-1",
					UserID:      "user-1",
					Namespace:   "default",
					Rules:       []securitygroups.SecurityGroupRule{},
					CreatedAt:   "2024-01-02T00:00:00Z",
					UpdatedAt:   "2024-01-02T00:00:00Z",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	result, err := sgClient.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 security groups, got %d", len(result))
	}

	if result[0].SecurityGroup.ID != "sg-123" {
		t.Errorf("expected ID sg-123, got %s", result[0].SecurityGroup.ID)
	}

	if result[0].SecurityGroup.Name != "default" {
		t.Errorf("expected name 'default', got %s", result[0].SecurityGroup.Name)
	}

	if len(result[0].SecurityGroup.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(result[0].SecurityGroup.Rules))
	}
}

func TestContract_ListSecurityGroups_WithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		if r.URL.Query().Get("name") != "web" {
			t.Errorf("expected name query param 'web', got %s", r.URL.Query().Get("name"))
		}
		if r.URL.Query().Get("detail") != "true" {
			t.Errorf("expected detail query param 'true', got %s", r.URL.Query().Get("detail"))
		}

		response := securitygroups.SecurityGroupListResponse{
			SecurityGroups: []securitygroups.SecurityGroup{
				{
					ID:          "sg-456",
					Name:        "web",
					Description: "Web server security group",
					ProjectID:   "proj-1",
					UserID:      "user-1",
					Namespace:   "default",
					Rules: []securitygroups.SecurityGroupRule{
						{
							ID:         "rule-2",
							Direction:  securitygroups.DirectionIngress,
							Protocol:   securitygroups.ProtocolTCP,
							PortMin:    80,
							PortMax:    80,
							RemoteCIDR: "0.0.0.0/0",
						},
					},
					CreatedAt: "2024-01-02T00:00:00Z",
					UpdatedAt: "2024-01-02T00:00:00Z",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	opts := &securitygroups.ListSecurityGroupsOptions{
		Name:   "web",
		Detail: true,
	}

	result, err := sgClient.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 security group, got %d", len(result))
	}

	if result[0].SecurityGroup.Name != "web" {
		t.Errorf("expected name 'web', got %s", result[0].SecurityGroup.Name)
	}
}

func TestContract_ListSecurityGroups_EmptyList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := securitygroups.SecurityGroupListResponse{
			SecurityGroups: []securitygroups.SecurityGroup{},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	result, err := sgClient.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 security groups, got %d", len(result))
	}
}

func TestContract_ListSecurityGroups_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	_, err := sgClient.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
