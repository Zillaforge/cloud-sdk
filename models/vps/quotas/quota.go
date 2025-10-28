package quotas

// Quota represents project resource quotas for VPS services.
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

// QuotaDetail contains the limit and current usage for a specific resource.
type QuotaDetail struct {
	Limit int `json:"limit"` // -1 = unlimited
	Usage int `json:"usage"`
}
