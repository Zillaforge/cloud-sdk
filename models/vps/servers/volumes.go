package servers

// VolumeAttachment represents a volume attached to a server.
type VolumeAttachment struct {
	VolumeID string `json:"volume_id"`
	ServerID string `json:"server_id"`
	Device   string `json:"device,omitempty"`
}
