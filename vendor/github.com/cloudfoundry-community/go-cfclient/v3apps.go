package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type V3App struct {
	Name          string                         `json:"name,omitempty"`
	State         string                         `json:"state,omitempty"`
	Lifecycle     V3Lifecycle                    `json:"lifecycle,omitempty"`
	GUID          string                         `json:"guid,omitempty"`
	CreatedAt     string                         `json:"created_at,omitempty"`
	UpdatedAt     string                         `json:"updated_at,omitempty"`
	Relationships map[string]V3ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link                `json:"links,omitempty"`
	Metadata      V3Metadata                     `json:"metadata,omitempty"`
}

type V3Lifecycle struct {
	Type          string               `json:"type,omitempty"`
	BuildpackData V3BuildpackLifecycle `json:"data,omitempty"`
}

type V3BuildpackLifecycle struct {
	Buildpacks []string `json:"buildpacks,omitempty"`
	Stack      string   `json:"stack,omitempty"`
}

type CreateV3AppRequest struct {
	Name                 string
	SpaceGUID            string
	EnvironmentVariables map[string]string
	Lifecycle            *V3Lifecycle
	Metadata             *V3Metadata
}

type UpdateV3AppRequest struct {
	Name      string       `json:"name"`
	Lifecycle *V3Lifecycle `json:"lifecycle"`
	Metadata  *V3Metadata  `json:"metadata"`
}

func (c *Client) CreateV3App(r CreateV3AppRequest) (*V3App, error) {
	req := c.NewRequest("POST", "/v3/apps")
	params := map[string]interface{}{
		"name": r.Name,
		"relationships": map[string]interface{}{
			"space": V3ToOneRelationship{
				Data: V3Relationship{
					GUID: r.SpaceGUID,
				},
			},
		},
	}
	if len(r.EnvironmentVariables) > 0 {
		params["environment_variables"] = r.EnvironmentVariables
	}
	if r.Lifecycle != nil {
		params["lifecycle"] = r.Lifecycle
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	req.obj = params
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating v3 app %s, response code: %d", r.Name, resp.StatusCode)
	}

	var app V3App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 app JSON")
	}

	return &app, nil
}

func (c *Client) GetV3AppByGUID(guid string) (*V3App, error) {
	req := c.NewRequest("GET", "/v3/apps/"+guid)

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting v3 app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting v3 app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	var app V3App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 app JSON")
	}

	return &app, nil
}

func (c *Client) StartV3App(guid string) (*V3App, error) {
	req := c.NewRequest("POST", "/v3/apps/"+guid+"/actions/start")
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while starting v3 app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error starting v3 app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	var app V3App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 app JSON")
	}

	return &app, nil
}

func (c *Client) DeleteV3App(guid string) error {
	req := c.NewRequest("DELETE", "/v3/apps/"+guid)
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting v3 app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting v3 app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateV3App(appGUID string, r UpdateV3AppRequest) (*V3App, error) {
	req := c.NewRequest("PATCH", "/v3/apps/"+appGUID)
	params := make(map[string]interface{})
	if r.Name != "" {
		params["name"] = r.Name
	}
	if r.Lifecycle != nil {
		params["lifecycle"] = r.Lifecycle
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}
	if len(params) > 0 {
		req.obj = params
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating v3 app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error updating v3 app %s, response code: %d", appGUID, resp.StatusCode)
	}

	var app V3App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 app JSON")
	}

	return &app, nil
}

type listV3AppsResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []V3App    `json:"resources,omitempty"`
}

func (c *Client) ListV3AppsByQuery(query url.Values) ([]V3App, error) {
	var apps []V3App
	requestURL := "/v3/apps"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 apps")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 apps, response code: %d", resp.StatusCode)
		}

		var data listV3AppsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 apps")
		}

		apps = append(apps, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 apps")
		}
	}

	return apps, nil
}

func extractPathFromURL(requestURL string) (string, error) {
	url, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}
	result := url.Path
	if q := url.Query().Encode(); q != "" {
		result = result + "?" + q
	}
	return result, nil
}

type V3AppEnvironment struct {
	EnvVars       map[string]string          `json:"environment_variables,omitempty"`
	StagingEnv    map[string]string          `json:"staging_env_json,omitempty"`
	RunningEnv    map[string]string          `json:"running_env_json,omitempty"`
	SystemEnvVars map[string]json.RawMessage `json:"system_env_json,omitempty"`      // VCAP_SERVICES
	AppEnvVars    map[string]json.RawMessage `json:"application_env_json,omitempty"` // VCAP_APPLICATION
}

func (c *Client) GetV3AppEnvironment(appGUID string) (V3AppEnvironment, error) {
	var result V3AppEnvironment

	resp, err := c.DoRequest(c.NewRequest("GET", "/v3/apps/"+appGUID+"/env"))
	if err != nil {
		return result, errors.Wrapf(err, "Error requesting app env for %s", appGUID)
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, errors.Wrap(err, "Error parsing JSON for app env")
	}

	return result, nil
}

type V3EnvVar struct {
	Var map[string]*string `json:"var"`
}

type v3EnvVarResponse struct {
	V3EnvVar
	Links map[string]Link `json:"links"`
}

func (c *Client) SetV3AppEnvVariables(appGUID string, envRequest V3EnvVar) (V3EnvVar, error) {
	var result v3EnvVarResponse

	req := c.NewRequest("PATCH", "/v3/apps/"+appGUID+"/environment_variables")
	req.obj = envRequest

	resp, err := c.DoRequest(req)
	if err != nil {
		return result.V3EnvVar, errors.Wrapf(err, "Error setting app env variables for %s", appGUID)
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result.V3EnvVar, errors.Wrap(err, "Error parsing JSON for app env")
	}

	return result.V3EnvVar, nil
}
