package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/meta"
)

type CorsCommand struct {
	meta.Meta
}

func (c *CorsCommand) Run(args []string) int {
	var allowedStr string
	var disable, status bool

	flags := c.Meta.FlagSet("cors", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	flags.StringVar(&allowedStr, "allowed-origins", "", "")
	flags.BoolVar(&disable, "disable", false, "")
	flags.BoolVar(&status, "status", false, "")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	corsRequest := &api.CORSRequest{}

	if !status {
		allowedOrigins, err := regexp.Compile(allowedStr)
		if err != nil {
			return 1
		}

		if !disable {
			corsRequest.AllowedOrigins = allowedOrigins.String()
			corsRequest.Enabled = true
		} else {
			corsRequest.Enabled = false
		}
	}

	return c.runCORS(status, corsRequest)
}

func (c *CorsCommand) runCORS(status bool, corsRequest *api.CORSRequest) int {
	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 1
	}

	// Get the current CORS configuration.
	if status {
		resp, err := client.Sys().CORSStatus()
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error getting CORS configuration: %s", err))
			return 1
		}
		c.Ui.Output(
			fmt.Sprintf("Enabled: %t\nAllowed Origins: %s",
				resp.Enabled,
				resp.AllowedOrigins,
			))
		return 0
	}

	// Disable (i.e. clear) the CORS configuration.
	if corsRequest.Enabled == false {
		_, err := client.Sys().DisableCORS()
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error disabling CORS: %s", err))
			return 1
		}
		return 0
	}

	// Update the CORS configuration.
	_, err = client.Sys().ConfigureCORS(corsRequest)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error configuring CORS: %s", err))
		return 1
	}

	return 0
}

func (c *CorsCommand) Help() string {
	helpText := `
Usage: vault cors [options]

  Configures the HTTP server to return CORS headers.

  This command connects to a Vault server and can enable CORS, disable CORS, or change
  the regular expression for origins that are allowed to make cross-origin requests.

General Options:
` + meta.GeneralOptionsUsage() + `
Cors Options:

  -status                   Returns the current CORS configuration.

  -allowed-origins=""       A regular expression that describes the origins
                            that should be allowed to make cross-origin
                            requests and be served CORS headers. A return code
                            of 0 means the regular expressions is valid, and
                            Vault will now serve CORS headers to clients from
                            matching origins; a return code of 1 means an error
                            was encountered.

  -disable                  Stop serving CORS headers for all origins.
`
	return strings.TrimSpace(helpText)
}

func (c *CorsCommand) Synopsis() string {
	return "Configure CORS settings"
}
