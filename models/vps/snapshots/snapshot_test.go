package snapshots

import "testing"

func TestCreateSnapshotRequest_Validate(t *testing.T) {
	req := &CreateSnapshotRequest{}
	if err := req.Validate(); err == nil {
		t.Fatal("expected validation error for empty request, got nil")
	}

	req = &CreateSnapshotRequest{Name: "snap", VolumeID: "vol-1"}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateSnapshotRequest_Validate(t *testing.T) {
	req := &UpdateSnapshotRequest{}
	if err := req.Validate(); err == nil {
		t.Fatal("expected validation error for empty update request, got nil")
	}

	req = &UpdateSnapshotRequest{Name: "new-name"}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
