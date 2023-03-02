package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/hashicorp/vault/helper/timeutil"
)

// Test_ActivityLog_ComputeCurrentMonthForBillingPeriodInternal creates 3 months of hyperloglogs and fills them with
// overlapping clients. The test calls computeCurrentMonthForBillingPeriodInternal with the current month map having
// some overlap with the previous months. The test then verifies that the results have the correct number of entity and
// non-entity clients. The test also calls computeCurrentMonthForBillingPeriodInternal with an empty current month map,
// and verifies that the results are all 0
func Test_ActivityLog_ComputeCurrentMonthForBillingPeriodInternal(t *testing.T) {
	// populate the first month with clients 1-10
	monthOneHLL := hyperloglog.New()
	// populate the second month with clients 5-15
	monthTwoHLL := hyperloglog.New()
	// populate the third month with clients 10-20
	monthThreeHLL := hyperloglog.New()

	for i := 0; i < 20; i++ {
		clientID := []byte(fmt.Sprintf("client_%d", i))
		if i < 10 {
			monthOneHLL.Insert(clientID)
		}
		if 5 <= i && i < 15 {
			monthTwoHLL.Insert(clientID)
		}
		if 10 <= i && i < 20 {
			monthThreeHLL.Insert(clientID)
		}
	}
	mockHLLGetFunc := func(ctx context.Context, startTime time.Time) (*hyperloglog.Sketch, error) {
		currMonthStart := timeutil.StartOfMonth(time.Now())
		if startTime.Equal(timeutil.MonthsPreviousTo(3, currMonthStart)) {
			return monthThreeHLL, nil
		}
		if startTime.Equal(timeutil.MonthsPreviousTo(2, currMonthStart)) {
			return monthTwoHLL, nil
		}
		if startTime.Equal(timeutil.MonthsPreviousTo(1, currMonthStart)) {
			return monthOneHLL, nil
		}
		return nil, fmt.Errorf("bad start time")
	}

	// Let's add 2 entities exclusive to month 1 (clients 0,1),
	// 2 entities shared by month 1 and 2 (clients 5,6),
	// 2 entities shared by month 2 and 3 (clients 10,11), and
	// 2 entities exclusive to month 3 (15,16). Furthermore, we can add
	// 3 new entities (clients 20,21, and 22).
	entitiesStruct := make(map[string]struct{}, 0)
	entitiesStruct["client_0"] = struct{}{}
	entitiesStruct["client_1"] = struct{}{}
	entitiesStruct["client_5"] = struct{}{}
	entitiesStruct["client_6"] = struct{}{}
	entitiesStruct["client_10"] = struct{}{}
	entitiesStruct["client_11"] = struct{}{}
	entitiesStruct["client_15"] = struct{}{}
	entitiesStruct["client_16"] = struct{}{}
	entitiesStruct["client_20"] = struct{}{}
	entitiesStruct["client_21"] = struct{}{}
	entitiesStruct["client_22"] = struct{}{}

	// We will add 3 nonentity clients from month 1 (clients 2,3,4),
	// 3 shared by months 1 and 2 (7,8,9),
	// 3 shared by months 2 and 3 (12,13,14), and
	// 3 exclusive to month 3 (17,18,19). We will also
	// add 4 new nonentity clients.
	nonEntitiesStruct := make(map[string]struct{}, 0)
	nonEntitiesStruct["client_2"] = struct{}{}
	nonEntitiesStruct["client_3"] = struct{}{}
	nonEntitiesStruct["client_4"] = struct{}{}
	nonEntitiesStruct["client_7"] = struct{}{}
	nonEntitiesStruct["client_8"] = struct{}{}
	nonEntitiesStruct["client_9"] = struct{}{}
	nonEntitiesStruct["client_12"] = struct{}{}
	nonEntitiesStruct["client_13"] = struct{}{}
	nonEntitiesStruct["client_14"] = struct{}{}
	nonEntitiesStruct["client_17"] = struct{}{}
	nonEntitiesStruct["client_18"] = struct{}{}
	nonEntitiesStruct["client_19"] = struct{}{}
	nonEntitiesStruct["client_23"] = struct{}{}
	nonEntitiesStruct["client_24"] = struct{}{}
	nonEntitiesStruct["client_25"] = struct{}{}
	nonEntitiesStruct["client_26"] = struct{}{}

	counts := &processCounts{
		Entities:    entitiesStruct,
		NonEntities: nonEntitiesStruct,
	}

	currentMonthClientsMap := make(map[int64]*processMonth, 1)
	currentMonthClients := &processMonth{
		Counts:     counts,
		NewClients: &processNewClients{Counts: counts},
	}
	// Technially I think currentMonthClientsMap should have the keys as
	// unix timestamps, but for the purposes of the unit test it doesn't
	// matter what the values actually are.
	currentMonthClientsMap[0] = currentMonthClients

	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	endTime := timeutil.StartOfMonth(time.Now())
	startTime := timeutil.MonthsPreviousTo(3, endTime)

	monthRecord, err := a.computeCurrentMonthForBillingPeriodInternal(context.Background(), currentMonthClientsMap, mockHLLGetFunc, startTime, endTime)
	if err != nil {
		t.Fatal(err)
	}

	// We should have 11 entity clients and 16 nonentity clients, and 3 new entity clients
	// and 4 new nonentity clients
	if monthRecord.Counts.EntityClients != 11 {
		t.Fatalf("wrong number of entity clients. Expected 11, got %d", monthRecord.Counts.EntityClients)
	}
	if monthRecord.Counts.NonEntityClients != 16 {
		t.Fatalf("wrong number of non entity clients. Expected 16, got %d", monthRecord.Counts.NonEntityClients)
	}
	if monthRecord.NewClients.Counts.EntityClients != 3 {
		t.Fatalf("wrong number of new entity clients. Expected 3, got %d", monthRecord.NewClients.Counts.EntityClients)
	}
	if monthRecord.NewClients.Counts.NonEntityClients != 4 {
		t.Fatalf("wrong number of new non entity clients. Expected 4, got %d", monthRecord.NewClients.Counts.NonEntityClients)
	}

	// Attempt to compute current month when no records exist
	endTime = time.Now().UTC()
	startTime = timeutil.StartOfMonth(endTime)
	emptyClientsMap := make(map[int64]*processMonth, 0)
	monthRecord, err = a.computeCurrentMonthForBillingPeriodInternal(context.Background(), emptyClientsMap, mockHLLGetFunc, startTime, endTime)
	if err != nil {
		t.Fatalf("failed to compute empty current month, err: %v", err)
	}

	// We should have 0 entity clients, nonentity clients,new entity clients
	// and new nonentity clients
	if monthRecord.Counts.EntityClients != 0 {
		t.Fatalf("wrong number of entity clients. Expected 0, got %d", monthRecord.Counts.EntityClients)
	}
	if monthRecord.Counts.NonEntityClients != 0 {
		t.Fatalf("wrong number of non entity clients. Expected 0, got %d", monthRecord.Counts.NonEntityClients)
	}
	if monthRecord.NewClients.Counts.EntityClients != 0 {
		t.Fatalf("wrong number of new entity clients. Expected 0, got %d", monthRecord.NewClients.Counts.EntityClients)
	}
	if monthRecord.NewClients.Counts.NonEntityClients != 0 {
		t.Fatalf("wrong number of new non entity clients. Expected 0, got %d", monthRecord.NewClients.Counts.NonEntityClients)
	}
}
