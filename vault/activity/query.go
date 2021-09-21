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
	"github.com/hashicorp/vault/sdk/logical"
)

// About 66 bytes per record:
//{"namespace_id":"xxxxx","entities":1234,"non_entity_tokens":1234},
// = approx 7900 namespaces in 512KiB
// So one storage entry is fine (for now).
type NamespaceRecord struct {
	NamespaceID     string `json:"namespace_id"`
	Entities        uint64 `json:"entities"`
	NonEntityTokens uint64 `json:"non_entity_tokens"`
}

type PrecomputedQuery struct {
	StartTime  time.Time
	EndTime    time.Time
	Namespaces []*NamespaceRecord `json:"namespaces"`
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
	err = s.view.Put(ctx, &logical.StorageEntry{
		Key:   path,
		Value: asJson,
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

	filteredList := make([]time.Time, 0, len(startTimes))
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
	contiguous := timeutil.GetMostRecentContiguousMonths(filteredList)
	actualStartTime := contiguous[len(contiguous)-1]

	s.logger.Trace("chose start time", "actualStartTime", actualStartTime, "contiguous", contiguous)

	endTimes, err := s.listEndTimes(ctx, actualStartTime)
	if err != nil {
		return nil, err
	}
	s.logger.Trace("retrieved end times from storage", "endTimes", endTimes)

	// Might happen if there's a race with GC
	if len(endTimes) == 0 {
		s.logger.Warn("missing end times", "start time", actualStartTime)
		return nil, nil
	}
	var actualEndTime time.Time
	for _, t := range endTimes {
		if timeutil.InRange(t, startTime, endTime) {
			if actualEndTime.IsZero() || t.After(actualEndTime) {
				actualEndTime = t
			}
		}
	}
	if actualEndTime.IsZero() {
		s.logger.Warn("no end time in range", "start time", actualStartTime)
		return nil, nil
	}

	path := fmt.Sprintf("%v/%v", actualStartTime.Unix(), actualEndTime.Unix())
	entry, err := s.view.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	p := &PrecomputedQuery{}
	err = json.Unmarshal(entry.Value, p)
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
