// +build !as_performance

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
	"errors"
	"reflect"

	. "github.com/aerospike/aerospike-client-go/internal/atomic"
	. "github.com/aerospike/aerospike-client-go/types"
	xornd "github.com/aerospike/aerospike-client-go/types/rand"
)

// PutObject writes record bin(s) to the server.
// The policy specifies the transaction timeout, record expiration and how the transaction is
// handled when the record already exists.
// If the policy is nil, the default relevant policy will be used.
// A struct can be tagged to influence the way the object is put in the database:
//
// type Person struct {
//		TTL uint32 `asm:"ttl"`
//		RecGen uint32 `asm:"gen"`
//		Name string `as:"name"`
// 		Address string `as:"desc,omitempty"`
// 		Age uint8 `as:",omitempty"`
// 		Password string `as:"-"`
// }
//
// Tag `as:` denotes Aerospike fields. The first value will be the alias for the field.
// `,omitempty` (without any spaces between the comma and the word) will act like the
// json package, and will not send the value of the field to the database if the value is zero value.
// Tag `asm:` denotes Aerospike Meta fields, and includes ttl and generation values.
// If a tag is marked with `-`, it will not be sent to the database at all.
// Note: Tag `as` can be replaced with any other user-defined tag via the function `SetAerospikeTag`.
func (clnt *Client) PutObject(policy *WritePolicy, key *Key, obj interface{}) (err error) {
	policy = clnt.getUsableWritePolicy(policy)

	binMap := marshal(obj, clnt.cluster.supportsFloat.Get())
	command, err := newWriteCommand(clnt.cluster, policy, key, nil, binMap, _WRITE)
	if err != nil {
		return err
	}

	res := command.Execute()
	return res
}

// GetObject reads a record for specified key and puts the result into the provided object.
// The policy can be used to specify timeouts.
// If the policy is nil, the default relevant policy will be used.
func (clnt *Client) GetObject(policy *BasePolicy, key *Key, obj interface{}) error {
	policy = clnt.getUsablePolicy(policy)

	rval := reflect.ValueOf(obj)
	binNames := objectMappings.getFields(rval.Type())

	command, err := newReadCommand(clnt.cluster, policy, key, binNames, nil)
	if err != nil {
		return err
	}

	command.object = &rval
	return command.Execute()
}

// BatchGetObject reads multiple record headers and bins for specified keys in one batch request.
// The returned objects are in positional order with the original key array order.
// If a key is not found, the positional object will not change, and the positional found boolean will be false.
// The policy can be used to specify timeouts.
// If the policy is nil, the default relevant policy will be used.
func (clnt *Client) BatchGetObjects(policy *BatchPolicy, keys []*Key, objects []interface{}) (found []bool, err error) {
	policy = clnt.getUsableBatchPolicy(policy)

	// check the size of  key and objects
	if len(keys) != len(objects) {
		return nil, errors.New("wrong number of arguments to BatchGetObjects: number of keys and objects do not match")
	}

	if len(keys) == 0 {
		return nil, errors.New("wrong number of arguments to BatchGetObjects: keys are empty")
	}

	binSet := map[string]struct{}{}
	objectsVal := make([]*reflect.Value, len(objects))
	for i := range objects {
		rval := reflect.ValueOf(objects[i])
		objectsVal[i] = &rval
		for _, bn := range objectMappings.getFields(rval.Type()) {
			binSet[bn] = struct{}{}
		}
	}
	binNames := make([]string, 0, len(binSet))
	for binName := range binSet {
		binNames = append(binNames, binName)
	}

	objectsFound := make([]bool, len(keys))
	cmd := newBatchCommandGet(nil, nil, policy, keys, binNames, nil, _INFO1_READ)
	cmd.objects = objectsVal
	cmd.objectsFound = objectsFound

	batchNodes, err := newBatchNodeList(clnt.cluster, policy, keys)
	if err != nil {
		return nil, err
	}

	err, filteredOut := clnt.batchExecute(policy, batchNodes, cmd)
	if err != nil {
		return nil, err
	}

	if filteredOut > 0 {
		err = ErrFilteredOut
	}

	return objectsFound, err
}

// ScanAllObjects reads all records in specified namespace and set from all nodes.
// If the policy's concurrentNodes is specified, each server node will be read in
// parallel. Otherwise, server nodes are read sequentially.
// If the policy is nil, the default relevant policy will be used.
func (clnt *Client) ScanAllObjects(apolicy *ScanPolicy, objChan interface{}, namespace string, setName string, binNames ...string) (*Recordset, error) {
	policy := *clnt.getUsableScanPolicy(apolicy)

	nodes := clnt.cluster.GetNodes()
	if len(nodes) == 0 {
		return nil, NewAerospikeError(SERVER_NOT_AVAILABLE, "Scan failed because cluster is empty.")
	}

	clusterKey := int64(0)
	if policy.FailOnClusterChange {
		var err error
		clusterKey, err = queryValidateBegin(nodes[0], namespace)
		if err != nil {
			return nil, err
		}
	}

	first := NewAtomicBool(true)

	// result recordset
	taskID := uint64(xornd.Int64())
	res := &Recordset{
		objectset: *newObjectset(reflect.ValueOf(objChan), len(nodes), taskID),
	}

	// the whole call should be wrapped in a goroutine
	if policy.ConcurrentNodes {
		for _, node := range nodes {
			go func(node *Node, first bool) {
				// Errors are handled inside the command itself
				clnt.scanNodeObjects(&policy, node, res, namespace, setName, taskID, clusterKey, first, binNames...)
			}(node, first.CompareAndToggle(true))
		}
	} else {
		// scan nodes one by one
		go func() {
			for _, node := range nodes {
				// Errors are handled inside the command itself
				clnt.scanNodeObjects(&policy, node, res, namespace, setName, taskID, clusterKey, first.CompareAndToggle(true), binNames...)
			}
		}()
	}

	return res, nil
}

// scanNodeObjects reads all records in specified namespace and set for one node only,
// and marshalls the results into the objects of the provided channel in Recordset.
// If the policy is nil, the default relevant policy will be used.
// The resulting records will be marshalled into the objChan.
// objChan will be closed after all the records are read.
func (clnt *Client) ScanNodeObjects(apolicy *ScanPolicy, node *Node, objChan interface{}, namespace string, setName string, binNames ...string) (*Recordset, error) {
	policy := *clnt.getUsableScanPolicy(apolicy)

	clusterKey := int64(0)
	if policy.FailOnClusterChange {
		var err error
		clusterKey, err = queryValidateBegin(node, namespace)
		if err != nil {
			return nil, err
		}
	}

	// results channel must be async for performance
	taskID := uint64(xornd.Int64())
	res := &Recordset{
		objectset: *newObjectset(reflect.ValueOf(objChan), 1, taskID),
	}

	go clnt.scanNodeObjects(&policy, node, res, namespace, setName, taskID, clusterKey, true, binNames...)
	return res, nil
}

// scanNodeObjects reads all records in specified namespace and set for one node only,
// and marshalls the results into the objects of the provided channel in Recordset.
// If the policy is nil, the default relevant policy will be used.
func (clnt *Client) scanNodeObjects(policy *ScanPolicy, node *Node, recordset *Recordset, namespace string, setName string, taskID uint64, clusterKey int64, first bool, binNames ...string) error {
	command := newScanObjectsCommand(node, policy, namespace, setName, binNames, recordset, taskID, clusterKey, first)
	return command.Execute()
}

// QueryNodeObjects executes a query on all nodes in the cluster and marshals the records into the given channel.
// The query executor puts records on the channel from separate goroutines.
// The caller can concurrently pop objects.
//
// This method is only supported by Aerospike 3+ servers.
// If the policy is nil, the default relevant policy will be used.
func (clnt *Client) QueryObjects(policy *QueryPolicy, statement *Statement, objChan interface{}) (*Recordset, error) {
	policy = clnt.getUsableQueryPolicy(policy)

	nodes := clnt.cluster.GetNodes()
	if len(nodes) == 0 {
		return nil, NewAerospikeError(SERVER_NOT_AVAILABLE, "Query failed because cluster is empty.")
	}

	clusterKey := int64(0)
	if policy.FailOnClusterChange {
		var err error
		clusterKey, err = queryValidateBegin(nodes[0], statement.Namespace)
		if err != nil {
			return nil, err
		}
	}

	first := NewAtomicBool(true)

	// results channel must be async for performance
	recSet := &Recordset{
		objectset: *newObjectset(reflect.ValueOf(objChan), len(nodes), statement.TaskId),
	}

	// the whole call sho
	// results channel must be async for performance
	for _, node := range nodes {
		// copy policies to avoid race conditions
		newPolicy := *policy
		command := newQueryObjectsCommand(node, &newPolicy, statement, recSet, clusterKey, first.CompareAndToggle(true))
		go func() {
			// Do not send the error to the channel; it is already handled in the Execute method
			command.Execute()
		}()
	}

	return recSet, nil
}

// QueryNodeObjects executes a query on a specific node and marshals the records into the given channel.
// The caller can concurrently pop records off the channel.
//
// This method is only supported by Aerospike 3+ servers.
// If the policy is nil, the default relevant policy will be used.
func (clnt *Client) QueryNodeObjects(policy *QueryPolicy, node *Node, statement *Statement, objChan interface{}) (*Recordset, error) {
	policy = clnt.getUsableQueryPolicy(policy)

	// results channel must be async for performance
	recSet := &Recordset{
		objectset: *newObjectset(reflect.ValueOf(objChan), 1, statement.TaskId),
	}

	clusterKey := int64(0)
	if policy.FailOnClusterChange {
		var err error
		clusterKey, err = queryValidateBegin(node, statement.Namespace)
		if err != nil {
			return nil, err
		}
	}

	// copy policies to avoid race conditions
	newPolicy := *policy
	command := newQueryRecordCommand(node, &newPolicy, statement, recSet, clusterKey, true)
	go func() {
		command.Execute()
	}()

	return recSet, nil
}
