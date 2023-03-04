package command

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"github.com/mitchellh/cli"
	"github.com/pkg/errors"
	"github.com/posener/complete"
	"golang.org/x/oauth2"

	hcpv "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/client/vault_service"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/models"
)

const (
	hcpVaultClusterDetailFile = "hcpv_cluster.json"
	defaultDirectory          = ".config/hcp"
)

var _ cli.Command = (*CloudSetClusterCommand)(nil)

type CloudSetClusterCommand struct {
	*BaseCommand

	flagOrganizationID string
	flagProjectID      string
	flagClusterName    string
}

func (c *CloudSetClusterCommand) Synopsis() string {
	return "Login to HCP"
}

func (c *CloudSetClusterCommand) Help() string {
	helpText := `
Usage: vault cloud set-cluster --organization-id <org-id> --project-id <project-id> --name <name> [options] [args]
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *CloudSetClusterCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "organization-id",
		Target:  &c.flagOrganizationID,
		Default: "",
		Usage:   "HCP Organization ID that contains the target HCP Vault cluster",
	})

	f.StringVar(&StringVar{
		Name:    "project-id",
		Target:  &c.flagProjectID,
		Default: "",
		Usage:   "HCP Project ID that contains the target HCP Vault cluster",
	})
	f.StringVar(&StringVar{
		Name:    "name",
		Target:  &c.flagClusterName,
		Default: "",
		Usage:   "HCP Vault Cluster name",
	})
	return set
}

func (c *CloudSetClusterCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *CloudSetClusterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *CloudSetClusterCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if c.flagOrganizationID == "" || c.flagProjectID == "" || c.flagClusterName == "" {
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

	// Get the cluster
	getParams := hcpv.NewGetParams().
		WithClusterID(c.flagClusterName).
		WithLocationOrganizationID(c.flagOrganizationID).
		WithLocationProjectID(c.flagProjectID)

	resp, err := hcpvClient.Get(getParams, nil)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to lookup HCP Vault cluster details: %s", err))
		return 1
	}

	if err := storeHCPVCluster(resp.GetPayload().Cluster); err != nil {
		c.UI.Error(fmt.Sprintf("Failed to store HCP Vault cluster details: %s", err))
		return 1
	}

	return 0
}

func storeHCPVCluster(c *models.HashicorpCloudVault20201125Cluster) error {
	bytes, err := c.MarshalBinary()
	if err != nil {
		return fmt.Errorf("error marshalling cluster details: %w", err)
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to retrieve user's home directory path: %v", err)
	}

	clusterFile := filepath.Join(userHome, defaultDirectory, hcpVaultClusterDetailFile)
	err = os.WriteFile(clusterFile, bytes, 0755)
	if err != nil {
		return fmt.Errorf("failed to write credentials to the cache file: %v", err)
	}

	return nil
}

func IsHCPVCluster() bool {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	clusterFile := filepath.Join(userHome, defaultDirectory, hcpVaultClusterDetailFile)
	_, err = os.Stat(clusterFile)
	if err != nil {
		return false
	}

	return true
}

func GetHCPVCluster() (*models.HashicorpCloudVault20201125Cluster, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	clusterFile := filepath.Join(userHome, defaultDirectory, hcpVaultClusterDetailFile)
	data, err := os.ReadFile(clusterFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read HCP Vault cluster detail file: %w", err)
	}

	m := &models.HashicorpCloudVault20201125Cluster{}
	if err := m.UnmarshalBinary(data); err != nil {
		os.Remove(clusterFile)
		return nil, fmt.Errorf("failed to unmarshall HCP Vault cluster detail file: %w", err)
	}

	return m, nil
}

type authRoundTripper struct {
	// Source supplies the token to add to outgoing requests'
	// Authorization headers.
	Source oauth2.TokenSource

	// Base is the base RoundTripper used to make HTTP requests.
	// If nil, http.DefaultTransport is used.
	Base http.RoundTripper
}

func (a *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}

	if a.Source == nil {
		return nil, errors.New("oauth2: Transport's Source is nil")
	}
	token, err := a.Source.Token()
	if err != nil {
		return nil, err
	}

	req2 := cloneRequest(req) // per RoundTripper contract

	cookie := &http.Cookie{
		//Domain:  domain,
		Name:    "hcp_access_token",
		Value:   token.AccessToken,
		Expires: token.Expiry,
	}
	req2.AddCookie(cookie)

	// req.Body is assumed to be closed by the base RoundTripper.
	reqBodyClosed = true
	return a.Base.RoundTrip(req2)
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

func HCPProxyRoundTripper(c *models.HashicorpCloudVault20201125Cluster, base http.RoundTripper) (http.RoundTripper, error) {
	cfg, err := config.NewHCPConfig(config.FromEnv())
	if err != nil {
		return nil, errors.Wrap(err, "failed retrieving HCP config")
	}

	rt := &authRoundTripper{
		Source: cfg,
		Base:   base,
	}

	return rt, nil
}
