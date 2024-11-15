[![GoDoc](https://godoc.org/github.com/couchbase/gocb?status.png)](https://pkg.go.dev/github.com/couchbase/gocb)

# Couchbase Go Client

The Go SDK library allows you to connect to a Couchbase cluster from Go. 
It is written in pure Go, and uses the included gocbcore library to handle communicating to the cluster over the Couchbase binary protocol.


## Useful Links

### Source
The project source is hosted at [https://github.com/couchbase/gocb](https://github.com/couchbase/gocb).

### Documentation
You can explore our API reference through godoc at [https://pkg.go.dev/github.com/couchbase/gocb](https://pkg.go.dev/github.com/couchbase/gocb).

You can also find documentation for the Go SDK on the [official Couchbase docs](https://docs.couchbase.com/go-sdk/current/hello-world/overview.html).

### Bug Tracker
Issues are tracked on Couchbase's public [issues.couchbase.com](http://www.couchbase.com/issues/browse/GOCBC).
Contact [the site admins](https://issues.couchbase.com/secure/ContactAdministrators!default.jspa) regarding login or other problems at issues.couchbase.com (officially) or ask around [on the forum](https://forums.couchbase.com/) (unofficially).

### Discussion
You can chat with us on [Discord](https://discord.com/invite/sQ5qbPZuTh) or the [official Couchbase forums](https://forums.couchbase.com/c/go-sdk/23).

## Installing

To install the latest stable version, run:
```bash
go get github.com/couchbase/gocb/v2@latest
```

To install the latest developer version, run:
```bash
go get github.com/couchbase/gocb/v2@master
```

## Testing

You can run tests in the usual Go way:

`go test -race ./...`

Which will execute both the unit test suite and the integration test suite.
By default, the integration test suite is run against a mock Couchbase Server.
See the `testmain_test.go` file for information on command line arguments for running tests against a real server instance.

## Release train

Releases are targeted for every third Tuesday of the month.
This is subject to change based on priorities.

## Linting

Linting is performed used `golangci-lint`.
To run:

`make lint`

## License
Copyright 2016 Couchbase Inc.

Licensed under the Apache License, Version 2.0.

See
[LICENSE](https://github.com/couchbase/gocb/blob/master/LICENSE)
for further details.
