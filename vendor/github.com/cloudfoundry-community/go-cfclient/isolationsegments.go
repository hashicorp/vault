package cfclient

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type IsolationSegment struct {
	GUID      string    `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	c         *Client
}

type IsolationSegementResponse struct {
	GUID      string    `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Links     struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Spaces struct {
			Href string `json:"href"`
		} `json:"spaces"`
		Organizations struct {
			Href string `json:"href"`
		} `json:"organizations"`
	} `json:"links"`
}

type ListIsolationSegmentsResponse struct {
	Pagination Pagination                  `json:"pagination"`
	Resources  []IsolationSegementResponse `json:"resources"`
}

func (c *Client) CreateIsolationSegment(name string) (*IsolationSegment, error) {
	req := c.NewRequest("POST", "/v3/isolation_segments")
	req.obj = map[string]interface{}{
		"name": name,
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating isolation segment")
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating isolation segment %s, response code: %d", name, resp.StatusCode)
	}
	return respBodyToIsolationSegment(resp.Body, c)
}

func respBodyToIsolationSegment(body io.ReadCloser, c *Client) (*IsolationSegment, error) {
	bodyRaw, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	isr := IsolationSegementResponse{}
	err = json.Unmarshal(bodyRaw, &isr)
	if err != nil {
		return nil, err
	}

	return &IsolationSegment{
		GUID:      isr.GUID,
		Name:      isr.Name,
		CreatedAt: isr.CreatedAt,
		UpdatedAt: isr.UpdatedAt,
		c:         c,
	}, nil
}

func (c *Client) GetIsolationSegmentByGUID(guid string) (*IsolationSegment, error) {
	var isr IsolationSegementResponse
	r := c.NewRequest("GET", "/v3/isolation_segments/"+guid)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting isolation segment by GUID")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading isolation segment response body")
	}

	err = json.Unmarshal(resBody, &isr)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling isolation segment response")
	}
	return &IsolationSegment{Name: isr.Name, GUID: isr.GUID, CreatedAt: isr.CreatedAt, UpdatedAt: isr.UpdatedAt, c: c}, nil
}

func (c *Client) ListIsolationSegmentsByQuery(query url.Values) ([]IsolationSegment, error) {
	var iss []IsolationSegment
	requestUrl := "/v3/isolation_segments?" + query.Encode()
	for {
		var isr ListIsolationSegmentsResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting isolation segments")
		}
		defer resp.Body.Close()
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading isolation segment request")
		}

		err = json.Unmarshal(resBody, &isr)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshalling isolation segment")
		}

		for _, is := range isr.Resources {
			iss = append(iss, IsolationSegment{
				Name:      is.Name,
				GUID:      is.GUID,
				CreatedAt: is.CreatedAt,
				UpdatedAt: is.UpdatedAt,
				c:         c,
			})
		}

		var ok bool
		requestUrl, ok = isr.Pagination.Next.(string)
		if !ok || requestUrl == "" {
			break
		}
	}
	return iss, nil
}

func (c *Client) ListIsolationSegments() ([]IsolationSegment, error) {
	return c.ListIsolationSegmentsByQuery(nil)
}

// TODO listOrgsForIsolationSegments
// TODO listSpacesForIsolationSegments
// TODO setDefaultIsolationSegmentForOrg

func (c *Client) DeleteIsolationSegmentByGUID(guid string) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v3/isolation_segments/%s", guid)))
	if err != nil {
		return errors.Wrap(err, "Error during sending DELETE request for isolation segments")
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Error deleting isolation segment %s, response code: %d", guid, resp.StatusCode)
	}
	return nil
}

func (i *IsolationSegment) Delete() error {
	return i.c.DeleteIsolationSegmentByGUID(i.GUID)
}

func (c *Client) AddIsolationSegmentToOrg(isolationSegmentGUID, orgGUID string) error {
	isoSegment := IsolationSegment{GUID: isolationSegmentGUID, c: c}
	return isoSegment.AddOrg(orgGUID)
}

func (c *Client) RemoveIsolationSegmentFromOrg(isolationSegmentGUID, orgGUID string) error {
	isoSegment := IsolationSegment{GUID: isolationSegmentGUID, c: c}
	return isoSegment.RemoveOrg(orgGUID)
}

func (c *Client) AddIsolationSegmentToSpace(isolationSegmentGUID, spaceGUID string) error {
	isoSegment := IsolationSegment{GUID: isolationSegmentGUID, c: c}
	return isoSegment.AddSpace(spaceGUID)
}

func (c *Client) RemoveIsolationSegmentFromSpace(isolationSegmentGUID, spaceGUID string) error {
	isoSegment := IsolationSegment{GUID: isolationSegmentGUID, c: c}
	return isoSegment.RemoveSpace(spaceGUID)
}

func (i *IsolationSegment) AddOrg(orgGuid string) error {
	if i == nil || i.c == nil {
		return errors.New("No communication handle.")
	}
	req := i.c.NewRequest("POST", fmt.Sprintf("/v3/isolation_segments/%s/relationships/organizations", i.GUID))
	type Entry struct {
		GUID string `json:"guid"`
	}
	req.obj = map[string]interface{}{
		"data": []Entry{{GUID: orgGuid}},
	}
	resp, err := i.c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error during adding org to isolation segment")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error adding org %s to isolation segment %s, response code: %d", orgGuid, i.Name, resp.StatusCode)
	}
	return nil
}

func (i *IsolationSegment) RemoveOrg(orgGuid string) error {
	if i == nil || i.c == nil {
		return errors.New("No communication handle.")
	}
	req := i.c.NewRequest("DELETE", fmt.Sprintf("/v3/isolation_segments/%s/relationships/organizations/%s", i.GUID, orgGuid))
	resp, err := i.c.DoRequest(req)
	if err != nil {
		return errors.Wrapf(err, "Error during removing org %s in isolation segment %s", orgGuid, i.Name)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Error deleting org %s in isolation segment %s, response code: %d", orgGuid, i.Name, resp.StatusCode)
	}
	return nil
}

func (i *IsolationSegment) AddSpace(spaceGuid string) error {
	if i == nil || i.c == nil {
		return errors.New("No communication handle.")
	}
	req := i.c.NewRequest("PUT", fmt.Sprintf("/v2/spaces/%s", spaceGuid))
	req.obj = map[string]interface{}{
		"isolation_segment_guid": i.GUID,
	}
	resp, err := i.c.DoRequest(req)
	if err != nil {
		return errors.Wrapf(err, "Error during adding space %s to isolation segment %s", spaceGuid, i.Name)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error adding space to isolation segment %s, response code: %d", i.Name, resp.StatusCode)
	}
	return nil
}

func (i *IsolationSegment) RemoveSpace(spaceGuid string) error {
	if i == nil || i.c == nil {
		return errors.New("No communication handle.")
	}
	req := i.c.NewRequest("DELETE", fmt.Sprintf("/v2/spaces/%s/isolation_segment", spaceGuid))
	resp, err := i.c.DoRequest(req)
	if err != nil {
		return errors.Wrapf(err, "Error during deleting space %s in isolation segment %s", spaceGuid, i.Name)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Error deleting space %s from isolation segment %s, response code: %d", spaceGuid, i.Name, resp.StatusCode)
	}
	return nil
}
