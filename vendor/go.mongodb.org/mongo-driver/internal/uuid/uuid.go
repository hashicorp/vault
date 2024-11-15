// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package uuid

import (
	"encoding/hex"
	"io"

	"go.mongodb.org/mongo-driver/internal/randutil"
)

// UUID represents a UUID.
type UUID [16]byte

// A source is a UUID generator that reads random values from a io.Reader.
// It should be safe to use from multiple goroutines.
type source struct {
	random io.Reader
}

// new returns a random UUIDv4 with bytes read from the source's random number generator.
func (s *source) new() (UUID, error) {
	var uuid UUID
	_, err := io.ReadFull(s.random, uuid[:])
	if err != nil {
		return UUID{}, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid, nil
}

// newSource returns a source that uses a pseudo-random number generator in reandutil package.
// It is intended to be used to initialize the package-global UUID generator.
func newSource() *source {
	return &source{
		random: randutil.NewLockedRand(),
	}
}

// globalSource is a package-global pseudo-random UUID generator.
var globalSource = newSource()

// New returns a random UUIDv4. It uses a global pseudo-random number generator in randutil
// at package initialization.
//
// New should not be used to generate cryptographically-secure random UUIDs.
func New() (UUID, error) {
	return globalSource.new()
}

func (uuid UUID) String() string {
	var str [36]byte
	hex.Encode(str[:], uuid[:4])
	str[8] = '-'
	hex.Encode(str[9:13], uuid[4:6])
	str[13] = '-'
	hex.Encode(str[14:18], uuid[6:8])
	str[18] = '-'
	hex.Encode(str[19:23], uuid[8:10])
	str[23] = '-'
	hex.Encode(str[24:], uuid[10:])
	return string(str[:])
}
