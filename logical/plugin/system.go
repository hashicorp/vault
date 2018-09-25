package plugin

import (
	"context"
	"time"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/logical"
)

type SystemViewServer struct {
	impl logical.SystemView
}

func (s *SystemViewServer) DefaultLeaseTTL(_ interface{}, reply *DefaultLeaseTTLReply) error {
	ttl := s.impl.DefaultLeaseTTL()
	*reply = DefaultLeaseTTLReply{
		DefaultLeaseTTL: ttl,
	}

	return nil
}

func (s *SystemViewServer) MaxLeaseTTL(_ interface{}, reply *MaxLeaseTTLReply) error {
	ttl := s.impl.MaxLeaseTTL()
	*reply = MaxLeaseTTLReply{
		MaxLeaseTTL: ttl,
	}

	return nil
}

func (s *SystemViewServer) SudoPrivilege(args *SudoPrivilegeArgs, reply *SudoPrivilegeReply) error {
	sudo := s.impl.SudoPrivilege(context.Background(), args.Path, args.Token)
	*reply = SudoPrivilegeReply{
		Sudo: sudo,
	}

	return nil
}

func (s *SystemViewServer) Tainted(_ interface{}, reply *TaintedReply) error {
	tainted := s.impl.Tainted()
	*reply = TaintedReply{
		Tainted: tainted,
	}

	return nil
}

func (s *SystemViewServer) CachingDisabled(_ interface{}, reply *CachingDisabledReply) error {
	cachingDisabled := s.impl.CachingDisabled()
	*reply = CachingDisabledReply{
		CachingDisabled: cachingDisabled,
	}

	return nil
}

func (s *SystemViewServer) ReplicationState(_ interface{}, reply *ReplicationStateReply) error {
	replicationState := s.impl.ReplicationState()
	*reply = ReplicationStateReply{
		ReplicationState: replicationState,
	}

	return nil
}

func (s *SystemViewServer) ResponseWrapData(args *ResponseWrapDataArgs, reply *ResponseWrapDataReply) error {
	// Do not allow JWTs to be returned
	info, err := s.impl.ResponseWrapData(context.Background(), args.Data, args.TTL, false)
	if err != nil {
		*reply = ResponseWrapDataReply{
			Error: wrapError(err),
		}
		return nil
	}
	*reply = ResponseWrapDataReply{
		ResponseWrapInfo: info,
	}

	return nil
}

func (s *SystemViewServer) MlockEnabled(_ interface{}, reply *MlockEnabledReply) error {
	enabled := s.impl.MlockEnabled()
	*reply = MlockEnabledReply{
		MlockEnabled: enabled,
	}

	return nil
}

func (s *SystemViewServer) LocalMount(_ interface{}, reply *LocalMountReply) error {
	local := s.impl.LocalMount()
	*reply = LocalMountReply{
		Local: local,
	}

	return nil
}

func (s *SystemViewServer) EntityInfo(args *EntityInfoArgs, reply *EntityInfoReply) error {
	entity, err := s.impl.EntityInfo(args.EntityID)
	if err != nil {
		*reply = EntityInfoReply{
			Error: wrapError(err),
		}
		return nil
	}
	*reply = EntityInfoReply{
		Entity: entity,
	}

	return nil
}

func (s *SystemViewServer) PluginEnv(_ interface{}, reply *PluginEnvReply) error {
	pluginEnv, err := s.impl.PluginEnv(context.Background())
	if err != nil {
		*reply = PluginEnvReply{
			Error: wrapError(err),
		}
		return nil
	}
	*reply = PluginEnvReply{
		PluginEnvironment: pluginEnv,
	}

	return nil
}

type DefaultLeaseTTLReply struct {
	DefaultLeaseTTL time.Duration
}

type MaxLeaseTTLReply struct {
	MaxLeaseTTL time.Duration
}

type SudoPrivilegeArgs struct {
	Path  string
	Token string
}

type SudoPrivilegeReply struct {
	Sudo bool
}

type TaintedReply struct {
	Tainted bool
}

type CachingDisabledReply struct {
	CachingDisabled bool
}

type ReplicationStateReply struct {
	ReplicationState consts.ReplicationState
}

type ResponseWrapDataArgs struct {
	Data map[string]interface{}
	TTL  time.Duration
	JWT  bool
}

type ResponseWrapDataReply struct {
	ResponseWrapInfo *wrapping.ResponseWrapInfo
	Error            error
}

type MlockEnabledReply struct {
	MlockEnabled bool
}

type LocalMountReply struct {
	Local bool
}

type EntityInfoArgs struct {
	EntityID string
}

type EntityInfoReply struct {
	Entity *logical.Entity
	Error  error
}

type PluginEnvReply struct {
	PluginEnvironment *logical.PluginEnvironment
	Error             error
}
