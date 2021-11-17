// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: Remove entire file when support for Go1.12 and lower has been dropped.
// +build !go1.13

package spanner

import "golang.org/x/xerrors"

// unwrap is a generic implementation of (errors|xerrors).Unwrap(error). This
// implementation uses xerrors and is included in Go 1.12 and earlier builds.
func unwrap(err error) error {
	return xerrors.Unwrap(err)
}

// errorAs is a generic implementation of
// (errors|xerrors).As(error, interface{}). This implementation uses xerrors
// and is included in Go 1.12 and earlier builds.
func errorAs(err error, target interface{}) bool {
	return xerrors.As(err, target)
}
