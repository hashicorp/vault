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
	"fmt"
	"reflect"
	"runtime"
	"sync"

	. "github.com/aerospike/aerospike-client-go/internal/atomic"
	. "github.com/aerospike/aerospike-client-go/types"
)

type Result struct {
	Record *Record
	Err    error
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

	// Errors is a channel on which all errors will be sent back.
	// NOTE: Do not use Errors directly. Range on channel returned by Results() instead.
	// This field is deprecated and will be unexported in the future
	Errors chan error

	wgGoroutines sync.WaitGroup
	goroutines   *AtomicInt

	closed, active *AtomicBool
	cancelled      chan struct{}

	chanLock sync.Mutex

	taskID uint64
}

// TaskId returns the transactionId/jobId sent to the server for this recordset.
func (os *objectset) TaskId() uint64 {
	return os.taskID
}

// Recordset encapsulates the result of Scan and Query commands.
type Recordset struct {
	objectset

	// Records is a channel on which the resulting records will be sent back.
	// NOTE: Do not use Records directly. Range on channel returned by Results() instead.
	// Will be unexported in the future
	Records chan *Record
}

// makes sure the recordset is closed eventually, even if it is not consumed
func recordsetFinalizer(rs *Recordset) {
	rs.Close()
}

// newObjectset generates a new RecordSet instance.
func newObjectset(objChan reflect.Value, goroutines int, taskID uint64) *objectset {

	if objChan.Kind() != reflect.Chan ||
		objChan.Type().Elem().Kind() != reflect.Ptr ||
		objChan.Type().Elem().Elem().Kind() != reflect.Struct {
		panic("Scan/Query object channels should be of type `chan *T`")
	}

	rs := &objectset{
		objChan:    objChan,
		Errors:     make(chan error, goroutines),
		active:     NewAtomicBool(true),
		closed:     NewAtomicBool(false),
		goroutines: NewAtomicInt(goroutines),
		cancelled:  make(chan struct{}),
		taskID:     taskID,
	}
	rs.wgGoroutines.Add(goroutines)

	return rs
}

// newRecordset generates a new RecordSet instance.
func newRecordset(recSize, goroutines int, taskID uint64) *Recordset {
	var nilChan chan *struct{}

	rs := &Recordset{
		Records:   make(chan *Record, recSize),
		objectset: *newObjectset(reflect.ValueOf(nilChan), goroutines, taskID),
	}

	runtime.SetFinalizer(rs, recordsetFinalizer)
	return rs
}

// IsActive returns true if the operation hasn't been finished or cancelled.
func (rcs *Recordset) IsActive() bool {
	return rcs.active.Get()
}

// Read reads the next record from the Recordset. If the Recordset has been
// closed, it returns ErrRecordsetClosed.
func (rcs *Recordset) Read() (record *Record, err error) {
	var ok bool

L:
	select {
	case record, ok = <-rcs.Records:
		if !ok {
			err = ErrRecordsetClosed
		}
	case err = <-rcs.Errors:
		if err == nil {
			// if err == nil, it means the Errors chan has been closed
			// we should not return nil as an error, so we should listen
			// to other chans again to determine either cancellation,
			// or normal EOR
			goto L
		}
	}

	return record, err
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
	recCap := cap(rcs.Records)
	if recCap < 1 {
		recCap = 1
	}
	res := make(chan *Result, recCap)

	select {
	case <-rcs.cancelled:
		// Bail early and give the caller a channel for nothing -- it's
		// functionally wasted memory, but the caller did something
		// after close, so it's their own doing.
		close(res)
		return res
	default:
	}

	go func(cancelled <-chan struct{}) {
		defer close(res)
		for {
			record, err := rcs.Read()
			if err == ErrRecordsetClosed {
				return
			}

			result := &Result{Record: record, Err: err}
			select {
			case <-cancelled:
				return
			case res <- result:

			}
		}
	}(rcs.cancelled)

	return res
}

// Close all streams from different nodes. A successful close return nil,
// subsequent calls to the method will return ErrRecordsetClosed.
func (rcs *Recordset) Close() error {
	// do it only once
	if !rcs.closed.CompareAndToggle(false) {
		return ErrRecordsetClosed
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

		if rcs.Records != nil {
			close(rcs.Records)
		} else if rcs.objChan.IsValid() {
			rcs.objChan.Close()
		}

		close(rcs.Errors)
	}
}

func (rcs *Recordset) sendError(err error) {
	rcs.chanLock.Lock()
	defer rcs.chanLock.Unlock()
	if rcs.IsActive() {
		rcs.Errors <- err
	}
}
