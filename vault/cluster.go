package vault

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/jsonutil"
)

const (
	// Storage path where the local cluster name and identifier are stored
	coreClusterPath = "core/cluster/local"
)

// Structure representing the storage entry that holds cluster information
type Cluster struct {
	// Name of the cluster
	Name string `json:"name" structs:"name" mapstructure:"name"`

	// Identifier of the cluster
	ID string `json:"id" structs:"id" mapstructure:"id"`
}

// Cluster fetches the details of either local or global cluster based on the
// input. This method errors out when Vault is sealed.
func (c *Core) Cluster() (*Cluster, error) {
	// Fetch the storage entry. This call fails when Vault is sealed.
	entry, err := c.barrier.Get(coreClusterPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	// Decode the cluster information
	var cluster Cluster
	if err = jsonutil.DecodeJSON(entry.Value, &cluster); err != nil {
		return nil, fmt.Errorf("failed to decode cluster details: %v", err)
	}

	return &cluster, nil
}

// setupCluster creates storage entries for holding Vault cluster information.
// Entries will be created only if they are not already present. If clusterName
// is not supplied, this method will auto-generate it.
func (c *Core) setupCluster() error {
	// Check if storage index is already present or not
	cluster, err := c.Cluster()
	if err != nil {
		c.logger.Printf("[ERR] core: failed to get cluster details: %v", err)
		return err
	}
	if cluster != nil {
		// If index is already present, don't update it
		return nil
	}

	// If cluster name is not supplied, generate one
	if c.clusterName == "" {
		clusterNameBytes, err := uuid.GenerateRandomBytes(4)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to generate cluster name: %v", err)
			return err
		}
		c.clusterName = fmt.Sprintf("vaule-cluster-%08x", clusterNameBytes)
	}

	// Generate a clusterID
	var err error
	clusterID, err := uuid.GenerateUUID()
	if err != nil {
		c.logger.Printf("[ERR] core: failed to generate cluster identifier: %v", err)
		return err
	}

	// Encode the cluster information into as a JSON string
	rawCluster, err := json.Marshal(&Cluster{
		Name: c.clusterName,
		ID:   clusterID,
	})
	if err != nil {
		c.logger.Printf("[ERR] core: failed to encode cluster details: %v", err)
		return err
	}

	// Store it
	return c.barrier.Put(&Entry{
		Key:   coreClusterPath,
		Value: rawCluster,
	})

}
