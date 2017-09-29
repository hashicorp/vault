package api

import (
	"fmt"
	"sort"
	"time"
)

var (
	// NodeDownErr marks an operation as not able to complete since the node is
	// down.
	NodeDownErr = fmt.Errorf("node down")
)

// Allocations is used to query the alloc-related endpoints.
type Allocations struct {
	client *Client
}

// Allocations returns a handle on the allocs endpoints.
func (c *Client) Allocations() *Allocations {
	return &Allocations{client: c}
}

// List returns a list of all of the allocations.
func (a *Allocations) List(q *QueryOptions) ([]*AllocationListStub, *QueryMeta, error) {
	var resp []*AllocationListStub
	qm, err := a.client.query("/v1/allocations", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(AllocIndexSort(resp))
	return resp, qm, nil
}

func (a *Allocations) PrefixList(prefix string) ([]*AllocationListStub, *QueryMeta, error) {
	return a.List(&QueryOptions{Prefix: prefix})
}

// Info is used to retrieve a single allocation.
func (a *Allocations) Info(allocID string, q *QueryOptions) (*Allocation, *QueryMeta, error) {
	var resp Allocation
	qm, err := a.client.query("/v1/allocation/"+allocID, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

func (a *Allocations) Stats(alloc *Allocation, q *QueryOptions) (*AllocResourceUsage, error) {
	nodeClient, err := a.client.GetNodeClient(alloc.NodeID, q)
	if err != nil {
		return nil, err
	}

	var resp AllocResourceUsage
	_, err = nodeClient.query("/v1/client/allocation/"+alloc.ID+"/stats", &resp, nil)
	return &resp, err
}

func (a *Allocations) GC(alloc *Allocation, q *QueryOptions) error {
	nodeClient, err := a.client.GetNodeClient(alloc.NodeID, q)
	if err != nil {
		return err
	}

	var resp struct{}
	_, err = nodeClient.query("/v1/client/allocation/"+alloc.ID+"/gc", &resp, nil)
	return err
}

// Allocation is used for serialization of allocations.
type Allocation struct {
	ID                 string
	Namespace          string
	EvalID             string
	Name               string
	NodeID             string
	JobID              string
	Job                *Job
	TaskGroup          string
	Resources          *Resources
	TaskResources      map[string]*Resources
	Services           map[string]string
	Metrics            *AllocationMetric
	DesiredStatus      string
	DesiredDescription string
	ClientStatus       string
	ClientDescription  string
	TaskStates         map[string]*TaskState
	DeploymentID       string
	DeploymentStatus   *AllocDeploymentStatus
	PreviousAllocation string
	CreateIndex        uint64
	ModifyIndex        uint64
	AllocModifyIndex   uint64
	CreateTime         int64
}

// AllocationMetric is used to deserialize allocation metrics.
type AllocationMetric struct {
	NodesEvaluated     int
	NodesFiltered      int
	NodesAvailable     map[string]int
	ClassFiltered      map[string]int
	ConstraintFiltered map[string]int
	NodesExhausted     int
	ClassExhausted     map[string]int
	DimensionExhausted map[string]int
	Scores             map[string]float64
	AllocationTime     time.Duration
	CoalescedFailures  int
}

// AllocationListStub is used to return a subset of an allocation
// during list operations.
type AllocationListStub struct {
	ID                 string
	EvalID             string
	Name               string
	NodeID             string
	JobID              string
	JobVersion         uint64
	TaskGroup          string
	DesiredStatus      string
	DesiredDescription string
	ClientStatus       string
	ClientDescription  string
	TaskStates         map[string]*TaskState
	DeploymentStatus   *AllocDeploymentStatus
	CreateIndex        uint64
	ModifyIndex        uint64
	CreateTime         int64
}

// AllocDeploymentStatus captures the status of the allocation as part of the
// deployment. This can include things like if the allocation has been marked as
// heatlhy.
type AllocDeploymentStatus struct {
	Healthy     *bool
	ModifyIndex uint64
}

// AllocIndexSort reverse sorts allocs by CreateIndex.
type AllocIndexSort []*AllocationListStub

func (a AllocIndexSort) Len() int {
	return len(a)
}

func (a AllocIndexSort) Less(i, j int) bool {
	return a[i].CreateIndex > a[j].CreateIndex
}

func (a AllocIndexSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
