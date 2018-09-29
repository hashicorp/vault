// +build !enterprise

package vault

import "github.com/hashicorp/vault/helper/consts"

type ReplicatedCluster struct {
	State              consts.ReplicationState
	ClusterID          string
	PrimaryClusterAddr string
}
