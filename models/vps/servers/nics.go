package servers

import (
	"github.com/Zillaforge/cloud-sdk/models/vps/common"
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
)

// ServerNIC represents a vNIC attached to a server (ServerNICInfo from pb).
type ServerNIC struct {
	ID             string                  `json:"id"`
	MAC            string                  `json:"mac"`
	NetworkID      string                  `json:"network_id"`
	Network        *common.IDName          `json:"network,omitempty"`
	Addresses      []string                `json:"addresses,omitempty"`
	FloatingIP     *floatingips.FloatingIP `json:"floating_ip,omitempty"`
	SecurityGroups []*common.IDName        `json:"security_groups,omitempty"`
	SGIDs          []string                `json:"sg_ids,omitempty"`
	IsProviderNet  bool                    `json:"is_provider_net,omitempty"`
}

// ServerNICCreateRequest is the body for AddServerNIC.
type ServerNICCreateRequest struct {
	NetworkID string   `json:"network_id"`
	SGIDs     []string `json:"sg_ids"`
	FixedIP   string   `json:"fixed_ip,omitempty"`
}

// ServerNICUpdateRequest is the body for UpdateServerNIC.
type ServerNICUpdateRequest struct {
	SGIDs []string `json:"sg_ids"`
}

// ServerNICsListResponse is the response from listing NICs.
type ServerNICsListResponse struct {
	NICs []*ServerNIC `json:"nics"`
}

// ServerNICAssociateFloatingIPRequest is the body for AssociateFloatingIPToNIC.
type ServerNICAssociateFloatingIPRequest struct {
	FIPID string `json:"fip_id,omitempty"`
}
