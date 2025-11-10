package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

// TestSecurityGroupsIntegration_FullLifecycle tests complete security group lifecycle
func TestSecurityGroupsIntegration_FullLifecycle(t *testing.T) {
	var createdSGID string
	var sgRules []securitygroups.SecurityGroupRule

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// LIST
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/security_groups") {
			response := securitygroups.SecurityGroupListResponse{
				SecurityGroups: []securitygroups.SecurityGroup{},
			}
			if createdSGID != "" {
				response.SecurityGroups = append(response.SecurityGroups, securitygroups.SecurityGroup{
					ID:          createdSGID,
					Name:        "integration-test-sg",
					Description: "Integration test security group",
					ProjectID:   "proj-1",
					UserID:      "user-1",
					Namespace:   "default",
					Rules:       sgRules,
					CreatedAt:   "2024-01-01T00:00:00Z",
					UpdatedAt:   "2024-01-01T00:00:00Z",
				})
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		// CREATE
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/security_groups") {
			var req securitygroups.SecurityGroupCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode create request: %v", err)
			}

			createdSGID = "sg-integration-test"
			sgRules = []securitygroups.SecurityGroupRule{}
			for i, rule := range req.Rules {
				ruleResp := securitygroups.SecurityGroupRule{
					ID:         "rule-" + string(rune('1'+i)),
					Direction:  rule.Direction,
					Protocol:   rule.Protocol,
					RemoteCIDR: rule.RemoteCIDR,
				}
				if rule.PortMin != nil {
					ruleResp.PortMin = *rule.PortMin
				}
				if rule.PortMax != nil {
					ruleResp.PortMax = *rule.PortMax
				}
				sgRules = append(sgRules, ruleResp)
			}

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(&securitygroups.SecurityGroup{
				ID:          createdSGID,
				Name:        req.Name,
				Description: req.Description,
				ProjectID:   "proj-1",
				UserID:      "user-1",
				Namespace:   "default",
				Rules:       sgRules,
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-01T00:00:00Z",
			})
			return
		}

		// GET
		if r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/security_groups/sg-integration-test") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(&securitygroups.SecurityGroup{
				ID:          "sg-integration-test",
				Name:        "integration-test-sg",
				Description: "Integration test security group",
				ProjectID:   "proj-1",
				UserID:      "user-1",
				Namespace:   "default",
				Rules:       sgRules,
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-01T00:00:00Z",
			})
			return
		}

		// UPDATE
		if r.Method == http.MethodPut && strings.Contains(r.URL.Path, "/security_groups/sg-integration-test") {
			var req securitygroups.SecurityGroupUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode update request: %v", err)
			}

			name := "integration-test-sg"
			if req.Name != nil {
				name = *req.Name
			}

			desc := "Integration test security group"
			if req.Description != nil {
				desc = *req.Description
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(&securitygroups.SecurityGroup{
				ID:          "sg-integration-test",
				Name:        name,
				Description: desc,
				ProjectID:   "proj-1",
				UserID:      "user-1",
				Namespace:   "default",
				Rules:       sgRules,
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-02T00:00:00Z",
			})
			return
		}

		// DELETE
		if r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/security_groups/sg-integration-test") {
			createdSGID = ""
			sgRules = nil
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Unexpected request
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	sgClient := vpsClient.SecurityGroups()
	ctx := context.Background()

	// Step 1: List empty
	t.Run("Step1_ListEmpty", func(t *testing.T) {
		list, err := sgClient.List(ctx, nil)
		if err != nil {
			t.Fatalf("list failed: %v", err)
		}
		if len(list) != 0 {
			t.Errorf("expected 0 security groups, got %d", len(list))
		}
	})

	// Step 2: Create security group with rules
	t.Run("Step2_Create", func(t *testing.T) {
		createReq := securitygroups.SecurityGroupCreateRequest{
			Name:        "integration-test-sg",
			Description: "Integration test security group",
			Rules: []securitygroups.SecurityGroupRuleCreateRequest{
				{
					Direction:  securitygroups.DirectionEgress,
					Protocol:   securitygroups.ProtocolAny,
					RemoteCIDR: "0.0.0.0/0",
				},
			},
		}
		sg, err := sgClient.Create(ctx, createReq)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}
		if sg.ID != "sg-integration-test" {
			t.Errorf("expected ID sg-integration-test, got %s", sg.ID)
		}
		if len(sg.SecurityGroup.Rules) != 1 {
			t.Errorf("expected 1 rule, got %d", len(sg.SecurityGroup.Rules))
		}
	})

	// Step 3: Get security group
	t.Run("Step3_Get", func(t *testing.T) {
		sg, err := sgClient.Get(ctx, "sg-integration-test")
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if sg.Name != "integration-test-sg" {
			t.Errorf("expected name 'integration-test-sg', got %s", sg.Name)
		}
	})

	// Step 4: List with created security group
	t.Run("Step4_ListWithSG", func(t *testing.T) {
		list, err := sgClient.List(ctx, nil)
		if err != nil {
			t.Fatalf("list failed: %v", err)
		}
		if len(list) != 1 {
			t.Errorf("expected 1 security group, got %d", len(list))
		}
	})

	// Step 5: Update security group
	t.Run("Step5_Update", func(t *testing.T) {
		newDesc := "Updated description"
		updateReq := securitygroups.SecurityGroupUpdateRequest{
			Description: &newDesc,
		}
		sg, err := sgClient.Update(ctx, "sg-integration-test", updateReq)
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}
		if sg.SecurityGroup.Description != "Updated description" {
			t.Errorf("expected description 'Updated description', got %s", sg.SecurityGroup.Description)
		}
	})

	// Step 6: Delete security group
	t.Run("Step6_Delete", func(t *testing.T) {
		err := sgClient.Delete(ctx, "sg-integration-test")
		if err != nil {
			t.Fatalf("delete failed: %v", err)
		}
	})

	// Step 7: Verify deletion
	t.Run("Step7_VerifyDeletion", func(t *testing.T) {
		list, err := sgClient.List(ctx, nil)
		if err != nil {
			t.Fatalf("list failed: %v", err)
		}
		if len(list) != 0 {
			t.Errorf("expected 0 security groups after deletion, got %d", len(list))
		}
	})
}
