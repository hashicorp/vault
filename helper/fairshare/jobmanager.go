package fairshare

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

/*
Future Work:
- track workers per queue. this will involve things like:
	- somehow wrap the Execute/OnFailure functions to increment counter when
		they start running, and decrement when they stop running
		-- put a queue.IncrementCounter() call at the beginning
		-- call the provided work function in the middle
		-- put a queue.DecrementCounter() call at the end
	- job has a queueID or reference to the queue
- queue only removed when empty AND no workers
*/

type JobManager struct {
	name              string
	queues            map[string]*list.List
	queuesIndex       []string
	lastQueueAccessed int
	quit              chan struct{}
	newWork           chan struct{} // must be buffered
	workerPool        *dispatcher
	onceStart         sync.Once
	onceStop          sync.Once
	logger            log.Logger
	wg                sync.WaitGroup
	totalJobs         int

	// protects `queues`, `queuesIndex`, `lastQueueAccessed`
	l sync.RWMutex
}

// NewJobManager creates a job manager, with an optional name
func NewJobManager(name string, numWorkers int, l log.Logger) *JobManager {
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
		logger:            l,
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

// Stop stops the job manager, and waits for the worker pool and
// job manager to quit gracefully
func (j *JobManager) Stop() {
	j.onceStop.Do(func() {
		j.logger.Trace("terminating job manager and waiting...")
		j.workerPool.stop()
		close(j.quit)
		j.wg.Wait()
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
	metrics.AddSampleWithLabels([]string{j.name, "job_manager", "queue_length"}, float32(j.queues[queueID].Len()), []metrics.Label{{"queue_id", queueID}})
	metrics.SetGauge([]string{j.name, "job_manager", "total_jobs"}, float32(j.totalJobs))
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
	// TODO implement with VLT-145
	return nil
}

// GetWorkQueueLengths() returns a map of queue ID to number of active workers
func (j *JobManager) GetWorkQueueLengths() map[string]int {
	out := make(map[string]int)

	j.l.RLock()
	defer j.l.RUnlock()

	for k, v := range j.queues {
		out[k] = v.Len()
	}

	return out
}

// getNextJob grabs the next job to be processed and prunes empty queues
func (j *JobManager) getNextJob() Job {
	j.l.Lock()
	defer j.l.Unlock()

	if len(j.queues) == 0 {
		return nil
	}

	j.lastQueueAccessed = (j.lastQueueAccessed + 1) % len(j.queuesIndex)
	queueID := j.queuesIndex[j.lastQueueAccessed]

	jobElement := j.queues[queueID].Front()
	out := j.queues[queueID].Remove(jobElement)

	j.totalJobs--
	metrics.AddSampleWithLabels([]string{j.name, "job_manager", "queue_length"}, float32(j.queues[queueID].Len()), []metrics.Label{{"queue_id", queueID}})
	metrics.SetGauge([]string{j.name, "job_manager", "total_jobs"}, float32(j.totalJobs))

	if j.queues[queueID].Len() == 0 {
		j.removeLastQueueAccessed()
	}

	return out.(Job)
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

				job := j.getNextJob()
				if job != nil {
					j.workerPool.dispatch(job)
				} else {
					break
				}
			}

			// listen for a wake-up when an emtpy job manager has been given
			// new work
			select {
			case <-j.quit:
				j.wg.Done()
				return
			case <-j.newWork:
				break
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
