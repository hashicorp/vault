package trace

// SentinelTrace contains callbacks which can be triggered for specific events
// during a Sentinel's runtime.
//
// All callbacks are called synchronously.
type SentinelTrace struct {
	// TopoChanged is called when the Sentinel's replica set's topology changes.
	TopoChanged func(SentinelTopoChanged)

	// InternalError is called whenever the Sentinel encounters an error which
	// is not otherwise communicated to the user.
	InternalError func(SentinelInternalError)
}

// SentinelNodeInfo describes the attributes of a node in a sentinel replica
// set's topology.
type SentinelNodeInfo struct {
	Addr      string
	IsPrimary bool
}

// SentinelTopoChanged is passed into the SentinelTrace.TopoChanged callback
// whenever the Sentinel's replica set's topology has changed.
type SentinelTopoChanged struct {
	Added   []SentinelNodeInfo
	Removed []SentinelNodeInfo
	Changed []SentinelNodeInfo
}

// SentinelInternalError is passed into the SentinelTrace.InternalError callback
// whenever Sentinel encounters an error which is not otherwise communicated to
// the user.
type SentinelInternalError struct {
	Err error
}
