package flavors

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/flavors"
)

// Client provides operations for managing flavors.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new flavors client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	basePath := "/api/v1/project/" + projectID
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// List retrieves a list of available flavors with optional filtering.
// GET /api/v1/project/{project-id}/flavors
func (c *Client) List(ctx context.Context, opts *flavors.ListFlavorsOptions) (*flavors.FlavorListResponse, error) {
	path := c.basePath + "/flavors"

	// Build query parameters
	query := url.Values{}
	if opts != nil {
		if opts.Name != "" {
			query.Set("name", opts.Name)
		}
		if opts.Public != nil {
			query.Set("public", strconv.FormatBool(*opts.Public))
		}
		if opts.Tag != "" {
			query.Set("tag", opts.Tag)
		}
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response flavors.FlavorListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list flavors: %w", err)
	}

	return &response, nil
}

// Get retrieves details of a specific flavor by ID.
// GET /api/v1/project/{project-id}/flavors/{flavor-id}
func (c *Client) Get(ctx context.Context, flavorID string) (*flavors.Flavor, error) {
	path := fmt.Sprintf("%s/flavors/%s", c.basePath, flavorID)

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var flavor flavors.Flavor
	if err := c.baseClient.Do(ctx, req, &flavor); err != nil {
		return nil, fmt.Errorf("failed to get flavor %s: %w", flavorID, err)
	}

	return &flavor, nil
}
