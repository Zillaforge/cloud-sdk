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

func TestContract_CreateRule_Success(t *testing.T) {
	// First, mock getting the security group to get the resource wrapper
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-1/security_groups/sg-123" {
			// Get security group response
			response := securitygroups.SecurityGroup{
				ID:          "sg-123",
				Name:        "test-sg",
				Description: "Test security group",
				ProjectID:   "proj-1",
				UserID:      "user-1",
				Namespace:   "default",
				Rules:       []securitygroups.SecurityGroupRule{},
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-01T00:00:00Z",
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		if r.Method == http.MethodPost && r.URL.Path == "/vps/api/v1/project/proj-1/security_groups/sg-123/rules" {
			// Create rule response
			body, _ := io.ReadAll(r.Body)
			var req securitygroups.SecurityGroupRuleCreateRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Fatalf("failed to unmarshal request body: %v", err)
			}

			// Verify request fields
			if req.Direction != securitygroups.DirectionIngress {
				t.Errorf("expected direction ingress, got %s", req.Direction)
			}
			if req.Protocol != securitygroups.ProtocolTCP {
				t.Errorf("expected protocol tcp, got %s", req.Protocol)
			}

			response := securitygroups.SecurityGroupRule{
				ID:         "rule-new-123",
				Direction:  req.Direction,
				Protocol:   req.Protocol,
				PortMin:    *req.PortMin,
				PortMax:    *req.PortMax,
				RemoteCIDR: req.RemoteCIDR,
			}

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	// Get security group resource
	sgResource, err := sgClient.Get(context.Background(), "sg-123")
	if err != nil {
		t.Fatalf("failed to get security group: %v", err)
	}

	// Create rule via Rules() sub-resource
	portMin := 80
	portMax := 80
	req := securitygroups.SecurityGroupRuleCreateRequest{
		Direction:  securitygroups.DirectionIngress,
		Protocol:   securitygroups.ProtocolTCP,
		PortMin:    &portMin,
		PortMax:    &portMax,
		RemoteCIDR: "0.0.0.0/0",
	}

	rule, err := sgResource.Rules().Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rule.ID != "rule-new-123" {
		t.Errorf("expected ID rule-new-123, got %s", rule.ID)
	}
	if rule.Protocol != securitygroups.ProtocolTCP {
		t.Errorf("expected protocol TCP, got %s", rule.Protocol)
	}
}

func TestContract_CreateRule_ICMP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-1/security_groups/sg-123" {
			response := securitygroups.SecurityGroup{
				ID:          "sg-123",
				Name:        "test-sg",
				Description: "Test security group",
				ProjectID:   "proj-1",
				UserID:      "user-1",
				Namespace:   "default",
				Rules:       []securitygroups.SecurityGroupRule{},
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-01T00:00:00Z",
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		if r.Method == http.MethodPost && r.URL.Path == "/vps/api/v1/project/proj-1/security_groups/sg-123/rules" {
			body, _ := io.ReadAll(r.Body)
			var req securitygroups.SecurityGroupRuleCreateRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Fatalf("failed to unmarshal request body: %v", err)
			}

			// Verify no ports for ICMP
			if req.PortMin != nil {
				t.Error("expected PortMin to be nil for ICMP protocol")
			}
			if req.PortMax != nil {
				t.Error("expected PortMax to be nil for ICMP protocol")
			}

			response := securitygroups.SecurityGroupRule{
				ID:         "rule-icmp-456",
				Direction:  req.Direction,
				Protocol:   req.Protocol,
				RemoteCIDR: req.RemoteCIDR,
			}

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	sgResource, err := sgClient.Get(context.Background(), "sg-123")
	if err != nil {
		t.Fatalf("failed to get security group: %v", err)
	}

	req := securitygroups.SecurityGroupRuleCreateRequest{
		Direction:  securitygroups.DirectionIngress,
		Protocol:   securitygroups.ProtocolICMP,
		RemoteCIDR: "0.0.0.0/0",
	}

	rule, err := sgResource.Rules().Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rule.Protocol != securitygroups.ProtocolICMP {
		t.Errorf("expected protocol ICMP, got %s", rule.Protocol)
	}
}

func TestContract_DeleteRule_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-1/security_groups/sg-123" {
			response := securitygroups.SecurityGroup{
				ID:          "sg-123",
				Name:        "test-sg",
				Description: "Test security group",
				ProjectID:   "proj-1",
				UserID:      "user-1",
				Namespace:   "default",
				Rules: []securitygroups.SecurityGroupRule{
					{
						ID:         "rule-to-delete",
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
			return
		}

		if r.Method == http.MethodDelete && r.URL.Path == "/vps/api/v1/project/proj-1/security_groups/sg-123/rules/rule-to-delete" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	sgResource, err := sgClient.Get(context.Background(), "sg-123")
	if err != nil {
		t.Fatalf("failed to get security group: %v", err)
	}

	err = sgResource.Rules().Delete(context.Background(), "rule-to-delete")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestContract_DeleteRule_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-1/security_groups/sg-123" {
			response := securitygroups.SecurityGroup{
				ID:          "sg-123",
				Name:        "test-sg",
				Description: "Test security group",
				ProjectID:   "proj-1",
				UserID:      "user-1",
				Namespace:   "default",
				Rules:       []securitygroups.SecurityGroupRule{},
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-01T00:00:00Z",
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "rule not found"})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()

	sgResource, err := sgClient.Get(context.Background(), "sg-123")
	if err != nil {
		t.Fatalf("failed to get security group: %v", err)
	}

	err = sgResource.Rules().Delete(context.Background(), "nonexistent-rule")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
