package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type V3Build struct {
	State     string          `json:"state,omitempty"`
	Error     string          `json:"error,omitempty"`
	Lifecycle V3Lifecycle     `json:"lifecycle,omitempty"`
	Package   V3Relationship  `json:"package,omitempty"`
	Droplet   V3Relationship  `json:"droplet,omitempty"`
	GUID      string          `json:"guid,omitempty"`
	CreatedAt string          `json:"created_at,omitempty"`
	UpdatedAt string          `json:"updated_at,omitempty"`
	CreatedBy V3CreatedBy     `json:"created_by,omitempty"`
	Links     map[string]Link `json:"links,omitempty"`
	Metadata  V3Metadata      `json:"metadata,omitempty"`
}

type V3CreatedBy struct {
	GUID  string `json:"guid,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func (c *Client) GetV3BuildByGUID(buildGUID string) (*V3Build, error) {
	resp, err := c.DoRequest(c.NewRequest("GET", "/v3/builds/"+buildGUID))
	if err != nil {
		return nil, errors.Wrap(err, "Error getting V3 build")
	}
	defer resp.Body.Close()

	var build V3Build
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, errors.Wrap(err, "Error reading V3 build JSON")
	}

	return &build, nil
}

func (c *Client) CreateV3Build(packageGUID string, lifecycle *V3Lifecycle, metadata *V3Metadata) (*V3Build, error) {
	req := c.NewRequest("POST", "/v3/builds")
	params := map[string]interface{}{
		"package": map[string]interface{}{
			"guid": packageGUID,
		},
	}
	if lifecycle != nil {
		params["lifecycle"] = lifecycle
	}
	if metadata != nil {
		params["metadata"] = metadata
	}
	req.obj = params

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 build")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating v3 build, response code: %d", resp.StatusCode)
	}

	var build V3Build
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, errors.Wrap(err, "Error reading V3 Build JSON")
	}

	return &build, nil
}
