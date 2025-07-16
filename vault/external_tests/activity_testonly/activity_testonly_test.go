// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package activity_testonly

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

var allClientTypeTestCases = []struct {
	clientType       string
	topLevelJSONKey  string
	responseCountsFn func(r vault.ResponseCounts) int
}{
	{
		clientType:      vault.ACMEActivityType,
		topLevelJSONKey: "acme_clients",
		responseCountsFn: func(r vault.ResponseCounts) int {
			return r.ACMEClients
		},
	},
	{
		clientType:      "secret-sync",
		topLevelJSONKey: "secret_syncs",
		responseCountsFn: func(r vault.ResponseCounts) int {
			return r.SecretSyncs
		},
	},
	{
		clientType:      "entity",
		topLevelJSONKey: "entity_clients",
		responseCountsFn: func(r vault.ResponseCounts) int {
			return r.EntityClients
		},
	},
	{
		clientType:      "non-entity-token",
		topLevelJSONKey: "non_entity_clients",
		responseCountsFn: func(r vault.ResponseCounts) int {
			return r.NonEntityClients
		},
	},
}

func getMonthsData(t *testing.T, resp *api.Secret) []vault.ResponseMonth {
	t.Helper()
	monthsRaw, ok := resp.Data["months"]
	require.True(t, ok)
	monthsResponse := make([]vault.ResponseMonth, 0)
	err := mapstructure.Decode(monthsRaw, &monthsResponse)
	require.NoError(t, err)
	return monthsResponse
}

func getNamespaceData(t *testing.T, resp *api.Secret) []vault.ResponseNamespace {
	t.Helper()
	nsRaw, ok := resp.Data["by_namespace"]
	require.True(t, ok)
	nsResponse := make([]vault.ResponseNamespace, 0)
	err := mapstructure.Decode(nsRaw, &nsResponse)
	require.NoError(t, err)
	return nsResponse
}

func getTotals(t *testing.T, resp *api.Secret) vault.ResponseCounts {
	t.Helper()
	totalRaw, ok := resp.Data["total"]
	require.True(t, ok)
	total := vault.ResponseCounts{}
	err := mapstructure.Decode(totalRaw, &total)
	require.NoError(t, err)
	return total
}

// Test_ActivityLogCurrentMonth_Response runs for each client type. The subtest
// creates 10 clients of the type and verifies that the activity log partial
// month response returns 10 clients of that type at every level of the response
// hierarchy
func Test_ActivityLogCurrentMonth_Response(t *testing.T) {
	t.Parallel()

	for _, tc := range allClientTypeTestCases {
		tc := tc
		t.Run(tc.clientType, func(t *testing.T) {
			t.Parallel()
			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client
			_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
				"enabled": "enable",
			})
			_, err = clientcountutil.NewActivityLogData(client).
				NewCurrentMonthData().
				NewClientsSeen(10, clientcountutil.WithClientType(tc.clientType)).
				Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES)
			require.NoError(t, err)

			resp, err := client.Logical().Read("sys/internal/counters/activity/monthly")
			require.NoError(t, err)

			clientsOfType, ok := resp.Data[tc.topLevelJSONKey]
			require.True(t, ok)
			require.Equal(t, json.Number("10"), clientsOfType)
			clients, ok := resp.Data["clients"]
			require.True(t, ok)
			require.Equal(t, json.Number("10"), clients)

			byNamespace := getNamespaceData(t, resp)
			require.Equal(t, 10, tc.responseCountsFn(byNamespace[0].Counts))
			require.Equal(t, 10, tc.responseCountsFn(*byNamespace[0].Mounts[0].Counts))
			require.Equal(t, 10, byNamespace[0].Counts.Clients)
			require.Equal(t, 10, byNamespace[0].Mounts[0].Counts.Clients)

			byMonth := getMonthsData(t, resp)
			require.Equal(t, 10, tc.responseCountsFn(*byMonth[0].NewClients.Counts))
			require.Equal(t, 10, tc.responseCountsFn(*byMonth[0].Counts))
			require.Equal(t, 10, tc.responseCountsFn(byMonth[0].Namespaces[0].Counts))
			require.Equal(t, 10, tc.responseCountsFn(*byMonth[0].Namespaces[0].Mounts[0].Counts))
			require.Equal(t, 10, byMonth[0].NewClients.Counts.Clients)
			require.Equal(t, 10, byMonth[0].Counts.Clients)
			require.Equal(t, 10, byMonth[0].Namespaces[0].Counts.Clients)
			require.Equal(t, 10, byMonth[0].Namespaces[0].Mounts[0].Counts.Clients)
		})
	}
}

// Test_ActivityLog_Deduplication runs for all client types. The subtest
// verifies that the clients of that type are deduplicated across months. The
// test creates 10 clients and repeats those clients in later months, then also
// registers 3 and then 2 new clients. The test verifies that the total number
// of clients is 15 (10 + 2 + 3), ensuring that the duplicates are not included
func Test_ActivityLog_Deduplication(t *testing.T) {
	t.Parallel()
	for _, tc := range allClientTypeTestCases {
		tc := tc
		t.Run(tc.clientType, func(t *testing.T) {
			t.Parallel()
			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client
			_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
				"enabled": "enable",
			})
			_, err = clientcountutil.NewActivityLogData(client).
				NewPreviousMonthData(3).
				NewClientsSeen(10, clientcountutil.WithClientType(tc.clientType)).
				NewPreviousMonthData(2).
				RepeatedClientsSeen(4, clientcountutil.WithClientType(tc.clientType)).
				NewClientsSeen(3, clientcountutil.WithClientType(tc.clientType)).
				NewPreviousMonthData(1).
				RepeatedClientsSeen(5, clientcountutil.WithClientType(tc.clientType)).
				NewClientsSeen(2, clientcountutil.WithClientType(tc.clientType)).
				Write(context.Background(), generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES)
			require.NoError(t, err)

			now := time.Now().UTC()
			resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", map[string][]string{
				"end_time":   {timeutil.StartOfMonth(now).Format(time.RFC3339)},
				"start_time": {timeutil.StartOfMonth(timeutil.MonthsPreviousTo(4, now)).Format(time.RFC3339)},
			},
			)
			require.NoError(t, err)

			total := getTotals(t, resp)
			require.Equal(t, 15, tc.responseCountsFn(total))
			require.Equal(t, 15, total.Clients)
		})
	}
}

// getJSONExport is used to fetch activity export records using json format.
// The records will be returned as a map keyed by client ID.
func getJSONExport(t *testing.T, client *api.Client, startTime time.Time, now time.Time) (map[string]vault.ActivityLogExportRecord, error) {
	t.Helper()

	resp, err := client.Logical().ReadRawWithData("sys/internal/counters/activity/export", map[string][]string{
		"start_time": {startTime.Format(time.RFC3339)},
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
		"format":     {"json"},
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(contents)
	decoder := json.NewDecoder(buf)
	clients := make(map[string]vault.ActivityLogExportRecord)

	for {
		if !decoder.More() {
			break
		}

		var record vault.ActivityLogExportRecord
		err := decoder.Decode(&record)
		if err != nil {
			return nil, err
		}

		clients[record.ClientID] = record
	}

	return clients, nil
}

// getCSVExport fetches activity export records using csv format. All flattened
// map and slice fields will be unflattened so that the a proper ActivityLogExportRecord
// can be formed. The records will returned as a map keyed by client ID.
func getCSVExport(t *testing.T, client *api.Client, startTime time.Time, now time.Time) (map[string]vault.ActivityLogExportRecord, error) {
	t.Helper()

	boolFields := map[string]struct{}{
		"local_entity_alias": {},
	}

	mapFields := map[string]struct{}{
		"entity_alias_custom_metadata": {},
		"entity_alias_metadata":        {},
		"entity_metadata":              {},
	}

	sliceFields := map[string]struct{}{
		"entity_group_ids": {},
		"policies":         {},
	}

	resp, err := client.Logical().ReadRawWithData("sys/internal/counters/activity/export", map[string][]string{
		"start_time": {startTime.Format(time.RFC3339)},
		"end_time":   {timeutil.EndOfMonth(now).Format(time.RFC3339)},
		"format":     {"csv"},
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	csvRdr := csv.NewReader(resp.Body)
	clients := make(map[string]vault.ActivityLogExportRecord)

	csvRecords, err := csvRdr.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(csvRecords) == 0 {
		return clients, nil
	}

	csvHeader := csvRecords[0]

	// skip initial row as it is header
	for rowIdx := 1; rowIdx < len(csvRecords); rowIdx++ {
		baseRecord := vault.ActivityLogExportRecord{
			Policies:                  []string{},
			EntityMetadata:            map[string]string{},
			EntityAliasMetadata:       map[string]string{},
			EntityAliasCustomMetadata: map[string]string{},
			EntityGroupIDs:            []string{},
		}

		recordMap := make(map[string]interface{})

		// create base map
		err = mapstructure.Decode(baseRecord, &recordMap)
		if err != nil {
			return nil, err
		}

		for columnIdx, columnName := range csvHeader {
			val := csvRecords[rowIdx][columnIdx]

			// determine if column has been flattened
			columnNameParts := strings.SplitN(columnName, ".", 2)

			if len(columnNameParts) > 1 {
				prefix := columnNameParts[0]

				if _, ok := mapFields[prefix]; ok {
					m := recordMap[prefix]

					// ignore empty CSV column value
					if val != "" {
						m.(map[string]string)[columnNameParts[1]] = val
						recordMap[prefix] = m
					}
				} else if _, ok := sliceFields[prefix]; ok {
					s := recordMap[prefix]

					// ignore empty CSV column value
					if val != "" {
						s = append(s.([]string), val)
						recordMap[prefix] = s
					}
				} else {
					t.Fatalf("unexpected CSV field: %q", columnName)
				}
			} else if _, ok := boolFields[columnName]; ok {
				recordMap[columnName], err = strconv.ParseBool(val)
				if err != nil {
					return nil, err
				}
			} else {
				recordMap[columnName] = val
			}
		}

		var record vault.ActivityLogExportRecord
		err = mapstructure.Decode(recordMap, &record)
		if err != nil {
			return nil, err
		}

		clients[record.ClientID] = record
	}

	return clients, nil
}

// Test_ActivityLog_Export_Sudo ensures that the export API is only accessible via
// a root token or a token with a sudo policy.
func Test_ActivityLog_Export_Sudo(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)
	t.Parallel()

	now := time.Now().UTC()
	var err error

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	_, err = client.Logical().Write("sys/internal/counters/config", map[string]interface{}{
		"enabled": "enable",
	})
	require.NoError(t, err)

	rootToken := client.Token()

	_, err = clientcountutil.NewActivityLogData(client).
		NewCurrentMonthData().
		NewClientsSeen(10).
		Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES)

	require.NoError(t, err)

	// Ensure access via root token
	startTime := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(1, now))
	clients, err := getJSONExport(t, client, startTime, now)
	require.NoError(t, err)
	require.Len(t, clients, 10)

	client.Sys().PutPolicy("non-sudo-export", `
path "sys/internal/counters/activity/export" {
	capabilities = ["read"]
}
	`)

	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"non-sudo-export"},
	})
	require.NoError(t, err)

	nonSudoToken := secret.Auth.ClientToken
	client.SetToken(nonSudoToken)

	// Ensure no access via token without sudo access
	clients, err = getJSONExport(t, client, startTime, now)
	require.ErrorContains(t, err, "permission denied")

	client.SetToken(rootToken)
	client.Sys().PutPolicy("sudo-export", `
path "sys/internal/counters/activity/export" {
	capabilities = ["read", "sudo"]
}
	`)

	secret, err = client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"sudo-export"},
	})
	require.NoError(t, err)

	sudoToken := secret.Auth.ClientToken
	client.SetToken(sudoToken)

	// Ensure access via token with sudo access
	clients, err = getJSONExport(t, client, startTime, now)
	require.NoError(t, err)
	require.Len(t, clients, 10)
}
