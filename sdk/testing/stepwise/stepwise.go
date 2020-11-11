// Package stepwise offers types and functions to enable black-box style tests
// that are executed in defined set of steps. Stepwise utilizes "Environments" which
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

// TestEnvVar must be set to a non-empty value for acceptance tests to run.
const TestEnvVar = "VAULT_ACC"

// Operation defines operations each step could perform. These are
// intentionally redefined from the logical package in the SDK, so users
// consistently use the stepwise package and not a combination of both stepwise
// and logical.
type Operation string

const (
	WriteOperation  Operation = "create"
	UpdateOperation           = "update"
	ReadOperation             = "read"
	DeleteOperation           = "delete"
	ListOperation             = "list"
	HelpOperation             = "help"
)

// Environment is the interface Environments need to implement to be used in
// Case to execute each Step
type Environment interface {
	// Setup is responsible for creating the Vault cluster for use in the test
	// case.
	Setup() error

	// Client should return a clone of a configured Vault API client to
	// communicate with the Vault cluster created in Setup and managed by this
	// Environment.
	Client() (*api.Client, error)

	// Teardown is responsible for destroying any and all infrastructure created
	// during Setup or otherwise over the course of executing test cases.
	Teardown() error

	// Name returns the name of the environment provider, e.g. Docker, Minikube,
	// et.al.
	Name() string

	// MountPath returns the path the plugin is mounted at
	MountPath() string

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

// MountOptions are a collection of options each step driver should
// support
type MountOptions struct {
	// MountPathPrefix is an optional prefix to use when mounting the plugin. If
	// omitted the mount path will default to the PluginName with a random suffix.
	MountPathPrefix string

	// Name is used to register the plugin. This can be arbitrary but should be a
	// reasonable value. For an example, if the plugin in test is a secret backend
	// that generates UUIDs with the name "vault-plugin-secrets-uuid", then "uuid"
	// or "test-uuid" would be reasonable. The name is used for lookups in the
	// catalog. See "name" in the "Register Plugin" endpoint docs:
	// - https://www.vaultproject.io/api-docs/system/plugins-catalog#register-plugin
	RegistryName string

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
	Operation Operation

	// Path is the localized request path. The mount prefix, namespace, and
	// optionally "auth" will be automatically added.
	Path string

	// Arguments to pass in the request. These arguments represent payloads sent
	// to the API.
	Data map[string]interface{}

	// Assert is a function that is called after this step is executed in order to
	// test that the step executed successfully. If this is not set, then the next
	// step will be called
	Assert AssertionFunc

	// Unauthenticated will make the request unauthenticated.
	Unauthenticated bool
}

// AssertionFunc is the callback used for Assert in Steps.
type AssertionFunc func(*api.Secret, error) error

// Case represents a scenario we want to test which involves a series of
// steps to be followed sequentially, evaluating the results after each step.
type Case struct {
	// Environment is used to setup the Vault instance and provide the client that
	// will be used to drive the tests
	Environment Environment

	// Precheck enabls a test case to determine if it should run or not
	Precheck func()

	// Steps are the set of operations that are run for this test case. During
	// execution each step will be logged to output with a 1-based index as it is
	// ran, with the first step logged as step '1' and not step '0'.
	Steps []Step

	// SkipTeardown allows the Environment TeardownFunc to be skipped, leaving any
	// infrastructure created after the test exists. This is useful for debugging
	// during plugin development to examine the state of the Vault cluster after a
	// test runs. Depending on the Environment used this could incur costs the
	// user is responsible for.
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
	checkShouldRun(tt)

	if c.Precheck != nil {
		c.Precheck()
	}

	if c.Environment == nil {
		tt.Fatal("nil driver in acceptance test")
		// return here only used during testing when using mockT type, otherwise
		// Fatal will exit
		return
	}

	logger := logging.NewVaultLogger(log.Trace)

	if err := c.Environment.Setup(); err != nil {
		tt.Fatal(err)
	}

	defer func() {
		if c.SkipTeardown {
			logger.Info("driver Teardown skipped")
			return
		}
		if err := c.Environment.Teardown(); err != nil {
			logger.Error("error in driver teardown:", "error", err)
		}
	}()

	// retrieve the root client from the Environment. If this returns an error,
	// fail immediately
	rootClient, err := c.Environment.Client()
	if err != nil {
		tt.Fatal(err)
	}

	// Trap the rootToken so that we can preform revocation or other tasks in the
	// event any steps remove the token during testing.
	rootToken := c.Environment.RootToken()

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
				tt.Error(fmt.Sprintf(
					"WARNING: Revoking the following secret failed. It may\n"+
						"still exist. Please verify:\n\n%#v",
					s))
			}
		}
	}()

	stepCount := len(c.Steps)
	for i, step := range c.Steps {
		if logger.IsWarn() {
			// range is zero based, so add 1 for a human friendly output of steps.
			progress := fmt.Sprintf("%d/%d", i+1, stepCount)
			logger.Warn("Executing test step", "step_number", progress)
		}

		// reset token in case it was cleared
		client, err := rootClient.Clone()
		if err != nil {
			tt.Fatal(err)
		}

		// TODO: support creating tokens with policies listed in each Step
		client.SetToken(rootToken)

		resp, respErr := makeRequest(tt, c.Environment, step)
		if resp != nil {
			responses = append(responses, resp)
		}

		// Run the associated AssertionFunc, if any. If an error was expected it is
		// sent to the Assert function to validate.
		if step.Assert != nil {
			if err := step.Assert(resp, respErr); err != nil {
				tt.Error(fmt.Errorf("failed step %d: %w", i+1, err))
			}
		}
	}
}

func makeRequest(tt TestT, env Environment, step Step) (*api.Secret, error) {
	tt.Helper()
	client, err := env.Client()
	if err != nil {
		return nil, err
	}

	if step.Unauthenticated {
		token := client.Token()
		client.ClearToken()
		// restore the client token after this request completes
		defer func() {
			client.SetToken(token)
		}()
	}

	path := fmt.Sprintf("%s/%s", env.MountPath(), step.Path)
	switch step.Operation {
	case WriteOperation, UpdateOperation:
		return client.Logical().Write(path, step.Data)
	case ReadOperation:
		// TODO support ReadWithData
		return client.Logical().Read(path)
	case ListOperation:
		return client.Logical().List(path)
	case DeleteOperation:
		return client.Logical().Delete(path)
	default:
		return nil, fmt.Errorf("invalid operation: %s", step.Operation)
	}
}

func checkShouldRun(tt TestT) {
	tt.Helper()
	if os.Getenv(TestEnvVar) == "" {
		tt.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			TestEnvVar))
		return
	}

	// We require verbose mode so that the user knows what is going on.
	if !testing.Verbose() {
		tt.Fatal("Acceptance tests must be run with the -v flag on tests")
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
