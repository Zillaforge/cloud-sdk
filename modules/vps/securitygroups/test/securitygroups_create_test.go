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

func TestContract_CreateSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var req securitygroups.SecurityGroupCreateRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		if req.Name != "web-sg" {
			t.Errorf("expected name 'web-sg', got %s", req.Name)
		}

		response := securitygroups.SecurityGroup{
			ID:          "sg-new-123",
			Name:        req.Name,
			Description: req.Description,
			ProjectID:   "proj-1",
			UserID:      "user-1",
			Namespace:   "default",
			Rules:       []securitygroups.SecurityGroupRule{},
			CreatedAt:   "2024-01-01T00:00:00Z",
			UpdatedAt:   "2024-01-01T00:00:00Z",
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	req := securitygroups.SecurityGroupCreateRequest{
		Name:        "web-sg",
		Description: "Web server security group",
	}

	result, err := sgClient.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "sg-new-123" {
		t.Errorf("expected ID sg-new-123, got %s", result.ID)
	}

	if result.Name != "web-sg" {
		t.Errorf("expected name 'web-sg', got %s", result.Name)
	}
}

func TestContract_CreateSecurityGroup_WithRules(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req securitygroups.SecurityGroupCreateRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		if len(req.Rules) != 2 {
			t.Errorf("expected 2 rules, got %d", len(req.Rules))
		}

		response := securitygroups.SecurityGroup{
			ID:          "sg-new-456",
			Name:        req.Name,
			Description: req.Description,
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
				{
					ID:         "rule-2",
					Direction:  securitygroups.DirectionIngress,
					Protocol:   securitygroups.ProtocolTCP,
					PortMin:    443,
					PortMax:    443,
					RemoteCIDR: "0.0.0.0/0",
				},
			},
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	req := securitygroups.SecurityGroupCreateRequest{
		Name:        "web-sg",
		Description: "Web server security group",
		Rules: []securitygroups.SecurityGroupRuleCreateRequest{
			{
				Direction:  securitygroups.DirectionIngress,
				Protocol:   securitygroups.ProtocolTCP,
				PortMin:    intPtr(80),
				PortMax:    intPtr(80),
				RemoteCIDR: "0.0.0.0/0",
			},
			{
				Direction:  securitygroups.DirectionIngress,
				Protocol:   securitygroups.ProtocolTCP,
				PortMin:    intPtr(443),
				PortMax:    intPtr(443),
				RemoteCIDR: "0.0.0.0/0",
			},
		},
	}

	result, err := sgClient.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.SecurityGroup.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(result.SecurityGroup.Rules))
	}
}

func TestContract_CreateSecurityGroup_ValidationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "name is required"})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	req := securitygroups.SecurityGroupCreateRequest{
		Name: "", // Invalid: empty name
	}

	_, err := sgClient.Create(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// intPtr returns a pointer to an int (helper for optional port fields in request models)
func intPtr(i int) *int {
	return &i
}
