// Package projects provides the client for IAM project operations.
package projects

import (
	"context"
	"fmt"
	"strconv"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/iam/projects"
)

// Client handles project operations for the IAM API.
type Client struct {
	baseClient *internalhttp.Client
	basePath   string
}

// NewClient creates a new projects client with the provided HTTP client.
func NewClient(baseClient *internalhttp.Client, basePath string) *Client {
	return &Client{
		baseClient: baseClient,
		basePath:   basePath,
	}
}

// List retrieves all projects the user belongs to with optional pagination.
func (c *Client) List(ctx context.Context, opts *projects.ListProjectsOptions) ([]*projects.ProjectMembership, error) {
	// Build query parameters
	path := c.basePath + "projects"
	if opts != nil {
		query := ""
		if opts.Offset != nil {
			if query != "" {
				query += "&"
			}
			query += "offset=" + strconv.Itoa(*opts.Offset)
		}
		if opts.Limit != nil {
			if query != "" {
				query += "&"
			}
			query += "limit=" + strconv.Itoa(*opts.Limit)
		}
		if opts.Order != nil {
			if query != "" {
				query += "&"
			}
			query += "order=" + *opts.Order
		}
		if query != "" {
			path += "?" + query
		}
	}

	var response projects.ListProjectsResponse

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return response.Projects, nil
}

// Get retrieves specific project details by project ID.
func (c *Client) Get(ctx context.Context, projectID string) (*projects.GetProjectResponse, error) {
	var response projects.GetProjectResponse

	// Build path with project ID
	path := c.basePath + "project/" + projectID

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to get project %s: %w", projectID, err)
	}

	return &response, nil
}
