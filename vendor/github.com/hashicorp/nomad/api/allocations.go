package api

import (
	"context"
	"fmt"
	"io"
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

// Exec is used to execute a command inside a running task.  The command is to run inside
// the task environment.
//
// The parameters are:
// * ctx: context to set deadlines or timeout
// * allocation: the allocation to execute command inside
// * task: the task's name to execute command in
// * tty: indicates whether to start a pseudo-tty for the command
// * stdin, stdout, stderr: the std io to pass to command.
//      If tty is true, then streams need to point to a tty that's alive for the whole process
// * terminalSizeCh: A channel to send new tty terminal sizes
//
// The call blocks until command terminates (or an error occurs), and returns the exit code.
func (a *Allocations) Exec(ctx context.Context,
	alloc *Allocation, task string, tty bool, command []string,
	stdin io.Reader, stdout, stderr io.Writer,
	terminalSizeCh <-chan TerminalSize, q *QueryOptions) (exitCode int, err error) {

	s := &execSession{
		client:  a.client,
		alloc:   alloc,
		task:    task,
		tty:     tty,
		command: command,

		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,

		terminalSizeCh: terminalSizeCh,
		q:              q,
	}

	return s.run(ctx)
}

func (a *Allocations) Stats(alloc *Allocation, q *QueryOptions) (*AllocResourceUsage, error) {
	var resp AllocResourceUsage
	path := fmt.Sprintf("/v1/client/allocation/%s/stats", alloc.ID)
	_, err := a.client.query(path, &resp, q)
	return &resp, err
}

func (a *Allocations) GC(alloc *Allocation, q *QueryOptions) error {
	var resp struct{}
	_, err := a.client.query("/v1/client/allocation/"+alloc.ID+"/gc", &resp, nil)
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

func (a *Allocations) Stop(alloc *Allocation, q *QueryOptions) (*AllocStopResponse, error) {
	var resp AllocStopResponse
	_, err := a.client.putQuery("/v1/allocation/"+alloc.ID+"/stop", nil, &resp, q)
	return &resp, err
}

// AllocStopResponse is the response to an `AllocStopRequest`
type AllocStopResponse struct {
	// EvalID is the id of the follow up evalution for the rescheduled alloc.
	EvalID string

	WriteMeta
}

func (a *Allocations) Signal(alloc *Allocation, q *QueryOptions, task, signal string) error {
	req := AllocSignalRequest{
		Signal: signal,
		Task:   task,
	}

	var resp GenericResponse
	_, err := a.client.putQuery("/v1/client/allocation/"+alloc.ID+"/signal", &req, &resp, q)
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
	ResourcesExhausted map[string]*Resources
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

// Stub returns a list stub for the allocation
func (a *Allocation) Stub() *AllocationListStub {
	return &AllocationListStub{
		ID:                    a.ID,
		EvalID:                a.EvalID,
		Name:                  a.Name,
		Namespace:             a.Namespace,
		NodeID:                a.NodeID,
		NodeName:              a.NodeName,
		JobID:                 a.JobID,
		JobType:               *a.Job.Type,
		JobVersion:            *a.Job.Version,
		TaskGroup:             a.TaskGroup,
		DesiredStatus:         a.DesiredStatus,
		DesiredDescription:    a.DesiredDescription,
		ClientStatus:          a.ClientStatus,
		ClientDescription:     a.ClientDescription,
		TaskStates:            a.TaskStates,
		DeploymentStatus:      a.DeploymentStatus,
		FollowupEvalID:        a.FollowupEvalID,
		RescheduleTracker:     a.RescheduleTracker,
		PreemptedAllocations:  a.PreemptedAllocations,
		PreemptedByAllocation: a.PreemptedByAllocation,
		CreateIndex:           a.CreateIndex,
		ModifyIndex:           a.ModifyIndex,
		CreateTime:            a.CreateTime,
		ModifyTime:            a.ModifyTime,
	}
}

// ServerTerminalStatus returns true if the desired state of the allocation is
// terminal.
func (a *Allocation) ServerTerminalStatus() bool {
	switch a.DesiredStatus {
	case AllocDesiredStatusStop, AllocDesiredStatusEvict:
		return true
	default:
		return false
	}
}

// ClientTerminalStatus returns true if the client status is terminal and will
// therefore no longer transition.
func (a *Allocation) ClientTerminalStatus() bool {
	switch a.ClientStatus {
	case AllocClientStatusComplete, AllocClientStatusFailed, AllocClientStatusLost:
		return true
	default:
		return false
	}
}

// AllocationListStub is used to return a subset of an allocation
// during list operations.
type AllocationListStub struct {
	ID                    string
	EvalID                string
	Name                  string
	Namespace             string
	NodeID                string
	NodeName              string
	JobID                 string
	JobType               string
	JobVersion            uint64
	TaskGroup             string
	AllocatedResources    *AllocatedResources `json:",omitempty"`
	DesiredStatus         string
	DesiredDescription    string
	ClientStatus          string
	ClientDescription     string
	TaskStates            map[string]*TaskState
	DeploymentStatus      *AllocDeploymentStatus
	FollowupEvalID        string
	RescheduleTracker     *RescheduleTracker
	PreemptedAllocations  []string
	PreemptedByAllocation string
	CreateIndex           uint64
	ModifyIndex           uint64
	CreateTime            int64
	ModifyTime            int64
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
	Devices  []*AllocatedDeviceResource
}

type AllocatedSharedResources struct {
	DiskMB   int64
	Networks []*NetworkResource
	Ports    []PortMapping
}

type PortMapping struct {
	Label  string
	Value  int
	To     int
	HostIP string
}

type AllocatedCpuResources struct {
	CpuShares int64
}

type AllocatedMemoryResources struct {
	MemoryMB    int64
	MemoryMaxMB int64
}

type AllocatedDeviceResource struct {
	Vendor    string
	Type      string
	Name      string
	DeviceIDs []string
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

func (a Allocation) GetTaskGroup() *TaskGroup {
	for _, tg := range a.Job.TaskGroups {
		if *tg.Name == a.TaskGroup {
			return tg
		}
	}
	return nil
}

// RescheduleInfo is used to calculate remaining reschedule attempts
// according to the given time and the task groups reschedule policy
func (a Allocation) RescheduleInfo(t time.Time) (int, int) {
	tg := a.GetTaskGroup()
	if tg == nil || tg.ReschedulePolicy == nil {
		return 0, 0
	}
	reschedulePolicy := tg.ReschedulePolicy
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

type AllocSignalRequest struct {
	Task   string
	Signal string
}

// GenericResponse is used to respond to a request where no
// specific response information is needed.
type GenericResponse struct {
	WriteMeta
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

// ExecStreamingIOOperation represents a stream write operation: either appending data or close (exclusively)
type ExecStreamingIOOperation struct {
	Data  []byte `json:"data,omitempty"`
	Close bool   `json:"close,omitempty"`
}

// TerminalSize represents the size of the terminal
type TerminalSize struct {
	Height int `json:"height,omitempty"`
	Width  int `json:"width,omitempty"`
}

var execStreamingInputHeartbeat = ExecStreamingInput{}

// ExecStreamingInput represents user input to be sent to nomad exec handler.
//
// At most one field should be set.
type ExecStreamingInput struct {
	Stdin   *ExecStreamingIOOperation `json:"stdin,omitempty"`
	TTYSize *TerminalSize             `json:"tty_size,omitempty"`
}

// ExecStreamingExitResults captures the exit code of just completed nomad exec command
type ExecStreamingExitResult struct {
	ExitCode int `json:"exit_code"`
}

// ExecStreamingInput represents an output streaming entity, e.g. stdout/stderr update or termination
//
// At most one of these fields should be set: `Stdout`, `Stderr`, or `Result`.
// If `Exited` is true, then `Result` is non-nil, and other fields are nil.
type ExecStreamingOutput struct {
	Stdout *ExecStreamingIOOperation `json:"stdout,omitempty"`
	Stderr *ExecStreamingIOOperation `json:"stderr,omitempty"`

	Exited bool                     `json:"exited,omitempty"`
	Result *ExecStreamingExitResult `json:"result,omitempty"`
}
