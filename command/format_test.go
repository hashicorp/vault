package command

import (
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/jsonutil"
)

var output string

type mockUi struct {
	t          *testing.T
	SampleData string
}

func (m mockUi) Ask(_ string) (string, error) {
	m.t.FailNow()
	return "", nil
}
func (m mockUi) AskSecret(_ string) (string, error) {
	m.t.FailNow()
	return "", nil
}
func (m mockUi) Output(s string) {
	output = s
}
func (m mockUi) Info(s string) {
	m.t.Log(s)
}
func (m mockUi) Error(s string) {
	m.t.Log(s)
}
func (m mockUi) Warn(s string) {
	m.t.Log(s)
}

func TestJsonFormatter(t *testing.T) {
	ui := mockUi{t: t, SampleData: "something"}
	if err := outputWithFormat(ui, "json", nil, ui); err != 0 {
		t.Fatal(err)
	}
	var newUi mockUi
	if err := jsonutil.DecodeJSON([]byte(output), &newUi); err != nil {
		t.Fatal(err)
	}
	if newUi.SampleData != ui.SampleData {
		t.Fatalf(`values not equal ("%s" != "%s")`,
			newUi.SampleData,
			ui.SampleData)
	}
}

func TestYamlFormatter(t *testing.T) {
	ui := mockUi{t: t, SampleData: "something"}
	if err := outputWithFormat(ui, "yaml", nil, ui); err != 0 {
		t.Fatal(err)
	}
	var newUi mockUi
	err := yaml.Unmarshal([]byte(output), &newUi)
	if err != nil {
		t.Fatal(err)
	}
	if newUi.SampleData != ui.SampleData {
		t.Fatalf(`values not equal ("%s" != "%s")`,
			newUi.SampleData,
			ui.SampleData)
	}
}

func TestTableFormatter(t *testing.T) {
	ui := mockUi{t: t}
	s := api.Secret{Data: map[string]interface{}{"k": "something"}}
	if err := outputWithFormat(ui, "table", &s, &s); err != 0 {
		t.Fatal(err)
	}
	if !strings.Contains(output, "something") {
		t.Fatal("did not find 'something'")
	}
}
