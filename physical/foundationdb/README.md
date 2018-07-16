# FoundationDB storage backend

Extra steps are required to produce a Vault build containing the FoundationDB
backend; attempts to use the backend on a build produced without following
this procedure will fail with a descriptive error message at runtime.

## Installing the Go bindings

You will need to install the FoundationDB Go bindings to build the FoundationDB
backend. Make sure you have the FoundationDB client library installed on your
system, along with Mono (core is enough), then install the Go bindings using
the `fdb-go-install.sh` script:

```
$ physical/foundationdb/fdb-go-install.sh
```

## Building Vault

To build Vault the FoundationDB backend, add FDB_ENABLED=1 when invoking
`make`, e.g.

```
$ make dev FDB_ENABLED=1
```

## Running tests

Similarly, add FDB_ENABLED=1 to your `make` invocation when running tests,
e.g.

```
$ make test TEST=./physical/foundationdb FDB_ENABLED=1
```
