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

	// ExpectError indicates if this step is expected to return an error. If the
	// step operation returns an error and ExpectError is not set, the test will
	// fail immediately and the StepCheckFunc will not be ran. If ExpectError is
	// true, the StepCheckFunc will be called (if any) and include the error from
	// the response. It is the responsibility of the StepCheckFunc to validate the
	// error is appropriate or not, if expected.
	ExpectError bool

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

	// RootToken returns the root token of the cluster, used for administrative
	// tasks
	RootToken() string
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
	// MountPathPrefix is an optional prefix to use when mounting the plugin. If
	// omitted the mount path will default to the PluginName with a random suffix.
	MountPathPrefix string

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
	//
	// Test Drivers should normally handle this by tearing down any infrastructure
	// created during the setup of the Vault cluster for testing. The test case
	// specific teardown should only be used if the Driver cluster supports
	// multiple tests running on a single cluster, which is not currently
	// supported. Until this feature is supported by the Driver being used, users
	// should not use this function.
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
		c.PreCheck()
	}

	// Defer on the test case teardown, regardless of pass/fail at this point.
	if c.Teardown != nil {
		defer func() {
			if err := c.Teardown(); err != nil {
				tt.Error("failed to tear down:", err)
			}
		}()
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
			logger.Info("driver Teardown skipped")
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

	// Trap the rootToken so that we can preform revocation or other tasks in the
	// event any steps remove the token during testing.
	rootToken := c.Driver.RootToken()

	// track all responses to revoke any secrets
	var responses []*api.Secret

	// Defer revocation of any secrets created. We intentionally enclose the
	// responses slice so in the event of a fatal error during test evaluation, we
	// are still able to revoke any leases/secrets created
	defer func() {
		// failedRevokes tracks any errors we get when attempting to revoke a lease
		// to log to users at the end of the test.
		var failedRevokes []*api.Secret
		for _, secret := range responses {
			if secret.LeaseID != "" {
				logger.Info("Revoking secret", "lease_id", fmt.Sprintf("%s", secret.LeaseID))
				if err := client.Sys().Revoke(secret.LeaseID); err != nil {
					logger.Warn("Error revoking secret", "lease_id", fmt.Sprintf("%s", secret.LeaseID))
					tt.Error(fmt.Errorf("error revoking lease: %w", err))
					failedRevokes = append(failedRevokes, secret)
					continue
				}
				logger.Info("Successfully revoked secret", "lease_id", fmt.Sprintf("%s", secret.LeaseID))
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
	}()

	stepCount := len(c.Steps)
	for i, step := range c.Steps {
		// range is zero based, so add 1 for a human friendly output of steps.
		// "index" here is only used for logging / output, and not to reference the
		// step in the slice of steps.
		index := i + 1
		if logger.IsWarn() {
			progress := fmt.Sprintf("%d/%d", index, stepCount)
			logger.Warn("Executing test step", "step_number", progress)
		}

		// preserve the root token, which may be cleared from the client if the this
		// step is meant to be unauthenticated. Restored after the request is made.
		// TODO: use a non-root token for tests, or allow a user configured one as
		// part of the test
		if step.Unauthenticated {
			client.ClearToken()
		}

		// ExpandPath will turn a test path into a full path, prefixing with the
		// correct mount, namespaces, or "auth" if needed based on mount path and
		// plugin type
		path := c.Driver.ExpandPath(step.Path)
		var respErr error
		var resp *api.Secret

		switch step.Operation {
		case WriteOperation, UpdateOperation:
			resp, respErr = client.Logical().Write(path, step.Data)
		case ReadOperation:
			// Some operations support reading with data given.
			// TODO: see how the CLI parses args and turns them into
			// map[string][]string, or change how step.Data is defined (currently
			// map[string]interface{})
			// resp, respErr = client.Logical().ReadWithData(path, step.Data)

			resp, respErr = client.Logical().Read(path)
		case ListOperation:
			resp, respErr = client.Logical().List(path)
		case DeleteOperation:
			resp, respErr = client.Logical().Delete(path)
		default:
			panic("bad operation")
		}
		if resp != nil {
			responses = append(responses, resp)
		}

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

		// If a step returned an unexpected error, fail the entire test immediately
		if respErr != nil && !step.ExpectError {
			tt.Fatal(fmt.Errorf("unexpected error in step %d: %w", index, respErr))
		}

		// run the associated StepCheckFunc, if any. If an error was expected it is
		// sent to the Check function to validate.
		if step.Check != nil {
			if err := step.Check(resp, respErr); err != nil {
				tt.Error(fmt.Errorf("Failed step %d: %w", index, err))
			}
		}

		// reset token in case it was cleared
		client.SetToken(rootToken)
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
