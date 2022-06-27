package configutil

import (
	"fmt"
	"github.com/armon/go-metrics"
	"strings"
	"testing"
	"time"
)

type mocIDogStatsd struct {
	tagsHistory [][]string
}

func (m *mocIDogStatsd) SetTags(tags []string) {
	m.tagsHistory = append(m.tagsHistory, tags)
}

func (m *mocIDogStatsd) EnableHostNamePropagation() {
}

func (m mocIDogStatsd) SetGauge(key []string, val float32) {
}

func (m mocIDogStatsd) IncrCounter(key []string, val float32) {
}

func (m mocIDogStatsd) EmitKey(key []string, val float32) {
}

func (m mocIDogStatsd) AddSample(key []string, val float32) {
}

func (m mocIDogStatsd) SetGaugeWithLabels(key []string, val float32, labels []metrics.Label) {
}

func (m mocIDogStatsd) IncrCounterWithLabels(key []string, val float32, labels []metrics.Label) {
}

func (m mocIDogStatsd) AddSampleWithLabels(key []string, val float32, labels []metrics.Label) {
}

func (m mocIDogStatsd) Shutdown() {
}

type mockUi struct {
	t        *testing.T
	warnings []string
	errors   []string
	infos    []string
}

func (m mockUi) Ask(_ string) (string, error) {
	m.t.FailNow()
	return "", nil
}

func (m mockUi) AskSecret(_ string) (string, error) {
	m.t.FailNow()
	return "", nil
}
func (m *mockUi) Output(s string) {}
func (m *mockUi) Info(s string) {
	m.infos = append(m.infos, s)
	m.t.Log(s)
}
func (m *mockUi) Error(s string) {
	m.errors = append(m.errors, s)
	m.t.Log(s)
}
func (m *mockUi) Warn(s string) {
	m.warnings = append(m.warnings, s)
	m.t.Log(s)
}

func TestDatadogSink_BadConnection(t *testing.T) {
	ui := mockUi{t: t}
	datadog := NewDatadogSink("", "", &ui)
	datadog.creator = func(addr string, hostName string) (dogStatsdSink, error) {
		return nil, fmt.Errorf("can not connect")
	}

	if datadog == nil {
		t.Fatalf("result can not be nil")
	}

	_ = datadog.getSink()
	_ = datadog.getSink()
	_ = datadog.getSink()

	if len(ui.warnings) == 0 {
		t.Fatalf("no warnings logged")
	}

	if len(ui.warnings) > 1 || len(ui.errors) > 0 || len(ui.infos) > 0 {
		t.Fatalf("excess logging")
	}

	if !strings.HasPrefix(ui.warnings[0], "failed to connect to datadog:") {
		t.Fatalf("incorrect message logged")
	}

	interval, err := time.ParseDuration("-15m")
	if err != nil {
		t.Fatal(err)
	}
	oldTime := datadog.attemptedToConnectAt.Add(interval)
	datadog.attemptedToConnectAt = &oldTime

	_ = datadog.getSink()
	_ = datadog.getSink()
	_ = datadog.getSink()

	if len(ui.warnings) < 2 {
		t.Fatalf("no warnings logged")
	}

	if len(ui.warnings) > 2 || len(ui.errors) > 0 || len(ui.infos) > 0 {
		t.Fatalf("excess logging")
	}

	if !strings.HasPrefix(ui.warnings[1], "failed to connect to datadog:") {
		t.Fatalf("incorrect message logged")
	}
}

func TestDatadogSink_Concurrency(t *testing.T) {
	ui := mockUi{t: t}
	statsd := mocIDogStatsd{}

	datadog := NewDatadogSink("", "", &ui)

	i := 0
	datadog.creator = func(addr string, hostName string) (dogStatsdSink, error) {
		i++
		if i < 2 {
			return nil, fmt.Errorf("bad connection")
		}
		return &statsd, nil
	}

	if datadog == nil {
		t.Fatalf("result can not be nil")
	}

	datadog.SetTags([]string{"tag1", "tag2"})
	if len(datadog.tags) != 2 {
		t.Fatalf("")
	}

	datadog.SetTags([]string{"tag3"})

	interval, err := time.ParseDuration("-15m")
	if err != nil {
		t.Fatal(err)
	}
	oldTime := datadog.attemptedToConnectAt.Add(interval)
	datadog.attemptedToConnectAt = &oldTime

	datadog.SetTags([]string{"tag4"})

	if len(datadog.tags) != 0 {
		t.Fatalf("didn't reset cached tags")
	}

	if len(statsd.tagsHistory) != 2 {
		t.Fatalf("didn't send all tags to server")
	}
}
