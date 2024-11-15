package gocb

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type internalProviderCore struct {
	provider mgmtProvider

	tracer RequestTracer
	meter  *meterWrapper
}

func (ic *internalProviderCore) GetNodesMetadata(opts *GetNodesMetadataOptions) ([]NodeMetadata, error) {
	path := "/pools/default"
	start := time.Now()
	defer ic.meter.ValueRecord(meterValueServiceManagement, "internal_get_nodes_metadata", start)

	span := createSpan(ic.tracer, opts.ParentSpan, "internal_get_nodes_metadata", "management")
	span.SetAttribute("db.operation", "GET "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Path:          path,
		Method:        "GET",
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := ic.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp, "")
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		return nil, makeMgmtBadStatusError("failed to get nodes metadata", &req, resp)
	}

	var nodesData jsonClusterCfg
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&nodesData)
	if err != nil {
		return nil, err
	}

	nodes := make([]NodeMetadata, len(nodesData.Nodes))
	for i, nodeData := range nodesData.Nodes {
		nodes[i] = NodeMetadata(nodeData)
	}

	return nodes, nil
}
