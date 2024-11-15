package gocb

import (
	"context"
	"sync"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v10"
)

func (p *kvProviderCore) observeOnceSeqNo(
	ctx context.Context,
	c *Collection,
	trace RequestSpan,
	docID string,
	mt gocbcore.MutationToken,
	replicaIdx int,
	cancelCh chan struct{},
	timeout time.Duration,
	user string,
) (didReplicate, didPersist bool, errOut error) {
	observedOpm := newKvOpManagerCore(c, "observe_once", trace, p)
	defer observedOpm.Finish(true)

	observedOpm.SetDocumentID(docID)
	observedOpm.SetCancelCh(cancelCh)
	observedOpm.SetTimeout(timeout)
	observedOpm.SetImpersonate(user)
	observedOpm.SetContext(ctx)

	err := observedOpm.Wait(p.agent.ObserveVb(gocbcore.ObserveVbOptions{
		VbID:         mt.VbID,
		VbUUID:       mt.VbUUID,
		ReplicaIdx:   replicaIdx,
		TraceContext: observedOpm.TraceSpanContext(),
		Deadline:     observedOpm.Deadline(),
		User:         observedOpm.Impersonate(),
	}, func(res *gocbcore.ObserveVbResult, err error) {
		if err != nil || res == nil {
			errOut = observedOpm.EnhanceErr(err)
			observedOpm.Reject()
			return
		}

		didReplicate = res.CurrentSeqNo >= mt.SeqNo
		didPersist = res.PersistSeqNo >= mt.SeqNo

		observedOpm.Resolve(nil)
	}))
	if err == nil {
		errOut = err
	}
	return
}

func (p *kvProviderCore) observeOne(
	ctx context.Context,
	c *Collection,
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

		didReplicate, didPersist, err := p.observeOnceSeqNo(ctx, c, trace, docID, mt, replicaIdx, cancelCh, timeout, user)
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

func (p *kvProviderCore) waitForDurability(
	ctx context.Context,
	c *Collection,
	trace RequestSpan,
	docID string,
	mt gocbcore.MutationToken,
	replicateTo uint,
	persistTo uint,
	deadline time.Time,
	cancelCh chan struct{},
	user string,
) error {
	observeOpm := newKvOpManagerCore(c, "observe", trace, p)
	defer observeOpm.Finish(true)

	observeOpm.SetDocumentID(docID)

	snapshot, err := p.snapshotProvider.WaitForConfigSnapshot(ctx, deadline)
	if err != nil {
		return err
	}

	numReplicas, err := snapshot.NumReplicas()
	if err != nil {
		return err
	}

	numServers := numReplicas + 1
	if replicateTo > uint(numServers-1) || persistTo > uint(numServers) {
		return observeOpm.EnhanceErr(ErrDurabilityImpossible)
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
			p.observeOne(ctx, c, observeOpm.TraceSpan(), docID, mt, ridx, replicaCh, persistCh, subOpCancelCh,
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
			return observeOpm.EnhanceErr(ErrAmbiguousTimeout)
		case <-cancelCh:
			// parent asked for cancellation
			close(subOpCancelCh)
			wg.Wait()
			return observeOpm.EnhanceErr(ErrRequestCanceled)
		}

		if numReplicated >= replicateTo && numPersisted >= persistTo {
			close(subOpCancelCh)
			wg.Wait()
			return nil
		}
	}
}
