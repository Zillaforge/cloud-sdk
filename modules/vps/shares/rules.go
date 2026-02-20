package shares

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	sharesmodel "github.com/Zillaforge/cloud-sdk/models/vps/shares"
)

// RulesClient handles access rule operations for a specific share.
type RulesClient struct {
	baseClient *internalhttp.Client
	projectID  string
	shareID    string
}

// Create adds an access rule to the share.
// POST /api/v1/project/{project-id}/shares/{share-id}/rules
func (c *RulesClient) Create(ctx context.Context, req *sharesmodel.CreateShareRuleRequest) (*sharesmodel.ShareRule, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/project/%s/shares/%s/rules", c.projectID, c.shareID)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var rule sharesmodel.ShareRule
	if err := c.baseClient.Do(ctx, httpReq, &rule); err != nil {
		return nil, fmt.Errorf("failed to create rule on share %s: %w", c.shareID, err)
	}

	return &rule, nil
}

// Delete removes an access rule from the share.
// DELETE /api/v1/project/{project-id}/shares/{share-id}/rules/{rule-id}
func (c *RulesClient) Delete(ctx context.Context, ruleID string) error {
	path := fmt.Sprintf("/api/v1/project/%s/shares/%s/rules/%s", c.projectID, c.shareID, ruleID)

	httpReq := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, httpReq, nil); err != nil {
		return fmt.Errorf("failed to delete rule %s from share %s: %w", ruleID, c.shareID, err)
	}

	return nil
}
