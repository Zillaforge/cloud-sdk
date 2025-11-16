// Package repositories provides data models for VRM repository operations.
package repositories

import (
	"fmt"
	"strings"

	"github.com/Zillaforge/cloud-sdk/models/vrm/common"
)

// Repository is an alias for common.Repository
type Repository = common.Repository

// Tag is an alias for common.Tag
type Tag = common.Tag

// CreateRepositoryRequest represents a request to create a new repository.
// It includes the required name and operating system, with an optional description.
// The operating system must be either "linux" or "windows".
type CreateRepositoryRequest struct {
	Name            string `json:"name"`
	OperatingSystem string `json:"operatingSystem"`
	Description     string `json:"description,omitempty"`
}

// Validate validates the CreateRepositoryRequest.
// It ensures that name and operating system are provided and valid,
// with operating system restricted to "linux" or "windows".
func (r *CreateRepositoryRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("createRepositoryRequest cannot be nil")
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required and must not be empty")
	}
	if strings.TrimSpace(r.OperatingSystem) == "" {
		return fmt.Errorf("operatingSystem is required and must not be empty")
	}
	if r.OperatingSystem != "linux" && r.OperatingSystem != "windows" {
		return fmt.Errorf("operatingSystem must be 'linux' or 'windows'")
	}
	return nil
}

// UpdateRepositoryRequest represents a request to update an existing repository.
// Currently only supports updating the repository description.
type UpdateRepositoryRequest struct {
	Description string `json:"description,omitempty"`
}

// Validate validates the UpdateRepositoryRequest.
// Currently performs basic nil checks but may be extended for future validation rules.
func (r *UpdateRepositoryRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("updateRepositoryRequest cannot be nil")
	}
	return nil
}

// CreateSnapshotRequest represents a request to create a snapshot from a server.
// It supports two modes: creating a snapshot into a new repository (name + operating system)
// or targeting an existing repository via repository ID. Version is always required.
type CreateSnapshotRequest struct {
	Version         string `json:"version"`
	Name            string `json:"name,omitempty"`
	OperatingSystem string `json:"operatingSystem,omitempty"`
	Description     string `json:"description,omitempty"`
	RepositoryID    string `json:"repositoryId,omitempty"`
}

// Validate ensures the CreateSnapshotRequest matches the supported modes and values.
func (r *CreateSnapshotRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("createSnapshotRequest cannot be nil")
	}

	if strings.TrimSpace(r.Version) == "" {
		return fmt.Errorf("version is required and must not be empty")
	}

	repositoryIDProvided := strings.TrimSpace(r.RepositoryID) != ""

	if repositoryIDProvided {
		return nil
	}

	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required when repositoryId is not provided")
	}
	if strings.TrimSpace(r.OperatingSystem) == "" {
		return fmt.Errorf("operatingSystem is required when repositoryId is not provided")
	}
	if r.OperatingSystem != "linux" && r.OperatingSystem != "windows" {
		return fmt.Errorf("operatingSystem must be 'linux' or 'windows'")
	}

	return nil
}

// ToCreateSnapshotRequest returns itself for compatibility.
func (r *CreateSnapshotRequest) ToCreateSnapshotRequest() CreateSnapshotRequest {
	return *r
}

// SnapshotRequester is an interface for snapshot request types.
type SnapshotRequester interface {
	ToCreateSnapshotRequest() CreateSnapshotRequest
	Validate() error
}

// CreateSnapshotFromNewRepositoryRequest represents a request to create a snapshot into a new repository.
type CreateSnapshotFromNewRepositoryRequest struct {
	Name            string `json:"name"`
	OperatingSystem string `json:"operatingSystem"`
	Version         string `json:"version"`
	Description     string `json:"description,omitempty"`
}

// Validate validates the CreateSnapshotFromNewRepositoryRequest.
func (r *CreateSnapshotFromNewRepositoryRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if r.OperatingSystem != "linux" && r.OperatingSystem != "windows" {
		return fmt.Errorf("operatingSystem must be 'linux' or 'windows'")
	}
	if strings.TrimSpace(r.Version) == "" {
		return fmt.Errorf("version is required")
	}
	return nil
}

// toCreateSnapshotRequest converts to CreateSnapshotRequest.
func (r *CreateSnapshotFromNewRepositoryRequest) ToCreateSnapshotRequest() CreateSnapshotRequest {
	return CreateSnapshotRequest{
		Name:            r.Name,
		OperatingSystem: r.OperatingSystem,
		Version:         r.Version,
		Description:     r.Description,
	}
}

// CreateSnapshotFromExistingRepositoryRequest represents a request to create a snapshot into an existing repository.
type CreateSnapshotFromExistingRepositoryRequest struct {
	RepositoryID string `json:"repositoryId"`
	Version      string `json:"version"`
}

// Validate validates the CreateSnapshotFromExistingRepositoryRequest.
func (r *CreateSnapshotFromExistingRepositoryRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if strings.TrimSpace(r.RepositoryID) == "" {
		return fmt.Errorf("repositoryId is required")
	}
	if strings.TrimSpace(r.Version) == "" {
		return fmt.Errorf("version is required")
	}
	return nil
}

// toCreateSnapshotRequest converts to CreateSnapshotRequest.
func (r *CreateSnapshotFromExistingRepositoryRequest) ToCreateSnapshotRequest() CreateSnapshotRequest {
	return CreateSnapshotRequest{
		RepositoryID: r.RepositoryID,
		Version:      r.Version,
	}
}

// CreateSnapshotResponse represents the API response after creating a snapshot.
// It returns the repository metadata (existing or newly created) and the tag
// generated for the snapshot version.
type CreateSnapshotResponse struct {
	Repository *common.Repository `json:"repository"`
	Tag        *common.Tag        `json:"tag"`
}

var (
	validDiskFormats = map[string]struct{}{
		"ami": {}, "ari": {}, "aki": {}, "vhd": {}, "vmdk": {}, "raw": {}, "qcow2": {}, "vdi": {}, "iso": {},
	}
	validContainerFormats = map[string]struct{}{
		"ami": {}, "ari": {}, "aki": {}, "bare": {}, "ovf": {},
	}
)

// UploadImageRequest represents a request to upload an image into VRM.
// The request supports three mutually exclusive modes:
//  1. Creating a new repository (name + operating system + metadata)
//  2. Targeting an existing repository via repository ID
//  3. Targeting an existing tag via tag ID
//
// The common field across all modes is the image filepath source.
type UploadImageRequest struct {
	Name            string `json:"name,omitempty"`
	OperatingSystem string `json:"operatingSystem,omitempty"`
	Version         string `json:"version,omitempty"`
	Type            string `json:"type,omitempty"`
	DiskFormat      string `json:"diskFormat,omitempty"`
	ContainerFormat string `json:"containerFormat,omitempty"`
	Filepath        string `json:"filepath"`
	RepositoryID    string `json:"repositoryId,omitempty"`
	TagID           string `json:"tagId,omitempty"`
}

// UploadToNewRepositoryRequest represents a request to upload an image and create a new repository.
type UploadToNewRepositoryRequest struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Type            string `json:"type"`
	DiskFormat      string `json:"diskFormat"`
	ContainerFormat string `json:"containerFormat"`
	OperatingSystem string `json:"operatingSystem"`
	Description     string `json:"description,omitempty"`
	Filepath        string `json:"filepath"`
}

// Validate validates the UploadToNewRepositoryRequest.
func (r *UploadToNewRepositoryRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("uploadToNewRepositoryRequest cannot be nil")
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(r.OperatingSystem) == "" {
		return fmt.Errorf("operatingSystem is required")
	}
	if r.OperatingSystem != "linux" && r.OperatingSystem != "windows" {
		return fmt.Errorf("operatingSystem must be 'linux' or 'windows'")
	}
	if strings.TrimSpace(r.Version) == "" {
		return fmt.Errorf("version is required")
	}
	if strings.TrimSpace(r.Type) == "" {
		return fmt.Errorf("type is required")
	}
	if strings.TrimSpace(r.DiskFormat) == "" {
		return fmt.Errorf("diskFormat is required")
	}
	if _, ok := validDiskFormats[r.DiskFormat]; !ok {
		return fmt.Errorf("invalid diskFormat: %s", r.DiskFormat)
	}
	if strings.TrimSpace(r.ContainerFormat) == "" {
		return fmt.Errorf("containerFormat is required")
	}
	if _, ok := validContainerFormats[r.ContainerFormat]; !ok {
		return fmt.Errorf("invalid containerFormat: %s", r.ContainerFormat)
	}
	if strings.TrimSpace(r.Filepath) == "" {
		return fmt.Errorf("filepath is required")
	}
	return nil
}

// ToUploadImageRequest converts to UploadImageRequest.
func (r *UploadToNewRepositoryRequest) ToUploadImageRequest() UploadImageRequest {
	return UploadImageRequest{
		Name:            r.Name,
		OperatingSystem: r.OperatingSystem,
		Version:         r.Version,
		Type:            r.Type,
		DiskFormat:      r.DiskFormat,
		ContainerFormat: r.ContainerFormat,
		Filepath:        r.Filepath,
	}
}

// UploadToExistingRepositoryRequest represents a request to upload an image to an existing repository.
type UploadToExistingRepositoryRequest struct {
	RepositoryID    string `json:"repositoryId"`
	Version         string `json:"version"`
	Type            string `json:"type"`
	DiskFormat      string `json:"diskFormat"`
	ContainerFormat string `json:"containerFormat"`
	Filepath        string `json:"filepath"`
}

// Validate validates the UploadToExistingRepositoryRequest.
func (r *UploadToExistingRepositoryRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("uploadToExistingRepositoryRequest cannot be nil")
	}
	if strings.TrimSpace(r.RepositoryID) == "" {
		return fmt.Errorf("repositoryId is required")
	}
	if strings.TrimSpace(r.Version) == "" {
		return fmt.Errorf("version is required")
	}
	if strings.TrimSpace(r.Type) == "" {
		return fmt.Errorf("type is required")
	}
	if strings.TrimSpace(r.DiskFormat) == "" {
		return fmt.Errorf("diskFormat is required")
	}
	if _, ok := validDiskFormats[r.DiskFormat]; !ok {
		return fmt.Errorf("invalid diskFormat: %s", r.DiskFormat)
	}
	if strings.TrimSpace(r.ContainerFormat) == "" {
		return fmt.Errorf("containerFormat is required")
	}
	if _, ok := validContainerFormats[r.ContainerFormat]; !ok {
		return fmt.Errorf("invalid containerFormat: %s", r.ContainerFormat)
	}
	if strings.TrimSpace(r.Filepath) == "" {
		return fmt.Errorf("filepath is required")
	}
	return nil
}

// ToUploadImageRequest converts to UploadImageRequest.
func (r *UploadToExistingRepositoryRequest) ToUploadImageRequest() UploadImageRequest {
	return UploadImageRequest{
		RepositoryID:    r.RepositoryID,
		Version:         r.Version,
		Type:            r.Type,
		DiskFormat:      r.DiskFormat,
		ContainerFormat: r.ContainerFormat,
		Filepath:        r.Filepath,
	}
}

// UploadToExistingTagRequest represents a request to upload an image to an existing tag.
type UploadToExistingTagRequest struct {
	TagID    string `json:"tagId"`
	Filepath string `json:"filepath"`
}

// Validate validates the UploadToExistingTagRequest.
func (r *UploadToExistingTagRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("uploadToExistingTagRequest cannot be nil")
	}
	if strings.TrimSpace(r.TagID) == "" {
		return fmt.Errorf("tagId is required")
	}
	if strings.TrimSpace(r.Filepath) == "" {
		return fmt.Errorf("filepath is required")
	}
	return nil
}

// ToUploadImageRequest converts to UploadImageRequest.
func (r *UploadToExistingTagRequest) ToUploadImageRequest() UploadImageRequest {
	return UploadImageRequest{
		TagID:    r.TagID,
		Filepath: r.Filepath,
	}
}

// UploadRequester is an interface for upload request types.
type UploadRequester interface {
	Validate() error
	ToUploadImageRequest() UploadImageRequest
}

// Validate ensures the UploadImageRequest satisfies one of the supported modes.
func (r *UploadImageRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("uploadImageRequest cannot be nil")
	}
	if strings.TrimSpace(r.Filepath) == "" {
		return fmt.Errorf("filepath is required and must not be empty")
	}

	hasRepo := strings.TrimSpace(r.RepositoryID) != ""
	hasTag := strings.TrimSpace(r.TagID) != ""

	if hasRepo && hasTag {
		return fmt.Errorf("only one of repositoryId or tagId can be provided")
	}

	if hasTag {
		return nil
	}

	if hasRepo {
		if strings.TrimSpace(r.Version) == "" {
			return fmt.Errorf("version is required when repositoryId is provided")
		}
		if strings.TrimSpace(r.Type) == "" {
			return fmt.Errorf("type is required when repositoryId is provided")
		}
		if strings.TrimSpace(r.DiskFormat) == "" {
			return fmt.Errorf("diskFormat is required when repositoryId is provided")
		}
		if _, ok := validDiskFormats[r.DiskFormat]; !ok {
			return fmt.Errorf("invalid diskFormat: %s", r.DiskFormat)
		}
		if strings.TrimSpace(r.ContainerFormat) == "" {
			return fmt.Errorf("containerFormat is required when repositoryId is provided")
		}
		if _, ok := validContainerFormats[r.ContainerFormat]; !ok {
			return fmt.Errorf("invalid containerFormat: %s", r.ContainerFormat)
		}
		return nil
	}

	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required when repositoryId and tagId are not provided")
	}
	if strings.TrimSpace(r.OperatingSystem) == "" {
		return fmt.Errorf("operatingSystem is required when repositoryId and tagId are not provided")
	}
	if r.OperatingSystem != "linux" && r.OperatingSystem != "windows" {
		return fmt.Errorf("operatingSystem must be 'linux' or 'windows'")
	}
	if strings.TrimSpace(r.Version) == "" {
		return fmt.Errorf("version is required when repositoryId and tagId are not provided")
	}
	if strings.TrimSpace(r.Type) == "" {
		return fmt.Errorf("type is required when repositoryId and tagId are not provided")
	}
	if strings.TrimSpace(r.DiskFormat) == "" {
		return fmt.Errorf("diskFormat is required when repositoryId and tagId are not provided")
	}
	if _, ok := validDiskFormats[r.DiskFormat]; !ok {
		return fmt.Errorf("invalid diskFormat: %s", r.DiskFormat)
	}
	if strings.TrimSpace(r.ContainerFormat) == "" {
		return fmt.Errorf("containerFormat is required when repositoryId and tagId are not provided")
	}
	if _, ok := validContainerFormats[r.ContainerFormat]; !ok {
		return fmt.Errorf("invalid containerFormat: %s", r.ContainerFormat)
	}

	return nil
}

// UploadImageResponse represents the response after uploading an image.
// Similar to snapshot creation, it returns the repository metadata and the tag
// generated (or updated) by the upload operation.
type UploadImageResponse struct {
	Repository *common.Repository `json:"repository"`
	Tag        *common.Tag        `json:"tag"`
}

// ListRepositoriesResponse represents the JSON response structure for listing repositories.
type ListRepositoriesResponse struct {
	Repositories []*common.Repository `json:"repositories"`
	Total        int                  `json:"total"`
}

// ListRepositoriesOptions represents options for listing repositories.
// It supports pagination with limit and offset, filtering with where conditions,
// and namespace specification for multi-tenant operations.
type ListRepositoriesOptions struct {
	Limit     int      // -1 for all, positive integer for limit
	Offset    int      // non-negative integer for offset
	Where     []string // Filter conditions (e.g., "namespace=public")
	Namespace string   // Namespace for X-Namespace header
}

// Validate validates the ListRepositoriesOptions.
// It ensures that limit and offset values are within acceptable ranges.
func (o *ListRepositoriesOptions) Validate() error {
	if o == nil {
		return fmt.Errorf("listRepositoriesOptions cannot be nil")
	}
	if o.Limit < -1 {
		return fmt.Errorf("limit must be >= -1")
	}
	if o.Offset < 0 {
		return fmt.Errorf("offset must be >= 0")
	}
	return nil
}
