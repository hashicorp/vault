package gocbcore

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

type memdPipelineClient struct {
	parent    *memdPipeline
	address   string
	client    *memdClient
	consumer  *memdOpConsumer
	lock      sync.Mutex
	closedSig chan struct{}
	state     uint32

	connectError error
}

func newMemdPipelineClient(parent *memdPipeline) *memdPipelineClient {
	return &memdPipelineClient{
		parent:    parent,
		address:   parent.address,
		closedSig: make(chan struct{}),
		state:     uint32(EndpointStateDisconnected),
	}
}

func (pipecli *memdPipelineClient) State() EndpointState {
	return EndpointState(atomic.LoadUint32(&pipecli.state))
}

func (pipecli *memdPipelineClient) Error() error {
	pipecli.lock.Lock()
	defer pipecli.lock.Unlock()
	return pipecli.connectError
}

func (pipecli *memdPipelineClient) ReassignTo(parent *memdPipeline) {
	pipecli.lock.Lock()
	pipecli.parent = parent
	oldConsumer := pipecli.consumer
	pipecli.consumer = nil
	pipecli.lock.Unlock()

	if oldConsumer != nil {
		oldConsumer.Close()
	}
}

func (pipecli *memdPipelineClient) ioLoop(client *memdClient) {
	pipecli.lock.Lock()
	if pipecli.parent == nil {
		logDebugf("Pipeline client ioLoop started with no parent pipeline")
		pipecli.lock.Unlock()

		err := client.Close()
		if err != nil {
			logErrorf("Failed to close client for shut down ioLoop (%s)", err)
		}

		return
	}

	pipecli.client = client
	pipecli.lock.Unlock()

	killSig := make(chan struct{})

	// This goroutine is responsible for monitoring the client and handling
	// the cleanup whenever it shuts down.  All cases of the client being
	// shut down flow through this goroutine, even cases where we may already
	// be aware that the client is shutdown, outside this scope.
	go func() {
		logDebugf("Pipeline client `%s/%p` client watcher starting...", pipecli.address, pipecli)

		<-client.CloseNotify()

		logDebugf("Pipeline client `%s/%p` client died", pipecli.address, pipecli)

		pipecli.lock.Lock()
		pipecli.client = nil
		activeConsumer := pipecli.consumer
		pipecli.consumer = nil
		pipecli.lock.Unlock()

		logDebugf("Pipeline client `%s/%p` closing consumer %p", pipecli.address, pipecli, activeConsumer)

		// If we have a consumer, we need to close it to signal the loop below that
		// something has happened.  If there is no consumer, we don't need to signal
		// as the loop below will already be in the process of fetching a new one,
		// where it will inevitably detect the problem.
		if activeConsumer != nil {
			activeConsumer.Close()
		}

		killSig <- struct{}{}
	}()

	logDebugf("Pipeline client `%s/%p` IO loop starting...", pipecli.address, pipecli)

	var localConsumer *memdOpConsumer
	for {
		if localConsumer == nil {
			logDebugf("Pipeline client `%s/%p` fetching new consumer", pipecli.address, pipecli)

			pipecli.lock.Lock()

			if pipecli.consumer != nil {
				// If we still have an active consumer, lets close it to make room for the new one
				pipecli.consumer.Close()
				pipecli.consumer = nil
			}

			if pipecli.client == nil {
				// The client has disconnected from the server, this only occurs AFTER the watcher
				// goroutine running above has detected the client is closed and has cleaned it up.
				pipecli.lock.Unlock()
				break
			}

			if pipecli.parent == nil {
				// This pipelineClient has been shut down
				logDebugf("Pipeline client `%s/%p` found no parent pipeline", pipecli.address, pipecli)
				pipecli.lock.Unlock()

				// Close our client to force the watcher goroutine above to clean it up
				err := client.Close()
				if err != nil {
					logErrorf("Pipeline client `%s/%p` failed to shut down client socket (%s)", pipecli.address, pipecli, err)
				}

				break
			}

			// Fetch a new consumer to use for this iteration
			localConsumer = pipecli.parent.queue.Consumer()
			pipecli.consumer = localConsumer

			pipecli.lock.Unlock()
		}

		req := localConsumer.Pop()
		if req == nil {
			// Set the local consumer to null, this will force our normal logic to run
			// which will clean up the original consumer and then attempt to acquire a
			// new one if we are not being cleaned up.  This is a minor code-optimization
			// to avoid having to do a lock/unlock just to lock above anyways.  It does
			// have the downside of not being able to detect where we've looped around
			// in error though.
			localConsumer = nil
			continue
		}

		err := client.SendRequest(req)
		if err != nil {
			logDebugf("Pipeline client `%s/%p` encountered a socket write error: %v", pipecli.address, pipecli, err)

			if !errors.Is(err, io.EOF) {
				// If we errored the write, and the client was not already closed,
				// lets go ahead and close it.  This will trigger the shutdown
				// logic via the client watcher above.  If the socket error was EOF
				// we already did shut down, and the watcher should already be
				// cleaning up.
				err := client.Close()
				if err != nil {
					logErrorf("Pipeline client `%s/%p` failed to shut down errored client socket (%s)", pipecli.address, pipecli, err)
				}
			}

			// Send this request upwards to be processed by the higher level processor
			shortCircuited, routeErr := client.postErrHandler(nil, req, err)
			if !shortCircuited {
				client.CancelRequest(req, err)
				req.tryCallback(nil, routeErr)
				break
			}

			// Stop looping
			break
		}
	}

	atomic.StoreUint32(&pipecli.state, uint32(EndpointStateDisconnecting))
	logDebugf("Pipeline client `%s/%p` waiting for client shutdown", pipecli.address, pipecli)

	// We must wait for the close wait goroutine to die as well before we can continue.
	<-killSig

	logDebugf("Pipeline client `%s/%p` received client shutdown notification", pipecli.address, pipecli)
}

func (pipecli *memdPipelineClient) Run() {
	for {
		logDebugf("Pipeline Client `%s/%p` preparing for new client loop", pipecli.address, pipecli)
		atomic.StoreUint32(&pipecli.state, uint32(EndpointStateConnecting))

		pipecli.lock.Lock()
		pipeline := pipecli.parent
		pipecli.lock.Unlock()

		if pipeline == nil {
			// If our pipeline is nil, it indicates that we need to shut down.
			logDebugf("Pipeline Client `%s/%p` is shutting down", pipecli.address, pipecli)
			break
		}

		logDebugf("Pipeline Client `%s/%p` retrieving new client connection for parent %p", pipecli.address, pipecli, pipeline)
		client, err := pipeline.getClientFn()
		if err != nil {
			atomic.StoreUint32(&pipecli.state, uint32(EndpointStateDisconnected))
			pipecli.lock.Lock()
			pipecli.connectError = err
			pipecli.lock.Unlock()
			continue
		}
		pipecli.lock.Lock()
		pipecli.connectError = nil
		pipecli.lock.Unlock()
		atomic.StoreUint32(&pipecli.state, uint32(EndpointStateConnected))

		// Runs until the connection has died (for whatever reason)
		logDebugf("Pipeline Client `%s/%p` starting new client loop for %p", pipecli.address, pipecli, client)
		pipecli.ioLoop(client)
	}

	// Lets notify anyone who is watching that we are now shut down
	close(pipecli.closedSig)
}

// Close will close this pipeline client.  Note that this method will not wait for
// everything to be cleaned up before returning.
func (pipecli *memdPipelineClient) Close() error {
	logDebugf("Pipeline Client `%s/%p` received close request", pipecli.address, pipecli)
	atomic.StoreUint32(&pipecli.state, uint32(EndpointStateDisconnecting))

	// To shut down the client, we remove our reference to the parent. This
	// causes our ioLoop see that we are being shut down and perform cleanup
	// before exiting.
	pipecli.lock.Lock()
	pipecli.parent = nil
	activeConsumer := pipecli.consumer
	pipecli.consumer = nil
	pipecli.lock.Unlock()

	// If we have an active consumer, we need to close it to cause the running
	// ioLoop to unpause and pick up that our parent has been removed.  Note
	// that in some cases, we might not have an active consumer. This means
	// that the ioLoop is about to try and fetch one, finding the missing
	// parent in doing so.
	if activeConsumer != nil {
		activeConsumer.Close()
	}

	// Lets wait till the ioLoop has shut everything down before returning.
	<-pipecli.closedSig
	atomic.StoreUint32(&pipecli.state, uint32(EndpointStateDisconnected))

	logDebugf("Pipeline Client `%s/%p` has exited", pipecli.address, pipecli)

	return nil
}
