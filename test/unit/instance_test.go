package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListInstances(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("linodes_list")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances", fixtureData)

	instances, err := base.Client.ListInstances(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing instances: %v", err)
	}

	assert.Equal(t, 1, len(instances))
	linode := instances[0]
	assert.Equal(t, 123, linode.ID)
	assert.Equal(t, "linode123", linode.Label)
	assert.Equal(t, "running", string(linode.Status))
	assert.Equal(t, "203.0.113.1", linode.IPv4[0].String())
	assert.Equal(t, "g6-standard-1", linode.Type)
	assert.Equal(t, "us-east", linode.Region)
	assert.Equal(t, 4096, linode.Specs.Memory)
}
