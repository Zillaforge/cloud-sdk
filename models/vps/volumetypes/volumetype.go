package volumetypes

// VolumeTypeListResponse represents the response from listing volume types.
// Matches pb.VolumeTypeListOutput from vps.yaml.
type VolumeTypeListResponse struct {
	VolumeTypes []string `json:"volume_types"`
}
