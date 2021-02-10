/*
 * doc.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2013-2018 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// FoundationDB Go API

/*
Package fdb provides an interface to FoundationDB databases (version 2.0 or higher).

To build and run programs using this package, you must have an installed copy of
the FoundationDB client libraries (version 2.0.0 or later), available for Linux,
Windows and OS X at https://www.foundationdb.org/download/.

This documentation specifically applies to the FoundationDB Go binding. For more
extensive guidance to programming with FoundationDB, as well as API
documentation for the other FoundationDB interfaces, please see
https://apple.github.io/foundationdb/index.html.

Basic Usage

A basic interaction with the FoundationDB API is demonstrated below:

    package main

    import (
        "github.com/apple/foundationdb/bindings/go/src/fdb"
        "log"
        "fmt"
    )

    func main() {
        // Different API versions may expose different runtime behaviors.
        fdb.MustAPIVersion(700)

        // Open the default database from the system cluster
        db := fdb.MustOpenDefault()

        // Database reads and writes happen inside transactions
        ret, e := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
            tr.Set(fdb.Key("hello"), []byte("world"))
            return tr.Get(fdb.Key("foo")).MustGet(), nil
            // db.Transact automatically commits (and if necessary,
            // retries) the transaction
        })
        if e != nil {
            log.Fatalf("Unable to perform FDB transaction (%v)", e)
        }

        fmt.Printf("hello is now world, foo was: %s\n", string(ret.([]byte)))
    }

Futures

Many functions in this package are asynchronous and return Future objects. A
Future represents a value (or error) to be available at some later
time. Functions documented as blocking on a Future will block the calling
goroutine until the Future is ready (although if the Future is already ready,
the call will not block at all). While a goroutine is blocked on a Future, other
goroutines are free to execute and interact with the FoundationDB API.

It is possible (and often recommended) to call several asynchronous operations
and have multiple Future objects outstanding inside a single goroutine. All
operations will execute in parallel, and the calling goroutine will not block
until a blocking method on any one of the Futures is called.

On Panics

Idiomatic Go code strongly frowns at panics that escape library/package
boundaries, in favor of explicitly returned errors. Idiomatic FoundationDB
client programs, however, are built around the idea of retryable
programmer-provided transactional functions. Retryable transactions can be
implemented using only error values:

    ret, e := db.Transact(func (tr Transaction) (interface{}, error) {
        // FoundationDB futures represent a value that will become available
        futureValueOne := tr.Get(fdb.Key("foo"))
        futureValueTwo := tr.Get(fdb.Key("bar"))

        // Both reads are being carried out in parallel

        // Get the first value (or any error)
        valueOne, e := futureValueOne.Get()
        if e != nil {
            return nil, e
        }

        // Get the second value (or any error)
        valueTwo, e := futureValueTwo.Get()
        if e != nil {
            return nil, e
        }

        // Return the two values
        return []string{valueOne, valueTwo}, nil
    })

If either read encounters an error, it will be returned to Transact, which will
determine if the error is retryable or not (using (Transaction).OnError). If the
error is an FDB Error and retryable (such as a conflict with with another
transaction), then the programmer-provided function will be run again. If the
error is fatal (or not an FDB Error), then the error will be returned to the
caller of Transact.

In practice, checking for an error from every asynchronous future type in the
FoundationDB API quickly becomes frustrating. As a convenience, every Future
type also has a MustGet method, which returns the same type and value as Get,
but exposes FoundationDB Errors via a panic rather than an explicitly returned
error. The above example may be rewritten as:

    ret, e := db.Transact(func (tr Transaction) (interface{}, error) {
        // FoundationDB futures represent a value that will become available
        futureValueOne := tr.Get(fdb.Key("foo"))
        futureValueTwo := tr.Get(fdb.Key("bar"))

        // Both reads are being carried out in parallel

        // Get the first value
        valueOne := futureValueOne.MustGet()
        // Get the second value
        valueTwo := futureValueTwo.MustGet()

        // Return the two values
        return []string{valueOne, valueTwo}, nil
    })

MustGet returns nil (which is different from empty slice []byte{}), when the
key doesn't exist, and hence non-existence can be checked as follows:

    val := tr.Get(fdb.Key("foobar")).MustGet()
    if val == nil {
      fmt.Println("foobar does not exist.")
    } else {
      fmt.Println("foobar exists.")
    }

Any panic that occurs during execution of the caller-provided function will be
recovered by the (Database).Transact method. If the error is an FDB Error, it
will either result in a retry of the function or be returned by Transact. If the
error is any other type (panics from code other than MustGet), Transact will
re-panic the original value.

Note that (Transaction).Transact also recovers panics, but does not itself
retry. If the recovered value is an FDB Error, it will be returned to the caller
of (Transaction).Transact; all other values will be re-panicked.

Transactions and Goroutines

When using a Transactor in the fdb package, particular care must be taken if
goroutines are created inside of the function passed to the Transact method. Any
panic from the goroutine will not be recovered by Transact, and (unless
otherwise recovered) will result in the termination of that goroutine.

Furthermore, any errors returned or panicked by fdb methods called in the
goroutine must be safely returned to the function passed to Transact, and either
returned or panicked, to allow Transact to appropriately retry or terminate the
transactional function.

Lastly, a transactional function may be retried indefinitely. It is advisable to
make sure any goroutines created during the transactional function have
completed before returning from the transactional function, or a potentially
unbounded number of goroutines may be created.

Given these complexities, it is generally best practice to use a single
goroutine for each logical thread of interaction with FoundationDB, and allow
each goroutine to block when necessary to wait for Futures to become ready.

Streaming Modes

When using GetRange methods in the FoundationDB API, clients can request large
ranges of the database to iterate over. Making such a request doesn't
necessarily mean that the client will consume all of the data in the range --
sometimes the client doesn't know how far it intends to iterate in
advance. FoundationDB tries to balance latency and bandwidth by requesting data
for iteration in batches.

The Mode field of the RangeOptions struct allows a client to customize this
performance tradeoff by providing extra information about how the iterator will
be used.

The default value of Mode is StreamingModeIterator, which tries to provide a
reasonable default balance. Other streaming modes that prioritize throughput or
latency are available -- see the documented StreamingMode values for specific
options.

Atomic Operations

The FDB package provides a number of atomic operations on the Database and
Transaction objects. An atomic operation is a single database command that
carries out several logical steps: reading the value of a key, performing a
transformation on that value, and writing the result. Different atomic
operations perform different transformations. Like other database operations, an
atomic operation is used within a transaction.

For more information on atomic operations in FoundationDB, please see
https://apple.github.io/foundationdb/developer-guide.html#atomic-operations. The
operands to atomic operations in this API must be provided as appropriately
encoded byte slices. To convert a Go type to a byte slice, see the binary
package.

The current atomic operations in this API are Add, BitAnd, BitOr, BitXor,
CompareAndClear, Max, Min, SetVersionstampedKey, SetVersionstampedValue
(all methods on Transaction).
*/
package fdb
