// +build tools

// This file ensures tool dependencies are kept in sync.  This is the
// recommended way of doing this according to
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// To install the following tools at the version used by this repo run:
// $ make bootstrap
// or
// $ go generate -tags tools tools/tools.go

package tools

// use this instead of google.golang.org/protobuf/cmd/protoc-gen-go since this supports grpc plugin while the other does not.
// see https://github.com/golang/protobuf/releases#v1.4-generated-code and
// https://github.com/protocolbuffers/protobuf-go/releases/tag/v1.20.0#v1.20-grpc-support
//go:generate go install github.com/golang/protobuf/protoc-gen-go
import _ "github.com/golang/protobuf/protoc-gen-go"

//go:generate go install golang.org/x/tools/cmd/goimports
import _ "golang.org/x/tools/cmd/goimports"

//go:generate go install github.com/mitchellh/gox
import _ "github.com/mitchellh/gox"

//go:generate go install github.com/hashicorp/go-bindata
import _ "github.com/hashicorp/go-bindata"

//go:generate go install github.com/elazarl/go-bindata-assetfs
import _ "github.com/elazarl/go-bindata-assetfs"

//go:generate go install github.com/client9/misspell/cmd/misspell
import _ "github.com/client9/misspell/cmd/misspell"
