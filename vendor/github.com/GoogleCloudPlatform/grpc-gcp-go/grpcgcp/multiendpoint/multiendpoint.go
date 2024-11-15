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

// Package multiendpoint implements multiendpoint feature. See [MultiEndpoint]
package multiendpoint

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type timerAlike interface {
	Reset(time.Duration) bool
	Stop() bool
}

// To be redefined in tests.
var (
	timeNow = func() time.Time {
		return time.Now()
	}
	timeAfterFunc = func(d time.Duration, f func()) timerAlike {
		return time.AfterFunc(d, f)
	}
)

// MultiEndpoint holds a list of endpoints, tracks their availability and defines the current
// endpoint. An endpoint has a priority defined by its position in the list (first item has top
// priority).
//
// The current endpoint is the highest available endpoint in the list. If no endpoint is available,
// MultiEndpoint sticks to the previously current endpoint.
//
// Sometimes switching between endpoints can be costly, and it is worth waiting for some time
// after current endpoint becomes unavailable. For this case, use
// [MultiEndpointOptions.RecoveryTimeout] to set the recovery timeout. MultiEndpoint will keep the
// current endpoint for up to recovery timeout after it became unavailable to give it some time to
// recover.
//
// The list of endpoints can be changed at any time with [MultiEndpoint.SetEndpoints] function.
// MultiEndpoint will:
//   - remove obsolete endpoints;
//   - preserve remaining endpoints and their states;
//   - add new endpoints;
//   - update all endpoints priority according to the new order;
//   - change current endpoint if necessary.
//
// After updating the list of endpoints, MultiEndpoint will switch the current endpoint to the
// highest available endpoint in the list. If you have many processes using MultiEndpoint, this may
// lead to immediate shift of all traffic which may be undesired. To smooth this transfer, use
// [MultiEndpointOptions.SwitchingDelay] with randomized value to introduce a jitter. Each
// MultiEndpoint will delay switching from an available endpoint to another endpoint for this amount
// of time. This delay is only applicable when switching from a lower priority available endpoint to
// a higher priority available endpoint.
type MultiEndpoint interface {
	// Current returns current endpoint.
	//
	// Note that the read is not synchronized and in case of a race condition there is a chance of
	// getting an outdated current endpoint.
	Current() string

	// SetEndpointAvailability informs MultiEndpoint when an endpoint becomes available or unavailable.
	// This may change the current endpoint.
	SetEndpointAvailability(e string, avail bool)

	// SetEndpoints updates a list of endpoints:
	//   - remove obsolete endpoints
	//   - preserve remaining endpoints and their states
	//   - add new endpoints
	//   - update all endpoints priority according to the new order
	// This may change the current endpoint.
	SetEndpoints(endpoints []string) error
}

// MultiEndpointOptions is used for configuring [MultiEndpoint].
type MultiEndpointOptions struct {
	// A list of endpoints ordered by priority (first endpoint has top priority).
	Endpoints []string
	// RecoveryTimeout sets the amount of time MultiEndpoint keeps endpoint as current after it
	// became unavailable.
	RecoveryTimeout time.Duration
	// When switching from a lower priority available endpoint to a higher priority available
	// endpoint the MultiEndpoint will delay the switch for this duration.
	SwitchingDelay time.Duration
}

// NewMultiEndpoint validates options and creates a new [MultiEndpoint].
func NewMultiEndpoint(b *MultiEndpointOptions) (MultiEndpoint, error) {
	if len(b.Endpoints) == 0 {
		return nil, fmt.Errorf("endpoints list cannot be empty")
	}

	me := &multiEndpoint{
		recoveryTimeout: b.RecoveryTimeout,
		switchingDelay:  b.SwitchingDelay,
		current:         b.Endpoints[0],
	}
	eMap := make(map[string]*endpoint)
	for i, e := range b.Endpoints {
		eMap[e] = me.newEndpoint(e, i)
	}
	me.endpoints = eMap
	return me, nil
}

type multiEndpoint struct {
	sync.RWMutex

	endpoints       map[string]*endpoint
	recoveryTimeout time.Duration
	switchingDelay  time.Duration
	current         string
	future          string
}

// Current returns current endpoint.
func (me *multiEndpoint) Current() string {
	me.RLock()
	defer me.RUnlock()
	return me.current
}

// SetEndpoints updates endpoints list:
//   - remove obsolete endpoints;
//   - preserve remaining endpoints and their states;
//   - add new endpoints;
//   - update all endpoints priority according to the new order;
//   - change current endpoint if necessary.
func (me *multiEndpoint) SetEndpoints(endpoints []string) error {
	me.Lock()
	defer me.Unlock()
	if len(endpoints) == 0 {
		return errors.New("endpoints list cannot be empty")
	}
	newEndpoints := make(map[string]struct{})
	for _, v := range endpoints {
		newEndpoints[v] = struct{}{}
	}
	// Remove obsolete endpoints.
	for e := range me.endpoints {
		if _, ok := newEndpoints[e]; !ok {
			delete(me.endpoints, e)
		}
	}
	// Add new endpoints and update priority.
	for i, e := range endpoints {
		if _, ok := me.endpoints[e]; !ok {
			me.endpoints[e] = me.newEndpoint(e, i)
		} else {
			me.endpoints[e].priority = i
		}
	}

	me.maybeUpdateCurrent()
	return nil
}

// Updates current to the top-priority available endpoint unless the current endpoint is
// recovering.
//
// Must be run under me.Lock.
func (me *multiEndpoint) maybeUpdateCurrent() {
	c, exists := me.endpoints[me.current]
	var topA *endpoint
	var top *endpoint
	for _, e := range me.endpoints {
		if e.status == available && (topA == nil || topA.priority > e.priority) {
			topA = e
		}
		if top == nil || top.priority > e.priority {
			top = e
		}
	}

	if exists && c.status == recovering && (topA == nil || topA.priority > c.priority) {
		// Let current endpoint recover while no higher priority endpoints available.
		return
	}

	// Always prefer top available endpoint.
	if topA != nil {
		me.switchFromTo(c, topA)
		return
	}

	// If no current endpoint exists, resort to the top priority endpoint immediately.
	if !exists {
		me.current = top.id
	}
}

func (me *multiEndpoint) newEndpoint(id string, priority int) *endpoint {
	s := unavailable
	if me.recoveryTimeout > 0 {
		s = recovering
	}
	e := &endpoint{
		id:       id,
		priority: priority,
		status:   s,
	}
	if e.status == recovering {
		me.scheduleUnavailable(e)
	}
	return e
}

// Changes or schedules a change of current to the endpoint t.
//
// Must be run under me.Lock.
func (me *multiEndpoint) switchFromTo(f, t *endpoint) {
	if me.current == t.id {
		return
	}

	if me.switchingDelay == 0 || f == nil || f.status == unavailable {
		// Switching immediately if no delay or no current or current is unavailable.
		me.current = t.id
		return
	}

	me.future = t.id
	timeAfterFunc(me.switchingDelay, func() {
		me.Lock()
		defer me.Unlock()
		if e, ok := me.endpoints[me.future]; ok && e.status == available {
			me.current = e.id
		}
	})
}

// SetEndpointAvailability updates the state of an endpoint.
func (me *multiEndpoint) SetEndpointAvailability(e string, avail bool) {
	me.Lock()
	defer me.Unlock()
	me.setEndpointAvailability(e, avail)
	me.maybeUpdateCurrent()
}

// Must be run under me.Lock.
func (me *multiEndpoint) setEndpointAvailability(e string, avail bool) {
	ee, ok := me.endpoints[e]
	if !ok {
		return
	}

	if avail {
		setState(ee, available)
		return
	}

	if ee.status != available {
		return
	}

	if me.recoveryTimeout == 0 {
		setState(ee, unavailable)
		return
	}

	setState(ee, recovering)
	me.scheduleUnavailable(ee)
}

// Change the state of endpoint e to state s.
//
// Must be run under me.Lock.
func setState(e *endpoint, s status) {
	if e.futureChange != nil {
		e.futureChange.Stop()
	}
	e.status = s
	e.lastChange = timeNow()
}

// Schedule endpoint e to become unavailable after recoveryTimeout.
func (me *multiEndpoint) scheduleUnavailable(e *endpoint) {
	stateChange := e.lastChange
	e.futureChange = timeAfterFunc(me.recoveryTimeout, func() {
		me.Lock()
		defer me.Unlock()
		if e.lastChange != stateChange {
			// This timer is outdated.
			return
		}
		setState(e, unavailable)
		me.maybeUpdateCurrent()
	})
}
