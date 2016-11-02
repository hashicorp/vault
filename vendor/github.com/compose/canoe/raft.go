package canoe

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/cenk/backoff"
	cTypes "github.com/compose/canoe/types"

	"github.com/coreos/etcd/etcdserver/stats"
	"github.com/coreos/etcd/pkg/fileutil"
	eTypes "github.com/coreos/etcd/pkg/types"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/rafthttp"
	"github.com/coreos/etcd/snap"
	"github.com/coreos/etcd/wal"
)

// LogData is the format of data you should expect in Apply operations on the FSM.
// It is also what you should pass to Propose calls to a Node
type LogData []byte

// because WAL and Snap look to see if ANY files exist in the dir
// for confirmation. Meaning that if one or the other is enabled
// but not the other, then checks will fail
var walDirExtension = "/wal"
var snapDirExtension = "/snap"

// Node is a raft node. It is responsible for communicating with all other nodes on the cluster,
// and in general doing all the rafty things
type Node struct {
	node           raft.Node
	raftStorage    *raft.MemoryStorage
	transport      *rafthttp.Transport
	bootstrapPeers []string
	bootstrapNode  bool
	peerMap        map[uint64]cTypes.Peer
	id             uint64
	cid            uint64
	raftPort       int

	configPort int

	raftConfig *raft.Config

	started     bool
	initialized bool
	running     bool

	proposeC chan string
	fsm      FSM

	observers     map[uint64]*Observer
	observersLock sync.RWMutex

	initBackoffArgs *InitializationBackoffArgs
	snapshotConfig  *SnapshotConfig

	dataDir string
	ss      *snap.Snapshotter
	wal     *wal.WAL

	lastConfState *raftpb.ConfState

	stopc chan struct{}

	logger Logger
}

// NodeConfig exposes all the configuration options of a Node
type NodeConfig struct {
	// If not specified or 0, will autogenerate a new UUID
	// It is typically safe to let canoe autogenerate a UUID
	ID uint64

	// If not specified 0x100 will be used
	ClusterID uint64

	FSM               FSM
	RaftPort          int
	ConfigurationPort int

	// BootstrapPeers is a list of peers which we believe to be part of a cluster we wish to join.
	// For now, this list is ignored if the node is marked as a BootstrapNode
	BootstrapPeers []string

	// BootstrapNode is currently needed when bootstrapping a new cluster, a single node must mark itself
	// as the bootstrap node.
	BootstrapNode bool

	// DataDir is where your data will be persisted to disk
	// for use when either you need to restart a node, or
	// it goes offline and needs to be restarted
	DataDir string

	InitBackoff *InitializationBackoffArgs
	// if nil, then default to no snapshotting
	SnapshotConfig *SnapshotConfig

	Logger Logger
}

// Logger is a clone of etcd.Logger interface. We have it cloned in case we want to add more functionality
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warning(v ...interface{})
	Warningf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

// SnapshotConfig defines when you want raft to take a snapshot and compact the WAL
type SnapshotConfig struct {

	// How often do you want to Snapshot and compact logs?
	Interval time.Duration

	// If the interval ticks but not enough logs have been commited then ignore
	// the snapshot this interval
	// This can be useful if you expect your snapshot procedure to have an expensive base cost
	MinCommittedLogs uint64

	// MaxRetainedSnapshots specifies how many snapshots you want to save from
	// purging at a given time
	MaxRetainedSnapshots uint
}

// DefaultSnapshotConfig is what is used for snapshotting when SnapshotConfig isn't specified
// Note: by default we do not snapshot
var DefaultSnapshotConfig = &SnapshotConfig{
	Interval:             -1 * time.Minute,
	MinCommittedLogs:     0,
	MaxRetainedSnapshots: 0,
}

// InitializationBackoffArgs defines the backoff arguments for initializing a Node into a cluster
// as attempts to join or bootstrap a cluster are dependent on other nodes
type InitializationBackoffArgs struct {
	InitialInterval     time.Duration
	Multiplier          float64
	MaxInterval         time.Duration
	MaxElapsedTime      time.Duration
	RandomizationFactor float64
}

// DefaultInitializationBackoffArgs are the default backoff args
var DefaultInitializationBackoffArgs = &InitializationBackoffArgs{
	InitialInterval:     500 * time.Millisecond,
	RandomizationFactor: .5,
	Multiplier:          2,
	MaxInterval:         5 * time.Second,
	MaxElapsedTime:      2 * time.Minute,
}

// UniqueID returns the unique id for the raft node.
// This can be useful to get when defining your state machine so you don't have to
// define a new ID for identification and ownership purposes if your application needs that
func (rn *Node) UniqueID() uint64 {
	return rn.id
}

// NewNode creates a new node from the config options
func NewNode(args *NodeConfig) (*Node, error) {
	// TODO: Look into which config options we want others to specify. For now hardcoded
	// TODO: Allow user to specify KV pairs of known nodes, and bypass the http discovery
	// NOTE: Peers are used EXCLUSIVELY to round-robin to other nodes and attempt to add
	//		ourselves to an existing cluster or bootstrap node
	rn, err := nonInitNode(args)
	if err != nil {
		return nil, err
	}

	return rn, nil
}

func (rn *Node) shouldRejoinCluster() bool {
	return wal.Exist(rn.walDir()) && rn.walDir() != ""
}

func (rn *Node) advanceTicksForElection() error {
	for i := 0; i < rn.raftConfig.ElectionTick-1; i++ {
		rn.node.Tick()
	}
	return nil
}

// Start starts the raft node
func (rn *Node) Start() error {
	// TODO: Intermittent issues with restoring disconnected member from snapshot

	walEnabled := rn.walDir() != ""
	rejoinCluster := rn.shouldRejoinCluster()
	if rn.started {
		return nil
	}

	if walEnabled {
		rn.logger.Info("Initializing persistent storage")
		if err := rn.initPersistentStorage(); err != nil {
			return errors.Wrap(err, "Error initializing persistent storage")
		}
		rn.logger.Info("Finished initializing persistent storage")
	}

	if rejoinCluster {
		rn.logger.Info("Restoring canoe from persistent storage")
		if err := rn.restoreRaft(); err != nil {
			return errors.Wrap(err, "Error restoring raft")
		}
		rn.logger.Info("Finished restoring canoe from persistent storage")

		rn.logger.Info("Restarting canoe node")
		rn.node = raft.RestartNode(rn.raftConfig)
		rn.logger.Info("Successfully restarted canoe node")
	} else {
		// TODO: Fix the mess that is transport initialization
		rn.logger.Info("Attaching transport layer")
		if err := rn.attachTransport(); err != nil {
			return errors.Wrap(err, "Error attaching raft transport")
		}
		rn.logger.Info("Successfully attached transport layer")

		rn.logger.Info("Starting transport layer")
		if err := rn.transport.Start(); err != nil {
			return errors.Wrap(err, "Error starting raft transport")
		}
		rn.logger.Info("Successfully Started transport layer")

		if rn.bootstrapNode {
			rn.logger.Info("Starting node as bootstrap")
			rn.node = raft.StartNode(rn.raftConfig, []raft.Peer{raft.Peer{ID: rn.id}})
		} else {
			rn.logger.Info("Starting node without bootstrap flag")
			rn.node = raft.StartNode(rn.raftConfig, nil)
		}
	}

	rn.logger.Debug("Advancing election ticks")
	if err := rn.advanceTicksForElection(); err != nil {
		return errors.Wrap(err, "Error optimizing election ticks")
	}
	rn.logger.Debug("Successfully advanced election ticks")

	rn.initialized = true

	go func(rn *Node) {
		rn.logger.Info("Scanning for new raft logs")
		if err := rn.scanReady(); err != nil {
			rn.logger.Errorf("%+v", err)
			if errors.Cause(err) == ErrorRemovedFromCluster {
				rn.logger.Info("Trying to destroy canoe data")
				if err := rn.Destroy(); err != nil {
					rn.logger.Fatalf("%+v", err)
				}
				rn.logger.Info("Canoe data destroyed")
				os.Exit(1)
			} else {
				rn.logger.Info("Trying to cleanly stop canoe")
				if err := rn.Stop(); err != nil {
					rn.logger.Fatalf("%+v", err)
				}
				rn.logger.Info("Canoe cleanly stopped")
				os.Exit(1)
			}
		}
	}(rn)

	// periodically cleanup old snapshots
	if rn.snapDir() != "" && rn.snapshotConfig.Interval > 0 && rn.snapshotConfig.MaxRetainedSnapshots > 0 {
		go func(rn *Node) {
			errc := fileutil.PurgeFile(rn.snapDir(), "snap", rn.snapshotConfig.MaxRetainedSnapshots, rn.snapshotConfig.Interval, rn.stopc)
			select {
			case e := <-errc:
				rn.logger.Fatalf("failed to purge snap file %+v", e)
			case <-rn.stopc:
				return
			}
		}(rn)
	}

	// Start config http service
	go func(rn *Node) {
		rn.logger.Info("Starting http config service")
		if err := rn.serveHTTP(); err != nil {
			rn.logger.Fatalf("%+v", err)
		}
	}(rn)

	// start raft
	go func(rn *Node) {
		rn.logger.Info("Starting raft server")
		if err := rn.serveRaft(); err != nil {
			rn.logger.Fatalf("%+v", err)
		}
	}(rn)
	rn.started = true

	// TODO: add case for when no peers or bootstrap specified it waits to get added.
	if rejoinCluster {
		rn.logger.Info("Rejoining canoe cluster")
		if err := rn.selfRejoinCluster(); err != nil {
			return errors.Wrap(err, "Error rejoining raft cluster")
		}
	} else if !rn.bootstrapNode {
		rn.logger.Info("Adding self to existing cluster")
		if err := rn.addSelfToCluster(); err != nil {
			return errors.Wrap(err, "Error adding self to existing raft cluster")
		}
	}

	// final step to mark node as initialized
	rn.running = true
	return nil
}

// IsRunning reports if the raft node is running
func (rn *Node) IsRunning() bool {
	return rn.running
}

// Stop will stop the raft node.
//
// Note: stopping will not remove this node from the cluster. This means that it will affect consensus and quorum
func (rn *Node) Stop() error {
	rn.logger.Info("Stopping canoe")
	close(rn.stopc)

	rn.logger.Debug("Stopping raft transporter")
	rn.transport.Stop()
	// TODO: Don't poll stuff here
	for rn.running {
		time.Sleep(200 * time.Millisecond)
	}
	rn.logger.Info("Canoe has stopped")
	rn.started = false
	rn.initialized = false
	return nil
}

// Destroy is a HARD stop. It first reconfigures the raft cluster
// to remove itself(ONLY do this if you are intending to permenantly leave the cluster and know consequences around consensus) - read the raft paper's reconfiguration section before using this.
// It then halts all running goroutines
//
// WARNING! - Destroy will recursively remove everything under <DataDir>/snap and <DataDir>/wal
func (rn *Node) Destroy() error {
	rn.logger.Debug("Removing self from canoe cluster")
	if err := rn.removeSelfFromCluster(); err != nil {
		return errors.Wrap(err, "Error removing self from existing cluster")
	}
	rn.logger.Debug("Successfully removed self from canoe cluster")

	if rn.running {
		close(rn.stopc)
		rn.logger.Debug("Stopping raft transport layer")
		rn.transport.Stop()
		// TODO: Have a stopped chan for triggering this action
		for rn.running {
			time.Sleep(200 * time.Millisecond)
		}
	}

	rn.logger.Debug("Deleting persistent data")
	rn.deletePersistentData()
	rn.logger.Debug("Successfully deleted persistent data")

	rn.started = false
	rn.initialized = false
	return nil
}

func (rn *Node) removeSelfFromCluster() error {
	notify := func(err error, t time.Duration) {
		rn.logger.Warningf("Couldn't remove self from cluster: %s Trying again in %v", err.Error(), t)
	}

	expBackoff := backoff.NewExponentialBackOff()

	expBackoff.InitialInterval = rn.initBackoffArgs.InitialInterval
	expBackoff.RandomizationFactor = rn.initBackoffArgs.RandomizationFactor
	expBackoff.Multiplier = rn.initBackoffArgs.Multiplier
	expBackoff.MaxInterval = rn.initBackoffArgs.MaxInterval
	expBackoff.MaxElapsedTime = rn.initBackoffArgs.MaxElapsedTime

	op := func() error {
		return rn.requestSelfDeletion()
	}

	return backoff.RetryNotify(op, expBackoff, notify)
}

func (rn *Node) addSelfToCluster() error {
	notify := func(err error, t time.Duration) {
		rn.logger.Warningf("Couldn't add self to cluster: %s Trying again in %v", err.Error(), t)
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = rn.initBackoffArgs.InitialInterval
	expBackoff.RandomizationFactor = rn.initBackoffArgs.RandomizationFactor
	expBackoff.Multiplier = rn.initBackoffArgs.Multiplier
	expBackoff.MaxInterval = rn.initBackoffArgs.MaxInterval
	expBackoff.MaxElapsedTime = rn.initBackoffArgs.MaxElapsedTime

	op := func() error {
		return rn.requestSelfAddition()
	}

	return backoff.RetryNotify(op, expBackoff, notify)
}

func (rn *Node) selfRejoinCluster() error {
	notify := func(err error, t time.Duration) {
		rn.logger.Warningf("Couldn't join cluster: %s Trying again in %v", err.Error(), t)
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = rn.initBackoffArgs.InitialInterval
	expBackoff.RandomizationFactor = rn.initBackoffArgs.RandomizationFactor
	expBackoff.Multiplier = rn.initBackoffArgs.Multiplier
	expBackoff.MaxInterval = rn.initBackoffArgs.MaxInterval
	expBackoff.MaxElapsedTime = rn.initBackoffArgs.MaxElapsedTime

	op := func() error {
		return rn.requestRejoinCluster()
	}

	return backoff.RetryNotify(op, expBackoff, notify)
}

func nonInitNode(args *NodeConfig) (*Node, error) {
	if args.BootstrapNode {
		args.BootstrapPeers = nil
	}

	if args.InitBackoff == nil {
		args.InitBackoff = DefaultInitializationBackoffArgs
	}

	if args.SnapshotConfig == nil {
		args.SnapshotConfig = DefaultSnapshotConfig
	}

	rn := &Node{
		proposeC:        make(chan string),
		raftStorage:     raft.NewMemoryStorage(),
		bootstrapPeers:  args.BootstrapPeers,
		bootstrapNode:   args.BootstrapNode,
		id:              args.ID,
		cid:             args.ClusterID,
		raftPort:        args.RaftPort,
		configPort:      args.ConfigurationPort,
		fsm:             args.FSM,
		initialized:     false,
		observers:       make(map[uint64]*Observer),
		peerMap:         make(map[uint64]cTypes.Peer),
		initBackoffArgs: args.InitBackoff,
		snapshotConfig:  args.SnapshotConfig,
		dataDir:         args.DataDir,
		logger:          args.Logger,
		stopc:           make(chan struct{}),
	}

	if rn.id == 0 {
		rn.id = Uint64UUID()
	}
	if rn.cid == 0 {
		rn.cid = 0x100
	}

	//TODO: Fix these magix numbers with user-specifiable config
	rn.raftConfig = &raft.Config{
		ID:              rn.id,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         rn.raftStorage,
		MaxSizePerMsg:   1024 * 1024,
		MaxInflightMsgs: 256,
		CheckQuorum:     true,
	}

	if rn.logger != nil {
		rn.raftConfig.Logger = raft.Logger(rn.logger)
	} else {
		rn.logger = DefaultLogger
		rn.raftConfig.Logger = rn.logger
	}

	return rn, nil
}

func (rn *Node) attachTransport() error {
	ss := &stats.ServerStats{}
	ss.Initialize()

	//ID TBA on raft restoration creation
	// due to unfortunate dependency on the restore process needing
	rn.transport = &rafthttp.Transport{
		ID:          eTypes.ID(rn.id),
		ClusterID:   eTypes.ID(rn.cid),
		Raft:        rn,
		Snapshotter: rn.ss,
		ServerStats: ss,
		LeaderStats: stats.NewLeaderStats(strconv.FormatUint(rn.id, 10)),
		ErrorC:      make(chan error),
	}

	return nil
}

func (rn *Node) proposePeerAddition(addReq *raftpb.ConfChange, async bool) error {
	addReq.Type = raftpb.ConfChangeAddNode

	observChan := make(chan Observation)
	// setup listener for node addition
	// before asking for node addition
	if !async {
		filterFn := func(o Observation) bool {

			switch o.(type) {
			case raftpb.Entry:
				entry := o.(raftpb.Entry)
				switch entry.Type {
				case raftpb.EntryConfChange:
					var cc raftpb.ConfChange
					cc.Unmarshal(entry.Data)
					rn.node.ApplyConfChange(cc)
					switch cc.Type {
					case raftpb.ConfChangeAddNode:
						// wait until we get a matching node id
						return addReq.NodeID == cc.NodeID
					default:
						return false
					}
				default:
					return false
				}
			default:
				return false
			}
		}

		observer := NewObserver(observChan, filterFn)
		rn.RegisterObserver(observer)
		defer rn.UnregisterObserver(observer)
	}

	if err := rn.node.ProposeConfChange(context.TODO(), *addReq); err != nil {
		return errors.Wrap(err, "Error proposing configuration change")
	}

	if async {
		return nil
	}

	select {
	case <-observChan:
		return nil
	case <-time.After(10 * time.Second):
		return errors.New("Timed out waiting for config change")
	}
}

func (rn *Node) proposePeerDeletion(delReq *raftpb.ConfChange, async bool) error {
	delReq.Type = raftpb.ConfChangeRemoveNode

	observChan := make(chan Observation)
	// setup listener for node addition
	// before asking for node addition
	if !async {
		filterFn := func(o Observation) bool {
			switch o.(type) {
			case raftpb.Entry:
				entry := o.(raftpb.Entry)
				switch entry.Type {
				case raftpb.EntryConfChange:
					var cc raftpb.ConfChange
					cc.Unmarshal(entry.Data)
					rn.node.ApplyConfChange(cc)
					switch cc.Type {
					case raftpb.ConfChangeRemoveNode:
						// wait until we get a matching node id
						return delReq.NodeID == cc.NodeID
					default:
						return false
					}
				default:
					return false
				}
			default:
				return false
			}
		}

		observer := NewObserver(observChan, filterFn)
		rn.RegisterObserver(observer)
		defer rn.UnregisterObserver(observer)
	}

	if err := rn.node.ProposeConfChange(context.TODO(), *delReq); err != nil {
		return errors.Wrap(err, "Error proposing configuration change to raft")
	}

	if async {
		return nil
	}

	select {
	case <-observChan:
		return nil
	case <-time.After(10 * time.Second):
		return errors.Wrap(rn.proposePeerDeletion(delReq, async), "Error proposing peer deletion")

	}
}

func (rn *Node) canAlterPeer() bool {
	return rn.isHealthy() && rn.initialized
}

// TODO: Define healthy better
func (rn *Node) isHealthy() bool {
	return rn.running
}

func (rn *Node) scanReady() error {
	defer func() {
		if rn.wal != nil {
			rn.logger.Info("Closed WAL")
			rn.wal.Close()
		}
	}()
	defer func(rn *Node) {
		rn.running = false
	}(rn)

	var snapTicker *time.Ticker

	// if non-interval based then create a ticker which will never post to a chan
	if rn.snapshotConfig.Interval <= 0 {
		snapTicker = time.NewTicker(1 * time.Second)
		snapTicker.Stop()
	} else {
		snapTicker = time.NewTicker(rn.snapshotConfig.Interval)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-rn.stopc:
			return nil
		case <-ticker.C:
			rn.node.Tick()
		case <-snapTicker.C:
			if err := rn.createSnapAndCompact(false); err != nil {
				return errors.Wrap(err, "Error creating snapshot and compacting WAL")
			}
		case rd := <-rn.node.Ready():
			if rn.wal != nil {
				rn.wal.Save(rd.HardState, rd.Entries)
			}
			rn.raftStorage.Append(rd.Entries)
			rn.transport.Send(rd.Messages)

			if !raft.IsEmptySnap(rd.Snapshot) {
				if err := rn.processSnapshot(rd.Snapshot); err != nil {
					return errors.Wrap(err, "Error processing raft snapshot")
				}
			}

			if err := rn.publishEntries(rd.CommittedEntries); err != nil {
				return errors.Wrap(err, "Error publishing raft entries")
			}

			rn.node.Advance()

		}
	}
}

func (rn *Node) restoreFSMFromSnapshot(raftSnap raftpb.Snapshot) error {
	if raft.IsEmptySnap(raftSnap) {
		return nil
	}

	rn.logger.Info("Restoring FSM from snapshot")
	var snapStruct snapshot
	if err := json.Unmarshal(raftSnap.Data, &snapStruct); err != nil {
		return errors.Wrap(err, "Error unmarshaling raft snapshot")
	}

	rn.logger.Debug("Scanning snapshot for peers")
	for id, info := range snapStruct.Metadata.Peers {
		raftURL := fmt.Sprintf("http://%s", net.JoinHostPort(info.IP, strconv.Itoa(info.RaftPort)))
		rn.logger.Debug("Adding transport peer from Snapshot: %x - %s", id, raftURL)
		rn.transport.AddPeer(eTypes.ID(id), []string{raftURL})
		rn.peerMap[id] = info
	}

	rn.logger.Debug("Inserting raw Snapshot data into FSM")
	if err := rn.fsm.Restore(SnapshotData(snapStruct.Data)); err != nil {
		return errors.Wrap(err, "Error restoring FSM from snapshot when calling external FSM")
	}

	return nil
}

func (rn *Node) processSnapshot(raftSnap raftpb.Snapshot) error {
	if err := rn.restoreFSMFromSnapshot(raftSnap); err != nil {
		return errors.Wrap(err, "Error restoring FSM from snapshot")
	}

	if err := rn.persistSnapshot(raftSnap); err != nil {
		return errors.Wrap(err, "Error persisting snapshot to storage")
	}
	if err := rn.raftStorage.ApplySnapshot(raftSnap); err != nil {
		return errors.Wrap(err, "Error applying snapshot to mem raft storage")
	}

	rn.ReportSnapshot(rn.id, raft.SnapshotFinish)

	return nil
}

type snapshot struct {
	Metadata *snapshotMetadata `json:"metadata"`
	Data     []byte            `json:"data"`
}

type snapshotMetadata struct {
	Peers map[uint64]cTypes.Peer `json:"peers"`
}

// MarshalJSON fulfills the JSON interface
func (p *snapshotMetadata) MarshalJSON() ([]byte, error) {
	tmpStruct := &struct {
		Peers map[string]cTypes.Peer `json:"peers"`
	}{
		Peers: make(map[string]cTypes.Peer),
	}

	for key, val := range p.Peers {
		tmpStruct.Peers[strconv.FormatUint(key, 10)] = val
	}

	return json.Marshal(tmpStruct)
}

// UnmarshalJSON fulfills the JSON interface
func (p *snapshotMetadata) UnmarshalJSON(data []byte) error {
	tmpStruct := &struct {
		Peers map[string]cTypes.Peer `json:"peers"`
	}{}

	if err := json.Unmarshal(data, tmpStruct); err != nil {
		return errors.Wrap(err, "Error unmarshaling snapshot metadata")
	}

	p.Peers = make(map[uint64]cTypes.Peer)

	for key, val := range tmpStruct.Peers {
		convKey, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return errors.Wrap(err, "Error parsing IDs from peer map")
		}
		p.Peers[convKey] = val
	}

	return nil
}

// TODO: Limit to only snapping after min committed
func (rn *Node) createSnapAndCompact(force bool) error {
	index := rn.node.Status().Applied
	lastSnap, err := rn.raftStorage.Snapshot()
	if err != nil {
		return errors.Wrap(err, "Error fetching last snapshot from in memory storage")
	}

	if index <= lastSnap.Metadata.Index && !force {
		return nil
	}

	fsmData, err := rn.fsm.Snapshot()
	if err != nil {
		return errors.Wrap(err, "Error getting snapshot from FSM")
	}

	finalSnap := &snapshot{
		Metadata: &snapshotMetadata{
			Peers: rn.peerMap,
		},
		Data: []byte(fsmData),
	}
	rn.logger.Debug("Snapshot Creating Peers: %v", finalSnap.Metadata.Peers)

	data, err := json.Marshal(finalSnap)
	if err != nil {
		return errors.Wrap(err, "Error marshalling wrapped snapshot")
	}

	rn.logger.Debug("Creating Snapsot")
	raftSnap, err := rn.raftStorage.CreateSnapshot(index, rn.lastConfState, []byte(data))
	if err != nil {
		return errors.Wrap(err, "Error creating snapshot in memory storage")
	}
	rn.logger.Debug("Successfully Created Snapsot")

	rn.logger.Debug("Compacting storage")
	if err = rn.raftStorage.Compact(raftSnap.Metadata.Index); err != nil {
		return errors.Wrap(err, "Error compacting memory storage after snapshot")
	}
	rn.logger.Debug("Successfully compacted storage")

	rn.logger.Debug("Persisting snapshot")
	if err = rn.persistSnapshot(raftSnap); err != nil {
		return errors.Wrap(err, "Error persisting snapshot")
	}
	rn.logger.Debug("Successfully persisted snapshot")

	return nil
}

func (rn *Node) commitsSinceLastSnap() uint64 {
	raftSnap, err := rn.raftStorage.Snapshot()
	if err != nil {
		// this should NEVER err
		panic(err)
	}
	curIndex, err := rn.raftStorage.LastIndex()
	if err != nil {
		// this should NEVER err
		panic(err)
	}
	return curIndex - raftSnap.Metadata.Index
}

// ErrorRemovedFromCluster is returned when an operation failed because this Node
// has been removed from the cluster
var ErrorRemovedFromCluster = errors.New("I have been removed from cluster")

func (rn *Node) publishEntries(ents []raftpb.Entry) error {
	for _, entry := range ents {
		switch entry.Type {
		case raftpb.EntryNormal:
			if len(entry.Data) == 0 {
				break
			}
			// Yes, this is probably a blocking call
			// An FSM should be responsible for being efficient
			// for high-load situations
			if err := rn.fsm.Apply(LogData(entry.Data)); err != nil {
				return errors.Wrap(err, "Error with FSM applying log entry")
			}

		case raftpb.EntryConfChange:
			var cc raftpb.ConfChange
			if err := cc.Unmarshal(entry.Data); err != nil {
				return errors.Wrap(err, "Error unmarshaling ConfChange")
			}
			confState := rn.node.ApplyConfChange(cc)
			rn.lastConfState = confState

			switch cc.Type {
			case raftpb.ConfChangeAddNode:
				if len(cc.Context) > 0 {
					var ctxData cTypes.Peer
					if err := json.Unmarshal(cc.Context, &ctxData); err != nil {
						return errors.Wrap(err, "Error unmarshalling add node request")
					}

					raftURL := fmt.Sprintf("http://%s", net.JoinHostPort(ctxData.IP, strconv.Itoa(ctxData.RaftPort)))

					if cc.NodeID != rn.id {
						rn.logger.Debug("Adding transport peer from raft entry: %x - %s", cc.NodeID, raftURL)
						rn.transport.AddPeer(eTypes.ID(cc.NodeID), []string{raftURL})
					}
					rn.peerMap[cc.NodeID] = ctxData
				}
			case raftpb.ConfChangeRemoveNode:
				if cc.NodeID == uint64(rn.id) {
					return ErrorRemovedFromCluster
				}
				rn.transport.RemovePeer(eTypes.ID(cc.NodeID))
				delete(rn.peerMap, cc.NodeID)
			}

		}
		rn.observe(entry)
	}
	return nil
}

// Propose asks raft to apply the data to the state machine
func (rn *Node) Propose(data []byte) error {
	return rn.node.Propose(context.TODO(), data)
}

// Process fulfills the requirement for rafthttp.Raft interface
func (rn *Node) Process(ctx context.Context, m raftpb.Message) error {
	return rn.node.Step(ctx, m)
}

// TODO: Get these defined

// IsIDRemoved fulfills the requirement for rafthttp.Raft interface
func (rn *Node) IsIDRemoved(id uint64) bool {
	return false
}

// ReportUnreachable fulfills the interface for rafthttp.Raft
func (rn *Node) ReportUnreachable(id uint64) {}

// ReportSnapshot fulfills the requirement for rafthttp.Raft
func (rn *Node) ReportSnapshot(id uint64, status raft.SnapshotStatus) {}
