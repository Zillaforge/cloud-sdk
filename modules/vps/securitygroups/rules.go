package securitygroups

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

// RulesClient provides operations for managing security group rules.
type RulesClient struct {
	baseClient      *internalhttp.Client
	projectID       string
	securityGroupID string
	basePath        string
}

// NewRulesClient creates a new rules client for a specific security group.
func NewRulesClient(baseClient *internalhttp.Client, projectID, securityGroupID string) *RulesClient {
	basePath := fmt.Sprintf("/api/v1/project/%s/security_groups/%s", projectID, securityGroupID)
	return &RulesClient{
		baseClient:      baseClient,
		projectID:       projectID,
		securityGroupID: securityGroupID,
		basePath:        basePath,
	}
}

// Create creates a new rule in the security group.
// POST /api/v1/project/{project-id}/security_groups/{sg-id}/rules
func (rc *RulesClient) Create(ctx context.Context, req securitygroups.SecurityGroupRuleCreateRequest) (*securitygroups.SecurityGroupRule, error) {
	path := fmt.Sprintf("%s/rules", rc.basePath)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var rule securitygroups.SecurityGroupRule
	if err := rc.baseClient.Do(ctx, httpReq, &rule); err != nil {
		return nil, fmt.Errorf("failed to create rule in security group %s: %w", rc.securityGroupID, err)
	}

	return &rule, nil
}

// Delete removes a rule from the security group.
// DELETE /api/v1/project/{project-id}/security_groups/{sg-id}/rules/{rule-id}
func (rc *RulesClient) Delete(ctx context.Context, ruleID string) error {
	path := fmt.Sprintf("%s/rules/%s", rc.basePath, ruleID)

	req := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := rc.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete rule %s from security group %s: %w", ruleID, rc.securityGroupID, err)
	}

	return nil
}
