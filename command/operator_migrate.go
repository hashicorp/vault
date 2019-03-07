package command

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
	"github.com/pkg/errors"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorMigrateCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorMigrateCommand)(nil)

var errAbort = errors.New("Migration aborted")

type OperatorMigrateCommand struct {
	*BaseCommand

	PhysicalBackends map[string]physical.Factory
	flagConfig       string
	flagStart        string
	flagReset        bool
	logger           log.Logger
	ShutdownCh       chan struct{}
}

type migratorConfig struct {
	StorageSource      *server.Storage `hcl:"-"`
	StorageDestination *server.Storage `hcl:"-"`
}

func (c *OperatorMigrateCommand) Synopsis() string {
	return "Migrates Vault data between storage backends"
}

func (c *OperatorMigrateCommand) Help() string {
	helpText := `
Usage: vault operator migrate [options]

  This command starts a storage backend migration process to copy all data
  from one backend to another. This operates directly on encrypted data and
  does not require a Vault server, nor any unsealing.

  Start a migration with a configuration file:

      $ vault operator migrate -config=migrate.hcl

  For more information, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorMigrateCommand) Flags() *FlagSets {
	set := NewFlagSets(c.UI)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:   "config",
		Target: &c.flagConfig,
		Completion: complete.PredictOr(
			complete.PredictFiles("*.hcl"),
		),
		Usage: "Path to a configuration file. This configuration file should " +
			"contain only migrator directives.",
	})

	f.StringVar(&StringVar{
		Name:   "start",
		Target: &c.flagStart,
		Usage:  "Only copy keys lexicographically at or after this value.",
	})

	f.BoolVar(&BoolVar{
		Name:   "reset",
		Target: &c.flagReset,
		Usage:  "Reset the migration lock. No migration will occur.",
	})

	return set
}

func (c *OperatorMigrateCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *OperatorMigrateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorMigrateCommand) Run(args []string) int {
	c.logger = logging.NewVaultLogger(log.Info)
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if c.flagConfig == "" {
		c.UI.Error("Must specify exactly one config path using -config")
		return 1
	}

	config, err := c.loadMigratorConfig(c.flagConfig)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error loading configuration from %s: %s", c.flagConfig, err))
		return 1
	}

	if err := c.migrate(config); err != nil {
		if err == errAbort {
			return 0
		}
		c.UI.Error(fmt.Sprintf("Error migrating: %s", err))
		return 2
	}

	if c.flagReset {
		c.UI.Output("Success! Migration lock reset (if it was set).")
	} else {
		c.UI.Output("Success! All of the keys have been migrated.")
	}

	return 0
}

// migrate attempts to instantiate the source and destinations backends,
// and then invoke the migration the the root of the keyspace.
func (c *OperatorMigrateCommand) migrate(config *migratorConfig) error {
	from, err := c.newBackend(config.StorageSource.Type, config.StorageSource.Config)
	if err != nil {
		return errwrap.Wrapf("error mounting 'storage_source': {{err}}", err)
	}

	if c.flagReset {
		if err := SetStorageMigration(from, false); err != nil {
			return errwrap.Wrapf("error reseting migration lock: {{err}}", err)
		}
		return nil
	}

	to, err := c.newBackend(config.StorageDestination.Type, config.StorageDestination.Config)
	if err != nil {
		return errwrap.Wrapf("error mounting 'storage_destination': {{err}}", err)
	}

	migrationStatus, err := CheckStorageMigration(from)
	if err != nil {
		return errwrap.Wrapf("error checking migration status: {{err}}", err)
	}

	if migrationStatus != nil {
		return fmt.Errorf("Storage migration in progress (started: %s).", migrationStatus.Start.Format(time.RFC3339))
	}

	if err := SetStorageMigration(from, true); err != nil {
		return errwrap.Wrapf("error setting migration lock: {{err}}", err)
	}

	defer SetStorageMigration(from, false)

	ctx, cancelFunc := context.WithCancel(context.Background())

	doneCh := make(chan error)
	go func() {
		doneCh <- c.migrateAll(ctx, from, to)
	}()

	select {
	case err := <-doneCh:
		cancelFunc()
		return err
	case <-c.ShutdownCh:
		c.UI.Output("==> Migration shutdown triggered\n")
		cancelFunc()
		<-doneCh
		return errAbort
	}
}

// migrateAll copies all keys in lexicographic order.
func (c *OperatorMigrateCommand) migrateAll(ctx context.Context, from physical.Backend, to physical.Backend) error {
	return dfsScan(ctx, from, func(ctx context.Context, path string) error {
		if path < c.flagStart || path == storageMigrationLock || path == vault.CoreLockPath {
			return nil
		}

		entry, err := from.Get(ctx, path)

		if err != nil {
			return errwrap.Wrapf("error reading entry: {{err}}", err)
		}

		if entry == nil {
			return nil
		}

		if err := to.Put(ctx, entry); err != nil {
			return errwrap.Wrapf("error writing entry: {{err}}", err)
		}
		c.logger.Info("copied key", "path", path)
		return nil
	})
}

func (c *OperatorMigrateCommand) newBackend(kind string, conf map[string]string) (physical.Backend, error) {
	factory, ok := c.PhysicalBackends[kind]
	if !ok {
		return nil, fmt.Errorf("no Vault storage backend named: %+q", kind)
	}

	return factory(conf, c.logger)
}

// loadMigratorConfig loads the configuration at the given path
func (c *OperatorMigrateCommand) loadMigratorConfig(path string) (*migratorConfig, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, fmt.Errorf("location is a directory, not a file")
	}

	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	obj, err := hcl.ParseBytes(d)
	if err != nil {
		return nil, err
	}

	var result migratorConfig
	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	// Look for storage_* stanzas
	for _, stanza := range []string{"storage_source", "storage_destination"} {
		o := list.Filter(stanza)
		if len(o.Items) != 1 {
			return nil, fmt.Errorf("exactly one '%s' block is required", stanza)
		}

		if err := parseStorage(&result, o, stanza); err != nil {
			return nil, errwrap.Wrapf("error parsing '%s': {{err}}", err)
		}
	}
	return &result, nil
}

// parseStorage reuses the existing storage parsing that's part of the main Vault
// config processing, but only keeps the storage result.
func parseStorage(result *migratorConfig, list *ast.ObjectList, name string) error {
	tmpConfig := new(server.Config)

	if err := server.ParseStorage(tmpConfig, list, name); err != nil {
		return err
	}

	switch name {
	case "storage_source":
		result.StorageSource = tmpConfig.Storage
	case "storage_destination":
		result.StorageDestination = tmpConfig.Storage
	default:
		return fmt.Errorf("unknown storage name: %s", name)
	}

	return nil
}

// dfsScan will invoke cb with every key from source.
// Keys will be traversed in lexicographic, depth-first order.
func dfsScan(ctx context.Context, source physical.Backend, cb func(ctx context.Context, path string) error) error {
	dfs := []string{""}

	for l := len(dfs); l > 0; l = len(dfs) {
		key := dfs[len(dfs)-1]
		if key == "" || strings.HasSuffix(key, "/") {
			children, err := source.List(ctx, key)
			if err != nil {
				return errwrap.Wrapf("failed to scan for children: {{err}}", err)
			}
			sort.Strings(children)

			// remove List-triggering key and add children in reverse order
			dfs = dfs[:len(dfs)-1]
			for i := len(children) - 1; i >= 0; i-- {
				if children[i] != "" {
					dfs = append(dfs, key+children[i])
				}
			}
		} else {
			err := cb(ctx, key)
			if err != nil {
				return err
			}

			dfs = dfs[:len(dfs)-1]
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
	return nil
}
