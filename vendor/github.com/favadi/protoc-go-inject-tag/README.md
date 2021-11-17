# protoc-go-inject-tag

[![Build Status](https://www.travis-ci.com/favadi/protoc-go-inject-tag.svg?branch=master)](https://www.travis-ci.com/favadi/protoc-go-inject-tag)
[![Go Report Card](https://goreportcard.com/badge/github.com/favadi/protoc-go-inject-tag)](https://goreportcard.com/report/github.com/favadi/protoc-go-inject-tag)
[![Coverage Status](https://coveralls.io/repos/github/favadi/protoc-go-inject-tag/badge.svg)](https://coveralls.io/github/favadi/protoc-go-inject-tag)

## Why?

Golang [protobuf](https://github.com/golang/protobuf) doesn't support
[custom tags to generated structs](https://github.com/golang/protobuf/issues/52).
This tool injects custom tags to generated protobuf files, which is commonly
used for validating fields, omitting fields from JSON data, etc.

## Install

- [protobuf version 3](https://github.com/google/protobuf)

  For OS X:

  ```console
  $ brew install protobuf
  ```

- go support for protobuf: `go get -u github.com/golang/protobuf/{proto,protoc-gen-go}`

- `go get github.com/favadi/protoc-go-inject-tag` or download the
  binaries from the releases page.

## Usage

```console
$ protoc-go-inject-tag -h
Usage of protoc-go-inject-tag:
  -XXX_skip string
        tags that should be skipped (applies 'tag:"-"') for unknown fields (deprecated since protoc-gen-go v1.4.0)
  -input string
        pattern to match input file(s)
  -verbose
        verbose logging
```

Add a comment with the following syntax before fields, and these will be
injected into the resulting `.pb.go` file. This can be specified above the
field, or trailing the field.

```proto
// @gotags: custom_tag:"custom_value"
```

## Example

```proto
// file: test.proto
syntax = "proto3";

package pb;
option go_package = "/pb";

message IP {
  // @gotags: valid:"ip"
  string Address = 1;

  // Or:
  string MAC = 2; // @gotags: validate:"omitempty"
}
```

Generate your `.pb.go` files with the protoc command as normal:

```console
$ protoc --proto_path=. --go_out=paths=source_relative:. test.proto
```

Then run `protoc-go-inject-tag` against the generated files (e.g `test.pb.go`):

```console
$ protoc-go-inject-tag -input=./test.pb.go
# or
$ protoc-go-inject-tag -input="*.pb.go"
```

The custom tags will be injected to `test.pb.go`:

```go
type IP struct {
	// @gotags: valid:"ip"
	Address string `protobuf:"bytes,1,opt,name=Address,json=address" json:"Address,omitempty" valid:"ip"`
}
```

## Deprecated functionality

#### Skip `XXX_*` fields

To skip the tag for the generated `XXX_*` fields (unknown fields), use the
`-XXX_skip=yaml,xml` flag. This is deprecated, as this functionality hasn't
existed in `protoc-gen-go` since v1.4.x.

#### `inject_tag` keyword

Since **v1.3.0**, we recommend using `@gotags:` rather than `@inject_tags:`,
as `@gotags` is more indicative of the language the comment is for. We don't
plan on removing `@inject_tags:` support anytime soon, however we strongly
recommend switching to `@gotags`.
