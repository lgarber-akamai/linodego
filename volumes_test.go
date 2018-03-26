package golinode

import (
	"testing"
)

func TestListVolumes(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	volumes, err := client.ListVolumes()
	if err != nil {
		t.Errorf("Error listing instances, expected struct, got error %v", err)
	}
	if len(volumes) != 1 {
		t.Errorf("Expected a list of instances, but got %v", volumes)
	}
}

func TestGetVolume(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	_, err = client.GetVolume(4880)
	if err != nil {
		t.Errorf("Error getting volume 4880, expected *LinodeVolume, got error %v", err)
	}
}