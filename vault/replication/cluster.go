// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package replication

import "github.com/hashicorp/vault/sdk/helper/consts"

type Cluster struct {
	State              consts.ReplicationState
	ClusterID          string
	PrimaryClusterAddr string
}

type Clusters struct {
	DR          *Cluster
	Performance *Cluster
}
