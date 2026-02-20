package loadbalancers

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	lbmodels "github.com/Zillaforge/cloud-sdk/models/vps/loadbalancers"
)

// ListenersClient handles listener operations for a specific load balancer.
type ListenersClient struct {
	baseClient *internalhttp.Client
	projectID  string
	lbID       string
}

// List retrieves all listeners for the load balancer.
// GET /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners
func (c *ListenersClient) List(ctx context.Context) ([]*lbmodels.Listener, error) {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/listeners", c.projectID, c.lbID)

	req := &internalhttp.Request{Method: "GET", Path: path}

	var response lbmodels.ListenerListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list listeners for load balancer %s: %w", c.lbID, err)
	}

	return response.Listeners, nil
}

// Get retrieves a specific listener by ID.
// GET /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners/{listener-id}
func (c *ListenersClient) Get(ctx context.Context, listenerID string) (*lbmodels.Listener, error) {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/listeners/%s", c.projectID, c.lbID, listenerID)

	req := &internalhttp.Request{Method: "GET", Path: path}

	var listener lbmodels.Listener
	if err := c.baseClient.Do(ctx, req, &listener); err != nil {
		return nil, fmt.Errorf("failed to get listener %s for load balancer %s: %w", listenerID, c.lbID, err)
	}

	return &listener, nil
}

// Create creates a new listener on the load balancer.
// POST /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners
func (c *ListenersClient) Create(ctx context.Context, request *lbmodels.ListenerCreateRequest) (*lbmodels.Listener, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/listeners", c.projectID, c.lbID)

	req := &internalhttp.Request{Method: "POST", Path: path, Body: request}

	var listener lbmodels.Listener
	if err := c.baseClient.Do(ctx, req, &listener); err != nil {
		return nil, fmt.Errorf("failed to create listener for load balancer %s: %w", c.lbID, err)
	}

	return &listener, nil
}

// Update updates an existing listener.
// PUT /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners/{listener-id}
func (c *ListenersClient) Update(ctx context.Context, listenerID string, request *lbmodels.ListenerUpdateRequest) (*lbmodels.Listener, error) {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/listeners/%s", c.projectID, c.lbID, listenerID)

	req := &internalhttp.Request{Method: "PUT", Path: path, Body: request}

	var listener lbmodels.Listener
	if err := c.baseClient.Do(ctx, req, &listener); err != nil {
		return nil, fmt.Errorf("failed to update listener %s for load balancer %s: %w", listenerID, c.lbID, err)
	}

	return &listener, nil
}

// Delete deletes a listener from the load balancer.
// DELETE /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners/{listener-id}
func (c *ListenersClient) Delete(ctx context.Context, listenerID string) error {
	path := fmt.Sprintf("/api/v1/project/%s/loadbalancers/%s/listeners/%s", c.projectID, c.lbID, listenerID)

	req := &internalhttp.Request{Method: "DELETE", Path: path}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete listener %s from load balancer %s: %w", listenerID, c.lbID, err)
	}

	return nil
}
