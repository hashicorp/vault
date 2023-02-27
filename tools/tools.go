// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build tools

// This file ensures tool dependencies are kept in sync.  This is the
// recommended way of doing this according to
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// To install the following tools at the version used by this repo run:
// $ make bootstrap
// or
// $ go generate -tags tools tools/tools.go

package tools

//go:generate go install golang.org/x/tools/cmd/goimports
//go:generate go install github.com/client9/misspell/cmd/misspell
//go:generate go install mvdan.cc/gofumpt
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go install github.com/favadi/protoc-go-inject-tag
//go:generate go install github.com/golangci/revgrep/cmd/revgrep
import (
	_ "golang.org/x/tools/cmd/goimports"

	_ "github.com/client9/misspell/cmd/misspell"

	_ "mvdan.cc/gofumpt"

	_ "google.golang.org/protobuf/cmd/protoc-gen-go"

	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"

	_ "github.com/favadi/protoc-go-inject-tag"

	_ "github.com/golangci/revgrep/cmd/revgrep"
)
