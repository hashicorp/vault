// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
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
	if timeutil.IsCurrentMonth(startTime, time.Now().UTC()) {
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

	// Now we will add the clients for the current month to a copy of the billing period's hll to
	// see how the cardinality grows.
	billingPeriodHLLWithCurrentMonthEntityClients := billingPeriodHLL.Clone()
	billingPeriodHLLWithCurrentMonthNonEntityClients := billingPeriodHLL.Clone()

	// There's at most one month of data here. We should validate this assumption explicitly
	if len(byMonth) > 1 {
		return nil, errors.New(fmt.Sprintf("multiple months of data found in partial month's client count breakdowns: %+v\n", byMonth))
	}

	totalEntities := 0
	totalNonEntities := 0
	for _, month := range byMonth {

		if month.NewClients == nil || month.NewClients.Counts == nil || month.Counts == nil {
			return nil, errors.New("malformed current month used to calculate current month's activity")
		}

		// Note that the following calculations assume that all clients seen are currently in
		// the NewClients section of byMonth. It is best to explicitly check this, just verify
		// our assumptions about the passed in byMonth argument.
		if len(month.Counts.Entities) != len(month.NewClients.Counts.Entities) ||
			len(month.Counts.NonEntities) != len(month.NewClients.Counts.NonEntities) {
			return nil, errors.New("current month clients cache assumes billing period")
		}

		// All the clients for the current month are in the newClients section, initially.
		// We need to deduplicate these clients across the billing period by adding them
		// into the billing period hyperloglogs.
		entities := month.NewClients.Counts.Entities
		nonEntities := month.NewClients.Counts.NonEntities
		if entities != nil {
			for entityID := range entities {
				billingPeriodHLLWithCurrentMonthEntityClients.Insert([]byte(entityID))
				totalEntities += 1
			}
		}
		if nonEntities != nil {
			for nonEntityID := range nonEntities {
				billingPeriodHLLWithCurrentMonthNonEntityClients.Insert([]byte(nonEntityID))
				totalNonEntities += 1
			}
		}
	}
	// The number of new entities for the current month is approximately the size of the hll with
	// the current month's entities minus the size of the initial billing period hll.
	currentMonthNewEntities := billingPeriodHLLWithCurrentMonthEntityClients.Estimate() - billingPeriodHLL.Estimate()
	currentMonthNewNonEntities := billingPeriodHLLWithCurrentMonthNonEntityClients.Estimate() - billingPeriodHLL.Estimate()
	return &activity.MonthRecord{
		Timestamp:  timeutil.StartOfMonth(endTime).UTC().Unix(),
		NewClients: &activity.NewClientRecord{Counts: &activity.CountsRecord{EntityClients: int(currentMonthNewEntities), NonEntityClients: int(currentMonthNewNonEntities)}},
		Counts:     &activity.CountsRecord{EntityClients: totalEntities, NonEntityClients: totalNonEntities},
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
			Entities:        uint64(len(ns.Counts.Entities)),
			NonEntityTokens: uint64(len(ns.Counts.NonEntities) + int(ns.Counts.Tokens)),
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
			Counts: &activity.CountsRecord{
				EntityClients:    len(mountCounts.Counts.Entities),
				NonEntityClients: len(mountCounts.Counts.NonEntities) + int(mountCounts.Counts.Tokens),
			},
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
