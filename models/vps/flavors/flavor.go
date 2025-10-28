package flavors

// Flavor represents a compute instance flavor/size.
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

// ListFlavorsOptions provides filtering options for listing flavors.
type ListFlavorsOptions struct {
	Name   string
	Public *bool // nil = all, true = public only, false = private only
	Tag    string
}

// FlavorListResponse represents the response from listing flavors.
type FlavorListResponse struct {
	Items []*Flavor `json:"items"`
}
