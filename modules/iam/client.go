// Package iam provides the IAM (Identity and Access Management) API client.
package iam

import (
	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/modules/iam/projects"
	"github.com/Zillaforge/cloud-sdk/modules/iam/users"
)

// Client represents the IAM API client.
// IAM is non-project-scoped and provides identity and access management operations.
type Client struct {
	baseClient *internalhttp.Client
	basePath   string
}

// NewClient creates a new IAM client with the provided HTTP client.
func NewClient(baseClient *internalhttp.Client) *Client {
	return &Client{
		baseClient: baseClient,
		basePath:   "/api/v1/",
	}
}

// Users returns a client for user operations.
func (c *Client) Users() *users.Client {
	return users.NewClient(c.baseClient, c.basePath)
}

// Projects returns a client for project operations.
func (c *Client) Projects() *projects.Client {
	return projects.NewClient(c.baseClient, c.basePath)
}
