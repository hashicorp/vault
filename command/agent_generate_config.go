// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*AgentGenerateConfigCommand)(nil)
	_ cli.CommandAutocomplete = (*AgentGenerateConfigCommand)(nil)
)

type AgentGenerateConfigCommand struct {
	*BaseCommand

	flagType    string
	flagSecrets []string
	flagExec    string
}

func (c *AgentGenerateConfigCommand) Synopsis() string {
	return "Generate a Vault Agent configuration file."
}

func (c *AgentGenerateConfigCommand) Help() string {
	helpText := `
Usage: vault agent generate-config [options] [args]

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AgentGenerateConfigCommand) Flags() *FlagSets {
	set := NewFlagSets(c.UI)

	// Common Options
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "type",
		Target:  &c.flagType,
		Default: "env-template",
		Usage:   "Type of configuration file to generate; currently, only 'env-template' is supported.",
		Completion: complete.PredictSet(
			"env-template",
		),
	})

	f.StringSliceVar(&StringSliceVar{
		Name:       "secret",
		Target:     &c.flagSecrets,
		Usage:      "Path to a kv-v1 or kv-v2 secret (e.g. secret/data/foo, kv-v2/prefix/*); multiple secrets and tail '*' wildcards are allowed.",
		Completion: c.PredictVaultFolders(),
	})

	f.StringVar(&StringVar{
		Name:    "exec",
		Target:  &c.flagExec,
		Default: "env",
		Usage:   "The command to execute in env-template mode.",
	})

	return set
}

func (c *AgentGenerateConfigCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AgentGenerateConfigCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AgentGenerateConfigCommand) Run(args []string) int {
	ctx := context.Background()

	flags := c.Flags()

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = flags.Args()

	if len(args) > 1 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected at most 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	templates, err := fetchTemplates(ctx, client, c.flagSecrets)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error generating templates: %v", err))
		return 2
	}

	var execCommand []string
	if c.flagExec != "" {
		execCommand = strings.Split(c.flagExec, " ")
	} else {
		execCommand = []string{"env"}
	}

	tokenPath, err := homedir.Expand("~/.vault-token")
	if err != nil {
		c.UI.Error(fmt.Sprintf("Could not expand home directory: %v", err))
		return 2
	}

	config := generatedConfig{
		Vault: &generatedConfigVault{
			Address: client.Address(),
		},
		AutoAuth: &generatedConfigAutoAuth{
			Method: &generatedConfigAutoAuthMethod{
				Type: "token_file",
				Config: generatedConfigAutoAuthMethodConfig{
					TokenFilePath: tokenPath,
				},
			},
		},
		EnvTemplates: templates,
		Exec: &generatedConfigExec{
			Command:                execCommand,
			RestartOnSecretChanges: true,
			RestartKillSignal:      "SIGTERM",
		},
	}

	var configPath string
	if len(args) == 1 {
		configPath = args[0]
	} else {
		configPath = "agent.hcl"
	}

	contents := hclwrite.NewEmptyFile()

	gohcl.EncodeIntoBody(&config, contents.Body())

	f, err := os.Create(configPath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Could not create configuration file %q: %v", configPath, err))
		return 1
	}
	defer func() {
		if err := f.Close(); err != nil {
			c.UI.Error(fmt.Sprintf("Could not close configuration file %q: %v", configPath, err))
		}
	}()

	if _, err := contents.WriteTo(f); err != nil {
		c.UI.Error(fmt.Sprintf("Could not write to configuration file %q: %v", configPath, err))
		return 1
	}

	c.UI.Info(fmt.Sprintf("Successfully generated %q configuration file!", configPath))

	return 0
}

func fetchTemplates(ctx context.Context, client *api.Client, secretPaths []string) ([]generatedConfigEnvTemplate, error) {
	var templates []generatedConfigEnvTemplate

	for _, path := range secretPaths {
		path = sanitizePath(path)

		mountPath, v2, err := isKVv2(path, client)
		if err != nil {
			return nil, fmt.Errorf("could not validate secret path %q: %w", path, err)
		}

		switch {
		// this path contains a tail wildcard, attempt to walk the tree
		case strings.HasSuffix(path, "/*"):
			t, err := fetchTemplatesFromTree(ctx, client, path[:len(path)-2], mountPath, v2)
			if err != nil {
				return nil, fmt.Errorf("could not traverse sercet at %q: %w", path, err)
			}
			templates = append(templates, t...)

		// don't allow any other wildcards
		case strings.Contains(path, "*"):
			return nil, fmt.Errorf("the path %q cannot contain '*' wildcard characters except as the last element of the path", path)

		// regular secret path
		default:
			t, err := fetchTemplatesFromSecret(ctx, client, path, mountPath, v2)
			if err != nil {
				return nil, fmt.Errorf("could not read secret at %q: %v", path, err)
			}
			templates = append(templates, t...)
		}
	}

	return templates, nil
}

func fetchTemplatesFromTree(ctx context.Context, client *api.Client, path, mountPath string, v2 bool) ([]generatedConfigEnvTemplate, error) {
	var templates []generatedConfigEnvTemplate

	if v2 {
		path = addPrefixToKVPath(path, mountPath, "metadata", true)
	}

	err := walkSecretsTree(ctx, client, path, func(child string, directory bool) error {
		if directory {
			return nil
		}

		t, err := fetchTemplatesFromSecret(ctx, client, child, mountPath, v2)
		if err != nil {
			return err
		}
		templates = append(templates, t...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return templates, nil
}

func fetchTemplatesFromSecret(ctx context.Context, client *api.Client, path, mountPath string, v2 bool) ([]generatedConfigEnvTemplate, error) {
	var templates []generatedConfigEnvTemplate

	if v2 {
		path = addPrefixToKVPath(path, mountPath, "data", true)
	}

	resp, err := client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("error querying: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("secret not found")
	}

	var data map[string]interface{}
	if v2 {
		internal, ok := resp.Data["data"]
		if !ok {
			return nil, fmt.Errorf("secret.Data not found")
		}
		data = internal.(map[string]interface{})
	} else {
		data = resp.Data
	}

	fields := make([]string, 0, len(data))

	for field := range data {
		fields = append(fields, field)
	}

	// sort for a deterministic output
	sort.Strings(fields)

	var dataContents string
	if v2 {
		dataContents = ".Data.data"
	} else {
		dataContents = ".Data"
	}

	for _, field := range fields {
		templates = append(templates, generatedConfigEnvTemplate{
			Name:              constructDefaultEnvironmentKey(path, field),
			Contents:          fmt.Sprintf(`{{ with secret "%s" }}{{ %s.%s }}{{ end }}`, path, dataContents, field),
			ErrorOnMissingKey: true,
		})
	}

	return templates, nil
}

func constructDefaultEnvironmentKey(path string, field string) string {
	pathParts := strings.Split(path, "/")
	pathPartsLast := pathParts[len(pathParts)-1]

	nonWordRegex := regexp.MustCompile(`[^\w]+`) // match a sequence of non-word characters

	p1 := nonWordRegex.Split(pathPartsLast, -1)
	p2 := nonWordRegex.Split(field, -1)

	keyParts := append(p1, p2...)

	return strings.ToUpper(strings.Join(keyParts, "_"))
}

// Below, we are redefining a subset of the configuration-related structures
// defined under command/agent/config. Using these structures we can tailor the
// output of the generated config, while using the original structures would
// have produced an HCL document with many empty fields. The structures below
// should not be used for anything other than config generation.
type generatedConfig struct {
	AutoAuth     *generatedConfigAutoAuth     `hcl:"auto_auth,block"`
	Vault        *generatedConfigVault        `hcl:"vault,block"`
	EnvTemplates []generatedConfigEnvTemplate `hcl:"env_template,block"`
	Exec         *generatedConfigExec         `hcl:"exec,block"`
}

type generatedConfigExec struct {
	Command                []string `hcl:"command"`
	RestartOnSecretChanges bool     `hcp:"restart_on_secret_changes"`
	RestartKillSignal      string   `hcp:"restart_kill_signal"`
}

type generatedConfigEnvTemplate struct {
	Name              string `hcl:"name,label"`
	Contents          string `hcl:"contents,attr"`
	ErrorOnMissingKey bool   `hcl:"error_on_missing_key,optional"`
}

type generatedConfigVault struct {
	Address string `hcl:"address"`
}

type generatedConfigAutoAuth struct {
	Method *generatedConfigAutoAuthMethod `hcl:"method,block"`
}

type generatedConfigAutoAuthMethod struct {
	Type   string                              `hcl:"type"`
	Config generatedConfigAutoAuthMethodConfig `hcl:"config,block"`
}

type generatedConfigAutoAuthMethodConfig struct {
	TokenFilePath string `hcl:"token_file_path"`
}
