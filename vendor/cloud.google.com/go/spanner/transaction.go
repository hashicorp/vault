/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/internal/trace"
	vkit "cloud.google.com/go/spanner/apiv1"
	"github.com/golang/protobuf/proto"
	"google.golang.org/api/iterator"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// transactionID stores a transaction ID which uniquely identifies a transaction
// in Cloud Spanner.
type transactionID []byte

// txReadEnv manages a read-transaction environment consisting of a session
// handle and a transaction selector.
type txReadEnv interface {
	// acquire returns a read-transaction environment that can be used to
	// perform a transactional read.
	acquire(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error)
	// sets the transaction's read timestamp
	setTimestamp(time.Time)
	// release should be called at the end of every transactional read to deal
	// with session recycling.
	release(error)
}

// txReadOnly contains methods for doing transactional reads.
type txReadOnly struct {
	// read-transaction environment for performing transactional read
	// operations.
	txReadEnv

	// Atomic. Only needed for DML statements, but used forall.
	sequenceNumber int64

	// replaceSessionFunc is a function that can be called to replace the
	// session that is used by the transaction. This function should only be
	// defined for single-use transactions that can safely be retried on a
	// different session. All other transactions will set this function to nil.
	replaceSessionFunc func(ctx context.Context) error

	// sp is the session pool for allocating a session to execute the read-only
	// transaction. It is set only once during initialization of the
	// txReadOnly.
	sp *sessionPool
	// sh is the sessionHandle allocated from sp.
	sh *sessionHandle

	// qo provides options for executing a sql query.
	qo QueryOptions
}

// TransactionOptions provides options for a transaction.
type TransactionOptions struct {
}

// errSessionClosed returns error for using a recycled/destroyed session
func errSessionClosed(sh *sessionHandle) error {
	return spannerErrorf(codes.FailedPrecondition,
		"session is already recycled / destroyed: session_id = %q, rpc_client = %v", sh.getID(), sh.getClient())
}

// Read returns a RowIterator for reading multiple rows from the database.
func (t *txReadOnly) Read(ctx context.Context, table string, keys KeySet, columns []string) *RowIterator {
	return t.ReadWithOptions(ctx, table, keys, columns, nil)
}

// ReadUsingIndex calls ReadWithOptions with ReadOptions{Index: index}.
func (t *txReadOnly) ReadUsingIndex(ctx context.Context, table, index string, keys KeySet, columns []string) (ri *RowIterator) {
	return t.ReadWithOptions(ctx, table, keys, columns, &ReadOptions{Index: index})
}

// ReadOptions provides options for reading rows from a database.
type ReadOptions struct {
	// The index to use for reading. If non-empty, you can only read columns
	// that are part of the index key, part of the primary key, or stored in the
	// index due to a STORING clause in the index definition.
	Index string

	// The maximum number of rows to read. A limit value less than 1 means no
	// limit.
	Limit int
}

// ReadWithOptions returns a RowIterator for reading multiple rows from the
// database. Pass a ReadOptions to modify the read operation.
func (t *txReadOnly) ReadWithOptions(ctx context.Context, table string, keys KeySet, columns []string, opts *ReadOptions) (ri *RowIterator) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.Read")
	defer func() { trace.EndSpan(ctx, ri.err) }()
	var (
		sh  *sessionHandle
		ts  *sppb.TransactionSelector
		err error
	)
	kset, err := keys.keySetProto()
	if err != nil {
		return &RowIterator{err: err}
	}
	if sh, ts, err = t.acquire(ctx); err != nil {
		return &RowIterator{err: err}
	}
	// Cloud Spanner will return "Session not found" on bad sessions.
	client := sh.getClient()
	if client == nil {
		// Might happen if transaction is closed in the middle of a API call.
		return &RowIterator{err: errSessionClosed(sh)}
	}
	index := ""
	limit := 0
	if opts != nil {
		index = opts.Index
		if opts.Limit > 0 {
			limit = opts.Limit
		}
	}
	return streamWithReplaceSessionFunc(
		contextWithOutgoingMetadata(ctx, sh.getMetadata()),
		sh.session.logger,
		func(ctx context.Context, resumeToken []byte) (streamingReceiver, error) {
			return client.StreamingRead(ctx,
				&sppb.ReadRequest{
					Session:     t.sh.getID(),
					Transaction: ts,
					Table:       table,
					Index:       index,
					Columns:     columns,
					KeySet:      kset,
					ResumeToken: resumeToken,
					Limit:       int64(limit),
				})
		},
		t.replaceSessionFunc,
		t.setTimestamp,
		t.release,
	)
}

// errRowNotFound returns error for not being able to read the row identified by
// key.
func errRowNotFound(table string, key Key) error {
	return spannerErrorf(codes.NotFound, "row not found(Table: %v, PrimaryKey: %v)", table, key)
}

// errRowNotFoundByIndex returns error for not being able to read the row by index.
func errRowNotFoundByIndex(table string, key Key, index string) error {
	return spannerErrorf(codes.NotFound, "row not found(Table: %v, IndexKey: %v, Index: %v)", table, key, index)
}

// errMultipleRowsFound returns error for receiving more than one row when reading a single row using an index.
func errMultipleRowsFound(table string, key Key, index string) error {
	return spannerErrorf(codes.FailedPrecondition, "more than one row found by index(Table: %v, IndexKey: %v, Index: %v)", table, key, index)
}

// ReadRow reads a single row from the database.
//
// If no row is present with the given key, then ReadRow returns an error where
// spanner.ErrCode(err) is codes.NotFound.
func (t *txReadOnly) ReadRow(ctx context.Context, table string, key Key, columns []string) (*Row, error) {
	iter := t.Read(ctx, table, key, columns)
	defer iter.Stop()
	row, err := iter.Next()
	switch err {
	case iterator.Done:
		return nil, errRowNotFound(table, key)
	case nil:
		return row, nil
	default:
		return nil, err
	}
}

// ReadRowUsingIndex reads a single row from the database using an index.
//
// If no row is present with the given index, then ReadRowUsingIndex returns an
// error where spanner.ErrCode(err) is codes.NotFound.
//
// If more than one row received with the given index, then ReadRowUsingIndex
// returns an error where spanner.ErrCode(err) is codes.FailedPrecondition.
func (t *txReadOnly) ReadRowUsingIndex(ctx context.Context, table string, index string, key Key, columns []string) (*Row, error) {
	iter := t.ReadUsingIndex(ctx, table, index, key, columns)
	defer iter.Stop()
	row, err := iter.Next()
	switch err {
	case iterator.Done:
		return nil, errRowNotFoundByIndex(table, key, index)
	case nil:
		// If more than one row found, return an error.
		_, err := iter.Next()
		switch err {
		case iterator.Done:
			return row, nil
		case nil:
			return nil, errMultipleRowsFound(table, key, index)
		default:
			return nil, err
		}
	default:
		return nil, err
	}
}

// QueryOptions provides options for executing a sql query from a database.
type QueryOptions struct {
	Mode    *sppb.ExecuteSqlRequest_QueryMode
	Options *sppb.ExecuteSqlRequest_QueryOptions
}

// merge combines two QueryOptions that the input parameter will have higher
// order of precedence.
func (qo QueryOptions) merge(opts QueryOptions) QueryOptions {
	merged := QueryOptions{
		Mode:    qo.Mode,
		Options: &sppb.ExecuteSqlRequest_QueryOptions{},
	}
	if opts.Mode != nil {
		merged.Mode = opts.Mode
	}
	proto.Merge(merged.Options, qo.Options)
	proto.Merge(merged.Options, opts.Options)
	return merged
}

// Query executes a query against the database. It returns a RowIterator for
// retrieving the resulting rows.
//
// Query returns only row data, without a query plan or execution statistics.
// Use QueryWithStats to get rows along with the plan and statistics. Use
// AnalyzeQuery to get just the plan.
func (t *txReadOnly) Query(ctx context.Context, statement Statement) *RowIterator {
	mode := sppb.ExecuteSqlRequest_NORMAL
	return t.query(ctx, statement, QueryOptions{
		Mode:    &mode,
		Options: t.qo.Options,
	})
}

// QueryWithOptions executes a SQL statment against the database. It returns
// a RowIterator for retrieving the resulting rows. The sql query execution
// will be optimized based on the given query options.
func (t *txReadOnly) QueryWithOptions(ctx context.Context, statement Statement, opts QueryOptions) *RowIterator {
	return t.query(ctx, statement, t.qo.merge(opts))
}

// QueryWithStats executes a SQL statement against the database. It returns
// a RowIterator for retrieving the resulting rows. The RowIterator will also
// be populated with a query plan and execution statistics.
func (t *txReadOnly) QueryWithStats(ctx context.Context, statement Statement) *RowIterator {
	mode := sppb.ExecuteSqlRequest_PROFILE
	return t.query(ctx, statement, QueryOptions{
		Mode:    &mode,
		Options: t.qo.Options,
	})
}

// AnalyzeQuery returns the query plan for statement.
func (t *txReadOnly) AnalyzeQuery(ctx context.Context, statement Statement) (*sppb.QueryPlan, error) {
	mode := sppb.ExecuteSqlRequest_PLAN
	iter := t.query(ctx, statement, QueryOptions{
		Mode:    &mode,
		Options: t.qo.Options,
	})
	defer iter.Stop()
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	if iter.QueryPlan == nil {
		return nil, spannerErrorf(codes.Internal, "query plan unavailable")
	}
	return iter.QueryPlan, nil
}

func (t *txReadOnly) query(ctx context.Context, statement Statement, options QueryOptions) (ri *RowIterator) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.Query")
	defer func() { trace.EndSpan(ctx, ri.err) }()
	req, sh, err := t.prepareExecuteSQL(ctx, statement, options)
	if err != nil {
		return &RowIterator{err: err}
	}
	client := sh.getClient()
	return streamWithReplaceSessionFunc(
		contextWithOutgoingMetadata(ctx, sh.getMetadata()),
		sh.session.logger,
		func(ctx context.Context, resumeToken []byte) (streamingReceiver, error) {
			req.ResumeToken = resumeToken
			req.Session = t.sh.getID()
			return client.ExecuteStreamingSql(ctx, req)
		},
		t.replaceSessionFunc,
		t.setTimestamp,
		t.release)
}

func (t *txReadOnly) prepareExecuteSQL(ctx context.Context, stmt Statement, options QueryOptions) (*sppb.ExecuteSqlRequest, *sessionHandle, error) {
	sh, ts, err := t.acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	// Cloud Spanner will return "Session not found" on bad sessions.
	sid := sh.getID()
	if sid == "" {
		// Might happen if transaction is closed in the middle of a API call.
		return nil, nil, errSessionClosed(sh)
	}
	params, paramTypes, err := stmt.convertParams()
	if err != nil {
		return nil, nil, err
	}
	mode := sppb.ExecuteSqlRequest_NORMAL
	if options.Mode != nil {
		mode = *options.Mode
	}
	req := &sppb.ExecuteSqlRequest{
		Session:      sid,
		Transaction:  ts,
		Sql:          stmt.SQL,
		QueryMode:    mode,
		Seqno:        atomic.AddInt64(&t.sequenceNumber, 1),
		Params:       params,
		ParamTypes:   paramTypes,
		QueryOptions: options.Options,
	}
	return req, sh, nil
}

// txState is the status of a transaction.
type txState int

const (
	// transaction is new, waiting to be initialized..
	txNew txState = iota
	// transaction is being initialized.
	txInit
	// transaction is active and can perform read/write.
	txActive
	// transaction is closed, cannot be used anymore.
	txClosed
)

// errRtsUnavailable returns error for read transaction's read timestamp being
// unavailable.
func errRtsUnavailable() error {
	return spannerErrorf(codes.Internal, "read timestamp is unavailable")
}

// errTxClosed returns error for using a closed transaction.
func errTxClosed() error {
	return spannerErrorf(codes.InvalidArgument, "cannot use a closed transaction")
}

// errUnexpectedTxState returns error for transaction enters an unexpected state.
func errUnexpectedTxState(ts txState) error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected transaction state: %v", ts)
}

// ReadOnlyTransaction provides a snapshot transaction with guaranteed
// consistency across reads, but does not allow writes.  Read-only transactions
// can be configured to read at timestamps in the past.
//
// Read-only transactions do not take locks. Instead, they work by choosing a
// Cloud Spanner timestamp, then executing all reads at that timestamp. Since
// they do not acquire locks, they do not block concurrent read-write
// transactions.
//
// Unlike locking read-write transactions, read-only transactions never abort.
// They can fail if the chosen read timestamp is garbage collected; however, the
// default garbage collection policy is generous enough that most applications
// do not need to worry about this in practice. See the documentation of
// TimestampBound for more details.
//
// A ReadOnlyTransaction consumes resources on the server until Close is called.
type ReadOnlyTransaction struct {
	// mu protects concurrent access to the internal states of ReadOnlyTransaction.
	mu sync.Mutex
	// txReadOnly contains methods for performing transactional reads.
	txReadOnly
	// singleUse indicates that the transaction can be used for only one read.
	singleUse bool
	// tx is the transaction ID in Cloud Spanner that uniquely identifies the
	// ReadOnlyTransaction.
	tx transactionID
	// txReadyOrClosed is for broadcasting that transaction ID has been returned
	// by Cloud Spanner or that transaction is closed.
	txReadyOrClosed chan struct{}
	// state is the current transaction status of the ReadOnly transaction.
	state txState
	// rts is the read timestamp returned by transactional reads.
	rts time.Time
	// tb is the read staleness bound specification for transactional reads.
	tb TimestampBound
}

// errTxInitTimeout returns error for timeout in waiting for initialization of
// the transaction.
func errTxInitTimeout() error {
	return spannerErrorf(codes.Canceled, "timeout/context canceled in waiting for transaction's initialization")
}

// getTimestampBound returns the read staleness bound specified for the
// ReadOnlyTransaction.
func (t *ReadOnlyTransaction) getTimestampBound() TimestampBound {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tb
}

// begin starts a snapshot read-only Transaction on Cloud Spanner.
func (t *ReadOnlyTransaction) begin(ctx context.Context) error {
	var (
		locked bool
		tx     transactionID
		rts    time.Time
		sh     *sessionHandle
		err    error
		res    *sppb.Transaction
	)
	defer func() {
		if !locked {
			t.mu.Lock()
			// Not necessary, just to make it clear that t.mu is being held when
			// locked == true.
			locked = true
		}
		if t.state != txClosed {
			// Signal other initialization routines.
			close(t.txReadyOrClosed)
			t.txReadyOrClosed = make(chan struct{})
		}
		t.mu.Unlock()
		if err != nil && sh != nil {
			// Got a valid session handle, but failed to initialize transaction=
			// on Cloud Spanner.
			if isSessionNotFoundError(err) {
				sh.destroy()
			}
			// If sh.destroy was already executed, this becomes a noop.
			sh.recycle()
		}
	}()
	// Retry the BeginTransaction call if a 'Session not found' is returned.
	for {
		sh, err = t.sp.take(ctx)
		if err != nil {
			return err
		}
		res, err = sh.getClient().BeginTransaction(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.BeginTransactionRequest{
			Session: sh.getID(),
			Options: &sppb.TransactionOptions{
				Mode: &sppb.TransactionOptions_ReadOnly_{
					ReadOnly: buildTransactionOptionsReadOnly(t.getTimestampBound(), true),
				},
			},
		})
		if isSessionNotFoundError(err) {
			sh.destroy()
			continue
		} else if err == nil {
			tx = res.Id
			if res.ReadTimestamp != nil {
				rts = time.Unix(res.ReadTimestamp.Seconds, int64(res.ReadTimestamp.Nanos))
			}
		} else {
			err = ToSpannerError(err)
		}
		break
	}
	t.mu.Lock()

	// defer function will be executed with t.mu being held.
	locked = true

	// During the execution of t.begin(), t.Close() was invoked.
	if t.state == txClosed {
		return errSessionClosed(sh)
	}

	// If begin() fails, this allows other queries to take over the
	// initialization.
	t.tx = nil
	if err == nil {
		t.tx = tx
		t.rts = rts
		t.sh = sh
		// State transite to txActive.
		t.state = txActive
	}
	return err
}

// acquire implements txReadEnv.acquire.
func (t *ReadOnlyTransaction) acquire(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	if err := checkNestedTxn(ctx); err != nil {
		return nil, nil, err
	}
	if t.singleUse {
		return t.acquireSingleUse(ctx)
	}
	return t.acquireMultiUse(ctx)
}

func (t *ReadOnlyTransaction) acquireSingleUse(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	switch t.state {
	case txClosed:
		// A closed single-use transaction can never be reused.
		return nil, nil, errTxClosed()
	case txNew:
		t.state = txClosed
		ts := &sppb.TransactionSelector{
			Selector: &sppb.TransactionSelector_SingleUse{
				SingleUse: &sppb.TransactionOptions{
					Mode: &sppb.TransactionOptions_ReadOnly_{
						ReadOnly: buildTransactionOptionsReadOnly(t.tb, true),
					},
				},
			},
		}
		sh, err := t.sp.take(ctx)
		if err != nil {
			return nil, nil, err
		}

		// Install session handle into t, which can be used for readonly
		// operations later.
		t.sh = sh
		return sh, ts, nil
	}
	us := t.state

	// SingleUse transaction should only be in either txNew state or txClosed
	// state.
	return nil, nil, errUnexpectedTxState(us)
}

func (t *ReadOnlyTransaction) acquireMultiUse(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	for {
		t.mu.Lock()
		switch t.state {
		case txClosed:
			t.mu.Unlock()
			return nil, nil, errTxClosed()
		case txNew:
			// State transit to txInit so that no further TimestampBound change
			// is accepted.
			t.state = txInit
			t.mu.Unlock()
			continue
		case txInit:
			if t.tx != nil {
				// Wait for a transaction ID to become ready.
				txReadyOrClosed := t.txReadyOrClosed
				t.mu.Unlock()
				select {
				case <-txReadyOrClosed:
					// Need to check transaction state again.
					continue
				case <-ctx.Done():
					// The waiting for initialization is timeout, return error
					// directly.
					return nil, nil, errTxInitTimeout()
				}
			}
			// Take the ownership of initializing the transaction.
			t.tx = transactionID{}
			t.mu.Unlock()
			// Begin a read-only transaction.
			//
			// TODO: consider adding a transaction option which allow queries to
			//  initiate transactions by themselves. Note that this option might
			//  not be always good because the ID of the new transaction won't
			//  be ready till the query returns some data or completes.
			if err := t.begin(ctx); err != nil {
				return nil, nil, err
			}

			// If t.begin() succeeded, t.state should have been changed to
			// txActive, so we can just continue here.
			continue
		case txActive:
			sh := t.sh
			ts := &sppb.TransactionSelector{
				Selector: &sppb.TransactionSelector_Id{
					Id: t.tx,
				},
			}
			t.mu.Unlock()
			return sh, ts, nil
		}
		state := t.state
		t.mu.Unlock()
		return nil, nil, errUnexpectedTxState(state)
	}
}

func (t *ReadOnlyTransaction) setTimestamp(ts time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.rts.IsZero() {
		t.rts = ts
	}
}

// release implements txReadEnv.release.
func (t *ReadOnlyTransaction) release(err error) {
	t.mu.Lock()
	sh := t.sh
	t.mu.Unlock()
	if sh != nil { // sh could be nil if t.acquire() fails.
		if isSessionNotFoundError(err) {
			sh.destroy()
		}
		if t.singleUse {
			// If session handle is already destroyed, this becomes a noop.
			sh.recycle()
		}
	}
}

// Close closes a ReadOnlyTransaction, the transaction cannot perform any reads
// after being closed.
func (t *ReadOnlyTransaction) Close() {
	if t.singleUse {
		return
	}
	t.mu.Lock()
	if t.state != txClosed {
		t.state = txClosed
		close(t.txReadyOrClosed)
	}
	sh := t.sh
	t.mu.Unlock()
	if sh == nil {
		return
	}
	// If session handle is already destroyed, this becomes a noop. If there are
	// still active queries and if the recycled session is reused before they
	// complete, Cloud Spanner will cancel them on behalf of the new transaction
	// on the session.
	if sh != nil {
		sh.recycle()
	}
}

// Timestamp returns the timestamp chosen to perform reads and queries in this
// transaction. The value can only be read after some read or query has either
// returned some data or completed without returning any data.
func (t *ReadOnlyTransaction) Timestamp() (time.Time, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.rts.IsZero() {
		return t.rts, errRtsUnavailable()
	}
	return t.rts, nil
}

// WithTimestampBound specifies the TimestampBound to use for read or query.
// This can only be used before the first read or query is invoked. Note:
// bounded staleness is not available with general ReadOnlyTransactions; use a
// single-use ReadOnlyTransaction instead.
//
// The returned value is the ReadOnlyTransaction so calls can be chained.
func (t *ReadOnlyTransaction) WithTimestampBound(tb TimestampBound) *ReadOnlyTransaction {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state == txNew {
		// Only allow to set TimestampBound before the first query.
		t.tb = tb
	}
	return t
}

// ReadWriteTransaction provides a locking read-write transaction.
//
// This type of transaction is the only way to write data into Cloud Spanner;
// (*Client).Apply, (*Client).ApplyAtLeastOnce, (*Client).PartitionedUpdate use
// transactions internally. These transactions rely on pessimistic locking and,
// if necessary, two-phase commit. Locking read-write transactions may abort,
// requiring the application to retry. However, the interface exposed by
// (*Client).ReadWriteTransaction eliminates the need for applications to write
// retry loops explicitly.
//
// Locking transactions may be used to atomically read-modify-write data
// anywhere in a database. This type of transaction is externally consistent.
//
// Clients should attempt to minimize the amount of time a transaction is
// active. Faster transactions commit with higher probability and cause less
// contention. Cloud Spanner attempts to keep read locks active as long as the
// transaction continues to do reads.  Long periods of inactivity at the client
// may cause Cloud Spanner to release a transaction's locks and abort it.
//
// Reads performed within a transaction acquire locks on the data being
// read. Writes can only be done at commit time, after all reads have been
// completed. Conceptually, a read-write transaction consists of zero or more
// reads or SQL queries followed by a commit.
//
// See (*Client).ReadWriteTransaction for an example.
//
// Semantics
//
// Cloud Spanner can commit the transaction if all read locks it acquired are
// still valid at commit time, and it is able to acquire write locks for all
// writes. Cloud Spanner can abort the transaction for any reason. If a commit
// attempt returns ABORTED, Cloud Spanner guarantees that the transaction has
// not modified any user data in Cloud Spanner.
//
// Unless the transaction commits, Cloud Spanner makes no guarantees about how
// long the transaction's locks were held for. It is an error to use Cloud
// Spanner locks for any sort of mutual exclusion other than between Cloud
// Spanner transactions themselves.
//
// Aborted transactions
//
// Application code does not need to retry explicitly; RunInTransaction will
// automatically retry a transaction if an attempt results in an abort. The lock
// priority of a transaction increases after each prior aborted transaction,
// meaning that the next attempt has a slightly better chance of success than
// before.
//
// Under some circumstances (e.g., many transactions attempting to modify the
// same row(s)), a transaction can abort many times in a short period before
// successfully committing. Thus, it is not a good idea to cap the number of
// retries a transaction can attempt; instead, it is better to limit the total
// amount of wall time spent retrying.
//
// Idle transactions
//
// A transaction is considered idle if it has no outstanding reads or SQL
// queries and has not started a read or SQL query within the last 10
// seconds. Idle transactions can be aborted by Cloud Spanner so that they don't
// hold on to locks indefinitely. In that case, the commit will fail with error
// ABORTED.
//
// If this behavior is undesirable, periodically executing a simple SQL query
// in the transaction (e.g., SELECT 1) prevents the transaction from becoming
// idle.
type ReadWriteTransaction struct {
	// txReadOnly contains methods for performing transactional reads.
	txReadOnly
	// tx is the transaction ID in Cloud Spanner that uniquely identifies the
	// ReadWriteTransaction. It is set only once in ReadWriteTransaction.begin()
	// during the initialization of ReadWriteTransaction.
	tx transactionID
	// mu protects concurrent access to the internal states of
	// ReadWriteTransaction.
	mu sync.Mutex
	// state is the current transaction status of the read-write transaction.
	state txState
	// wb is the set of buffered mutations waiting to be committed.
	wb []*Mutation
}

// BufferWrite adds a list of mutations to the set of updates that will be
// applied when the transaction is committed. It does not actually apply the
// write until the transaction is committed, so the operation does not block.
// The effects of the write won't be visible to any reads (including reads done
// in the same transaction) until the transaction commits.
//
// See the example for Client.ReadWriteTransaction.
func (t *ReadWriteTransaction) BufferWrite(ms []*Mutation) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state == txClosed {
		return errTxClosed()
	}
	if t.state != txActive {
		return errUnexpectedTxState(t.state)
	}
	t.wb = append(t.wb, ms...)
	return nil
}

// Update executes a DML statement against the database. It returns the number
// of affected rows. Update returns an error if the statement is a query.
// However, the query is executed, and any data read will be validated upon
// commit.
func (t *ReadWriteTransaction) Update(ctx context.Context, stmt Statement) (rowCount int64, err error) {
	mode := sppb.ExecuteSqlRequest_NORMAL
	return t.update(ctx, stmt, QueryOptions{
		Mode:    &mode,
		Options: t.qo.Options,
	})
}

// UpdateWithOptions executes a DML statement against the database. It returns
// the number of affected rows. The sql query execution will be optimized
// based on the given query options.
func (t *ReadWriteTransaction) UpdateWithOptions(ctx context.Context, stmt Statement, opts QueryOptions) (rowCount int64, err error) {
	return t.update(ctx, stmt, t.qo.merge(opts))
}

func (t *ReadWriteTransaction) update(ctx context.Context, stmt Statement, opts QueryOptions) (rowCount int64, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.Update")
	defer func() { trace.EndSpan(ctx, err) }()
	req, sh, err := t.prepareExecuteSQL(ctx, stmt, opts)
	if err != nil {
		return 0, err
	}
	resultSet, err := sh.getClient().ExecuteSql(contextWithOutgoingMetadata(ctx, sh.getMetadata()), req)
	if err != nil {
		return 0, ToSpannerError(err)
	}
	if resultSet.Stats == nil {
		return 0, spannerErrorf(codes.InvalidArgument, "query passed to Update: %q", stmt.SQL)
	}
	return extractRowCount(resultSet.Stats)
}

// BatchUpdate groups one or more DML statements and sends them to Spanner in a
// single RPC. This is an efficient way to execute multiple DML statements.
//
// A slice of counts is returned, where each count represents the number of
// affected rows for the given query at the same index. If an error occurs,
// counts will be returned up to the query that encountered the error.
func (t *ReadWriteTransaction) BatchUpdate(ctx context.Context, stmts []Statement) (_ []int64, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.BatchUpdate")
	defer func() { trace.EndSpan(ctx, err) }()

	sh, ts, err := t.acquire(ctx)
	if err != nil {
		return nil, err
	}
	// Cloud Spanner will return "Session not found" on bad sessions.
	sid := sh.getID()
	if sid == "" {
		// Might happen if transaction is closed in the middle of a API call.
		return nil, errSessionClosed(sh)
	}

	var sppbStmts []*sppb.ExecuteBatchDmlRequest_Statement
	for _, st := range stmts {
		params, paramTypes, err := st.convertParams()
		if err != nil {
			return nil, err
		}
		sppbStmts = append(sppbStmts, &sppb.ExecuteBatchDmlRequest_Statement{
			Sql:        st.SQL,
			Params:     params,
			ParamTypes: paramTypes,
		})
	}

	resp, err := sh.getClient().ExecuteBatchDml(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.ExecuteBatchDmlRequest{
		Session:     sh.getID(),
		Transaction: ts,
		Statements:  sppbStmts,
		Seqno:       atomic.AddInt64(&t.sequenceNumber, 1),
	})
	if err != nil {
		return nil, ToSpannerError(err)
	}

	var counts []int64
	for _, rs := range resp.ResultSets {
		count, err := extractRowCount(rs.Stats)
		if err != nil {
			return nil, err
		}
		counts = append(counts, count)
	}
	if resp.Status != nil && resp.Status.Code != 0 {
		return counts, spannerErrorf(codes.Code(uint32(resp.Status.Code)), resp.Status.Message)
	}
	return counts, nil
}

// acquire implements txReadEnv.acquire.
func (t *ReadWriteTransaction) acquire(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	ts := &sppb.TransactionSelector{
		Selector: &sppb.TransactionSelector_Id{
			Id: t.tx,
		},
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	switch t.state {
	case txClosed:
		return nil, nil, errTxClosed()
	case txActive:
		return t.sh, ts, nil
	}
	return nil, nil, errUnexpectedTxState(t.state)
}

// release implements txReadEnv.release.
func (t *ReadWriteTransaction) release(err error) {
	t.mu.Lock()
	sh := t.sh
	t.mu.Unlock()
	if sh != nil && isSessionNotFoundError(err) {
		sh.destroy()
	}
}

func beginTransaction(ctx context.Context, sid string, client *vkit.Client) (transactionID, error) {
	res, err := client.BeginTransaction(ctx, &sppb.BeginTransactionRequest{
		Session: sid,
		Options: &sppb.TransactionOptions{
			Mode: &sppb.TransactionOptions_ReadWrite_{
				ReadWrite: &sppb.TransactionOptions_ReadWrite{},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if res.Id == nil {
		return nil, spannerErrorf(codes.Unknown, "BeginTransaction returned a transaction with a nil ID.")
	}
	return res.Id, nil
}

// begin starts a read-write transacton on Cloud Spanner, it is always called
// before any of the public APIs.
func (t *ReadWriteTransaction) begin(ctx context.Context) error {
	if t.tx != nil {
		t.state = txActive
		return nil
	}
	tx, err := beginTransaction(contextWithOutgoingMetadata(ctx, t.sh.getMetadata()), t.sh.getID(), t.sh.getClient())
	if err == nil {
		t.tx = tx
		t.state = txActive
		return nil
	}
	if isSessionNotFoundError(err) {
		t.sh.destroy()
	}
	return err
}

// CommitResponse provides a response of a transaction commit in a database.
type CommitResponse struct {
	// CommitTs is the commit time for a transaction.
	CommitTs time.Time
}

// commit tries to commit a readwrite transaction to Cloud Spanner. It also
// returns the commit response for the transactions.
func (t *ReadWriteTransaction) commit(ctx context.Context) (CommitResponse, error) {
	resp := CommitResponse{}
	t.mu.Lock()
	t.state = txClosed // No further operations after commit.
	mPb, err := mutationsProto(t.wb)
	t.mu.Unlock()
	if err != nil {
		return resp, err
	}
	// In case that sessionHandle was destroyed but transaction body fails to
	// report it.
	sid, client := t.sh.getID(), t.sh.getClient()
	if sid == "" || client == nil {
		return resp, errSessionClosed(t.sh)
	}

	res, e := client.Commit(contextWithOutgoingMetadata(ctx, t.sh.getMetadata()), &sppb.CommitRequest{
		Session: sid,
		Transaction: &sppb.CommitRequest_TransactionId{
			TransactionId: t.tx,
		},
		Mutations: mPb,
	})
	if e != nil {
		return resp, toSpannerErrorWithCommitInfo(e, true)
	}
	if tstamp := res.GetCommitTimestamp(); tstamp != nil {
		resp.CommitTs = time.Unix(tstamp.Seconds, int64(tstamp.Nanos))
	}
	if isSessionNotFoundError(err) {
		t.sh.destroy()
	}
	return resp, err
}

// rollback is called when a commit is aborted or the transaction body runs
// into error.
func (t *ReadWriteTransaction) rollback(ctx context.Context) {
	t.mu.Lock()
	// Forbid further operations on rollbacked transaction.
	t.state = txClosed
	t.mu.Unlock()
	// In case that sessionHandle was destroyed but transaction body fails to
	// report it.
	sid, client := t.sh.getID(), t.sh.getClient()
	if sid == "" || client == nil {
		return
	}
	err := client.Rollback(contextWithOutgoingMetadata(ctx, t.sh.getMetadata()), &sppb.RollbackRequest{
		Session:       sid,
		TransactionId: t.tx,
	})
	if isSessionNotFoundError(err) {
		t.sh.destroy()
	}
}

// runInTransaction executes f under a read-write transaction context.
func (t *ReadWriteTransaction) runInTransaction(ctx context.Context, f func(context.Context, *ReadWriteTransaction) error) (CommitResponse, error) {
	var (
		resp            CommitResponse
		err             error
		errDuringCommit bool
	)
	if err = f(context.WithValue(ctx, transactionInProgressKey{}, 1), t); err == nil {
		// Try to commit if transaction body returns no error.
		resp, err = t.commit(ctx)
		errDuringCommit = err != nil
	}
	if err != nil {
		if isAbortedErr(err) {
			// Retry the transaction using the same session on ABORT error.
			// Cloud Spanner will create the new transaction with the previous
			// one's wound-wait priority.
			return resp, err
		}
		if isSessionNotFoundError(err) {
			t.sh.destroy()
			return resp, err
		}
		// Rollback the transaction unless the error occurred during the
		// commit. Executing a rollback after a commit has failed will
		// otherwise cause an error. Note that transient errors, such as
		// UNAVAILABLE, are already handled in the gRPC layer and do not show
		// up here. Context errors (deadline exceeded / canceled) during
		// commits are also not rolled back.
		if !errDuringCommit {
			t.rollback(ctx)
		}
		return resp, err
	}
	// err == nil, return commit response.
	return resp, nil
}

// ReadWriteStmtBasedTransaction provides a wrapper of ReadWriteTransaction in
// order to run a read-write transaction in a statement-based way.
//
// This struct is returned by NewReadWriteStmtBasedTransaction and contains
// Commit() and Rollback() methods to end a transaction.
type ReadWriteStmtBasedTransaction struct {
	// ReadWriteTransaction contains methods for performing transactional reads.
	ReadWriteTransaction
}

// NewReadWriteStmtBasedTransaction starts a read-write transaction. Commit() or
// Rollback() must be called to end a transaction. If Commit() or Rollback() is
// not called, the session that is used by the transaction will not be returned
// to the pool and cause a session leak.
//
// This method should only be used when manual error handling and retry
// management is needed. Cloud Spanner may abort a read/write transaction at any
// moment, and each statement that is executed on the transaction should be
// checked for an Aborted error, including queries and read operations.
//
// For most use cases, client.ReadWriteTransaction should be used, as it will
// handle all Aborted and 'Session not found' errors automatically.
func NewReadWriteStmtBasedTransaction(ctx context.Context, c *Client) (*ReadWriteStmtBasedTransaction, error) {
	return NewReadWriteStmtBasedTransactionWithOptions(ctx, c, TransactionOptions{})
}

// NewReadWriteStmtBasedTransactionWithOptions starts a read-write transaction
// with configurable options. Commit() or Rollback() must be called to end a
// transaction. If Commit() or Rollback() is not called, the session that is
// used by the transaction will not be returned to the pool and cause a session
// leak.
//
// NewReadWriteStmtBasedTransactionWithOptions is a configurable version of
// NewReadWriteStmtBasedTransaction.
func NewReadWriteStmtBasedTransactionWithOptions(ctx context.Context, c *Client, options TransactionOptions) (*ReadWriteStmtBasedTransaction, error) {
	var (
		sh  *sessionHandle
		err error
		t   *ReadWriteStmtBasedTransaction
	)
	sh, err = c.idleSessions.takeWriteSession(ctx)
	if err != nil {
		// If session retrieval fails, just fail the transaction.
		return nil, err
	}
	t = &ReadWriteStmtBasedTransaction{
		ReadWriteTransaction: ReadWriteTransaction{
			tx: sh.getTransactionID(),
		},
	}
	t.txReadOnly.sh = sh
	t.txReadOnly.txReadEnv = t
	t.txReadOnly.qo = c.qo

	if err = t.begin(ctx); err != nil {
		if sh != nil {
			sh.recycle()
		}
		return nil, err
	}
	return t, err
}

// Commit tries to commit a readwrite transaction to Cloud Spanner. It also
// returns the commit timestamp for the transactions.
func (t *ReadWriteStmtBasedTransaction) Commit(ctx context.Context) (time.Time, error) {
	resp, err := t.commit(ctx)
	// Rolling back an aborted transaction is not necessary.
	if err != nil && status.Code(err) != codes.Aborted {
		t.rollback(ctx)
	}
	if t.sh != nil {
		t.sh.recycle()
	}
	return resp.CommitTs, err
}

// Rollback is called to cancel the ongoing transaction that has not been
// committed yet.
func (t *ReadWriteStmtBasedTransaction) Rollback(ctx context.Context) {
	t.rollback(ctx)
	if t.sh != nil {
		t.sh.recycle()
	}
}

// writeOnlyTransaction provides the most efficient way of doing write-only
// transactions. It essentially does blind writes to Cloud Spanner.
type writeOnlyTransaction struct {
	// sp is the session pool which writeOnlyTransaction uses to get Cloud
	// Spanner sessions for blind writes.
	sp *sessionPool
}

// applyAtLeastOnce commits a list of mutations to Cloud Spanner at least once,
// unless one of the following happens:
//
//     1) Context times out.
//     2) An unretryable error (e.g. database not found) occurs.
//     3) There is a malformed Mutation object.
func (t *writeOnlyTransaction) applyAtLeastOnce(ctx context.Context, ms ...*Mutation) (time.Time, error) {
	var (
		ts time.Time
		sh *sessionHandle
	)
	mPb, err := mutationsProto(ms)
	if err != nil {
		// Malformed mutation found, just return the error.
		return ts, err
	}

	// Retry-loop for aborted transactions.
	// TODO: Replace with generic retryer.
	for {
		if sh == nil || sh.getID() == "" || sh.getClient() == nil {
			// No usable session for doing the commit, take one from pool.
			sh, err = t.sp.take(ctx)
			if err != nil {
				// sessionPool.Take already retries for session
				// creations/retrivals.
				return ts, err
			}
			defer sh.recycle()
		}
		res, err := sh.getClient().Commit(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.CommitRequest{
			Session: sh.getID(),
			Transaction: &sppb.CommitRequest_SingleUseTransaction{
				SingleUseTransaction: &sppb.TransactionOptions{
					Mode: &sppb.TransactionOptions_ReadWrite_{
						ReadWrite: &sppb.TransactionOptions_ReadWrite{},
					},
				},
			},
			Mutations: mPb,
		})
		if err != nil && !isAbortedErr(err) {
			if isSessionNotFoundError(err) {
				// Discard the bad session.
				sh.destroy()
			}
			return ts, toSpannerErrorWithCommitInfo(err, true)
		} else if err == nil {
			if tstamp := res.GetCommitTimestamp(); tstamp != nil {
				ts = time.Unix(tstamp.Seconds, int64(tstamp.Nanos))
			}
			break
		}
	}
	return ts, ToSpannerError(err)
}

// isAbortedErr returns true if the error indicates that an gRPC call is
// aborted on the server side.
func isAbortedErr(err error) bool {
	if err == nil {
		return false
	}
	if ErrCode(err) == codes.Aborted {
		return true
	}
	return false
}
