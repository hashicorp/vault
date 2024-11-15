package trace

// PersistentPubSubTrace contains callbacks which can be triggered for specific
// events during a persistent PubSubConn's runtime.
//
// All callbacks are called synchronously.
type PersistentPubSubTrace struct {
	// InternalError is called whenever the PersistentPubSub encounters an error
	// which is not otherwise communicated to the user.
	InternalError func(PersistentPubSubInternalError)
}

// PersistentPubSubInternalError is passed into the
// PersistentPubSubTrace.InternalError callback whenever PersistentPubSub
// encounters an error which is not otherwise communicated to the user.
type PersistentPubSubInternalError struct {
	Err error
}
