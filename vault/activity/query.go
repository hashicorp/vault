// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package activity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/helper/compressutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type NamespaceRecord struct {
	NamespaceID     string         `json:"namespace_id"`
	Entities        uint64         `json:"entities"`
	NonEntityTokens uint64         `json:"non_entity_tokens"`
	SecretSyncs     uint64         `json:"secret_syncs"`
	Mounts          []*MountRecord `json:"mounts"`
	ACMEClients     uint64         `json:"acme_clients"`
}

type CountsRecord struct {
	EntityClients    int `json:"entity_clients"`
	NonEntityClients int `json:"non_entity_clients"`
	SecretSyncs      int `json:"secret_syncs"`
	ACMEClients      int `json:"acme_clients"`
}

// HasCounts returns true when any of the record's fields have a non-zero value
func (c *CountsRecord) HasCounts() bool {
	return c.EntityClients+c.NonEntityClients+c.SecretSyncs+c.ACMEClients != 0
}

type NewClientRecord struct {
	Counts     *CountsRecord             `json:"counts"`
	Namespaces []*MonthlyNamespaceRecord `json:"namespaces"`
}

type MonthRecord struct {
	Timestamp  int64                     `json:"timestamp"`
	Counts     *CountsRecord             `json:"counts"`
	Namespaces []*MonthlyNamespaceRecord `json:"namespaces"`
	NewClients *NewClientRecord          `json:"new_clients"`
}

type MonthlyNamespaceRecord struct {
	NamespaceID string         `json:"namespace_id"`
	Counts      *CountsRecord  `json:"counts"`
	Mounts      []*MountRecord `json:"mounts"`
}

type MountRecord struct {
	MountPath string        `json:"mount_path"`
	MountType string        `json:"mount_type"`
	Counts    *CountsRecord `json:"counts"`
}

type PrecomputedQuery struct {
	StartTime  time.Time
	EndTime    time.Time
	Namespaces []*NamespaceRecord `json:"namespaces"`
	Months     []*MonthRecord     `json:"months"`
}

type PrecomputedQueryStore struct {
	logger log.Logger
	view   logical.Storage
}

// The query store should be initialized with a view to the subdirectory
// it should use, like "queries".
func NewPrecomputedQueryStore(logger log.Logger, view logical.Storage, retentionMonths int) *PrecomputedQueryStore {
	return &PrecomputedQueryStore{
		logger: logger,
		view:   view,
	}
}

func (s *PrecomputedQueryStore) Put(ctx context.Context, p *PrecomputedQuery) error {
	path := fmt.Sprintf("%v/%v", p.StartTime.Unix(), p.EndTime.Unix())
	asJson, err := json.Marshal(p)
	if err != nil {
		return err
	}

	compressedPq, err := compressutil.Compress(asJson, &compressutil.CompressionConfig{
		Type: compressutil.CompressionTypeLZ4,
	})
	if err != nil {
		return err
	}

	err = s.view.Put(ctx, &logical.StorageEntry{
		Key:   path,
		Value: compressedPq,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PrecomputedQueryStore) listStartTimes(ctx context.Context) ([]time.Time, error) {
	// We could cache this to save a storage operation on each query,
	// but that seems like a marginal improvment.
	rawStartTimes, err := s.view.List(ctx, "")
	if err != nil {
		return nil, err
	}
	startTimes := make([]time.Time, 0, len(rawStartTimes))

	for _, raw := range rawStartTimes {
		t, err := timeutil.ParseTimeFromPath(raw)
		if err != nil {
			s.logger.Warn("could not parse precomputed query subdirectory", "key", raw)
			continue
		}
		startTimes = append(startTimes, t)
	}
	return startTimes, nil
}

func (s *PrecomputedQueryStore) listEndTimes(ctx context.Context, startTime time.Time) ([]time.Time, error) {
	rawEndTimes, err := s.view.List(ctx, fmt.Sprintf("%v/", startTime.Unix()))
	if err != nil {
		return nil, err
	}
	endTimes := make([]time.Time, 0, len(rawEndTimes))

	for _, raw := range rawEndTimes {
		val, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			s.logger.Warn("could not parse precomputed query end time", "key", raw)
			continue
		}
		endTimes = append(endTimes, time.Unix(val, 0).UTC())
	}
	return endTimes, nil
}

func (s *PrecomputedQueryStore) getMaxEndTime(ctx context.Context, startTime time.Time, endTimeBound time.Time) (time.Time, error) {
	rawEndTimes, err := s.view.List(ctx, fmt.Sprintf("%v/", startTime.Unix()))
	if err != nil {
		return time.Time{}, err
	}

	maxEndTime := time.Time{}
	for _, raw := range rawEndTimes {
		val, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			s.logger.Warn("could not parse precomputed query end time", "key", raw)
			continue
		}
		endTime := time.Unix(val, 0).UTC()
		s.logger.Trace("end time in consideration is", "end time", endTime, "end time bound", endTimeBound)
		if endTime.After(maxEndTime) && !endTime.After(endTimeBound) {
			s.logger.Trace("end time has been updated")
			maxEndTime = endTime
		}

	}
	return maxEndTime, nil
}

func (s *PrecomputedQueryStore) QueriesAvailable(ctx context.Context) (bool, error) {
	startTimes, err := s.listStartTimes(ctx)
	if err != nil {
		return false, err
	}
	return len(startTimes) > 0, nil
}

func (s *PrecomputedQueryStore) Get(ctx context.Context, startTime, endTime time.Time) (*PrecomputedQuery, error) {
	if startTime.After(endTime) {
		return nil, errors.New("start time is after end time")
	}
	startTime = timeutil.StartOfMonth(startTime)
	endTime = timeutil.EndOfMonth(endTime)
	s.logger.Trace("searching for matching queries", "startTime", startTime, "endTime", endTime)

	// Find the oldest continuous region which overlaps with the given range.
	// We only have to handle some collection of lower triangles like this,
	// not arbitrary sets of endpoints (except in the middle of writes or GC):
	//
	//     start ->
	// end   #
	//  |    ##
	//  V    ###
	//
	//           #
	//           ##
	//           ###
	//
	// (1) find all saved start times T that are
	//     in [startTime,endTime]
	//     (if there is some report that overlaps, it will
	//      have a start time in the range-- an overlap
	//      only at the end is impossible.)
	// (2) take the latest continguous region within
	//     that set
	// i.e., walk up the diagonal as far as we can in a single
	// triangle.
	// (These could be combined into a single pass, but
	// that seems more complicated to understand.)

	startTimes, err := s.listStartTimes(ctx)
	if err != nil {
		return nil, err
	}
	s.logger.Trace("retrieved start times from storage", "startTimes", startTimes)

	filteredList := make([]time.Time, 0)
	for _, t := range startTimes {
		if timeutil.InRange(t, startTime, endTime) {
			filteredList = append(filteredList, t)
		}
	}
	s.logger.Trace("filtered to range", "startTimes", filteredList)

	if len(filteredList) == 0 {
		return nil, nil
	}
	// Descending order, as required by the timeutil function
	sort.Slice(filteredList, func(i, j int) bool {
		return filteredList[i].After(filteredList[j])
	})

	closestStartTime := time.Time{}
	closestEndTime := time.Time{}
	maxTimeDifference := time.Duration(0)
	for i := len(filteredList) - 1; i >= 0; i-- {
		testStartTime := filteredList[i]
		s.logger.Trace("trying test start times", "startTime", testStartTime, "filteredList", filteredList)
		testEndTime, err := s.getMaxEndTime(ctx, testStartTime, endTime)
		if err != nil {
			return nil, err
		}
		if testEndTime.IsZero() {
			// Might happen if there's a race with GC
			s.logger.Warn("missing end times", "start time", testStartTime)
			continue
		}
		s.logger.Trace("retrieved max end time from storage", "endTime", testEndTime)
		diff := testEndTime.Sub(testStartTime)
		if diff >= maxTimeDifference {
			closestStartTime = testStartTime
			closestEndTime = testEndTime
			maxTimeDifference = diff
			s.logger.Trace("updating closest times")
		}
	}
	s.logger.Trace("chose start/end times", "startTime", closestStartTime, "endTime", closestEndTime)

	if closestStartTime.IsZero() || closestEndTime.IsZero() {
		s.logger.Warn("no start or end time in range", "start time", closestStartTime, "end time", closestEndTime)
		return nil, nil
	}

	path := fmt.Sprintf("%v/%v", closestStartTime.Unix(), closestEndTime.Unix())
	entry, err := s.view.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		s.logger.Warn("no end time entry found", "start time", closestStartTime, "end time", closestEndTime)
		return nil, nil
	}

	value, notCompressed, err := compressutil.Decompress(entry.Value)
	if err != nil {
		return nil, err
	}
	if notCompressed {
		value = entry.Value
	}

	p := &PrecomputedQuery{}
	err = json.Unmarshal(value, p)
	if err != nil {
		s.logger.Warn("failed query lookup at", "path", path)
		return nil, err
	}

	return p, nil
}

func (s *PrecomputedQueryStore) DeleteQueriesBefore(ctx context.Context, retentionThreshold time.Time) error {
	startTimes, err := s.listStartTimes(ctx)
	if err != nil {
		return err
	}

	for _, t := range startTimes {
		path := fmt.Sprintf("%v/", t.Unix())
		if t.Before(retentionThreshold) {
			rawEndTimes, err := s.view.List(ctx, path)
			if err != nil {
				return err
			}

			s.logger.Trace("deleting queries", "startTime", t)
			// Don't care about what the end time is,
			// the start time along determines deletion.
			for _, end := range rawEndTimes {
				err = s.view.Delete(ctx, path+end)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *MonthlyNamespaceRecord) ToNamespaceRecord() *NamespaceRecord {
	return &NamespaceRecord{
		NamespaceID:     m.NamespaceID,
		Entities:        uint64(m.Counts.EntityClients),
		NonEntityTokens: uint64(m.Counts.NonEntityClients),
		SecretSyncs:     uint64(m.Counts.SecretSyncs),
		Mounts:          m.Mounts,
		ACMEClients:     uint64(m.Counts.ACMEClients),
	}
}

func (n *NamespaceRecord) CombineWithMonthlyNamespaceRecord(nsRecord *MonthlyNamespaceRecord) {
	existingMounts := make(map[string]*MountRecord)
	for _, mountRecord := range n.Mounts {
		existingMounts[mountRecord.MountPath] = mountRecord
	}

	for _, mountRecord := range nsRecord.Mounts {
		if existingMountRecord, ok := existingMounts[mountRecord.MountPath]; ok {
			existingMountRecord.Add(mountRecord)
		} else {
			n.Mounts = append(n.Mounts, mountRecord)
		}
	}

	n.SecretSyncs += uint64(nsRecord.Counts.SecretSyncs)
	n.Entities += uint64(nsRecord.Counts.EntityClients)
	n.NonEntityTokens += uint64(nsRecord.Counts.NonEntityClients)
	n.ACMEClients += uint64(nsRecord.Counts.ACMEClients)
}

func (m *MountRecord) Add(m2 *MountRecord) {
	m.Counts.ACMEClients += m2.Counts.ACMEClients
	m.Counts.NonEntityClients += m2.Counts.NonEntityClients
	m.Counts.EntityClients += m2.Counts.EntityClients
	m.Counts.SecretSyncs += m2.Counts.SecretSyncs
}

func (q *PrecomputedQuery) CombineWithCurrentMonth(currentMonth *MonthRecord) {
	// Append the current months data to the precomputed query month's data
	q.Months = append(q.Months, currentMonth)

	existingNamespaceMounts := make(map[string]*NamespaceRecord)
	// Store the existing namespaces and mounts in the precomputed query for easy access
	for _, monthlyNamespaceRecord := range q.Namespaces {
		existingNamespaceMounts[monthlyNamespaceRecord.NamespaceID] = monthlyNamespaceRecord
	}

	// Get the counts of new clients each mount per namespace in the current month, and increment
	// its total count in the precomputed query. These total values will be visible in the
	// by_namespace grouping in the final response data
	for _, nsRecord := range currentMonth.NewClients.Namespaces {
		namespaceId := nsRecord.NamespaceID

		// If the namespace already exists in the previous months, iterate through the mounts and increment the counts
		if existingNsRecord, ok := existingNamespaceMounts[namespaceId]; ok {
			existingNsRecord.CombineWithMonthlyNamespaceRecord(nsRecord)
		} else {
			// Else just add the new namespace record to the slice in the precomputed query's namespace slice
			q.Namespaces = append(q.Namespaces, nsRecord.ToNamespaceRecord())
		}
	}
}
