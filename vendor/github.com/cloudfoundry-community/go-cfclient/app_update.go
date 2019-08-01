package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UpdateResponse struct {
	Metadata Meta                 `json:"metadata"`
	Entity   UpdateResponseEntity `json:"entity"`
}
type AppUpdateResource struct {
	Name                     string                 `json:"name,omitempty"`
	Memory                   int                    `json:"memory,omitempty"`
	Instances                int                    `json:"instances,omitempty"`
	DiskQuota                int                    `json:"disk_quota,omitempty"`
	SpaceGuid                string                 `json:"space_guid,omitempty"`
	StackGuid                string                 `json:"stack_guid,omitempty"`
	State                    AppState               `json:"state,omitempty"`
	Command                  string                 `json:"command,omitempty"`
	Buildpack                string                 `json:"buildpack,omitempty"`
	HealthCheckHttpEndpoint  string                 `json:"health_check_http_endpoint,omitempty"`
	HealthCheckType          string                 `json:"health_check_type,omitempty"`
	HealthCheckTimeout       int                    `json:"health_check_timeout,omitempty"`
	Diego                    bool                   `json:"diego,omitempty"`
	EnableSSH                bool                   `json:"enable_ssh,omitempty"`
	DockerImage              string                 `json:"docker_image,omitempty"`
	DockerCredentials        map[string]interface{} `json:"docker_credentials_json,omitempty"`
	Environment              map[string]interface{} `json:"environment_json,omitempty"`
	StagingFailedReason      string                 `json:"staging_failed_reason,omitempty"`
	StagingFailedDescription string                 `json:"staging_failed_description,omitempty"`
	Ports                    []int                  `json:"ports,omitempty"`
}

type UpdateResponseEntity struct {
	Name                     string                 `json:"name"`
	Production               bool                   `json:"production"`
	SpaceGuid                string                 `json:"space_guid"`
	StackGuid                string                 `json:"stack_guid"`
	Buildpack                string                 `json:"buildpack"`
	DetectedBuildpack        string                 `json:"detected_buildpack"`
	DetectedBuildpackGuid    string                 `json:"detected_buildpack_guid"`
	Environment              map[string]interface{} `json:"environment_json"`
	Memory                   int                    `json:"memory"`
	Instances                int                    `json:"instances"`
	DiskQuota                int                    `json:"disk_quota"`
	State                    string                 `json:"state"`
	Version                  string                 `json:"version"`
	Command                  string                 `json:"command"`
	Console                  bool                   `json:"console"`
	Debug                    string                 `json:"debug"`
	StagingTaskId            string                 `json:"staging_task_id"`
	PackageState             string                 `json:"package_state"`
	HealthCheckHttpEndpoint  string                 `json:"health_check_http_endpoint"`
	HealthCheckType          string                 `json:"health_check_type"`
	HealthCheckTimeout       int                    `json:"health_check_timeout"`
	StagingFailedReason      string                 `json:"staging_failed_reason"`
	StagingFailedDescription string                 `json:"staging_failed_description"`
	Diego                    bool                   `json:"diego,omitempty"`
	DockerImage              string                 `json:"docker_image"`
	DockerCredentials        struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"docker_credentials"`
	PackageUpdatedAt     string `json:"package_updated_at"`
	DetectedStartCommand string `json:"detected_start_command"`
	EnableSSH            bool   `json:"enable_ssh"`
	Ports                []int  `json:"ports"`
	SpaceURL             string `json:"space_url"`
	StackURL             string `json:"stack_url"`
	RoutesURL            string `json:"routes_url"`
	EventsURL            string `json:"events_url"`
	ServiceBindingsUrl   string `json:"service_bindings_url"`
	RouteMappingsUrl     string `json:"route_mappings_url"`
}

func (c *Client) UpdateApp(guid string, aur AppUpdateResource) (UpdateResponse, error) {
	var updateResponse UpdateResponse

	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(aur)
	if err != nil {
		return UpdateResponse{}, err
	}
	req := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/apps/%s", guid), buf)
	resp, err := c.DoRequest(req)
	if err != nil {
		return UpdateResponse{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return UpdateResponse{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return UpdateResponse{}, err
	}
	err = json.Unmarshal(body, &updateResponse)
	if err != nil {
		return UpdateResponse{}, err
	}
	return updateResponse, nil
}
