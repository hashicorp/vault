package command

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	hcpv "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/client/vault_service"
	hcpvmodels "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/models"
	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
)

var _ cli.Command = (*CloudListClustersCommand)(nil)

type CloudListClustersCommand struct {
	*BaseCommand

	flagOrganizationID string
	flagProjectID      string
}

func (c *CloudListClustersCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if c.flagOrganizationID == "" || c.flagProjectID == "" {
		c.UI.Error("HCP Organization, Project ID and cluster name must be set")
		return 1
	}

	cfg, err := config.NewHCPConfig(config.FromEnv())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating HCP Config: %s", err))
		return 1
	}

	hcpClient, err := httpclient.New(httpclient.Config{HCPConfig: cfg})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating HCP HTTP Client: %s", err))
		return 1
	}

	hcpvClient := hcpv.New(hcpClient, nil)

	listParams := hcpv.NewListParams().
		WithLocationProjectID(c.flagProjectID).
		WithLocationOrganizationID(c.flagOrganizationID)

	// TODO: deal with pagination, right now just list whats returned
	resp, err := hcpvClient.List(listParams, nil)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to list HCP Vault clusters: %s", err))
		return 1
	}
	clusters := resp.GetPayload().Clusters

	// TODO: usually commands use `OutputList`, but the table formatter works for vault specific responses
	format := Format(c.UI)
	if format == "table" {

		isProxyEnabled := func(clusterDNS *hcpvmodels.HashicorpCloudVault20201125ClusterDNSNames) bool {
			if clusterDNS == nil {
				return false
			}
			return clusterDNS.Proxy != ""
		}

		header := strings.Join(
			[]string{"ID", "Version", "State", "Proxy Enabled"},
			hopeDelim,
		)
		rows := []string{header}
		for _, cluster := range clusters {
			row := []string{
				cluster.ID,
				cluster.CurrentVersion,
				string(*cluster.State),
				strconv.FormatBool(isProxyEnabled(cluster.DNSNames)),
			}
			rows = append(rows, strings.Join(row, hopeDelim))
		}
		output := tableOutput(rows, &columnize.Config{
			Delim: hopeDelim,
		})
		c.UI.Output(output)
		return 0
	}

	return OutputList(c.UI, clusters)
}

func (c *CloudListClustersCommand) Help() string {
	helpText := `
Usage: vault cloud list-clusters --organization-id [org-id] --project-id [project-id]
` + c.Flags().Help()

	return strings.TrimSpace(helpText)

}

func (c *CloudListClustersCommand) Synopsis() string {
	return "List HCP Vault clusters"
}

func (c *CloudListClustersCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:   "organization-id",
		Target: &c.flagOrganizationID,
		// https://github.com/hashicorp/hcp-sdk-go#user-profile
		Default: os.Getenv("HCP_ORGANIZATION_ID"),
		Usage:   "HCP Organization ID that contains the target HCP Vault cluster",
	})

	f.StringVar(&StringVar{
		Name:   "project-id",
		Target: &c.flagProjectID,
		// https://github.com/hashicorp/hcp-sdk-go#user-profile
		Default: os.Getenv("HCP_PROJECT_ID"),
		Usage:   "HCP Project ID that contains the target HCP Vault cluster",
	})

	return set
}
