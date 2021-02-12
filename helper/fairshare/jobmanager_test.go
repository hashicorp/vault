package fairshare

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestJobManager_NewJobManager(t *testing.T) {
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
			numWorkers:         5,
			expectedNumWorkers: 5,
		},
		{
			name:               "",
			numWorkers:         5,
			expectedNumWorkers: 5,
		},
		{
			name:               "",
			numWorkers:         5,
			expectedNumWorkers: 5,
		},
		{
			name:               "",
			numWorkers:         5,
			expectedNumWorkers: 5,
		},
	}

	l := newTestLogger("jobmanager-test")
	for tcNum, tc := range testCases {
		j := NewJobManager(tc.name, tc.numWorkers, l)

		if tc.name != "" && tc.name != j.name {
			t.Errorf("tc %d: expected name %s, got %s", tcNum, tc.name, j.name)
		}
		if j.queues == nil {
			t.Errorf("tc %d: queues not set up properly", tcNum)
		}
		if j.queuesIndex == nil {
			t.Errorf("tc %d: queues index not set up properly", tcNum)
		}
		if j.quit == nil {
			t.Errorf("tc %d: quit channel not set up properly", tcNum)
		}
		if j.workerPool.numWorkers != tc.expectedNumWorkers {
			t.Errorf("tc %d: expected %d workers, got %d", tcNum, tc.expectedNumWorkers, j.workerPool.numWorkers)
		}
		if j.logger == nil {
			t.Errorf("tc %d: logger not set up properly", tcNum)
		}
	}
}

func TestJobManager_Start(t *testing.T) {
	numJobs := 10
	j := NewJobManager("job-mgr-test", 3, newTestLogger("jobmanager-test"))

	var wg sync.WaitGroup
	wg.Add(numJobs)
	j.Start()
	defer j.Stop()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	ex := func(_ string) error {
		wg.Done()
		return nil
	}
	onFail := func(_ error) {}

	for i := 0; i < numJobs; i++ {
		// distribute jobs between 3 queues in the job manager
		job := newTestJob(t, fmt.Sprintf("test-job-%d", i), ex, onFail)
		j.AddJob(&job, fmt.Sprintf("queue-%d", i%3))
	}

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}
}

func TestJobManager_StartAndPause(t *testing.T) {
	numJobs := 10
	j := NewJobManager("job-mgr-test", 3, newTestLogger("jobmanager-test"))

	var wg sync.WaitGroup
	wg.Add(numJobs)
	j.Start()
	defer j.Stop()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	ex := func(_ string) error {
		wg.Done()
		return nil
	}
	onFail := func(_ error) {}

	for i := 0; i < numJobs; i++ {
		// distribute jobs between 3 queues in the job manager
		job := newTestJob(t, fmt.Sprintf("test-job-%d", i), ex, onFail)
		j.AddJob(&job, fmt.Sprintf("queue-%d", i%3))
	}

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}

	// now that the work queue is empty, let's add more jobs and make sure
	// we pick up where we left off

	for i := 0; i < 5; i++ {
		numAdditionalJobs := 5
		wg.Add(numAdditionalJobs)

		timeout = time.After(5 * time.Second)
		go func() {
			wg.Wait()
			doneCh <- struct{}{}
		}()

		for i := numJobs; i < numJobs+numAdditionalJobs; i++ {
			// distribute jobs between 3 queues in the job manager
			job := newTestJob(t, fmt.Sprintf("test-job-%d", i), ex, onFail)
			j.AddJob(&job, fmt.Sprintf("queue-%d", i%3))
		}

		select {
		case <-doneCh:
			break
		case <-timeout:
			t.Fatal("timed out")
		}

		numJobs += numAdditionalJobs
	}
}

func TestJobManager_Stop(t *testing.T) {
	j := NewJobManager("job-mgr-test", 5, newTestLogger("jobmanager-test"))

	j.Start()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		j.Stop()
		doneCh <- struct{}{}
	}()

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}
}

func TestFairshare_StopMultiple(t *testing.T) {
	j := NewJobManager("job-mgr-test", 5, newTestLogger("jobmanager-test"))

	j.Start()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		j.Stop()
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

		j.Stop()
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

func TestJobManager_AddJob(t *testing.T) {
	testCases := []struct {
		name    string
		queueID string
	}{
		{
			name:    "test1",
			queueID: "q1",
		},
		{
			name:    "test2",
			queueID: "q1",
		},
		{
			name:    "test3",
			queueID: "q1",
		},
		{
			name:    "test4",
			queueID: "q2",
		},
		{
			name:    "test5",
			queueID: "q3",
		},
	}

	j := NewJobManager("job-mgr-test", 3, newTestLogger("jobmanager-test"))

	expectedCount := make(map[string]int)
	for _, tc := range testCases {
		if _, ok := expectedCount[tc.queueID]; !ok {
			expectedCount[tc.queueID] = 1
		} else {
			expectedCount[tc.queueID]++
		}

		job := newDefaultTestJob(t, tc.name)
		j.AddJob(&job, tc.queueID)
	}

	if len(expectedCount) != len(j.queues) {
		t.Fatalf("expected %d queues, got %d", len(expectedCount), len(j.queues))
	}

	for k, v := range j.queues {
		if v.Len() != expectedCount[k] {
			t.Fatalf("queue %s has bad count. expected %d, got %d", k, expectedCount[k], v.Len())
		}
	}
}

func TestJobManager_GetPendingJobCount(t *testing.T) {
	numJobs := 15
	j := NewJobManager("test-job-mgr", 3, newTestLogger("jobmanager-test"))

	for i := 0; i < numJobs; i++ {
		job := newDefaultTestJob(t, fmt.Sprintf("job-%d", i))
		j.AddJob(&job, fmt.Sprintf("queue-%d", i%4))
	}

	pendingJobs := j.GetPendingJobCount()
	if pendingJobs != numJobs {
		t.Errorf("expected %d jobs, got %d", numJobs, pendingJobs)
	}
}

func TestJobManager_GetWorkQueueLengths(t *testing.T) {
	j := NewJobManager("test-job-mgr", 3, newTestLogger("jobmanager-test"))

	expected := make(map[string]int)
	for i := 0; i < 25; i++ {
		queueID := fmt.Sprintf("queue-%d", i%4)
		job := newDefaultTestJob(t, fmt.Sprintf("job-%d", i))

		j.AddJob(&job, queueID)

		if _, ok := expected[queueID]; !ok {
			expected[queueID] = 1
		} else {
			expected[queueID]++
		}
	}

	pendingJobs := j.GetWorkQueueLengths()
	if !reflect.DeepEqual(pendingJobs, expected) {
		t.Errorf("expected %v job count, got %v", expected, pendingJobs)
	}
}

func TestJobManager_removeLastQueueAccessed(t *testing.T) {
	j := NewJobManager("job-mgr-test", 1, newTestLogger("jobmanager-test"))

	j.addQueue("queue-0")
	j.addQueue("queue-1")
	j.addQueue("queue-2")

	testCases := []struct {
		lastQueueAccessed        int
		updatedLastQueueAccessed int
		len                      int
		expectedQueues           []string
	}{
		{
			// remove with bad index (too low)
			lastQueueAccessed:        -1,
			updatedLastQueueAccessed: -1,
			len:                      3,
			expectedQueues:           []string{"queue-0", "queue-1", "queue-2"},
		},
		{
			// remove with bad index (too high)
			lastQueueAccessed:        3,
			updatedLastQueueAccessed: 3,
			len:                      3,
			expectedQueues:           []string{"queue-0", "queue-1", "queue-2"},
		},
		{
			// remove queue-1 (index 1)
			lastQueueAccessed:        1,
			updatedLastQueueAccessed: 0,
			len:                      2,
			expectedQueues:           []string{"queue-0", "queue-2"},
		},
		{
			// remove queue-0 (index 0)
			lastQueueAccessed:        0,
			updatedLastQueueAccessed: 0,
			len:                      1,
			expectedQueues:           []string{"queue-2"},
		},
		{
			// remove queue-1 (index 1)
			lastQueueAccessed:        0,
			updatedLastQueueAccessed: -1,
			len:                      0,
			expectedQueues:           []string{},
		},
	}

	j.l.Lock()
	defer j.l.Unlock()
	for _, tc := range testCases {
		j.lastQueueAccessed = tc.lastQueueAccessed
		j.removeLastQueueAccessed()

		if j.lastQueueAccessed != tc.updatedLastQueueAccessed {
			t.Errorf("last queue access update failed. expected %d, got %d", tc.updatedLastQueueAccessed, j.lastQueueAccessed)
		}
		if len(j.queuesIndex) != tc.len {
			t.Fatalf("queue index update failed. expected %d elements, found %v", tc.len, j.queues)
		}
		if len(j.queues) != len(tc.expectedQueues) {
			t.Fatalf("bad amount of queues. expected %d, found %v", len(tc.expectedQueues), j.queues)
		}

		for _, q := range tc.expectedQueues {
			if _, ok := j.queues[q]; !ok {
				t.Errorf("bad queue. expected %s in %v", q, j.queues)
			}
		}
	}
}

func TestJobManager_getNextJob(t *testing.T) {
	testCases := []struct {
		name    string
		queueID string
	}{
		{
			name:    "job-1",
			queueID: "queue-1",
		},
		{
			name:    "job-2",
			queueID: "queue-2",
		},
		{
			name:    "job-3",
			queueID: "queue-1",
		},
		{
			name:    "job-4",
			queueID: "queue-3",
		},
		{
			name:    "job-5",
			queueID: "queue-3",
		},
	}

	// we add the jobs before starting the workers, so we'd expect the round
	// robin to pick the least-recently-added job from each queue, and cycle
	// through queues in a round-robin fashion. jobs would appear on the queues
	// as illustrated below, and we expect to round robin as:
	// queue-1 -> queue-2 -> queue-3 -> queue-1 ...
	//
	// queue-1 [job-3, job-1]
	// queue-2 [job-2]
	// queue-3 [job-5, job-4]

	// ... where jobs are pushed to the left side and popped from the right side

	expectedOrder := []string{"job-1", "job-2", "job-4", "job-3", "job-5"}

	// use one worker to guarantee ordering
	numWorkers := 1
	resultsCh := make(chan string)
	defer close(resultsCh)

	j := NewJobManager("test-job-mgr", numWorkers, newTestLogger("jobmanager-test"))

	doneCh := make(chan struct{})
	var mu sync.Mutex
	order := make([]string, 0)
	timeout := time.After(5 * time.Second)

	go func() {
		for {
			select {
			case res, ok := <-resultsCh:
				if !ok {
					return
				}

				mu.Lock()
				order = append(order, res)
				mu.Unlock()
			}
		}
	}()

	var wg sync.WaitGroup
	ex := func(name string) error {
		resultsCh <- name
		time.Sleep(50 * time.Millisecond)
		wg.Done()
		return nil
	}
	onFail := func(_ error) {}

	for _, tc := range testCases {
		wg.Add(1)
		job := newTestJob(t, tc.name, ex, onFail)
		j.AddJob(&job, tc.queueID)
	}

	j.Start()
	defer j.Stop()

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

	mu.Lock()
	defer mu.Unlock()
	if !reflect.DeepEqual(order, expectedOrder) {
		t.Fatalf("results out of order. \nexpected: %v\ngot: %v", expectedOrder, order)
	}
}

func TestFairshare_nilLoggerJobManager(t *testing.T) {
	j := NewJobManager("test-job-mgr", 1, nil)
	if j.logger == nil {
		t.Error("logger not set up properly")
	}
}
