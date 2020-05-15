// Package stepwise offers types and functions to enable black-box style tests
// that are executed in defined set of steps
package stepwise

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/y0ssar1an/q"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestEnvVar must be set to a non-empty value for acceptance tests to run.
const TestEnvVar = "VAULT_ACC"

// TestTeardownFunc is the callback used for Teardown in TestCase.
type TestTeardownFunc func() error

// Define operations
type StepOperation string

const (
	// The operations below are called per path
	WriteOperation StepOperation = "create"
	ReadOperation                = "read"
	// UpdateOperation                   = "update"
	DeleteOperation = "delete"
	ListOperation   = "list"
	HelpOperation   = "help"
	// AliasLookaheadOperation           = "alias-lookahead"

	// The operations below are called globally, the path is less relevant.
	// RevokeOperation   Operation = "revoke"
	// RenewOperation              = "renew"
	// RollbackOperation           = "rollback"
)

// Step represents a single step of a test Case
type Step struct {
	Operation StepOperation
	// Path is the request path. The mount prefix will be automatically added.
	Path string

	// Arguments to pass in
	Data map[string]interface{}

	// Check is called after this step is executed in order to test that
	// the step executed successfully. If this is not set, then the next
	// step will be called
	Check StepCheckFunc

	// PreFlight is called directly before execution of the request, allowing
	// modification of the request parameters (e.g. Path) with dynamic values.
	// PreFlight PreFlightFunc

	// ErrorOk, if true, will let erroneous responses through to the check
	ErrorOk bool

	// Unauthenticated, if true, will make the request unauthenticated.
	Unauthenticated bool
}

// StepCheckFunc is the callback used for Check in TestStep.
type StepCheckFunc func(*api.Secret) error

// StepDriver is the interface Drivers need to implement to be used in
// Case to execute each Step
type StepDriver interface {
	Setup() error
	Client() (*api.Client, error)
	Teardown() error
	Name() string // maybe?
}

// Case is a single set of tests to run for a backend. A test Case
// should generally map 1:1 to each test method for your integration
// tests.
type Case struct {
	Driver StepDriver

	// Precheck, if non-nil, will be called once before the test case
	// runs at all. This can be used for some validation prior to the
	// test running.
	PreCheck func()

	// Steps are the set of operations that are run for this test case.
	Steps []Step

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
func Run(tt TestT, c Case) {
	var err error
	// q.Q("==> here in testing.Test")
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
		q.Q(">> testing defer teardown")
		defer func() {
			err := c.Teardown()
			if err != nil {
				tt.Error("failed to tear down:", err)
			}
		}()
	}

	// TODO setup on driver here
	if c.Driver != nil {
		q.Q("Found driver:", c.Driver)
		err := c.Driver.Setup()
		if err != nil {
			c.Driver.Teardown()
			tt.Fatal(err)
		}
	} else {
		q.Q("nil driver")
		tt.Fatal("nil driver")
	}

	// Create an in-memory Vault core
	logger := logging.NewVaultLogger(log.Trace)
	q.Q("==> Run here")

	prefix := "mnt"
	isAuthBackend := false
	// if c.CredentialBackend != nil || c.CredentialFactory != nil {
	// 	isAuthBackend = true

	// 	// Enable the test auth method
	// 	opts := &api.EnableAuthOptions{
	// 		Type: "test",
	// 	}
	// 	if err := client.Sys().EnableAuthWithOptions(prefix, opts); err != nil {
	// 		tt.Fatal("error enabling backend: ", err)
	// 		return
	// 	}
	// }

	// tokenInfo, err := client.Auth().Token().LookupSelf()
	// if err != nil {
	// 	tt.Fatal("error looking up token: ", err)
	// 	return
	// }
	// var tokenPolicies []string
	// if tokenPoliciesRaw, ok := tokenInfo.Data["policies"]; ok {
	// 	if tokenPoliciesSliceRaw, ok := tokenPoliciesRaw.([]interface{}); ok {
	// 		for _, p := range tokenPoliciesSliceRaw {
	// 			tokenPolicies = append(tokenPolicies, p.(string))
	// 		}
	// 	}
	// }

	// TODO
	//	Steps here
	// TODO - go through steps after setting up docker
	//

	// Make requests
	var revoke []*logical.Request
	for i, s := range c.Steps {
		q.Q("==> step:", s)
		if logger.IsWarn() {
			logger.Warn("Executing test step", "step_number", i+1)
		}

		// // Create the request
		// req := &logical.Request{
		// 	// Operation: s.Operation,
		// 	// Path:      s.Path,
		// 	// Data:      s.Data,
		// }

		// TODO hard coded path here, need mount point. Will it be dynamic? probabaly
		// needs to be
		path := fmt.Sprintf("transit/%s", s.Path)
		var err error
		var resp *api.Secret
		client, cerr := c.Driver.Client()
		if cerr != nil {
			tt.Fatal(cerr)
		}
		// TODO should check expect none here?
		// var lr *logical.Response
		switch s.Operation {
		case WriteOperation, UpdateOperation:
			resp, err = client.Logical().Write(path, s.Data)
		case ReadOperation:
			// resp, err = client.Logical().ReadWithData(path, s.Data)
			resp, err = client.Logical().Read(path)
		case ListOperation:
			resp, err = client.Logical().List(path)
			// TODO why though
			// lr = &logical.Response{}
		case DeleteOperation:
			resp, err = client.Logical().Delete(path)
		default:
			panic("bad operation")
		}
		q.Q("test resp:", resp)
		// if !s.Unauthenticated {
		// 	// req.ClientToken = client.Token()
		// 	// req.SetTokenEntry(&logical.TokenEntry{
		// 	// 	ID:          req.ClientToken,
		// 	// 	NamespaceID: namespace.RootNamespaceID,
		// 	// 	Policies:    tokenPolicies,
		// 	// 	DisplayName: tokenInfo.Data["display_name"].(string),
		// 	// })
		// }
		// req.Connection = &logical.Connection{RemoteAddr: s.RemoteAddr}
		// if s.ConnState != nil {
		// 	req.Connection.ConnState = s.ConnState
		// }

		// if s.PreFlight != nil {
		// 	// ct := req.ClientToken
		// 	// req.ClientToken = ""
		// 	// if err := s.PreFlight(req); err != nil {
		// 	// 	tt.Error(fmt.Sprintf("Failed preflight for step %d: %s", i+1, err))
		// 	// 	break
		// 	// }
		// 	// req.ClientToken = ct
		// }

		// Make sure to prefix the path with where we mounted the thing
		// req.Path = fmt.Sprintf("%s/%s", prefix, req.Path)

		if isAuthBackend {
			// Prepend the path with "auth"
			// req.Path = "auth/" + req.Path
		}

		// Make the request
		// resp, err := core.HandleRequest(namespace.RootContext(nil), req)
		// if resp != nil && resp.Secret != nil {
		// 	// Revoke this secret later
		// 	revoke = append(revoke, &logical.Request{
		// 		Operation: logical.UpdateOperation,
		// 		Path:      "sys/revoke/" + resp.Secret.LeaseID,
		// 	})
		// }

		// TODO
		// - test returned error check here
		//

		// Test step returned an error.
		if err != nil {
			// But if an error is expected, do not fail the test step,
			// regardless of whether the error is a 'logical.ErrorResponse'
			// or not. Set the err to nil. If the error is a logical.ErrorResponse,
			// it will be handled later.
			if s.ErrorOk {
				// err = nil
			} else {
				// // If the error is not expected, fail right away.
				// tt.Error(fmt.Sprintf("Failed step %d: %s", i+1, err))
				// break
			}
		}

		// If the error is a 'logical.ErrorResponse' and if error was not expected,
		// set the error so that this can be caught below.
		// TODO resp error check
		// if resp.IsError() && !s.ErrorOk {
		// 	err = fmt.Errorf("erroneous response:\n\n%#v", resp)
		// }

		// TODO
		// - test check func here
		//

		// Either the 'err' was nil or if an error was expected, it was set to nil.
		// Call the 'Check' function if there is one.
		//
		// TODO: This works perfectly for now, but it would be better if 'Check'
		// function takes in both the response object and the error, and decide on
		// the action on its own.
		// if err == nil && s.Check != nil {
		// 	// Call the test method
		// 	// TODO check here
		// 	// err = s.Check(resp)
		// }

		if err != nil {
			tt.Error(fmt.Sprintf("Failed step %d: %s", i+1, err))
			break
		}
	}

	// TODO
	// TODO
	// - Revoking things here
	//

	// Revoke any secrets we might have.
	var failedRevokes []*logical.Secret
	for _, req := range revoke {
		q.Q("==>==> revoke req:", req)
		// TODO do we revoke?
		// if logger.IsWarn() {
		// 	logger.Warn("Revoking secret", "secret", fmt.Sprintf("%#v", req))
		// }
		// req.ClientToken = client.Token()
		// resp, err := core.HandleRequest(namespace.RootContext(nil), req)
		// if err == nil && resp.IsError() {
		// 	err = fmt.Errorf("erroneous response:\n\n%#v", resp)
		// }
		// if err != nil {
		// 	failedRevokes = append(failedRevokes, req.Secret)
		// 	tt.Error(fmt.Sprintf("Revoke error: %s", err))
		// }
	}

	// TODO
	// - Rollbacks here
	//

	// Perform any rollbacks. This should no-op if there aren't any.
	// We set the "immediate" flag here that any backend can pick up on
	// to do all rollbacks immediately even if the WAL entries are new.
	// logger.Warn("Requesting RollbackOperation")
	// rollbackPath := prefix + "/"
	// if c.CredentialFactory != nil || c.CredentialBackend != nil {
	// 	rollbackPath = "auth/" + rollbackPath
	// }
	// req := logical.RollbackRequest(rollbackPath)
	// req.Data["immediate"] = true
	// req.ClientToken = client.Token()
	// resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	// if err == nil && resp.IsError() {
	// 	err = fmt.Errorf("erroneous response:\n\n%#v", resp)
	// }
	// if err != nil {
	// 	if !errwrap.Contains(err, logical.ErrUnsupportedOperation.Error()) {
	// 		tt.Error(fmt.Sprintf("[ERR] Rollback error: %s", err))
	// 	}
	// }

	// If we have any failed revokes, log it.
	if len(failedRevokes) > 0 {
		for _, s := range failedRevokes {
			tt.Error(fmt.Sprintf(
				"WARNING: Revoking the following secret failed. It may\n"+
					"still exist. Please verify:\n\n%#v",
				s))
		}
	}

	q.Q("==> calling driver teardown()")
	err = c.Driver.Teardown()
	if err != nil {
		tt.Fatal(err)
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
