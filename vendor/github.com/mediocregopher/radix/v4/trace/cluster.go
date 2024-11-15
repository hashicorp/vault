package trace

import "context"

// ClusterTrace contains callbacks which can be triggered for specific events
// during a Cluster's runtime.
//
// All callbacks are called synchronously.
type ClusterTrace struct {
	// StateChange is called when the Cluster becomes down or becomes available
	// again.
	StateChange func(ClusterStateChange)

	// TopoChanged is called when the Cluster's topology changes.
	TopoChanged func(ClusterTopoChanged)

	// Redirected is called when redis responds to an Action with a 'MOVED' or
	// 'ASK' error.
	Redirected func(ClusterRedirected)

	// InternalError is called whenever the Cluster encounters an error which is
	// not otherwise communicated to the user.
	InternalError func(ClusterInternalError)
}

// ClusterStateChange is passed into the ClusterTrace.StateChange callback
// whenever the Cluster's state has changed.
type ClusterStateChange struct {
	IsDown bool
}

// ClusterNodeInfo describes the attributes of a node in a redis cluster's
// topology.
type ClusterNodeInfo struct {
	Addr      string
	Slots     [][2]uint16
	IsPrimary bool
}

// ClusterTopoChanged is passed into the ClusterTrace.TopoChanged callback
// whenever the Cluster's topology has changed.
type ClusterTopoChanged struct {
	Added   []ClusterNodeInfo
	Removed []ClusterNodeInfo
	Changed []ClusterNodeInfo
}

// ClusterRedirected is passed into the ClusterTrace.Redirected callback
// whenever redis responds to an Action with a 'MOVED' or 'ASK' error.
type ClusterRedirected struct {
	// Context is the Context passed into the Do call which is performing the
	// Action which received a MOVED/ASK error.
	Context context.Context

	// Addr is the address of the redis instance the Action was performed
	// against.
	Addr string

	// Key is the key that the Action would operate on.
	Key string

	// Moved and Ask denote which kind of error was returned. One will be true.
	Moved, Ask bool

	// RedirectCount denotes how many times the Action has been redirected so
	// far.
	RedirectCount int

	// Final indicates that the MOVED/ASK error which was received will not be
	// honored, and the call to Do will be returning the MOVED/ASK error.
	Final bool
}

// ClusterInternalError is passed into the ClusterTrace.InternalError callback
// whenever Cluster encounters an error which is not otherwise communicated to
// the user.
type ClusterInternalError struct {
	Err error
}
