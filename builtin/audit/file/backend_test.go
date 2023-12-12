// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package file

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestAuditFile_fileModeNew(t *testing.T) {
	modeStr := "0777"
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(t.TempDir(), "auditTest.txt")
	config := map[string]string{
		"path": file,
		"mode": modeStr,
	}

	_, err = Factory(context.Background(), &audit.BackendConfig{
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Config:     config,
	}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(file)
	if err != nil {
		t.Fatalf("Cannot retrieve file mode from `Stat`")
	}
	if info.Mode() != os.FileMode(mode) {
		t.Fatalf("File mode does not match.")
	}
}

func TestAuditFile_fileModeExisting(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("Failure to create test file.")
	}
	defer os.Remove(f.Name())

	err = os.Chmod(f.Name(), 0o777)
	if err != nil {
		t.Fatalf("Failure to chmod temp file for testing.")
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("Failure to close temp file for test.")
	}

	config := map[string]string{
		"path": f.Name(),
	}

	_, err = Factory(context.Background(), &audit.BackendConfig{
		Config:     config,
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
	}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(f.Name())
	if err != nil {
		t.Fatalf("cannot retrieve file mode from `Stat`")
	}
	if info.Mode() != os.FileMode(0o600) {
		t.Fatalf("File mode does not match.")
	}
}

func TestAuditFile_fileMode0000(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("Failure to create test file. The error is %v", err)
	}
	defer os.Remove(f.Name())

	err = os.Chmod(f.Name(), 0o777)
	if err != nil {
		t.Fatalf("Failure to chmod temp file for testing. The error is %v", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("Failure to close temp file for test. The error is %v", err)
	}

	config := map[string]string{
		"path": f.Name(),
		"mode": "0000",
	}

	_, err = Factory(context.Background(), &audit.BackendConfig{
		Config:     config,
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
	}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(f.Name())
	if err != nil {
		t.Fatalf("cannot retrieve file mode from `Stat`. The error is %v", err)
	}
	if info.Mode() != os.FileMode(0o777) {
		t.Fatalf("File mode does not match.")
	}
}

// TestAuditFile_EventLogger_fileModeNew verifies that the Factory function
// correctly sets the file mode when the useEventLogger argument is set to
// true.
func TestAuditFile_EventLogger_fileModeNew(t *testing.T) {
	modeStr := "0777"
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(t.TempDir(), "auditTest.txt")
	config := map[string]string{
		"path": file,
		"mode": modeStr,
	}

	_, err = Factory(context.Background(), &audit.BackendConfig{
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Config:     config,
	}, true, nil)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(file)
	if err != nil {
		t.Fatalf("Cannot retrieve file mode from `Stat`")
	}
	if info.Mode() != os.FileMode(mode) {
		t.Fatalf("File mode does not match.")
	}
}

func BenchmarkAuditFile_request(b *testing.B) {
	config := map[string]string{
		"path": "/dev/null",
	}
	sink, err := Factory(context.Background(), &audit.BackendConfig{
		Config:     config,
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
	}, false, nil)
	if err != nil {
		b.Fatal(err)
	}

	in := &logical.LogInput{
		Auth: &logical.Auth{
			ClientToken:     "foo",
			Accessor:        "bar",
			EntityID:        "foobarentity",
			DisplayName:     "testtoken",
			NoDefaultPolicy: true,
			Policies:        []string{"root"},
			TokenType:       logical.TokenTypeService,
		},
		Request: &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "/foo",
			Connection: &logical.Connection{
				RemoteAddr: "127.0.0.1",
			},
			WrapInfo: &logical.RequestWrapInfo{
				TTL: 60 * time.Second,
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
	}

	ctx := namespace.RootContext(context.Background())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := sink.LogRequest(ctx, in); err != nil {
				panic(err)
			}
		}
	})
}

// TestBackend_formatterConfig ensures that all the configuration values are parsed correctly.
func TestBackend_formatterConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config          map[string]string
		want            audit.FormatterConfig
		wantErr         bool
		expectedMessage string
	}{
		"happy-path-json": {
			config: map[string]string{
				"format":               audit.JSONFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want: audit.FormatterConfig{
				Raw:                true,
				HMACAccessor:       true,
				ElideListResponses: true,
				RequiredFormat:     "json",
			}, wantErr: false,
		},
		"happy-path-jsonx": {
			config: map[string]string{
				"format":               audit.JSONxFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want: audit.FormatterConfig{
				Raw:                true,
				HMACAccessor:       true,
				ElideListResponses: true,
				RequiredFormat:     "jsonx",
			},
			wantErr: false,
		},
		"invalid-format": {
			config: map[string]string{
				"format":               " squiggly ",
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "true",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "audit.NewFormatterConfig: error applying options: audit.(format).validate: 'squiggly' is not a valid format: invalid parameter",
		},
		"invalid-hmac-accessor": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "maybe",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "file.formatterConfig: unable to parse 'hmac_accessor': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"invalid-log-raw": {
			config: map[string]string{
				"format":        audit.JSONFormat.String(),
				"hmac_accessor": "true",
				"log_raw":       "maybe",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "file.formatterConfig: unable to parse 'log_raw': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
		"invalid-elide-bool": {
			config: map[string]string{
				"format":               audit.JSONFormat.String(),
				"hmac_accessor":        "true",
				"log_raw":              "true",
				"elide_list_responses": "maybe",
			},
			want:            audit.FormatterConfig{},
			wantErr:         true,
			expectedMessage: "file.formatterConfig: unable to parse 'elide_list_responses': strconv.ParseBool: parsing \"maybe\": invalid syntax",
		},
	}
	for name, testCase := range tests {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := formatterConfig(testCase.config)
			if testCase.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, testCase.expectedMessage)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, testCase.want, got)
		})
	}
}

// TestBackend_configureFilterNode ensures that configureFilterNode handles various
// filter values as expected. Empty (including whitespace) strings should return
// no error but skip configuration of the node.
func TestBackend_configureFilterNode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		filter           string
		shouldSkipNode   bool
		wantErr          bool
		expectedErrorMsg string
	}{
		"happy": {
			filter: "foo == bar",
		},
		"empty": {
			filter:         "",
			shouldSkipNode: true,
		},
		"spacey": {
			filter:         "    ",
			shouldSkipNode: true,
		},
		"bad": {
			filter:           "___qwerty",
			wantErr:          true,
			expectedErrorMsg: "file.(Backend).configureFilterNode: error creating filter node: audit.NewEntryFilter: cannot create new audit filter",
		},
	}
	for name, testCase := range tests {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := &Backend{
				nodeIDList: []eventlogger.NodeID{},
				nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
			}

			err := b.configureFilterNode(testCase.filter)

			switch {
			case testCase.wantErr:
				require.Error(t, err)
				require.ErrorContains(t, err, testCase.expectedErrorMsg)
				require.Len(t, b.nodeIDList, 0)
				require.Len(t, b.nodeMap, 0)
			case testCase.shouldSkipNode:
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 0)
				require.Len(t, b.nodeMap, 0)
			default:
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 1)
				require.Len(t, b.nodeMap, 1)
				id := b.nodeIDList[0]
				node := b.nodeMap[id]
				require.Equal(t, eventlogger.NodeTypeFilter, node.Type())
			}
		})
	}
}

// TestBackend_configureFormatterNode ensures that configureFormatterNode
// populates the nodeIDList and nodeMap on Backend when given valid formatConfig.
func TestBackend_configureFormatterNode(t *testing.T) {
	t.Parallel()

	b := &Backend{
		nodeIDList: []eventlogger.NodeID{},
		nodeMap:    map[eventlogger.NodeID]eventlogger.Node{},
	}

	formatConfig, err := audit.NewFormatterConfig()
	require.NoError(t, err)

	err = b.configureFormatterNode(formatConfig)

	require.NoError(t, err)
	require.Len(t, b.nodeIDList, 1)
	require.Len(t, b.nodeMap, 1)
	id := b.nodeIDList[0]
	node := b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFormatter, node.Type())
}

/*
SINKS:
	name: (should never be bad really as its mount path, but better safe than sorry)
		empty
		spaces
		happy
	filePath:
		empty
		spaces
		stdout
		discard
		some-legit-value
		some-value-we-dont-have-permission-to
	mode:
		bs mode
		empty mode
		spaces
		legit mode
	format:
		json
		jsonx
		bs value
		nothing
		spaces
*/
