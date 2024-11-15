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

/*
Package spanner provides a client for reading and writing to Cloud Spanner
databases. See the packages under admin for clients that operate on databases
and instances.

See https://cloud.google.com/spanner/docs/getting-started/go/ for an
introduction to Cloud Spanner and additional help on using this API.

See https://godoc.org/cloud.google.com/go for authentication, timeouts,
connection pooling and similar aspects of this package.

# Creating a Client

To start working with this package, create a client that refers to the database
of interest:

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, "projects/P/instances/I/databases/D")
	if err != nil {
	    // TODO: Handle error.
	}
	defer client.Close()

Remember to close the client after use to free up the sessions in the session
pool.

To use an emulator with this library, you can set the SPANNER_EMULATOR_HOST
environment variable to the address at which your emulator is running. This will
send requests to that address instead of to Cloud Spanner. You can then create
and use a client as usual:

	// Set SPANNER_EMULATOR_HOST environment variable.
	err := os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")
	if err != nil {
	    // TODO: Handle error.
	}
	// Create client as usual.
	client, err := spanner.NewClient(ctx, "projects/P/instances/I/databases/D")
	if err != nil {
	    // TODO: Handle error.
	}

# Simple Reads and Writes

Two Client methods, Apply and Single, work well for simple reads and writes. As
a quick introduction, here we write a new row to the database and read it back:

	_, err := client.Apply(ctx, []*spanner.Mutation{
	    spanner.Insert("Users",
	        []string{"name", "email"},
	        []interface{}{"alice", "a@example.com"})})
	if err != nil {
	    // TODO: Handle error.
	}
	row, err := client.Single().ReadRow(ctx, "Users",
	    spanner.Key{"alice"}, []string{"email"})
	if err != nil {
	    // TODO: Handle error.
	}

All the methods used above are discussed in more detail below.

# Keys

Every Cloud Spanner row has a unique key, composed of one or more columns.
Construct keys with a literal of type Key:

	key1 := spanner.Key{"alice"}

# KeyRanges

The keys of a Cloud Spanner table are ordered. You can specify ranges of keys
using the KeyRange type:

	kr1 := spanner.KeyRange{Start: key1, End: key2}

By default, a KeyRange includes its start key but not its end key. Use
the Kind field to specify other boundary conditions:

	// include both keys
	kr2 := spanner.KeyRange{Start: key1, End: key2, Kind: spanner.ClosedClosed}

# KeySets

A KeySet represents a set of keys. A single Key or KeyRange can act as a KeySet.
Use the KeySets function to build the union of several KeySets:

	ks1 := spanner.KeySets(key1, key2, kr1, kr2)

AllKeys returns a KeySet that refers to all the keys in a table:

	ks2 := spanner.AllKeys()

# Transactions

All Cloud Spanner reads and writes occur inside transactions. There are two
types of transactions, read-only and read-write. Read-only transactions cannot
change the database, do not acquire locks, and may access either the current
database state or states in the past. Read-write transactions can read the
database before writing to it, and always apply to the most recent database
state.

# Single Reads

The simplest and fastest transaction is a ReadOnlyTransaction that supports a
single read operation. Use Client.Single to create such a transaction. You can
chain the call to Single with a call to a Read method.

When you only want one row whose key you know, use ReadRow. Provide the table
name, key, and the columns you want to read:

	row, err := client.Single().ReadRow(ctx, "Accounts", spanner.Key{"alice"}, []string{"balance"})

Read multiple rows with the Read method. It takes a table name, KeySet, and list
of columns:

	iter := client.Single().Read(ctx, "Accounts", keyset1, columns)

Read returns a RowIterator. You can call the Do method on the iterator and pass
a callback:

	err := iter.Do(func(row *Row) error {
	   // TODO: use row
	   return nil
	})

RowIterator also follows the standard pattern for the Google
Cloud Client Libraries:

	defer iter.Stop()
	for {
	    row, err := iter.Next()
	    if err == iterator.Done {
	        break
	    }
	    if err != nil {
	        // TODO: Handle error.
	    }
	    // TODO: use row
	}

Always call Stop when you finish using an iterator this way, whether or not you
iterate to the end. (Failing to call Stop could lead you to exhaust the
database's session quota.)

To read rows with an index, use ReadUsingIndex.

# Statements

The most general form of reading uses SQL statements. Construct a Statement
with NewStatement, setting any parameters using the Statement's Params map:

	stmt := spanner.NewStatement("SELECT First, Last FROM SINGERS WHERE Last >= @start")
	stmt.Params["start"] = "Dylan"

You can also construct a Statement directly with a struct literal, providing
your own map of parameters.

Use the Query method to run the statement and obtain an iterator:

	iter := client.Single().Query(ctx, stmt)

# Rows

Once you have a Row, via an iterator or a call to ReadRow, you can extract
column values in several ways. Pass in a pointer to a Go variable of the
appropriate type when you extract a value.

You can extract by column position or name:

	err := row.Column(0, &name)
	err = row.ColumnByName("balance", &balance)

You can extract all the columns at once:

	err = row.Columns(&name, &balance)

Or you can define a Go struct that corresponds to your columns, and extract
into that:

	var s struct { Name string; Balance int64 }
	err = row.ToStruct(&s)

For Cloud Spanner columns that may contain NULL, use one of the NullXXX types,
like NullString:

	var ns spanner.NullString
	if err := row.Column(0, &ns); err != nil {
	    // TODO: Handle error.
	}
	if ns.Valid {
	    fmt.Println(ns.StringVal)
	} else {
	    fmt.Println("column is NULL")
	}

# Multiple Reads

To perform more than one read in a transaction, use ReadOnlyTransaction:

	txn := client.ReadOnlyTransaction()
	defer txn.Close()
	iter := txn.Query(ctx, stmt1)
	// ...
	iter =  txn.Query(ctx, stmt2)
	// ...

You must call Close when you are done with the transaction.

# Timestamps and Timestamp Bounds

Cloud Spanner read-only transactions conceptually perform all their reads at a
single moment in time, called the transaction's read timestamp. Once a read has
started, you can call ReadOnlyTransaction's Timestamp method to obtain the read
timestamp.

By default, a transaction will pick the most recent time (a time where all
previously committed transactions are visible) for its reads. This provides the
freshest data, but may involve some delay. You can often get a quicker response
if you are willing to tolerate "stale" data. You can control the read timestamp
selected by a transaction by calling the WithTimestampBound method on the
transaction before using it. For example, to perform a query on data that is at
most one minute stale, use

	client.Single().
	    WithTimestampBound(spanner.MaxStaleness(1*time.Minute)).
	    Query(ctx, stmt)

See the documentation of TimestampBound for more details.

# Mutations

To write values to a Cloud Spanner database, construct a Mutation. The spanner
package has functions for inserting, updating and deleting rows. Except for the
Delete methods, which take a Key or KeyRange, each mutation-building function
comes in three varieties.

One takes lists of columns and values along with the table name:

	m1 := spanner.Insert("Users",
	    []string{"name", "email"},
	    []interface{}{"alice", "a@example.com"})

One takes a map from column names to values:

	m2 := spanner.InsertMap("Users", map[string]interface{}{
	    "name":  "alice",
	    "email": "a@example.com",
	})

And the third accepts a struct value, and determines the columns from the
struct field names:

	type User struct { Name, Email string }
	u := User{Name: "alice", Email: "a@example.com"}
	m3, err := spanner.InsertStruct("Users", u)

# Writes

To apply a list of mutations to the database, use Apply:

	_, err := client.Apply(ctx, []*spanner.Mutation{m1, m2, m3})

If you need to read before writing in a single transaction, use a
ReadWriteTransaction. ReadWriteTransactions may be aborted automatically by the
backend and need to be retried. You pass in a function to ReadWriteTransaction,
and the client will handle the retries automatically. Use the transaction's
BufferWrite method to buffer mutations, which will all be executed at the end
of the transaction:

	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
	    var balance int64
	    row, err := txn.ReadRow(ctx, "Accounts", spanner.Key{"alice"}, []string{"balance"})
	    if err != nil {
	        // The transaction function will be called again if the error code
	        // of this error is Aborted. The backend may automatically abort
	        // any read/write transaction if it detects a deadlock or other
	        // problems.
	        return err
	    }
	    if err := row.Column(0, &balance); err != nil {
	        return err
	    }

	    if balance <= 10 {
	        return errors.New("insufficient funds in account")
	    }
	    balance -= 10
	    m := spanner.Update("Accounts", []string{"user", "balance"}, []interface{}{"alice", balance})
	    // The buffered mutation will be committed.  If the commit
	    // fails with an Aborted error, this function will be called
	    // again.
	    return txn.BufferWrite([]*spanner.Mutation{m})
	})

# Structs

Cloud Spanner STRUCT (aka STRUCT) values
(https://cloud.google.com/spanner/docs/data-types#struct-type) can be
represented by a Go struct value.

A proto StructType is built from the field types and field tag information of
the Go struct. If a field in the struct type definition has a
"spanner:<field_name>" tag, then the value of the "spanner" key in the tag is
used as the name for that field in the built StructType, otherwise the field
name in the struct definition is used. To specify a field with an empty field
name in a Cloud Spanner STRUCT type, use the `spanner:""` tag annotation against
the corresponding field in the Go struct's type definition.

A STRUCT value can contain STRUCT-typed and Array-of-STRUCT typed fields and
these can be specified using named struct-typed and []struct-typed fields inside
a Go struct. However, embedded struct fields are not allowed. Unexported struct
fields are ignored.

NULL STRUCT values in Cloud Spanner are typed. A nil pointer to a Go struct
value can be used to specify a NULL STRUCT value of the corresponding
StructType.  Nil and empty slices of a Go STRUCT type can be used to specify
NULL and empty array values respectively of the corresponding StructType. A
slice of pointers to a Go struct type can be used to specify an array of
NULL-able STRUCT values.

# DML and Partitioned DML

Spanner supports DML statements like INSERT, UPDATE and DELETE. Use
ReadWriteTransaction.Update to run DML statements. It returns the number of rows
affected. (You can call use ReadWriteTransaction.Query with a DML statement. The
first call to Next on the resulting RowIterator will return iterator.Done, and
the RowCount field of the iterator will hold the number of affected rows.)

For large databases, it may be more efficient to partition the DML statement.
Use client.PartitionedUpdate to run a DML statement in this way. Not all DML
statements can be partitioned.

# Tracing

This client has been instrumented to use OpenCensus tracing
(http://opencensus.io). To enable tracing, see "Enabling Tracing for a Program"
at https://godoc.org/go.opencensus.io/trace. OpenCensus tracing requires Go 1.8
or higher.
*/
package spanner // import "cloud.google.com/go/spanner"
