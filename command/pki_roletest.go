package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PKIRoleTestCommand)(nil)
	_ cli.CommandAutocomplete = (*PKIRoleTestCommand)(nil)
)

type PKIRoleTestCommand struct {
	*BaseCommand

	flagQuiet bool
	flagMount string
}

func (c *PKIRoleTestCommand) Synopsis() string {
	return "Test PKI Secrets Engine role issuance restrictions"
}

func (c *PKIRoleTestCommand) Help() string {
	helpText := `
Usage: vault pki role-test [options] ROLE COMMON_NAME [DNS-SANS...] [K=V]

  Reports whether or not issuance will succeed against the PKI role with the
  specified common name and optionally, DNS SAN names. Other parameters to
  /pki/issue/:ROLE can be specified in K=V format (mirroring vault write).

  Check whether the role will issue for localhost:

      $ vault pki role-test -mount=pki-int server-role localhost

  To withhold all output and only use return codes for indicating issuance
  status:

      $ vault pki role-test -quiet server-role foo.example.com

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIRoleTestCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "quiet",
		Target:  &c.flagQuiet,
		Default: false,
		EnvVar:  "",
		Usage: "Suppress CLI output; use return status to indicate success" +
			" (return code 0) or failure (return code 3).",
	})

	f.StringVar(&StringVar{
		Name:    "mount",
		Target:  &c.flagMount,
		Default: "pki",
		EnvVar:  "",
		Usage:   "PKI mount to test issuance under.",
	})

	return set
}

func (c *PKIRoleTestCommand) AutocompleteArgs() complete.Predictor {
	// Return an anything predictor here, similar to `vault write`. We
	// don't know what values are valid for the role and/or common names.
	return complete.PredictAnything
}

func (c *PKIRoleTestCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKIRoleTestCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) < 2 {
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2+, got %d)", len(args)))
		return 1
	}

	mount := sanitizePath(c.flagMount)
	role := sanitizePath(args[0])
	path := mount + "/issue/" + role
	commonName := args[1]

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var dnsSans []string
	remainder := args[2:]

	for _, value := range args[2:] {
		if strings.Contains(value, "=") {
			// Start of K=V data.
			break
		}

		dnsSans = append(dnsSans, value)
		remainder = remainder[1:]
	}

	data, err := parseArgsData(nil, remainder)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	if data == nil {
		data = make(map[string]interface{})
	}

	data["dry_run"] = true
	data["common_name"] = commonName
	data["ttl"] = "1s"
	if len(dnsSans) > 0 {
		data["alt_names"] = strings.Join(dnsSans, ",")
	}

	_, err = client.Logical().Write(path, data)
	if err != nil {
		if !c.flagQuiet {
			c.UI.Error(fmt.Sprintf("Error issuing certificate: %v", err))
		}

		return 3
	}

	if !c.flagQuiet {
		c.UI.Info(fmt.Sprintf("Success! Certificate would be issued."))
	}

	return 0
}
