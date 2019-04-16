/*
 * keyselector.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2013-2018 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// FoundationDB Go API

package fdb

// A Selectable can be converted to a FoundationDB KeySelector. All functions in
// the FoundationDB API that resolve a key selector to a key accept Selectable.
type Selectable interface {
	FDBKeySelector() KeySelector
}

// KeySelector represents a description of a key in a FoundationDB database. A
// KeySelector may be resolved to a specific key with the GetKey method, or used
// as the endpoints of a SelectorRange to be used with a GetRange function.
//
// The most common key selectors are constructed with the functions documented
// below. For details of how KeySelectors are specified and resolved, see
// https://apple.github.io/foundationdb/developer-guide.html#key-selectors.
type KeySelector struct {
	Key     KeyConvertible
	OrEqual bool
	Offset  int
}

func (ks KeySelector) FDBKeySelector() KeySelector {
	return ks
}

// LastLessThan returns the KeySelector specifying the lexigraphically greatest
// key present in the database which is lexigraphically strictly less than the
// given key.
func LastLessThan(key KeyConvertible) KeySelector {
	return KeySelector{key, false, 0}
}

// LastLessOrEqual returns the KeySelector specifying the lexigraphically
// greatest key present in the database which is lexigraphically less than or
// equal to the given key.
func LastLessOrEqual(key KeyConvertible) KeySelector {
	return KeySelector{key, true, 0}
}

// FirstGreaterThan returns the KeySelector specifying the lexigraphically least
// key present in the database which is lexigraphically strictly greater than
// the given key.
func FirstGreaterThan(key KeyConvertible) KeySelector {
	return KeySelector{key, true, 1}
}

// FirstGreaterOrEqual returns the KeySelector specifying the lexigraphically
// least key present in the database which is lexigraphically greater than or
// equal to the given key.
func FirstGreaterOrEqual(key KeyConvertible) KeySelector {
	return KeySelector{key, false, 1}
}
