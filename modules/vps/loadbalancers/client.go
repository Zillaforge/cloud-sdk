package loadbalancers

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
	lbmodels "github.com/Zillaforge/cloud-sdk/models/vps/loadbalancers"
)

// Client provides operations for managing load balancers.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new load balancers client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   "/api/v1/project/" + projectID,
	}
}

// List retrieves a list of load balancers with optional filtering.
// GET /api/v1/project/{project-id}/loadbalancers
func (c *Client) List(ctx context.Context, opts *lbmodels.ListLoadBalancersOptions) ([]*LoadBalancerResource, error) {
	path := c.basePath + "/loadbalancers"

	query := url.Values{}
	if opts != nil {
		if opts.Name != "" {
			query.Set("name", opts.Name)
		}
		if opts.UserID != "" {
			query.Set("user_id", opts.UserID)
		}
		if opts.Detail {
			query.Set("detail", strconv.FormatBool(opts.Detail))
		}
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req := &internalhttp.Request{Method: "GET", Path: path}

	var response lbmodels.LoadBalancerListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list load balancers: %w", err)
	}

	resources := make([]*LoadBalancerResource, len(response.LoadBalancers))
	for i, lb := range response.LoadBalancers {
		resources[i] = newLoadBalancerResource(c.baseClient, c.projectID, lb)
	}
	return resources, nil
}

// Get retrieves details of a specific load balancer by ID.
// GET /api/v1/project/{project-id}/loadbalancers/{lb-id}
func (c *Client) Get(ctx context.Context, lbID string) (*LoadBalancerResource, error) {
	path := fmt.Sprintf("%s/loadbalancers/%s", c.basePath, lbID)

	req := &internalhttp.Request{Method: "GET", Path: path}

	var lb lbmodels.LoadBalancer
	if err := c.baseClient.Do(ctx, req, &lb); err != nil {
		return nil, fmt.Errorf("failed to get load balancer %s: %w", lbID, err)
	}

	return newLoadBalancerResource(c.baseClient, c.projectID, &lb), nil
}

// Create creates a new load balancer.
// POST /api/v1/project/{project-id}/loadbalancers
func (c *Client) Create(ctx context.Context, request *lbmodels.LoadBalancerCreateRequest) (*LoadBalancerResource, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := c.basePath + "/loadbalancers"

	req := &internalhttp.Request{Method: "POST", Path: path, Body: request}

	var lb lbmodels.LoadBalancer
	if err := c.baseClient.Do(ctx, req, &lb); err != nil {
		return nil, fmt.Errorf("failed to create load balancer: %w", err)
	}

	return newLoadBalancerResource(c.baseClient, c.projectID, &lb), nil
}

// Update updates load balancer metadata (name, description).
// PUT /api/v1/project/{project-id}/loadbalancers/{lb-id}
func (c *Client) Update(ctx context.Context, lbID string, request *lbmodels.LoadBalancerUpdateRequest) (*LoadBalancerResource, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("%s/loadbalancers/%s", c.basePath, lbID)

	req := &internalhttp.Request{Method: "PUT", Path: path, Body: request}

	var lb lbmodels.LoadBalancer
	if err := c.baseClient.Do(ctx, req, &lb); err != nil {
		return nil, fmt.Errorf("failed to update load balancer %s: %w", lbID, err)
	}

	return newLoadBalancerResource(c.baseClient, c.projectID, &lb), nil
}

// Delete deletes a load balancer.
// DELETE /api/v1/project/{project-id}/loadbalancers/{lb-id}
func (c *Client) Delete(ctx context.Context, lbID string) error {
	path := fmt.Sprintf("%s/loadbalancers/%s", c.basePath, lbID)

	req := &internalhttp.Request{Method: "DELETE", Path: path}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete load balancer %s: %w", lbID, err)
	}

	return nil
}

// AssociateFloatingIP associates a floating IP to the load balancer's VIP.
// If FIPID is empty in the request, a new floating IP will be created.
// POST /api/v1/project/{project-id}/loadbalancers/{lb-id}/floatingip
func (c *Client) AssociateFloatingIP(ctx context.Context, lbID string, request *lbmodels.AssociateFloatingIPRequest) (*floatingips.FloatingIP, error) {
	path := fmt.Sprintf("%s/loadbalancers/%s/floatingip", c.basePath, lbID)

	req := &internalhttp.Request{Method: "POST", Path: path, Body: request}

	var fip floatingips.FloatingIP
	if err := c.baseClient.Do(ctx, req, &fip); err != nil {
		return nil, fmt.Errorf("failed to associate floating IP to load balancer %s: %w", lbID, err)
	}

	return &fip, nil
}

// ListenerOperations defines operations on load balancer listeners (sub-resource).
type ListenerOperations interface {
	List(ctx context.Context) ([]*lbmodels.Listener, error)
	Get(ctx context.Context, listenerID string) (*lbmodels.Listener, error)
	Create(ctx context.Context, req *lbmodels.ListenerCreateRequest) (*lbmodels.Listener, error)
	Update(ctx context.Context, listenerID string, req *lbmodels.ListenerUpdateRequest) (*lbmodels.Listener, error)
	Delete(ctx context.Context, listenerID string) error
}

// PoolOperations defines operations on load balancer pools (sub-resource).
type PoolOperations interface {
	List(ctx context.Context) ([]*PoolResource, error)
	Get(ctx context.Context, poolID string) (*PoolResource, error)
	Create(ctx context.Context, req *lbmodels.PoolCreateRequest) (*PoolResource, error)
	Update(ctx context.Context, poolID string, req *lbmodels.PoolUpdateRequest) (*PoolResource, error)
	Delete(ctx context.Context, poolID string) error
}

// MemberOperations defines operations on pool members (sub-sub-resource).
type MemberOperations interface {
	Create(ctx context.Context, req *lbmodels.MemberCreateRequest) (*lbmodels.Pool, error)
	BatchUpdate(ctx context.Context, req *lbmodels.MembersBatchUpdateRequest) error
	Delete(ctx context.Context, memberID string) error
}

// Compile-time assertions that clients satisfy their interfaces.
var _ ListenerOperations = (*ListenersClient)(nil)
var _ PoolOperations = (*PoolsClient)(nil)
var _ MemberOperations = (*MembersClient)(nil)

// LoadBalancerResource wraps a LoadBalancer with sub-resource operations.
type LoadBalancerResource struct {
	*lbmodels.LoadBalancer
	listenerOps ListenerOperations
	poolOps     PoolOperations
}

// Listeners returns the listener operations for this load balancer.
func (r *LoadBalancerResource) Listeners() ListenerOperations {
	return r.listenerOps
}

// Pools returns the pool operations for this load balancer.
func (r *LoadBalancerResource) Pools() PoolOperations {
	return r.poolOps
}

// PoolResource wraps a Pool with member sub-resource operations.
type PoolResource struct {
	*lbmodels.Pool
	memberOps MemberOperations
}

// Members returns the member operations for this pool.
func (r *PoolResource) Members() MemberOperations {
	return r.memberOps
}

// newLoadBalancerResource constructs a LoadBalancerResource with wired sub-clients.
func newLoadBalancerResource(baseClient *internalhttp.Client, projectID string, lb *lbmodels.LoadBalancer) *LoadBalancerResource {
	return &LoadBalancerResource{
		LoadBalancer: lb,
		listenerOps: &ListenersClient{
			baseClient: baseClient,
			projectID:  projectID,
			lbID:       lb.ID,
		},
		poolOps: &PoolsClient{
			baseClient: baseClient,
			projectID:  projectID,
			lbID:       lb.ID,
		},
	}
}

// newPoolResource constructs a PoolResource with wired members sub-client.
func newPoolResource(baseClient *internalhttp.Client, projectID, lbID string, pool *lbmodels.Pool) *PoolResource {
	return &PoolResource{
		Pool: pool,
		memberOps: &MembersClient{
			baseClient: baseClient,
			projectID:  projectID,
			lbID:       lbID,
			poolID:     pool.ID,
		},
	}
}
