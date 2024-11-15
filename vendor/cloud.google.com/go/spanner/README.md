## Cloud Spanner [![Go Reference](https://pkg.go.dev/badge/cloud.google.com/go/spanner.svg)](https://pkg.go.dev/cloud.google.com/go/spanner)

- [About Cloud Spanner](https://cloud.google.com/spanner/)
- [API documentation](https://cloud.google.com/spanner/docs)
- [Go client documentation](https://pkg.go.dev/cloud.google.com/go/spanner)

### Example Usage

First create a `spanner.Client` to use throughout your application:

[snip]:# (spanner-1)
```go
client, err := spanner.NewClient(ctx, "projects/P/instances/I/databases/D")
if err != nil {
	log.Fatal(err)
}
```

[snip]:# (spanner-2)
```go
// Simple Reads And Writes
_, err = client.Apply(ctx, []*spanner.Mutation{
	spanner.Insert("Users",
		[]string{"name", "email"},
		[]interface{}{"alice", "a@example.com"})})
if err != nil {
	log.Fatal(err)
}
row, err := client.Single().ReadRow(ctx, "Users",
	spanner.Key{"alice"}, []string{"email"})
if err != nil {
	log.Fatal(err)
}
```

### Session Leak
A `Client` object of the Client Library has a limit on the number of maximum sessions. For example the
default value of `MaxOpened`, which is the maximum number of sessions allowed by the session pool in the
Golang Client Library, is 400. You can configure these values at the time of
creating a `Client` by passing custom `SessionPoolConfig` as part of `ClientConfig`. When all the sessions are checked
out of the session pool, every new transaction has to wait until a session is returned to the pool.
If a session is never returned to the pool (hence causing a session leak), the transactions will have to wait
indefinitely and your application will be blocked.

#### Common Root Causes
The most common reason for session leaks in the Golang client library are:
1. Not stopping a `RowIterator` that is returned by `Query`, `Read` and other methods. Always use `RowIterator.Stop()` to ensure that the `RowIterator` is always closed.
2. Not closing a `ReadOnlyTransaction` when you no longer need it. Always call `ReadOnlyTransaction.Close()` after use, to ensure that the `ReadOnlyTransaction` is always closed.

As shown in the example below, the `txn.Close()` statement releases the session after it is complete.
If you fail to call `txn.Close()`, the session is not released back to the pool. The recommended way is to use `defer` as shown below.
```go
client, err := spanner.NewClient(ctx, "projects/P/instances/I/databases/D")
if err != nil {
  log.Fatal(err)
}
txn := client.ReadOnlyTransaction()
defer txn.Close()
```

#### Debugging and Resolving Session Leaks

##### Logging inactive transactions
This option logs warnings when you have exhausted >95% of your session pool. It is enabled by default.
This could mean two things; either you need to increase the max sessions in your session pool (as the number
of queries run using the client side database object is greater than your session pool can serve), or you may
have a session leak. To help debug which transactions may be causing this session leak, the logs will also contain stack traces of
transactions which have been running longer than expected if `TrackSessionHandles` under `SessionPoolConfig` is enabled.

```go
sessionPoolConfig := spanner.SessionPoolConfig{
    TrackSessionHandles: true,
    InactiveTransactionRemovalOptions: spanner.InactiveTransactionRemovalOptions{
      ActionOnInactiveTransaction: spanner.Warn,
    },
}
client, err := spanner.NewClientWithConfig(
	ctx, database, spanner.ClientConfig{SessionPoolConfig: sessionPoolConfig},
)
if err != nil {
	log.Fatal(err)
}
defer client.Close()

// Example Log message to warn presence of long running transactions
// session <session-info> checked out of pool at <session-checkout-time> is long running due to possible session leak for goroutine
// <Stack Trace of transaction>

```

##### Automatically clean inactive transactions
When the option to automatically clean inactive transactions is enabled, the client library will automatically detect
problematic transactions that are running for a very long time (thus causing session leaks) and close them.
The session will be removed from the pool and be replaced by a new session. To dig deeper into which transactions are being
closed, you can check the logs to see the stack trace of the transactions which might be causing these leaks and further
debug them.

```go
sessionPoolConfig := spanner.SessionPoolConfig{
    TrackSessionHandles: true,
    InactiveTransactionRemovalOptions: spanner.InactiveTransactionRemovalOptions{
      ActionOnInactiveTransaction: spanner.WarnAndClose,
    },
}
client, err := spanner.NewClientWithConfig(
	ctx, database, spanner.ClientConfig{SessionPoolConfig: sessionPoolConfig},
)
if err != nil {
log.Fatal(err)
}
defer client.Close()

// Example Log message for when transaction is recycled
// session <session-info> checked out of pool at <session-checkout-time> is long running and will be removed due to possible session leak for goroutine 
// <Stack Trace of transaction>
```