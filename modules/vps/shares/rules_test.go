package shares

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	sharesmodel "github.com/Zillaforge/cloud-sdk/models/vps/shares"
)

func TestRulesClient_Create(t *testing.T) {
	mockRule := &sharesmodel.ShareRule{
		ID:     "rule-1",
		CIDR:   "10.0.0.0/24",
		Level:  "rw",
		Status: "active",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/shares/share-1/rules" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockRule)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	rulesClient := &RulesClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		shareID:    "share-1",
	}

	req := &sharesmodel.CreateShareRuleRequest{CIDR: "10.0.0.0/24", Level: "rw"}
	rule, err := rulesClient.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.ID != mockRule.ID {
		t.Errorf("expected ID %s, got %s", mockRule.ID, rule.ID)
	}
	if rule.CIDR != mockRule.CIDR {
		t.Errorf("expected CIDR %s, got %s", mockRule.CIDR, rule.CIDR)
	}
}

func TestRulesClient_Create_MissingCIDR(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	rulesClient := &RulesClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		shareID:    "share-1",
	}

	req := &sharesmodel.CreateShareRuleRequest{CIDR: "", Level: "rw"}
	_, err := rulesClient.Create(context.Background(), req)
	if err == nil {
		t.Fatal("expected validation error for missing CIDR, got nil")
	}
}

func TestRulesClient_Create_InvalidLevel(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	rulesClient := &RulesClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		shareID:    "share-1",
	}

	req := &sharesmodel.CreateShareRuleRequest{CIDR: "10.0.0.0/24", Level: "admin"}
	_, err := rulesClient.Create(context.Background(), req)
	if err == nil {
		t.Fatal("expected validation error for invalid level, got nil")
	}
}

func TestRulesClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/shares/share-1/rules/rule-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	rulesClient := &RulesClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		shareID:    "share-1",
	}

	if err := rulesClient.Delete(context.Background(), "rule-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
