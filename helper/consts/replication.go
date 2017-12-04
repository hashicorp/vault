package consts

type ReplicationState uint32

const (
	_ ReplicationState = iota
	OldReplicationPrimary
	OldReplicationSecondary
	OldReplicationBootstrapping

	ReplicationDisabled           ReplicationState = 0
	ReplicationPerformancePrimary ReplicationState = 1 << iota
	ReplicationPerformanceSecondary
	ReplicationBootstrapping
	ReplicationDRPrimary
	ReplicationDRSecondary
)

func (r ReplicationState) String() string {
	switch r {
	case ReplicationPerformanceSecondary:
		return "perf-secondary"
	case ReplicationPerformancePrimary:
		return "perf-primary"
	case ReplicationBootstrapping:
		return "bootstrapping"
	case ReplicationDRPrimary:
		return "dr-primary"
	case ReplicationDRSecondary:
		return "dr-secondary"
	}

	return "disabled"
}

func (r ReplicationState) GetDRString() string {
	switch {
	case r.HasState(ReplicationDRPrimary):
		return ReplicationDRPrimary.String()
	case r.HasState(ReplicationDRSecondary):
		return ReplicationDRSecondary.String()
	default:
		return ReplicationDisabled.String()
	}
}

func (r ReplicationState) GetPerformanceString() string {
	switch {
	case r.HasState(ReplicationPerformancePrimary):
		return ReplicationPerformancePrimary.String()
	case r.HasState(ReplicationPerformanceSecondary):
		return ReplicationPerformanceSecondary.String()
	default:
		return ReplicationDisabled.String()
	}
}

func (r ReplicationState) HasState(flag ReplicationState) bool { return r&flag != 0 }
func (r *ReplicationState) AddState(flag ReplicationState)     { *r |= flag }
func (r *ReplicationState) ClearState(flag ReplicationState)   { *r &= ^flag }
func (r *ReplicationState) ToggleState(flag ReplicationState)  { *r ^= flag }
