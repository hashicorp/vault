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

// Package crdb provides helpers for using CockroachDB in client
// applications.
package crdb

import (
	"context"
	"database/sql"
	"errors"
)

// Execute runs fn and retries it as needed. It is used to add retry handling to
// the execution of a single statement. If a multi-statement transaction is
// being run, use ExecuteTx instead.
//
// Retry handling for individual statements (implicit transactions) is usually
// performed automatically on the CockroachDB SQL gateway. As such, use of this
// function is generally not necessary. The exception to this rule is that
// automatic retries for individual statements are disabled once CockroachDB
// begins streaming results for the statements back to the client. By default,
// result streaming does not begin until the size of the result being produced
// for the client, including protocol overhead, exceeds 16KiB. As long as the
// results of a single statement or batch of statements are known to stay clear
// of this limit, the client does not need to worry about retries and should not
// need to use this function.
//
// For more information about automatic transaction retries in CockroachDB, see
// https://cockroachlabs.com/docs/stable/transactions.html#automatic-retries.
//
// NOTE: the supplied fn closure should not have external side effects beyond
// changes to the database.
//
// fn must take care when wrapping errors returned from the database driver with
// additional context. For example, if the SELECT statement fails in the
// following snippet, the original retryable error will be masked by the call to
// fmt.Errorf, and the transaction will not be automatically retried.
//
//    crdb.Execute(func () error {
//        rows, err := db.QueryContext(ctx, "SELECT ...")
//        if err != nil {
//            return fmt.Errorf("scanning row: %s", err)
//        }
//        defer rows.Close()
//        for rows.Next() {
//            // ...
//        }
//        if err := rows.Err(); err != nil {
//            return fmt.Errorf("scanning row: %s", err)
//        }
//        return nil
//    })
//
// Instead, add context by returning an error that implements either:
// - a `Cause() error` method, in the manner of github.com/pkg/errors, or
// - an `Unwrap() error` method, in the manner of the Go 1.13 standard
// library.
//
// To achieve this, you can implement your own error type, or use
// `errors.Wrap()` from github.com/pkg/errors or
// github.com/cockroachdb/errors or a similar package, or use go
// 1.13's special `%w` formatter with fmt.Errorf(), for example
// fmt.Errorf("scanning row: %w", err).
//
//    import "github.com/pkg/errors"
//
//    crdb.Execute(func () error {
//        rows, err := db.QueryContext(ctx, "SELECT ...")
//        if err != nil {
//            return errors.Wrap(err, "scanning row")
//        }
//        defer rows.Close()
//        for rows.Next() {
//            // ...
//        }
//        if err := rows.Err(); err != nil {
//            return errors.Wrap(err, "scanning row")
//        }
//        return nil
//    })
//
func Execute(fn func() error) (err error) {
	for {
		err = fn()
		if err == nil || !errIsRetryable(err) {
			return err
		}
	}
}

type txConfigKey struct{}

// WithMaxRetries configures context so that ExecuteTx retries tx specified
// number of times when encountering retryable errors.
// Setting retries to 0 will retry indefinitely.
func WithMaxRetries(ctx context.Context, retries int) context.Context {
	return context.WithValue(ctx, txConfigKey{}, retries)
}

const defaultRetries = 50

func numRetriesFromContext(ctx context.Context) int {
	if v := ctx.Value(txConfigKey{}); v != nil {
		if retries, ok := v.(int); ok && retries >= 0 {
			return retries
		}
	}
	return defaultRetries
}

// ExecuteTx runs fn inside a transaction and retries it as needed. On
// non-retryable failures, the transaction is aborted and rolled back; on
// success, the transaction is committed.
//
// There are cases where the state of a transaction is inherently ambiguous: if
// we err on RELEASE with a communication error it's unclear if the transaction
// has been committed or not (similar to erroring on COMMIT in other databases).
// In that case, we return AmbiguousCommitError.
//
// There are cases when restarting a transaction fails: we err on ROLLBACK to
// the SAVEPOINT. In that case, we return a TxnRestartError.
//
// For more information about CockroachDB's transaction model, see
// https://cockroachlabs.com/docs/stable/transactions.html.
//
// NOTE: the supplied fn closure should not have external side effects beyond
// changes to the database.
//
// fn must take care when wrapping errors returned from the database driver with
// additional context. For example, if the UPDATE statement fails in the
// following snippet, the original retryable error will be masked by the call to
// fmt.Errorf, and the transaction will not be automatically retried.
//
//    crdb.ExecuteTx(ctx, db, txopts, func (tx *sql.Tx) error {
//        if err := tx.ExecContext(ctx, "UPDATE..."); err != nil {
//            return fmt.Errorf("updating record: %s", err)
//        }
//        return nil
//    })
//
// Instead, add context by returning an error that implements either:
// - a `Cause() error` method, in the manner of github.com/pkg/errors, or
// - an `Unwrap() error` method, in the manner of the Go 1.13 standard
// library.
//
// To achieve this, you can implement your own error type, or use
// `errors.Wrap()` from github.com/pkg/errors or
// github.com/cockroachdb/errors or a similar package, or use go
// 1.13's special `%w` formatter with fmt.Errorf(), for example
// fmt.Errorf("scanning row: %w", err).
//
//    import "github.com/pkg/errors"
//
//    crdb.ExecuteTx(ctx, db, txopts, func (tx *sql.Tx) error {
//        if err := tx.ExecContext(ctx, "UPDATE..."); err != nil {
//            return errors.Wrap(err, "updating record")
//        }
//        return nil
//    })
//
func ExecuteTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, fn func(*sql.Tx) error) error {
	// Start a transaction.
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	return ExecuteInTx(ctx, stdlibTxnAdapter{tx}, func() error { return fn(tx) })
}

type stdlibTxnAdapter struct {
	tx *sql.Tx
}

var _ Tx = stdlibTxnAdapter{}

// Exec is part of the tx interface.
func (tx stdlibTxnAdapter) Exec(ctx context.Context, q string, args ...interface{}) error {
	_, err := tx.tx.ExecContext(ctx, q, args...)
	return err
}

// Commit is part of the tx interface.
func (tx stdlibTxnAdapter) Commit(context.Context) error {
	return tx.tx.Commit()
}

// Commit is part of the tx interface.
func (tx stdlibTxnAdapter) Rollback(context.Context) error {
	return tx.tx.Rollback()
}

func errIsRetryable(err error) bool {
	// We look for either:
	//  - the standard PG errcode SerializationFailureError:40001 or
	//  - the Cockroach extension errcode RetriableError:CR000. This extension
	//    has been removed server-side, but support for it has been left here for
	//    now to maintain backwards compatibility.
	code := errCode(err)
	return code == "CR000" || code == "40001"
}

func errCode(err error) string {
	var sqlErr errWithSQLState
	if errors.As(err, &sqlErr) {
		return sqlErr.SQLState()
	}

	return ""
}

// errWithSQLState is implemented by pgx (pgconn.PgError) and lib/pq
type errWithSQLState interface {
	SQLState() string
}
