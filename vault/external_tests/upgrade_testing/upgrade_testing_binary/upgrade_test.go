package upgrade_testing_binary

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
)

func DockerPrivilegedNamespaceClient(t *testing.T, nodeCount int) *docker.DockerCluster {
	os.Setenv("VAULT_NEW_BINARY", "/Users/divya.chandrasekaran/code/vault-enterprise/bin/vault")
	os.Setenv("VAULT_BINARY", "/Users/divya.chandrasekaran/vault-old/vault")
	binary := os.Getenv("VAULT_BINARY")
	new_binary := os.Getenv("VAULT_NEW_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}
	opts := &docker.DockerClusterOptions{
		ImageRepo: "hashicorp/vault",
		// We're replacing the binary anyway, so we're not too particular about
		// the docker image version tag.
		ImageTag:       "latest",
		VaultBinary:    binary,
		VaultBinaryNew: new_binary,
		ClusterOptions: testcluster.ClusterOptions{
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				LogLevel: "TRACE",
				// If you want the test to run faster locally, you could
				// uncomment this performance_multiplier change.
				StorageOptions: map[string]string{
					"performance_multiplier": "1",
				},
			},
			AdministrativeNamespacePath: "admin/",
			NumCores:                    nodeCount,
		},
	}
	cluster := docker.NewTestDockerCluster(t, opts)
	return cluster
}

func TestQuotas_LeaseCount_Remount(t *testing.T) {
	cluster := DockerPrivilegedNamespaceClient(t, 1)
	defer cluster.Cleanup()
	// rootClient := cluster.Nodes()[0].APIClient()

	//// Create privileged namespace
	//_, err := rootClient.Logical().Write(fmt.Sprintf("sys/namespaces/%s", "admin"), nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//// Read privileged namespace
	//_, err = rootClient.Logical().Read(fmt.Sprintf("sys/namespaces/%s", "admin"))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//// Restart docker cluster
	//
	//// Read privileged namespace
	//_, err = rootClient.Logical().Read(fmt.Sprintf("sys/namespaces/%s", "admin"))
	//if err != nil {
	//	t.Fatal(err)
	//}
}
