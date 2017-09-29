package api

import (
	"sort"
)

// Deployments is used to query the deployments endpoints.
type Deployments struct {
	client *Client
}

// Deployments returns a new handle on the deployments.
func (c *Client) Deployments() *Deployments {
	return &Deployments{client: c}
}

// List is used to dump all of the deployments.
func (d *Deployments) List(q *QueryOptions) ([]*Deployment, *QueryMeta, error) {
	var resp []*Deployment
	qm, err := d.client.query("/v1/deployments", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(DeploymentIndexSort(resp))
	return resp, qm, nil
}

func (d *Deployments) PrefixList(prefix string) ([]*Deployment, *QueryMeta, error) {
	return d.List(&QueryOptions{Prefix: prefix})
}

// Info is used to query a single deployment by its ID.
func (d *Deployments) Info(deploymentID string, q *QueryOptions) (*Deployment, *QueryMeta, error) {
	var resp Deployment
	qm, err := d.client.query("/v1/deployment/"+deploymentID, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Allocations is used to retrieve a set of allocations that are part of the
// deployment
func (d *Deployments) Allocations(deploymentID string, q *QueryOptions) ([]*AllocationListStub, *QueryMeta, error) {
	var resp []*AllocationListStub
	qm, err := d.client.query("/v1/deployment/allocations/"+deploymentID, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(AllocIndexSort(resp))
	return resp, qm, nil
}

// Fail is used to fail the given deployment.
func (d *Deployments) Fail(deploymentID string, q *WriteOptions) (*DeploymentUpdateResponse, *WriteMeta, error) {
	var resp DeploymentUpdateResponse
	req := &DeploymentFailRequest{
		DeploymentID: deploymentID,
	}
	wm, err := d.client.write("/v1/deployment/fail/"+deploymentID, req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Pause is used to pause or unpause the given deployment.
func (d *Deployments) Pause(deploymentID string, pause bool, q *WriteOptions) (*DeploymentUpdateResponse, *WriteMeta, error) {
	var resp DeploymentUpdateResponse
	req := &DeploymentPauseRequest{
		DeploymentID: deploymentID,
		Pause:        pause,
	}
	wm, err := d.client.write("/v1/deployment/pause/"+deploymentID, req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// PromoteAll is used to promote all canaries in the given deployment
func (d *Deployments) PromoteAll(deploymentID string, q *WriteOptions) (*DeploymentUpdateResponse, *WriteMeta, error) {
	var resp DeploymentUpdateResponse
	req := &DeploymentPromoteRequest{
		DeploymentID: deploymentID,
		All:          true,
	}
	wm, err := d.client.write("/v1/deployment/promote/"+deploymentID, req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// PromoteGroups is used to promote canaries in the passed groups in the given deployment
func (d *Deployments) PromoteGroups(deploymentID string, groups []string, q *WriteOptions) (*DeploymentUpdateResponse, *WriteMeta, error) {
	var resp DeploymentUpdateResponse
	req := &DeploymentPromoteRequest{
		DeploymentID: deploymentID,
		Groups:       groups,
	}
	wm, err := d.client.write("/v1/deployment/promote/"+deploymentID, req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// SetAllocHealth is used to set allocation health for allocs that are part of
// the given deployment
func (d *Deployments) SetAllocHealth(deploymentID string, healthy, unhealthy []string, q *WriteOptions) (*DeploymentUpdateResponse, *WriteMeta, error) {
	var resp DeploymentUpdateResponse
	req := &DeploymentAllocHealthRequest{
		DeploymentID:           deploymentID,
		HealthyAllocationIDs:   healthy,
		UnhealthyAllocationIDs: unhealthy,
	}
	wm, err := d.client.write("/v1/deployment/allocation-health/"+deploymentID, req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Deployment is used to serialize an deployment.
type Deployment struct {
	ID                string
	Namespace         string
	JobID             string
	JobVersion        uint64
	JobModifyIndex    uint64
	JobCreateIndex    uint64
	TaskGroups        map[string]*DeploymentState
	Status            string
	StatusDescription string
	CreateIndex       uint64
	ModifyIndex       uint64
}

// DeploymentState tracks the state of a deployment for a given task group.
type DeploymentState struct {
	PlacedCanaries  []string
	AutoRevert      bool
	Promoted        bool
	DesiredCanaries int
	DesiredTotal    int
	PlacedAllocs    int
	HealthyAllocs   int
	UnhealthyAllocs int
}

// DeploymentIndexSort is a wrapper to sort deployments by CreateIndex. We
// reverse the test so that we get the highest index first.
type DeploymentIndexSort []*Deployment

func (d DeploymentIndexSort) Len() int {
	return len(d)
}

func (d DeploymentIndexSort) Less(i, j int) bool {
	return d[i].CreateIndex > d[j].CreateIndex
}

func (d DeploymentIndexSort) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// DeploymentUpdateResponse is used to respond to a deployment change. The
// response will include the modify index of the deployment as well as details
// of any triggered evaluation.
type DeploymentUpdateResponse struct {
	EvalID                string
	EvalCreateIndex       uint64
	DeploymentModifyIndex uint64
	RevertedJobVersion    *uint64
	WriteMeta
}

// DeploymentAllocHealthRequest is used to set the health of a set of
// allocations as part of a deployment.
type DeploymentAllocHealthRequest struct {
	DeploymentID string

	// Marks these allocations as healthy, allow further allocations
	// to be rolled.
	HealthyAllocationIDs []string

	// Any unhealthy allocations fail the deployment
	UnhealthyAllocationIDs []string

	WriteRequest
}

// DeploymentPromoteRequest is used to promote task groups in a deployment
type DeploymentPromoteRequest struct {
	DeploymentID string

	// All is to promote all task groups
	All bool

	// Groups is used to set the promotion status per task group
	Groups []string

	WriteRequest
}

// DeploymentPauseRequest is used to pause a deployment
type DeploymentPauseRequest struct {
	DeploymentID string

	// Pause sets the pause status
	Pause bool

	WriteRequest
}

// DeploymentSpecificRequest is used to make a request specific to a particular
// deployment
type DeploymentSpecificRequest struct {
	DeploymentID string
	QueryOptions
}

// DeploymentFailRequest is used to fail a particular deployment
type DeploymentFailRequest struct {
	DeploymentID string
	WriteRequest
}

// SingleDeploymentResponse is used to respond with a single deployment
type SingleDeploymentResponse struct {
	Deployment *Deployment
	QueryMeta
}
