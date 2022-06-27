package configutil

import (
	"strings"
	"testing"
)

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
	m.errors = append(m.infos, s)
	m.t.Log(s)
}
func (m *mockUi) Warn(s string) {
	m.warnings = append(m.infos, s)
	m.t.Log(s)
}

func TestDatadogSink(t *testing.T) {
	ui := mockUi{t: t}
	datadog := NewDatadogSink("", "", &ui)

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
}
