package gocbcore

import (
	"errors"
	"fmt"
	"sync"
)

var (
	errPipelineClosed = errors.New("pipeline has been closed")
	errPipelineFull   = errors.New("pipeline is too full")
)

type memdGetClientFn func() (*memdClient, error)

type memdPipeline struct {
	address     string
	getClientFn memdGetClientFn
	maxItems    int
	queue       *memdOpQueue
	maxClients  int
	clients     []*memdPipelineClient
	clientsLock sync.Mutex
}

func newPipeline(address string, maxClients, maxItems int, getClientFn memdGetClientFn) *memdPipeline {
	return &memdPipeline{
		address:     address,
		getClientFn: getClientFn,
		maxClients:  maxClients,
		maxItems:    maxItems,
		queue:       newMemdOpQueue(),
	}
}

func newDeadPipeline(maxItems int) *memdPipeline {
	return newPipeline("", 0, maxItems, nil)
}

// nolint: unused
func (pipeline *memdPipeline) debugString() string {
	var outStr string

	if pipeline.address != "" {
		outStr += fmt.Sprintf("Address: %s\n", pipeline.address)
		outStr += fmt.Sprintf("Max Clients: %d\n", pipeline.maxClients)
		outStr += fmt.Sprintf("Num Clients: %d\n", len(pipeline.clients))
		outStr += fmt.Sprintf("Max Items: %d\n", pipeline.maxItems)
	} else {
		outStr += "Dead-Server Queue\n"
	}

	outStr += "Op Queue:\n"
	outStr += reindentLog("  ", pipeline.queue.debugString())

	return outStr
}

func (pipeline *memdPipeline) Clients() []*memdPipelineClient {
	pipeline.clientsLock.Lock()
	defer pipeline.clientsLock.Unlock()
	return pipeline.clients
}

func (pipeline *memdPipeline) Address() string {
	return pipeline.address
}

func (pipeline *memdPipeline) StartClients() {
	pipeline.clientsLock.Lock()
	defer pipeline.clientsLock.Unlock()

	for len(pipeline.clients) < pipeline.maxClients {
		client := newMemdPipelineClient(pipeline)
		pipeline.clients = append(pipeline.clients, client)

		go client.Run()
	}
}

func (pipeline *memdPipeline) sendRequest(req *memdQRequest, maxItems int) error {
	err := pipeline.queue.Push(req, maxItems)
	if err == errOpQueueClosed {
		return errPipelineClosed
	} else if err == errOpQueueFull {
		return errPipelineFull
	} else if err != nil {
		return err
	}

	return nil
}

func (pipeline *memdPipeline) RequeueRequest(req *memdQRequest) error {
	return pipeline.sendRequest(req, 0)
}

func (pipeline *memdPipeline) SendRequest(req *memdQRequest) error {
	return pipeline.sendRequest(req, pipeline.maxItems)
}

// Performs a takeover of another pipeline.  Note that this does not
//  take over the requests queued in the old pipeline, and those must
//  be drained and processed separately.
func (pipeline *memdPipeline) Takeover(oldPipeline *memdPipeline) {
	if oldPipeline.address != pipeline.address {
		logErrorf("Attempted pipeline takeover for differing address")

		// We try to 'gracefully' error here by resolving all the requests as
		//  errors, but allowing the application to continue.
		err := oldPipeline.Close()
		if err != nil {
			// Log and continue with this non-fatal error.
			logDebugf("Failed to shutdown old pipeline (%s)", err)
		}

		// Drain all the requests as an internal error so they are not lost
		oldPipeline.Drain(func(req *memdQRequest) {
			req.tryCallback(nil, errCliInternalError)
		})

		return
	}

	// Migrate all the clients to the new pipeline
	oldPipeline.clientsLock.Lock()
	clients := oldPipeline.clients
	oldPipeline.clients = nil
	oldPipeline.clientsLock.Unlock()

	pipeline.clientsLock.Lock()
	pipeline.clients = clients
	for _, client := range pipeline.clients {
		client.ReassignTo(pipeline)
	}
	pipeline.clientsLock.Unlock()

	// Shut down the old pipelines queue, this will force all the
	//  clients to 'refresh' their consumer, and pick up the new
	//  pipeline queue from the new pipeline.  This will also block
	//  any writers from sending new requests here if they have an
	//  out of date route config.
	oldPipeline.queue.Close()
}

func (pipeline *memdPipeline) Close() error {
	// Shut down all the clients
	pipeline.clientsLock.Lock()
	clients := pipeline.clients
	pipeline.clients = nil
	pipeline.clientsLock.Unlock()

	hadErrors := false
	for _, pipecli := range clients {
		err := pipecli.Close()
		if err != nil {
			logErrorf("failed to shutdown pipeline client: %s", err)
			hadErrors = true
		}
	}

	// Kill the queue, forcing everyone to stop
	pipeline.queue.Close()

	if hadErrors {
		return errCliInternalError
	}

	return nil
}

func (pipeline *memdPipeline) Drain(cb func(*memdQRequest)) {
	pipeline.queue.Drain(cb)
}
