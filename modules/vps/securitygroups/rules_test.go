package securitygroups

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

// TestRulesClient_Create_TCP tests creating a TCP rule with ports
func TestRulesClient_Create_TCP(t *testing.T) {
	portMin := 80
	portMax := 80
	mockRule := &securitygroups.SecurityGroupRule{
		ID:         "rule-123",
		Direction:  securitygroups.DirectionIngress,
		Protocol:   securitygroups.ProtocolTCP,
		PortMin:    portMin,
		PortMax:    portMax,
		RemoteCIDR: "0.0.0.0/0",
	}

	expectedPath := "/api/v1/project/proj-123/security_groups/sg-456/rules"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockRule)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	rulesClient := NewRulesClient(baseClient, "proj-123", "sg-456")

	ctx := context.Background()
	req := securitygroups.SecurityGroupRuleCreateRequest{
		Direction:  securitygroups.DirectionIngress,
		Protocol:   securitygroups.ProtocolTCP,
		PortMin:    &portMin,
		PortMax:    &portMax,
		RemoteCIDR: "0.0.0.0/0",
	}
	rule, err := rulesClient.Create(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule == nil {
		t.Fatal("expected rule, got nil")
	}
	if rule.ID != "rule-123" {
		t.Errorf("expected ID rule-123, got %s", rule.ID)
	}
	if rule.Protocol != securitygroups.ProtocolTCP {
		t.Errorf("expected protocol TCP, got %s", rule.Protocol)
	}
	if rule.PortMin != portMin {
		t.Errorf("expected port min %d, got %d", portMin, rule.PortMin)
	}
}

// TestRulesClient_Create_ICMP tests creating an ICMP rule without ports
func TestRulesClient_Create_ICMP(t *testing.T) {
	mockRule := &securitygroups.SecurityGroupRule{
		ID:         "rule-789",
		Direction:  securitygroups.DirectionIngress,
		Protocol:   securitygroups.ProtocolICMP,
		RemoteCIDR: "0.0.0.0/0",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Verify request body has no port fields for ICMP
		var reqBody securitygroups.SecurityGroupRuleCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if reqBody.PortMin != nil {
			t.Error("expected PortMin to be nil for ICMP protocol")
		}
		if reqBody.PortMax != nil {
			t.Error("expected PortMax to be nil for ICMP protocol")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockRule)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	rulesClient := NewRulesClient(baseClient, "proj-123", "sg-456")

	ctx := context.Background()
	req := securitygroups.SecurityGroupRuleCreateRequest{
		Direction:  securitygroups.DirectionIngress,
		Protocol:   securitygroups.ProtocolICMP,
		RemoteCIDR: "0.0.0.0/0",
		// No PortMin/PortMax for ICMP
	}
	rule, err := rulesClient.Create(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Protocol != securitygroups.ProtocolICMP {
		t.Errorf("expected protocol ICMP, got %s", rule.Protocol)
	}
}

// TestRulesClient_Create_PortRange tests creating a rule with port range
func TestRulesClient_Create_PortRange(t *testing.T) {
	portMin := 8000
	portMax := 9000
	mockRule := &securitygroups.SecurityGroupRule{
		ID:         "rule-range",
		Direction:  securitygroups.DirectionIngress,
		Protocol:   securitygroups.ProtocolTCP,
		PortMin:    portMin,
		PortMax:    portMax,
		RemoteCIDR: "10.0.0.0/8",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody securitygroups.SecurityGroupRuleCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify port range
		if reqBody.PortMin == nil || *reqBody.PortMin != portMin {
			t.Errorf("expected PortMin %d, got %v", portMin, reqBody.PortMin)
		}
		if reqBody.PortMax == nil || *reqBody.PortMax != portMax {
			t.Errorf("expected PortMax %d, got %v", portMax, reqBody.PortMax)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockRule)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	rulesClient := NewRulesClient(baseClient, "proj-123", "sg-456")

	ctx := context.Background()
	req := securitygroups.SecurityGroupRuleCreateRequest{
		Direction:  securitygroups.DirectionIngress,
		Protocol:   securitygroups.ProtocolTCP,
		PortMin:    &portMin,
		PortMax:    &portMax,
		RemoteCIDR: "10.0.0.0/8",
	}
	rule, err := rulesClient.Create(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.PortMin != portMin {
		t.Errorf("expected port min %d, got %d", portMin, rule.PortMin)
	}
	if rule.PortMax != portMax {
		t.Errorf("expected port max %d, got %d", portMax, rule.PortMax)
	}
}

// TestRulesClient_Delete tests successful rule deletion
func TestRulesClient_Delete(t *testing.T) {
	expectedPath := "/api/v1/project/proj-123/security_groups/sg-456/rules/rule-789"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	rulesClient := NewRulesClient(baseClient, "proj-123", "sg-456")

	ctx := context.Background()
	err := rulesClient.Delete(ctx, "rule-789")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestRulesClient_Delete_Error tests error handling in Delete
func TestRulesClient_Delete_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "rule not found"})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	rulesClient := NewRulesClient(baseClient, "proj-123", "sg-456")

	ctx := context.Background()
	err := rulesClient.Delete(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
