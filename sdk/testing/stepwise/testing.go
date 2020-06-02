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
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestEnvVar must be set to a non-empty value for acceptance tests to run.
const TestEnvVar = "VAULT_ACC"

// TestTeardownFunc is the callback used for Teardown in TestCase.
type TestTeardownFunc func() error

// StepOperation defines operations each step could preform
type StepOperation string

const (
	// WriteOperation and UpdateOperation should be the same
	WriteOperation  StepOperation = "create"
	UpdateOperation               = "update"
	ReadOperation                 = "read"
	DeleteOperation               = "delete"
	ListOperation                 = "list"
	HelpOperation                 = "help"
	// AliasLookaheadOperation           = "alias-lookahead"

	// The operations below are called globally, the path is less relevant.
	RevokeOperation   StepOperation = "revoke"
	RenewOperation                  = "renew"
	RollbackOperation               = "rollback"
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
type StepCheckFunc func(*api.Secret, error) error

// StepDriver is the interface Drivers need to implement to be used in
// Case to execute each Step
type StepDriver interface {
	Setup() error
	Client() (*api.Client, error)
	Teardown() error
	Name() string // maybe?

	// ExpandPath adds any Namespace or mount path to the user defined path
	ExpandPath(string) string

	// MountPath returns the path the plugin is mounted at
	MountPath() string

	// BarrierKeys returns the keys used to seal/unseal the cluster. Used for
	// debugging. TODO verify we should provide this
	//BarrierKeys()        [][]byte
}

// PluginType defines the types of plugins supported
// This type re-create constants as a convienence so users don't need to import/use
// the consts package. Example:
//	driverOptions := &stepwise.DriverOptions{
//		PluginType: stepwise.PluginTypeCredential,
//	}
// versus:
//	driverOptions := &stepwise.DriverOptions{
//		PluginType: consts.PluginTypeCredential,
//	}
// These are originally defined in sdk/helper/consts/plugin_types.go
type PluginType consts.PluginType

const (
	PluginTypeUnknown PluginType = iota
	PluginTypeCredential
	PluginTypeDatabase
	PluginTypeSecrets
)

func (p PluginType) String() string {
	switch p {
	case PluginTypeUnknown:
		return "unknown"
	case PluginTypeCredential:
		return "auth"
	case PluginTypeDatabase:
		return "database"
	case PluginTypeSecrets:
		return "secret"
	default:
		return "unsupported"
	}
}

// DriverOptions are a collection of options each step driver should
// support
type DriverOptions struct {
	// MountPath is an optional string to specify the mount path for the plugin.
	// Defaults to a random string
	// TODO make mount path default to random
	MountPath string

	// Name is used to register the plugin. This can be arbitray but should be a
	// reasonable value. For an example, if the plugin in test is a secret backend
	// that generates UUIDs with the name "vault-plugin-secrets-uuid", then "uuid"
	// or "test-uuid" would be reasonable. The name is used for lookups in the
	// catalog. See "name" in the "Register Plugin" endpoint docs:
	// - https://www.vaultproject.io/api-docs/system/plugins-catalog#register-plugin
	Name string

	// PluginType is the optional type of plugin. See PluginType const defined
	// above
	PluginType PluginType

	// PluginName represents the name of the plugin that gets compiled. In the
	// standard plugin project file layout, it represents the folder under the
	// cmd/ folder. In the below example UUID project, the PluginName would be
	// "uuid":
	//
	// vault-plugin-secrets-uuid/
	// - backend.go
	// - cmd/
	// ----uuid/
	// ------main.go
	// - path_generate.go
	//
	PluginName string
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

	// SkipTeardown allows the TestTeardownFunc to be skipped, leaving any
	// infrastructure created during Driver setup to remain. Depending on the
	// Driver used this could incur costs the user is responsible for.
	// TODO maybe better wording here
	SkipTeardown bool
}

// Run performs an acceptance test on a backend with the given test case.
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
	tt.Helper()
	q.Q("---------")
	q.Q("Stepwise starting...")
	q.Q("---------")
	defer func() {
		q.Q("---------")
		q.Q("end")
		q.Q("---------")
		q.Q("")
	}()

	// We only run acceptance tests if an env var is set because they're
	// slow and generally require some outside configuration.
	if os.Getenv(TestEnvVar) == "" {
		tt.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			TestEnvVar))
		return
	}

	// We require verbose mode so that the user knows what is going on.
	if !testing.Verbose() {
		tt.Fatal("Acceptance tests must be run with the -v flag on tests")
		return
	}

	// Run the PreCheck if we have it
	if c.PreCheck != nil {
		q.Q("--> running precheck")
		c.PreCheck()
	}

	// Defer on the teardown, regardless of pass/fail at this point
	// TODO - checkErr is seperate right now b/c I wanted to stop Vault from being
	// torndown if a test failed. This needs to be removed, or configurable
	var checkErr error
	if c.Teardown != nil {
		defer func(testError error) {
			q.Q("## teardown error check err:", testError)
			if testError != nil {
				q.Q("## test check err is not nil, skipping tearing down")
				return
			}
			q.Q("## test check err is nil, tearing down...")
			err := c.Teardown()
			if err != nil {
				tt.Error("failed to tear down:", err)
			}
		}(checkErr)
	}

	// Create an in-memory Vault core
	// TODO use test logger if available?
	logger := logging.NewVaultLogger(log.Trace)
	if c.Driver == nil {
		tt.Fatal("nil driver in acceptance test")
	}

	err := c.Driver.Setup()
	if err != nil {
		driverErr := fmt.Errorf("error setting up driver: %w", err)
		if err := c.Driver.Teardown(); err != nil {
			driverErr = fmt.Errorf("error during driver teardown: %w", driverErr)
		}
		tt.Fatal(driverErr)
	}

	defer func() {
		if c.SkipTeardown {
			logger.Error("driver Teardown skipped")
			return
		}
		if err := c.Driver.Teardown(); err != nil {
			logger.Error("error in driver teardown:", "error", err)
		}
	}()

	// retrieve the client from the Driver. If this returns an error, fail
	// immediately
	client, err := c.Driver.Client()
	if err != nil {
		tt.Fatal(err)
	}

	// track all responses to revoke any secrets
	var responses []*api.Secret
	q.Q("mount path:", c.Driver.MountPath())
	for i, step := range c.Steps {
		// range is zero based, so add 1 for a human friendly output of steps.
		// "index" here is only used for logging / output, and not to reference the
		// step in the slice of steps.
		index := i + 1
		if logger.IsWarn() {
			logger.Warn("Executing test step", "step_number", index)
		}

		// ExpandPath will turn a test path into a full path, prefixing with the
		// correct mount, namespaces, or "auth" if needed based on mount path and
		// plugin type
		path := c.Driver.ExpandPath(step.Path)
		var err error
		var resp *api.Secret
		// TODO should check expect none here?
		// var lr *logical.Response
		switch step.Operation {
		case WriteOperation, UpdateOperation:
			q.Q("===> Write/Update operation")
			resp, err = client.Logical().Write(path, step.Data)
		case ReadOperation:
			q.Q("===> Read operation")
			// resp, err = client.Logical().ReadWithData(path, step.Data)
			resp, err = client.Logical().Read(path)
		case ListOperation:
			q.Q("===> List operation")
			resp, err = client.Logical().List(path)
		case DeleteOperation:
			q.Q("===> Delete operation")
			resp, err = client.Logical().Delete(path)
		default:
			panic("bad operation")
		}
		if resp != nil {
			responses = append(responses, resp)
		}
		// q.Q("test resp,err:", resp, err)
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

		// TODO
		// - test returned error check here
		//

		// Test step returned an error.
		// if err != nil {
		// 	// But if an error is expected, do not fail the test step,
		// 	// regardless of whether the error is a 'logical.ErrorResponse'
		// 	// or not. Set the err to nil. If the error is a logical.ErrorResponse,
		// 	// it will be handled later.
		// 	if s.ErrorOk {
		// 		q.Q("===> error ok, setting to nil")
		// 		err = nil
		// 	} else {
		// 		// // If the error is not expected, fail right away.
		// 		tt.Error(fmt.Sprintf("Failed step %d: %s", i+1, err))
		// 		break
		// 	}
		// }

		// Either the 'err' was nil or if an error was expected, it was set to nil.
		// Call the 'Check' function if there is one.
		if step.Check != nil {
			checkErr = step.Check(resp, err)
			// TODO allow error
			if checkErr != nil {
				// tt.Fatal("test check error:", checkErr)
				tt.Error(fmt.Sprintf("Failed step %d: %s", index, checkErr))
			}
		}

		// TODO which error is this?
		if err != nil {
			tt.Error(fmt.Sprintf("Failed step %d: %s", index, err))
			break
		}
	}

	// TODO
	// - Revoking things here
	//
	for _, secret := range responses {
		if secret.LeaseID != "" {
			if err := client.Sys().Revoke(secret.LeaseID); err != nil {
				tt.Error(fmt.Sprintf("===>> error revoking lease: %s", err))
			}
		}
	}

	// Revoke any secrets we might have.
	var failedRevokes []*logical.Secret
	for _, req := range responses {
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

	// failsafe - revoke by mount path
	q.Q("==<> failsafe")
	if err := client.Sys().RevokePrefix(c.Driver.MountPath()); err != nil {
		q.Q("==<> error in failsafe:", err)
		revokeErr := fmt.Errorf("[WARN] error revoking by prefix at tend of test: %w", err)
		tt.Error(revokeErr)
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
	Helper()
}
