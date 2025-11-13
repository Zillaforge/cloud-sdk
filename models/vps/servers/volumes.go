package servers

import (
	"github.com/Zillaforge/cloud-sdk/models/vps/volumes"
)

// ServerVolume represents a volume attachment with detailed volume information.
type ServerVolume struct {
	System   bool            `json:"system"`
	VolumeID string          `json:"volume_id"`
	Device   string          `json:"device"`
	Volume   *volumes.Volume `json:"volume"`
}

// ServerVolumesResponse represents the response for listing server volumes.
type ServerVolumesResponse struct {
	Disks []*ServerVolume `json:"disks"`
}
