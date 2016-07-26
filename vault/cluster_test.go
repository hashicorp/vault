package vault

import "testing"

func TestCluster(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	cluster, err := c.Cluster(true)
	if err != nil {
		t.Fatal(err)
	}
	if cluster == nil || cluster.Name == "" || cluster.ID == "" {
		t.Fatalf("local cluster information missing: cluster:%#v", cluster)
	}

	cluster, err = c.Cluster(false)
	if err != nil {
		t.Fatal(err)
	}
	if cluster == nil || cluster.Name == "" || cluster.ID == "" {
		t.Fatalf("global cluster information missing: cluster:%#v", cluster)
	}
}
