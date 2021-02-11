package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type CreateV3DeploymentOptionalParameters struct {
	Droplet  *V3Relationship       `json:"droplet,omitempty"`
	Revision *V3DeploymentRevision `json:"revision,omitempty"`
	Strategy *string               `json:"strategy,omitempty"`
	Metadata *V3Metadata           `json:"metadata,omitempty"`
}

type createV3DeploymentRequest struct {
	*CreateV3DeploymentOptionalParameters `json:",inline"`
	Relationships                         struct {
		App V3ToOneRelationship `json:"app"`
	} `json:"relationships"`
}

type V3DeploymentRevision struct {
	GUID    string `json:"guid"`
	Version int    `json:"version"`
}

type V3ProcessReference struct {
	GUID string `json:"guid"`
	Type string `type:"type"`
}

type V3DeploymentStatus struct {
	Value   string            `json:"value"`
	Reason  string            `json:"reason"`
	Details map[string]string `json:"details"`
}

type V3Deployment struct {
	GUID            string                         `json:"guid"`
	State           string                         `json:"state"`
	Status          V3DeploymentStatus             `json:"status"`
	Strategy        string                         `json:"strategy"`
	Droplet         V3Relationship                 `json:"droplet"`
	PreviousDroplet V3Relationship                 `json:"previous_droplet"`
	NewProcesses    []V3ProcessReference           `json:"new_processes"`
	Revision        V3DeploymentRevision           `json:"revision"`
	CreatedAt       string                         `json:"created_at,omitempty"`
	UpdatedAt       string                         `json:"updated_at,omitempty"`
	Links           map[string]Link                `json:"links,omitempty"`
	Metadata        V3Metadata                     `json:"metadata,omitempty"`
	Relationships   map[string]V3ToOneRelationship `json:"relationships,omitempty"`
}

func (c *Client) GetV3Deployment(deploymentGUID string) (*V3Deployment, error) {
	req := c.NewRequest("GET", "/v3/deployments/"+deploymentGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting deployment")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting deployment with GUID [%s], response code: %d", deploymentGUID, resp.StatusCode)
	}

	var r V3Deployment
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading deployment response JSON")
	}

	return &r, nil
}

func (c *Client) CreateV3Deployment(appGUID string, optionalParams *CreateV3DeploymentOptionalParameters) (*V3Deployment, error) {
	// validate the params
	if optionalParams != nil {
		if optionalParams.Droplet != nil && optionalParams.Revision != nil {
			return nil, errors.New("droplet and revision cannot both be set")
		}
	}

	requestBody := createV3DeploymentRequest{}
	requestBody.CreateV3DeploymentOptionalParameters = optionalParams

	requestBody.Relationships = struct {
		App V3ToOneRelationship "json:\"app\""
	}{
		App: V3ToOneRelationship{
			Data: V3Relationship{
				GUID: appGUID,
			},
		},
	}

	req := c.NewRequest("POST", "/v3/deployments")
	req.obj = requestBody

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating deployment")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating deployment for app GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r V3Deployment
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading deployment response JSON")
	}

	return &r, nil
}

func (c *Client) CancelV3Deployment(deploymentGUID string) error {
	req := c.NewRequest("POST", "/v3/deployments/"+deploymentGUID+"/actions/cancel")
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error canceling deployment")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error canceling deployment [%s], response code: %d", deploymentGUID, resp.StatusCode)
	}

	return nil
}
