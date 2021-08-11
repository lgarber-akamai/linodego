package integration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
)

var testLKEClusterPoolCreateOpts = linodego.LKEClusterPoolCreateOptions{
	Type:  "g6-standard-2",
	Count: 2,
	Disks: []linodego.LKEClusterPoolDisk{
		{
			Size: 1000,
			Type: "ext4",
		},
	},
	Tags: []string{"testing"},
}

func TestGetLKEClusterPool_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetLKEClusterPool_missing")
	defer teardown()

	i, err := client.GetLKEClusterPool(context.Background(), 0, 0)
	if err == nil {
		t.Errorf("should have received an error requesting a missing lkeClusterPool, got %v", i)
	}
	e, ok := err.(*linodego.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing lkeClusterPool, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing lkeClusterPool, got %v", e.Code)
	}
}

func TestGetLKEClusterPool_found(t *testing.T) {
	client, lkeCluster, pool, teardown, err := setupLKEClusterPool(t, "fixtures/TestGetLKEClusterPool_found")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	i, err := client.GetLKEClusterPool(context.Background(), lkeCluster.ID, pool.ID)
	if err != nil {
		t.Errorf("Error getting lkeClusterPool, expected struct, got %v and error %v", i, err)
	}

	if i.ID != pool.ID {
		t.Errorf("Expected a specific lkeClusterPool, but got a different one %v", i)
	}
	if i.Count != 2 {
		t.Errorf("expected count to be 2; got %d", i.Count)
	}
	if i.Type != "g6-standard-2" {
		t.Errorf("expected type to be g6-standard-2; got %s", i.Type)
	}
	if diff := cmp.Diff(linodego.LKEClusterPoolAutoscaler{
		Min:     2,
		Max:     2,
		Enabled: false,
	}, pool.Autoscaler); diff != "" {
		t.Errorf("unexpected autoscaler:\n%s", diff)
	}
	if diff := cmp.Diff([]string{"testing"}, i.Tags); diff != "" {
		t.Errorf("unexpected tags:\n%s", diff)
	}
}

func TestListLKEClusterPools(t *testing.T) {
	client, lkeCluster, _, teardown, err := setupLKEClusterPool(t, "fixtures/TestListLKEClusterPools")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	i, err := client.ListLKEClusterPools(context.Background(), lkeCluster.ID, nil)
	if err != nil {
		t.Errorf("Error listing lkeClusterPools, expected struct, got error %v", err)
	}
	if len(i) != 2 {
		t.Errorf("Expected two lkeClusterPools, but got %#v", i)
	}
}

func TestDeleteLKEClusterPoolNode(t *testing.T) {
	client, lkeCluster, clusterPool, teardown, err := setupLKEClusterPool(t, "fixtures/TestDeleteLKEClusterPoolNode")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	linodes := clusterPool.Linodes
	err = client.DeleteLKEClusterPoolNode(context.TODO(), lkeCluster.ID, linodes[0].ID)
	if err != nil {
		t.Errorf("failed to delete node %q: %s", linodes[0].ID, err)
	}

	clusterPool, err = client.GetLKEClusterPool(context.TODO(), lkeCluster.ID, clusterPool.ID)
	if err != nil {
		t.Errorf("failed to get updated node pool: %s", err)
	}

	if !(len(clusterPool.Linodes) == 1 && clusterPool.Linodes[0].ID == linodes[1].ID) {
		t.Errorf("expected cluster pool to have 1 linode (%s); got %v", linodes[1].ID, clusterPool.Linodes)
	}
}

func TestUpdateLKEClusterPool(t *testing.T) {
	client, lkeCluster, clusterPool, teardown, err := setupLKEClusterPool(t, "fixtures/TestUpdateLKEClusterPool")
	if err != nil {
		t.Error(err)
	}
	defer teardown()

	updatedAutoscaler := linodego.LKEClusterPoolAutoscaler{
		Enabled: true,
		Min:     2,
		Max:     5,
	}
	updatedTags := []string{}
	updated, err := client.UpdateLKEClusterPool(context.TODO(), lkeCluster.ID, clusterPool.ID, linodego.LKEClusterPoolUpdateOptions{
		Count:      2,            // downsize
		Tags:       &updatedTags, // remove all tags
		Autoscaler: &updatedAutoscaler,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(updated.Tags) != 0 {
		t.Errorf("expected tags to be empty; got %v", updated.Tags)
	}
	if updated.Count != 2 {
		t.Errorf("expected count to be 2; got %d", updated.Count)
	}

	updatedTags = []string{"bar", "foo", "test"}
	updated, err = client.UpdateLKEClusterPool(context.TODO(), lkeCluster.ID, clusterPool.ID, linodego.LKEClusterPoolUpdateOptions{
		Count: 3,            // upsize
		Tags:  &updatedTags, // repopulate tags
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(updatedTags, updated.Tags); diff != "" {
		t.Errorf("unexpected tags:\n%s", diff)
	}
	if diff := cmp.Diff(updatedAutoscaler, updated.Autoscaler); diff != "" {
		t.Errorf("unexpected autoscaler:\n%s", diff)
	}
	if updated.Count != 3 {
		t.Errorf("expected count to be 3; got %d", updated.Count)
	}
}

func setupLKEClusterPool(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.LKECluster, *linodego.LKEClusterPool, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, lkeCluster, fixtureTeardown, err := setupLKECluster(t, []clusterModifier{func(createOpts *linodego.LKEClusterCreateOptions) {
		createOpts.Label = randString(12, lowerBytes, digits) + "-linodego-testing"
	}}, fixturesYaml)
	if err != nil {
		t.Errorf("Error creating lkeCluster, got error %v", err)
	}

	pool, err := client.CreateLKEClusterPool(context.Background(), lkeCluster.ID, testLKEClusterPoolCreateOpts)
	if err != nil {
		t.Errorf("Error creating LKECluster Pool, got error %v", err)
	}

	teardown := func() {
		// delete the LKEClusterPool to exercise the code
		if err := client.DeleteLKEClusterPool(context.Background(), lkeCluster.ID, pool.ID); err != nil {
			t.Errorf("Expected to delete a LKECluster Pool, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, lkeCluster, pool, teardown, err
}
