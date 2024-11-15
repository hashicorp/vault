// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.23

//go:debug gotypesalias=1

package main

// Materialize aliases whenever the go toolchain version is after 1.23 (#69772).
// Remove this file after go.mod >= 1.23 (which implies gotypesalias=1).
