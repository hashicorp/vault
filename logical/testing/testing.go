package testing

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical/inmem"
	"github.com/hashicorp/vault/vault"
)

// TestEnvVar must be set to a non-empty value for acceptance tests to run.
const TestEnvVar = "VAULT_ACC"

// TestCase is a single set of tests to run for a backend. A TestCase
// should generally map 1:1 to each test method for your acceptance
// tests.
type TestCase struct {
	// Precheck, if non-nil, will be called once before the test case
	// runs at all. This can be used for some validation prior to the
	// test running.
	PreCheck func()

	// LogicalBackend is the backend that will be mounted.
	LogicalBackend logical.Backend

	// LogicalFactory can be used instead of LogicalBackend if the
	// backend requires more construction
	LogicalFactory logical.Factory

	// CredentialBackend is the backend that will be mounted.
	CredentialBackend logical.Backend

	// CredentialFactory can be used instead of CredentialBackend if the
	// backend requires more construction
	CredentialFactory logical.Factory

	// Steps are the set of operations that are run for this test case.
	Steps []TestStep

	// Teardown will be called before the test case is over regardless
	// of if the test succeeded or failed. This should return an error
	// in the case that the test can't guarantee all resources were
	// properly cleaned up.
	Teardown TestTeardownFunc

	// AcceptanceTest, if set, the test case will be run only if
	// the environment variable VAULT_ACC is set. If not this test case
	// will be run as a unit test.
	AcceptanceTest bool
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

	// PreFlight is called directly before execution of the request, allowing
	// modification of the request parameters (e.g. Path) with dynamic values.
	PreFlight PreFlightFunc

	// ErrorOk, if true, will let erroneous responses through to the check
	ErrorOk bool

	// Unauthenticated, if true, will make the request unauthenticated.
	Unauthenticated bool

	// RemoteAddr, if set, will set the remote addr on the request.
	RemoteAddr string

	// ConnState, if set, will set the tls connection state
	ConnState *tls.ConnectionState
}

// TestCheckFunc is the callback used for Check in TestStep.
type TestCheckFunc func(*logical.Response) error

// PreFlightFunc is used to modify request parameters directly before execution
// in each TestStep.
type PreFlightFunc func(*logical.Request) error

// TestTeardownFunc is the callback used for Teardown in TestCase.
type TestTeardownFunc func() error

// Test performs an acceptance test on a backend with the given test case.
//
// Tests are not run unless an environmental variable "VAULT_ACC" is
// set to some non-empty value. This is to avoid test cases surprising
// a user by creating real resources.
//
// Tests will fail unless the verbose flag (`go test -v`, or explicitly
// the "-test.v" flag) is set. Because some acceptance tests take quite
// long, we require the verbose flag so users are able to see progress
// output.
func Test(tt TestT, c TestCase) {
	// We only run acceptance tests if an env var is set because they're
	// slow and generally require some outside configuration.
	if c.AcceptanceTest && os.Getenv(TestEnvVar) == "" {
		tt.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			TestEnvVar))
		return
	}

	// We require verbose mode so that the user knows what is going on.
	if c.AcceptanceTest && !testTesting && !testing.Verbose() {
		tt.Fatal("Acceptance tests must be run with the -v flag on tests")
		return
	}

	// Run the PreCheck if we have it
	if c.PreCheck != nil {
		c.PreCheck()
	}

	// Defer on the teardown, regardless of pass/fail at this point
	if c.Teardown != nil {
		defer c.Teardown()
	}

	// Check that something is provided
	if c.LogicalBackend == nil && c.LogicalFactory == nil {
		if c.CredentialBackend == nil && c.CredentialFactory == nil {
			tt.Fatal("Must provide either Backend or Factory")
			return
		}
	}
	// We currently only support doing one logical OR one credential test at a time.
	if (c.LogicalFactory != nil || c.LogicalBackend != nil) && (c.CredentialFactory != nil || c.CredentialBackend != nil) {
		tt.Fatal("Must provide only one backend or factory")
		return
	}

	// Create an in-memory Vault core
	logger := logging.NewVaultLogger(log.Trace)

	phys, err := inmem.NewInmem(nil, logger)
	if err != nil {
		tt.Fatal(err)
		return
	}

	config := &vault.CoreConfig{
		Physical:        phys,
		DisableMlock:    true,
		BuiltinRegistry: vault.NewMockBuiltinRegistry(),
	}

	if c.LogicalBackend != nil || c.LogicalFactory != nil {
		config.LogicalBackends = map[string]logical.Factory{
			"test": func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
				if c.LogicalBackend != nil {
					return c.LogicalBackend, nil
				}
				return c.LogicalFactory(ctx, conf)
			},
		}
	}
	if c.CredentialBackend != nil || c.CredentialFactory != nil {
		config.CredentialBackends = map[string]logical.Factory{
			"test": func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
				if c.CredentialBackend != nil {
					return c.CredentialBackend, nil
				}
				return c.CredentialFactory(ctx, conf)
			},
		}
	}

	core, err := vault.NewCore(config)
	if err != nil {
		tt.Fatal("error initializing core: ", err)
		return
	}

	// Initialize the core
	init, err := core.Initialize(context.Background(), &vault.InitParams{
		BarrierConfig: &vault.SealConfig{
			SecretShares:    1,
			SecretThreshold: 1,
		},
		RecoveryConfig: nil,
	})
	if err != nil {
		tt.Fatal("error initializing core: ", err)
		return
	}

	// Unseal the core
	if unsealed, err := core.Unseal(init.SecretShares[0]); err != nil {
		tt.Fatal("error unsealing core: ", err)
		return
	} else if !unsealed {
		tt.Fatal("vault shouldn't be sealed")
		return
	}

	// Create an HTTP API server and client
	ln, addr := http.TestServer(nil, core)
	defer ln.Close()
	clientConfig := api.DefaultConfig()
	clientConfig.Address = addr
	client, err := api.NewClient(clientConfig)
	if err != nil {
		tt.Fatal("error initializing HTTP client: ", err)
		return
	}

	// Set the token so we're authenticated
	client.SetToken(init.RootToken)

	prefix := "mnt"
	if c.LogicalBackend != nil || c.LogicalFactory != nil {
		// Mount the backend
		mountInfo := &api.MountInput{
			Type:        "test",
			Description: "acceptance test",
		}
		if err := client.Sys().Mount(prefix, mountInfo); err != nil {
			tt.Fatal("error mounting backend: ", err)
			return
		}
	}

	isAuthBackend := false
	if c.CredentialBackend != nil || c.CredentialFactory != nil {
		isAuthBackend = true

		// Enable the test auth method
		opts := &api.EnableAuthOptions{
			Type: "test",
		}
		if err := client.Sys().EnableAuthWithOptions(prefix, opts); err != nil {
			tt.Fatal("error enabling backend: ", err)
			return
		}
	}

	tokenInfo, err := client.Auth().Token().LookupSelf()
	if err != nil {
		tt.Fatal("error looking up token: ", err)
		return
	}
	var tokenPolicies []string
	if tokenPoliciesRaw, ok := tokenInfo.Data["policies"]; ok {
		if tokenPoliciesSliceRaw, ok := tokenPoliciesRaw.([]interface{}); ok {
			for _, p := range tokenPoliciesSliceRaw {
				tokenPolicies = append(tokenPolicies, p.(string))
			}
		}
	}

	// Make requests
	var revoke []*logical.Request
	for i, s := range c.Steps {
		if logger.IsWarn() {
			logger.Warn("Executing test step", "step_number", i+1)
		}

		// Create the request
		req := &logical.Request{
			Operation: s.Operation,
			Path:      s.Path,
			Data:      s.Data,
		}
		if !s.Unauthenticated {
			req.ClientToken = client.Token()
			req.SetTokenEntry(&logical.TokenEntry{
				ID:          req.ClientToken,
				NamespaceID: namespace.RootNamespaceID,
				Policies:    tokenPolicies,
				DisplayName: tokenInfo.Data["display_name"].(string),
			})
		}
		if s.RemoteAddr != "" {
			req.Connection = &logical.Connection{RemoteAddr: s.RemoteAddr}
		}
		if s.ConnState != nil {
			req.Connection = &logical.Connection{ConnState: s.ConnState}
		}

		if s.PreFlight != nil {
			ct := req.ClientToken
			req.ClientToken = ""
			if err := s.PreFlight(req); err != nil {
				tt.Error(fmt.Sprintf("Failed preflight for step %d: %s", i+1, err))
				break
			}
			req.ClientToken = ct
		}

		// Make sure to prefix the path with where we mounted the thing
		req.Path = fmt.Sprintf("%s/%s", prefix, req.Path)

		if isAuthBackend {
			// Prepend the path with "auth"
			req.Path = "auth/" + req.Path
		}

		// Make the request
		resp, err := core.HandleRequest(namespace.RootContext(nil), req)
		if resp != nil && resp.Secret != nil {
			// Revoke this secret later
			revoke = append(revoke, &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "sys/revoke/" + resp.Secret.LeaseID,
			})
		}

		// Test step returned an error.
		if err != nil {
			// But if an error is expected, do not fail the test step,
			// regardless of whether the error is a 'logical.ErrorResponse'
			// or not. Set the err to nil. If the error is a logical.ErrorResponse,
			// it will be handled later.
			if s.ErrorOk {
				err = nil
			} else {
				// If the error is not expected, fail right away.
				tt.Error(fmt.Sprintf("Failed step %d: %s", i+1, err))
				break
			}
		}

		// If the error is a 'logical.ErrorResponse' and if error was not expected,
		// set the error so that this can be caught below.
		if resp.IsError() && !s.ErrorOk {
			err = fmt.Errorf("erroneous response:\n\n%#v", resp)
		}

		// Either the 'err' was nil or if an error was expected, it was set to nil.
		// Call the 'Check' function if there is one.
		//
		// TODO: This works perfectly for now, but it would be better if 'Check'
		// function takes in both the response object and the error, and decide on
		// the action on its own.
		if err == nil && s.Check != nil {
			// Call the test method
			err = s.Check(resp)
		}

		if err != nil {
			tt.Error(fmt.Sprintf("Failed step %d: %s", i+1, err))
			break
		}
	}

	// Revoke any secrets we might have.
	var failedRevokes []*logical.Secret
	for _, req := range revoke {
		if logger.IsWarn() {
			logger.Warn("Revoking secret", "secret", fmt.Sprintf("%#v", req))
		}
		req.ClientToken = client.Token()
		resp, err := core.HandleRequest(namespace.RootContext(nil), req)
		if err == nil && resp.IsError() {
			err = fmt.Errorf("erroneous response:\n\n%#v", resp)
		}
		if err != nil {
			failedRevokes = append(failedRevokes, req.Secret)
			tt.Error(fmt.Sprintf("Revoke error: %s", err))
		}
	}

	// Perform any rollbacks. This should no-op if there aren't any.
	// We set the "immediate" flag here that any backend can pick up on
	// to do all rollbacks immediately even if the WAL entries are new.
	logger.Warn("Requesting RollbackOperation")
	rollbackPath := prefix + "/"
	if c.CredentialFactory != nil || c.CredentialBackend != nil {
		rollbackPath = "auth/" + rollbackPath
	}
	req := logical.RollbackRequest(rollbackPath)
	req.Data["immediate"] = true
	req.ClientToken = client.Token()
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err == nil && resp.IsError() {
		err = fmt.Errorf("erroneous response:\n\n%#v", resp)
	}
	if err != nil {
		if !errwrap.Contains(err, logical.ErrUnsupportedOperation.Error()) {
			tt.Error(fmt.Sprintf("[ERR] Rollback error: %s", err))
		}
	}

	// If we have any failed revokes, log it.
	if len(failedRevokes) > 0 {
		for _, s := range failedRevokes {
			tt.Error(fmt.Sprintf(
				"WARNING: Revoking the following secret failed. It may\n"+
					"still exist. Please verify:\n\n%#v",
				s))
		}
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
		expected := make([]string, len(policies))
		copy(expected, policies)
		sort.Strings(expected)
		ret := make([]string, len(resp.Auth.Policies))
		copy(ret, resp.Auth.Policies)
		sort.Strings(ret)
		if !reflect.DeepEqual(ret, expected) {
			return fmt.Errorf("invalid policies: expected %#v, got %#v", expected, ret)
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
