// +build !enterprise

package replication

import "github.com/hashicorp/vault/helper/consts"

type Cluster struct {
	State              consts.ReplicationState
	ClusterID          string
	PrimaryClusterAddr string
}

type Clusters struct {
	DR          *Cluster
	Performance *Cluster
}
