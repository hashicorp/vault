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

const (
	AllocDesiredStatusRun   = "run"   // Allocation should run
	AllocDesiredStatusStop  = "stop"  // Allocation should stop
	AllocDesiredStatusEvict = "evict" // Allocation should stop, and was evicted
)

const (
	AllocClientStatusPending  = "pending"
	AllocClientStatusRunning  = "running"
	AllocClientStatusComplete = "complete"
	AllocClientStatusFailed   = "failed"
	AllocClientStatusLost     = "lost"
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
	var resp AllocResourceUsage
	path := fmt.Sprintf("/v1/client/allocation/%s/stats", alloc.ID)
	_, err := a.client.query(path, &resp, q)
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

func (a *Allocations) Restart(alloc *Allocation, taskName string, q *QueryOptions) error {
	req := AllocationRestartRequest{
		TaskName: taskName,
	}

	var resp struct{}
	_, err := a.client.putQuery("/v1/client/allocation/"+alloc.ID+"/restart", &req, &resp, q)
	return err
}

// Allocation is used for serialization of allocations.
type Allocation struct {
	ID                    string
	Namespace             string
	EvalID                string
	Name                  string
	NodeID                string
	NodeName              string
	JobID                 string
	Job                   *Job
	TaskGroup             string
	Resources             *Resources
	TaskResources         map[string]*Resources
	AllocatedResources    *AllocatedResources
	Services              map[string]string
	Metrics               *AllocationMetric
	DesiredStatus         string
	DesiredDescription    string
	DesiredTransition     DesiredTransition
	ClientStatus          string
	ClientDescription     string
	TaskStates            map[string]*TaskState
	DeploymentID          string
	DeploymentStatus      *AllocDeploymentStatus
	FollowupEvalID        string
	PreviousAllocation    string
	NextAllocation        string
	RescheduleTracker     *RescheduleTracker
	PreemptedAllocations  []string
	PreemptedByAllocation string
	CreateIndex           uint64
	ModifyIndex           uint64
	AllocModifyIndex      uint64
	CreateTime            int64
	ModifyTime            int64
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
	QuotaExhausted     []string
	// Deprecated, replaced with ScoreMetaData
	Scores            map[string]float64
	AllocationTime    time.Duration
	CoalescedFailures int
	ScoreMetaData     []*NodeScoreMeta
}

// NodeScoreMeta is used to serialize node scoring metadata
// displayed in the CLI during verbose mode
type NodeScoreMeta struct {
	NodeID    string
	Scores    map[string]float64
	NormScore float64
}

// AllocationListStub is used to return a subset of an allocation
// during list operations.
type AllocationListStub struct {
	ID                 string
	EvalID             string
	Name               string
	Namespace          string
	NodeID             string
	NodeName           string
	JobID              string
	JobType            string
	JobVersion         uint64
	TaskGroup          string
	DesiredStatus      string
	DesiredDescription string
	ClientStatus       string
	ClientDescription  string
	TaskStates         map[string]*TaskState
	DeploymentStatus   *AllocDeploymentStatus
	FollowupEvalID     string
	RescheduleTracker  *RescheduleTracker
	CreateIndex        uint64
	ModifyIndex        uint64
	CreateTime         int64
	ModifyTime         int64
}

// AllocDeploymentStatus captures the status of the allocation as part of the
// deployment. This can include things like if the allocation has been marked as
// healthy.
type AllocDeploymentStatus struct {
	Healthy     *bool
	Timestamp   time.Time
	Canary      bool
	ModifyIndex uint64
}

type AllocatedResources struct {
	Tasks  map[string]*AllocatedTaskResources
	Shared AllocatedSharedResources
}

type AllocatedTaskResources struct {
	Cpu      AllocatedCpuResources
	Memory   AllocatedMemoryResources
	Networks []*NetworkResource
}

type AllocatedSharedResources struct {
	DiskMB int64
}

type AllocatedCpuResources struct {
	CpuShares int64
}

type AllocatedMemoryResources struct {
	MemoryMB int64
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

// RescheduleInfo is used to calculate remaining reschedule attempts
// according to the given time and the task groups reschedule policy
func (a Allocation) RescheduleInfo(t time.Time) (int, int) {
	var reschedulePolicy *ReschedulePolicy
	for _, tg := range a.Job.TaskGroups {
		if *tg.Name == a.TaskGroup {
			reschedulePolicy = tg.ReschedulePolicy
		}
	}
	if reschedulePolicy == nil {
		return 0, 0
	}
	availableAttempts := *reschedulePolicy.Attempts
	interval := *reschedulePolicy.Interval
	attempted := 0

	// Loop over reschedule tracker to find attempts within the restart policy's interval
	if a.RescheduleTracker != nil && availableAttempts > 0 && interval > 0 {
		for j := len(a.RescheduleTracker.Events) - 1; j >= 0; j-- {
			lastAttempt := a.RescheduleTracker.Events[j].RescheduleTime
			timeDiff := t.UTC().UnixNano() - lastAttempt
			if timeDiff < interval.Nanoseconds() {
				attempted += 1
			}
		}
	}
	return attempted, availableAttempts
}

type AllocationRestartRequest struct {
	TaskName string
}

// RescheduleTracker encapsulates previous reschedule events
type RescheduleTracker struct {
	Events []*RescheduleEvent
}

// RescheduleEvent is used to keep track of previous attempts at rescheduling an allocation
type RescheduleEvent struct {
	// RescheduleTime is the timestamp of a reschedule attempt
	RescheduleTime int64

	// PrevAllocID is the ID of the previous allocation being restarted
	PrevAllocID string

	// PrevNodeID is the node ID of the previous allocation
	PrevNodeID string
}

// DesiredTransition is used to mark an allocation as having a desired state
// transition. This information can be used by the scheduler to make the
// correct decision.
type DesiredTransition struct {
	// Migrate is used to indicate that this allocation should be stopped and
	// migrated to another node.
	Migrate *bool

	// Reschedule is used to indicate that this allocation is eligible to be
	// rescheduled.
	Reschedule *bool
}

// ShouldMigrate returns whether the transition object dictates a migration.
func (d DesiredTransition) ShouldMigrate() bool {
	return d.Migrate != nil && *d.Migrate
}
