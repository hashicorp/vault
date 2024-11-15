// Copyright 2014-2021 Aerospike, Inc.
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
	"errors"
	"io"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	iatomic "github.com/aerospike/aerospike-client-go/v5/internal/atomic"
	"github.com/aerospike/aerospike-client-go/v5/logger"
	"github.com/aerospike/aerospike-client-go/v5/types"
)

const (
	_PARTITIONS = 4096

	_SUPPORTS_PARTITION_SCAN = 1 << 0
	_SUPPORTS_QUERY_SHOW     = 1 << 1
)

// Node represents an Aerospike Database Server Node
type Node struct {
	cluster     *Cluster
	name        string
	host        *Host
	aliases     atomic.Value //[]*Host
	stats       nodeStats
	sessionInfo atomic.Value //*sessionInfo

	racks atomic.Value //map[string]int

	// tendConn reserves a connection for tend so that it won't have to
	// wait in queue for connections, since that will cause starvation
	// and the node being dropped under load.
	tendConn     *Connection
	tendConnLock sync.Mutex // All uses of tend connection should be synchronized

	peersGeneration iatomic.Int
	peersCount      iatomic.Int

	connections     connectionHeap
	connectionCount iatomic.Int

	partitionGeneration iatomic.Int
	referenceCount      iatomic.Int
	failures            iatomic.Int
	partitionChanged    iatomic.Bool
	errorCount          iatomic.Int
	rebalanceGeneration iatomic.Int

	features int

	active iatomic.Bool
}

// NewNode initializes a server node with connection parameters.
func newNode(cluster *Cluster, nv *nodeValidator) *Node {
	newNode := &Node{
		cluster: cluster,
		name:    nv.name,
		host:    nv.primaryHost,

		features: nv.features,

		// Assign host to first IP alias because the server identifies nodes
		// by IP address (not hostname).
		connections:         *newConnectionHeap(cluster.clientPolicy.MinConnectionsPerNode, cluster.clientPolicy.ConnectionQueueSize),
		connectionCount:     *iatomic.NewInt(0),
		peersGeneration:     *iatomic.NewInt(-1),
		partitionGeneration: *iatomic.NewInt(-2),
		referenceCount:      *iatomic.NewInt(0),
		failures:            *iatomic.NewInt(0),
		active:              *iatomic.NewBool(true),
		partitionChanged:    *iatomic.NewBool(false),
		errorCount:          *iatomic.NewInt(0),
		rebalanceGeneration: *iatomic.NewInt(-1),
	}

	newNode.aliases.Store(nv.aliases)
	newNode.sessionInfo.Store(nv.sessionInfo)
	newNode.racks.Store(map[string]int{})

	// this will reset to zero on first aggregation on the cluster,
	// therefore will only be counted once.
	atomic.AddInt64(&newNode.stats.NodeAdded, 1)

	return newNode
}

// Refresh requests current status from server node, and updates node with the result.
func (nd *Node) SupportsQueryShow() bool {
	return (nd.features & _SUPPORTS_QUERY_SHOW) != 0
}

// Refresh requests current status from server node, and updates node with the result.
func (nd *Node) Refresh(peers *peers) Error {
	if !nd.active.Get() {
		return nil
	}

	atomic.AddInt64(&nd.stats.TendsTotal, 1)

	// Close idleConnections
	defer nd.dropIdleConnections()

	nd.referenceCount.Set(0)

	var infoMap map[string]string
	commands := []string{"node", "peers-generation", "partition-generation"}
	if nd.cluster.clientPolicy.RackAware {
		commands = append(commands, "racks:")
	}

	infoMap, err := nd.RequestInfo(&nd.cluster.infoPolicy, commands...)
	if err != nil {
		nd.refreshFailed(err)
		return err
	}

	if err = nd.verifyNodeName(infoMap); err != nil {
		nd.refreshFailed(err)
		return err
	}

	if err = nd.verifyPeersGeneration(infoMap, peers); err != nil {
		nd.refreshFailed(err)
		return err
	}

	if err = nd.verifyPartitionGeneration(infoMap); err != nil {
		nd.refreshFailed(err)
		return err
	}

	if err = nd.updateRackInfo(infoMap); err != nil {
		// Update rack info should fail if the feature is not supported on the server
		if err.Matches(types.UNSUPPORTED_FEATURE) {
			nd.refreshFailed(err)
			return err
		}
		// Should not fail in other cases
		logger.Logger.Warn("Updating node rack info failed with error: %s (racks: `%s`)", err, infoMap["racks:"])
	}

	nd.failures.Set(0)
	peers.refreshCount.IncrementAndGet()
	nd.referenceCount.IncrementAndGet()
	atomic.AddInt64(&nd.stats.TendsSuccessful, 1)

	if err = nd.refreshSessionToken(); err != nil {
		logger.Logger.Error("Error refreshing session token: %s", err.Error())
	}

	if _, err = nd.fillMinConns(); err != nil {
		logger.Logger.Error("Error filling up the connection queue to the minimum required")
	}

	return nil
}

// refreshSessionToken refreshes the session token if it has been expired
func (nd *Node) refreshSessionToken() Error {
	// no session token to refresh
	if !nd.cluster.clientPolicy.RequiresAuthentication() {
		return nil
	}

	st := nd.sessionInfo.Load().(*sessionInfo)

	// Consider when the next tend will be in this calculation. If the next tend will be too late,
	// refresh the sessionInfo now.
	if st.expiration.IsZero() || time.Now().Before(st.expiration.Add(-nd.cluster.clientPolicy.TendInterval)) {
		return nil
	}

	nd.tendConnLock.Lock()
	defer nd.tendConnLock.Unlock()

	if err := nd.initTendConn(nd.cluster.clientPolicy.LoginTimeout); err != nil {
		return err
	}

	command := newLoginCommand(nd.tendConn.dataBuffer)
	if err := command.login(&nd.cluster.clientPolicy, nd.tendConn, nd.cluster.Password()); err != nil {
		// force new connections to use default creds until a new valid session token is acquired
		nd.resetSessionInfo()
		// Socket not authenticated. Do not put back into pool.
		nd.tendConn.Close()
		return err
	}

	nd.sessionInfo.Store(command.sessionInfo())
	return nil
}

func (nd *Node) updateRackInfo(infoMap map[string]string) Error {
	if !nd.cluster.clientPolicy.RackAware {
		return nil
	}

	// Do not raise an error if the server does not support rackaware
	if strings.HasPrefix(strings.ToUpper(infoMap["racks:"]), "ERROR") {
		return newError(types.UNSUPPORTED_FEATURE, "You have set the ClientPolicy.RackAware = true, but the server does not support this feature.")
	}

	ss := strings.Split(infoMap["racks:"], ";")
	racks := map[string]int{}
	for _, s := range ss {
		in := bufio.NewReader(strings.NewReader(s))
		_, err := in.ReadString('=')
		if err != nil {
			return newErrorAndWrap(err, types.PARSE_ERROR)
		}

		ns, err := in.ReadString(':')
		if err != nil {
			return newErrorAndWrap(err, types.PARSE_ERROR)
		}

		for {
			_, err = in.ReadString('_')
			if err != nil {
				return newErrorAndWrap(err, types.PARSE_ERROR)
			}

			rackStr, err := in.ReadString('=')
			if err != nil {
				return newErrorAndWrap(err, types.PARSE_ERROR)
			}

			rack, err := strconv.Atoi(rackStr[:len(rackStr)-1])
			if err != nil {
				return newErrorAndWrap(err, types.PARSE_ERROR)
			}

			nodesList, err := in.ReadString(':')
			if err != nil && err != io.EOF {
				return newErrorAndWrap(err, types.PARSE_ERROR)
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

func (nd *Node) verifyNodeName(infoMap map[string]string) Error {
	infoName, exists := infoMap["node"]

	if !exists || len(infoName) == 0 {
		return newError(types.INVALID_NODE_ERROR, "Node name is empty")
	}

	if !(nd.name == infoName) {
		// Set node to inactive immediately.
		nd.active.Set(false)
		return newError(types.INVALID_NODE_ERROR, "Node name has changed. Old="+nd.name+" New="+infoName)
	}
	return nil
}

func (nd *Node) verifyPeersGeneration(infoMap map[string]string, peers *peers) Error {
	genString := infoMap["peers-generation"]
	if len(genString) == 0 {
		return newError(types.PARSE_ERROR, "peers-generation is empty")
	}

	gen, err := strconv.Atoi(genString)
	if err != nil {
		return newError(types.PARSE_ERROR, "peers-generation is not a number: "+genString)
	}

	peers.genChanged.Or(nd.peersGeneration.Get() != gen)
	return nil
}

func (nd *Node) verifyPartitionGeneration(infoMap map[string]string) Error {
	genString := infoMap["partition-generation"]

	if len(genString) == 0 {
		return newError(types.PARSE_ERROR, "partition-generation is empty")
	}

	gen, err := strconv.Atoi(genString)
	if err != nil {
		return newError(types.PARSE_ERROR, "partition-generation is not a number:"+genString)
	}

	if nd.partitionGeneration.Get() != gen {
		nd.partitionChanged.Set(true)
	}
	return nil
}

func (nd *Node) refreshPeers(peers *peers) {
	// Do not refresh peers when node connection has already failed during this cluster tend iteration.
	if nd.failures.Get() > 0 || !nd.active.Get() {
		return
	}

	peerParser, err := parsePeers(nd.cluster, nd)
	if err != nil {
		logger.Logger.Debug("Parsing peers failed: %s", err)
		nd.refreshFailed(err)
		return
	}

	peers.appendPeers(peerParser.peers)
	nd.peersGeneration.Set(int(peerParser.generation()))
	nd.peersCount.Set(len(peers.peers()))
	peers.refreshCount.IncrementAndGet()
}

func (nd *Node) refreshPartitions(peers *peers, partitions partitionMap, freshlyAdded bool) {
	// Do not refresh peers when node connection has already failed during this cluster tend iteration.
	// Also, avoid "split cluster" case where this node thinks it's a 1-node cluster.
	// Unchecked, such a node can dominate the partition map and cause all other
	// nodes to be dropped.
	if !freshlyAdded {
		if nd.failures.Get() > 0 || !nd.active.Get() || (nd.peersCount.Get() == 0 && peers.refreshCount.Get() > 1) {
			return
		}
	}

	parser, err := newPartitionParser(nd, partitions, _PARTITIONS)
	if err != nil {
		nd.refreshFailed(err)
		return
	}

	if parser.generation != nd.partitionGeneration.Get() {
		logger.Logger.Info("Node %s partition generation changed from %d to %d", nd.host.String(), nd.partitionGeneration.Get(), parser.getGeneration())
		nd.partitionChanged.Set(true)
		nd.partitionGeneration.Set(parser.getGeneration())
		atomic.AddInt64(&nd.stats.PartitionMapUpdates, 1)
	}
}

func (nd *Node) refreshFailed(e Error) {
	nd.peersGeneration.Set(-1)
	nd.partitionGeneration.Set(-1)

	if nd.cluster.clientPolicy.RackAware {
		nd.rebalanceGeneration.Set(-1)
	}

	nd.failures.IncrementAndGet()
	atomic.AddInt64(&nd.stats.TendsFailed, 1)

	// Only log message if cluster is still active.
	if nd.cluster.IsConnected() {
		logger.Logger.Warn("Node `%s` refresh failed: `%s`", nd, e)
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
func (nd *Node) GetConnection(timeout time.Duration) (conn *Connection, err Error) {
	if timeout <= 0 {
		timeout = _DEFAULT_TIMEOUT
	}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err = nd.getConnection(deadline, timeout)
		if err == nil && conn != nil {
			return conn, nil
		}

		if errors.Is(err, ErrServerNotAvailable) {
			return nil, err
		}

		time.Sleep(5 * time.Millisecond)
	}

	// in case the block didn't run at all
	if err == nil {
		err = ErrConnectionPoolEmpty.err()
	}

	return nil, err
}

// getConnection gets a connection to the node.
// If no pooled connection is available, a new connection will be created.
func (nd *Node) getConnection(deadline time.Time, timeout time.Duration) (conn *Connection, err Error) {
	return nd.getConnectionWithHint(deadline, timeout, 0)
}

// newConnectionAllowed will tentatively check if the client is allowed to make a new connection
// based on the ClientPolicy passed to it.
// This is more or less a copy of the logic in the beginning of newConnection function.
func (nd *Node) newConnectionAllowed() Error {
	if !nd.active.Get() {
		return ErrServerNotAvailable.err()
	}

	// if connection count is limited and enough connections are already created, don't create a new one
	cc := nd.connectionCount.IncrementAndGet()
	defer nd.connectionCount.DecrementAndGet()
	if nd.cluster.clientPolicy.LimitConnectionsToQueueSize && cc > nd.cluster.clientPolicy.ConnectionQueueSize {
		return ErrTooManyConnectionsForNode.err()
	}

	// Check for opening connection threshold
	if nd.cluster.clientPolicy.OpeningConnectionThreshold > 0 {
		ct := nd.cluster.connectionThreshold.IncrementAndGet()
		defer nd.cluster.connectionThreshold.DecrementAndGet()
		if ct > nd.cluster.clientPolicy.OpeningConnectionThreshold {
			return ErrTooManyOpeningConnections.err()
		}
	}

	return nil
}

// newConnection will make a new connection for the node.
func (nd *Node) newConnection(overrideThreshold bool) (*Connection, Error) {
	if !nd.active.Get() {
		return nil, ErrServerNotAvailable.err()
	}

	// if connection count is limited and enough connections are already created, don't create a new one
	cc := nd.connectionCount.IncrementAndGet()
	if nd.cluster.clientPolicy.LimitConnectionsToQueueSize && cc > nd.cluster.clientPolicy.ConnectionQueueSize {
		nd.connectionCount.DecrementAndGet()
		atomic.AddInt64(&nd.stats.ConnectionsPoolEmpty, 1)

		return nil, ErrTooManyConnectionsForNode.err()
	}

	// Check for opening connection threshold
	if !overrideThreshold && nd.cluster.clientPolicy.OpeningConnectionThreshold > 0 {
		ct := nd.cluster.connectionThreshold.IncrementAndGet()
		if ct > nd.cluster.clientPolicy.OpeningConnectionThreshold {
			nd.cluster.connectionThreshold.DecrementAndGet()
			nd.connectionCount.DecrementAndGet()

			return nil, ErrTooManyOpeningConnections.err()
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

	sessionInfo := nd.sessionInfo.Load().(*sessionInfo)
	// need to authenticate
	if err = conn.login(&nd.cluster.clientPolicy, nd.cluster.Password(), sessionInfo); err != nil {
		// increment node errors if authentication hit a network error
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
		logger.Logger.Debug("Error trying to make a connection to the node %s: %s", nd.String(), err.Error())
		return
	}

	nd.putConnectionWithHint(conn, hint)
}

// getConnectionWithHint gets a connection to the node.
// If no pooled connection is available, a new connection will be created.
func (nd *Node) getConnectionWithHint(deadline time.Time, timeout time.Duration, hint byte) (conn *Connection, err Error) {
	if !nd.active.Get() {
		return nil, ErrServerNotAvailable.err()
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
		// tentatively check if a connection is allowed to avoid launching too many goroutines.
		err = nd.newConnectionAllowed()
		if err == nil {
			go nd.makeConnectionForPool(hint)
		} else if errors.Is(err, ErrTooManyConnectionsForNode) {
			return nil, ErrConnectionPoolExhausted.err()
		}
		return nil, ErrConnectionPoolEmpty.err()
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
	if nd != nil {
		return nd.name + " " + nd.host.String()
	}
	return "<nil>"
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
func (nd *Node) MigrationInProgress() (bool, Error) {
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
func (nd *Node) WaitUntillMigrationIsFinished(timeout time.Duration) Error {
	if timeout <= 0 {
		timeout = _NO_TIMEOUT
	}
	done := make(chan Error)

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
		return newError(types.TIMEOUT)
	case err := <-done:
		return err
	}
}

// initTendConn sets up a connection to be used for info requests.
// The same connection will be used for tend.
func (nd *Node) initTendConn(timeout time.Duration) Error {
	if timeout <= 0 {
		timeout = _DEFAULT_TIMEOUT
	}
	deadline := time.Now().Add(timeout)

	if nd.tendConn == nil || !nd.tendConn.IsConnected() {
		var tendConn *Connection
		var err Error
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
func (nd *Node) requestInfoWithRetry(policy *InfoPolicy, n int, name ...string) (res map[string]string, err Error) {
	for i := 0; i < n; i++ {
		if res, err = nd.requestInfo(policy.Timeout, name...); err == nil {
			return res, nil
		}

		logger.Logger.Error("Error occurred while fetching info from the server node %s: %s", nd.host.String(), err.Error())
		time.Sleep(100 * time.Millisecond)
	}

	// return the last error
	return nil, err
}

// RequestInfo gets info values by name from the specified database server node.
func (nd *Node) RequestInfo(policy *InfoPolicy, name ...string) (map[string]string, Error) {
	return nd.requestInfo(policy.Timeout, name...)
}

// RequestInfo gets info values by name from the specified database server node.
func (nd *Node) requestInfo(timeout time.Duration, name ...string) (map[string]string, Error) {
	nd.tendConnLock.Lock()
	defer nd.tendConnLock.Unlock()

	if err := nd.initTendConn(timeout); err != nil {
		return nil, err
	}

	response, err := nd.tendConn.RequestInfo(name...)
	if err != nil {
		nd.tendConn.Close()
		return nil, err
	}
	return response, nil
}

// requestRawInfo gets info values by name from the specified database server node.
// It won't parse the results.
func (nd *Node) requestRawInfo(policy *InfoPolicy, name ...string) (*info, Error) {
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
func (nd *Node) RequestStats(policy *InfoPolicy) (map[string]string, Error) {
	infoMap, err := nd.RequestInfo(policy, "statistics")
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

// resetSessionInfo resets the sessionInfo after an
// unsuccessful authentication with token
func (nd *Node) resetSessionInfo() {
	si := &sessionInfo{}
	nd.sessionInfo.Store(si)
}

// sessionToken returns the session token for the node.
// It will return nil if the session has expired.
func (nd *Node) sessionToken() []byte {
	si := nd.sessionInfo.Load().(*sessionInfo)
	if !si.isValid() {
		return nil
	}

	return si.token
}

// Rack returns the rack number for the namespace.
func (nd *Node) Rack(namespace string) (int, Error) {
	racks := nd.racks.Load().(map[string]int)
	v, exists := racks[namespace]

	if exists {
		return v, nil
	}

	return -1, newCustomNodeError(nd, types.RACK_NOT_DEFINED)
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
func (nd *Node) WarmUp(count int) (int, Error) {
	var g errgroup.Group
	cnt := iatomic.NewInt(0)

	toAlloc := nd.connections.Cap() - nd.connectionCount.Get()
	if count < toAlloc && count > 0 {
		toAlloc = count
	}

	for i := 0; i < toAlloc; i++ {
		g.Go(func() error {
			conn, err := nd.newConnection(true)
			if err != nil {
				if errors.Is(err, ErrTooManyConnectionsForNode) {
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
	if err != nil {
		return cnt.Get(), err.(Error)
	}
	return cnt.Get(), nil
}

// fillMinCounts will fill the connection pool to the minimum required
// by the ClientPolicy.MinConnectionsPerNode
func (nd *Node) fillMinConns() (int, Error) {
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
func (nd *Node) validateErrorCount() Error {
	if !nd.errorCountWithinLimit() {
		return newError(types.MAX_ERROR_RATE)
	}
	return nil
}

// PeersGeneration returns node's Peers Generation
func (nd *Node) PeersGeneration() int {
	return nd.peersGeneration.Get()
}

// PartitionGeneration returns node's Partition Generation
func (nd *Node) PartitionGeneration() int {
	return nd.partitionGeneration.Get()
}

// RebalanceGeneration returns node's Rebalance Generation
func (nd *Node) RebalanceGeneration() int {
	return nd.rebalanceGeneration.Get()
}
