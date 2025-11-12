package common

import (
	"fmt"
	"strings"
)

// IDName represents a lightweight reference to an entity with id and name.
// Used for Creator, Project, User references throughout the VRM API.
type IDName struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Account     string `json:"account,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// Validate validates the IDName structure.
// At minimum, ID must be present.
func (n *IDName) Validate() error {
	if n == nil {
		return fmt.Errorf("IDName cannot be nil")
	}
	if strings.TrimSpace(n.ID) == "" {
		return fmt.Errorf("IDName.ID is required and must not be empty")
	}
	return nil
}

// DiskFormat represents the disk image format for virtual machine images.
// Valid values include ami, ari, aki, vhd, vmdk, raw, qcow2, vdi, and iso.
type DiskFormat string

const (
	DiskFormatAMI   DiskFormat = "ami"
	DiskFormatARI   DiskFormat = "ari"
	DiskFormatAKI   DiskFormat = "aki"
	DiskFormatVHD   DiskFormat = "vhd"
	DiskFormatVMDK  DiskFormat = "vmdk"
	DiskFormatRaw   DiskFormat = "raw"
	DiskFormatQcow2 DiskFormat = "qcow2"
	DiskFormatVDI   DiskFormat = "vdi"
	DiskFormatISO   DiskFormat = "iso"
)

// IsValid checks if the DiskFormat is a valid value.
func (f DiskFormat) IsValid() bool {
	switch f {
	case DiskFormatAMI, DiskFormatARI, DiskFormatAKI, DiskFormatVHD, DiskFormatVMDK,
		DiskFormatRaw, DiskFormatQcow2, DiskFormatVDI, DiskFormatISO:
		return true
	}
	return false
}

// String returns the string representation of DiskFormat.
func (f DiskFormat) String() string {
	return string(f)
}

// ContainerFormat represents the container image format for virtual machine images.
// Valid values include ami, ari, aki, bare, and ovf.
type ContainerFormat string

const (
	ContainerFormatAMI  ContainerFormat = "ami"
	ContainerFormatARI  ContainerFormat = "ari"
	ContainerFormatAKI  ContainerFormat = "aki"
	ContainerFormatBare ContainerFormat = "bare"
	ContainerFormatOVF  ContainerFormat = "ovf"
)

// IsValid checks if the ContainerFormat is a valid value.
func (f ContainerFormat) IsValid() bool {
	switch f {
	case ContainerFormatAMI, ContainerFormatARI, ContainerFormatAKI, ContainerFormatBare, ContainerFormatOVF:
		return true
	}
	return false
}

// String returns the string representation of ContainerFormat.
func (f ContainerFormat) String() string {
	return string(f)
}
