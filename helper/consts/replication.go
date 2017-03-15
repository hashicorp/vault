package consts

type ReplicationState uint32

const (
	ReplicationDisabled ReplicationState = iota
	ReplicationPrimary
	ReplicationSecondary
)

func (r ReplicationState) String() string {
	switch r {
	case ReplicationSecondary:
		return "secondary"
	case ReplicationPrimary:
		return "primary"
	}

	return "disabled"
}
