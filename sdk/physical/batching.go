package physical

import (
	"context"
	"time"
)

type batchValue struct{}

var contextBatch batchValue = struct{}{}

func ContextBatch(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextBatch, "1")
}

func IsBatchContext(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	boolRaw := ctx.Value(contextBatch)
	return boolRaw != nil && boolRaw.(string) != ""
}

type batchRequest struct {
	op      Operation
	entry   Entry
	errChan chan error
}

type Batcher struct {
	timer         *time.Ticker
	batchMax      int
	submit        chan batchRequest
	activeContext context.Context
	storage       Backend
}

var _ Backend = &Batcher{}

func NewBatcher(storage Backend) *Batcher {
	return &Batcher{
		timer:    time.NewTicker(100 * time.Millisecond),
		batchMax: 64, // consul only allows 64
		submit:   make(chan batchRequest),
		storage:  storage,
	}
}

func (d *Batcher) Start(ctx context.Context) {
	d.activeContext = ctx
	var reqs []batchRequest
	for {
		select {
		case <-d.activeContext.Done():
			// Should we write to accumulated reqs errChans?
			return
		case req := <-d.submit:
			reqs = append(reqs, req)
			if len(reqs) >= d.batchMax {
				d.run(reqs)
				reqs = reqs[:0]
			}
		case <-d.timer.C:
			if len(reqs) > 0 {
				d.run(reqs)
				reqs = reqs[:0]
			}
		}
	}
}

func (d *Batcher) List(ctx context.Context, key string) ([]string, error) {
	return d.storage.List(ctx, key)
}

func (d *Batcher) Get(ctx context.Context, key string) (*Entry, error) {
	return d.storage.Get(ctx, key)
}

func (d *Batcher) Put(ctx context.Context, entry *Entry) error {
	if !IsBatchContext(ctx) || !d.isTransactional() {
		return d.storage.Put(ctx, entry)
	}

	errChan := make(chan error, 1)
	d.submit <- batchRequest{
		op:      PutOperation,
		entry:   *entry,
		errChan: errChan,
	}
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (d *Batcher) isTransactional() bool {
	_, ok := d.storage.(Transactional)
	return ok
}

func (d *Batcher) Delete(ctx context.Context, key string) error {
	if !IsBatchContext(ctx) || !d.isTransactional() {
		return d.storage.Delete(ctx, key)
	}

	errChan := make(chan error)
	d.submit <- batchRequest{
		op:      DeleteOperation,
		entry:   Entry{Key: key},
		errChan: errChan,
	}
	select {
	case err := <-errChan:
		close(errChan)
		return err
	case <-ctx.Done():
		// Make sure we don't block the run goroutine's attempt to send the error
		go func() { <-errChan }()
		return ctx.Err()
	}
}

func (d *Batcher) run(reqs []batchRequest) {
	// TODO should we eliminate dups?
	txn := make([]*TxnEntry, len(reqs))
	for i, req := range reqs {
		txn[i] = &TxnEntry{
			Operation: req.op,
			Entry:     &req.entry,
		}
	}
	err := d.storage.(Transactional).Transaction(d.activeContext, txn)
	for i := range reqs {
		reqs[i].errChan <- err
	}
}
