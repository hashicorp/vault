package gocbcore

type pipelineSnapshot struct {
	state *kvMuxState

	idx int
}

func (pi pipelineSnapshot) RevID() int64 {
	return pi.state.revID
}

func (pi pipelineSnapshot) NumPipelines() int {
	return pi.state.NumPipelines()
}

func (pi pipelineSnapshot) PipelineAt(idx int) *memdPipeline {
	return pi.state.GetPipeline(idx)
}

func (pi pipelineSnapshot) Iterate(offset int, cb func(*memdPipeline) bool) {
	l := pi.state.NumPipelines()
	pi.idx = offset
	for iters := 0; iters < l; iters++ {
		pi.idx = (pi.idx + 1) % l
		p := pi.state.GetPipeline(pi.idx)

		if cb(p) {
			return
		}
	}
}

func (pi pipelineSnapshot) NodeByVbucket(vbID uint16, replicaID uint32) (int, error) {
	if pi.state.vbMap == nil {
		return 0, errUnsupportedOperation
	}

	return pi.state.vbMap.NodeByVbucket(vbID, replicaID)
}
