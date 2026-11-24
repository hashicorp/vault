// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package fairshare

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

// TestFairshare_newDispatcher verifies that a new dispatcher is created with
// the correct configuration and initial state.
func TestFairshare_newDispatcher(t *testing.T) {
	testCases := []struct {
		name               string
		numWorkers         int
		expectedNumWorkers int
	}{
		{
			name:               "",
			numWorkers:         0,
			expectedNumWorkers: 1,
		},
		{
			name:               "",
			numWorkers:         10,
			expectedNumWorkers: 10,
		},
		{
			name:               "test-dispatcher",
			numWorkers:         10,
			expectedNumWorkers: 10,
		},
	}

	l := newTestLogger("workerpool-test")
	for tcNum, tc := range testCases {
		d := newDispatcher(tc.name, tc.numWorkers, l)

		if tc.name != "" && d.name != tc.name {
			t.Errorf("tc %d: expected name %s, got %s", tcNum, tc.name, d.name)
		}
		if len(d.workers) != tc.expectedNumWorkers {
			t.Errorf("tc %d: expected %d workers, got %d", tcNum, tc.expectedNumWorkers, len(d.workers))
		}
		if d.jobCh == nil {
			t.Errorf("tc %d: work channel not set up properly", tcNum)
		}
	}
}

// TestFairshare_createDispatcher verifies that the dispatcher creation process
// properly initializes all required components and worker pools.
func TestFairshare_createDispatcher(t *testing.T) {
	testCases := []struct {
		name               string
		numWorkers         int
		expectedNumWorkers int
	}{
		{
			name:               "",
			numWorkers:         -1,
			expectedNumWorkers: 1,
		},
		{
			name:               "",
			numWorkers:         0,
			expectedNumWorkers: 1,
		},
		{
			name:               "",
			numWorkers:         10,
			expectedNumWorkers: 10,
		},
		{
			name:               "",
			numWorkers:         10,
			expectedNumWorkers: 10,
		},
		{
			name:               "test-dispatcher",
			numWorkers:         10,
			expectedNumWorkers: 10,
		},
	}

	l := newTestLogger("workerpool-test")
	for tcNum, tc := range testCases {
		d := createDispatcher(tc.name, tc.numWorkers, l)
		if d == nil {
			t.Fatalf("tc %d: expected non-nil object", tcNum)
		}

		if tc.name != "" && d.name != tc.name {
			t.Errorf("tc %d: expected name %s, got %s", tcNum, tc.name, d.name)
		}
		if len(d.name) == 0 {
			t.Errorf("tc %d: expected name to be set", tcNum)
		}
		if d.numWorkers != tc.expectedNumWorkers {
			t.Errorf("tc %d: expected %d workers, got %d", tcNum, tc.expectedNumWorkers, d.numWorkers)
		}
		if d.workers == nil {
			t.Errorf("tc %d: expected non-nil workers", tcNum)
		}
		if d.jobCh == nil {
			t.Errorf("tc %d: work channel not set up properly", tcNum)
		}
		if d.quit == nil {
			t.Errorf("tc %d: expected non-nil quit channel", tcNum)
		}
		if d.logger == nil {
			t.Errorf("tc %d: expected non-nil logger", tcNum)
		}
	}
}

// TestFairshare_initDispatcher verifies that dispatcher initialization sets up
// the worker pool and job channels correctly.
func TestFairshare_initDispatcher(t *testing.T) {
	testCases := []struct {
		numWorkers int
	}{
		{
			numWorkers: 1,
		},
		{
			numWorkers: 10,
		},
		{
			numWorkers: 100,
		},
		{
			numWorkers: 1000,
		},
	}

	l := newTestLogger("workerpool-test")
	for tcNum, tc := range testCases {
		d := createDispatcher("", tc.numWorkers, l)

		d.init()
		if len(d.workers) != tc.numWorkers {
			t.Fatalf("tc %d: expected %d workers, got %d", tcNum, tc.numWorkers, len(d.workers))
		}
	}
}

// TestFairshare_initializeWorker verifies that individual workers are properly
// initialized with the correct configuration and channels.
func TestFairshare_initializeWorker(t *testing.T) {
	numWorkers := 3

	d := createDispatcher("", numWorkers, newTestLogger("workerpool-test"))

	for workerNum := 0; workerNum < numWorkers; workerNum++ {
		d.initializeWorker()

		w := d.workers[workerNum]
		expectedName := fmt.Sprint("worker-", workerNum)
		if w.name != expectedName {
			t.Errorf("tc %d: expected name %s, got %s", workerNum, expectedName, w.name)
		}
		if w.jobCh != d.jobCh {
			t.Errorf("tc %d: work channel not set up properly", workerNum)
		}
		if w.quit == nil || w.quit != d.quit {
			t.Errorf("tc %d: quit channel not set up properly", workerNum)
		}
		if w.logger == nil || w.logger != d.logger {
			t.Errorf("tc %d: logger not set up properly", workerNum)
		}
	}
}

// TestFairshare_startWorker verifies that workers can be started and begin
// processing jobs from their assigned channels.
func TestFairshare_startWorker(t *testing.T) {
	d := newDispatcher("", 1, newTestLogger("workerpool-test"))

	d.workers[0].start()
	defer d.stop()

	var wg sync.WaitGroup
	ex := func(_ string) error {
		wg.Done()
		return nil
	}
	onFail := func(_ error) {}

	job := newTestJob(t, "test job", ex, onFail)

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)

	wg.Add(1)
	d.dispatch(&job, nil, nil)
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}
}

// TestFairshare_start verifies that the dispatcher can be started and all
// workers begin processing jobs correctly.
func TestFairshare_start(t *testing.T) {
	numJobs := 10
	var wg sync.WaitGroup
	ex := func(_ string) error {
		wg.Done()
		return nil
	}
	onFail := func(_ error) {}

	wg.Add(numJobs)
	d := newDispatcher("", 3, newTestLogger("workerpool-test"))

	d.start()
	defer d.stop()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	for i := 0; i < numJobs; i++ {
		job := newTestJob(t, fmt.Sprintf("job-%d", i), ex, onFail)
		d.dispatch(&job, nil, nil)
	}

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}
}

// TestFairshare_stop verifies that the dispatcher can be stopped cleanly,
// shutting down all workers and cleaning up resources.
func TestFairshare_stop(t *testing.T) {
	d := newDispatcher("", 5, newTestLogger("workerpool-test"))

	d.start()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)

	go func() {
		d.stop()
		d.wg.Wait()
		doneCh <- struct{}{}
	}()

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}
}

// TestFairshare_stopMultiple verifies that calling stop multiple times on a
// dispatcher is safe and doesn't cause panics or errors.
func TestFairshare_stopMultiple(t *testing.T) {
	d := newDispatcher("", 5, newTestLogger("workerpool-test"))

	d.start()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)

	go func() {
		d.stop()
		d.wg.Wait()
		doneCh <- struct{}{}
	}()

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}

	// essentially, we don't want to panic here
	var r interface{}
	go func() {
		t.Helper()

		defer func() {
			r = recover()
			doneCh <- struct{}{}
		}()

		d.stop()
		d.wg.Wait()
	}()

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}

	if r != nil {
		t.Fatalf("panic during second stop: %v", r)
	}
}

// TestFairshare_dispatch verifies that jobs are correctly dispatched to
// available workers and processed in the expected order.
func TestFairshare_dispatch(t *testing.T) {
	d := newDispatcher("", 1, newTestLogger("workerpool-test"))

	var wg sync.WaitGroup
	accumulatedIDs := make([]string, 0)
	ex := func(id string) error {
		accumulatedIDs = append(accumulatedIDs, id)
		wg.Done()
		return nil
	}
	onFail := func(_ error) {}

	expectedIDs := []string{"job-1", "job-2", "job-3", "job-4"}
	go func() {
		for _, id := range expectedIDs {
			job := newTestJob(t, id, ex, onFail)
			d.dispatch(&job, nil, nil)
		}
	}()

	wg.Add(len(expectedIDs))
	d.start()
	defer d.stop()

	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	timeout := time.After(5 * time.Second)
	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}

	if !reflect.DeepEqual(accumulatedIDs, expectedIDs) {
		t.Fatalf("bad job ids. expected %v, got %v", expectedIDs, accumulatedIDs)
	}
}

// TestFairshare_jobFailure verifies that the dispatcher handles job failures
// gracefully and continues processing subsequent jobs.
func TestFairshare_jobFailure(t *testing.T) {
	numJobs := 10
	testErr := fmt.Errorf("test error")
	var wg sync.WaitGroup

	ex := func(_ string) error {
		return testErr
	}
	onFail := func(err error) {
		if err != testErr {
			t.Errorf("got unexpected error. expected %v, got %v", testErr, err)
		}

		wg.Done()
	}

	wg.Add(numJobs)
	d := newDispatcher("", 3, newTestLogger("workerpool-test"))

	d.start()
	defer d.stop()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	for i := 0; i < numJobs; i++ {
		job := newTestJob(t, fmt.Sprintf("job-%d", i), ex, onFail)
		d.dispatch(&job, nil, nil)
	}

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}
}

// TestFairshare_nilLoggerDispatcher verifies that the dispatcher can operate
// correctly even when initialized with a nil logger.
func TestFairshare_nilLoggerDispatcher(t *testing.T) {
	d := newDispatcher("test-job-mgr", 1, nil)
	if d.logger == nil {
		t.Error("logger not set up properly")
	}
}
