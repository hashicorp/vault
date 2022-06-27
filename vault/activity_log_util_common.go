package vault

import (
	"sort"
	"time"

	"github.com/hashicorp/vault/vault/activity"
)

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
	for mountpath, mountCounts := range mts {
		mount := activity.MountRecord{
			MountPath: mountpath,
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

// TODO
// computeCurrentMonthForBillingPeriod computes the current month's data with respect
// to a billing period. This function is currently a stub with the bare minimum amount
// of data to get the pre-existing tests to pass. It will be filled out in a separate PR
// and this comment will be removed.
func (a *ActivityLog) computeCurrentMonthForBillingPeriod(byMonth map[int64]*processMonth, startTime time.Time, endTime time.Time) (*activity.MonthRecord, error) {
	return &activity.MonthRecord{
		NewClients: &activity.NewClientRecord{Counts: &activity.CountsRecord{EntityClients: 0, NonEntityClients: 0}},
		Counts:     &activity.CountsRecord{EntityClients: 0, NonEntityClients: 0},
	}, nil
}
