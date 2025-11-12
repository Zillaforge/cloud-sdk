package common

import (
	"testing"
)

func TestIDNameValidate(t *testing.T) {
	tests := []struct {
		name    string
		idname  *IDName
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid IDName",
			idname: &IDName{
				ID:   "test-id-123",
				Name: "test",
			},
			wantErr: false,
		},
		{
			name: "valid IDName with all fields",
			idname: &IDName{
				ID:          "user-456",
				Name:        "john",
				Account:     "john@example.com",
				DisplayName: "John Doe",
			},
			wantErr: false,
		},
		{
			name: "valid IDName with only ID",
			idname: &IDName{
				ID: "repo-789",
			},
			wantErr: false,
		},
		{
			name:    "nil IDName",
			idname:  nil,
			wantErr: true,
			errMsg:  "IDName cannot be nil",
		},
		{
			name: "empty ID",
			idname: &IDName{
				ID:   "",
				Name: "test",
			},
			wantErr: true,
			errMsg:  "IDName.ID is required and must not be empty",
		},
		{
			name: "whitespace-only ID",
			idname: &IDName{
				ID:   "   ",
				Name: "test",
			},
			wantErr: true,
			errMsg:  "IDName.ID is required and must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.idname.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestDiskFormatIsValid(t *testing.T) {
	tests := []struct {
		format DiskFormat
		valid  bool
	}{
		{DiskFormatAMI, true},
		{DiskFormatARI, true},
		{DiskFormatAKI, true},
		{DiskFormatVHD, true},
		{DiskFormatVMDK, true},
		{DiskFormatRaw, true},
		{DiskFormatQcow2, true},
		{DiskFormatVDI, true},
		{DiskFormatISO, true},
		{DiskFormat("invalid"), false},
		{DiskFormat(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestContainerFormatIsValid(t *testing.T) {
	tests := []struct {
		format ContainerFormat
		valid  bool
	}{
		{ContainerFormatAMI, true},
		{ContainerFormatARI, true},
		{ContainerFormatAKI, true},
		{ContainerFormatBare, true},
		{ContainerFormatOVF, true},
		{ContainerFormat("invalid"), false},
		{ContainerFormat(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestDiskFormatString(t *testing.T) {
	tests := []struct {
		format   DiskFormat
		expected string
	}{
		{DiskFormatAMI, "ami"},
		{DiskFormatQcow2, "qcow2"},
		{DiskFormatISO, "iso"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContainerFormatString(t *testing.T) {
	tests := []struct {
		format   ContainerFormat
		expected string
	}{
		{ContainerFormatAMI, "ami"},
		{ContainerFormatBare, "bare"},
		{ContainerFormatOVF, "ovf"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
