// Copyright 2016 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package crdb

import "context"

// Tx abstracts the operations needed by ExecuteInTx so that different
// frameworks (e.g. go's sql package, pgx, gorm) can be used with ExecuteInTx.
type Tx interface {
	Exec(context.Context, string, ...interface{}) error
	Commit(context.Context) error
	Rollback(context.Context) error
}

// ExecuteInTx runs fn inside tx. This method is primarily intended for internal
// use. See other packages for higher-level, framework-specific ExecuteTx()
// functions.
//
// *WARNING*: It is assumed that no statements have been executed on the
// supplied Tx. ExecuteInTx will only retry statements that are performed within
// the supplied closure (fn). Any statements performed on the tx before
// ExecuteInTx is invoked will *not* be re-run if the transaction needs to be
// retried.
//
// fn is subject to the same restrictions as the fn passed to ExecuteTx.
func ExecuteInTx(ctx context.Context, tx Tx, fn func() error) (err error) {
	defer func() {
		r := recover()

		if r == nil && err == nil {
			// Ignore commit errors. The tx has already been committed by RELEASE.
			_ = tx.Commit(ctx)
			return
		}

		// We always need to execute a Rollback() so sql.DB releases the
		// connection.
		_ = tx.Rollback(ctx)

		if r != nil {
			panic(r)
		}
	}()

	// Specify that we intend to retry this txn in case of CockroachDB retryable
	// errors.
	if err = tx.Exec(ctx, "SAVEPOINT cockroach_restart"); err != nil {
		return err
	}

	maxRetries := numRetriesFromContext(ctx)
	retryCount := 0
	for {
		releaseFailed := false
		err = fn()
		if err == nil {
			// RELEASE acts like COMMIT in CockroachDB. We use it since it gives us an
			// opportunity to react to retryable errors, whereas tx.Commit() doesn't.
			if err = tx.Exec(ctx, "RELEASE SAVEPOINT cockroach_restart"); err == nil {
				return nil
			}
			releaseFailed = true
		}

		// We got an error; let's see if it's a retryable one and, if so, restart.
		if !errIsRetryable(err) {
			if releaseFailed {
				err = newAmbiguousCommitError(err)
			}
			return err
		}

		if rollbackErr := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT cockroach_restart"); rollbackErr != nil {
			return newTxnRestartError(rollbackErr, err)
		}

		retryCount++
		if maxRetries > 0 && retryCount > maxRetries {
			return newMaxRetriesExceededError(err, maxRetries)
		}
	}
}
