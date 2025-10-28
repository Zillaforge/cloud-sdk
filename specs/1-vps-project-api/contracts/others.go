// Package contracts defines VPS router, security group, keypair, flavor, and quota operations.
// Swagger references: /api/v1/project/{project-id}/{routers,security_groups,keypairs,flavors,quotas}*
package contracts

import "context"

// RouterOperations defines all router-related API operations.
type RouterOperations interface {
	List(ctx context.Context, opts *ListRoutersOptions) (*RouterListResponse, error)
	Create(ctx context.Context, req *RouterCreateRequest) (*Router, error)
	Get(ctx context.Context, routerID string) (*RouterResource, error)
	Update(ctx context.Context, routerID string, req *RouterUpdateRequest) (*Router, error)
	Delete(ctx context.Context, routerID string) error
	SetState(ctx context.Context, routerID string, req *RouterSetStateRequest) error
}

type Router struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	State       string   `json:"state"` // enabled, disabled
	ProjectID   string   `json:"project_id"`
	Networks    []string `json:"networks,omitempty"`
}

// RouterResource wraps a Router with sub-resource operations.
type RouterResource struct {
	*Router
	networkOps RouterNetworkOperations
}

// Networks returns the network operations for this router.
func (r *RouterResource) Networks() RouterNetworkOperations {
	return r.networkOps
}

// RouterNetworkOperations defines operations on router networks (sub-resource).
type RouterNetworkOperations interface {
	// List lists all networks associated with the router.
	// GET /api/v1/project/{project-id}/routers/{router-id}/networks
	// Returns: []*RouterNetwork
	// Errors: 401, 403, 404, 500
	List(ctx context.Context) ([]*RouterNetwork, error)

	// Associate associates a network with the router.
	// POST /api/v1/project/{project-id}/routers/{router-id}/networks/{network-id}
	// Returns: 204
	// Errors: 400, 401, 403, 404, 409, 500
	Associate(ctx context.Context, networkID string) error

	// Disassociate disassociates a network from the router.
	// DELETE /api/v1/project/{project-id}/routers/{router-id}/networks/{network-id}
	// Returns: 204
	// Errors: 401, 403, 404, 500
	Disassociate(ctx context.Context, networkID string) error
}

type RouterCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type RouterUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type RouterSetStateRequest struct {
	State string `json:"state"` // enabled or disabled
}

type RouterNetwork struct {
	NetworkID string `json:"network_id"`
}

type ListRoutersOptions struct {
	Name string
}

type RouterListResponse struct {
	Items []*Router `json:"items"`
}

// SecurityGroupOperations defines all security group-related API operations.
type SecurityGroupOperations interface {
	List(ctx context.Context, opts *ListSecurityGroupsOptions) (*SecurityGroupListResponse, error)
	Create(ctx context.Context, req *SecurityGroupCreateRequest) (*SecurityGroup, error)
	Get(ctx context.Context, sgID string) (*SecurityGroupResource, error)
	Update(ctx context.Context, sgID string, req *SecurityGroupUpdateRequest) (*SecurityGroup, error)
	Delete(ctx context.Context, sgID string) error
}

type SecurityGroup struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	ProjectID   string              `json:"project_id"`
	Rules       []SecurityGroupRule `json:"rules,omitempty"`
}

// SecurityGroupResource wraps a SecurityGroup with sub-resource operations.
type SecurityGroupResource struct {
	*SecurityGroup
	ruleOps SecurityGroupRuleOperations
}

// Rules returns the rule operations for this security group.
func (sg *SecurityGroupResource) Rules() SecurityGroupRuleOperations {
	return sg.ruleOps
}

// SecurityGroupRuleOperations defines operations on security group rules (sub-resource).
type SecurityGroupRuleOperations interface {
	// Add adds a new rule to the security group.
	// POST /api/v1/project/{project-id}/security_groups/{sg-id}/rules
	// Body: SecurityGroupRuleCreateRequest
	// Returns: 201 + SecurityGroupRule
	// Errors: 400, 401, 403, 404, 409, 500
	Add(ctx context.Context, req *SecurityGroupRuleCreateRequest) (*SecurityGroupRule, error)

	// Delete deletes a rule from the security group.
	// DELETE /api/v1/project/{project-id}/security_groups/{sg-id}/rules/{rule-id}
	// Returns: 204
	// Errors: 401, 403, 404, 500
	Delete(ctx context.Context, ruleID string) error
}

type SecurityGroupRule struct {
	ID          string `json:"id"`
	Direction   string `json:"direction"` // ingress, egress
	Protocol    string `json:"protocol"`  // tcp, udp, icmp
	PortMin     int    `json:"port_range_min,omitempty"`
	PortMax     int    `json:"port_range_max,omitempty"`
	RemoteCIDR  string `json:"remote_ip_prefix,omitempty"`
	RemoteGroup string `json:"remote_group_id,omitempty"`
}

type SecurityGroupCreateRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Rules       []SecurityGroupRule `json:"rules,omitempty"`
}

type SecurityGroupUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type SecurityGroupRuleCreateRequest struct {
	Direction   string `json:"direction"`
	Protocol    string `json:"protocol"`
	PortMin     int    `json:"port_range_min,omitempty"`
	PortMax     int    `json:"port_range_max,omitempty"`
	RemoteCIDR  string `json:"remote_ip_prefix,omitempty"`
	RemoteGroup string `json:"remote_group_id,omitempty"`
}

type ListSecurityGroupsOptions struct {
	Name string
}

type SecurityGroupListResponse struct {
	Items []*SecurityGroup `json:"items"`
}

// KeypairOperations defines all keypair-related API operations.
type KeypairOperations interface {
	List(ctx context.Context, opts *ListKeypairsOptions) (*KeypairListResponse, error)
	Create(ctx context.Context, req *KeypairCreateRequest) (*Keypair, error)
	Get(ctx context.Context, keypairID string) (*Keypair, error)
	Update(ctx context.Context, keypairID string, req *KeypairUpdateRequest) (*Keypair, error)
	Delete(ctx context.Context, keypairID string) error
}

type Keypair struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PublicKey   string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
	UserID      string `json:"user_id"`
}

type KeypairCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PublicKey   string `json:"public_key,omitempty"` // import; omit to generate
}

type KeypairUpdateRequest struct {
	Description string `json:"description,omitempty"`
}

type ListKeypairsOptions struct {
	Name string
}

type KeypairListResponse struct {
	Items []*Keypair `json:"items"`
}

// FlavorOperations defines all flavor-related API operations.
type FlavorOperations interface {
	List(ctx context.Context, opts *ListFlavorsOptions) (*FlavorListResponse, error)
	Get(ctx context.Context, flavorID string) (*Flavor, error)
}

type Flavor struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	VCPUs       int      `json:"vcpus"`
	RAM         int      `json:"ram"`  // MiB
	Disk        int      `json:"disk"` // GiB
	Public      bool     `json:"public"`
	Tags        []string `json:"tags,omitempty"`
}

type ListFlavorsOptions struct {
	Name   string
	Public *bool // nil = all, true = public only, false = private only
	Tag    string
}

type FlavorListResponse struct {
	Items []*Flavor `json:"items"`
}

// QuotaOperations defines quota retrieval operations.
type QuotaOperations interface {
	Get(ctx context.Context) (*Quota, error)
}

type Quota struct {
	VM         QuotaDetail `json:"vm"`
	VCPU       QuotaDetail `json:"vcpu"`
	RAM        QuotaDetail `json:"ram"`
	GPU        QuotaDetail `json:"gpu"`
	BlockSize  QuotaDetail `json:"block_size"`
	Network    QuotaDetail `json:"network"`
	Router     QuotaDetail `json:"router"`
	FloatingIP QuotaDetail `json:"floating_ip"`
	Share      QuotaDetail `json:"share,omitempty"`
	ShareSize  QuotaDetail `json:"share_size,omitempty"`
}

type QuotaDetail struct {
	Limit int `json:"limit"` // -1 = unlimited
	Usage int `json:"usage"`
}
