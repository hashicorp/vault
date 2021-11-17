package gocbcore

// ConfigSnapshot is a snapshot of the underlying configuration currently in use.
type ConfigSnapshot struct {
	state *kvMuxState
}

// RevID returns the config revision for this snapshot.
func (pi ConfigSnapshot) RevID() int64 {
	return pi.state.RevID()
}

// KeyToVbucket translates a particular key to its assigned vbucket.
func (pi ConfigSnapshot) KeyToVbucket(key []byte) (uint16, error) {
	if pi.state.VBMap() == nil {
		return 0, errUnsupportedOperation
	}
	return pi.state.VBMap().VbucketByKey(key), nil
}

// KeyToServer translates a particular key to its assigned server index.
func (pi ConfigSnapshot) KeyToServer(key []byte, replicaIdx uint32) (int, error) {
	if pi.state.VBMap() != nil {
		serverIdx, err := pi.state.VBMap().NodeByKey(key, replicaIdx)
		if err != nil {
			return 0, err
		}

		return serverIdx, nil
	}

	if pi.state.KetamaMap() != nil {
		serverIdx, err := pi.state.KetamaMap().NodeByKey(key)
		if err != nil {
			return 0, err
		}

		return serverIdx, nil
	}

	return 0, errCliInternalError
}

// VbucketToServer returns the server index for a particular vbucket.
func (pi ConfigSnapshot) VbucketToServer(vbID uint16, replicaIdx uint32) (int, error) {
	if pi.state.VBMap() == nil {
		return 0, errUnsupportedOperation
	}

	serverIdx, err := pi.state.VBMap().NodeByVbucket(vbID, replicaIdx)
	if err != nil {
		return 0, err
	}

	return serverIdx, nil
}

// VbucketsOnServer returns the list of VBuckets for a server.
func (pi ConfigSnapshot) VbucketsOnServer(index int) ([]uint16, error) {
	if pi.state.VBMap() == nil {
		return nil, errUnsupportedOperation
	}

	return pi.state.VBMap().VbucketsOnServer(index)
}

// NumVbuckets returns the number of VBuckets configured on the
// connected cluster.
func (pi ConfigSnapshot) NumVbuckets() (int, error) {
	if pi.state.VBMap() == nil {
		return 0, errUnsupportedOperation
	}

	return pi.state.VBMap().NumVbuckets(), nil
}

// NumReplicas returns the number of replicas configured on the
// connected cluster.
func (pi ConfigSnapshot) NumReplicas() (int, error) {
	if pi.state.VBMap() == nil {
		return 0, errUnsupportedOperation
	}

	return pi.state.VBMap().NumReplicas(), nil
}

// NumServers returns the number of servers accessible for K/V.
func (pi ConfigSnapshot) NumServers() (int, error) {
	return pi.state.NumPipelines(), nil
}

// BucketUUID returns the UUID of the bucket we are connected to.
func (pi ConfigSnapshot) BucketUUID() string {
	return pi.state.UUID()
}
