package command

import (
	"os"
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
func (m mockUi) Output(s string) { output = s }
func (m mockUi) Info(s string)   { m.t.Log(s) }
func (m mockUi) Error(s string)  { m.t.Log(s) }
func (m mockUi) Warn(s string)   { m.t.Log(s) }

func TestJsonFormatter(t *testing.T) {
	os.Setenv(EnvVaultFormat, "json")
	ui := mockUi{t: t, SampleData: "something"}
	if err := outputWithFormat(ui, nil, ui); err != 0 {
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
	os.Setenv(EnvVaultFormat, "yaml")
	ui := mockUi{t: t, SampleData: "something"}
	if err := outputWithFormat(ui, nil, ui); err != 0 {
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
	os.Setenv(EnvVaultFormat, "table")
	ui := mockUi{t: t}
	s := api.Secret{Data: map[string]interface{}{"k": "something"}}
	if err := outputWithFormat(ui, &s, &s); err != 0 {
		t.Fatal(err)
	}
	if !strings.Contains(output, "something") {
		t.Fatal("did not find 'something'")
	}
}

func Test_Format_Parsing(t *testing.T) {
	defer func() {
		os.Setenv(EnvVaultCLINoColor, "")
		os.Setenv(EnvVaultFormat, "")
	}()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"format",
			[]string{"-format", "json"},
			"{",
			0,
		},
		{
			"format_bad",
			[]string{"-format", "nope-not-real"},
			"Invalid output format",
			1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			client, closer := testVaultServer(t)
			defer closer()

			// Login with the token so we can renew-self.
			token, _ := testTokenAndAccessor(t, client)
			client.SetToken(token)

			ui, cmd := testTokenRenewCommand(t)
			cmd.client = client

			tc.args = setupEnv(tc.args)

			code := cmd.Run(tc.args)
			if code != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, tc.out) {
				t.Errorf("expected %q to contain %q", combined, tc.out)
			}
		})
	}
}
