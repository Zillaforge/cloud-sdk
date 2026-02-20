package loadbalancers

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	lbmodels "github.com/Zillaforge/cloud-sdk/models/vps/loadbalancers"
)

// PoolsClient handles pool operations for a specific load balancer.
type PoolsClient struct {
	baseClient *internalhttp.Client
	projectID  string
	lbID       string
}

// List retrieves all pools for the load balancer.
// GET /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools
func (c *PoolsClient) List(ctx context.Context) ([]*PoolResource, error) {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/pools", c.projectID, c.lbID)

	req := &internalhttp.Request{Method: "GET", Path: path}

	var response lbmodels.PoolListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list pools for load balancer %s: %w", c.lbID, err)
	}

	resources := make([]*PoolResource, len(response.Pools))
	for i, pool := range response.Pools {
		resources[i] = newPoolResource(c.baseClient, c.projectID, c.lbID, pool)
	}
	return resources, nil
}

// Get retrieves a specific pool by ID.
// GET /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}
func (c *PoolsClient) Get(ctx context.Context, poolID string) (*PoolResource, error) {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/pools/%s", c.projectID, c.lbID, poolID)

	req := &internalhttp.Request{Method: "GET", Path: path}

	var pool lbmodels.Pool
	if err := c.baseClient.Do(ctx, req, &pool); err != nil {
		return nil, fmt.Errorf("failed to get pool %s for load balancer %s: %w", poolID, c.lbID, err)
	}

	return newPoolResource(c.baseClient, c.projectID, c.lbID, &pool), nil
}

// Create creates a new pool on the load balancer.
// POST /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools
func (c *PoolsClient) Create(ctx context.Context, request *lbmodels.PoolCreateRequest) (*PoolResource, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/pools", c.projectID, c.lbID)

	req := &internalhttp.Request{Method: "POST", Path: path, Body: request}

	var pool lbmodels.Pool
	if err := c.baseClient.Do(ctx, req, &pool); err != nil {
		return nil, fmt.Errorf("failed to create pool for load balancer %s: %w", c.lbID, err)
	}

	return newPoolResource(c.baseClient, c.projectID, c.lbID, &pool), nil
}

// Update updates an existing pool.
// PUT /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}
func (c *PoolsClient) Update(ctx context.Context, poolID string, request *lbmodels.PoolUpdateRequest) (*PoolResource, error) {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/pools/%s", c.projectID, c.lbID, poolID)

	req := &internalhttp.Request{Method: "PUT", Path: path, Body: request}

	var pool lbmodels.Pool
	if err := c.baseClient.Do(ctx, req, &pool); err != nil {
		return nil, fmt.Errorf("failed to update pool %s for load balancer %s: %w", poolID, c.lbID, err)
	}

	return newPoolResource(c.baseClient, c.projectID, c.lbID, &pool), nil
}

// Delete deletes a pool from the load balancer.
// DELETE /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}
func (c *PoolsClient) Delete(ctx context.Context, poolID string) error {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/pools/%s", c.projectID, c.lbID, poolID)

	req := &internalhttp.Request{Method: "DELETE", Path: path}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete pool %s from load balancer %s: %w", poolID, c.lbID, err)
	}

	return nil
}
