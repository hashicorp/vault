# Couchbase Connection Strings for Go

This library allows you to parse and resolve Couchbase Connection Strings in Go.
This is used by the Couchbase Go SDK, as well as various tools throughout the
Couchbase infrastructure.


## Using the Library

To parse a connection string, simply call `Parse` with your connection string.
You will receive a `ConnSpec` structure representing the connection string`:

```go
type Address struct {
	Host string
	Port int
}

type ConnSpec struct {
	Scheme string
	Addresses []Address
	Bucket string
	Options map[string][]string
}
```

One you have a parsed connection string, you can also use our resolver to take
the `ConnSpec` and resolve any DNS SRV records as well as generate a list of
endpoints for the Couchbase server.  You will receive a `ResolvedConnSpec`
structure in return:

```go
type ResolvedConnSpec struct {
	UseSsl bool
	MemdHosts []Address
	HttpHosts []Address
	Bucket string
	Options map[string][]string
}
```

## License
Copyright 2016 Couchbase Inc.

Licensed under the Apache License, Version 2.0.

See
[LICENSE](https://github.com/couchbaselabs/gocbconnstr/blob/master/LICENSE)
for further details.