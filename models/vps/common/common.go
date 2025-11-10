package common

// IDName represents a reference to another resource with ID and name
type IDName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
