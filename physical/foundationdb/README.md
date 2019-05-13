# FoundationDB storage backend

Extra steps are required to produce a Vault build containing the FoundationDB
backend; attempts to use the backend on a build produced without following
this procedure will fail with a descriptive error message at runtime.

## Installing the Go bindings

### Picking a version

The version of the Go bindings and the FoundationDB client library used to
build them must match.

This version will determine the minimum API version that can be used, hence
it should be no higher than the version of FoundationDB used in your cluster,
and must also satisfy the requirements of the backend code.

The minimum required API version for the FoundationDB backend is 520.

### Installation

Make sure you have Mono installed (core is enough), then install the
Go bindings using the `fdb-go-install.sh` script:

```
$ physical/foundationdb/fdb-go-install.sh install --fdbver x.y.z
```

By default, if `--fdbver x.y.z` is not specified, version 5.2.4 will be used.

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
