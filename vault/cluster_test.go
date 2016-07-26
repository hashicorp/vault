package vault

import "testing"

func TestCluster(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	cluster, err := c.Cluster()
	if err != nil {
		t.Fatal(err)
	}
	if cluster == nil || cluster.Name == "" || cluster.ID == "" {
		t.Fatalf("cluster information missing: cluster:%#v", cluster)
	}
}
