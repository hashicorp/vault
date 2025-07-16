// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build tools

// This file is here for backwards compat only. You can now use make instead of go generate to
// install tools.

// You can replace
// $ go generate -tags tools tools/tools.go
// with
// $ make tools

package tools

//go:generate ./tools.sh install-tools
