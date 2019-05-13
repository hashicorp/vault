# pointerstructure [![GoDoc](https://godoc.org/github.com/mitchellh/pointerstructure?status.svg)](https://godoc.org/github.com/mitchellh/pointerstructure)

pointerstructure is a Go library for identifying a specific value within
any Go structure using a string syntax.

pointerstructure is based on
[JSON Pointer (RFC 6901)](https://tools.ietf.org/html/rfc6901), but
reimplemented for Go.

The goal of pointerstructure is to provide a single, well-known format
for addressing a specific value. This can be useful for user provided
input on structures, diffs of structures, etc.

## Features

  * Get the value for an address

  * Set the value for an address within an existing structure

  * Delete the value at an address

  * Sorting a list of addresses

## Installation

Standard `go get`:

```
$ go get github.com/mitchellh/pointerstructure
```

## Usage & Example

For usage and examples see the [Godoc](http://godoc.org/github.com/mitchellh/pointerstructure).

A quick code example is shown below:

```go
complex := map[string]interface{}{
	"alice": 42,
	"bob": []interface{}{
		map[string]interface{}{
			"name": "Bob",
		},
	},
}

value, err := pointerstructure.Get(complex, "/bob/0/name")
if err != nil {
	panic(err)
}

fmt.Printf("%s", value)
// Output:
// Bob
```

Continuing the example above, you can also set values:

```go
value, err = pointerstructure.Set(complex, "/bob/0/name", "Alice")
if err != nil {
	panic(err)
}

value, err = pointerstructure.Get(complex, "/bob/0/name")
if err != nil {
	panic(err)
}

fmt.Printf("%s", value)
// Output:
// Alice
```
