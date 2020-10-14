package metrics

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/shared-secure-libs/metricsutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

func TestLeaderReElectionMetrics(t *testing.T) {
	inm := metrics.NewInmemSink(time.Minute, time.Minute*10)
	clusterSink := metricsutil.NewClusterMetricSink("mycluster", inm)
	clusterSink.GaugeInterval = time.Second
	conf := &vault.CoreConfig{
		BuiltinRegistry: vault.NewMockBuiltinRegistry(),
		MetricsHelper:   metricsutil.NewMetricsHelper(inm, true),
		MetricSink:      clusterSink,
	}
	cluster := vault.NewTestCluster(t, conf, &vault.TestClusterOptions{
		KeepStandbysSealed: false,
		HandlerFunc:        vaulthttp.Handler,
		NumCores:           2,
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
		if gauge.Name == "core.leader" {
			coreLeaderMetric = true
			if gauge.Value != 1 {
				t.Errorf("metric incorrectly reports leader status")
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
			if gauge.Name == "core.leader" {
				coreLeaderMetric = true
				if gauge.Value != 1 {
					t.Errorf("metric incorrectly reports leader status")
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
	Name  string `json:"Name"`
	Value int    `json:"Value"`
}
