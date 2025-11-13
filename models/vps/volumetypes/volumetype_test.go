package volumetypes

import (
	"encoding/json"
	"testing"
)

func TestVolumeTypeListResponse_JSONUnmarshaling(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    []string
		wantErr bool
	}{
		{
			name:    "valid response with multiple types",
			json:    `{"volume_types": ["SSD", "HDD", "NVMe"]}`,
			want:    []string{"SSD", "HDD", "NVMe"},
			wantErr: false,
		},
		{
			name:    "valid response with single type",
			json:    `{"volume_types": ["SSD"]}`,
			want:    []string{"SSD"},
			wantErr: false,
		},
		{
			name:    "empty types array",
			json:    `{"volume_types": []}`,
			want:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response VolumeTypeListResponse
			err := json.Unmarshal([]byte(tt.json), &response)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(response.VolumeTypes) != len(tt.want) {
					t.Errorf("got %d types, want %d", len(response.VolumeTypes), len(tt.want))
					return
				}

				for i, volumeType := range response.VolumeTypes {
					if volumeType != tt.want[i] {
						t.Errorf("VolumeTypes[%d] = %v, want %v", i, volumeType, tt.want[i])
					}
				}
			}
		})
	}
}

func TestVolumeTypeListResponse_JSONMarshaling(t *testing.T) {
	response := VolumeTypeListResponse{
		VolumeTypes: []string{"SSD", "HDD", "NVMe"},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled VolumeTypeListResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(unmarshaled.VolumeTypes) != len(response.VolumeTypes) {
		t.Errorf("got %d types, want %d", len(unmarshaled.VolumeTypes), len(response.VolumeTypes))
	}

	for i, volumeType := range unmarshaled.VolumeTypes {
		if volumeType != response.VolumeTypes[i] {
			t.Errorf("VolumeTypes[%d] = %v, want %v", i, volumeType, response.VolumeTypes[i])
		}
	}
}
