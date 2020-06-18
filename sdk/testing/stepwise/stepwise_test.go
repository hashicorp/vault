package stepwise

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

type testRun struct {
	expectedTestT *mockT
	environment   *mockEnvironment
	steps         []Step
	skipTeardown  bool
}

// TestStepwise_Run_SkipIfNotAcc tests if the Stepwise Run function skips tests
// if the VAULT_ACC environment variable is not set. This test is seperate from
// the table tests due to the unsetting/re-setting of the environment variable,
// which is assumed/needed for all other tests.
func TestStepwise_Run_SkipIfNotAcc(t *testing.T) {
	if err := os.Setenv(TestEnvVar, ""); err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Setenv(TestEnvVar, "1")
	skipCase := Case{
		Environment: new(mockEnvironment),
		Steps:       []Step{Step{}},
	}

	expected := mockT{
		SkipCalled: true,
	}

	testT := new(mockT)
	Run(testT, skipCase)

	if testT.SkipCalled != expected.SkipCalled {
		t.Fatalf("expected SkipCalled (%t), got (%t)", expected.SkipCalled, testT.SkipCalled)
	}
}

func TestStepwise_Run_Basic(t *testing.T) {
	testRuns := map[string]testRun{
		"basic": {
			steps: []Step{
				stepFunc("keys", ListOperation, false),
			},
			environment: new(mockEnvironment),
		},
		"error": {
			expectedTestT: &mockT{
				ErrorCalled: true,
			},
			steps: []Step{
				stepFunc("keys", ListOperation, false),
				stepFunc("keys", ReadOperation, true),
			},
			environment: new(mockEnvironment),
		},
		"nil-env": {
			expectedTestT: &mockT{
				FatalCalled: true,
			},
			steps: []Step{
				stepFunc("keys", ListOperation, false),
			},
		},
		"skipTeardown": {
			steps: []Step{
				stepFunc("keys", ListOperation, false),
				stepFunc("keys", ReadOperation, false),
			},
			skipTeardown: true,
			environment:  new(mockEnvironment),
		},
	}

	for name, tr := range testRuns {
		t.Run(name, func(t *testing.T) {
			testT := new(mockT)
			expectedT := tr.expectedTestT
			if expectedT == nil {
				expectedT = new(mockT)
			}
			testCase := Case{
				Steps:        tr.steps,
				SkipTeardown: tr.skipTeardown,
			}

			if tr.environment != nil {
				testCase.Environment = tr.environment
			}

			Run(testT, testCase)

			if tr.environment == nil && !testT.FatalCalled {
				t.Fatal("expected FatalCalled with nil environment, but wasn't")
			}

			if tr.environment != nil {
				if tr.skipTeardown && tr.environment.teardownCalled {
					t.Fatal("SkipTeardown is true, but Teardown was called")
				}
				if !tr.skipTeardown && !tr.environment.teardownCalled {
					t.Fatal("SkipTeardown is false, but Teardown was not called")
				}
			}

			if expectedT.ErrorCalled != testT.ErrorCalled {
				t.Fatalf("expected ErrorCalled (%t), got (%t)", expectedT.ErrorCalled, testT.ErrorCalled)
			}
		})
	}
}

func TestStepwise_makeRequest(t *testing.T) {
	me := new(mockEnvironment)
	me.Setup()
	testT := new(mockT)

	type testRequest struct {
		Operation         Operation
		Path              string
		ExpectedRequestID string
		ExpectErr         bool
	}
	testRequests := map[string]testRequest{
		"list": testRequest{
			Operation:         ListOperation,
			Path:              "keys",
			ExpectedRequestID: "list-request",
		},
		"read": testRequest{
			Operation:         ReadOperation,
			Path:              "keys/name",
			ExpectedRequestID: "read-request",
		},
		"error": testRequest{
			Operation: ReadOperation,
			Path:      "error",
			ExpectErr: true,
		},
	}

	for name, tc := range testRequests {
		t.Run(name, func(t *testing.T) {
			step := Step{
				Operation: tc.Operation,
				Path:      tc.Path,
			}

			secret, err := makeRequest(testT, me, step)
			if err != nil && !tc.ExpectErr {
				t.Fatalf("unexpected error: %s", err)
			}
			if err == nil && tc.ExpectErr {
				t.Fatal("expected error but got none:")
			}

			if err != nil && tc.ExpectErr {
				return
			}
			if secret.RequestID != tc.ExpectedRequestID {
				t.Fatalf("expected (%s), got (%s)", tc.ExpectedRequestID, secret.RequestID)
			}
		})
	}
}

type mockEnvironment struct {
	ts     *httptest.Server
	client *api.Client
	l      sync.Mutex

	teardownCalled bool
}

// Setup creates the mock environment, establishing a test HTTP server
func (m *mockEnvironment) Setup() error {
	mux := http.NewServeMux()
	// LIST
	mux.HandleFunc("/v1/test/keys", func(w http.ResponseWriter, req *http.Request) {
		listResp := api.Secret{
			RequestID: "list-request",
			LeaseID:   "lease-id",
		}
		out, err := jsonutil.EncodeJSON(listResp)
		if err != nil {
			panic(err)
		}
		w.Write(out)
	})
	// lease revoke
	mux.HandleFunc("/v1/sys/leases/revoke", func(w http.ResponseWriter, req *http.Request) {
		w.Write(nil)
	})
	// READ
	mux.HandleFunc("/v1/test/keys/name", func(w http.ResponseWriter, req *http.Request) {
		readResp := api.Secret{
			RequestID: "read-request",
			LeaseID:   "lease-id",
		}
		out, err := jsonutil.EncodeJSON(readResp)
		if err != nil {
			panic(err)
		}
		w.Write(out)
	})
	// fall through that rejects any other url than "/"
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "{}")
	})
	m.ts = httptest.NewServer(mux)

	return nil
}

// Client creates a Vault API client configured to the mock environment's test
// server
func (m *mockEnvironment) Client() (*api.Client, error) {
	m.l.Lock()
	defer m.l.Unlock()
	if m.ts == nil {
		return nil, errors.New("client() called but test server is nil")
	}

	if m.client != nil {
		return m.client, nil
	}

	cfg := api.Config{
		Address:    m.ts.URL,
		HttpClient: cleanhttp.DefaultPooledClient(),
		Timeout:    time.Second * 5,
		MaxRetries: 2,
	}

	client, err := api.NewClient(&cfg)
	if err != nil {
		return nil, err
	}

	m.client = client
	return client, nil
}

func (m *mockEnvironment) Teardown() error {
	m.teardownCalled = true
	m.ts.Close()
	return nil
}

func (m *mockEnvironment) Name() string { return "" }
func (m *mockEnvironment) ExpandPath(path string) string {
	return "/test/" + path
}

func (m *mockEnvironment) MountPath() string { return "" }
func (m *mockEnvironment) RootToken() string { return "" }

func stepFunc(path string, operation Operation, shouldError bool) Step {
	s := Step{
		Operation: operation,
		Path:      path,
	}
	if shouldError {
		s.Assert = func(resp *api.Secret, err error) error {
			return errors.New("some error")
		}
	}
	return s
}
