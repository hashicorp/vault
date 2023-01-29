package node_status

import (
	"context"

	"github.com/hashicorp/hcp-link/pkg/nodestatus"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/vault/hcp_link/internal"
	"github.com/hashicorp/vault/vault/hcp_link/proto/node_status"
	"github.com/shirou/gopsutil/v3/host"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return status, err
	}

	listenerAddresses, err := c.NodeStatusGetter.ListenerAddresses()
	if err != nil {
		return status, err
	}

	logLevel, err := logging.ParseLogLevel(c.NodeStatusGetter.LogLevel())
	if err != nil {
		return status, err
	}

	raftStatus := &node_status.RaftStatus{}
	if sealStatus.StorageType == "raft" {
		raftStatus.IsVoter = c.NodeStatusGetter.IsRaftVoter()
	}

	protoRes := &node_status.LinkedClusterNodeStatusResponse{
		Type:                   sealStatus.Type,
		Initialized:            sealStatus.Initialized,
		Sealed:                 sealStatus.Sealed,
		T:                      int64(sealStatus.T),
		N:                      int64(sealStatus.N),
		Progress:               int64(sealStatus.Progress),
		Nonce:                  sealStatus.Nonce,
		Version:                sealStatus.Version,
		BuildDate:              sealStatus.BuildDate,
		Migration:              sealStatus.Migration,
		ClusterID:              sealStatus.ClusterID,
		ClusterName:            sealStatus.ClusterName,
		RecoverySeal:           sealStatus.RecoverySeal,
		StorageType:            sealStatus.StorageType,
		ReplicationState:       replState.StateStrings(),
		Hostname:               hostInfo.Hostname,
		ListenerAddresses:      listenerAddresses,
		OperatingSystem:        hostInfo.OS,
		OperatingSystemVersion: hostInfo.PlatformVersion,
		LogLevel:               node_status.LogLevel(logLevel),
		ActiveTime:             timestamppb.New(c.NodeStatusGetter.ActiveTime()),
		RaftStatus:             raftStatus,
	}

	status.StatusVersion = uint32(Version)
	status.Status = protoRes

	return status, nil
}
