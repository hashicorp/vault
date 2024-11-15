// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package randutil provides common random number utilities.
package randutil

import (
	crand "crypto/rand"
	"fmt"
	"io"

	xrand "go.mongodb.org/mongo-driver/internal/rand"
)

// NewLockedRand returns a new "x/exp/rand" pseudo-random number generator seeded with a
// cryptographically-secure random number.
// It is safe to use from multiple goroutines.
func NewLockedRand() *xrand.Rand {
	var randSrc = new(xrand.LockedSource)
	randSrc.Seed(cryptoSeed())
	return xrand.New(randSrc)
}

// cryptoSeed returns a random uint64 read from the "crypto/rand" random number generator. It is
// intended to be used to seed pseudorandom number generators at package initialization. It panics
// if it encounters any errors.
func cryptoSeed() uint64 {
	var b [8]byte
	_, err := io.ReadFull(crand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("failed to read 8 bytes from a \"crypto/rand\".Reader: %v", err))
	}

	return (uint64(b[0]) << 0) | (uint64(b[1]) << 8) | (uint64(b[2]) << 16) | (uint64(b[3]) << 24) |
		(uint64(b[4]) << 32) | (uint64(b[5]) << 40) | (uint64(b[6]) << 48) | (uint64(b[7]) << 56)
}
