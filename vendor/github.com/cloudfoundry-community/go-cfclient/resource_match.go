package cfclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// The Resource match Api retruns the response in the following data structure
type Resource struct {
	Sha1 string `json:"sha1"`
	Size int    `json:"size"`
}

// ResourceMatch matches given resource list of SHA / file size pairs against
// the Cloud Controller cache, and reports the subset which describes already
// existing files
func (c *Client) ResourceMatch(resources []Resource) ([]Resource, error) {

	var resourcesList []Resource
	emptyResource := make([]Resource, 0)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(resources)
	if err != nil {
		return emptyResource, errors.Wrapf(err, "Error reading Resource List")
	}
	r := c.NewRequestWithBody("PUT", "/v2/resource_match", buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return emptyResource, errors.Wrapf(err, "Error uploading Resource List")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return emptyResource, errors.Wrapf(err, "Error reading Resources http response body")
	}
	err = json.Unmarshal(resBody, &resourcesList)
	if err != nil {
		return emptyResource, errors.Wrapf(err, "Error reading Resources http response body")
	}
	return resourcesList, nil

}
