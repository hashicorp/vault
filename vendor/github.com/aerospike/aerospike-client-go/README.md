# Aerospike Go Client

[![Aerospike Client Go](https://goreportcard.com/badge/github.com/aerospike/aerospike-client-go)](https://goreportcard.com/report/github.com/aerospike/aerospike-client-go)
[![Build Status](https://travis-ci.org/aerospike/aerospike-client-go.svg?branch=master)](https://travis-ci.org/aerospike/aerospike-client-go)
[![Godoc](https://godoc.org/github.com/aerospike/aerospike-client-go?status.svg)](https://pkg.go.dev/github.com/aerospike/aerospike-client-go)

An Aerospike library for Go.

This library is compatible with Go 1.9+ and supports the following operating systems: Linux, Mac OS X (Windows builds are possible, but untested).

Up-to-date documentation is available in the [![Godoc](https://godoc.org/github.com/aerospike/aerospike-client-go?status.svg)](https://pkg.go.dev/github.com/aerospike/aerospike-client-go).

You can refer to the test files for idiomatic use cases.

Please refer to [`CHANGELOG.md`](CHANGELOG.md) for release notes, or if you encounter breaking changes.

## Notice:

We have released the Go Client v2, with some breaking API changes. Most changes are minor, and can be fixed with relative ease.

The only major issue is that the behavior of the client when a key does not exist has changed. 

It used to return no error, but a `nil` `Record.Bins`. Now it returns `ErrKeyNotFound` error.

This is a significant changes, and you should search your code for all instances of `Bins == nil` and adapt the code accordingly.

Please refer to the [`CHANGELOG.md`](CHANGELOG.md) for details.

- [Usage](#Usage)
- [Prerequisites](#Prerequisites)
- [Installation](#Installation)
- [Tweaking Performance](#Performance)
- [Benchmarks](#Benchmarks)
- [API Documentation](#API-Documentation)
- [Google App Engine](#App-Engine)
- [Reflection](#Reflection)
- [Tests](#Tests)
- [Examples](#Examples)
  - [Tools](#Tools)


## Usage:

The following is a very simple example of CRUD operations in an Aerospike database.

```go
package main

import (
  "fmt"

  aero "github.com/aerospike/aerospike-client-go"
)

// This is only for this example.
// Please handle errors properly.
func panicOnError(err error) {
  if err != nil {
    panic(err)
  }
}

func main() {
  // define a client to connect to
  client, err := aero.NewClient("127.0.0.1", 3000)
  panicOnError(err)

  key, err := aero.NewKey("test", "aerospike", "key")
  panicOnError(err)

  // define some bins with data
  bins := aero.BinMap{
    "bin1": 42,
    "bin2": "An elephant is a mouse with an operating system",
    "bin3": []interface{}{"Go", 2009},
  }

  // write the bins
  err = client.Put(nil, key, bins)
  panicOnError(err)

  // read it back!
  rec, err := client.Get(nil, key)
  panicOnError(err)

  // delete the key, and check if key exists
  existed, err := client.Delete(nil, key)
  panicOnError(err)
  fmt.Printf("Record existed before delete? %v\n", existed)
}
```

More examples illustrating the use of the API are located in the
[`examples`](examples) directory.

Details about the API are available in the [`docs`](docs) directory.

<a name="Prerequisites"></a>
## Prerequisites

[Go](http://golang.org) version v1.12+ is required.

To install the latest stable version of Go, visit
[http://golang.org/dl/](http://golang.org/dl/)


Aerospike Go client implements the wire protocol, and does not depend on the C client.
It is goroutine friendly, and works asynchronously.

Supported operating systems:

- Major Linux distributions (Ubuntu, Debian, Red Hat)
- Mac OS X
- Windows (untested)

<a name="Installation"></a>
## Installation:

1. Install Go 1.9+ and setup your environment as [Documented](http://golang.org/doc/code.html#GOPATH) here.
2. Get the client in your ```GOPATH``` : ```go get github.com/aerospike/aerospike-client-go```
  * To update the client library: ```go get -u github.com/aerospike/aerospike-client-go```

Using [gopkg.in](https://gopkg.in/) is also supported: `go get -u gopkg.in/aerospike/aerospike-client-go.v1`

### Some Hints:

 * To run a go program directly: ```go run <filename.go>```
 * to build:  ```go build -o <output> <filename.go>```
  * example: ```go build -o benchmark tools/benchmark/benchmark.go```

<a name="Performance"></a>
## Performance Tweaking

We are bending all efforts to improve the client's performance. In our reference benchmarks, Go client performs almost as good as the C client.

To read about performance variables, please refer to [`docs/performance.md`](docs/performance.md)

<a name="Tests"></a>
## Tests

This library is packaged with a number of tests. Tests require Ginkgo and Gomega library.

Before running the tests, you need to update the dependencies:

    $ go get .

To run all the test cases with race detection:

    $ ginkgo -r -race


<a name="Examples"></a>
## Examples

A variety of example applications are provided in the [`examples`](examples) directory.

<a name="Tools"></a>
### Tools

A variety of clones of original tools are provided in the [`tools`](tools) directory.
They show how to use more advanced features of the library to re-implement the same functionality in a more concise way.

<a name="Benchmarks"></a>
## Benchmarks

Benchmark utility is provided in the [`tools/benchmark`](tools/benchmark) directory.
See the [`tools/benchmark/README.md`](tools/benchmark/README.md) for details.

<a name="API-Documentation"></a>
## API Documentation

A simple API documentation is available in the [`docs`](docs/README.md) directory. The latest up-to-date docs can be found in [![Godoc](https://godoc.org/github.com/aerospike/aerospike-client-go?status.svg)](https://pkg.go.dev/github.com/aerospike/aerospike-client-go).

<a name="App-Engine"></a>
## Google App Engine

To build the library for App Engine, build it with the build tag `app_engine`. Aggregation functionality is not available in this build.


<a name="Reflection"></a>
## Reflection, and Object API

To make the library both flexible and fast, we had to integrate the reflection API (methods with `[Get/Put/...]Object` names) tightly in the library. In case you wanted to avoid mixing those API in your app inadvertently, you can use the build tag `as_performance` to remove those APIs from the build.


## License

The Aerospike Go Client is made available under the terms of the Apache License, Version 2, as stated in the file `LICENSE`.

Individual files may be made available under their own specific license,
all compatible with Apache License, Version 2. Please see individual files for details.

