package keypairs

import "github.com/Zillaforge/cloud-sdk/models/vps/common"

// Keypair represents an SSH keypair for server access.
type Keypair struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	PublicKey   string         `json:"public_key"`
	Fingerprint string         `json:"fingerprint"`
	UserID      string         `json:"user_id"`
	PrivateKey  string         `json:"private_key,omitempty"`
	User        *common.IDName `json:"user,omitempty"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
}

// KeypairCreateRequest represents a request to create a new keypair.
type KeypairCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PublicKey   string `json:"public_key,omitempty"` // import existing key; omit to generate new key
}

// KeypairUpdateRequest represents a request to update a keypair.
type KeypairUpdateRequest struct {
	Description string `json:"description,omitempty"`
}

// ListKeypairsOptions provides filtering options for listing keypairs.
type ListKeypairsOptions struct {
	Name string
}

// KeypairListResponse represents the response from listing keypairs.
type KeypairListResponse struct {
	Keypairs []Keypair `json:"keypairs"`
}
