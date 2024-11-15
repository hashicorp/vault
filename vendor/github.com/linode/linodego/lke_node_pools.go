package linodego

import (
	"context"
)

// LKELinodeStatus constants start with LKELinode and include
// Linode API LKENodePool Linode Status values
type LKELinodeStatus string

// LKENodePoolStatus constants reflect the current status of an LKENodePool
const (
	LKELinodeReady    LKELinodeStatus = "ready"
	LKELinodeNotReady LKELinodeStatus = "not_ready"
)

// LKENodePoolDisk represents a Node disk in an LKENodePool object
type LKENodePoolDisk struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

type LKENodePoolAutoscaler struct {
	Enabled bool `json:"enabled"`
	Min     int  `json:"min"`
	Max     int  `json:"max"`
}

// LKENodePoolLinode represents a LKENodePoolLinode object
type LKENodePoolLinode struct {
	ID         string          `json:"id"`
	InstanceID int             `json:"instance_id"`
	Status     LKELinodeStatus `json:"status"`
}

// LKENodePoolTaintEffect represents the effect value of a taint
type LKENodePoolTaintEffect string

const (
	LKENodePoolTaintEffectNoSchedule       LKENodePoolTaintEffect = "NoSchedule"
	LKENodePoolTaintEffectPreferNoSchedule LKENodePoolTaintEffect = "PreferNoSchedule"
	LKENodePoolTaintEffectNoExecute        LKENodePoolTaintEffect = "NoExecute"
)

// LKENodePoolTaint represents a corev1.Taint to add to an LKENodePool
type LKENodePoolTaint struct {
	Key    string                 `json:"key"`
	Value  string                 `json:"value,omitempty"`
	Effect LKENodePoolTaintEffect `json:"effect"`
}

// LKENodePoolLabels represents Kubernetes labels to add to an LKENodePool
type LKENodePoolLabels map[string]string

// LKENodePool represents a LKENodePool object
type LKENodePool struct {
	ID      int                 `json:"id"`
	Count   int                 `json:"count"`
	Type    string              `json:"type"`
	Disks   []LKENodePoolDisk   `json:"disks"`
	Linodes []LKENodePoolLinode `json:"nodes"`
	Tags    []string            `json:"tags"`
	Labels  LKENodePoolLabels   `json:"labels"`
	Taints  []LKENodePoolTaint  `json:"taints"`

	Autoscaler LKENodePoolAutoscaler `json:"autoscaler"`

	// NOTE: Disk encryption may not currently be available to all users.
	DiskEncryption InstanceDiskEncryption `json:"disk_encryption,omitempty"`
}

// LKENodePoolCreateOptions fields are those accepted by CreateLKENodePool
type LKENodePoolCreateOptions struct {
	Count  int                `json:"count"`
	Type   string             `json:"type"`
	Disks  []LKENodePoolDisk  `json:"disks"`
	Tags   []string           `json:"tags"`
	Labels LKENodePoolLabels  `json:"labels"`
	Taints []LKENodePoolTaint `json:"taints"`

	Autoscaler *LKENodePoolAutoscaler `json:"autoscaler,omitempty"`
}

// LKENodePoolUpdateOptions fields are those accepted by UpdateLKENodePoolUpdate
type LKENodePoolUpdateOptions struct {
	Count  int                 `json:"count,omitempty"`
	Tags   *[]string           `json:"tags,omitempty"`
	Labels *LKENodePoolLabels  `json:"labels,omitempty"`
	Taints *[]LKENodePoolTaint `json:"taints,omitempty"`

	Autoscaler *LKENodePoolAutoscaler `json:"autoscaler,omitempty"`
}

// GetCreateOptions converts a LKENodePool to LKENodePoolCreateOptions for
// use in CreateLKENodePool
func (l LKENodePool) GetCreateOptions() (o LKENodePoolCreateOptions) {
	o.Count = l.Count
	o.Disks = l.Disks
	o.Tags = l.Tags
	o.Labels = l.Labels
	o.Taints = l.Taints
	o.Autoscaler = &l.Autoscaler
	return
}

// GetUpdateOptions converts a LKENodePool to LKENodePoolUpdateOptions for use in UpdateLKENodePoolUpdate
func (l LKENodePool) GetUpdateOptions() (o LKENodePoolUpdateOptions) {
	o.Count = l.Count
	o.Tags = &l.Tags
	o.Labels = &l.Labels
	o.Taints = &l.Taints
	o.Autoscaler = &l.Autoscaler
	return
}

// ListLKENodePools lists LKENodePools
func (c *Client) ListLKENodePools(ctx context.Context, clusterID int, opts *ListOptions) ([]LKENodePool, error) {
	response, err := getPaginatedResults[LKENodePool](ctx, c, formatAPIPath("lke/clusters/%d/pools", clusterID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetLKENodePool gets the LKENodePool with the provided ID
func (c *Client) GetLKENodePool(ctx context.Context, clusterID, poolID int) (*LKENodePool, error) {
	e := formatAPIPath("lke/clusters/%d/pools/%d", clusterID, poolID)
	response, err := doGETRequest[LKENodePool](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateLKENodePool creates a LKENodePool
func (c *Client) CreateLKENodePool(ctx context.Context, clusterID int, opts LKENodePoolCreateOptions) (*LKENodePool, error) {
	e := formatAPIPath("lke/clusters/%d/pools", clusterID)
	response, err := doPOSTRequest[LKENodePool](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateLKENodePool updates the LKENodePool with the specified id
func (c *Client) UpdateLKENodePool(ctx context.Context, clusterID, poolID int, opts LKENodePoolUpdateOptions) (*LKENodePool, error) {
	e := formatAPIPath("lke/clusters/%d/pools/%d", clusterID, poolID)
	response, err := doPUTRequest[LKENodePool](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteLKENodePool deletes the LKENodePool with the specified id
func (c *Client) DeleteLKENodePool(ctx context.Context, clusterID, poolID int) error {
	e := formatAPIPath("lke/clusters/%d/pools/%d", clusterID, poolID)
	err := doDELETERequest(ctx, c, e)
	return err
}

// DeleteLKENodePoolNode deletes a given node from a node pool
func (c *Client) DeleteLKENodePoolNode(ctx context.Context, clusterID int, nodeID string) error {
	e := formatAPIPath("lke/clusters/%d/nodes/%s", clusterID, nodeID)
	err := doDELETERequest(ctx, c, e)
	return err
}
