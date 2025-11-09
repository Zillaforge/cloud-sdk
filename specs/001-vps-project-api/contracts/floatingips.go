// Package contracts defines the VPS floating IP operations interface.
// Swagger reference: /api/v1/project/{project-id}/floatingips*
package contracts

import "context"

// FloatingIPOperations defines all floating IP-related API operations.
type FloatingIPOperations interface {
	// List retrieves all floating IPs for the project.
	// GET /api/v1/project/{project-id}/floatingips
	// Returns: FloatingIPListResponse
	// Errors: 401, 403, 500
	List(ctx context.Context, opts *ListFloatingIPsOptions) (*FloatingIPListResponse, error)

	// Create allocates a new floating IP.
	// POST /api/v1/project/{project-id}/floatingips
	// Body: FloatingIPCreateRequest
	// Returns: 201 + FloatingIP (may be PENDING if requires approval)
	// Errors: 400, 401, 403, 409 (quota), 500
	Create(ctx context.Context, req *FloatingIPCreateRequest) (*FloatingIP, error)

	// Get retrieves a specific floating IP.
	// GET /api/v1/project/{project-id}/floatingips/{fip-id}
	// Returns: 200 + FloatingIP
	// Errors: 401, 403, 404, 500
	Get(ctx context.Context, fipID string) (*FloatingIP, error)

	// Update updates floating IP description.
	// PUT /api/v1/project/{project-id}/floatingips/{fip-id}
	// Body: FloatingIPUpdateRequest
	// Returns: 200 + FloatingIP
	// Errors: 400, 401, 403, 404, 500
	Update(ctx context.Context, fipID string, req *FloatingIPUpdateRequest) (*FloatingIP, error)

	// Delete releases a floating IP.
	// DELETE /api/v1/project/{project-id}/floatingips/{fip-id}
	// Returns: 204
	// Errors: 401, 403, 404, 409 (associated), 500
	Delete(ctx context.Context, fipID string) error

	// Approve approves a pending floating IP request (admin only).
	// POST /api/v1/project/{project-id}/floatingips/{fip-id}/approve
	// Returns: 202
	// Errors: 400, 401, 403, 404, 500
	Approve(ctx context.Context, fipID string) error

	// Reject rejects a pending floating IP request (admin only).
	// POST /api/v1/project/{project-id}/floatingips/{fip-id}/reject
	// Returns: 202
	// Errors: 400, 401, 403, 404, 500
	Reject(ctx context.Context, fipID string) error

	// Disassociate disassociates a floating IP from its port.
	// DELETE /api/v1/project/{project-id}/floatingips/{fip-id}/disassociate
	// Returns: 202
	// Errors: 400, 401, 403, 404, 500
	Disassociate(ctx context.Context, fipID string) error
}

type ListFloatingIPsOptions struct {
	Status string // filter by status
}

type FloatingIPListResponse struct {
	Items []*FloatingIP `json:"items"`
}

type FloatingIP struct {
	ID          string `json:"id"`
	Address     string `json:"address"`
	Status      string `json:"status"` // ACTIVE, PENDING, DOWN, REJECTED
	ProjectID   string `json:"project_id"`
	PortID      string `json:"port_id,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type FloatingIPCreateRequest struct {
	Description string `json:"description,omitempty"`
	// ExtNetworkID omitted if default
}

type FloatingIPUpdateRequest struct {
	Description string `json:"description,omitempty"`
}
