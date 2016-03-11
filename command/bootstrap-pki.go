package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/mlock"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
)

// BootstrapPKICommand is a Command that mounts a new mount.
type BootstrapPKICommand struct {
	Meta
}

func (c *BootstrapPKICommand) Run(args []string) int {
	var ttl, outputDir string
	var enableMlock bool
	flags := c.Meta.FlagSet("bootstrap-pki", FlagSetDefault)
	flags.StringVar(&ttl, "ttl", "720h", "")
	flags.StringVar(&outputDir, "output-dir", "", "")
	flags.BoolVar(&enableMlock, "enable-mlock", false, "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if outputDir == "" {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf("\nAn output directory via -output-dir is required"))
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nbootstrap-pki expects at least one certificate string"))
		return 1
	}

	// If mlock isn't supported, error. We disable this in
	// dev because it is quite scary to see when first using Vault.
	if enableMlock && !mlock.Supported() {
		c.Ui.Output("==> Error: mlock was enabled, but is not supported on this system!\n")
		c.Ui.Output("    The `mlock` syscall to prevent memory from being swapped to")
		c.Ui.Output("    disk is not supported on this system. Enabling mlock or")
		c.Ui.Output("    running Vault on a system with mlock is much more secure.\n")
		return 1
	}

	// Create an in-memory Vault core
	core, err := vault.NewCore(&vault.CoreConfig{
		Physical: physical.NewInmem(),
		LogicalBackends: map[string]logical.Factory{
			"pki": pki.Factory,
		},
		Logger:       nil,
		DisableMlock: !enableMlock,
	})

	if err != nil {
		c.Ui.Error(fmt.Sprintf("\nerror initializing core: ", err))
		return 1
	}

	// Initialize the core
	init, err := core.Initialize(&vault.SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("\nerror initializing core: ", err))
		return 1
	}

	// Unseal the core
	if unsealed, err := core.Unseal(init.SecretShares[0]); err != nil {
		c.Ui.Error(fmt.Sprintf("\nerror unsealing core: ", err))
		return 1
	} else if !unsealed {
		c.Ui.Error("\ncore could not be unsealed")
		return 1
	}

	systemBackend := vault.NewSystemBackend(core,
		&logical.BackendConfig{
			Logger: nil,
			System: logical.TestSystemView(),
		})

	// Create an HTTP API server and client
	ln, addr := http.TestServer(nil, core)
	defer ln.Close()
	clientConfig := api.DefaultConfig()
	clientConfig.Address = addr
	client, err := api.NewClient(clientConfig)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error initializing HTTP client: %s", err))
		return 1
	}

	// Set the token so we're authenticated
	client.SetToken(init.RootToken)

	// Mount the backend
	prefix := "pki"
	mountInfo := &api.MountInput{
		Type: "pki",
	}
	if err := client.Sys().Mount(prefix, mountInfo); err != nil {
		c.Ui.Error(fmt.Sprintf("error mounting backend: %s", err))
		return 1
	}

	req := &logical.Request{
		ClientToken: init.RootToken,
		Operation:   logical.UpdateOperation,
		Path:        "mounts/pki/tune",
		Data: map[string]interface{}{
			"default_lease_ttl": ttl,
			// Make this super long; we'll always use the default, but this
			// needs to be sure to be longer
			"max_lease_ttl": "876600h",
		},
	}
	_, err = systemBackend.HandleRequest(req)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error tuning PKI mount: %s", err))
		return 1
	}

	rando, err := uuid.GenerateUUID()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error generating random UUID for CA cert: %s", err))
		return 1
	}

	req.Path = "pki/root/generate/internal"
	req.Data = map[string]interface{}{
		"ttl":         ttl,
		"common_name": fmt.Sprintf("Vault Bootstrap Root CA %s", rando),
	}
	resp, err = core.HandleRequest(req)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error generating CA cert: %s", err))
	}
	caCert := resp.Data["certificate"].(string)

	c.Ui.Error(fmt.Sprintf("%#v", caCert))

	return 0
}

/*
	// Make requests
	var revoke []*logical.Request
	for i, s := range c.Steps {
		log.Printf("[WARN] Executing test step %d", i+1)

		// Create the request
		req := &logical.Request{
			Operation: s.Operation,
			Path:      s.Path,
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

		if s.PreFlight != nil {
			ct := req.ClientToken
			req.ClientToken = ""
			if err := s.PreFlight(req); err != nil {
				t.Error(fmt.Sprintf("Failed preflight for step %d: %s", i+1, err))
				break
			}
			req.ClientToken = ct
		}

		// Make sure to prefix the path with where we mounted the thing
		req.Path = fmt.Sprintf("%s/%s", prefix, req.Path)

		// Make the request
		resp, err := core.HandleRequest(req)
		if resp != nil && resp.Secret != nil {
			// Revoke this secret later
			revoke = append(revoke, &logical.Request{
	// If no path is specified, we default the path to the backend type
	if path == "" {
		path = mountType
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	mountInfo := &api.MountInput{
		Type:        mountType,
		Description: description,
		Config: api.MountConfigInput{
			DefaultLeaseTTL: defaultLeaseTTL,
			MaxLeaseTTL:     maxLeaseTTL,
		},
	}

	if err := client.Sys().Mount(path, mountInfo); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Mount error: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully mounted '%s' at '%s'!",
		mountType, path))

	return 0
}
*/

func (c *BootstrapPKICommand) Synopsis() string {
	return "Mount a logical backend"
}

func (c *BootstrapPKICommand) Help() string {
	helpText := `
Usage: vault mount [options] type

  Mount a logical backend.

  This command mounts a logical backend for storing and/or generating
  secrets.

General Options:

  ` + generalOptionsUsage() + `

Mount Options:

  -description=<desc>            Human-friendly description of the purpose for
                                 the mount. This shows up in the mounts command.

  -path=<path>                   Mount point for the logical backend. This
                                 defauls to the type of the mount.

  -default-lease-ttl=<duration>  Default lease time-to-live for this backend.
                                 If not specified, uses the global default, or
                                 the previously set value. Set to '0' to
                                 explicitly set it to use the global default.

  -max-lease-ttl=<duration>      Max lease time-to-live for this backend.
                                 If not specified, uses the global default, or
                                 the previously set value. Set to '0' to
                                 explicitly set it to use the global default.

`
	return strings.TrimSpace(helpText)
}
