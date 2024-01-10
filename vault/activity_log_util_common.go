// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"google.golang.org/protobuf/proto"
)

type HLLGetter func(ctx context.Context, startTime time.Time) (*hyperloglog.Sketch, error)

// computeCurrentMonthForBillingPeriod computes the current month's data with respect
// to a billing period.
func (a *ActivityLog) computeCurrentMonthForBillingPeriod(ctx context.Context, byMonth map[int64]*processMonth, startTime time.Time, endTime time.Time) (*activity.MonthRecord, error) {
	return a.computeCurrentMonthForBillingPeriodInternal(ctx, byMonth, a.CreateOrFetchHyperlogLog, startTime, endTime)
}

// CreateOrFetchHyperlogLog creates a new hyperlogLog for each startTime (month) if it does not exist in storage.
// hyperlogLog is used here to solve count-distinct problem i.e, to count the number of distinct clients
// In activity log, hyperloglog is a sketch containing clientID's in a given month
func (a *ActivityLog) CreateOrFetchHyperlogLog(ctx context.Context, startTime time.Time) (*hyperloglog.Sketch, error) {
	monthlyHLLPath := fmt.Sprintf("%s%d", distinctClientsBasePath, startTime.Unix())
	hll := hyperloglog.New()
	data, err := a.view.Get(ctx, monthlyHLLPath)
	if err != nil {
		// If there is no hll, we should log the error, as having this fire multiple times
		// is a sign that something is wrong with hll store/get. However, this is not a
		// critical failure (in fact it is expected during the first month rotation after
		// this code is deployed), so we will not throw an error.
		a.logger.Warn("fetch of hyperloglog threw an error at path", monthlyHLLPath, "error", err)
	}
	if data == nil {
		a.logger.Trace("creating hyperloglog ", "path", monthlyHLLPath)
		err = a.StoreHyperlogLog(ctx, startTime, hll)
		if err != nil {
			return hll, fmt.Errorf("error storing hyperloglog at path %s: error %w", monthlyHLLPath, err)
		}
	} else {
		err = hll.UnmarshalBinary(data.Value)
		if err != nil {
			return hll, fmt.Errorf("error unmarshaling hyperloglog at path %s: error %w", monthlyHLLPath, err)
		}
	}
	return hll, nil
}

// StoreHyperlogLog stores the hyperloglog (a sketch containing client IDs) for startTime (month) in storage
func (a *ActivityLog) StoreHyperlogLog(ctx context.Context, startTime time.Time, newHll *hyperloglog.Sketch) error {
	monthlyHLLPath := fmt.Sprintf("%s%d", distinctClientsBasePath, startTime.Unix())
	a.logger.Trace("storing hyperloglog ", "path", monthlyHLLPath)
	marshalledHll, err := newHll.MarshalBinary()
	if err != nil {
		return err
	}
	err = a.view.Put(ctx, &logical.StorageEntry{
		Key:   monthlyHLLPath,
		Value: marshalledHll,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *ActivityLog) computeCurrentMonthForBillingPeriodInternal(ctx context.Context, byMonth map[int64]*processMonth, hllGetFunc HLLGetter, startTime time.Time, endTime time.Time) (*activity.MonthRecord, error) {
	if timeutil.IsCurrentMonth(startTime, a.clock.Now().UTC()) {
		monthlyComputation := a.transformMonthBreakdowns(byMonth)
		if len(monthlyComputation) > 1 {
			a.logger.Warn("monthly in-memory activitylog computation returned multiple months of data", "months returned", len(byMonth))
		}
		if len(monthlyComputation) > 0 {
			return monthlyComputation[0], nil
		}
	}
	// Fetch all hyperloglogs for months from startMonth to endMonth. If a month doesn't have an associated
	// hll, warn and continue.

	// hllMonthlyTimestamp is the start time of the month corresponding to which a hyperloglog of that month's
	// client data is stored. The path at which the hyperloglog for a month is stored containes this timestamp.
	hllMonthlyTimestamp := timeutil.StartOfMonth(startTime)
	billingPeriodHLL := hyperloglog.New()
	for hllMonthlyTimestamp.Before(timeutil.StartOfMonth(endTime)) {
		monthSketch, err := hllGetFunc(ctx, hllMonthlyTimestamp)
		// If there's an error with the hyperloglog fetch, we should still deduplicate on
		// the hlls that we have so we will warn that we couldn't find a hll for the month
		// and continue.
		if err != nil {
			a.logger.Warn("no hyperloglog associated with timestamp", "timestamp", hllMonthlyTimestamp)
			hllMonthlyTimestamp = timeutil.StartOfNextMonth(hllMonthlyTimestamp)
			continue
		}
		// Union the monthly hll into the billing period's hll
		err = billingPeriodHLL.Merge(monthSketch)
		if err != nil {
			// In this case we can't afford to fail silently. Since this error indicates
			// data corruption, we should not try to do any further deduplication
			return nil, err
		}
		hllMonthlyTimestamp = timeutil.StartOfNextMonth(hllMonthlyTimestamp)
	}
	// There's at most one month of data here. We should validate this assumption explicitly
	if len(byMonth) > 1 {
		return nil, errors.New(fmt.Sprintf("multiple months of data found in partial month's client count breakdowns: %+v\n", byMonth))
	}

	activityTypes := []string{entityActivityType, nonEntityTokenActivityType, secretSyncActivityType}

	// Now we will add the clients for the current month to a copy of the billing period's hll to
	// see how the cardinality grows.
	hllByType := make(map[string]*hyperloglog.Sketch, len(activityTypes))
	totalByType := make(map[string]int, len(activityTypes))
	for _, typ := range activityTypes {
		hllByType[typ] = billingPeriodHLL.Clone()
	}

	for _, month := range byMonth {
		if month.NewClients == nil || month.NewClients.Counts == nil || month.Counts == nil {
			return nil, errors.New("malformed current month used to calculate current month's activity")
		}

		for _, typ := range activityTypes {
			// Note that the following calculations assume that all clients seen are currently in
			// the NewClients section of byMonth. It is best to explicitly check this, just verify
			// our assumptions about the passed in byMonth argument.
			if month.Counts.countByType(typ) != month.NewClients.Counts.countByType(typ) {
				return nil, errors.New("current month clients cache assumes billing period")
			}
			for clientID := range month.NewClients.Counts.clientsByType(typ) {
				// All the clients for the current month are in the newClients section, initially.
				// We need to deduplicate these clients across the billing period by adding them
				// into the billing period hyperloglogs.
				hllByType[typ].Insert([]byte(clientID))
				totalByType[typ] += 1
			}
		}
	}
	currentMonthNewByType := make(map[string]int, len(activityTypes))
	for _, typ := range activityTypes {
		// The number of new entities for the current month is approximately the size of the hll with
		// the current month's entities minus the size of the initial billing period hll.
		currentMonthNewByType[typ] = int(hllByType[typ].Estimate() - billingPeriodHLL.Estimate())
	}

	return &activity.MonthRecord{
		Timestamp: timeutil.StartOfMonth(endTime).UTC().Unix(),
		NewClients: &activity.NewClientRecord{Counts: &activity.CountsRecord{
			EntityClients:    currentMonthNewByType[entityActivityType],
			NonEntityClients: currentMonthNewByType[nonEntityTokenActivityType],
			SecretSyncs:      currentMonthNewByType[secretSyncActivityType],
		}},
		Counts: &activity.CountsRecord{
			EntityClients:    totalByType[entityActivityType],
			NonEntityClients: totalByType[nonEntityTokenActivityType],
			SecretSyncs:      totalByType[secretSyncActivityType],
		},
	}, nil
}

// sortALResponseNamespaces sorts the namespaces for activity log responses.
func (a *ActivityLog) sortALResponseNamespaces(byNamespaceResponse []*ResponseNamespace) {
	sort.Slice(byNamespaceResponse, func(i, j int) bool {
		return byNamespaceResponse[i].Counts.Clients > byNamespaceResponse[j].Counts.Clients
	})
}

// transformALNamespaceBreakdowns takes the namespace breakdowns stored in the intermediary
// struct used in precomputation segment traversal and to store the current month data and
// reorganizes it into query structs. This helper is used by the partial month endpoint so as to
// not have to maintain two separate response data computations for two separate APIs.
func (a *ActivityLog) transformALNamespaceBreakdowns(nsData map[string]*processByNamespace) []*activity.NamespaceRecord {
	byNamespace := make([]*activity.NamespaceRecord, 0)
	for nsID, ns := range nsData {

		nsRecord := activity.NamespaceRecord{
			NamespaceID:     nsID,
			Entities:        uint64(ns.Counts.countByType(entityActivityType)),
			NonEntityTokens: uint64(ns.Counts.countByType(nonEntityTokenActivityType)),
			SecretSyncs:     uint64(ns.Counts.countByType(secretSyncActivityType)),
			Mounts:          a.transformActivityLogMounts(ns.Mounts),
		}
		byNamespace = append(byNamespace, &nsRecord)
	}
	return byNamespace
}

// limitNamespacesInALResponse will truncate the number of namespaces shown in the activity
// endpoints to the number specified in limitNamespaces (the API filtering parameter)
func (a *ActivityLog) limitNamespacesInALResponse(byNamespaceResponse []*ResponseNamespace, limitNamespaces int) (int, int, []*ResponseNamespace) {
	if limitNamespaces > len(byNamespaceResponse) {
		limitNamespaces = len(byNamespaceResponse)
	}
	byNamespaceResponse = byNamespaceResponse[:limitNamespaces]
	// recalculate total entities and tokens
	totalEntities := 0
	totalTokens := 0
	for _, namespaceData := range byNamespaceResponse {
		totalEntities += namespaceData.Counts.DistinctEntities
		totalTokens += namespaceData.Counts.NonEntityTokens
	}
	return totalEntities, totalTokens, byNamespaceResponse
}

// transformActivityLogMounts is a helper used to reformat data for transformMonthlyNamespaceBreakdowns.
// For more details, please see the function comment for transformMonthlyNamespaceBreakdowns
func (a *ActivityLog) transformActivityLogMounts(mts map[string]*processMount) []*activity.MountRecord {
	mounts := make([]*activity.MountRecord, 0)
	for mountAccessor, mountCounts := range mts {
		mount := activity.MountRecord{
			MountPath: a.mountAccessorToMountPath(mountAccessor),
			Counts:    mountCounts.Counts.toCountsRecord(),
		}
		mounts = append(mounts, &mount)
	}
	return mounts
}

// sortActivityLogMonthsResponse contains the sorting logic for the months
// portion of the activity log response.
func (a *ActivityLog) sortActivityLogMonthsResponse(months []*ResponseMonth) {
	// Sort the months in ascending order of timestamps
	sort.Slice(months, func(i, j int) bool {
		firstTimestamp, errOne := time.Parse(time.RFC3339, months[i].Timestamp)
		secondTimestamp, errTwo := time.Parse(time.RFC3339, months[j].Timestamp)
		if errOne == nil && errTwo == nil {
			return firstTimestamp.Before(secondTimestamp)
		}
		// Keep the nondeterministic ordering in storage
		a.logger.Error("unable to parse activity log timestamps", "timestamp",
			months[i].Timestamp, "error", errOne, "timestamp", months[j].Timestamp, "error", errTwo)
		return i < j
	})

	// Within each month sort everything by descending order of activity
	for _, month := range months {
		sort.Slice(month.Namespaces, func(i, j int) bool {
			return month.Namespaces[i].Counts.Clients > month.Namespaces[j].Counts.Clients
		})

		for _, ns := range month.Namespaces {
			sort.Slice(ns.Mounts, func(i, j int) bool {
				return ns.Mounts[i].Counts.Clients > ns.Mounts[j].Counts.Clients
			})
		}

		sort.Slice(month.NewClients.Namespaces, func(i, j int) bool {
			return month.NewClients.Namespaces[i].Counts.Clients > month.NewClients.Namespaces[j].Counts.Clients
		})

		for _, ns := range month.NewClients.Namespaces {
			sort.Slice(ns.Mounts, func(i, j int) bool {
				return ns.Mounts[i].Counts.Clients > ns.Mounts[j].Counts.Clients
			})
		}
	}
}

const (
	noMountAccessor = "no mount accessor (pre-1.10 upgrade?)"
	deletedMountFmt = "deleted mount; accessor %q"
)

// mountAccessorToMountPath transforms the mount accessor to the mount path
// returns a placeholder string if the mount accessor is empty or deleted
func (a *ActivityLog) mountAccessorToMountPath(mountAccessor string) string {
	var displayPath string
	if mountAccessor == "" {
		displayPath = noMountAccessor
	} else {
		valResp := a.core.router.ValidateMountByAccessor(mountAccessor)
		if valResp == nil {
			displayPath = fmt.Sprintf(deletedMountFmt, mountAccessor)
		} else {
			displayPath = valResp.MountPath
			if !strings.HasSuffix(displayPath, "/") {
				displayPath += "/"
			}
		}
	}
	return displayPath
}

type singleTypeSegmentReader struct {
	basePath         string
	startTime        time.Time
	paths            []string
	currentPathIndex int
	a                *ActivityLog
}
type segmentReader struct {
	tokens   *singleTypeSegmentReader
	entities *singleTypeSegmentReader
}

// SegmentReader is an interface that provides methods to read tokens and entities in order
type SegmentReader interface {
	ReadToken(ctx context.Context) (*activity.TokenCount, error)
	ReadEntity(ctx context.Context) (*activity.EntityActivityLog, error)
}

func (a *ActivityLog) NewSegmentFileReader(ctx context.Context, startTime time.Time) (SegmentReader, error) {
	entities, err := a.newSingleTypeSegmentReader(ctx, startTime, activityEntityBasePath)
	if err != nil {
		return nil, err
	}
	tokens, err := a.newSingleTypeSegmentReader(ctx, startTime, activityTokenBasePath)
	if err != nil {
		return nil, err
	}
	return &segmentReader{entities: entities, tokens: tokens}, nil
}

func (a *ActivityLog) newSingleTypeSegmentReader(ctx context.Context, startTime time.Time, prefix string) (*singleTypeSegmentReader, error) {
	basePath := prefix + fmt.Sprint(startTime.Unix()) + "/"
	pathList, err := a.view.List(ctx, basePath)
	if err != nil {
		return nil, err
	}
	return &singleTypeSegmentReader{
		basePath:         basePath,
		startTime:        startTime,
		paths:            pathList,
		currentPathIndex: 0,
		a:                a,
	}, nil
}

func (s *singleTypeSegmentReader) nextValue(ctx context.Context, out proto.Message) error {
	var raw *logical.StorageEntry
	var path string
	for raw == nil {
		if s.currentPathIndex >= len(s.paths) {
			return io.EOF
		}
		path = s.paths[s.currentPathIndex]
		// increment the index to continue iterating for the next read call, even if an error occurs during this call
		s.currentPathIndex++
		var err error
		raw, err = s.a.view.Get(ctx, s.basePath+path)
		if err != nil {
			return err
		}
		if raw == nil {
			s.a.logger.Warn("expected log segment file has been deleted", "startTime", s.startTime, "segmentPath", path)
		}
	}
	err := proto.Unmarshal(raw.Value, out)
	if err != nil {
		return fmt.Errorf("unable to parse segment file %v%v: %w", s.basePath, path, err)
	}
	return nil
}

// ReadToken reads a token from the segment
// If there is none available, then the error will be io.EOF
func (e *segmentReader) ReadToken(ctx context.Context) (*activity.TokenCount, error) {
	out := &activity.TokenCount{}
	err := e.tokens.nextValue(ctx, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReadEntity reads an entity from the segment
// If there is none available, then the error will be io.EOF
func (e *segmentReader) ReadEntity(ctx context.Context) (*activity.EntityActivityLog, error) {
	out := &activity.EntityActivityLog{}
	err := e.entities.nextValue(ctx, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
