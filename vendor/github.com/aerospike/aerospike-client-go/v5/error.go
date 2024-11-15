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
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/aerospike/aerospike-client-go/v5/types"
)

// Error is the internal error interface for the Aerospike client's errors.
// All the public API return this error type. This interface is compatible
// with error interface, including errors.Is and errors.As.
type Error interface {
	error

	// Matches will return true is the ResultCode of the error or
	// any of the errors wrapped down the chain has any of the
	// provided codes.
	Matches(rcs ...types.ResultCode) bool

	// resultCode returns the error result code.
	resultCode() types.ResultCode

	// Unwrap returns the error inside
	Unwrap() error

	// Trace returns a stack trace of where the error originates from
	Trace() string

	iter(int) Error
	setInDoubt(bool, int) Error
	setNode(*Node) Error
	markInDoubt() Error
	wrap(error) Error
}

// AerospikeError implements Error interface for aerospike specific errors.
// All errors returning from the library are of this type.
// Errors resulting from Go's stdlib are not translated to this type, unless
// they are a net.Timeout error. Refer to errors_test.go for examples.
// To be able to check for error type, you could use the idiomatic
// errors.Is and errors.As patterns:
//   if errors.Is(err, as.ErrTimeout) {
//       ...
//   }
// or
//   if errors.Is(err, &as.AerospikeError{ResultCode: ast.PARAMETER_ERROR}) {
//       ...
//   }
// or
//   if err.Matches(ast.TIMEOUT, ast.NETWORK_ERROR, ast.PARAMETER_ERROR) {
//       ...
//   }
// or
//   ae := &as.AerospikeError{}
//   if errors.As(err, &ae) {
//       println(ae.ResultCode)
//   }
type AerospikeError struct {
	wrapped error

	// error message
	msg string

	// Node where the error occurred
	Node *Node

	// ResultCode determines the type of error
	ResultCode types.ResultCode

	// InDoubt determines if the command was sent to the server, but
	// there is doubt if the server received and executed the command
	// and changed the data. Only applies to commands that change data
	InDoubt bool

	// Iteration determies on which retry the error occurred
	Iteration int

	// Includes stack frames for the error
	stackFrames []stackFrame
}

var _ error = &AerospikeError{}
var _ Error = &AerospikeError{}

// newError generates a new AerospikeError instance.
// If no message is provided, the result code will be translated into the default
// error message automatically.
func newError(code types.ResultCode, messages ...string) Error {
	if len(messages) == 0 {
		messages = []string{types.ResultCodeToString(code)}
	}

	return &AerospikeError{msg: strings.Join(messages, " "), ResultCode: code, stackFrames: stackTrace(nil)}
}

func newErrorAndWrap(e error, code types.ResultCode, messages ...string) Error {
	ne := newError(code, messages...)
	ne.wrap(e)
	return ne
}

func newTimeoutError(e error, messages ...string) Error {
	ne := newError(types.TIMEOUT, messages...)
	ne.wrap(e)
	return ne
}

func newCommonError(e error, messages ...string) Error {
	ne := newError(types.COMMON_ERROR, messages...)
	ne.wrap(e)
	return ne
}

// SetInDoubt sets whether it is possible that the write transaction may have completed
// even though this error was generated.  This may be the case when a
// client error occurs (like timeout) after the command was sent to the server.
func (ase *AerospikeError) setInDoubt(isRead bool, commandSentCounter int) Error {
	if !isRead && (commandSentCounter > 1 || (commandSentCounter == 1 && (ase.ResultCode == types.TIMEOUT || ase.ResultCode <= 0))) {
		ase.InDoubt = true
	}
	return ase
}

func (ase *AerospikeError) setNode(node *Node) Error {
	ase.Node = node
	return ase
}

func (ase *AerospikeError) markInDoubt() Error {
	ase.InDoubt = true
	return ase
}

func (ase *AerospikeError) resultCode() types.ResultCode {
	return ase.ResultCode
}

// Trace returns a stack trace of where the error originates from
func (ase *AerospikeError) Trace() string {
	var sb strings.Builder
	for i := range ase.stackFrames {
		sb.WriteString(ase.stackFrames[i].String())
		sb.WriteString("\n")
	}

	if ase.wrapped != nil {
		ae := new(AerospikeError)
		if errors.As(ase.wrapped, &ae) {
			sb.WriteString("Embedded:\n")
			sb.WriteString(ae.Trace())
		}
	}

	return sb.String()
}

// Error implements the error interface
func (ase *AerospikeError) Error() string {
	const cErr = "ResultCode: %s, Iteration: %d, InDoubt: %t, Node: %s: %s"
	const cErrNL = cErr + "\n%s"
	if ase.wrapped != nil {
		return fmt.Sprintf(cErrNL, ase.ResultCode.String(), ase.Iteration, ase.InDoubt, ase.Node, ase.msg, ase.wrapped.Error())
	}
	return fmt.Sprintf(cErr, ase.ResultCode.String(), ase.Iteration, ase.InDoubt, ase.Node, ase.msg)
}

func (ase *AerospikeError) wrap(err error) Error {
	ase.wrapped = err
	return ase
}

func (ase *AerospikeError) iter(i int) Error {
	if ase == nil {
		return nil
	}
	ase.Iteration = i
	return ase
}

// Matches returns true if the error or any of its wrapped errors contains
// any of the passed results codes.
// For convenience, it will return false if the error is nil.
func (ase *AerospikeError) Matches(rcs ...types.ResultCode) bool {
	// don't panic on nil error, and don't go ahead
	// if no result codes are provided
	if ase == nil || len(rcs) == 0 {
		return false
	}

	for i := range rcs {
		if ase.ResultCode == rcs[i] {
			return true
		}
	}

	ae := &AerospikeError{}
	if ase.wrapped != nil && errors.As(ase.wrapped, &ae) {
		return ae.Matches(rcs...)
	}

	return false
}

// As implements the interface for errors.As function.
func (ase *AerospikeError) As(target interface{}) bool {
	ae, ok := target.(*AerospikeError)
	if !ok {
		return false
	}

	ae.wrapped = ase.wrapped
	ae.msg = ase.msg
	ae.ResultCode = ase.ResultCode
	ae.InDoubt = ase.InDoubt
	ae.Node = ase.Node
	return true
}

// Is compares an error with the AerospikeError.
// If the error is not of type *AerospikeError, it will return false.
// Otherwise, it will compare ResultCode and Node (if it exists), and
// will return a result accordingly.
// If passed error's InDoubt is set to true, the InDoubt property will
// also be checked. You should not check if the error's InDoubt is false, since
// it is not checked when the passed error's InDoubt is false.
func (ase *AerospikeError) Is(e error) bool {
	if ase == nil || e == nil {
		return false
	}

	var target *AerospikeError

	switch t := e.(type) {
	case *AerospikeError:
		target = t
	case *constAerospikeError:
		target = &t.AerospikeError
	default:
		return false
	}

	res := (ase.ResultCode == target.ResultCode) &&
		(ase.Node == target.Node || target.Node == nil)

	if target.InDoubt {
		res = res && (ase.InDoubt == target.InDoubt)
	}

	return res
}

// Unwrap will return the error wrapped inside the error, or nil.
func (ase *AerospikeError) Unwrap() error {
	return ase.wrapped
}

/*
	Node Error
*/

func newNodeError(node *Node, err Error) Error {
	if err == nil {
		return nil
	}

	ae := new(AerospikeError)
	errors.As(err, &ae)

	res := *ae
	res.Node = node
	res.wrap(err)
	return &res
}

func newCustomNodeError(node *Node, code types.ResultCode, messages ...string) Error {
	ne := newError(code, messages...)
	ne.setNode(node)
	return ne
}

func newWrapNetworkError(err error, messages ...string) Error {
	ne := newError(types.NETWORK_ERROR, messages...)
	ne.wrap(err)
	return ne
}

func newInvalidNodeError(clusterSize int, partition *Partition) Error {
	// important to check for clusterSize first, since partition may be nil sometimes
	if clusterSize == 0 {
		return ErrClusterIsEmpty.err()
	}
	res := newError(types.INVALID_NODE_ERROR, "Node not found for partition "+partition.String()+" in partition table.")
	res.wrap(nil)
	return res
}

/*
	constAerospikeError
*/

var _ Error = newError(0)

// constAerospikeError makes sure that constant errors are not chained and invalidated.
// By having a new type, the compiler will enforce the constants.
type constAerospikeError struct {
	AerospikeError
}

func newConstError(code types.ResultCode, messages ...string) *constAerospikeError {
	if len(messages) == 0 {
		messages = []string{types.ResultCodeToString(code)}
	}

	return &constAerospikeError{AerospikeError{msg: strings.Join(messages, " "), ResultCode: code}}
}

func (ase *constAerospikeError) err() Error {
	v := ase.AerospikeError
	v.wrap(nil)
	return &v
}

//revive:disable

var (
	ErrServerNotAvailable              = newConstError(types.SERVER_NOT_AVAILABLE)
	ErrInvalidPartitionMap             = newConstError(types.INVALID_CLUSTER_PARTITION_MAP, "Partition map errors normally occur when the cluster has partitioned due to network anomaly or node crash, or is not configured properly. Refer to https://www.aerospike.com/docs/operations/configure for more information.")
	ErrKeyNotFound                     = newConstError(types.KEY_NOT_FOUND_ERROR)
	ErrRecordsetClosed                 = newConstError(types.RECORDSET_CLOSED)
	ErrConnectionPoolEmpty             = newConstError(types.NO_AVAILABLE_CONNECTIONS_TO_NODE, "connection pool is empty. This happens when no connections were available")
	ErrConnectionPoolExhausted         = newConstError(types.NO_AVAILABLE_CONNECTIONS_TO_NODE, "Connection pool is exhausted. This happens when all connection are in-use already, and opening more connections is not allowed due to the limits set in policy.ConnectionQueueSize and policy.LimitConnectionsToQueueSize")
	ErrTooManyConnectionsForNode       = newConstError(types.NO_AVAILABLE_CONNECTIONS_TO_NODE, "connection limit reached for this node. This value is controlled via ClientPolicy.LimitConnectionsToQueueSize")
	ErrTooManyOpeningConnections       = newConstError(types.NO_AVAILABLE_CONNECTIONS_TO_NODE, "too many connections are trying to open at once. This value is controlled via ClientPolicy.OpeningConnectionThreshold")
	ErrTimeout                         = newConstError(types.TIMEOUT, "command execution timed out on client: See `Policy.Timeout`")
	ErrNetTimeout                      = newConstError(types.TIMEOUT, "network timeout")
	ErrUDFBadResponse                  = newConstError(types.UDF_BAD_RESPONSE, "invalid UDF return value")
	ErrNoOperationsSpecified           = newConstError(types.INVALID_COMMAND, "no operations were passed to QueryExecute")
	ErrNoBinNamesAllowedInQueryExecute = newConstError(types.INVALID_COMMAND, "`Statement.BinNames` must be empty for QueryExecute")
	ErrFilteredOut                     = newConstError(types.FILTERED_OUT)
	ErrPartitionScanQueryNotSupported  = newConstError(types.PARAMETER_ERROR, "partition Scans/Queries are not supported by all nodes in this cluster")
	ErrScanTerminated                  = newConstError(types.SCAN_TERMINATED)
	ErrQueryTerminated                 = newConstError(types.QUERY_TERMINATED)
	ErrClusterIsEmpty                  = newConstError(types.INVALID_NODE_ERROR, "cluster is empty")
	ErrInvalidUser                     = newConstError(types.INVALID_USER)
	ErrNotAuthenticated                = newConstError(types.NOT_AUTHENTICATED)
	ErrNetwork                         = newConstError(types.NOT_AUTHENTICATED)
	ErrInvalidObjectType               = newConstError(types.SERIALIZE_ERROR, "invalid type for result object. It should be of type struct pointer or addressable")
	ErrMaxRetriesExceeded              = newConstError(types.MAX_RETRIES_EXCEEDED, "command execution timed out on client: Exceeded number of retries. See `Policy.MaxRetries`.")
	ErrInvalidParam                    = newConstError(types.PARAMETER_ERROR)
	ErrLuaPoolEmpty                    = newConstError(types.COMMON_ERROR, "Error fetching a lua instance from pool")
)

//revive:enable

// chainErrors wraps an error inside a new error. The new (outer) error cannot be nil.
// if the old error is nil, the new error will be returned.
func chainErrors(outer Error, inner error) Error {
	if inner == nil && outer == nil {
		return nil
	} else if inner == nil {
		return outer
	} else if outer == nil {
		if e, ok := inner.(Error); ok {
			return e
		}
		return newCommonError(inner)
	}

	var ae *AerospikeError
	switch outer.(type) {
	case *constAerospikeError:
		t := outer.(*constAerospikeError).AerospikeError
		ae = &t
	case *AerospikeError:
		// copy the reference to avoid issues with checking the last error
		// when it is chained.
		t := *outer.(*AerospikeError)
		ae = &t
	}

	if inner == nil {
		return ae
	}

	ae.wrapped = inner
	return ae
}

type stackFrame struct {
	fl, fn string
	ln     int
}

func (st *stackFrame) String() string {
	return st.fl + ":" + strconv.Itoa(st.ln) + " " + st.fn + "()"
}

func stackTrace(err Error) []stackFrame {
	const maxDepth = 10
	sFrames := make([]stackFrame, 0, maxDepth)
	for i := 3; i <= maxDepth+3; i++ {
		pc, fl, ln, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		sFrame := stackFrame{
			fl: fl,
			fn: fn.Name(),
			ln: ln,
		}
		sFrames = append(sFrames, sFrame)
	}

	if len(sFrames) > 0 {
		return sFrames
	}
	return nil
}
