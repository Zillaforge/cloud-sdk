package loadbalancers

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	lbmodels "github.com/Zillaforge/cloud-sdk/models/vps/loadbalancers"
)

// MembersClient handles member operations for a specific pool.
type MembersClient struct {
	baseClient *internalhttp.Client
	projectID  string
	lbID       string
	poolID     string
}

// Create adds a new member to the pool.
// POST /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}/members
func (c *MembersClient) Create(ctx context.Context, request *lbmodels.MemberCreateRequest) (*lbmodels.Pool, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/pools/%s/members",
		c.projectID, c.lbID, c.poolID)

	req := &internalhttp.Request{Method: "POST", Path: path, Body: request}

	var pool lbmodels.Pool
	if err := c.baseClient.Do(ctx, req, &pool); err != nil {
		return nil, fmt.Errorf("failed to create member in pool %s: %w", c.poolID, err)
	}

	return &pool, nil
}

// BatchUpdate replaces the entire member list for the pool.
// PUT /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}/members
func (c *MembersClient) BatchUpdate(ctx context.Context, request *lbmodels.MembersBatchUpdateRequest) error {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/pools/%s/members",
		c.projectID, c.lbID, c.poolID)

	req := &internalhttp.Request{Method: "PUT", Path: path, Body: request}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to batch update members in pool %s: %w", c.poolID, err)
	}

	return nil
}

// Delete removes a member from the pool.
// DELETE /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}/members/{member-id}
func (c *MembersClient) Delete(ctx context.Context, memberID string) error {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/pools/%s/members/%s",
		c.projectID, c.lbID, c.poolID, memberID)

	req := &internalhttp.Request{Method: "DELETE", Path: path}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete member %s from pool %s: %w", memberID, c.poolID, err)
	}

	return nil
}
