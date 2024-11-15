package linodego

import (
	"context"
)

// NodeBalancerNode objects represent a backend that can accept traffic for a NodeBalancer Config
type NodeBalancerNode struct {
	ID             int      `json:"id"`
	Address        string   `json:"address"`
	Label          string   `json:"label"`
	Status         string   `json:"status"`
	Weight         int      `json:"weight"`
	Mode           NodeMode `json:"mode"`
	ConfigID       int      `json:"config_id"`
	NodeBalancerID int      `json:"nodebalancer_id"`
}

// NodeMode is the mode a NodeBalancer should use when sending traffic to a NodeBalancer Node
type NodeMode string

var (
	// ModeAccept is the NodeMode indicating a NodeBalancer Node is accepting traffic
	ModeAccept NodeMode = "accept"

	// ModeReject is the NodeMode indicating a NodeBalancer Node is not receiving traffic
	ModeReject NodeMode = "reject"

	// ModeDrain is the NodeMode indicating a NodeBalancer Node is not receiving new traffic, but may continue receiving traffic from pinned connections
	ModeDrain NodeMode = "drain"

	// ModeBackup is the NodeMode indicating a NodeBalancer Node will only receive traffic if all "accept" Nodes are down
	ModeBackup NodeMode = "backup"
)

// NodeBalancerNodeCreateOptions fields are those accepted by CreateNodeBalancerNode
type NodeBalancerNodeCreateOptions struct {
	Address string   `json:"address"`
	Label   string   `json:"label"`
	Weight  int      `json:"weight,omitempty"`
	Mode    NodeMode `json:"mode,omitempty"`
}

// NodeBalancerNodeUpdateOptions fields are those accepted by UpdateNodeBalancerNode
type NodeBalancerNodeUpdateOptions struct {
	Address string   `json:"address,omitempty"`
	Label   string   `json:"label,omitempty"`
	Weight  int      `json:"weight,omitempty"`
	Mode    NodeMode `json:"mode,omitempty"`
}

// GetCreateOptions converts a NodeBalancerNode to NodeBalancerNodeCreateOptions for use in CreateNodeBalancerNode
func (i NodeBalancerNode) GetCreateOptions() NodeBalancerNodeCreateOptions {
	return NodeBalancerNodeCreateOptions{
		Address: i.Address,
		Label:   i.Label,
		Weight:  i.Weight,
		Mode:    i.Mode,
	}
}

// GetUpdateOptions converts a NodeBalancerNode to NodeBalancerNodeUpdateOptions for use in UpdateNodeBalancerNode
func (i NodeBalancerNode) GetUpdateOptions() NodeBalancerNodeUpdateOptions {
	return NodeBalancerNodeUpdateOptions{
		Address: i.Address,
		Label:   i.Label,
		Weight:  i.Weight,
		Mode:    i.Mode,
	}
}

// ListNodeBalancerNodes lists NodeBalancerNodes
func (c *Client) ListNodeBalancerNodes(ctx context.Context, nodebalancerID int, configID int, opts *ListOptions) ([]NodeBalancerNode, error) {
	response, err := getPaginatedResults[NodeBalancerNode](ctx, c, formatAPIPath("nodebalancers/%d/configs/%d/nodes", nodebalancerID, configID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetNodeBalancerNode gets the template with the provided ID
func (c *Client) GetNodeBalancerNode(ctx context.Context, nodebalancerID int, configID int, nodeID int) (*NodeBalancerNode, error) {
	e := formatAPIPath("nodebalancers/%d/configs/%d/nodes/%d", nodebalancerID, configID, nodeID)
	response, err := doGETRequest[NodeBalancerNode](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateNodeBalancerNode creates a NodeBalancerNode
func (c *Client) CreateNodeBalancerNode(ctx context.Context, nodebalancerID int, configID int, opts NodeBalancerNodeCreateOptions) (*NodeBalancerNode, error) {
	e := formatAPIPath("nodebalancers/%d/configs/%d/nodes", nodebalancerID, configID)
	response, err := doPOSTRequest[NodeBalancerNode](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateNodeBalancerNode updates the NodeBalancerNode with the specified id
func (c *Client) UpdateNodeBalancerNode(ctx context.Context, nodebalancerID int, configID int, nodeID int, opts NodeBalancerNodeUpdateOptions) (*NodeBalancerNode, error) {
	e := formatAPIPath("nodebalancers/%d/configs/%d/nodes/%d", nodebalancerID, configID, nodeID)
	response, err := doPUTRequest[NodeBalancerNode](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteNodeBalancerNode deletes the NodeBalancerNode with the specified id
func (c *Client) DeleteNodeBalancerNode(ctx context.Context, nodebalancerID int, configID int, nodeID int) error {
	e := formatAPIPath("nodebalancers/%d/configs/%d/nodes/%d", nodebalancerID, configID, nodeID)
	err := doDELETERequest(ctx, c, e)
	return err
}
