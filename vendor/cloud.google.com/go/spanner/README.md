## Cloud Spanner [![GoDoc](https://godoc.org/cloud.google.com/go/spanner?status.svg)](https://godoc.org/cloud.google.com/go/spanner)

- [About Cloud Spanner](https://cloud.google.com/spanner/)
- [API documentation](https://cloud.google.com/spanner/docs)
- [Go client documentation](https://godoc.org/cloud.google.com/go/spanner)

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