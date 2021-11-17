package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type MetadataHolder struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Annotations map[string]interface{} `json:"annotations"`
	Labels      map[string]interface{} `json:"labels"`
}

func (m *Metadata) AddAnnotation(key string, value string) {
	if m.Annotations == nil {
		m.Annotations = make(map[string]interface{})
	}
	m.Annotations[key] = value
}

func (m *Metadata) RemoveAnnotation(key string) {
	if m.Annotations == nil {
		m.Annotations = make(map[string]interface{})
	}
	m.Annotations[key] = nil
}

func (m *Metadata) AddLabel(prefix, key string, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]interface{})
	}
	if len(prefix) > 0 {
		m.Labels[fmt.Sprintf("%s/%s", prefix, key)] = value
	} else {
		m.Labels[key] = value
	}
}

func (m *Metadata) RemoveLabel(prefix, key string) {
	if m.Labels == nil {
		m.Labels = make(map[string]interface{})
	}
	if len(prefix) > 0 {
		m.Labels[fmt.Sprintf("%s/%s", prefix, key)] = nil
	} else {
		m.Labels[key] = nil
	}
}

func (m *Metadata) Clear() *Metadata {
	metadata := &Metadata{}
	for key := range m.Annotations {
		if strings.Contains(key, "/") {
			metadata.RemoveAnnotation(strings.Split(key, "/")[1])
		}
		metadata.RemoveAnnotation(key)
	}
	for key := range m.Labels {
		metadata.RemoveLabel("", key)
	}
	return metadata
}

func (c *Client) UpdateOrgMetadata(orgGUID string, metadata Metadata) error {
	holder := MetadataHolder{}
	holder.Metadata = metadata
	requestURL := fmt.Sprintf("/v3/organizations/%s", orgGUID)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(holder)
	if err != nil {
		return err
	}
	r := c.NewRequestWithBody("PATCH", requestURL, buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "Error updating metadata for org %s, response code: %d", orgGUID, resp.StatusCode)
	}
	return nil
}

func (c *Client) UpdateSpaceMetadata(spaceGUID string, metadata Metadata) error {
	holder := MetadataHolder{}
	holder.Metadata = metadata
	requestURL := fmt.Sprintf("/v3/spaces/%s", spaceGUID)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(holder)
	if err != nil {
		return err
	}
	r := c.NewRequestWithBody("PATCH", requestURL, buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "Error updating metadata for space %s, response code: %d", spaceGUID, resp.StatusCode)
	}
	return nil
}

func (c *Client) OrgMetadata(orgGUID string) (*Metadata, error) {
	requestURL := fmt.Sprintf("/v3/organizations/%s", orgGUID)
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return &Metadata{}, errors.Wrap(err, "Error requesting space info")
	}
	defer resp.Body.Close()
	return c.handleMetadataResp(resp)
}

func (c *Client) SpaceMetadata(spaceGUID string) (*Metadata, error) {
	requestURL := fmt.Sprintf("/v3/spaces/%s", spaceGUID)
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return &Metadata{}, errors.Wrap(err, "Error requesting space info")
	}
	defer resp.Body.Close()
	return c.handleMetadataResp(resp)
}

func (c *Client) RemoveOrgMetadata(orgGUID string) error {
	metadata, err := c.OrgMetadata(orgGUID)
	if err != nil {
		return err
	}
	return c.UpdateOrgMetadata(orgGUID, *metadata.Clear())
}

func (c *Client) RemoveSpaceMetadata(spaceGUID string) error {
	metadata, err := c.SpaceMetadata(spaceGUID)
	if err != nil {
		return err
	}
	return c.UpdateSpaceMetadata(spaceGUID, *metadata.Clear())
}

func (c *Client) handleMetadataResp(resp *http.Response) (*Metadata, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return &Metadata{}, err
	}
	var metadataResource MetadataHolder
	err = json.Unmarshal(body, &metadataResource)
	if err != nil {
		return &Metadata{}, err
	}
	return &metadataResource.Metadata, nil
}
