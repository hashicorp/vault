package fairshare

import (
	"fmt"
	"io/ioutil"
	"sync"

	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

// Job is an interface for jobs used with this job manager
type Job interface {
	Execute() error
	OnFailure(err error)
}

// worker represents a single worker in a pool
type worker struct {
	name   string
	jobCh  <-chan Job
	quit   chan struct{}
	logger log.Logger

	// waitgroup for testing stop functionality
	wg *sync.WaitGroup
}

// start starts the worker listening and working until the quit channel is closed
func (w *worker) start() {
	w.wg.Add(1)

	go func() {
		for {
			select {
			case <-w.quit:
				w.wg.Done()
				return
			case job := <-w.jobCh:
				err := job.Execute()
				if err != nil {
					job.OnFailure(err)
				}
			}
		}
	}()
}

// dispatcher represents a worker pool
type dispatcher struct {
	name       string
	numWorkers int
	workers    []worker
	jobCh      chan Job
	onceStart  sync.Once
	onceStop   sync.Once
	quit       chan struct{}
	logger     log.Logger
	wg         *sync.WaitGroup
}

// newDispatcher generates a new worker dispatcher and populates it with workers
func newDispatcher(name string, numWorkers int, l log.Logger) *dispatcher {
	d := createDispatcher(name, numWorkers, l)

	d.init()
	return d
}

// dispatch dispatches a job to the worker pool
func (d *dispatcher) dispatch(job Job) {
	select {
	case d.jobCh <- job:
	case <-d.quit:
		d.logger.Info("shutting down during dispatch")
	}
}

// start starts all the workers listening on the job channel
// this will only start the workers for this dispatch once
func (d *dispatcher) start() {
	d.onceStart.Do(func() {
		d.logger.Trace("starting dispatcher")
		for _, w := range d.workers {
			worker := w
			worker.start()
		}
	})
}

// stop stops the worker pool asynchronously
func (d *dispatcher) stop() {
	d.onceStop.Do(func() {
		d.logger.Trace("terminating dispatcher")
		close(d.quit)
	})
}

// createDispatcher generates a new Dispatcher object, but does not initialize the
// worker pool
func createDispatcher(name string, numWorkers int, l log.Logger) *dispatcher {
	if l == nil {
		l = logging.NewVaultLoggerWithWriter(ioutil.Discard, log.NoLevel)
	}
	if numWorkers <= 0 {
		numWorkers = 1
		l.Warn("must have 1 or more workers. setting number of workers to 1")
	}

	if name == "" {
		guid, err := uuid.GenerateUUID()
		if err != nil {
			l.Warn("uuid generator failed, using 'no-uuid'", "err", err)
			guid = "no-uuid"
		}

		name = fmt.Sprintf("dispatcher-%s", guid)
	}

	var wg sync.WaitGroup
	d := dispatcher{
		name:       name,
		numWorkers: numWorkers,
		workers:    make([]worker, 0),
		jobCh:      make(chan Job),
		quit:       make(chan struct{}),
		logger:     l,
		wg:         &wg,
	}

	d.logger.Trace("created dispatcher", "name", d.name, "num_workers", d.numWorkers)
	return &d
}

func (d *dispatcher) init() {
	for len(d.workers) < d.numWorkers {
		d.initializeWorker()
	}

	d.logger.Trace("initialized dispatcher", "num_workers", d.numWorkers)
}

// initializeWorker initializes and adds a new worker, with an optional name
func (d *dispatcher) initializeWorker() {
	w := worker{
		name:   fmt.Sprint("worker-", len(d.workers)),
		jobCh:  d.jobCh,
		quit:   d.quit,
		logger: d.logger,
		wg:     d.wg,
	}

	d.workers = append(d.workers, w)
}
