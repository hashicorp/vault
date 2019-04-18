package command

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"

	consulapi "github.com/hashicorp/consul/api"
)

var _ cli.Command = (*OperatorInitCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorInitCommand)(nil)

type OperatorInitCommand struct {
	*BaseCommand

	flagStatus          bool
	flagKeyShares       int
	flagKeyThreshold    int
	flagPGPKeys         []string
	flagRootTokenPGPKey string

	// Auto Unseal
	flagRecoveryShares    int
	flagRecoveryThreshold int
	flagRecoveryPGPKeys   []string

	// Consul
	flagConsulAuto    bool
	flagConsulService string
}

const (
	defKeyShares    = 5
	defKeyThreshold = 3
)

func (c *OperatorInitCommand) Synopsis() string {
	return "Initializes a server"
}

func (c *OperatorInitCommand) Help() string {
	helpText := `
Usage: vault operator init [options]

  Initializes a Vault server. Initialization is the process by which Vault's
  storage backend is prepared to receive data. Since Vault servers share the
  same storage backend in HA mode, you only need to initialize one Vault to
  initialize the storage backend.

  During initialization, Vault generates an in-memory master key and applies
  Shamir's secret sharing algorithm to disassemble that master key into a
  configuration number of key shares such that a configurable subset of those
  key shares must come together to regenerate the master key. These keys are
  often called "unseal keys" in Vault's documentation.

  This command cannot be run against an already-initialized Vault cluster.

  Start initialization with the default options:

      $ vault operator init

  Initialize, but encrypt the unseal keys with pgp keys:

      $ vault operator init \
          -key-shares=3 \
          -key-threshold=2 \
          -pgp-keys="keybase:hashicorp,keybase:jefferai,keybase:sethvargo"

  Encrypt the initial root token using a pgp key:

      $ vault operator init -root-token-pgp-key="keybase:hashicorp"

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *OperatorInitCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.BoolVar(&BoolVar{
		Name:    "status",
		Target:  &c.flagStatus,
		Default: false,
		Usage: "Print the current initialization status. An exit code of 0 means " +
			"the Vault is already initialized. An exit code of 1 means an error " +
			"occurred. An exit code of 2 means the mean is not initialized.",
	})

	f.IntVar(&IntVar{
		Name:       "key-shares",
		Aliases:    []string{"n"},
		Target:     &c.flagKeyShares,
		Default:    defKeyShares,
		Completion: complete.PredictAnything,
		Usage: "Number of key shares to split the generated master key into. " +
			"This is the number of \"unseal keys\" to generate.",
	})

	f.IntVar(&IntVar{
		Name:       "key-threshold",
		Aliases:    []string{"t"},
		Target:     &c.flagKeyThreshold,
		Default:    defKeyThreshold,
		Completion: complete.PredictAnything,
		Usage: "Number of key shares required to reconstruct the master key. " +
			"This must be less than or equal to -key-shares.",
	})

	f.VarFlag(&VarFlag{
		Name:       "pgp-keys",
		Value:      (*pgpkeys.PubKeyFilesFlag)(&c.flagPGPKeys),
		Completion: complete.PredictAnything,
		Usage: "Comma-separated list of paths to files on disk containing " +
			"public GPG keys OR a comma-separated list of Keybase usernames using " +
			"the format \"keybase:<username>\". When supplied, the generated " +
			"unseal keys will be encrypted and base64-encoded in the order " +
			"specified in this list. The number of entries must match -key-shares, " +
			"unless -store-shares are used.",
	})

	f.VarFlag(&VarFlag{
		Name:       "root-token-pgp-key",
		Value:      (*pgpkeys.PubKeyFileFlag)(&c.flagRootTokenPGPKey),
		Completion: complete.PredictAnything,
		Usage: "Path to a file on disk containing a binary or base64-encoded " +
			"public GPG key. This can also be specified as a Keybase username " +
			"using the format \"keybase:<username>\". When supplied, the generated " +
			"root token will be encrypted and base64-encoded with the given public " +
			"key.",
	})

	// Consul Options
	f = set.NewFlagSet("Consul Options")

	f.BoolVar(&BoolVar{
		Name:    "consul-auto",
		Target:  &c.flagConsulAuto,
		Default: false,
		Usage: "Perform automatic service discovery using Consul in HA mode. " +
			"When all nodes in a Vault HA cluster are registered with Consul, " +
			"enabling this option will trigger automatic service discovery based " +
			"on the provided -consul-service value. When Consul is Vault's HA " +
			"backend, this functionality is automatically enabled. Ensure the " +
			"proper Consul environment variables are set (CONSUL_HTTP_ADDR, etc). " +
			"When only one Vault server is discovered, it will be initialized " +
			"automatically. When more than one Vault server is discovered, they " +
			"will each be output for selection.",
	})

	f.StringVar(&StringVar{
		Name:       "consul-service",
		Target:     &c.flagConsulService,
		Default:    "vault",
		Completion: complete.PredictAnything,
		Usage: "Name of the service in Consul under which the Vault servers are " +
			"registered.",
	})

	// Auto Unseal Options
	f = set.NewFlagSet("Auto Unseal Options")

	f.IntVar(&IntVar{
		Name:       "recovery-shares",
		Target:     &c.flagRecoveryShares,
		Default:    5,
		Completion: complete.PredictAnything,
		Usage: "Number of key shares to split the recovery key into. " +
			"This is only used in auto-unseal mode.",
	})

	f.IntVar(&IntVar{
		Name:       "recovery-threshold",
		Target:     &c.flagRecoveryThreshold,
		Default:    3,
		Completion: complete.PredictAnything,
		Usage: "Number of key shares required to reconstruct the recovery key. " +
			"This is only used in Auto Unseal mode.",
	})

	f.VarFlag(&VarFlag{
		Name:       "recovery-pgp-keys",
		Value:      (*pgpkeys.PubKeyFilesFlag)(&c.flagRecoveryPGPKeys),
		Completion: complete.PredictAnything,
		Usage: "Behaves like -pgp-keys, but for the recovery key shares. This " +
			"is only used in Auto Unseal mode.",
	})

	return set
}

func (c *OperatorInitCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *OperatorInitCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorInitCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	// Build the initial init request
	initReq := &api.InitRequest{
		SecretShares:    c.flagKeyShares,
		SecretThreshold: c.flagKeyThreshold,
		PGPKeys:         c.flagPGPKeys,
		RootTokenPGPKey: c.flagRootTokenPGPKey,

		RecoveryShares:    c.flagRecoveryShares,
		RecoveryThreshold: c.flagRecoveryThreshold,
		RecoveryPGPKeys:   c.flagRecoveryPGPKeys,
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Check auto mode
	switch {
	case c.flagStatus:
		return c.status(client)
	case c.flagConsulAuto:
		return c.consulAuto(client, initReq)
	default:
		return c.init(client, initReq)
	}
}

// consulAuto enables auto-joining via Consul.
func (c *OperatorInitCommand) consulAuto(client *api.Client, req *api.InitRequest) int {
	// Capture the client original address and reset it
	originalAddr := client.Address()
	defer client.SetAddress(originalAddr)

	// Create a client to communicate with Consul
	consulClient, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to create Consul client:%v", err))
		return 1
	}

	// Pull the scheme from the Vault client to determine if the Consul agent
	// should talk via HTTP or HTTPS.
	addr := client.Address()
	clientURL, err := url.Parse(addr)
	if err != nil || clientURL == nil {
		c.UI.Error(fmt.Sprintf("Failed to parse Vault address %s: %s", addr, err))
		return 1
	}

	var uninitedVaults []string
	var initedVault string

	// Query the nodes belonging to the cluster
	services, _, err := consulClient.Catalog().Service(c.flagConsulService, "", &consulapi.QueryOptions{
		AllowStale: true,
	})
	if err == nil {
		for _, service := range services {
			// Set the address on the client temporarily
			vaultAddr := (&url.URL{
				Scheme: clientURL.Scheme,
				Host:   fmt.Sprintf("%s:%d", service.ServiceAddress, service.ServicePort),
			}).String()
			client.SetAddress(vaultAddr)

			// Check the initialization status of the discovered node
			inited, err := client.Sys().InitStatus()
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error checking init status of %q: %s", vaultAddr, err))
			}
			if inited {
				initedVault = vaultAddr
				break
			}

			// If we got this far, we communicated successfully with Vault, but it
			// was not initialized.
			uninitedVaults = append(uninitedVaults, vaultAddr)
		}
	}

	// Get the correct export keywords and quotes for *nix vs Windows
	export := "export"
	quote := "\""
	if runtime.GOOS == "windows" {
		export = "set"
		quote = ""
	}

	if initedVault != "" {
		vaultURL, err := url.Parse(initedVault)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Failed to parse Vault address %q: %s", initedVault, err))
			return 2
		}
		vaultAddr := vaultURL.String()

		c.UI.Output(wrapAtLength(fmt.Sprintf(
			"Discovered an initialized Vault node at %q with Consul service name "+
				"%q. Set the following environment variable to target the discovered "+
				"Vault server:",
			vaultURL.String(), c.flagConsulService)))
		c.UI.Output("")
		c.UI.Output(fmt.Sprintf("    $ %s VAULT_ADDR=%s%s%s", export, quote, vaultAddr, quote))
		c.UI.Output("")
		return 0
	}

	switch len(uninitedVaults) {
	case 0:
		c.UI.Error(fmt.Sprintf("No Vault nodes registered as %q in Consul", c.flagConsulService))
		return 2
	case 1:
		// There was only one node found in the Vault cluster and it was
		// uninitialized.
		vaultURL, err := url.Parse(uninitedVaults[0])
		if err != nil {
			c.UI.Error(fmt.Sprintf("Failed to parse Vault address %q: %s", initedVault, err))
			return 2
		}
		vaultAddr := vaultURL.String()

		// Update the client to connect to this Vault server
		client.SetAddress(vaultAddr)

		// Let the client know that initialization is performed on the
		// discovered node.
		c.UI.Output(wrapAtLength(fmt.Sprintf(
			"Discovered an initialized Vault node at %q with Consul service name "+
				"%q. Set the following environment variable to target the discovered "+
				"Vault server:",
			vaultURL.String(), c.flagConsulService)))
		c.UI.Output("")
		c.UI.Output(fmt.Sprintf("    $ %s VAULT_ADDR=%s%s%s", export, quote, vaultAddr, quote))
		c.UI.Output("")
		c.UI.Output("Attempting to initialize it...")
		c.UI.Output("")

		// Attempt to initialize it
		return c.init(client, req)
	default:
		// If more than one Vault node were discovered, print out all of them,
		// requiring the client to update VAULT_ADDR and to run init again.
		c.UI.Output(wrapAtLength(fmt.Sprintf(
			"Discovered %d uninitialized Vault servers with Consul service name "+
				"%q. To initialize these Vaults, set any one of the following "+
				"environment variables and run \"vault operator init\":",
			len(uninitedVaults), c.flagConsulService)))
		c.UI.Output("")

		// Print valid commands to make setting the variables easier
		for _, node := range uninitedVaults {
			vaultURL, err := url.Parse(node)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Failed to parse Vault address %q: %s", initedVault, err))
				return 2
			}
			vaultAddr := vaultURL.String()

			c.UI.Output(fmt.Sprintf("    $ %s VAULT_ADDR=%s%s%s", export, quote, vaultAddr, quote))
		}

		c.UI.Output("")
		return 0
	}
}

func (c *OperatorInitCommand) init(client *api.Client, req *api.InitRequest) int {
	resp, err := client.Sys().Init(req)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing: %s", err))
		return 2
	}

	switch Format(c.UI) {
	case "table":
	default:
		return OutputData(c.UI, newMachineInit(req, resp))
	}

	for i, key := range resp.Keys {
		if resp.KeysB64 != nil && len(resp.KeysB64) == len(resp.Keys) {
			c.UI.Output(fmt.Sprintf("Unseal Key %d: %s", i+1, resp.KeysB64[i]))
		} else {
			c.UI.Output(fmt.Sprintf("Unseal Key %d: %s", i+1, key))
		}
	}
	for i, key := range resp.RecoveryKeys {
		if resp.RecoveryKeysB64 != nil && len(resp.RecoveryKeysB64) == len(resp.RecoveryKeys) {
			c.UI.Output(fmt.Sprintf("Recovery Key %d: %s", i+1, resp.RecoveryKeysB64[i]))
		} else {
			c.UI.Output(fmt.Sprintf("Recovery Key %d: %s", i+1, key))
		}
	}

	c.UI.Output("")
	c.UI.Output(fmt.Sprintf("Initial Root Token: %s", resp.RootToken))

	if len(resp.Keys) > 0 {
		c.UI.Output("")
		c.UI.Output(wrapAtLength(fmt.Sprintf(
			"Vault initialized with %d key shares and a key threshold of %d. Please "+
				"securely distribute the key shares printed above. When the Vault is "+
				"re-sealed, restarted, or stopped, you must supply at least %d of "+
				"these keys to unseal it before it can start servicing requests.",
			req.SecretShares,
			req.SecretThreshold,
			req.SecretThreshold)))

		c.UI.Output("")
		c.UI.Output(wrapAtLength(fmt.Sprintf(
			"Vault does not store the generated master key. Without at least %d "+
				"key to reconstruct the master key, Vault will remain permanently "+
				"sealed!",
			req.SecretThreshold)))

		c.UI.Output("")
		c.UI.Output(wrapAtLength(
			"It is possible to generate new unseal keys, provided you have a quorum " +
				"of existing unseal keys shares. See \"vault operator rekey\" for " +
				"more information."))
	} else {
		c.UI.Output("")
		c.UI.Output("Success! Vault is initialized")
	}

	if len(resp.RecoveryKeys) > 0 {
		c.UI.Output("")
		c.UI.Output(wrapAtLength(fmt.Sprintf(
			"Recovery key initialized with %d key shares and a key threshold of %d. "+
				"Please securely distribute the key shares printed above.",
			req.RecoveryShares,
			req.RecoveryThreshold)))
	}

	if len(resp.RecoveryKeys) > 0 && (req.SecretShares != defKeyShares || req.SecretThreshold != defKeyThreshold) {
		c.UI.Output("")
		c.UI.Warn(wrapAtLength(
			"WARNING! -key-shares and -key-threshold is ignored when " +
				"Auto Unseal is used. Use -recovery-shares and -recovery-threshold instead.",
		))
	}

	return 0
}

// status inspects the init status of vault and returns an appropriate error
// code and message.
func (c *OperatorInitCommand) status(client *api.Client) int {
	inited, err := client.Sys().InitStatus()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error checking init status: %s", err))
		return 1 // Normally we'd return 2, but 2 means something special here
	}

	if inited {
		c.UI.Output("Vault is initialized")
		return 0
	}

	c.UI.Output("Vault is not initialized")
	return 2
}

// machineInit is used to output information about the init command.
type machineInit struct {
	UnsealKeysB64     []string `json:"unseal_keys_b64"`
	UnsealKeysHex     []string `json:"unseal_keys_hex"`
	UnsealShares      int      `json:"unseal_shares"`
	UnsealThreshold   int      `json:"unseal_threshold"`
	RecoveryKeysB64   []string `json:"recovery_keys_b64"`
	RecoveryKeysHex   []string `json:"recovery_keys_hex"`
	RecoveryShares    int      `json:"recovery_keys_shares"`
	RecoveryThreshold int      `json:"recovery_keys_threshold"`
	RootToken         string   `json:"root_token"`
}

func newMachineInit(req *api.InitRequest, resp *api.InitResponse) *machineInit {
	init := &machineInit{}

	init.UnsealKeysHex = make([]string, len(resp.Keys))
	for i, v := range resp.Keys {
		init.UnsealKeysHex[i] = v
	}

	init.UnsealKeysB64 = make([]string, len(resp.KeysB64))
	for i, v := range resp.KeysB64 {
		init.UnsealKeysB64[i] = v
	}

	// If we don't get a set of keys back, it means that we are storing the keys,
	// so the key shares and threshold has been set to 1.
	if len(resp.Keys) == 0 {
		init.UnsealShares = 1
		init.UnsealThreshold = 1
	} else {
		init.UnsealShares = req.SecretShares
		init.UnsealThreshold = req.SecretThreshold
	}

	init.RecoveryKeysHex = make([]string, len(resp.RecoveryKeys))
	for i, v := range resp.RecoveryKeys {
		init.RecoveryKeysHex[i] = v
	}

	init.RecoveryKeysB64 = make([]string, len(resp.RecoveryKeysB64))
	for i, v := range resp.RecoveryKeysB64 {
		init.RecoveryKeysB64[i] = v
	}

	init.RecoveryShares = req.RecoveryShares
	init.RecoveryThreshold = req.RecoveryThreshold

	init.RootToken = resp.RootToken

	return init
}
