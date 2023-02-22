package command

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

type mockUi struct {
	t          *testing.T
	SampleData string
	outputData *string
}

func (m mockUi) Ask(_ string) (string, error) {
	m.t.FailNow()
	return "", nil
}

func (m mockUi) AskSecret(_ string) (string, error) {
	m.t.FailNow()
	return "", nil
}
func (m mockUi) Output(s string) { *m.outputData = s }
func (m mockUi) Info(s string)   { m.t.Log(s) }
func (m mockUi) Error(s string)  { m.t.Log(s) }
func (m mockUi) Warn(s string)   { m.t.Log(s) }

func TestJsonFormatter(t *testing.T) {
	os.Setenv(EnvVaultFormat, "json")
	var output string
	ui := mockUi{t: t, SampleData: "something", outputData: &output}
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
	var output string
	ui := mockUi{t: t, SampleData: "something", outputData: &output}
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
	var output string
	ui := mockUi{t: t, outputData: &output}

	// Testing secret formatting
	s := api.Secret{Data: map[string]interface{}{"k": "something"}}
	if err := outputWithFormat(ui, &s, &s); err != 0 {
		t.Fatal(err)
	}
	if !strings.Contains(output, "something") {
		t.Fatal("did not find 'something'")
	}
}

// TestStatusFormat tests to verify that the embedded struct
// SealStatusOutput ignores omitEmpty fields and prints out
// fields in the embedded struct explicitly. It also checks the spacing,
// indentation, and delimiters of table formatting explicitly.
func TestStatusFormat(t *testing.T) {
	var output string
	ui := mockUi{t: t, outputData: &output}
	os.Setenv(EnvVaultFormat, "table")

	statusHA := getMockStatusData(false)
	statusOmitEmpty := getMockStatusData(true)

	// Testing that HA fields are formatted properly for table.
	// All fields (including new HA fields) are expected
	if err := outputWithFormat(ui, nil, statusHA); err != 0 {
		t.Fatal(err)
	}

	expectedOutputString := `Key                           Value
---                           -----
Recovery Seal Type            type
Initialized                   true
Sealed                        true
Total Recovery Shares         2
Threshold                     1
Unseal Progress               3/1
Unseal Nonce                  nonce
Seal Migration in Progress    true
Version                       version
Build Date                    build date
Storage Type                  storage type
Cluster Name                  cluster name
Cluster ID                    cluster id
HA Enabled                    true
Raft Committed Index          3
Raft Applied Index            4
Last WAL                      2
Warnings                      [warning]`

	if expectedOutputString != output {
		fmt.Printf("%s\n%+v\n %s\n%+v\n", "output found was: ", output, "versus", expectedOutputString)
		t.Fatal("format output for status does not match expected format. Check print statements above.")
	}

	// Testing that omitEmpty fields are omitted from status
	// no HA fields are expected, except HA Enabled
	if err := outputWithFormat(ui, nil, statusOmitEmpty); err != 0 {
		t.Fatal(err)
	}

	expectedOutputString = `Key                           Value
---                           -----
Recovery Seal Type            type
Initialized                   true
Sealed                        true
Total Recovery Shares         2
Threshold                     1
Unseal Progress               3/1
Unseal Nonce                  nonce
Seal Migration in Progress    true
Version                       version
Build Date                    build date
Storage Type                  n/a
HA Enabled                    false`

	if expectedOutputString != output {
		fmt.Printf("%s\n%+v\n %s\n%+v\n", "output found was: ", output, "versus", expectedOutputString)
		t.Fatal("format output for status does not match expected format. Check print statements above.")
	}
}

// getMockStatusData outputs a SealStatusOutput struct from format.go to be used
// for testing. The emptyfields parameter specifies whether the struct will be
// initialized with all the omitempty fields as empty or not.
func getMockStatusData(emptyFields bool) SealStatusOutput {
	var status SealStatusOutput
	var sealStatusResponseMock api.SealStatusResponse
	if !emptyFields {
		sealStatusResponseMock = api.SealStatusResponse{
			Type:         "type",
			Initialized:  true,
			Sealed:       true,
			T:            1,
			N:            2,
			Progress:     3,
			Nonce:        "nonce",
			Version:      "version",
			BuildDate:    "build date",
			Migration:    true,
			ClusterName:  "cluster name",
			ClusterID:    "cluster id",
			RecoverySeal: true,
			StorageType:  "storage type",
			Warnings:     []string{"warning"},
		}

		// must initialize this struct without explicit field names due to embedding
		status = SealStatusOutput{
			sealStatusResponseMock,
			true,                     // HAEnabled
			true,                     // IsSelf
			time.Time{}.UTC(),        // ActiveTime
			"leader address",         // LeaderAddress
			"leader cluster address", // LeaderClusterAddress
			true,                     // PerfStandby
			1,                        // PerfStandbyLastRemoteWAL
			2,                        // LastWAL
			3,                        // RaftCommittedIndex
			4,                        // RaftAppliedIndex
		}
	} else {
		sealStatusResponseMock = api.SealStatusResponse{
			Type:         "type",
			Initialized:  true,
			Sealed:       true,
			T:            1,
			N:            2,
			Progress:     3,
			Nonce:        "nonce",
			Version:      "version",
			BuildDate:    "build date",
			Migration:    true,
			ClusterName:  "",
			ClusterID:    "",
			RecoverySeal: true,
			StorageType:  "",
		}

		// must initialize this struct without explicit field names due to embedding
		status = SealStatusOutput{
			sealStatusResponseMock,
			false,             // HAEnabled
			false,             // IsSelf
			time.Time{}.UTC(), // ActiveTime
			"",                // LeaderAddress
			"",                // LeaderClusterAddress
			false,             // PerfStandby
			0,                 // PerfStandbyLastRemoteWAL
			0,                 // LastWAL
			0,                 // RaftCommittedIndex
			0,                 // RaftAppliedIndex
		}
	}
	return status
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
			[]string{"token", "renew", "-format", "json"},
			"{",
			0,
		},
		{
			"format_bad",
			[]string{"token", "renew", "-format", "nope-not-real"},
			"Invalid output format",
			1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			client, closer := testVaultServer(t)
			defer closer()

			stdout := bytes.NewBuffer(nil)
			stderr := bytes.NewBuffer(nil)
			runOpts := &RunOptions{
				Stdout: stdout,
				Stderr: stderr,
				Client: client,
			}

			// Login with the token so we can renew-self.
			token, _ := testTokenAndAccessor(t, client)
			client.SetToken(token)

			code := RunCustom(tc.args, runOpts)
			if code != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}

			combined := stdout.String() + stderr.String()
			if !strings.Contains(combined, tc.out) {
				t.Errorf("expected %q to contain %q", combined, tc.out)
			}
		})
	}
}
