package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type V3PackageState string

const (
	AwaitingUpload   V3PackageState = "AWAITING_UPLOAD"
	ProcessingUpload V3PackageState = "PROCESSING_UPLOAD"
	Ready            V3PackageState = "READY"
	Failed           V3PackageState = "FAILED"
	Copying          V3PackageState = "COPYING"
	Expired          V3PackageState = "EXPIRED"
)

type V3Package struct {
	Type      string          `json:"type,omitempty"` // bits or docker
	Data      json.RawMessage `json:"data,omitempty"` // depends on value of Type
	State     V3PackageState  `json:"state,omitempty"`
	GUID      string          `json:"guid,omitempty"`
	CreatedAt string          `json:"created_at,omitempty"`
	UpdatedAt string          `json:"updated_at,omitempty"`
	Links     map[string]Link `json:"links,omitempty"`
	Metadata  V3Metadata      `json:"metadata,omitempty"`
}

func (v *V3Package) BitsData() (V3BitsPackage, error) {
	var bits V3BitsPackage
	if v.Type != "bits" {
		return bits, errors.New("this package is not of type bits")
	}

	if err := json.Unmarshal(v.Data, &bits); err != nil {
		return bits, err
	}

	return bits, nil
}

func (v *V3Package) DockerData() (V3DockerPackage, error) {
	var docker V3DockerPackage
	if v.Type != "docker" {
		return docker, errors.New("this package is not of type docker")
	}

	if err := json.Unmarshal(v.Data, &docker); err != nil {
		return docker, err
	}

	return docker, nil
}

// V3BitsPackage is the data for V3Packages of type bits.
// It provides an upload link to which a zip file should be uploaded.
type V3BitsPackage struct {
	Error    string `json:"error,omitempty"`
	Checksum struct {
		Type  string `json:"type,omitempty"`  // eg. sha256
		Value string `json:"value,omitempty"` // populated after the bits are uploaded
	} `json:"checksum,omitempty"`
}

// V3DockerPackage is the data for V3Packages of type docker.
// It references a docker image from a registry.
type V3DockerPackage struct {
	Image    string `json:"image,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type listV3PackagesResponse struct {
	Pagination Pagination  `json:"pagination,omitempty"`
	Resources  []V3Package `json:"resources,omitempty"`
}

func (c *Client) ListPackagesForAppV3(appGUID string, query url.Values) ([]V3Package, error) {
	var packages []V3Package
	requestURL := "/v3/apps/" + appGUID + "/packages"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		resp, err := c.DoRequest(c.NewRequest("GET", requestURL))
		if err != nil {
			return nil, errors.Wrapf(err, "Error requesting packages for app %s", appGUID)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 app packages, response code: %d", resp.StatusCode)
		}

		var data listV3PackagesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 app packages")
		}

		packages = append(packages, data.Resources...)
		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 packages")
		}
	}
	return packages, nil
}

// CopyPackageV3 makes a copy of a package that is associated with one app
// and associates the copy with a new app.
func (c *Client) CopyPackageV3(packageGUID, appGUID string) (*V3Package, error) {
	req := c.NewRequest("POST", "/v3/packages?source_guid="+packageGUID)
	req.obj = map[string]interface{}{
		"relationships": map[string]interface{}{
			"app": V3ToOneRelationship{
				Data: V3Relationship{
					GUID: appGUID,
				},
			},
		},
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while copying v3 package")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error copying v3 package %s, response code: %d", packageGUID, resp.StatusCode)
	}

	var pkg V3Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 app package")
	}

	return &pkg, nil
}

type v3DockerPackageData struct {
	Image string `json:"image"`
	*DockerCredentials
}

type createV3DockerPackageRequest struct {
	Type          string                         `json:"type"`
	Relationships map[string]V3ToOneRelationship `json:"relationships"`
	Data          v3DockerPackageData            `json:"data"`
}

// CreateV3DockerPackage creates a Docker package
func (c *Client) CreateV3DockerPackage(image string, appGUID string, dockerCredentials *DockerCredentials) (*V3Package, error) {
	req := c.NewRequest("POST", "/v3/packages")
	req.obj = createV3DockerPackageRequest{
		Type: "docker",
		Relationships: map[string]V3ToOneRelationship{
			"app": {Data: V3Relationship{GUID: appGUID}},
		},
		Data: v3DockerPackageData{
			Image:             image,
			DockerCredentials: dockerCredentials,
		},
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while copying v3 package")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating v3 docker package, response code: %d", resp.StatusCode)
	}

	var pkg V3Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 app package")
	}

	return &pkg, nil
}
