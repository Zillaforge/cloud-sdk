package common_test

import (
	"encoding/json"
	"testing"

	"github.com/Zillaforge/cloud-sdk/models/iam/common"
)

func TestPermission_JSONParsing(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    common.Permission
		wantErr bool
	}{
		{
			name: "valid permission",
			json: `{"id": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36", "label": "DEFAULT"}`,
			want: common.Permission{
				ID:    "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
				Label: "DEFAULT",
			},
			wantErr: false,
		},
		{
			name: "permission with empty label",
			json: `{"id": "test-id", "label": ""}`,
			want: common.Permission{
				ID:    "test-id",
				Label: "",
			},
			wantErr: false,
		},
		{
			name: "permission with unknown fields",
			json: `{"id": "test-id", "label": "ADMIN", "unknownField": "value"}`,
			want: common.Permission{
				ID:    "test-id",
				Label: "ADMIN",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got common.Permission
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.ID != tt.want.ID || got.Label != tt.want.Label {
					t.Errorf("json.Unmarshal() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestTenantRole_String(t *testing.T) {
	tests := []struct {
		name string
		role common.TenantRole
		want string
	}{
		{
			name: "member role",
			role: common.TenantRoleMember,
			want: "TENANT_MEMBER",
		},
		{
			name: "admin role",
			role: common.TenantRoleAdmin,
			want: "TENANT_ADMIN",
		},
		{
			name: "owner role",
			role: common.TenantRoleOwner,
			want: "TENANT_OWNER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.String(); got != tt.want {
				t.Errorf("TenantRole.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTenantRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role common.TenantRole
		want bool
	}{
		{
			name: "valid member role",
			role: common.TenantRoleMember,
			want: true,
		},
		{
			name: "valid admin role",
			role: common.TenantRoleAdmin,
			want: true,
		},
		{
			name: "valid owner role",
			role: common.TenantRoleOwner,
			want: true,
		},
		{
			name: "invalid role",
			role: common.TenantRole("INVALID_ROLE"),
			want: false,
		},
		{
			name: "empty role",
			role: common.TenantRole(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("TenantRole.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTenantRole_JSONParsing(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    common.TenantRole
		wantErr bool
	}{
		{
			name:    "valid member role",
			json:    `"TENANT_MEMBER"`,
			want:    common.TenantRoleMember,
			wantErr: false,
		},
		{
			name:    "valid admin role",
			json:    `"TENANT_ADMIN"`,
			want:    common.TenantRoleAdmin,
			wantErr: false,
		},
		{
			name:    "valid owner role",
			json:    `"TENANT_OWNER"`,
			want:    common.TenantRoleOwner,
			wantErr: false,
		},
		{
			name:    "unknown role (forward compatibility)",
			json:    `"TENANT_SUPER_ADMIN"`,
			want:    common.TenantRole("TENANT_SUPER_ADMIN"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got common.TenantRole
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
