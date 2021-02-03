// Copyright 2013-2020 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	. "github.com/aerospike/aerospike-client-go/internal/atomic"
	. "github.com/aerospike/aerospike-client-go/logger"
	. "github.com/aerospike/aerospike-client-go/types"
)

const (
	_PARTITIONS = 4096
)

// Node represents an Aerospike Database Server Node
type Node struct {
	cluster            *Cluster
	name               string
	host               *Host
	aliases            atomic.Value //[]*Host
	stats              nodeStats
	_sessionToken      atomic.Value //[]byte
	_sessionExpiration atomic.Value //time.Time

	racks atomic.Value //map[string]int

	// tendConn reserves a connection for tend so that it won't have to
	// wait in queue for connections, since that will cause starvation
	// and the node being dropped under load.
	tendConn     *Connection
	tendConnLock sync.Mutex // All uses of tend connection should be synchronized

	peersGeneration AtomicInt
	peersCount      AtomicInt

	connections     connectionHeap
	connectionCount AtomicInt
	health          AtomicInt //AtomicInteger

	partitionGeneration AtomicInt
	referenceCount      AtomicInt
	failures            AtomicInt
	partitionChanged    AtomicBool
	errorCount          AtomicInt
	rebalanceGeneration AtomicInt

	active AtomicBool

	supportsFloat, supportsBatchIndex, supportsReplicas, supportsGeo, supportsPeers, supportsLUTNow, supportsTruncateNamespace, supportsClusterStable, supportsBitwiseOps AtomicBool
}

// NewNode initializes a server node with connection parameters.
func newNode(cluster *Cluster, nv *nodeValidator) *Node {
	newNode := &Node{
		cluster: cluster,
		name:    nv.name,
		// address: nv.primaryAddress,
		host: nv.primaryHost,

		// Assign host to first IP alias because the server identifies nodes
		// by IP address (not hostname).
		connections:         *newConnectionHeap(cluster.clientPolicy.MinConnectionsPerNode, cluster.clientPolicy.ConnectionQueueSize),
		connectionCount:     *NewAtomicInt(0),
		peersGeneration:     *NewAtomicInt(-1),
		partitionGeneration: *NewAtomicInt(-2),
		referenceCount:      *NewAtomicInt(0),
		failures:            *NewAtomicInt(0),
		active:              *NewAtomicBool(true),
		partitionChanged:    *NewAtomicBool(false),
		errorCount:          *NewAtomicInt(0),
		rebalanceGeneration: *NewAtomicInt(-1),

		supportsFloat:             *NewAtomicBool(nv.supportsFloat),
		supportsBatchIndex:        *NewAtomicBool(nv.supportsBatchIndex),
		supportsReplicas:          *NewAtomicBool(nv.supportsReplicas),
		supportsGeo:               *NewAtomicBool(nv.supportsGeo),
		supportsPeers:             *NewAtomicBool(nv.supportsPeers),
		supportsLUTNow:            *NewAtomicBool(nv.supportsLUTNow),
		supportsTruncateNamespace: *NewAtomicBool(nv.supportsTruncateNamespace),
		supportsClusterStable:     *NewAtomicBool(nv.supportsClusterStable),
		supportsBitwiseOps:        *NewAtomicBool(nv.supportsBitwiseOps),
	}

	newNode.aliases.Store(nv.aliases)
	newNode._sessionToken.Store(nv.sessionToken)
	newNode.racks.Store(map[string]int{})

	// this will reset to zero on first aggregation on the cluster,
	// therefore will only be counted once.
	atomic.AddInt64(&newNode.stats.NodeAdded, 1)

	return newNode
}

// Refresh requests current status from server node, and updates node with the result.
func (nd *Node) Refresh(peers *peers) error {
	if !nd.active.Get() {
		return nil
	}

	atomic.AddInt64(&nd.stats.TendsTotal, 1)

	// Close idleConnections
	defer nd.dropIdleConnections()

	nd.referenceCount.Set(0)

	var infoMap map[string]string
	var err error
	if peers.usePeers.Get() {
		commands := []string{"node", "peers-generation", "partition-generation"}
		if nd.cluster.clientPolicy.RackAware {
			commands = append(commands, "racks:")
		}

		infoMap, err = nd.RequestInfo(&nd.cluster.infoPolicy, commands...)
		if err != nil {
			nd.refreshFailed(err)
			return err
		}

		if err := nd.verifyNodeName(infoMap); err != nil {
			nd.refreshFailed(err)
			return err
		}

		if err := nd.verifyPeersGeneration(infoMap, peers); err != nil {
			nd.refreshFailed(err)
			return err
		}

		if err := nd.verifyPartitionGeneration(infoMap); err != nil {
			nd.refreshFailed(err)
			return err
		}
	} else {
		commands := []string{"node", "partition-generation", nd.cluster.clientPolicy.servicesString()}
		if nd.cluster.clientPolicy.RackAware {
			commands = append(commands, "racks:")
		}

		infoMap, err = nd.RequestInfo(&nd.cluster.infoPolicy, commands...)
		if err != nil {
			nd.refreshFailed(err)
			return err
		}

		if err := nd.verifyNodeName(infoMap); err != nil {
			nd.refreshFailed(err)
			return err
		}

		if err = nd.verifyPartitionGeneration(infoMap); err != nil {
			nd.refreshFailed(err)
			return err
		}

		if err = nd.addFriends(infoMap, peers); err != nil {
			nd.refreshFailed(err)
			return err
		}
	}

	if err := nd.updateRackInfo(infoMap); err != nil {
		// Update rack info should fail if the feature is not supported on the server
		if aerr, ok := err.(AerospikeError); ok && aerr.ResultCode() == UNSUPPORTED_FEATURE {
			nd.refreshFailed(err)
			return err
		}
		// Should not fail in other cases
		Logger.Warn("Updating node rack info failed with error: %s (racks: `%s`)", err, infoMap["racks:"])
	}

	nd.failures.Set(0)
	peers.refreshCount.IncrementAndGet()
	nd.referenceCount.IncrementAndGet()
	atomic.AddInt64(&nd.stats.TendsSuccessful, 1)

	if err := nd.refreshSessionToken(); err != nil {
		Logger.Error("Error refreshing session token: %s", err.Error())
	}

	if _, err := nd.fillMinConns(); err != nil {
		Logger.Error("Error filling up the connection queue to the minimum required")
	}

	return nil
}

// refreshSessionToken refreshes the session token if it has been expired
func (nd *Node) refreshSessionToken() error {
	// no session token to refresh
	if !nd.cluster.clientPolicy.RequiresAuthentication() || nd.cluster.clientPolicy.AuthMode != AuthModeExternal {
		return nil
	}

	var deadline time.Time
	deadlineIfc := nd._sessionExpiration.Load()
	if deadlineIfc != nil {
		deadline = deadlineIfc.(time.Time)
	}

	if deadline.IsZero() || time.Now().Before(deadline) {
		return nil
	}

	nd.tendConnLock.Lock()
	defer nd.tendConnLock.Unlock()

	if err := nd.initTendConn(nd.cluster.clientPolicy.LoginTimeout); err != nil {
		return err
	}

	command := newLoginCommand(nd.tendConn.dataBuffer)
	if err := command.login(&nd.cluster.clientPolicy, nd.tendConn, nd.cluster.Password()); err != nil {
		// Socket not authenticated. Do not put back into pool.
		nd.tendConn.Close()
		return err
	}

	nd._sessionToken.Store(command.SessionToken)
	nd._sessionExpiration.Store(command.SessionExpiration)

	return nil
}

func (nd *Node) updateRackInfo(infoMap map[string]string) error {
	if !nd.cluster.clientPolicy.RackAware {
		return nil
	}

	// Do not raise an error if the server does not support rackaware
	if strings.HasPrefix(strings.ToUpper(infoMap["racks:"]), "ERROR") {
		return NewAerospikeError(UNSUPPORTED_FEATURE, "You have set the ClientPolicy.RackAware = true, but the server does not support this feature.")
	}

	ss := strings.Split(infoMap["racks:"], ";")
	racks := map[string]int{}
	for _, s := range ss {
		in := bufio.NewReader(strings.NewReader(s))
		_, err := in.ReadString('=')
		if err != nil {
			return err
		}

		ns, err := in.ReadString(':')
		if err != nil {
			return err
		}

		for {
			_, err = in.ReadString('_')
			if err != nil {
				return err
			}

			rackStr, err := in.ReadString('=')
			if err != nil {
				return err
			}

			rack, err := strconv.Atoi(rackStr[:len(rackStr)-1])
			if err != nil {
				return err
			}

			nodesList, err := in.ReadString(':')
			if err != nil && err != io.EOF {
				return err
			}

			nodes := strings.Split(strings.Trim(nodesList, ":"), ",")
			for i := range nodes {
				if nodes[i] == nd.name {
					racks[ns[:len(ns)-1]] = rack
				}
			}

			if err == io.EOF {
				break
			}
		}
	}

	nd.racks.Store(racks)

	return nil
}

func (nd *Node) verifyNodeName(infoMap map[string]string) error {
	infoName, exists := infoMap["node"]

	if !exists || len(infoName) == 0 {
		return NewAerospikeError(INVALID_NODE_ERROR, "Node name is empty")
	}

	if !(nd.name == infoName) {
		// Set node to inactive immediately.
		nd.active.Set(false)
		return NewAerospikeError(INVALID_NODE_ERROR, "Node name has changed. Old="+nd.name+" New="+infoName)
	}
	return nil
}

func (nd *Node) verifyPeersGeneration(infoMap map[string]string, peers *peers) error {
	genString := infoMap["peers-generation"]
	if len(genString) == 0 {
		return NewAerospikeError(PARSE_ERROR, "peers-generation is empty")
	}

	gen, err := strconv.Atoi(genString)
	if err != nil {
		return NewAerospikeError(PARSE_ERROR, "peers-generation is not a number: "+genString)
	}

	peers.genChanged.Or(nd.peersGeneration.Get() != gen)
	return nil
}

func (nd *Node) verifyPartitionGeneration(infoMap map[string]string) error {
	genString := infoMap["partition-generation"]

	if len(genString) == 0 {
		return NewAerospikeError(PARSE_ERROR, "partition-generation is empty")
	}

	gen, err := strconv.Atoi(genString)
	if err != nil {
		return NewAerospikeError(PARSE_ERROR, "partition-generation is not a number:"+genString)
	}

	if nd.partitionGeneration.Get() != gen {
		nd.partitionChanged.Set(true)
	}
	return nil
}

func (nd *Node) addFriends(infoMap map[string]string, peers *peers) error {
	friendString, exists := infoMap[nd.cluster.clientPolicy.servicesString()]

	if !exists || len(friendString) == 0 {
		nd.peersCount.Set(0)
		return nil
	}

	friendNames := strings.Split(friendString, ";")
	nd.peersCount.Set(len(friendNames))

	for _, friend := range friendNames {
		friendInfo := strings.Split(friend, ":")

		if len(friendInfo) != 2 {
			Logger.Error("Node info from asinfo:services is malformed. Expected HOST:PORT, but got `%s`", friend)
			continue
		}

		hostName := friendInfo[0]
		port, _ := strconv.Atoi(friendInfo[1])

		if len(nd.cluster.clientPolicy.IpMap) > 0 {
			if alternativeHost, ok := nd.cluster.clientPolicy.IpMap[hostName]; ok {
				hostName = alternativeHost
			}
		}

		host := NewHost(hostName, port)
		node := nd.cluster.findAlias(host)

		if node != nil {
			node.referenceCount.IncrementAndGet()
		} else {
			if !peers.hostExists(*host) {
				nd.prepareFriend(host, peers)
			}
		}
	}

	return nil
}

func (nd *Node) prepareFriend(host *Host, peers *peers) bool {
	nv := &nodeValidator{}
	if err := nv.validateNode(nd.cluster, host); err != nil {
		Logger.Warn("Adding node `%s` failed: %s", host, err)
		return false
	}

	node := peers.nodeByName(nv.name)

	if node != nil {
		// Duplicate node name found.  This usually occurs when the server
		// services list contains both internal and external IP addresses
		// for the same node.
		peers.addHost(*host)
		node.addAlias(host)
		return true
	}

	// Check for duplicate nodes in cluster.
	node = nd.cluster.nodesMap.Get().(map[string]*Node)[nv.name]

	if node != nil {
		peers.addHost(*host)
		node.addAlias(host)
		node.referenceCount.IncrementAndGet()
		nd.cluster.addAlias(host, node)
		return true
	}

	node = nd.cluster.createNode(nv)
	peers.addHost(*host)
	peers.addNode(nv.name, node)
	return true
}

func (nd *Node) refreshPeers(peers *peers) {
	// Do not refresh peers when node connection has already failed during this cluster tend iteration.
	if nd.failures.Get() > 0 || !nd.active.Get() {
		return
	}

	peerParser, err := parsePeers(nd.cluster, nd)
	if err != nil {
		Logger.Debug("Parsing peers failed: %s", err)
		nd.refreshFailed(err)
		return
	}

	peers.appendPeers(peerParser.peers)
	nd.peersGeneration.Set(int(peerParser.generation()))
	nd.peersCount.Set(len(peers.peers()))
	peers.refreshCount.IncrementAndGet()
}

func (nd *Node) refreshPartitions(peers *peers, partitions partitionMap) {
	// Do not refresh peers when node connection has already failed during this cluster tend iteration.
	// Also, avoid "split cluster" case where this node thinks it's a 1-node cluster.
	// Unchecked, such a node can dominate the partition map and cause all other
	// nodes to be dropped.
	if nd.failures.Get() > 0 || !nd.active.Get() || (nd.peersCount.Get() == 0 && peers.refreshCount.Get() > 1) {
		return
	}

	parser, err := newPartitionParser(nd, partitions, _PARTITIONS)
	if err != nil {
		nd.refreshFailed(err)
		return
	}

	if parser.generation != nd.partitionGeneration.Get() {
		Logger.Info("Node %s partition generation changed from %d to %d", nd.host.String(), nd.partitionGeneration.Get(), parser.getGeneration())
		nd.partitionChanged.Set(true)
		nd.partitionGeneration.Set(parser.getGeneration())
		atomic.AddInt64(&nd.stats.PartitionMapUpdates, 1)
	}
}

func (nd *Node) refreshFailed(e error) {
	nd.peersGeneration.Set(-1)
	nd.partitionGeneration.Set(-1)

	if nd.cluster.clientPolicy.RackAware {
		nd.rebalanceGeneration.Set(-1)
	}

	nd.failures.IncrementAndGet()
	atomic.AddInt64(&nd.stats.TendsFailed, 1)

	// Only log message if cluster is still active.
	if nd.cluster.IsConnected() {
		Logger.Warn("Node `%s` refresh failed: `%s`", nd, e)
	}
}

// dropIdleConnections picks a connection from the head of the connection pool queue
// if that connection is idle, it drops it and takes the next one until it picks
// a fresh connection or exhaust the queue.
func (nd *Node) dropIdleConnections() {
	nd.connections.DropIdle()
}

// GetConnection gets a connection to the node.
// If no pooled connection is available, a new connection will be created, unless
// ClientPolicy.MaxQueueSize number of connections are already created.
// This method will retry to retrieve a connection in case the connection pool
// is empty, until timeout is reached.
func (nd *Node) GetConnection(timeout time.Duration) (conn *Connection, err error) {
	if timeout <= 0 {
		timeout = _DEFAULT_TIMEOUT
	}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err = nd.getConnection(deadline, timeout)
		if err == nil && conn != nil {
			return conn, nil
		}

		if err == ErrServerNotAvailable {
			return nil, err
		}

		time.Sleep(5 * time.Millisecond)
	}

	// in case the block didn't run at all
	if err == nil {
		err = ErrConnectionPoolEmpty
	}

	return nil, err
}

// getConnection gets a connection to the node.
// If no pooled connection is available, a new connection will be created.
func (nd *Node) getConnection(deadline time.Time, timeout time.Duration) (conn *Connection, err error) {
	return nd.getConnectionWithHint(deadline, timeout, 0)
}

// newConnection will make a new connection for the node.
func (nd *Node) newConnection(overrideThreshold bool) (*Connection, error) {
	if !nd.active.Get() {
		return nil, ErrServerNotAvailable
	}

	// if connection count is limited and enough connections are already created, don't create a new one
	cc := nd.connectionCount.IncrementAndGet()
	if nd.cluster.clientPolicy.LimitConnectionsToQueueSize && cc > nd.cluster.clientPolicy.ConnectionQueueSize {
		nd.connectionCount.DecrementAndGet()
		atomic.AddInt64(&nd.stats.ConnectionsPoolEmpty, 1)

		return nil, ErrTooManyConnectionsForNode
	}

	// Check for opening connection threshold
	if !overrideThreshold && nd.cluster.clientPolicy.OpeningConnectionThreshold > 0 {
		ct := nd.cluster.connectionThreshold.IncrementAndGet()
		if ct > nd.cluster.clientPolicy.OpeningConnectionThreshold {
			nd.cluster.connectionThreshold.DecrementAndGet()
			nd.connectionCount.DecrementAndGet()

			return nil, ErrTooManyOpeningConnections
		}

		defer nd.cluster.connectionThreshold.DecrementAndGet()
	}

	atomic.AddInt64(&nd.stats.ConnectionsAttempts, 1)
	conn, err := NewConnection(&nd.cluster.clientPolicy, nd.host)
	if err != nil {
		nd.incrErrorCount()
		nd.connectionCount.DecrementAndGet()
		atomic.AddInt64(&nd.stats.ConnectionsFailed, 1)
		return nil, err
	}
	conn.node = nd

	// need to authenticate
	if err = conn.login(&nd.cluster.clientPolicy, nd.cluster.Password(), nd.sessionToken()); err != nil {
		// increment node errors if authenitocation hit a network error
		if networkError(err) {
			nd.incrErrorCount()
		}
		atomic.AddInt64(&nd.stats.ConnectionsFailed, 1)

		// Socket not authenticated. Do not put back into pool.
		conn.Close()
		return nil, err
	}

	atomic.AddInt64(&nd.stats.ConnectionsSuccessful, 1)
	conn.setIdleTimeout(nd.cluster.clientPolicy.IdleTimeout)

	return conn, nil
}

// makeConnectionForPool will try to open a connection until deadline.
// if no deadline is defined, it will only try for _DEFAULT_TIMEOUT.
func (nd *Node) makeConnectionForPool(hint byte) {
	conn, err := nd.newConnection(false)
	if err != nil {
		Logger.Debug("Error trying to make a connection to the node %s: %s", nd.String(), err.Error())
		return
	}

	nd.putConnectionWithHint(conn, hint)
}

// getConnectionWithHint gets a connection to the node.
// If no pooled connection is available, a new connection will be created.
func (nd *Node) getConnectionWithHint(deadline time.Time, timeout time.Duration, hint byte) (conn *Connection, err error) {
	if !nd.active.Get() {
		return nil, ErrServerNotAvailable
	}

	// try to get a valid connection from the connection pool
	for conn = nd.connections.Poll(hint); conn != nil; conn = nd.connections.Poll(hint) {
		if conn.IsConnected() {
			break
		}
		conn.Close()
		conn = nil
	}

	if conn == nil {
		go nd.makeConnectionForPool(hint)
		return nil, ErrConnectionPoolEmpty
	}

	if err = conn.SetTimeout(deadline, timeout); err != nil {
		atomic.AddInt64(&nd.stats.ConnectionsFailed, 1)

		// Do not put back into pool.
		conn.Close()
		return nil, err
	}

	conn.refresh()

	return conn, nil
}

// PutConnection puts back a connection to the pool.
// If connection pool is full, the connection will be
// closed and discarded.
func (nd *Node) putConnectionWithHint(conn *Connection, hint byte) bool {
	conn.refresh()
	if !nd.active.Get() || !nd.connections.Offer(conn, hint) {
		conn.Close()
		return false
	}
	return true
}

// PutConnection puts back a connection to the pool.
// If connection pool is full, the connection will be
// closed and discarded.
func (nd *Node) PutConnection(conn *Connection) {
	nd.putConnectionWithHint(conn, 0)
}

// InvalidateConnection closes and discards a connection from the pool.
func (nd *Node) InvalidateConnection(conn *Connection) {
	conn.Close()
}

// GetHost retrieves host for the node.
func (nd *Node) GetHost() *Host {
	return nd.host
}

// IsActive Checks if the node is active.
func (nd *Node) IsActive() bool {
	return nd != nil && nd.active.Get() && nd.partitionGeneration.Get() >= -1
}

// GetName returns node name.
func (nd *Node) GetName() string {
	return nd.name
}

// GetAliases returns node aliases.
func (nd *Node) GetAliases() []*Host {
	return nd.aliases.Load().([]*Host)
}

// Sets node aliases
func (nd *Node) setAliases(aliases []*Host) {
	nd.aliases.Store(aliases)
}

// AddAlias adds an alias for the node
func (nd *Node) addAlias(aliasToAdd *Host) {
	// Aliases are only referenced in the cluster tend goroutine,
	// so synchronization is not necessary.
	aliases := nd.GetAliases()
	if aliases == nil {
		aliases = []*Host{}
	}

	aliases = append(aliases, aliasToAdd)
	nd.setAliases(aliases)
}

// Close marks node as inactive and closes all of its pooled connections.
func (nd *Node) Close() {
	if nd.active.Get() {
		nd.active.Set(false)
		atomic.AddInt64(&nd.stats.NodeRemoved, 1)
	}
	nd.closeConnections()
	nd.connections.cleanup()
}

// String implements stringer interface
func (nd *Node) String() string {
	return nd.name + " " + nd.host.String()
}

func (nd *Node) closeConnections() {
	for conn := nd.connections.Poll(0); conn != nil; conn = nd.connections.Poll(0) {
		conn.Close()
	}

	// close the tend connection
	nd.tendConnLock.Lock()
	defer nd.tendConnLock.Unlock()
	if nd.tendConn != nil {
		nd.tendConn.Close()
	}
}

// Equals compares equality of two nodes based on their names.
func (nd *Node) Equals(other *Node) bool {
	return nd != nil && other != nil && (nd == other || nd.name == other.name)
}

// MigrationInProgress determines if the node is participating in a data migration
func (nd *Node) MigrationInProgress() (bool, error) {
	values, err := nd.RequestStats(&nd.cluster.infoPolicy)
	if err != nil {
		return false, err
	}

	// if the migrate_partitions_remaining exists and is not `0`, then migration is in progress
	if migration, exists := values["migrate_partitions_remaining"]; exists && migration != "0" {
		return true, nil
	}

	// migration not in progress
	return false, nil
}

// WaitUntillMigrationIsFinished will block until migration operations are finished.
func (nd *Node) WaitUntillMigrationIsFinished(timeout time.Duration) (err error) {
	if timeout <= 0 {
		timeout = _NO_TIMEOUT
	}
	done := make(chan error)

	go func() {
		// this function is guaranteed to return after timeout
		// no go routines will be leaked
		for {
			if res, err := nd.MigrationInProgress(); err != nil || !res {
				done <- err
				return
			}
		}
	}()

	dealine := time.After(timeout)
	select {
	case <-dealine:
		return NewAerospikeError(TIMEOUT)
	case err = <-done:
		return err
	}
}

// initTendConn sets up a connection to be used for info requests.
// The same connection will be used for tend.
func (nd *Node) initTendConn(timeout time.Duration) error {
	if timeout <= 0 {
		timeout = _DEFAULT_TIMEOUT
	}
	deadline := time.Now().Add(timeout)

	if nd.tendConn == nil || !nd.tendConn.IsConnected() {
		var tendConn *Connection
		var err error
		if nd.connectionCount.Get() == 0 {
			// if there are no connections in the pool, create a new connection synchronously.
			// this will make sure the initial tend will get a connection without multiple retries.
			tendConn, err = nd.newConnection(true)
		} else {
			tendConn, err = nd.GetConnection(timeout)
		}

		if err != nil {
			return err
		}
		nd.tendConn = tendConn
	}

	// Set timeout for tend conn
	return nd.tendConn.SetTimeout(deadline, timeout)
}

// requestInfoWithRetry gets info values by name from the specified database server node.
// It will try at least N times before returning an error.
func (nd *Node) requestInfoWithRetry(policy *InfoPolicy, n int, name ...string) (res map[string]string, err error) {
	for i := 0; i < n; i++ {
		if res, err = nd.requestInfo(policy.Timeout, name...); err == nil {
			return res, nil
		}

		Logger.Error("Error occurred while fetching info from the server node %s: %s", nd.host.String(), err.Error())
		time.Sleep(100 * time.Millisecond)
	}

	// return the last error
	return nil, err
}

// RequestInfo gets info values by name from the specified database server node.
func (nd *Node) RequestInfo(policy *InfoPolicy, name ...string) (map[string]string, error) {
	return nd.requestInfo(policy.Timeout, name...)
}

// RequestInfo gets info values by name from the specified database server node.
func (nd *Node) requestInfo(timeout time.Duration, name ...string) (map[string]string, error) {
	nd.tendConnLock.Lock()
	defer nd.tendConnLock.Unlock()

	if err := nd.initTendConn(timeout); err != nil {
		return nil, err
	}

	response, err := RequestInfo(nd.tendConn, name...)
	if err != nil {
		nd.tendConn.Close()
		return nil, err
	}
	return response, nil
}

// requestRawInfo gets info values by name from the specified database server node.
// It won't parse the results.
func (nd *Node) requestRawInfo(policy *InfoPolicy, name ...string) (*info, error) {
	nd.tendConnLock.Lock()
	defer nd.tendConnLock.Unlock()

	if err := nd.initTendConn(policy.Timeout); err != nil {
		return nil, err
	}

	response, err := newInfo(nd.tendConn, name...)
	if err != nil {
		nd.tendConn.Close()
		return nil, err
	}
	return response, nil
}

// RequestStats returns statistics for the specified node as a map
func (node *Node) RequestStats(policy *InfoPolicy) (map[string]string, error) {
	infoMap, err := node.RequestInfo(policy, "statistics")
	if err != nil {
		return nil, err
	}

	res := map[string]string{}

	v, exists := infoMap["statistics"]
	if !exists {
		return res, nil
	}

	values := strings.Split(v, ";")
	for i := range values {
		kv := strings.Split(values[i], "=")
		if len(kv) > 1 {
			res[kv[0]] = kv[1]
		}
	}

	return res, nil
}

// sessionToken returns the session token for the node.
// It will return nil if the session has expired.
func (nd *Node) sessionToken() []byte {
	var deadline time.Time
	deadlineIfc := nd._sessionExpiration.Load()
	if deadlineIfc != nil {
		deadline = deadlineIfc.(time.Time)
	}

	if deadline.IsZero() || time.Now().After(deadline) {
		return nil
	}

	st := nd._sessionToken.Load()
	if st != nil {
		return st.([]byte)
	}
	return nil
}

// Rack returns the rack number for the namespace.
func (nd *Node) Rack(namespace string) (int, error) {
	racks := nd.racks.Load().(map[string]int)
	v, exists := racks[namespace]

	if exists {
		return v, nil
	}

	return -1, newAerospikeNodeError(nd, RACK_NOT_DEFINED)
}

// Rack returns the rack number for the namespace.
func (nd *Node) hasRack(namespace string, rack int) bool {
	racks := nd.racks.Load().(map[string]int)
	v, exists := racks[namespace]

	if !exists {
		return false
	}

	return v == rack
}

// WarmUp fills the node's connection pool with connections.
// This is necessary on startup for high traffic programs.
// If the count is <= 0, the connection queue will be filled.
// If the count is more than the size of the pool, the pool will be filled.
// Note: One connection per node is reserved for tend operations and is not used for transactions.
func (nd *Node) WarmUp(count int) (int, error) {
	var g errgroup.Group
	cnt := NewAtomicInt(0)

	toAlloc := nd.connections.Cap() - nd.connectionCount.Get()
	if count < toAlloc && count > 0 {
		toAlloc = count
	}

	for i := 0; i < toAlloc; i++ {
		g.Go(func() error {
			conn, err := nd.newConnection(true)
			if err != nil {
				if err == ErrTooManyConnectionsForNode {
					return nil
				}
				return err
			}

			if nd.putConnectionWithHint(conn, 0) {
				cnt.IncrementAndGet()
			} else {
				conn.Close()
			}

			return nil
		})
	}

	err := g.Wait()
	return cnt.Get(), err
}

// fillMinCounts will fill the connection pool to the minimum required
// by the ClientPolicy.MinConnectionsPerNode
func (nd *Node) fillMinConns() (int, error) {
	if nd.cluster.clientPolicy.MinConnectionsPerNode > 0 {
		toFill := nd.cluster.clientPolicy.MinConnectionsPerNode - nd.connectionCount.Get()
		if toFill > 0 {
			return nd.WarmUp(toFill)
		}
	}
	return 0, nil
}

// Increments error count for the node. If errorCount goes above the threshold,
// the node will not accept any more requests until the next window.
func (nd *Node) incrErrorCount() {
	if nd.cluster.clientPolicy.MaxErrorRate > 0 {
		nd.errorCount.GetAndIncrement()
	}
}

// Resets the error count
func (nd *Node) resetErrorCount() {
	nd.errorCount.Set(0)
}

// checks if the errorCount is within set limits
func (nd *Node) errorCountWithinLimit() bool {
	return nd.cluster.clientPolicy.MaxErrorRate <= 0 || nd.errorCount.Get() <= nd.cluster.clientPolicy.MaxErrorRate
}

// returns error if errorCount has gone above the threshold set in the policy
func (nd *Node) validateErrorCount() error {
	if !nd.errorCountWithinLimit() {
		return NewAerospikeError(MAX_ERROR_RATE)
	}
	return nil
}
