package fairshare

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"sync"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

type JobManager struct {
	// TODO check if queuesIndex is still needed
	name              string
	queues            map[string]*list.List
	queuesIndex       []string
	lastQueueAccessed int

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

	// protects `queues`, `workerCount`, `queuesIndex`, `lastQueueAccessed`
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
		name:              name,
		queues:            make(map[string]*list.List),
		queuesIndex:       make([]string, 0),
		lastQueueAccessed: -1,
		quit:              make(chan struct{}),
		newWork:           make(chan struct{}, 1),
		workerPool:        wp,
		workerCount:       make(map[string]int),
		logger:            l,
		metricSink:        metricSink,
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
		j.workerCount[queueID] = 0
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

	queueID, queueIdx, _ := j.getNextQueueLegacy()
	j.lastQueueAccessed = queueIdx

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

		// worker count cleanup is handled in *JobManager.decrementWorkerCount

		j.removeLastQueueAccessed()
	}

	return out.(Job), queueID
}

// returns the next queue to assign work from, and a bool if there is work to be assigned
// note: this must be called with j.l held
// TODO remove during legacy strategy cleanup
func (j *JobManager) getNextQueueLegacy() (string, int, bool) {
	queueIdx := (j.lastQueueAccessed + 1) % len(j.queuesIndex)
	return j.queuesIndex[queueIdx], queueIdx, true
}

// returns the next queue to assign work from, and a bool if there is a queue that
// can have a worker assigned.
// the intent is to avoid over-allocating work from specific queues, as
// outlined in RFC VLT-145
// TODO update doc
// note: this must be called with j.l held
func (j *JobManager) getNextQueueFairshare() (string, int, bool) {
	var nextQueue string
	var haveWork bool

	queueIDsByIncreasingWorkers := j.sortByNumWorkers()
	for _, queueID := range queueIDsByIncreasingWorkers {
		if j.queues[queueID].Len() < 1 {
			// TODO this shouldn't happen as we prune empty queues - verify and remove
			continue
		}

		if !j.queueWorkersSaturated(queueID) {
			nextQueue = queueID
			haveWork = true
			break
		}
	}

	// TODO update this 0 return - we don't have the queue index here
	return nextQueue, 0, haveWork
}

// returns true if there are already too many workers on this queue
// note: this must be called with j.l held (at least for read)
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
	out := make([]string, len(j.queues))
	copy(out, j.queuesIndex)

	workersPerQueue := j.workerCount

	sort.Slice(out, func(i, j int) bool {
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
					j.workerPool.dispatch(job, func() {
						j.incrementWorkerCount(queueID)
					}, func() {
						j.decrementWorkerCount(queueID)
					})
				} else {
					break
				}
			}

			// listen for wake-up when an emtpy job manager has been given work
			select {
			case <-j.quit:
				j.wg.Done()
				return
			case <-j.newWork:
			}
		}
	}()
}

// addQueue generates a new queue if a queue for `queueID` doesn't exist
// note: this must be called with l held for write
func (j *JobManager) addQueue(queueID string) {
	if _, ok := j.queues[queueID]; !ok {
		j.queues[queueID] = list.New()
		j.queuesIndex = append(j.queuesIndex, queueID)
	}
}

// removeLastQueueAccessed removes the queue and index map for the last queue
// accessed. It is to be used when the last queue accessed has emptied.
// note: this must be called with l held for write
func (j *JobManager) removeLastQueueAccessed() {
	if j.lastQueueAccessed == -1 || j.lastQueueAccessed > len(j.queuesIndex)-1 {
		j.logger.Warn("call to remove queue out of bounds", "idx", j.lastQueueAccessed)
		return
	}

	queueID := j.queuesIndex[j.lastQueueAccessed]

	// remove the queue
	delete(j.queues, queueID)

	// remove the index for the queue
	j.queuesIndex = append(j.queuesIndex[:j.lastQueueAccessed], j.queuesIndex[j.lastQueueAccessed+1:]...)

	// correct the last queue accessed for round robining
	if j.lastQueueAccessed > 0 {
		j.lastQueueAccessed--
	} else {
		j.lastQueueAccessed = len(j.queuesIndex) - 1
	}
}
