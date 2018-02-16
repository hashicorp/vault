// Copyright 2017, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package view

import (
	"fmt"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

type command interface {
	handleCommand(w *worker)
}

// getViewByNameReq is the command to get a view given its name.
type getViewByNameReq struct {
	name string
	c    chan *getViewByNameResp
}

type getViewByNameResp struct {
	v *View
}

func (cmd *getViewByNameReq) handleCommand(w *worker) {
	cmd.c <- &getViewByNameResp{w.views[cmd.name]}
}

// registerViewReq is the command to register a view with the library.
type registerViewReq struct {
	v   *View
	err chan error
}

func (cmd *registerViewReq) handleCommand(w *worker) {
	cmd.err <- w.tryRegisterView(cmd.v)
}

// unregisterViewReq is the command to unregister a view from the library.
type unregisterViewReq struct {
	v   *View
	err chan error
}

func (cmd *unregisterViewReq) handleCommand(w *worker) {
	v, ok := w.views[cmd.v.Name()]
	if !ok {
		cmd.err <- nil
		return
	}
	if v != cmd.v {
		cmd.err <- nil
		return
	}
	if v.isSubscribed() {
		cmd.err <- fmt.Errorf("cannot unregister view %q; all subscriptions must be unsubscribed first", cmd.v.Name())
		return
	}
	delete(w.views, cmd.v.Name())
	ref := w.getMeasureRef(v.Measure())
	delete(ref.views, v)
	cmd.err <- nil
}

// subscribeToViewReq is the command to subscribe to a view.
type subscribeToViewReq struct {
	v   *View
	err chan error
}

func (cmd *subscribeToViewReq) handleCommand(w *worker) {
	if cmd.v.isSubscribed() {
		cmd.err <- nil
		return
	}
	if err := w.tryRegisterView(cmd.v); err != nil {
		cmd.err <- fmt.Errorf("cannot subscribe to view: %v", err)
		return
	}
	cmd.v.subscribe()
	cmd.err <- nil
}

// unsubscribeFromViewReq is the command to unsubscribe to a view. Has no
// impact on the data collection for client that are pulling data from the
// library.
type unsubscribeFromViewReq struct {
	v   *View
	err chan error
}

func (cmd *unsubscribeFromViewReq) handleCommand(w *worker) {
	cmd.v.unsubscribe()
	if !cmd.v.isSubscribed() {
		// this was the last subscription and view is not collecting anymore.
		// The collected data can be cleared.
		cmd.v.clearRows()
	}
	// we always return nil because this operation never fails. However we
	// still need to return something on the channel to signal to the waiting
	// go routine that the operation completed.
	cmd.err <- nil
}

// retrieveDataReq is the command to retrieve data for a view.
type retrieveDataReq struct {
	now time.Time
	v   *View
	c   chan *retrieveDataResp
}

type retrieveDataResp struct {
	rows []*Row
	err  error
}

func (cmd *retrieveDataReq) handleCommand(w *worker) {
	if _, ok := w.views[cmd.v.Name()]; !ok {
		cmd.c <- &retrieveDataResp{
			nil,
			fmt.Errorf("cannot retrieve data; view %q is not registered", cmd.v.Name()),
		}
		return
	}

	if !cmd.v.isSubscribed() {
		cmd.c <- &retrieveDataResp{
			nil,
			fmt.Errorf("cannot retrieve data; view %q has no subscriptions or collection is not forcibly started", cmd.v.Name()),
		}
		return
	}
	cmd.c <- &retrieveDataResp{
		cmd.v.collectedRows(cmd.now),
		nil,
	}
}

// recordReq is the command to record data related to multiple measures
// at once.
type recordReq struct {
	now time.Time
	tm  *tag.Map
	ms  []stats.Measurement
}

func (cmd *recordReq) handleCommand(w *worker) {
	for _, m := range cmd.ms {
		ref := w.getMeasureRef(m.Measure)
		for v := range ref.views {
			v.addSample(cmd.tm, m.Value, cmd.now)
		}
	}
}

// setReportingPeriodReq is the command to modify the duration between
// reporting the collected data to the subscribed clients.
type setReportingPeriodReq struct {
	d time.Duration
	c chan bool
}

func (cmd *setReportingPeriodReq) handleCommand(w *worker) {
	w.timer.Stop()
	if cmd.d <= 0 {
		w.timer = time.NewTicker(defaultReportingDuration)
	} else {
		w.timer = time.NewTicker(cmd.d)
	}
	cmd.c <- true
}
