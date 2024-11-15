package cfclient

import (
	"encoding/json"

	"github.com/Masterminds/semver"
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

type V3Version struct {
	Links struct {
		CCV3 struct {
			Meta struct {
				Version string `json:"version"`
			} `json:"meta"`
		} `json:"cloud_controller_v3"`
	} `json:"links"`
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

func (c *Client) SupportsMetadataAPI() (bool, error) {
	r := c.NewRequest("GET", "/")
	resp, err := c.DoRequest(r)
	if err != nil {
		return false, errors.Wrap(err, "Error requesting info")
	}
	defer resp.Body.Close()
	var v3 V3Version
	err = json.NewDecoder(resp.Body).Decode(&v3)
	if err != nil {
		return false, errors.Wrap(err, "Error unmarshalling info")
	}

	minimumSupportedVersion := semver.MustParse("3.66.0")
	actualVersion, err := semver.NewVersion(v3.Links.CCV3.Meta.Version)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing semver")
	}
	if !actualVersion.LessThan(minimumSupportedVersion) {
		return true, nil
	}

	return false, nil
}

func (c *Client) SupportsSpaceSupporterRole() (bool, error) {
	r := c.NewRequest("GET", "/")
	resp, err := c.DoRequest(r)
	if err != nil {
		return false, errors.Wrap(err, "Error requesting info")
	}
	defer resp.Body.Close()
	var v3 V3Version
	err = json.NewDecoder(resp.Body).Decode(&v3)
	if err != nil {
		return false, errors.Wrap(err, "Error unmarshalling info")
	}

	minimumSupportedVersion := semver.MustParse("3.102.0")
	actualVersion, err := semver.NewVersion(v3.Links.CCV3.Meta.Version)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing semver")
	}
	if !actualVersion.LessThan(minimumSupportedVersion) {
		return true, nil
	}

	return false, nil
}
