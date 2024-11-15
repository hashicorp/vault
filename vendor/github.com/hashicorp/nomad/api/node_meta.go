// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

// NodeMetaApplyRequest contains the Node meta update.
type NodeMetaApplyRequest struct {
	NodeID string
	Meta   map[string]*string
}

// NodeMetaResponse contains the merged Node metadata.
type NodeMetaResponse struct {
	// Meta is the merged static + dynamic Node metadata
	Meta map[string]string

	// Dynamic is the dynamic Node metadata (set via API)
	Dynamic map[string]*string

	// Static is the static Node metadata (set via agent configuration)
	Static map[string]string
}

// NodeMeta is a client for manipulating dynamic Node metadata.
type NodeMeta struct {
	client *Client
}

// Meta returns a NodeMeta client.
func (n *Nodes) Meta() *NodeMeta {
	return &NodeMeta{client: n.client}
}

// Apply dynamic Node metadata updates to a Node. If NodeID is unset then Node
// receiving the request is modified.
func (n *NodeMeta) Apply(meta *NodeMetaApplyRequest, qo *QueryOptions) (*NodeMetaResponse, error) {
	var out NodeMetaResponse
	_, err := n.client.postQuery("/v1/client/metadata", meta, &out, qo)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Read Node metadata (dynamic and static merged) from a Node directly. May
// differ from Node.Info as dynamic Node metadata updates are batched and may
// be delayed up to 10 seconds.
//
// If nodeID is empty then the metadata for the Node receiving the request is
// returned.
func (n *NodeMeta) Read(nodeID string, qo *QueryOptions) (*NodeMetaResponse, error) {
	if qo == nil {
		qo = &QueryOptions{}
	}

	if qo.Params == nil {
		qo.Params = make(map[string]string)
	}

	if nodeID != "" {
		qo.Params["node_id"] = nodeID
	}

	var out NodeMetaResponse
	_, err := n.client.query("/v1/client/metadata", &out, qo)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
