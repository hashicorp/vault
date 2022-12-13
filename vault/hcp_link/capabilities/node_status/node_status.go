package node_status

import (
	"context"

	"github.com/hashicorp/hcp-link/pkg/nodestatus"
	"github.com/hashicorp/vault/vault/hcp_link/internal"
	"github.com/hashicorp/vault/vault/hcp_link/proto/node_status"
)

var (
	_       nodestatus.Reporter = &NodeStatusReporter{}
	Version                     = 1
)

type NodeStatusReporter struct {
	NodeStatusGetter internal.WrappedCoreNodeStatus
}

func (c *NodeStatusReporter) GetNodeStatus(ctx context.Context) (nodestatus.NodeStatus, error) {
	var status nodestatus.NodeStatus

	sealStatus, err := c.NodeStatusGetter.GetSealStatus(ctx)
	if err != nil {
		return status, err
	}

	replState := c.NodeStatusGetter.ReplicationState()

	protoRes := &node_status.LinkedClusterNodeStatusResponse{
		Type:             sealStatus.Type,
		Initialized:      sealStatus.Initialized,
		Sealed:           sealStatus.Sealed,
		T:                int64(sealStatus.T),
		N:                int64(sealStatus.N),
		Progress:         int64(sealStatus.Progress),
		Nonce:            sealStatus.Nonce,
		Version:          sealStatus.Version,
		BuildDate:        sealStatus.BuildDate,
		Migration:        sealStatus.Migration,
		ClusterID:        sealStatus.ClusterID,
		ClusterName:      sealStatus.ClusterName,
		RecoverySeal:     sealStatus.RecoverySeal,
		StorageType:      sealStatus.StorageType,
		ReplicationState: replState.StateStrings(),
	}

	ns := nodestatus.NodeStatus{
		StatusVersion: uint32(Version),
		Status:        protoRes,
	}

	return ns, nil
}
