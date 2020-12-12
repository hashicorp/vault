package stepwise

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

// testTesting is used for testing the legacy testing framework
var testTesting = false

type testRun struct {
	expectedTestT *mockT
	environment   *mockEnvironment
	steps         []Step
	skipTeardown  bool
	requests      *requestCounts
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
		"basic_list": {
			steps: []Step{
				stepFunc("keys", ListOperation, false),
			},
			environment: new(mockEnvironment),
			requests: &requestCounts{
				listRequests: 1,
			},
		},
		"basic_list_read": {
			steps: []Step{
				stepFunc("keys", ListOperation, false),
				stepFunc("keys/name", ReadOperation, false),
			},
			environment: new(mockEnvironment),
			requests: &requestCounts{
				listRequests:   1,
				readRequests:   1,
				revokeRequests: 1,
			},
		},
		"basic_unauth": {
			steps: []Step{
				stepFuncWithoutAuth("keys", ListOperation, true),
			},
			expectedTestT: &mockT{
				ErrorCalled: true,
			},
			environment: new(mockEnvironment),
		},
		"error": {
			steps: []Step{
				stepFunc("keys", ListOperation, false),
				stepFunc("keys/something", ReadOperation, true),
			},
			expectedTestT: &mockT{
				ErrorCalled: true,
			},
			environment: new(mockEnvironment),
			requests: &requestCounts{
				listRequests: 1,
			},
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
				stepFunc("keys/name", ReadOperation, false),
				stepFunc("keys/name", ReadOperation, false),
				stepFunc("keys/name", DeleteOperation, false),
			},
			skipTeardown: true,
			environment:  new(mockEnvironment),
			requests: &requestCounts{
				listRequests:   1,
				readRequests:   2,
				revokeRequests: 2,
				deleteRequests: 1,
			},
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
			if tr.requests != nil {
				if !reflect.DeepEqual(*tr.requests, tr.environment.requests) {
					t.Fatalf("request counts do not match: %#v / %#v", tr.requests, tr.environment.requests)
				}
			}
		})
	}
}

type requestCounts struct {
	writeRequests  int
	readRequests   int
	deleteRequests int
	revokeRequests int
	listRequests   int
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
		UnAuth            bool
	}
	testRequests := map[string]testRequest{
		"list": {
			Operation:         ListOperation,
			Path:              "keys",
			ExpectedRequestID: "list-request",
		},
		"read": {
			Operation:         ReadOperation,
			Path:              "keys/name",
			ExpectedRequestID: "read-request",
		},
		"write": {
			Operation:         WriteOperation,
			Path:              "keys/name",
			ExpectedRequestID: "write-request",
		},
		"update": {
			Operation:         UpdateOperation,
			Path:              "keys/name",
			ExpectedRequestID: "write-request",
		},
		"update_unauth": {
			Operation: UpdateOperation,
			Path:      "keys/name",
			UnAuth:    true,
			ExpectErr: true,
		},
		"delete": {
			Operation:         DeleteOperation,
			Path:              "keys/name",
			ExpectedRequestID: "delete-request",
		},
		"error": {
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

			if tc.UnAuth {
				step.Unauthenticated = tc.UnAuth
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
	requests       requestCounts
}

// Setup creates the mock environment, establishing a test HTTP server
func (m *mockEnvironment) Setup() error {
	mux := http.NewServeMux()
	// LIST
	mux.HandleFunc("/v1/test/keys", func(w http.ResponseWriter, req *http.Request) {
		checkAuth(w, req)
		switch req.Method {
		case "GET":
			m.requests.listRequests++
			respondCommon("list", true, w, req)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	// lease revoke
	mux.HandleFunc("/v1/sys/leases/revoke", func(w http.ResponseWriter, req *http.Request) {
		checkAuth(w, req)
		m.requests.revokeRequests++
		w.WriteHeader(http.StatusOK)
	})
	// READ, DELETE, WRITE
	mux.HandleFunc("/v1/test/keys/name", func(w http.ResponseWriter, req *http.Request) {
		checkAuth(w, req)
		var method string
		// indicate if the common response should include a lease id
		var excludeLease bool
		switch req.Method {
		case "GET":
			m.requests.readRequests++
			method = "read"
		case "POST":
		case "PUT":
			m.requests.writeRequests++
			method = "write"
		case "DELETE":
			m.requests.deleteRequests++
			excludeLease = true
			method = "delete"
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
		respondCommon(method, excludeLease, w, req)
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

// respondCommon returns a mock secret with a request ID that matches the
// request method that was used to invoke it. A true Vault server would not
// respond with a request id / lease id for DELETE or REVOKE, but we do that
// here to verify that the makeRequest method translates the Step Operation
// and calls delete/revoke correctly
func respondCommon(id string, excludeLease bool, w http.ResponseWriter, req *http.Request) {
	resp := api.Secret{
		RequestID: id + "-request",
		LeaseID:   "lease-id",
	}
	if excludeLease {
		resp.LeaseID = ""
	}
	out, err := jsonutil.EncodeJSON(resp)
	if err != nil {
		panic(err)
	}
	w.Write(out)
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

	// need to set root token here to mimic an actual root token of a cluster
	client.SetToken(m.RootToken())
	m.client = client
	return client, nil
}

func (m *mockEnvironment) Teardown() error {
	m.teardownCalled = true
	m.ts.Close()
	return nil
}

func (m *mockEnvironment) Name() string { return "" }

func (m *mockEnvironment) MountPath() string {
	return "/test/"
}

func (m *mockEnvironment) RootToken() string { return "root-token" }

func stepFuncWithoutAuth(path string, operation Operation, shouldError bool) Step {
	return stepFuncCommon(path, operation, shouldError, true)
}

func stepFunc(path string, operation Operation, shouldError bool) Step {
	return stepFuncCommon(path, operation, shouldError, false)
}

func stepFuncCommon(path string, operation Operation, shouldError bool, unauth bool) Step {
	s := Step{
		Operation:       operation,
		Path:            path,
		Unauthenticated: unauth,
	}
	if shouldError {
		s.Assert = func(resp *api.Secret, err error) error {
			return errors.New("some error")
		}
	}
	return s
}

// mockT implements TestT for testing
type mockT struct {
	ErrorCalled bool
	ErrorArgs   []interface{}
	FatalCalled bool
	FatalArgs   []interface{}
	SkipCalled  bool
	SkipArgs    []interface{}

	f bool
}

func (t *mockT) Error(args ...interface{}) {
	t.ErrorCalled = true
	t.ErrorArgs = args
	t.f = true
}

func (t *mockT) Fatal(args ...interface{}) {
	t.FatalCalled = true
	t.FatalArgs = args
	t.f = true
}

func (t *mockT) Skip(args ...interface{}) {
	t.SkipCalled = true
	t.SkipArgs = args
	t.f = true
}

func (t *mockT) Helper() {}

// validates that X-Vault-Token is set on the requets to the mock endpoints
func checkAuth(w http.ResponseWriter, r *http.Request) {
	if token := r.Header.Get("X-Vault-Token"); token == "" {
		// not authenticated
		w.WriteHeader(http.StatusForbidden)
	}
}
