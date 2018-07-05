package api

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/nomad/nomad/structs"
)

// Nodes is used to query node-related API endpoints
type Nodes struct {
	client *Client
}

// Nodes returns a handle on the node endpoints.
func (c *Client) Nodes() *Nodes {
	return &Nodes{client: c}
}

// List is used to list out all of the nodes
func (n *Nodes) List(q *QueryOptions) ([]*NodeListStub, *QueryMeta, error) {
	var resp NodeIndexSort
	qm, err := n.client.query("/v1/nodes", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(resp)
	return resp, qm, nil
}

func (n *Nodes) PrefixList(prefix string) ([]*NodeListStub, *QueryMeta, error) {
	return n.List(&QueryOptions{Prefix: prefix})
}

// Info is used to query a specific node by its ID.
func (n *Nodes) Info(nodeID string, q *QueryOptions) (*Node, *QueryMeta, error) {
	var resp Node
	qm, err := n.client.query("/v1/node/"+nodeID, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// NodeUpdateDrainRequest is used to update the drain specification for a node.
type NodeUpdateDrainRequest struct {
	// NodeID is the node to update the drain specification for.
	NodeID string

	// DrainSpec is the drain specification to set for the node. A nil DrainSpec
	// will disable draining.
	DrainSpec *DrainSpec

	// MarkEligible marks the node as eligible for scheduling if removing
	// the drain strategy.
	MarkEligible bool
}

// NodeDrainUpdateResponse is used to respond to a node drain update
type NodeDrainUpdateResponse struct {
	NodeModifyIndex uint64
	EvalIDs         []string
	EvalCreateIndex uint64
	WriteMeta
}

// UpdateDrain is used to update the drain strategy for a given node. If
// markEligible is true and the drain is being removed, the node will be marked
// as having its scheduling being eligible
func (n *Nodes) UpdateDrain(nodeID string, spec *DrainSpec, markEligible bool, q *WriteOptions) (*NodeDrainUpdateResponse, error) {
	req := &NodeUpdateDrainRequest{
		NodeID:       nodeID,
		DrainSpec:    spec,
		MarkEligible: markEligible,
	}

	var resp NodeDrainUpdateResponse
	wm, err := n.client.write("/v1/node/"+nodeID+"/drain", req, &resp, q)
	if err != nil {
		return nil, err
	}
	resp.WriteMeta = *wm
	return &resp, nil
}

// MonitorMsgLevels represents the severity log level of a MonitorMessage.
type MonitorMsgLevel int

const (
	MonitorMsgLevelNormal MonitorMsgLevel = 0
	MonitorMsgLevelInfo   MonitorMsgLevel = 1
	MonitorMsgLevelWarn   MonitorMsgLevel = 2
	MonitorMsgLevelError  MonitorMsgLevel = 3
)

// MonitorMessage contains a message and log level.
type MonitorMessage struct {
	Level   MonitorMsgLevel
	Message string
}

// Messagef formats a new MonitorMessage.
func Messagef(lvl MonitorMsgLevel, msg string, args ...interface{}) *MonitorMessage {
	return &MonitorMessage{
		Level:   lvl,
		Message: fmt.Sprintf(msg, args...),
	}
}

func (m *MonitorMessage) String() string {
	return m.Message
}

// MonitorDrain emits drain related events on the returned string channel. The
// channel will be closed when all allocations on the draining node have
// stopped or the context is canceled.
func (n *Nodes) MonitorDrain(ctx context.Context, nodeID string, index uint64, ignoreSys bool) <-chan *MonitorMessage {
	outCh := make(chan *MonitorMessage, 8)
	nodeCh := make(chan *MonitorMessage, 1)
	allocCh := make(chan *MonitorMessage, 8)

	// Multiplex node and alloc chans onto outCh. This goroutine closes
	// outCh when other chans have been closed or context canceled.
	multiplexCtx, cancel := context.WithCancel(ctx)
	go n.monitorDrainMultiplex(multiplexCtx, cancel, outCh, nodeCh, allocCh)

	// Monitor node for updates
	go n.monitorDrainNode(multiplexCtx, cancel, nodeID, index, nodeCh)

	// Monitor allocs on node for updates
	go n.monitorDrainAllocs(multiplexCtx, nodeID, ignoreSys, allocCh)

	return outCh
}

// monitorDrainMultiplex multiplexes node and alloc updates onto the out chan.
// Closes out chan when either the context is canceled, both update chans are
// closed, or an error occurs.
func (n *Nodes) monitorDrainMultiplex(ctx context.Context, cancel func(),
	outCh chan<- *MonitorMessage, nodeCh, allocCh <-chan *MonitorMessage) {

	defer cancel()
	defer close(outCh)

	nodeOk := true
	allocOk := true
	var msg *MonitorMessage
	for {
		// If both chans have been closed, close the output chan
		if !nodeOk && !allocOk {
			return
		}

		select {
		case msg, nodeOk = <-nodeCh:
			if !nodeOk {
				// nil chan to prevent further recvs
				nodeCh = nil
			}

		case msg, allocOk = <-allocCh:
			if !allocOk {
				// nil chan to prevent further recvs
				allocCh = nil
			}

		case <-ctx.Done():
			return
		}

		if msg == nil {
			continue
		}

		select {
		case outCh <- msg:
		case <-ctx.Done():

			// If we are exiting but we have a message, attempt to send it
			// so we don't lose a message but do not block.
			select {
			case outCh <- msg:
			default:
			}

			return
		}

		// Abort on error messages
		if msg.Level == MonitorMsgLevelError {
			return
		}
	}
}

// monitorDrainNode emits node updates on nodeCh and closes the channel when
// the node has finished draining.
func (n *Nodes) monitorDrainNode(ctx context.Context, cancel func(),
	nodeID string, index uint64, nodeCh chan<- *MonitorMessage) {
	defer close(nodeCh)

	var lastStrategy *DrainStrategy
	var strategyChanged bool
	q := QueryOptions{
		AllowStale: true,
		WaitIndex:  index,
	}
	for {
		node, meta, err := n.Info(nodeID, &q)
		if err != nil {
			msg := Messagef(MonitorMsgLevelError, "Error monitoring node: %v", err)
			select {
			case nodeCh <- msg:
			case <-ctx.Done():
			}
			return
		}

		if node.DrainStrategy == nil {
			var msg *MonitorMessage
			if strategyChanged {
				msg = Messagef(MonitorMsgLevelInfo, "Node %q has marked all allocations for migration", nodeID)
			} else {
				msg = Messagef(MonitorMsgLevelInfo, "No drain strategy set for node %s", nodeID)
				defer cancel()
			}
			select {
			case nodeCh <- msg:
			case <-ctx.Done():
			}
			return
		}

		if node.Status == structs.NodeStatusDown {
			msg := Messagef(MonitorMsgLevelWarn, "Node %q down", nodeID)
			select {
			case nodeCh <- msg:
			case <-ctx.Done():
			}
		}

		// DrainStrategy changed
		if lastStrategy != nil && !node.DrainStrategy.Equal(lastStrategy) {
			msg := Messagef(MonitorMsgLevelInfo, "Node %q drain updated: %s", nodeID, node.DrainStrategy)
			select {
			case nodeCh <- msg:
			case <-ctx.Done():
				return
			}
		}

		lastStrategy = node.DrainStrategy
		strategyChanged = true

		// Drain still ongoing, update index and block for updates
		q.WaitIndex = meta.LastIndex
	}
}

// monitorDrainAllocs emits alloc updates on allocCh and closes the channel
// when the node has finished draining.
func (n *Nodes) monitorDrainAllocs(ctx context.Context, nodeID string, ignoreSys bool, allocCh chan<- *MonitorMessage) {
	defer close(allocCh)

	q := QueryOptions{AllowStale: true}
	initial := make(map[string]*Allocation, 4)

	for {
		allocs, meta, err := n.Allocations(nodeID, &q)
		if err != nil {
			msg := Messagef(MonitorMsgLevelError, "Error monitoring allocations: %v", err)
			select {
			case allocCh <- msg:
			case <-ctx.Done():
			}
			return
		}

		q.WaitIndex = meta.LastIndex

		runningAllocs := 0
		for _, a := range allocs {
			// Get previous version of alloc
			orig, existing := initial[a.ID]

			// Update local alloc state
			initial[a.ID] = a

			migrating := a.DesiredTransition.ShouldMigrate()

			var msg string
			switch {
			case !existing:
				// Should only be possible if response
				// from initial Allocations call was
				// stale. No need to output

			case orig.ClientStatus != a.ClientStatus:
				// Alloc status has changed; output
				msg = fmt.Sprintf("status %s -> %s", orig.ClientStatus, a.ClientStatus)

			case migrating && !orig.DesiredTransition.ShouldMigrate():
				// Alloc was marked for migration
				msg = "marked for migration"

			case migrating && (orig.DesiredStatus != a.DesiredStatus) && a.DesiredStatus == structs.AllocDesiredStatusStop:
				// Alloc has already been marked for migration and is now being stopped
				msg = "draining"
			}

			if msg != "" {
				select {
				case allocCh <- Messagef(MonitorMsgLevelNormal, "Alloc %q %s", a.ID, msg):
				case <-ctx.Done():
					return
				}
			}

			// Ignore malformed allocs
			if a.Job == nil || a.Job.Type == nil {
				continue
			}

			// Track how many allocs are still running
			if ignoreSys && a.Job.Type != nil && *a.Job.Type == structs.JobTypeSystem {
				continue
			}

			switch a.ClientStatus {
			case structs.AllocClientStatusPending, structs.AllocClientStatusRunning:
				runningAllocs++
			}
		}

		// Exit if all allocs are terminal
		if runningAllocs == 0 {
			msg := Messagef(MonitorMsgLevelInfo, "All allocations on node %q have stopped.", nodeID)
			select {
			case allocCh <- msg:
			case <-ctx.Done():
			}
			return
		}
	}
}

// NodeUpdateEligibilityRequest is used to update the drain specification for a node.
type NodeUpdateEligibilityRequest struct {
	// NodeID is the node to update the drain specification for.
	NodeID      string
	Eligibility string
}

// NodeEligibilityUpdateResponse is used to respond to a node eligibility update
type NodeEligibilityUpdateResponse struct {
	NodeModifyIndex uint64
	EvalIDs         []string
	EvalCreateIndex uint64
	WriteMeta
}

// ToggleEligibility is used to update the scheduling eligibility of the node
func (n *Nodes) ToggleEligibility(nodeID string, eligible bool, q *WriteOptions) (*NodeEligibilityUpdateResponse, error) {
	e := structs.NodeSchedulingEligible
	if !eligible {
		e = structs.NodeSchedulingIneligible
	}

	req := &NodeUpdateEligibilityRequest{
		NodeID:      nodeID,
		Eligibility: e,
	}

	var resp NodeEligibilityUpdateResponse
	wm, err := n.client.write("/v1/node/"+nodeID+"/eligibility", req, &resp, q)
	if err != nil {
		return nil, err
	}
	resp.WriteMeta = *wm
	return &resp, nil
}

// Allocations is used to return the allocations associated with a node.
func (n *Nodes) Allocations(nodeID string, q *QueryOptions) ([]*Allocation, *QueryMeta, error) {
	var resp []*Allocation
	qm, err := n.client.query("/v1/node/"+nodeID+"/allocations", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(AllocationSort(resp))
	return resp, qm, nil
}

// ForceEvaluate is used to force-evaluate an existing node.
func (n *Nodes) ForceEvaluate(nodeID string, q *WriteOptions) (string, *WriteMeta, error) {
	var resp nodeEvalResponse
	wm, err := n.client.write("/v1/node/"+nodeID+"/evaluate", nil, &resp, q)
	if err != nil {
		return "", nil, err
	}
	return resp.EvalID, wm, nil
}

func (n *Nodes) Stats(nodeID string, q *QueryOptions) (*HostStats, error) {
	var resp HostStats
	path := fmt.Sprintf("/v1/client/stats?node_id=%s", nodeID)
	if _, err := n.client.query(path, &resp, q); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (n *Nodes) GC(nodeID string, q *QueryOptions) error {
	var resp struct{}
	path := fmt.Sprintf("/v1/client/gc?node_id=%s", nodeID)
	_, err := n.client.query(path, &resp, q)
	return err
}

// TODO Add tests
func (n *Nodes) GcAlloc(allocID string, q *QueryOptions) error {
	var resp struct{}
	path := fmt.Sprintf("/v1/client/allocation/%s/gc", allocID)
	_, err := n.client.query(path, &resp, q)
	return err
}

// DriverInfo is used to deserialize a DriverInfo entry
type DriverInfo struct {
	Attributes        map[string]string
	Detected          bool
	Healthy           bool
	HealthDescription string
	UpdateTime        time.Time
}

// Node is used to deserialize a node entry.
type Node struct {
	ID                    string
	Datacenter            string
	Name                  string
	HTTPAddr              string
	TLSEnabled            bool
	Attributes            map[string]string
	Resources             *Resources
	Reserved              *Resources
	Links                 map[string]string
	Meta                  map[string]string
	NodeClass             string
	Drain                 bool
	DrainStrategy         *DrainStrategy
	SchedulingEligibility string
	Status                string
	StatusDescription     string
	StatusUpdatedAt       int64
	Events                []*NodeEvent
	Drivers               map[string]*DriverInfo
	CreateIndex           uint64
	ModifyIndex           uint64
}

// DrainStrategy describes a Node's drain behavior.
type DrainStrategy struct {
	// DrainSpec is the user declared drain specification
	DrainSpec

	// ForceDeadline is the deadline time for the drain after which drains will
	// be forced
	ForceDeadline time.Time
}

// DrainSpec describes a Node's drain behavior.
type DrainSpec struct {
	// Deadline is the duration after StartTime when the remaining
	// allocations on a draining Node should be told to stop.
	Deadline time.Duration

	// IgnoreSystemJobs allows systems jobs to remain on the node even though it
	// has been marked for draining.
	IgnoreSystemJobs bool
}

func (d *DrainStrategy) Equal(o *DrainStrategy) bool {
	if d == nil || o == nil {
		return d == o
	}

	if d.ForceDeadline != o.ForceDeadline {
		return false
	}
	if d.Deadline != o.Deadline {
		return false
	}
	if d.IgnoreSystemJobs != o.IgnoreSystemJobs {
		return false
	}

	return true
}

// String returns a human readable version of the drain strategy.
func (d *DrainStrategy) String() string {
	if d.IgnoreSystemJobs {
		return fmt.Sprintf("drain ignoring system jobs and deadline at %s", d.ForceDeadline)
	}
	return fmt.Sprintf("drain with deadline at %s", d.ForceDeadline)
}

const (
	NodeEventSubsystemDrain     = "Drain"
	NodeEventSubsystemDriver    = "Driver"
	NodeEventSubsystemHeartbeat = "Heartbeat"
	NodeEventSubsystemCluster   = "Cluster"
)

// NodeEvent is a single unit representing a node’s state change
type NodeEvent struct {
	Message     string
	Subsystem   string
	Details     map[string]string
	Timestamp   time.Time
	CreateIndex uint64
}

// HostStats represents resource usage stats of the host running a Nomad client
type HostStats struct {
	Memory           *HostMemoryStats
	CPU              []*HostCPUStats
	DiskStats        []*HostDiskStats
	Uptime           uint64
	CPUTicksConsumed float64
}

type HostMemoryStats struct {
	Total     uint64
	Available uint64
	Used      uint64
	Free      uint64
}

type HostCPUStats struct {
	CPU    string
	User   float64
	System float64
	Idle   float64
}

type HostDiskStats struct {
	Device            string
	Mountpoint        string
	Size              uint64
	Used              uint64
	Available         uint64
	UsedPercent       float64
	InodesUsedPercent float64
}

// NodeListStub is a subset of information returned during
// node list operations.
type NodeListStub struct {
	Address               string
	ID                    string
	Datacenter            string
	Name                  string
	NodeClass             string
	Version               string
	Drain                 bool
	SchedulingEligibility string
	Status                string
	StatusDescription     string
	Drivers               map[string]*DriverInfo
	CreateIndex           uint64
	ModifyIndex           uint64
}

// NodeIndexSort reverse sorts nodes by CreateIndex
type NodeIndexSort []*NodeListStub

func (n NodeIndexSort) Len() int {
	return len(n)
}

func (n NodeIndexSort) Less(i, j int) bool {
	return n[i].CreateIndex > n[j].CreateIndex
}

func (n NodeIndexSort) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// nodeEvalResponse is used to decode a force-eval.
type nodeEvalResponse struct {
	EvalID string
}

// AllocationSort reverse sorts allocs by CreateIndex.
type AllocationSort []*Allocation

func (a AllocationSort) Len() int {
	return len(a)
}

func (a AllocationSort) Less(i, j int) bool {
	return a[i].CreateIndex > a[j].CreateIndex
}

func (a AllocationSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
