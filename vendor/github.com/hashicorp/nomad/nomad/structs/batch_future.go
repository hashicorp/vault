package structs

// BatchFuture is used to wait on a batch update to complete
type BatchFuture struct {
	doneCh chan struct{}
	err    error
	index  uint64
}

// NewBatchFuture creates a new batch future
func NewBatchFuture() *BatchFuture {
	return &BatchFuture{
		doneCh: make(chan struct{}),
	}
}

// Wait is used to block for the future to complete and returns the error
func (b *BatchFuture) Wait() error {
	<-b.doneCh
	return b.err
}

// WaitCh is used to block for the future to complete
func (b *BatchFuture) WaitCh() <-chan struct{} {
	return b.doneCh
}

// Error is used to return the error of the batch, only after Wait()
func (b *BatchFuture) Error() error {
	return b.err
}

// Index is used to return the index of the batch, only after Wait()
func (b *BatchFuture) Index() uint64 {
	return b.index
}

// Respond is used to unblock the future
func (b *BatchFuture) Respond(index uint64, err error) {
	b.index = index
	b.err = err
	close(b.doneCh)
}
