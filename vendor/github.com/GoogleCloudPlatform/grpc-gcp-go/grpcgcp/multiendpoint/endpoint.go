/*
 *
 * Copyright 2023 gRPC authors.
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
 *
 */

package multiendpoint

import (
	"fmt"
	"time"
)

type status int

// Status of an endpoint.
const (
	unavailable status = iota
	available
	recovering
)

func (s status) String() string {
	switch s {
	case unavailable:
		return "Unavailable"
	case available:
		return "Available"
	case recovering:
		return "Recovering"
	default:
		return fmt.Sprintf("%d", s)
	}
}

type endpoint struct {
	id           string
	priority     int
	status       status
	lastChange   time.Time
	futureChange timerAlike
}
