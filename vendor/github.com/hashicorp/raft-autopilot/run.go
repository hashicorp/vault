package autopilot

import (
	"context"
	"time"
)

// Start will launch the go routines in the background to perform Autopilot.
// When the context passed in is cancelled or the Stop method is called
// then these routines will exit.
func (a *Autopilot) Start(ctx context.Context) {
	a.runLock.Lock()
	defer a.runLock.Unlock()

	// already running so there is nothing to do
	if a.running {
		return
	}

	ctx, shutdown := context.WithCancel(ctx)
	a.shutdown = shutdown
	a.startTime = a.time.Now()
	a.done = make(chan struct{})

	// While a go routine executed by a.run below will periodically
	// update the state, we want to go ahead and force updating it now
	// so that during a leadership transfer we don't report an empty
	// autopilot state. We put a pretty small timeout on this though
	// so as to prevent leader establishment from taking too long
	updateCtx, updateCancel := context.WithTimeout(ctx, time.Second)
	defer updateCancel()
	a.updateState(updateCtx)

	go a.run(ctx)
	a.running = true
}

// Stop will terminate the go routines being executed to perform autopilot.
func (a *Autopilot) Stop() <-chan struct{} {
	a.runLock.Lock()
	defer a.runLock.Unlock()

	// Nothing to do
	if !a.running {
		done := make(chan struct{})
		close(done)
		return done
	}

	a.shutdown()
	return a.done
}

func (a *Autopilot) run(ctx context.Context) {
	a.logger.Debug("autopilot is now running")
	// autopilot needs to do 3 things
	//
	// 1. periodically update the cluster state
	// 2. periodically check for and perform promotions and demotions
	// 3. Respond to servers leaving and prune dead servers
	//
	// We could attempt to do all of this in a single go routine except that
	// updating the cluster health could potentially take long enough to impact
	// the periodicity of the promotions and demotions performed by task 2/3.
	// So instead this go routine will spawn a second go routine to manage
	// updating the cluster health in the background. This go routine is still
	// in control of the overall running status and will not exit until the
	// child go routine has exited.

	// child go routine for cluster health updating
	stateUpdaterDone := make(chan struct{})
	go a.runStateUpdater(ctx, stateUpdaterDone)

	// cleanup for once we are stopped
	defer func() {
		// block waiting for our child go routine to also finish
		<-stateUpdaterDone

		a.logger.Debug("autopilot is now stopped")

		a.runLock.Lock()
		a.shutdown = nil
		a.running = false
		// this should be the final cleanup task as it is what notifies the rest
		// of the world that we are now done
		close(a.done)
		a.done = nil
		a.runLock.Unlock()
	}()

	reconcileTicker := time.NewTicker(a.reconcileInterval)
	defer reconcileTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-reconcileTicker.C:
			if err := a.reconcile(); err != nil {
				a.logger.Error("Failed to reconcile current state with the desired state")
			}

			if err := a.pruneDeadServers(); err != nil {
				a.logger.Error("Failed to prune dead servers", "error", err)
			}
		case <-a.removeDeadCh:
			if err := a.pruneDeadServers(); err != nil {
				a.logger.Error("Failed to prune dead servers", "error", err)
			}
		}
	}
}

// runStateUpdated will periodically update the autopilot state until the context
// passed in is cancelled. When finished the provide done chan will be closed.
func (a *Autopilot) runStateUpdater(ctx context.Context, done chan struct{}) {
	a.logger.Debug("state update routine is now running")
	defer func() {
		a.logger.Debug("state update routine is now stopped")
		close(done)
	}()

	ticker := time.NewTicker(a.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.updateState(ctx)
		}
	}
}
