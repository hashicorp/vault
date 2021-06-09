package fairshare

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

type JobManager struct {
	name   string
	queues map[string]*list.List

	quit    chan struct{}
	newWork chan struct{} // must be buffered

	workerPool  *dispatcher
	workerCount map[string]int

	onceStart sync.Once
	onceStop  sync.Once

	logger log.Logger

	totalJobs  int
	metricSink *metricsutil.ClusterMetricSink

	// waitgroup for testing stop functionality
	wg sync.WaitGroup

	// protects `queues`, `workerCount`
	l sync.RWMutex
}

// NewJobManager creates a job manager, with an optional name
func NewJobManager(name string, numWorkers int, l log.Logger, metricSink *metricsutil.ClusterMetricSink) *JobManager {
	if l == nil {
		l = logging.NewVaultLoggerWithWriter(ioutil.Discard, log.NoLevel)
	}
	if name == "" {
		guid, err := uuid.GenerateUUID()
		if err != nil {
			l.Warn("uuid generator failed, using 'no-uuid'", "err", err)
			guid = "no-uuid"
		}

		name = fmt.Sprintf("jobmanager-%s", guid)
	}

	wp := newDispatcher(fmt.Sprintf("%s-dispatcher", name), numWorkers, l)

	j := JobManager{
		name:        name,
		queues:      make(map[string]*list.List),
		quit:        make(chan struct{}),
		newWork:     make(chan struct{}, 1),
		workerPool:  wp,
		workerCount: make(map[string]int),
		logger:      l,
		metricSink:  metricSink,
	}

	j.logger.Trace("created job manager", "name", name, "pool_size", numWorkers)
	return &j
}

// Start starts the job manager
// note: a given job manager cannot be restarted after it has been stopped
func (j *JobManager) Start() {
	j.onceStart.Do(func() {
		j.logger.Trace("starting job manager", "name", j.name)
		j.workerPool.start()
		j.assignWork()
	})
}

// Stop stops the job manager asynchronously
func (j *JobManager) Stop() {
	j.onceStop.Do(func() {
		j.logger.Trace("terminating job manager...")
		close(j.quit)
		j.workerPool.stop()
	})
}

// AddJob adds a job to the given queue, creating the queue if it doesn't exist
func (j *JobManager) AddJob(job Job, queueID string) {
	j.l.Lock()
	if len(j.queues) == 0 {
		defer func() {
			// newWork must be buffered to avoid deadlocks if work is added
			// before the job manager is started
			j.newWork <- struct{}{}
		}()
	}
	defer j.l.Unlock()

	if _, ok := j.queues[queueID]; !ok {
		j.addQueue(queueID)
	}

	j.queues[queueID].PushBack(job)
	j.totalJobs++

	if j.metricSink != nil {
		j.metricSink.AddSampleWithLabels([]string{j.name, "job_manager", "queue_length"}, float32(j.queues[queueID].Len()), []metrics.Label{{"queue_id", queueID}})
		j.metricSink.AddSample([]string{j.name, "job_manager", "total_jobs"}, float32(j.totalJobs))
	}
}

// GetCurrentJobCount returns the total number of pending jobs in the job manager
func (j *JobManager) GetPendingJobCount() int {
	j.l.RLock()
	defer j.l.RUnlock()

	cnt := 0
	for _, q := range j.queues {
		cnt += q.Len()
	}

	return cnt
}

// GetWorkerCounts() returns a map of queue ID to number of active workers
func (j *JobManager) GetWorkerCounts() map[string]int {
	j.l.RLock()
	defer j.l.RUnlock()
	return j.workerCount
}

// GetWorkQueueLengths() returns a map of queue ID to number of jobs in the queue
func (j *JobManager) GetWorkQueueLengths() map[string]int {
	out := make(map[string]int)

	j.l.RLock()
	defer j.l.RUnlock()

	for k, v := range j.queues {
		out[k] = v.Len()
	}

	return out
}

// getNextJob pops the next job to be processed and prunes empty queues
// it also returns the ID of the queue the job is associated with
func (j *JobManager) getNextJob() (Job, string) {
	j.l.Lock()
	defer j.l.Unlock()

	if len(j.queues) == 0 {
		return nil, ""
	}

	queueID, canAssignWorker := j.getNextQueue()
	if !canAssignWorker {
		return nil, ""
	}

	jobElement := j.queues[queueID].Front()
	out := j.queues[queueID].Remove(jobElement)

	j.totalJobs--

	if j.metricSink != nil {
		j.metricSink.AddSampleWithLabels([]string{j.name, "job_manager", "queue_length"}, float32(j.queues[queueID].Len()), []metrics.Label{{"queue_id", queueID}})
		j.metricSink.AddSample([]string{j.name, "job_manager", "total_jobs"}, float32(j.totalJobs))
	}

	if j.queues[queueID].Len() == 0 {
		// we remove the empty queue, but we don't remove the worker count
		// in case we are still working on previous jobs from this queue.
		// worker count cleanup is handled in j.decrementWorkerCount
		delete(j.queues, queueID)
	}

	return out.(Job), queueID
}

// returns the next queue to assign work from, and a bool if there is a queue
// that can have a worker assigned.
// the intent is to avoid over-allocating work from specific queues, as
// outlined in RFC VLT-145
// note: this must be called with j.l held
func (j *JobManager) getNextQueue() (string, bool) {
	var nextQueue string
	var canAssignWorker bool

	queueIDsByIncreasingWorkers := j.sortByNumWorkers()
	for _, queueID := range queueIDsByIncreasingWorkers {
		if !j.queueWorkersSaturated(queueID) {
			nextQueue = queueID
			canAssignWorker = true
			break
		}
	}

	return nextQueue, canAssignWorker
}

// returns true if there are already too many workers on this queue
// note: this must be called with j.l held (at least for read)
// down the road we may want to factor in queue length relative to num queues
func (j *JobManager) queueWorkersSaturated(queueID string) bool {
	numActiveQueues := float64(len(j.queues))
	numTotalWorkers := float64(j.workerPool.numWorkers)
	maxWorkersPerQueue := math.Ceil(0.9 * numTotalWorkers / numActiveQueues)

	numWorkersPerQueue := j.workerCount

	return numWorkersPerQueue[queueID] >= int(maxWorkersPerQueue)
}

// sortByNumWorkers returns queueIDs in order of increasing number of workers
// note: this must be called with j.l held
func (j *JobManager) sortByNumWorkers() []string {
	out := make([]string, 0, len(j.queues))
	for queueID := range j.queues {
		out = append(out, queueID)
	}

	workersPerQueue := j.workerCount

	sort.Slice(out, func(i, j int) bool {
		// TODO do we want this explicity, or do we want some randomness?
		// I think it's fine since it only breaks ties between number of workers,
		// and it makes it easier to test
		if workersPerQueue[out[i]] == workersPerQueue[out[j]] {
			return out[i] < out[j]
		}

		return workersPerQueue[out[i]] < workersPerQueue[out[j]]
	})

	return out
}

// increment the worker count for this queue
func (j *JobManager) incrementWorkerCount(queueID string) {
	j.l.Lock()
	defer j.l.Unlock()

	j.workerCount[queueID]++
}

// decrement the worker count for this queue
// this also removes worker tracking for this queue if needed
func (j *JobManager) decrementWorkerCount(queueID string) {
	j.l.Lock()
	defer j.l.Unlock()

	j.workerCount[queueID]--

	_, queueEmpty := j.queues[queueID]
	if queueEmpty && j.workerCount[queueID] < 1 {
		delete(j.workerCount, queueID)
	}
}

// assignWork continually loops checks for new jobs and dispatches them to the
// worker pool
func (j *JobManager) assignWork() {
	j.wg.Add(1)

	go func() {
		for {
			for {
				// assign work while there are jobs to distribute
				select {
				case <-j.quit:
					j.wg.Done()
					return
				case <-j.newWork:
					// keep the channel empty since we're already processing work
				default:
				}

				job, queueID := j.getNextJob()
				if job != nil {
					j.workerPool.dispatch(job,
						func() {
							j.incrementWorkerCount(queueID)
						},
						func() {
							j.decrementWorkerCount(queueID)
						})
				} else {
					break
				}
			}

			select {
			case <-j.quit:
				j.wg.Done()
				return
			case <-j.newWork:
				// listen for wake-up when an emtpy job manager has been given work
			case <-time.After(50 * time.Millisecond):
				// periodically check if new workers can be assigned. with the
				// fairsharing worker distribution it can be the case that there
				// is work waiting, but no queues are elligible for another worker
			}
		}
	}()
}

// addQueue generates a new queue if a queue for `queueID` doesn't exist
// it also starts tracking workers on that queue, if not already tracked
// note: this must be called with l held for write
func (j *JobManager) addQueue(queueID string) {
	if _, ok := j.queues[queueID]; !ok {
		j.queues[queueID] = list.New()
	}

	// it's possible the queue ran out of work and was pruned, but there were
	// still workers operating on data formerly in that queue, which were still
	// being tracked. if that is the case, we don't want to wipe out that worker
	// count when the queue is re-initialized.
	if _, ok := j.workerCount[queueID]; !ok {
		j.workerCount[queueID] = 0
	}
}
