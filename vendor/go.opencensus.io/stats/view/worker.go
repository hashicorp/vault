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
	"errors"
	"fmt"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/internal"
	"go.opencensus.io/tag"
)

func init() {
	defaultWorker = newWorker()
	go defaultWorker.start()
	internal.DefaultRecorder = record
}

type measureRef struct {
	measure stats.Measure
	views   map[*View]struct{}
}

type worker struct {
	measures   map[string]*measureRef
	views      map[string]*View
	startTimes map[*View]time.Time

	timer      *time.Ticker
	c          chan command
	quit, done chan bool
}

var defaultWorker *worker

var defaultReportingDuration = 10 * time.Second

// FindMeasure returns a registered view associated with this name.
// If no registered view is found, nil is returned.
func Find(name string) (v *View) {
	req := &getViewByNameReq{
		name: name,
		c:    make(chan *getViewByNameResp),
	}
	defaultWorker.c <- req
	resp := <-req.c
	return resp.v
}

// Register registers view. It returns an error if the view is already registered.
//
// Subscription automatically registers a view.
// Most users will not register directly but register via subscription.
// Registration can be used by libraries to claim a view name.
//
// Unregister the view once the view is not required anymore.
func Register(v *View) error {
	req := &registerViewReq{
		v:   v,
		err: make(chan error),
	}
	defaultWorker.c <- req
	return <-req.err
}

// Unregister removes the previously registered view. It returns an error
// if the view wasn't registered. All data collected and not reported for the
// corresponding view will be lost. The view is automatically be unsubscribed.
func Unregister(v *View) error {
	req := &unregisterViewReq{
		v:   v,
		err: make(chan error),
	}
	defaultWorker.c <- req
	return <-req.err
}

// Subscribe subscribes a view. Once a view is subscribed, it reports data
// via the exporters.
// During subscription, if the view wasn't registered, it will be automatically
// registered. Once the view is no longer needed to export data,
// user should unsubscribe from the view.
func (v *View) Subscribe() error {
	req := &subscribeToViewReq{
		v:   v,
		err: make(chan error),
	}
	defaultWorker.c <- req
	return <-req.err
}

// Unsubscribe unsubscribes a previously subscribed view.
// Data will not be exported from this view once unsubscription happens.
func (v *View) Unsubscribe() error {
	req := &unsubscribeFromViewReq{
		v:   v,
		err: make(chan error),
	}
	defaultWorker.c <- req
	return <-req.err
}

// RetrieveData returns the current collected data for the view.
func (v *View) RetrieveData() ([]*Row, error) {
	if v == nil {
		return nil, errors.New("cannot retrieve data from nil view")
	}
	req := &retrieveDataReq{
		now: time.Now(),
		v:   v,
		c:   make(chan *retrieveDataResp),
	}
	defaultWorker.c <- req
	resp := <-req.c
	return resp.rows, resp.err
}

func record(tags *tag.Map, now time.Time, ms interface{}) {
	req := &recordReq{
		now: now,
		tm:  tags,
		ms:  ms.([]stats.Measurement),
	}
	defaultWorker.c <- req
}

// SetReportingPeriod sets the interval between reporting aggregated views in
// the program. If duration is less than or
// equal to zero, it enables the default behavior.
func SetReportingPeriod(d time.Duration) {
	// TODO(acetechnologist): ensure that the duration d is more than a certain
	// value. e.g. 1s
	req := &setReportingPeriodReq{
		d: d,
		c: make(chan bool),
	}
	defaultWorker.c <- req
	<-req.c // don't return until the timer is set to the new duration.
}

func newWorker() *worker {
	return &worker{
		measures:   make(map[string]*measureRef),
		views:      make(map[string]*View),
		startTimes: make(map[*View]time.Time),
		timer:      time.NewTicker(defaultReportingDuration),
		c:          make(chan command),
		quit:       make(chan bool),
		done:       make(chan bool),
	}
}

func (w *worker) start() {
	for {
		select {
		case cmd := <-w.c:
			if cmd != nil {
				cmd.handleCommand(w)
			}
		case <-w.timer.C:
			w.reportUsage(time.Now())
		case <-w.quit:
			w.timer.Stop()
			close(w.c)
			w.done <- true
			return
		}
	}
}

func (w *worker) stop() {
	w.quit <- true
	<-w.done
}

func (w *worker) getMeasureRef(m stats.Measure) *measureRef {
	if mr, ok := w.measures[m.Name()]; ok {
		return mr
	}
	mr := &measureRef{
		measure: m,
		views:   make(map[*View]struct{}),
	}
	w.measures[m.Name()] = mr
	return mr
}

func (w *worker) tryRegisterView(v *View) error {
	if err := checkViewName(v.name); err != nil {
		return err
	}
	if x, ok := w.views[v.Name()]; ok {
		if x != v {
			return fmt.Errorf("cannot register view %q; another view with the same name is already registered", v.Name())
		}

		// the view is already registered so there is nothing to do and the
		// command is considered successful.
		return nil
	}

	if v.Measure() == nil {
		return fmt.Errorf("cannot register view %q: measure not defined", v.Name())
	}

	w.views[v.Name()] = v
	ref := w.getMeasureRef(v.Measure())
	ref.views[v] = struct{}{}

	return nil
}

func (w *worker) reportUsage(start time.Time) {
	for _, v := range w.views {
		if !v.isSubscribed() {
			continue
		}
		rows := v.collectedRows(start)
		s, ok := w.startTimes[v]
		if !ok {
			w.startTimes[v] = start
		} else {
			start = s
		}
		// Make sure collector is never going
		// to mutate the exported data.
		rows = deepCopyRowData(rows)
		viewData := &Data{
			View:  v,
			Start: start,
			End:   time.Now(),
			Rows:  rows,
		}
		exportersMu.Lock()
		for e := range exporters {
			e.ExportView(viewData)
		}
		exportersMu.Unlock()
	}
}

func deepCopyRowData(rows []*Row) []*Row {
	newRows := make([]*Row, 0, len(rows))
	for _, r := range rows {
		newRows = append(newRows, &Row{
			Data: r.Data.clone(),
			Tags: r.Tags,
		})
	}
	return newRows
}
