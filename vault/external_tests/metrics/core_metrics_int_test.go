package metrics

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/metricsutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

func TestMountTableMetrics(t *testing.T) {
	clusterName := "mycluster"
	conf := &vault.CoreConfig{
		BuiltinRegistry: vault.NewMockBuiltinRegistry(),
		ClusterName:     clusterName,
	}
	cluster := vault.NewTestCluster(t, conf, &vault.TestClusterOptions{
		KeepStandbysSealed:     false,
		HandlerFunc:            vaulthttp.Handler,
		NumCores:               2,
		CoreMetricSinkProvider: testMetricSinkProvider(time.Minute),
	})

	cluster.Start()
	defer cluster.Cleanup()

	// Wait for core to become active
	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	// Verify that the nonlocal logical mount table has 3 entries -- cubbyhole, identity, and kv

	data, err := sysMetricsReq(client, cluster)
	if err != nil {
		t.Fatal(err)
	}

	nonlocalLogicalMountsize, err := gaugeSearchHelper(data, 3)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Mount new kv
	if err = client.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "2",
		},
	}); err != nil {
		t.Fatal(err)
	}

	data, err = sysMetricsReq(client, cluster)
	if err != nil {
		t.Fatal(err)
	}

	nonlocalLogicalMountsizeAfterMount, err := gaugeSearchHelper(data, 4)
	if err != nil {
		t.Errorf(err.Error())
	}

	if nonlocalLogicalMountsizeAfterMount <= nonlocalLogicalMountsize {
		t.Errorf("Mount size does not change after new mount is mounted")
	}
}

func sysMetricsReq(client *api.Client, cluster *vault.TestCluster) (*SysMetricsJSON, error) {
	r := client.NewRequest("GET", "/v1/sys/metrics")
	r.Headers.Set("X-Vault-Token", cluster.RootToken)
	var data SysMetricsJSON
	mountAddResp, err := client.RawRequestWithContext(context.Background(), r)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(mountAddResp.Response.Body)
	if err != nil {
		return nil, err
	}
	if mountAddResp != nil {
		defer mountAddResp.Body.Close()
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, errors.New("failed to unmarshal:" + err.Error())
	}
	return &data, nil
}

func gaugeSearchHelper(data *SysMetricsJSON, expectedValue int) (int, error) {
	foundFlag := false
	tablesize := int(^uint(0) >> 1)
	for _, gauge := range data.Gauges {
		labels := gauge.Labels
		if loc, ok := labels["local"]; ok && loc.(string) == "false" {
			if tp, ok := labels["type"]; ok && tp.(string) == "logical" {
				if gauge.Name == "core.mount_table.num_entries" {
					foundFlag = true
					if err := gaugeConditionCheck("eq", expectedValue, gauge.Value); err != nil {
						return int(^uint(0) >> 1), err
					}
				} else if gauge.Name == "core.mount_table.size" {
					tablesize = gauge.Value
				}
			}
		}
	}
	if !foundFlag {
		return int(^uint(0) >> 1), errors.New("No metrics reported for mount sizes")
	}
	return tablesize, nil
}

func gaugeConditionCheck(comparator string, compareVal int, compareToVal int) error {
	if comparator == "eq" && compareVal != compareToVal {
		return errors.New("equality gauge check for comparison failed")
	}
	return nil
}

func testMetricSinkProvider(gaugeInterval time.Duration) func(string) (*metricsutil.ClusterMetricSink, *metricsutil.MetricsHelper) {
	return func(clusterName string) (*metricsutil.ClusterMetricSink, *metricsutil.MetricsHelper) {
		inm := metrics.NewInmemSink(1000000*time.Hour, 2000000*time.Hour)
		clusterSink := metricsutil.NewClusterMetricSink(clusterName, inm)
		clusterSink.GaugeInterval = gaugeInterval
		return clusterSink, metricsutil.NewMetricsHelper(inm, false)
	}
}

func TestLeaderReElectionMetrics(t *testing.T) {
	clusterName := "mycluster"
	conf := &vault.CoreConfig{
		BuiltinRegistry: vault.NewMockBuiltinRegistry(),
		ClusterName:     clusterName,
	}
	cluster := vault.NewTestCluster(t, conf, &vault.TestClusterOptions{
		KeepStandbysSealed:     false,
		HandlerFunc:            vaulthttp.Handler,
		NumCores:               2,
		CoreMetricSinkProvider: testMetricSinkProvider(time.Minute),
	})

	cluster.Start()
	defer cluster.Cleanup()

	// Wait for core to become active
	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client
	standbyClient := cores[1].Client

	r := client.NewRequest("GET", "/v1/sys/metrics")
	r2 := standbyClient.NewRequest("GET", "/v1/sys/metrics")
	r.Headers.Set("X-Vault-Token", cluster.RootToken)
	r2.Headers.Set("X-Vault-Token", cluster.RootToken)
	respo, err := client.RawRequestWithContext(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(respo.Response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if respo != nil {
		defer respo.Body.Close()
	}
	var data SysMetricsJSON
	var coreLeaderMetric bool = false
	var coreUnsealMetric bool = false
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		t.Fatal("failed to unmarshal:", err)
	}
	for _, gauge := range data.Gauges {
		if gauge.Name == "core.active" {
			coreLeaderMetric = true
			if gauge.Value != 1 {
				t.Errorf("metric incorrectly reports active status")
			}
		}
		if gauge.Name == "core.unsealed" {
			coreUnsealMetric = true
			if gauge.Value != 1 {
				t.Errorf("metric incorrectly reports unseal status of leader")
			}
		}
	}
	if !coreLeaderMetric || !coreUnsealMetric {
		t.Errorf("unseal metric or leader metric are missing")
	}

	err = client.Sys().StepDown()
	if err != nil {
		t.Fatal(err)
	}
	// Wait for core to become active
	vault.TestWaitActive(t, cores[1].Core)

	r = standbyClient.NewRequest("GET", "/v1/sys/metrics")
	r.Headers.Set("X-Vault-Token", cluster.RootToken)
	respo, err = standbyClient.RawRequestWithContext(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	bodyBytes, err = ioutil.ReadAll(respo.Response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		t.Fatal("failed to unmarshal:", err)
	} else {
		coreLeaderMetric = false
		coreUnsealMetric = false
		for _, gauge := range data.Gauges {
			if gauge.Name == "core.active" {
				coreLeaderMetric = true
				if gauge.Value != 1 {
					t.Errorf("metric incorrectly reports active status")
				}
			}
			if gauge.Name == "core.unsealed" {
				coreUnsealMetric = true
				if gauge.Value != 1 {
					t.Errorf("metric incorrectly reports unseal status of leader")
				}
			}
		}
		if !coreLeaderMetric || !coreUnsealMetric {
			t.Errorf("unseal metric or leader metric are missing")
		}
	}
	if respo != nil {
		defer respo.Body.Close()
	}
}

type SysMetricsJSON struct {
	Gauges []GaugeJSON `json:"Gauges"`
}

type GaugeJSON struct {
	Name   string                 `json:"Name"`
	Value  int                    `json:"Value"`
	Labels map[string]interface{} `json:"Labels"`
}
