package vault

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/jsonutil"
)

const (
	// Storage path where the local cluster name and identifier are stored
	coreClusterLocalPath = "core/cluster/local"

	// Storage path where the global cluster name and identifier are stored
	coreClusterGlobalPath = "core/cluster/global"
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
func (c *Core) Cluster(isLocal bool) (*Cluster, error) {
	var key string
	if isLocal {
		key = coreClusterLocalPath
	} else {
		key = coreClusterGlobalPath
	}

	// Fetch the storage entry. This call fails when Vault is sealed.
	entry, err := c.barrier.Get(key)
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
// Entries will be created only if they are not already present.
func (c *Core) setupCluster() error {
	// Create or store a local name and local cluster ID, if not already stored
	if err := c.setCluster(true, c.localClusterName, c.localClusterID); err != nil {
		return err
	}

	// Create or store a global name and global cluster ID, if not already stored
	if err := c.setCluster(false, c.globalClusterName, c.globalClusterID); err != nil {
		return err
	}
	return nil
}

// setCluster creates a local or global storage index for a set of cluster
// information. If cluster name or cluster ID is not supplied, this method will
// auto-generate them respectively.
func (c *Core) setCluster(isLocal bool, clusterName, clusterID string) error {
	// Check if storage index is already present or not
	cluster, err := c.Cluster(isLocal)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to get cluster details: %v", err)
		return err
	}
	if cluster != nil {
		// If index is already present, don't update it
		return nil
	}

	// If clusterName is not supplied, generate one
	if clusterName == "" {
		clusterNameBytes, err := uuid.GenerateRandomBytes(4)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to generate cluster name: %v", err)
			return err
		}
		prefix := "vault-"
		if isLocal {
			prefix += "local-"
		} else {
			prefix += "global-"
		}
		clusterName = fmt.Sprintf("%s%08x", prefix, clusterNameBytes)
	}

	// If clusterName is not supplied, generate one
	if clusterID == "" {
		var err error
		clusterID, err = uuid.GenerateUUID()
		if err != nil {
			c.logger.Printf("[ERR] core: failed to generate cluster identifier: %v", err)
			return err
		}
	}

	// Encode the cluster information into as a JSON string
	rawCluster, err := json.Marshal(&Cluster{
		Name: clusterName,
		ID:   clusterID,
	})
	if err != nil {
		c.logger.Printf("[ERR] core: failed to encode cluster details: %v", err)
		return err
	}

	// Determine the storage path
	var key string
	if isLocal {
		key = coreClusterLocalPath
	} else {
		key = coreClusterGlobalPath
	}

	// Store it
	return c.barrier.Put(&Entry{
		Key:   key,
		Value: rawCluster,
	})
}
