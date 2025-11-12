// Package users provides data models for IAM user operations.
package users

// User represents an authenticated user in the IAM system.
type User struct {
	UserID      string                 `json:"userId"`
	Account     string                 `json:"account"`
	DisplayName string                 `json:"displayName"`
	Description string                 `json:"description"`
	Extra       map[string]interface{} `json:"extra"`
	Namespace   string                 `json:"namespace"`
	Email       string                 `json:"email"`
	Frozen      bool                   `json:"frozen"`
	MFA         bool                   `json:"mfa"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
	LastLoginAt string                 `json:"lastLoginAt"`
}

// GetUserResponse represents the response from the GET /user endpoint.
type GetUserResponse struct {
	*User
}
