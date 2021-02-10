package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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

	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	errCh := make(chan error, 4)

	sender, output := a.execFrames(ctx, alloc, task, tty, command, errCh, q)

	select {
	case err := <-errCh:
		return -2, err
	default:
	}

	// Errors resulting from sending input (in goroutines) are silently dropped.
	// To mitigate this, extra care is needed to distinguish between actual send errors
	// and from send errors due to command terminating and our race to detect failures.
	// If we have an actual network failure or send a bad input, we'd get an
	// error in the reading side of websocket.

	go func() {

		bytes := make([]byte, 2048)
		for {
			if ctx.Err() != nil {
				return
			}

			input := ExecStreamingInput{Stdin: &ExecStreamingIOOperation{}}

			n, err := stdin.Read(bytes)

			// always send data if we read some
			if n != 0 {
				input.Stdin.Data = bytes[:n]
				sender(&input)
			}

			// then handle error
			if err == io.EOF {
				// if n != 0, send data and we'll get n = 0 on next read
				if n == 0 {
					input.Stdin.Close = true
					sender(&input)
					return
				}
			} else if err != nil {
				errCh <- err
				return
			}
		}
	}()

	// forwarding terminal size
	go func() {
		for {
			resizeInput := ExecStreamingInput{}

			select {
			case <-ctx.Done():
				return
			case size, ok := <-terminalSizeCh:
				if !ok {
					return
				}
				resizeInput.TTYSize = &size
				sender(&resizeInput)
			}

		}
	}()

	// send a heartbeat every 10 seconds
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			// heartbeat message
			case <-time.After(10 * time.Second):
				sender(&execStreamingInputHeartbeat)
			}

		}
	}()

	for {
		select {
		case err := <-errCh:
			// drop websocket code, not relevant to user
			if wsErr, ok := err.(*websocket.CloseError); ok && wsErr.Text != "" {
				return -2, errors.New(wsErr.Text)
			}
			return -2, err
		case <-ctx.Done():
			return -2, ctx.Err()
		case frame, ok := <-output:
			if !ok {
				return -2, errors.New("disconnected without receiving the exit code")
			}

			switch {
			case frame.Stdout != nil:
				if len(frame.Stdout.Data) != 0 {
					stdout.Write(frame.Stdout.Data)
				}
				// don't really do anything if stdout is closing
			case frame.Stderr != nil:
				if len(frame.Stderr.Data) != 0 {
					stderr.Write(frame.Stderr.Data)
				}
				// don't really do anything if stderr is closing
			case frame.Exited && frame.Result != nil:
				return frame.Result.ExitCode, nil
			default:
				// noop - heartbeat
			}
		}
	}
}

func (a *Allocations) execFrames(ctx context.Context, alloc *Allocation, task string, tty bool, command []string,
	errCh chan<- error, q *QueryOptions) (sendFn func(*ExecStreamingInput) error, output <-chan *ExecStreamingOutput) {
	nodeClient, _ := a.client.GetNodeClientWithTimeout(alloc.NodeID, ClientConnTimeout, q)

	if q == nil {
		q = &QueryOptions{}
	}
	if q.Params == nil {
		q.Params = make(map[string]string)
	}

	commandBytes, err := json.Marshal(command)
	if err != nil {
		errCh <- fmt.Errorf("failed to marshal command: %s", err)
		return nil, nil
	}

	q.Params["tty"] = strconv.FormatBool(tty)
	q.Params["task"] = task
	q.Params["command"] = string(commandBytes)

	reqPath := fmt.Sprintf("/v1/client/allocation/%s/exec", alloc.ID)

	var conn *websocket.Conn

	if nodeClient != nil {
		conn, _, _ = nodeClient.websocket(reqPath, q)
	}

	if conn == nil {
		conn, _, err = a.client.websocket(reqPath, q)
		if err != nil {
			errCh <- err
			return nil, nil
		}
	}

	// Create the output channel
	frames := make(chan *ExecStreamingOutput, 10)

	go func() {
		defer conn.Close()
		for ctx.Err() == nil {

			// Decode the next frame
			var frame ExecStreamingOutput
			err := conn.ReadJSON(&frame)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				close(frames)
				return
			} else if err != nil {
				errCh <- err
				return
			}

			frames <- &frame
		}
	}()

	var sendLock sync.Mutex
	send := func(v *ExecStreamingInput) error {
		sendLock.Lock()
		defer sendLock.Unlock()

		return conn.WriteJSON(v)
	}

	return send, frames

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
