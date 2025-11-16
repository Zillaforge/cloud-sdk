package snapshots

import (
	"context"
	"fmt"
	"net/url"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	snapshotsmodel "github.com/Zillaforge/cloud-sdk/models/vps/snapshots"
)

// Client provides operations for managing snapshots.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new snapshots client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	return &Client{baseClient: baseClient, projectID: projectID, basePath: "/api/v1/project/" + projectID}
}

// Create creates a new snapshot.
func (c *Client) Create(ctx context.Context, req *snapshotsmodel.CreateSnapshotRequest) (*snapshotsmodel.Snapshot, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := c.basePath + "/snapshots"

	r := &internalhttp.Request{Method: "POST", Path: path, Body: req}

	var resp snapshotsmodel.Snapshot
	if err := c.baseClient.Do(ctx, r, &resp); err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	return &resp, nil
}

// List retrieves a list of snapshots with optional filters.
func (c *Client) List(ctx context.Context, opts *snapshotsmodel.ListSnapshotsOptions) ([]*snapshotsmodel.Snapshot, error) {
	path := c.basePath + "/snapshots"

	query := url.Values{}
	if opts != nil {
		if opts.Name != "" {
			query.Set("name", opts.Name)
		}
		if opts.VolumeID != "" {
			query.Set("volume_id", opts.VolumeID)
		}
		if opts.UserID != "" {
			query.Set("user_id", opts.UserID)
		}
		if opts.Status != "" {
			query.Set("status", opts.Status)
		}
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req := &internalhttp.Request{Method: "GET", Path: path}

	var response snapshotsmodel.SnapshotListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}

	return response.Snapshots, nil
}

// Get retrieves a specific snapshot by id.
func (c *Client) Get(ctx context.Context, id string) (*snapshotsmodel.Snapshot, error) {
	path := fmt.Sprintf("%s/snapshots/%s", c.basePath, id)

	req := &internalhttp.Request{Method: "GET", Path: path}

	var response snapshotsmodel.Snapshot
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to get snapshot %s: %w", id, err)
	}

	return &response, nil
}

// Update updates snapshot metadata (e.g., rename).
func (c *Client) Update(ctx context.Context, id string, reqBody *snapshotsmodel.UpdateSnapshotRequest) (*snapshotsmodel.Snapshot, error) {
	if err := reqBody.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	path := fmt.Sprintf("%s/snapshots/%s", c.basePath, id)

	req := &internalhttp.Request{Method: "PUT", Path: path, Body: reqBody}

	var response snapshotsmodel.Snapshot
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to update snapshot %s: %w", id, err)
	}

	return &response, nil
}

// Delete deletes a snapshot by id.
func (c *Client) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("%s/snapshots/%s", c.basePath, id)
	req := &internalhttp.Request{Method: "DELETE", Path: path}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete snapshot %s: %w", id, err)
	}
	return nil
}

// Implement other methods later (List, Get, Update, Delete)
