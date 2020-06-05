// Package stepwise offers types and functions to enable black-box style tests
// that are executed in defined set of steps. Stepwise utilizes "Drivers" which
// setup a running instance of Vault and provide a valid API client to execute
// user defined steps against.
package stepwise

import (
	"fmt"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

// StepOperation defines operations each step could preform
type StepOperation string

const (
	WriteOperation  StepOperation = "create"
	UpdateOperation               = "update"
	ReadOperation                 = "read"
	DeleteOperation               = "delete"
	ListOperation                 = "list"
	HelpOperation                 = "help"
)

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

	// RootToken returns the root token of the cluster, used for making requests
	// as well as administrative tasks
	RootToken() string
}

// PluginType defines the types of plugins supported
// This type re-create constants as a convienence so users don't need to import/use
// the consts package.
type PluginType consts.PluginType

// These are originally defined in sdk/helper/consts/plugin_types.go
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

// Step represents a single step of a test Case
type Step struct {
	// Operation defines what action is being taken in this step; write, read,
	// delete, et. al.
	Operation StepOperation

	// Path is the localized request path. The mount prefix, namespace, and
	// optionally "auth" will be automatically added.
	Path string

	// Arguments to pass in the request. These arguments represent payloads sent
	// to the API.
	Data map[string]interface{}

	// Check is a function that is called after this step is executed in order to
	// test that the step executed successfully. If this is not set, then the next
	// step will be called
	Check StepCheckFunc

	// ExpectError indicates if this step is expected to return an error. If the
	// step operation returns an error and ExpectError is not set, the test will
	// fail immediately and the StepCheckFunc will not be ran. If ExpectError is
	// true, the StepCheckFunc will be called (if any) and include the error from
	// the response. It is the responsibility of the StepCheckFunc to validate the
	// error is appropriate or not, if expected.
	ExpectError bool

	// Unauthenticated will make the request unauthenticated.
	Unauthenticated bool
}

// StepCheckFunc is the callback used for Check in Steps.
type StepCheckFunc func(*api.Secret, error) error

// Case is a collection of tests to run for a plugin
type Case struct {
	// Driver is used to setup the Vault instance and provide the client used to
	// drive the tests
	Driver StepDriver

	// Precheck, if non-nil, will be called once before the test case
	// runs at all. This can be used for some validation prior to the
	// test running.
	PreCheck func()

	// Steps are the set of operations that are run for this test case.
	Steps []Step

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

	// retrieve the root client from the Driver. If this returns an error, fail
	// immediately
	rootClient, err := c.Driver.Client()
	if err != nil {
		tt.Fatal(err)
	}

	// Trap the rootToken so that we can preform revocation or other tasks in the
	// event any steps remove the token during testing.
	rootToken := c.Driver.RootToken()

	// Defer revocation of any secrets created. We intentionally enclose the
	// responses slice so in the event of a fatal error during test evaluation, we
	// are still able to revoke any leases/secrets created
	var responses []*api.Secret
	defer func() {
		// restore root token for admin tasks
		rootClient.SetToken(rootToken)
		// failedRevokes tracks any errors we get when attempting to revoke a lease
		// to log to users at the end of the test.
		var failedRevokes []*api.Secret
		for _, secret := range responses {
			if secret.LeaseID == "" {
				continue
			}

			if err := rootClient.Sys().Revoke(secret.LeaseID); err != nil {
				tt.Error(fmt.Errorf("error revoking lease: %w", err))
				failedRevokes = append(failedRevokes, secret)
				continue
			}
		}

		// If we have any failed revokes, log it.
		if len(failedRevokes) > 0 {
			for _, s := range failedRevokes {
				tt.Fatal(fmt.Sprintf(
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

		// reset token in case it was cleared
		client, err := rootClient.Clone()
		if err != nil {
			tt.Fatal(err)
		}

		// TODO: support creating tokens with policies listed in each Step
		client.SetToken(rootToken)

		resp, respErr := makeRequest(tt, c.Driver, step)

		if resp != nil {
			responses = append(responses, resp)
		}

		// If a step returned an unexpected error, fail the entire test immediately
		if respErr != nil && !step.ExpectError {
			tt.Fatal(fmt.Errorf("unexpected error in step %d: %w", index, respErr))
		}

		// run the associated StepCheckFunc, if any. If an error was expected it is
		// sent to the Check function to validate.
		if step.Check != nil {
			if err := step.Check(resp, respErr); err != nil {
				tt.Error(fmt.Errorf("failed step %d: %w", index, err))
			}
		}
	}
}

func makeRequest(tt TestT, driver StepDriver, step Step) (*api.Secret, error) {
	client, err := driver.Client()
	if err != nil {
		return nil, err
	}

	if step.Unauthenticated {
		client.ClearToken()
	}

	// ExpandPath will turn a test path into a full path, prefixing with the
	// correct mount, namespaces, or "auth" if needed based on mount path and
	// plugin type
	path := driver.ExpandPath(step.Path)
	switch step.Operation {
	case WriteOperation, UpdateOperation:
		return client.Logical().Write(path, step.Data)
	case ReadOperation:
		// Some operations support reading with data given.
		// TODO: see how the CLI parses args and turns them into
		// map[string][]string, or change how step.Data is defined (currently
		// map[string]interface{})
		// resp, respErr = client.Logical().ReadWithData(path, step.Data)
		return client.Logical().Read(path)
	case ListOperation:
		return client.Logical().List(path)
	case DeleteOperation:
		return client.Logical().Delete(path)
	default:
		return nil, fmt.Errorf("invalid operation: %s", step.Operation)
	}
}
