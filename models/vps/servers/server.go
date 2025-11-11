package servers

import (
	"github.com/Zillaforge/cloud-sdk/models/vps/common"
	"github.com/Zillaforge/cloud-sdk/models/vps/flavors"
)

// ServerStatus represents the status of a server.
type ServerStatus string

// Server status constants
const (
	ServerStatusActive    ServerStatus = "ACTIVE"
	ServerStatusBuild     ServerStatus = "BUILD"
	ServerStatusShutoff   ServerStatus = "SHUTOFF"
	ServerStatusError     ServerStatus = "ERROR"
	ServerStatusReboot    ServerStatus = "REBOOT"
	ServerStatusDeleted   ServerStatus = "DELETED"
	ServerStatusSuspended ServerStatus = "SUSPENDED"
)

// VRMImgInfo contains VRM image repository information.
type VRMImgInfo struct {
	RepositoryID   string `json:"repository_id"`
	RepositoryName string `json:"repository_name"`
	TagID          string `json:"tag_id"`
	TagName        string `json:"tag_name"`
}

// Server represents a compute instance with all pb.ServerInfo fields.
type Server struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Status       ServerStatus      `json:"status"`
	StatusReason string            `json:"status_reason,omitempty"`
	FlavorID     string            `json:"flavor_id"`
	Flavor       *common.IDName    `json:"flavor,omitempty"`
	FlavorDetail *flavors.Flavor   `json:"flavor_detail,omitempty"`
	ImageID      string            `json:"image_id"`
	Image        *VRMImgInfo       `json:"image,omitempty"`
	ProjectID    string            `json:"project_id"`
	Project      *common.IDName    `json:"project,omitempty"`
	UserID       string            `json:"user_id,omitempty"`
	User         *common.IDName    `json:"user,omitempty"`
	KeypairID    string            `json:"keypair_id,omitempty"`
	Keypair      *common.IDName    `json:"keypair,omitempty"`
	Metadatas    map[string]string `json:"metadatas,omitempty"`
	PrivateIPs   []string          `json:"private_ips,omitempty"`
	PublicIPs    []string          `json:"public_ips,omitempty"`
	AZ           string            `json:"az,omitempty"`
	Namespace    string            `json:"namespace,omitempty"`
	RootDiskID   string            `json:"root_disk_id,omitempty"`
	RootDiskSize int               `json:"root_disk_size,omitempty"`
	BootScript   string            `json:"boot_script,omitempty"`
	ApprovedAt   string            `json:"approvedAt,omitempty"`
	CreatedAt    string            `json:"createdAt"`
	UpdatedAt    string            `json:"updatedAt"`
	UUID         string            `json:"uuid,omitempty"`
}

// ServersListRequest contains filter/pagination options for List.
type ServersListRequest struct {
	Name     string `json:"name,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	Status   string `json:"status,omitempty"`
	FlavorID string `json:"flavor_id,omitempty"`
	ImageID  string `json:"image_id,omitempty"`
	Detail   bool   `json:"detail,omitempty"`
}

// ServersListResponse is the response from List.
type ServersListResponse struct {
	Servers []*Server `json:"servers"`
}

// ServerDiskRequest specifies disk/volume for server creation.
type ServerDiskRequest struct {
	Name     string `json:"name,omitempty"`
	VolumeID string `json:"volume_id,omitempty"`
	Type     string `json:"type,omitempty"`
	Size     int    `json:"size,omitempty"`
}

// ServerCreateRequest is the body for Create.
type ServerCreateRequest struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	FlavorID    string                   `json:"flavor_id"`
	ImageID     string                   `json:"image_id"`
	NICs        []ServerNICCreateRequest `json:"nics"`
	SGIDs       []string                 `json:"sg_ids,omitempty"`
	KeypairID   string                   `json:"keypair_id,omitempty"`
	Password    string                   `json:"password,omitempty"`    // Base64
	BootScript  string                   `json:"boot_script,omitempty"` // Base64
	VolumeIDs   []string                 `json:"volume_ids,omitempty"`  // deprecated
	Volumes     []ServerDiskRequest      `json:"volumes,omitempty"`
}

// ServerUpdateRequest is the body for Update.
type ServerUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
