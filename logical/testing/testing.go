package testing

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
)

// TestEnvVar must be set to a non-empty value for acceptance tests to run.
const TestEnvVar = "TF_ACC"

// TestCase is a single set of tests to run for a backend. A TestCase
// should generally map 1:1 to each test method for your acceptance
// tests.
type TestCase struct {
	// Precheck, if non-nil, will be called once before the test case
	// runs at all. This can be used for some validation prior to the
	// test running.
	PreCheck func()

	// Backend is the backend that will be mounted.
	Backend logical.Backend

	// Factory can be used instead of Backend if the
	// backend requires more construction
	Factory logical.Factory

	// Steps are the set of operations that are run for this test case.
	Steps []TestStep

	// Teardown will be called before the test case is over regardless
	// of if the test succeeded or failed. This should return an error
	// in the case that the test can't guarantee all resources were
	// properly cleaned up.
	Teardown TestTeardownFunc
}

// TestStep is a single step within a TestCase.
type TestStep struct {
	// Operation is the operation to execute
	Operation logical.Operation

	// Path is the request path. The mount prefix will be automatically added.
	Path string

	// Arguments to pass in
	Data map[string]interface{}

	// Check is called after this step is executed in order to test that
	// the step executed successfully. If this is not set, then the next
	// step will be called
	Check TestCheckFunc

	// ErrorOk, if true, will let erroneous responses through to the check
	ErrorOk bool

	// Unauthenticated, if true, will make the request unauthenticated.
	Unauthenticated bool

	// RemoteAddr, if set, will set the remote addr on the request.
	RemoteAddr string

	// ConnState, if set, will set the tls conneciton state
	ConnState *tls.ConnectionState
}

// TestCheckFunc is the callback used for Check in TestStep.
type TestCheckFunc func(*logical.Response) error

// TestTeardownFunc is the callback used for Teardown in TestCase.
type TestTeardownFunc func() error

// Test performs an acceptance test on a backend with the given test case.
//
// Tests are not run unless an environmental variable "TF_ACC" is
// set to some non-empty value. This is to avoid test cases surprising
// a user by creating real resources.
//
// Tests will fail unless the verbose flag (`go test -v`, or explicitly
// the "-test.v" flag) is set. Because some acceptance tests take quite
// long, we require the verbose flag so users are able to see progress
// output.
func Test(t TestT, c TestCase) {
	// We only run acceptance tests if an env var is set because they're
	// slow and generally require some outside configuration.
	if os.Getenv(TestEnvVar) == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			TestEnvVar))
		return
	}

	// We require verbose mode so that the user knows what is going on.
	if !testTesting && !testing.Verbose() {
		t.Fatal("Acceptance tests must be run with the -v flag on tests")
		return
	}

	// Run the PreCheck if we have it
	if c.PreCheck != nil {
		c.PreCheck()
	}

	// Check that something is provided
	if c.Backend == nil && c.Factory == nil {
		t.Fatal("Must provide either Backend or Factory")
	}

	// Create an in-memory Vault core
	core, err := vault.NewCore(&vault.CoreConfig{
		Physical: physical.NewInmem(),
		LogicalBackends: map[string]logical.Factory{
			"test": func(conf *logical.BackendConfig) (logical.Backend, error) {
				if c.Backend != nil {
					return c.Backend, nil
				}
				return c.Factory(conf)
			},
		},
		DisableMlock: true,
	})
	if err != nil {
		t.Fatal("error initializing core: ", err)
		return
	}

	// Initialize the core
	init, err := core.Initialize(&vault.SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		t.Fatal("error initializing core: ", err)
	}

	// Unseal the core
	if unsealed, err := core.Unseal(init.SecretShares[0]); err != nil {
		t.Fatal("error unsealing core: ", err)
		return
	} else if !unsealed {
		t.Fatal("vault shouldn't be sealed")
		return
	}

	// Create an HTTP API server and client
	ln, addr := http.TestServer(nil, core)
	defer ln.Close()
	clientConfig := api.DefaultConfig()
	clientConfig.Address = addr
	client, err := api.NewClient(clientConfig)
	if err != nil {
		t.Fatal("error initializing HTTP client: ", err)
		return
	}

	// Set the token so we're authenticated
	client.SetToken(init.RootToken)

	// Mount the backend
	prefix := "mnt"
	mountInfo := &api.Mount{
		Type:        "test",
		Description: "acceptance test",
	}
	if err := client.Sys().Mount(prefix, mountInfo); err != nil {
		t.Fatal("error mounting backend: ", err)
		return
	}

	// Make requests
	var revoke []*logical.Request
	for i, s := range c.Steps {
		log.Printf("[WARN] Executing test step %d", i+1)

		// Make sure to prefix the path with where we mounted the thing
		path := fmt.Sprintf("%s/%s", prefix, s.Path)

		// Create the request
		req := &logical.Request{
			Operation: s.Operation,
			Path:      path,
			Data:      s.Data,
		}
		if !s.Unauthenticated {
			req.ClientToken = client.Token()
		}
		if s.RemoteAddr != "" {
			req.Connection = &logical.Connection{RemoteAddr: s.RemoteAddr}
		}
		if s.ConnState != nil {
			req.Connection = &logical.Connection{ConnState: s.ConnState}
		}

		// Make the request
		resp, err := core.HandleRequest(req)
		if resp != nil && resp.Secret != nil {
			// Revoke this secret later
			revoke = append(revoke, &logical.Request{
				Operation: logical.WriteOperation,
				Path:      "sys/revoke/" + resp.Secret.LeaseID,
			})
		}
		if err == nil && resp.IsError() && !s.ErrorOk {
			err = fmt.Errorf("Erroneous response:\n\n%#v", resp)
		}
		if err == nil && s.Check != nil {
			// Call the test method
			err = s.Check(resp)
		}
		if err != nil {
			t.Error(fmt.Sprintf("Failed step %d: %s", i+1, err))
			break
		}
	}

	// Revoke any secrets we might have.
	var failedRevokes []*logical.Secret
	for _, req := range revoke {
		log.Printf("[WARN] Revoking secret: %#v", req)
		req.ClientToken = client.Token()
		resp, err := core.HandleRequest(req)
		if err == nil && resp.IsError() {
			err = fmt.Errorf("Erroneous response:\n\n%#v", resp)
		}
		if err != nil {
			failedRevokes = append(failedRevokes, req.Secret)
			t.Error(fmt.Sprintf("[ERR] Revoke error: %s", err))
		}
	}

	// Perform any rollbacks. This should no-op if there aren't any.
	// We set the "immediate" flag here that any backend can pick up on
	// to do all rollbacks immediately even if the WAL entries are new.
	log.Printf("[WARN] Requesting RollbackOperation")
	req := logical.RollbackRequest(prefix + "/")
	req.Data["immediate"] = true
	req.ClientToken = client.Token()
	resp, err := core.HandleRequest(req)
	if err == nil && resp.IsError() {
		err = fmt.Errorf("Erroneous response:\n\n%#v", resp)
	}
	if err != nil && err != logical.ErrUnsupportedOperation {
		t.Error(fmt.Sprintf("[ERR] Rollback error: %s", err))
	}

	// If we have any failed revokes, log it.
	if len(failedRevokes) > 0 {
		for _, s := range failedRevokes {
			t.Error(fmt.Sprintf(
				"WARNING: Revoking the following secret failed. It may\n"+
					"still exist. Please verify:\n\n%#v",
				s))
		}
	}

	// Cleanup
	if c.Teardown != nil {
		c.Teardown()
	}
}

// TestCheckMulti is a helper to have multiple checks.
func TestCheckMulti(fs ...TestCheckFunc) TestCheckFunc {
	return func(resp *logical.Response) error {
		for _, f := range fs {
			if err := f(resp); err != nil {
				return err
			}
		}

		return nil
	}
}

// TestCheckAuth is a helper to check that a request generated an
// auth token with the proper policies.
func TestCheckAuth(policies []string) TestCheckFunc {
	return func(resp *logical.Response) error {
		if resp == nil || resp.Auth == nil {
			return fmt.Errorf("no auth in response")
		}
		if !reflect.DeepEqual(resp.Auth.Policies, policies) {
			return fmt.Errorf("invalid policies: %#v", resp.Auth.Policies)
		}

		return nil
	}
}

// TestCheckAuthDisplayName is a helper to check that a request generated a
// valid display name.
func TestCheckAuthDisplayName(n string) TestCheckFunc {
	return func(resp *logical.Response) error {
		if resp.Auth == nil {
			return fmt.Errorf("no auth in response")
		}
		if n != "" && resp.Auth.DisplayName != "mnt-"+n {
			return fmt.Errorf("invalid display name: %#v", resp.Auth.DisplayName)
		}

		return nil
	}
}

// TestCheckError is a helper to check that a response is an error.
func TestCheckError() TestCheckFunc {
	return func(resp *logical.Response) error {
		if !resp.IsError() {
			return fmt.Errorf("response should be error")
		}

		return nil
	}
}

// TestT is the interface used to handle the test lifecycle of a test.
//
// Users should just use a *testing.T object, which implements this.
type TestT interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Skip(args ...interface{})
}

var testTesting = false
