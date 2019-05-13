package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*TokenCreateCommand)(nil)
var _ cli.CommandAutocomplete = (*TokenCreateCommand)(nil)

type TokenCreateCommand struct {
	*BaseCommand

	flagID              string
	flagDisplayName     string
	flagTTL             time.Duration
	flagExplicitMaxTTL  time.Duration
	flagPeriod          time.Duration
	flagRenewable       bool
	flagOrphan          bool
	flagNoDefaultPolicy bool
	flagUseLimit        int
	flagRole            string
	flagType            string
	flagMetadata        map[string]string
	flagPolicies        []string
}

func (c *TokenCreateCommand) Synopsis() string {
	return "Create a new token"
}

func (c *TokenCreateCommand) Help() string {
	helpText := `
Usage: vault token create [options]

  Creates a new token that can be used for authentication. This token will be
  created as a child of the currently authenticated token. The generated token
  will inherit all policies and permissions of the currently authenticated
  token unless you explicitly define a subset list policies to assign to the
  token.

  A ttl can also be associated with the token. If a ttl is not associated
  with the token, then it cannot be renewed. If a ttl is associated with
  the token, it will expire after that amount of time unless it is renewed.

  Metadata associated with the token (specified with "-metadata") is written
  to the audit log when the token is used.

  If a role is specified, the role may override parameters specified here.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TokenCreateCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "id",
		Target:     &c.flagID,
		Completion: complete.PredictAnything,
		Usage: "Value for the token. By default, this is an auto-generated 36 " +
			"character UUID. Specifying this value requires sudo permissions.",
	})

	f.StringVar(&StringVar{
		Name:       "display-name",
		Target:     &c.flagDisplayName,
		Completion: complete.PredictAnything,
		Usage: "Name to associate with this token. This is a non-sensitive value " +
			"that can be used to help identify created secrets (e.g. prefixes).",
	})

	f.DurationVar(&DurationVar{
		Name:       "ttl",
		Target:     &c.flagTTL,
		Completion: complete.PredictAnything,
		Usage: "Initial TTL to associate with the token. Token renewals may be " +
			"able to extend beyond this value, depending on the configured maximum " +
			"TTLs. This is specified as a numeric string with suffix like \"30s\" " +
			"or \"5m\".",
	})

	f.DurationVar(&DurationVar{
		Name:       "explicit-max-ttl",
		Target:     &c.flagExplicitMaxTTL,
		Completion: complete.PredictAnything,
		Usage: "Explicit maximum lifetime for the token. Unlike normal TTLs, the " +
			"maximum TTL is a hard limit and cannot be exceeded. This is specified " +
			"as a numeric string with suffix like \"30s\" or \"5m\".",
	})

	f.DurationVar(&DurationVar{
		Name:       "period",
		Target:     &c.flagPeriod,
		Completion: complete.PredictAnything,
		Usage: "If specified, every renewal will use the given period. Periodic " +
			"tokens do not expire (unless -explicit-max-ttl is also provided). " +
			"Setting this value requires sudo permissions. This is specified as a " +
			"numeric string with suffix like \"30s\" or \"5m\".",
	})

	f.BoolVar(&BoolVar{
		Name:    "renewable",
		Target:  &c.flagRenewable,
		Default: true,
		Usage:   "Allow the token to be renewed up to it's maximum TTL.",
	})

	f.BoolVar(&BoolVar{
		Name:    "orphan",
		Target:  &c.flagOrphan,
		Default: false,
		Usage: "Create the token with no parent. This prevents the token from " +
			"being revoked when the token which created it expires. Setting this " +
			"value requires sudo permissions.",
	})

	f.BoolVar(&BoolVar{
		Name:    "no-default-policy",
		Target:  &c.flagNoDefaultPolicy,
		Default: false,
		Usage: "Detach the \"default\" policy from the policy set for this " +
			"token.",
	})

	f.IntVar(&IntVar{
		Name:    "use-limit",
		Target:  &c.flagUseLimit,
		Default: 0,
		Usage: "Number of times this token can be used. After the last use, the " +
			"token is automatically revoked. By default, tokens can be used an " +
			"unlimited number of times until their expiration.",
	})

	f.StringVar(&StringVar{
		Name:    "role",
		Target:  &c.flagRole,
		Default: "",
		Usage: "Name of the role to create the token against. Specifying -role " +
			"may override other arguments. The locally authenticated Vault token " +
			"must have permission for \"auth/token/create/<role>\".",
	})

	f.StringVar(&StringVar{
		Name:    "type",
		Target:  &c.flagType,
		Default: "service",
		Usage:   `The type of token to create. Can be "service" or "batch".`,
	})

	f.StringMapVar(&StringMapVar{
		Name:       "metadata",
		Target:     &c.flagMetadata,
		Completion: complete.PredictAnything,
		Usage: "Arbitrary key=value metadata to associate with the token. " +
			"This metadata will show in the audit log when the token is used. " +
			"This can be specified multiple times to add multiple pieces of " +
			"metadata.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:       "policy",
		Target:     &c.flagPolicies,
		Completion: c.PredictVaultPolicies(),
		Usage: "Name of a policy to associate with this token. This can be " +
			"specified multiple times to attach multiple policies.",
	})

	return set
}

func (c *TokenCreateCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TokenCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TokenCreateCommand) Run(args []string) int {
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

	if c.flagType == "batch" {
		c.flagRenewable = false
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	tcr := &api.TokenCreateRequest{
		ID:              c.flagID,
		Policies:        c.flagPolicies,
		Metadata:        c.flagMetadata,
		TTL:             c.flagTTL.String(),
		NoParent:        c.flagOrphan,
		NoDefaultPolicy: c.flagNoDefaultPolicy,
		DisplayName:     c.flagDisplayName,
		NumUses:         c.flagUseLimit,
		Renewable:       &c.flagRenewable,
		ExplicitMaxTTL:  c.flagExplicitMaxTTL.String(),
		Period:          c.flagPeriod.String(),
		Type:            c.flagType,
	}

	var secret *api.Secret
	if c.flagRole != "" {
		secret, err = client.Auth().Token().CreateWithRole(tcr, c.flagRole)
	} else {
		secret, err = client.Auth().Token().Create(tcr)
	}
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating token: %s", err))
		return 2
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}
