// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Package rand implements random value functions.
package rand

import (
	"crypto/rand"
)

const (
	// alpa numeric character set
	csAlphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// RandomBytes returns a random byte slice of size n and panics if crypto random reader returns an error.
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err.Error()) // rand should never fail
	}
	return b
}

// RandomString returns a random string of alphanumeric characters and panics if crypto random reader returns an error.
func RandomString(n int) string {
	bytes := RandomBytes(n)
	size := byte(len(csAlphanum)) // len character sets <= max(byte)
	for i, b := range bytes {
		bytes[i] = csAlphanum[b%size]
	}
	return string(bytes)
}
