// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import "time"

// AdminPolicy contains attributes used for user administration commands.
type AdminPolicy struct {

	// User administration command socket timeout.
	// Default is 2 seconds.
	Timeout time.Duration
}

// NewAdminPolicy generates a new AdminPolicy with default values.
func NewAdminPolicy() *AdminPolicy {
	return &AdminPolicy{
		Timeout: _DEFAULT_TIMEOUT,
	}
}

func (ap *AdminPolicy) deadline() (deadline time.Time) {
	if ap != nil && ap.Timeout > 0 {
		deadline = time.Now().Add(ap.Timeout)
	}

	return deadline
}

func (ap *AdminPolicy) timeout() time.Duration {
	if ap != nil && ap.Timeout > 0 {
		return ap.Timeout
	}

	return _DEFAULT_TIMEOUT
}
