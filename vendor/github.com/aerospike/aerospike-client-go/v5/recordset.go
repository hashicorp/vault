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
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/aerospike/aerospike-client-go/v5/internal/atomic"

	xornd "github.com/aerospike/aerospike-client-go/v5/types/rand"
)

// Result is the value returned by Recordset's Results() function.
type Result struct {
	Record *Record
	Err    Error
}

// String implements the Stringer interface
func (res *Result) String() string {
	if res.Record != nil {
		return fmt.Sprintf("%v", res.Record)
	}
	return fmt.Sprintf("%v", res.Err)
}

// Objectset encapsulates the result of Scan and Query commands.
type objectset struct {
	// a reference to the object channel to close on end signal
	objChan reflect.Value

	// errors is a channel on which all errors will be sent back.
	errors chan Error

	wgGoroutines sync.WaitGroup
	goroutines   *atomic.Int

	closed, active *atomic.Bool
	cancelled      chan struct{}

	chanLock sync.Mutex

	taskID uint64
}

// TaskId returns the transactionId/jobId sent to the server for this recordset.
func (os *objectset) TaskId() uint64 {
	os.chanLock.Lock()
	defer os.chanLock.Unlock()
	return os.taskID
}

// Always set the taskID client-side to a non-zero random value
func (os *objectset) resetTaskID() {
	os.chanLock.Lock()
	defer os.chanLock.Unlock()
	os.taskID = uint64(xornd.Int64())
}

// Recordset encapsulates the result of Scan and Query commands.
type Recordset struct {
	objectset

	// Records is a channel on which the resulting records will be sent back.
	// NOTE: Do not use Records directly. Range on channel returned by Results() instead.
	// Will be unexported in the future
	records chan *Result
}

// makes sure the recordset is closed eventually, even if it is not consumed
func recordsetFinalizer(rs *Recordset) {
	rs.Close()
}

// newObjectset generates a new RecordSet instance.
func newObjectset(objChan reflect.Value, goroutines int) *objectset {

	if objChan.Kind() != reflect.Chan ||
		objChan.Type().Elem().Kind() != reflect.Ptr ||
		objChan.Type().Elem().Elem().Kind() != reflect.Struct {
		panic("Scan/Query object channels should be of type `chan *T`")
	}

	rs := &objectset{
		objChan:    objChan,
		errors:     make(chan Error, goroutines),
		active:     atomic.NewBool(true),
		closed:     atomic.NewBool(false),
		goroutines: atomic.NewInt(goroutines),
		cancelled:  make(chan struct{}),
	}
	rs.wgGoroutines.Add(goroutines)
	rs.resetTaskID()
	return rs
}

// newRecordset generates a new RecordSet instance.
func newRecordset(recSize, goroutines int) *Recordset {
	var nilChan chan *struct{}

	rs := &Recordset{
		records:   make(chan *Result, recSize),
		objectset: *newObjectset(reflect.ValueOf(nilChan), goroutines),
	}

	runtime.SetFinalizer(rs, recordsetFinalizer)
	return rs
}

// IsActive returns true if the operation hasn't been finished or cancelled.
func (rcs *Recordset) IsActive() bool {
	return rcs.active.Get()
}

// Errors returns a read-only Error channel for the objectset. It will panic
// for recordsets returned for non-reflection APIs.
func (rcs *Recordset) Errors() <-chan Error {
	if rcs.records == nil {
		return (<-chan Error)(rcs.errors)
	}
	panic("Errors chan not valid for non-reflection API")
}

// Results returns a new receive-only channel with the results of the Scan/Query.
// This is a more idiomatic approach to the iterator pattern in getting the
// results back from the recordset, and doesn't require the user to write the
// ugly select in their code.
// Result contains a Record and an error reference.
//
// Example:
//
//  recordset, err := client.ScanAll(nil, namespace, set)
//  handleError(err)
//  for res := range recordset.Results() {
//    if res.Err != nil {
//      // handle error here
//    } else {
//      // process record here
//      fmt.Println(res.Record.Bins)
//    }
//  }
func (rcs *Recordset) Results() <-chan *Result {
	return (<-chan *Result)(rcs.records)
}

// Close all streams from different nodes. A successful close return nil,
// subsequent calls to the method will return ErrRecordsetClosed.err().
func (rcs *Recordset) Close() Error {
	// do it only once
	if !rcs.closed.CompareAndToggle(false) {
		return ErrRecordsetClosed.err()
	}

	// mark the recordset as inactive
	rcs.active.Set(false)

	close(rcs.cancelled)

	// wait till all goroutines are done, and signalEnd is called by the scan command
	rcs.wgGoroutines.Wait()

	return nil
}

func (rcs *Recordset) signalEnd() {
	rcs.wgGoroutines.Done()
	if rcs.goroutines.DecrementAndGet() == 0 {
		// mark the recordset as inactive
		rcs.active.Set(false)

		rcs.chanLock.Lock()
		defer rcs.chanLock.Unlock()

		if rcs.records != nil {
			close(rcs.records)
		} else if rcs.objChan.IsValid() {
			rcs.objChan.Close()
		}

		close(rcs.errors)
	}
}

func (rcs *Recordset) sendError(err Error) {
	rcs.chanLock.Lock()
	defer rcs.chanLock.Unlock()
	if rcs.IsActive() {
		if rcs.records != nil {
			rcs.records <- &Result{Err: err}
		} else {
			rcs.errors <- err
		}
	}
}
