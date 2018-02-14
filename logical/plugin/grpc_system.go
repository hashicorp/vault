package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"google.golang.org/grpc"

	"fmt"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/pb"
)

func newGRPCSystemView(conn *grpc.ClientConn) *gRPCSystemViewClient {
	return &gRPCSystemViewClient{
		client: pb.NewSystemViewClient(conn),
	}
}

type gRPCSystemViewClient struct {
	client pb.SystemViewClient
}

func (s *gRPCSystemViewClient) DefaultLeaseTTL() time.Duration {
	reply, err := s.client.DefaultLeaseTTL(context.Background(), &pb.Empty{})
	if err != nil {
		return 0
	}

	return time.Duration(reply.TTL)
}

func (s *gRPCSystemViewClient) MaxLeaseTTL() time.Duration {
	reply, err := s.client.MaxLeaseTTL(context.Background(), &pb.Empty{})
	if err != nil {
		return 0
	}

	return time.Duration(reply.TTL)
}

func (s *gRPCSystemViewClient) SudoPrivilege(ctx context.Context, path string, token string) bool {
	reply, err := s.client.SudoPrivilege(ctx, &pb.SudoPrivilegeArgs{
		Path:  path,
		Token: token,
	})
	if err != nil {
		return false
	}

	return reply.Sudo
}

func (s *gRPCSystemViewClient) Tainted() bool {
	reply, err := s.client.Tainted(context.Background(), &pb.Empty{})
	if err != nil {
		return false
	}

	return reply.Tainted
}

func (s *gRPCSystemViewClient) CachingDisabled() bool {
	reply, err := s.client.CachingDisabled(context.Background(), &pb.Empty{})
	if err != nil {
		return false
	}

	return reply.Disabled
}

func (s *gRPCSystemViewClient) ReplicationState() consts.ReplicationState {
	reply, err := s.client.ReplicationState(context.Background(), &pb.Empty{})
	if err != nil {
		return consts.ReplicationUnknown
	}

	return consts.ReplicationState(reply.State)
}

func (s *gRPCSystemViewClient) ResponseWrapData(ctx context.Context, data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reply, err := s.client.ResponseWrapData(ctx, &pb.ResponseWrapDataArgs{
		Data: string(buf[:]),
		TTL:  int64(ttl),
		JWT:  false,
	})
	if err != nil {
		return nil, err
	}
	if reply.Err != "" {
		return nil, errors.New(reply.Err)
	}

	info, err := pb.ProtoResponseWrapInfoToLogicalResponseWrapInfo(reply.WrapInfo)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *gRPCSystemViewClient) LookupPlugin(ctx context.Context, name string) (*pluginutil.PluginRunner, error) {
	return nil, fmt.Errorf("cannot call LookupPlugin from a plugin backend")
}

func (s *gRPCSystemViewClient) MlockEnabled() bool {
	reply, err := s.client.MlockEnabled(context.Background(), &pb.Empty{})
	if err != nil {
		return false
	}

	return reply.Enabled
}

func (s *gRPCSystemViewClient) LocalMount() bool {
	reply, err := s.client.LocalMount(context.Background(), &pb.Empty{})
	if err != nil {
		return false
	}

	return reply.Local
}

type gRPCSystemViewServer struct {
	impl logical.SystemView
}

func (s *gRPCSystemViewServer) DefaultLeaseTTL(ctx context.Context, _ *pb.Empty) (*pb.TTLReply, error) {
	ttl := s.impl.DefaultLeaseTTL()
	return &pb.TTLReply{
		TTL: int64(ttl),
	}, nil
}

func (s *gRPCSystemViewServer) MaxLeaseTTL(ctx context.Context, _ *pb.Empty) (*pb.TTLReply, error) {
	ttl := s.impl.MaxLeaseTTL()
	return &pb.TTLReply{
		TTL: int64(ttl),
	}, nil
}

func (s *gRPCSystemViewServer) SudoPrivilege(ctx context.Context, args *pb.SudoPrivilegeArgs) (*pb.SudoPrivilegeReply, error) {
	sudo := s.impl.SudoPrivilege(ctx, args.Path, args.Token)
	return &pb.SudoPrivilegeReply{
		Sudo: sudo,
	}, nil
}

func (s *gRPCSystemViewServer) Tainted(ctx context.Context, _ *pb.Empty) (*pb.TaintedReply, error) {
	tainted := s.impl.Tainted()
	return &pb.TaintedReply{
		Tainted: tainted,
	}, nil
}

func (s *gRPCSystemViewServer) CachingDisabled(ctx context.Context, _ *pb.Empty) (*pb.CachingDisabledReply, error) {
	cachingDisabled := s.impl.CachingDisabled()
	return &pb.CachingDisabledReply{
		Disabled: cachingDisabled,
	}, nil
}

func (s *gRPCSystemViewServer) ReplicationState(ctx context.Context, _ *pb.Empty) (*pb.ReplicationStateReply, error) {
	replicationState := s.impl.ReplicationState()
	return &pb.ReplicationStateReply{
		State: int32(replicationState),
	}, nil
}

func (s *gRPCSystemViewServer) ResponseWrapData(ctx context.Context, args *pb.ResponseWrapDataArgs) (*pb.ResponseWrapDataReply, error) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(args.Data), &data)
	if err != nil {
		return &pb.ResponseWrapDataReply{}, err
	}

	// Do not allow JWTs to be returned
	info, err := s.impl.ResponseWrapData(ctx, data, time.Duration(args.TTL), false)
	if err != nil {
		return &pb.ResponseWrapDataReply{
			Err: pb.ErrToString(err),
		}, nil
	}

	pbInfo, err := pb.LogicalResponseWrapInfoToProtoResponseWrapInfo(info)
	if err != nil {
		return &pb.ResponseWrapDataReply{}, err
	}

	return &pb.ResponseWrapDataReply{
		WrapInfo: pbInfo,
	}, nil
}

func (s *gRPCSystemViewServer) MlockEnabled(ctx context.Context, _ *pb.Empty) (*pb.MlockEnabledReply, error) {
	enabled := s.impl.MlockEnabled()
	return &pb.MlockEnabledReply{
		Enabled: enabled,
	}, nil
}

func (s *gRPCSystemViewServer) LocalMount(ctx context.Context, _ *pb.Empty) (*pb.LocalMountReply, error) {
	local := s.impl.LocalMount()
	return &pb.LocalMountReply{
		Local: local,
	}, nil
}
