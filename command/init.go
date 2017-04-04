package command

import (
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/physical"
)

// InitCommand is a Command that initializes a new Vault server.
type InitCommand struct {
	meta.Meta
}

func (c *InitCommand) Run(args []string) int {
	var threshold, shares, storedShares, recoveryThreshold, recoveryShares int
	var pgpKeys, recoveryPgpKeys, rootTokenPgpKey pgpkeys.PubKeyFilesFlag
	var auto, check bool
	var consulServiceName string
	flags := c.Meta.FlagSet("init", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	flags.IntVar(&shares, "key-shares", 5, "")
	flags.IntVar(&threshold, "key-threshold", 3, "")
	flags.IntVar(&storedShares, "stored-shares", 0, "")
	flags.Var(&pgpKeys, "pgp-keys", "")
	flags.Var(&rootTokenPgpKey, "root-token-pgp-key", "")
	flags.IntVar(&recoveryShares, "recovery-shares", 5, "")
	flags.IntVar(&recoveryThreshold, "recovery-threshold", 3, "")
	flags.Var(&recoveryPgpKeys, "recovery-pgp-keys", "")
	flags.BoolVar(&check, "check", false, "")
	flags.BoolVar(&auto, "auto", false, "")
	flags.StringVar(&consulServiceName, "consul-service", physical.DefaultServiceName, "")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	initRequest := &api.InitRequest{
		SecretShares:      shares,
		SecretThreshold:   threshold,
		StoredShares:      storedShares,
		PGPKeys:           pgpKeys,
		RecoveryShares:    recoveryShares,
		RecoveryThreshold: recoveryThreshold,
		RecoveryPGPKeys:   recoveryPgpKeys,
	}

	switch len(rootTokenPgpKey) {
	case 0:
	case 1:
		initRequest.RootTokenPGPKey = rootTokenPgpKey[0]
	default:
		c.Ui.Error("Only one PGP key can be specified for encrypting the root token")
		return 1
	}

	// If running in 'auto' mode, run service discovery based on environment
	// variables of Consul.
	if auto {

		// Create configuration for Consul
		consulConfig := consulapi.DefaultConfig()

		// Create a client to communicate with Consul
		consulClient, err := consulapi.NewClient(consulConfig)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to create Consul client:%v", err))
			return 1
		}

		// Fetch Vault's protocol scheme from the client
		vaultclient, err := c.Client()
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to fetch Vault client: %v", err))
			return 1
		}

		if vaultclient.Address() == "" {
			c.Ui.Error("Failed to fetch Vault client address")
			return 1
		}

		clientURL, err := url.Parse(vaultclient.Address())
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to parse Vault address: %v", err))
			return 1
		}

		if clientURL == nil {
			c.Ui.Error("Failed to parse Vault client address")
			return 1
		}

		var uninitializedVaults []string
		var initializedVault string

		// Query the nodes belonging to the cluster
		if services, _, err := consulClient.Catalog().Service(consulServiceName, "", &consulapi.QueryOptions{AllowStale: true}); err == nil {
		Loop:
			for _, service := range services {
				vaultAddress := &url.URL{
					Scheme: clientURL.Scheme,
					Host:   fmt.Sprintf("%s:%d", service.ServiceAddress, service.ServicePort),
				}

				// Set VAULT_ADDR to the discovered node
				os.Setenv(api.EnvVaultAddress, vaultAddress.String())

				// Create a client to communicate with the discovered node
				client, err := c.Client()
				if err != nil {
					c.Ui.Error(fmt.Sprintf("Error initializing client: %v", err))
					return 1
				}

				// Check the initialization status of the discovered node
				inited, err := client.Sys().InitStatus()
				switch {
				case err != nil:
					c.Ui.Error(fmt.Sprintf("Error checking initialization status of discovered node: %+q. Err: %v", vaultAddress.String(), err))
					return 1
				case inited:
					// One of the nodes in the cluster is initialized. Break out.
					initializedVault = vaultAddress.String()
					break Loop
				default:
					// Vault is uninitialized.
					uninitializedVaults = append(uninitializedVaults, vaultAddress.String())
				}
			}
		}

		export := "export"
		quote := "'"
		if runtime.GOOS == "windows" {
			export = "set"
			quote = ""
		}

		if initializedVault != "" {
			vaultURL, err := url.Parse(initializedVault)
			if err != nil {
				c.Ui.Error(fmt.Sprintf("Failed to parse Vault address: %+q. Err: %v", initializedVault, err))
			}
			c.Ui.Output(fmt.Sprintf("Discovered an initialized Vault node at %+q, using Consul service name %+q", vaultURL.String(), consulServiceName))
			c.Ui.Output("\nSet the following environment variable to operate on the discovered Vault:\n")
			c.Ui.Output(fmt.Sprintf("\t%s VAULT_ADDR=%s%s%s", export, quote, vaultURL.String(), quote))
			return 0
		}

		switch len(uninitializedVaults) {
		case 0:
			c.Ui.Error(fmt.Sprintf("Failed to discover Vault nodes using Consul service name %+q", consulServiceName))
			return 1
		case 1:
			// There was only one node found in the Vault cluster and it
			// was uninitialized.

			vaultURL, err := url.Parse(uninitializedVaults[0])
			if err != nil {
				c.Ui.Error(fmt.Sprintf("Failed to parse Vault address: %+q. Err: %v", uninitializedVaults[0], err))
			}

			// Set the VAULT_ADDR to the discovered node. This will ensure
			// that the client created will operate on the discovered node.
			os.Setenv(api.EnvVaultAddress, vaultURL.String())

			// Let the client know that initialization is perfomed on the
			// discovered node.
			c.Ui.Output(fmt.Sprintf("Discovered Vault at %+q using Consul service name %+q\n", vaultURL.String(), consulServiceName))

			// Attempt initializing it
			ret := c.runInit(check, initRequest)

			// Regardless of success or failure, instruct client to update VAULT_ADDR
			c.Ui.Output("\nSet the following environment variable to operate on the discovered Vault:\n")
			c.Ui.Output(fmt.Sprintf("\t%s VAULT_ADDR=%s%s%s", export, quote, vaultURL.String(), quote))

			return ret
		default:
			// If more than one Vault node were discovered, print out all of them,
			// requiring the client to update VAULT_ADDR and to run init again.
			c.Ui.Output(fmt.Sprintf("Discovered more than one uninitialized Vaults using Consul service name %+q\n", consulServiceName))
			c.Ui.Output("To initialize these Vaults, set any *one* of the following environment variables and run 'vault init':")

			// Print valid commands to make setting the variables easier
			for _, vaultNode := range uninitializedVaults {
				vaultURL, err := url.Parse(vaultNode)
				if err != nil {
					c.Ui.Error(fmt.Sprintf("Failed to parse Vault address: %+q. Err: %v", vaultNode, err))
				}
				c.Ui.Output(fmt.Sprintf("\t%s VAULT_ADDR=%s%s%s", export, quote, vaultURL.String(), quote))

			}
			return 0
		}
	}

	return c.runInit(check, initRequest)
}

func (c *InitCommand) runInit(check bool, initRequest *api.InitRequest) int {
	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 1
	}

	if check {
		return c.checkStatus(client)
	}

	resp, err := client.Sys().Init(initRequest)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing Vault: %s", err))
		return 1
	}

	for i, key := range resp.Keys {
		if resp.KeysB64 != nil && len(resp.KeysB64) == len(resp.Keys) {
			c.Ui.Output(fmt.Sprintf("Unseal Key %d: %s", i+1, resp.KeysB64[i]))
		} else {
			c.Ui.Output(fmt.Sprintf("Unseal Key %d: %s", i+1, key))
		}
	}
	for i, key := range resp.RecoveryKeys {
		if resp.RecoveryKeysB64 != nil && len(resp.RecoveryKeysB64) == len(resp.RecoveryKeys) {
			c.Ui.Output(fmt.Sprintf("Recovery Key %d: %s", i+1, resp.RecoveryKeysB64[i]))
		} else {
			c.Ui.Output(fmt.Sprintf("Recovery Key %d: %s", i+1, key))
		}
	}

	c.Ui.Output(fmt.Sprintf("Initial Root Token: %s", resp.RootToken))

	if initRequest.StoredShares < 1 {
		c.Ui.Output(fmt.Sprintf(
			"\n"+
				"Vault initialized with %d keys and a key threshold of %d. Please\n"+
				"securely distribute the above keys. When the vault is re-sealed,\n"+
				"restarted, or stopped, you must provide at least %d of these keys\n"+
				"to unseal it again.\n\n"+
				"Vault does not store the master key. Without at least %d keys,\n"+
				"your vault will remain permanently sealed.",
			initRequest.SecretShares,
			initRequest.SecretThreshold,
			initRequest.SecretThreshold,
			initRequest.SecretThreshold,
		))
	} else {
		c.Ui.Output(
			"\n" +
				"Vault initialized successfully.",
		)
	}
	if len(resp.RecoveryKeys) > 0 {
		c.Ui.Output(fmt.Sprintf(
			"\n"+
				"Recovery key initialized with %d keys and a key threshold of %d. Please\n"+
				"securely distribute the above keys.",
			initRequest.RecoveryShares,
			initRequest.RecoveryThreshold,
		))
	}

	return 0
}

func (c *InitCommand) checkStatus(client *api.Client) int {
	inited, err := client.Sys().InitStatus()
	switch {
	case err != nil:
		c.Ui.Error(fmt.Sprintf(
			"Error checking initialization status: %s", err))
		return 1
	case inited:
		c.Ui.Output("Vault has been initialized")
		return 0
	default:
		c.Ui.Output("Vault is not initialized")
		return 2
	}
}

func (c *InitCommand) Synopsis() string {
	return "Initialize a new Vault server"
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: vault init [options]

  Initialize a new Vault server.

  This command connects to a Vault server and initializes it for the
  first time. This sets up the initial set of master keys and the
  backend data store structure.

  This command can't be called on an already-initialized Vault server.

General Options:
` + meta.GeneralOptionsUsage() + `
Init Options:

  -check                    Don't actually initialize, just check if Vault is
                            already initialized. A return code of 0 means Vault
                            is initialized; a return code of 2 means Vault is not
                            initialized; a return code of 1 means an error was
                            encountered.

  -key-shares=5             The number of key shares to split the master key
                            into.

  -key-threshold=3          The number of key shares required to reconstruct
                            the master key.

  -stored-shares=0          The number of unseal keys to store. Only used with 
                            Vault HSM. Must currently be equivalent to the
                            number of shares.

  -pgp-keys                 If provided, must be a comma-separated list of
                            files on disk containing binary- or base64-format
                            public PGP keys, or Keybase usernames specified as
                            "keybase:<username>". The output unseal keys will
                            be encrypted and base64-encoded, in order, with the
                            given public keys. If you want to use them with the
                            'vault unseal' command, you will need to base64-
                            decode and decrypt; this will be the plaintext
                            unseal key. When 'stored-shares' are not used, the
                            number of entries in this field must match 'key-shares'.
                            When 'stored-shares' are used, the number of entries
                            should match the difference between 'key-shares'
                            and 'stored-shares'.

  -root-token-pgp-key       If provided, a file on disk with a binary- or
                            base64-format public PGP key, or a Keybase username
                            specified as "keybase:<username>". The output root
                            token will be encrypted and base64-encoded, in
                            order, with the given public key. You will need
                            to base64-decode and decrypt the result.

  -recovery-shares=5        The number of key shares to split the recovery key
                            into. Only used with Vault HSM.

  -recovery-threshold=3     The number of key shares required to reconstruct
                            the recovery key. Only used with Vault HSM.

  -recovery-pgp-keys        If provided, behaves like "pgp-keys" but for the
                            recovery key shares. Only used with Vault HSM.

  -auto                     If set, performs service discovery using Consul. 
                            When all the nodes of a Vault cluster are
                            registered with Consul, setting this flag will
                            trigger service discovery using the service name
                            with which Vault nodes are registered. This option
                            works well when each Vault cluster is registered
                            under a unique service name. Note that, when Consul
                            is serving as Vault's HA backend, Vault nodes are
                            registered with Consul by default. The service name
                            can be changed using 'consul-service' flag. Ensure
                            that environment variables required to communicate
                            with Consul, like (CONSUL_HTTP_ADDR,
                            CONSUL_HTTP_TOKEN, CONSUL_HTTP_SSL, et al) are
                            properly set. When only one Vault node is
                            discovered, it will be initialized and when more
                            than one Vault node is discovered, they will be
                            output for easy selection.

  -consul-service           Service name under which all the nodes of a Vault
                            cluster are registered with Consul. Note that, when
                            Vault uses Consul as its HA backend, by default,
                            Vault will register itself as a service with Consul
                            with the service name "vault". This name can be
                            modified in Vault's configuration file, using the
                            "service" option for the Consul backend.
`
	return strings.TrimSpace(helpText)
}
