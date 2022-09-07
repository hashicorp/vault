package consts

const VaultAllowPendingRemovalMountsEnv = "VAULT_ALLOW_PENDING_REMOVAL_MOUNTS"

// DeprecationStatus represents the current deprecation state for builtins
type DeprecationStatus uint32

// These are the states of deprecation for builtin plugins
const (
	Supported = iota
	Deprecated
	PendingRemoval
	Removed
	Unknown
)

// String returns the string representation of a builtin deprecation status
func (s DeprecationStatus) String() string {
	switch s {
	case Supported:
		return "supported"
	case Deprecated:
		return "deprecated"
	case PendingRemoval:
		return "pending removal"
	case Removed:
		return "removed"
	default:
		return ""
	}
}
