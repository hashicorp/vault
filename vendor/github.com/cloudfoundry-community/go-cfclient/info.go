package cfclient

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Info is metadata about a Cloud Foundry deployment
type Info struct {
	Name                     string `json:"name"`
	Build                    string `json:"build"`
	Support                  string `json:"support"`
	Version                  int    `json:"version"`
	Description              string `json:"description"`
	AuthorizationEndpoint    string `json:"authorization_endpoint"`
	TokenEndpoint            string `json:"token_endpoint"`
	MinCLIVersion            string `json:"min_cli_version"`
	MinRecommendedCLIVersion string `json:"min_recommended_cli_version"`
	APIVersion               string `json:"api_version"`
	AppSSHEndpoint           string `json:"app_ssh_endpoint"`
	AppSSHHostKeyFingerprint string `json:"app_ssh_host_key_fingerprint"`
	AppSSHOauthClient        string `json:"app_ssh_oauth_client"`
	DopplerLoggingEndpoint   string `json:"doppler_logging_endpoint"`
	RoutingEndpoint          string `json:"routing_endpoint"`
	User                     string `json:"user,omitempty"`
}

// GetInfo retrieves Info from the Cloud Controller API
func (c *Client) GetInfo() (*Info, error) {
	r := c.NewRequest("GET", "/v2/info")
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting info")
	}
	defer resp.Body.Close()
	var i Info
	err = json.NewDecoder(resp.Body).Decode(&i)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling info")
	}
	return &i, nil
}
