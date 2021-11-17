package workerpool

import (
	"github.com/gammazero/deque"
	"sync"
	"time"
)

const (
	// This value is the size of the queue that workers register their
	// availability to the dispatcher.  There may be hundreds of workers, but
	// only a small channel is needed to register some of the workers.
	readyQueueSize = 16

	// If worker pool receives no new work for this period of time, then stop
	// a worker goroutine.
	idleTimeoutSec = 5
)

// New creates and starts a pool of worker goroutines.
//
// The maxWorkers parameter specifies the maximum number of workers that will
// execute tasks concurrently.  After each timeout period, a worker goroutine
// is stopped until there are no remaining workers.
func New(maxWorkers int) *WorkerPool {
	// There must be at least one worker.
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	pool := &WorkerPool{
		taskQueue:    make(chan func(), 1),
		maxWorkers:   maxWorkers,
		readyWorkers: make(chan chan func(), readyQueueSize),
		timeout:      time.Second * idleTimeoutSec,
		stoppedChan:  make(chan struct{}),
	}

	// Start the task dispatcher.
	go pool.dispatch()

	return pool
}

// WorkerPool is a collection of goroutines, where the number of concurrent
// goroutines processing requests does not exceed the specified maximum.
type WorkerPool struct {
	maxWorkers   int
	timeout      time.Duration
	taskQueue    chan func()
	readyWorkers chan chan func()
	stoppedChan  chan struct{}
	waitingQueue deque.Deque
	stopMutex    sync.Mutex
	stopped      bool
}

// Stop stops the worker pool and waits for only currently running tasks to
// complete.  Pending tasks that are not currently running are abandoned.
// Tasks must not be submitted to the worker pool after calling stop.
//
// Since creating the worker pool starts at least one goroutine, for the
// dispatcher, Stop() or StopWait() should be called when the worker pool is no
// longer needed.
func (p *WorkerPool) Stop() {
	p.stop(false)
}

// StopWait stops the worker pool and waits for all queued tasks tasks to
// complete.  No additional tasks may be submitted, but all pending tasks are
// executed by workers before this function returns.
func (p *WorkerPool) StopWait() {
	p.stop(true)
}

// Stopped returns true if this worker pool has been stopped.
func (p *WorkerPool) Stopped() bool {
	p.stopMutex.Lock()
	defer p.stopMutex.Unlock()
	return p.stopped
}

// Submit enqueues a function for a worker to execute.
//
// Any external values needed by the task function must be captured in a
// closure.  Any return values should be returned over a channel that is
// captured in the task function closure.
//
// Submit will not block regardless of the number of tasks submitted.  Each
// task is immediately given to an available worker or passed to a goroutine to
// be given to the next available worker.  If there are no available workers,
// the dispatcher adds a worker, until the maximum number of workers is
// running.
//
// After the maximum number of workers are running, and no workers are
// available, incoming tasks are put onto a queue and will be executed as
// workers become available.
//
// When no new tasks have been submitted for time period and a worker is
// available, the worker is shutdown.  As long as no new tasks arrive, one
// available worker is shutdown each time period until there are no more idle
// workers.  Since the time to start new goroutines is not significant, there
// is no need to retain idle workers.
func (p *WorkerPool) Submit(task func()) {
	if task != nil {
		p.taskQueue <- task
	}
}

// SubmitWait enqueues the given function and waits for it to be executed.
func (p *WorkerPool) SubmitWait(task func()) {
	if task == nil {
		return
	}
	doneChan := make(chan struct{})
	p.taskQueue <- func() {
		task()
		close(doneChan)
	}
	<-doneChan
}

// dispatch sends the next queued task to an available worker.
func (p *WorkerPool) dispatch() {
	defer close(p.stoppedChan)
	timeout := time.NewTimer(p.timeout)
	var (
		workerCount    int
		task           func()
		ok, wait       bool
		workerTaskChan chan func()
	)
	startReady := make(chan chan func())
Loop:
	for {
		// As long as tasks are in the waiting queue, remove and execute these
		// tasks as workers become available, and place new incoming tasks on
		// the queue.  Once the queue is empty, then go back to submitting
		// incoming tasks directly to available workers.
		if p.waitingQueue.Len() != 0 {
			select {
			case task, ok = <-p.taskQueue:
				if !ok {
					break Loop
				}
				if task == nil {
					wait = true
					break Loop
				}
				p.waitingQueue.PushBack(task)
			case workerTaskChan = <-p.readyWorkers:
				// A worker is ready, so give task to worker.
				workerTaskChan <- p.waitingQueue.PopFront().(func())
			}
			continue
		}
		timeout.Reset(p.timeout)
		select {
		case task, ok = <-p.taskQueue:
			if !ok || task == nil {
				break Loop
			}
			// Got a task to do.
			select {
			case workerTaskChan = <-p.readyWorkers:
				// A worker is ready, so give task to worker.
				workerTaskChan <- task
			default:
				// No workers ready.
				// Create a new worker, if not at max.
				if workerCount < p.maxWorkers {
					workerCount++
					go func(t func()) {
						startWorker(startReady, p.readyWorkers)
						// Submit the task when the new worker.
						taskChan := <-startReady
						taskChan <- t
					}(task)
				} else {
					// Enqueue task to be executed by next available worker.
					p.waitingQueue.PushBack(task)
				}
			}
		case <-timeout.C:
			// Timed out waiting for work to arrive.  Kill a ready worker.
			if workerCount > 0 {
				select {
				case workerTaskChan = <-p.readyWorkers:
					// A worker is ready, so kill.
					close(workerTaskChan)
					workerCount--
				default:
					// No work, but no ready workers.  All workers are busy.
				}
			}
		}
	}

	// If instructed to wait for all queued tasks, then remove from queue and
	// give to workers until queue is empty.
	if wait {
		for p.waitingQueue.Len() != 0 {
			workerTaskChan = <-p.readyWorkers
			// A worker is ready, so give task to worker.
			workerTaskChan <- p.waitingQueue.PopFront().(func())
		}
	}

	// Stop all remaining workers as they become ready.
	for workerCount > 0 {
		workerTaskChan = <-p.readyWorkers
		close(workerTaskChan)
		workerCount--
	}
}

// startWorker starts a goroutine that executes tasks given by the dispatcher.
//
// When a new worker starts, it registers its availability on the startReady
// channel.  This ensures that the goroutine associated with starting the
// worker gets to use the worker to execute its task.  Otherwise, the main
// dispatcher loop could steal the new worker and not know to start up another
// worker for the waiting goroutine.  The task would then have to wait for
// another existing worker to become available, even though capacity is
// available to start additional workers.
//
// A worker registers that is it available to do work by putting its task
// channel on the readyWorkers channel.  The dispatcher reads a worker's task
// channel from the readyWorkers channel, and writes a task to the worker over
// the worker's task channel.  To stop a worker, the dispatcher closes a
// worker's task channel, instead of writing a task to it.
func startWorker(startReady, readyWorkers chan chan func()) {
	go func() {
		taskChan := make(chan func())
		var task func()
		var ok bool
		// Register availability on starReady channel.
		startReady <- taskChan
		for {
			// Read task from dispatcher.
			task, ok = <-taskChan
			if !ok {
				// Dispatcher has told worker to stop.
				break
			}

			// Execute the task.
			task()

			// Register availability on readyWorkers channel.
			readyWorkers <- taskChan
		}
	}()
}

// stop tells the dispatcher to exit, and whether or not to complete queued
// tasks.
func (p *WorkerPool) stop(wait bool) {
	p.stopMutex.Lock()
	defer p.stopMutex.Unlock()
	if p.stopped {
		return
	}
	p.stopped = true
	if wait {
		p.taskQueue <- nil
	}
	// Close task queue and wait for currently running tasks to finish.
	close(p.taskQueue)
	<-p.stoppedChan
}

// WaitingQueueSize will return the size of the waiting queue
func (p *WorkerPool) WaitingQueueSize() int {
	return p.waitingQueue.Len()
}
