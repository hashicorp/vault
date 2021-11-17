package gocb

import (
	"context"
	"sync"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v10"
)

func (c *Collection) observeOnceSeqNo(
	ctx context.Context,
	trace RequestSpan,
	docID string,
	mt gocbcore.MutationToken,
	replicaIdx int,
	cancelCh chan struct{},
	timeout time.Duration,
	user string,
) (didReplicate, didPersist bool, errOut error) {
	opm := c.newKvOpManager("observe_once", trace)
	defer opm.Finish(true)

	opm.SetDocumentID(docID)
	opm.SetCancelCh(cancelCh)
	opm.SetTimeout(timeout)
	opm.SetImpersonate(user)
	opm.SetContext(ctx)

	agent, err := c.getKvProvider()
	if err != nil {
		return false, false, err
	}
	err = opm.Wait(agent.ObserveVb(gocbcore.ObserveVbOptions{
		VbID:         mt.VbID,
		VbUUID:       mt.VbUUID,
		ReplicaIdx:   replicaIdx,
		TraceContext: opm.TraceSpanContext(),
		Deadline:     opm.Deadline(),
		User:         opm.Impersonate(),
	}, func(res *gocbcore.ObserveVbResult, err error) {
		if err != nil || res == nil {
			errOut = opm.EnhanceErr(err)
			opm.Reject()
			return
		}

		didReplicate = res.CurrentSeqNo >= mt.SeqNo
		didPersist = res.PersistSeqNo >= mt.SeqNo

		opm.Resolve(nil)
	}))
	if err == nil {
		errOut = err
	}
	return
}

func (c *Collection) observeOne(
	ctx context.Context,
	trace RequestSpan,
	docID string,
	mt gocbcore.MutationToken,
	replicaIdx int,
	replicaCh, persistCh, cancelCh chan struct{},
	timeout time.Duration,
	user string,
) {
	sentReplicated := false
	sentPersisted := false

	calc := gocbcore.ExponentialBackoff(10*time.Microsecond, 100*time.Millisecond, 0)
	retries := uint32(0)

ObserveLoop:
	for {
		select {
		case <-cancelCh:
			break ObserveLoop
		default:
			// not cancelled yet
		}

		didReplicate, didPersist, err := c.observeOnceSeqNo(ctx, trace, docID, mt, replicaIdx, cancelCh, timeout, user)
		if err != nil {
			logDebugf("ObserveOnce failed unexpected: %s", err)
			return
		}

		if didReplicate && !sentReplicated {
			replicaCh <- struct{}{}
			sentReplicated = true
		}

		if didPersist && !sentPersisted {
			persistCh <- struct{}{}
			sentPersisted = true
		}

		// If we've got persisted and replicated, we can just stop
		if sentPersisted && sentReplicated {
			break ObserveLoop
		}

		waitTmr := gocbcore.AcquireTimer(calc(retries))
		retries++
		select {
		case <-waitTmr.C:
			gocbcore.ReleaseTimer(waitTmr, true)
		case <-cancelCh:
			gocbcore.ReleaseTimer(waitTmr, false)
		}
	}
}

func (c *Collection) waitForDurability(
	ctx context.Context,
	trace RequestSpan,
	docID string,
	mt gocbcore.MutationToken,
	replicateTo uint,
	persistTo uint,
	deadline time.Time,
	cancelCh chan struct{},
	user string,
) error {
	opm := c.newKvOpManager("observe", trace)
	defer opm.Finish(true)

	opm.SetDocumentID(docID)

	agent, err := c.getKvProvider()
	if err != nil {
		return err
	}

	snapshot, err := agent.ConfigSnapshot()
	if err != nil {
		return err
	}

	numReplicas, err := snapshot.NumReplicas()
	if err != nil {
		return err
	}

	numServers := numReplicas + 1
	if replicateTo > uint(numServers-1) || persistTo > uint(numServers) {
		return opm.EnhanceErr(ErrDurabilityImpossible)
	}

	subOpCancelCh := make(chan struct{}, 1)
	replicaCh := make(chan struct{}, numServers)
	persistCh := make(chan struct{}, numServers)

	// If we cancel the sub ops then we need to wait for cancellation to complete before we exit, otherwise
	// we will attempt to close our span before the child spans complete.
	var wg sync.WaitGroup
	for replicaIdx := 0; replicaIdx < numServers; replicaIdx++ {
		wg.Add(1)
		go func(ridx int) {
			c.observeOne(ctx, opm.TraceSpan(), docID, mt, ridx, replicaCh, persistCh, subOpCancelCh,
				time.Until(deadline), user)
			wg.Done()
		}(replicaIdx)
	}

	numReplicated := uint(0)
	numPersisted := uint(0)

	for {
		select {
		case <-replicaCh:
			numReplicated++
		case <-persistCh:
			numPersisted++
		case <-time.After(time.Until(deadline)):
			// deadline exceeded
			close(subOpCancelCh)
			wg.Wait()
			return opm.EnhanceErr(ErrAmbiguousTimeout)
		case <-cancelCh:
			// parent asked for cancellation
			close(subOpCancelCh)
			wg.Wait()
			return opm.EnhanceErr(ErrRequestCanceled)
		}

		if numReplicated >= replicateTo && numPersisted >= persistTo {
			close(subOpCancelCh)
			wg.Wait()
			return nil
		}
	}
}
