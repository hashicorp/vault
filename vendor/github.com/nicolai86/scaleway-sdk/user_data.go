package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Userdatas represents the response of a GET /user_data
type Userdatas struct {
	UserData []string `json:"user_data"`
}

// Userdata represents []byte
type Userdata []byte

// GetUserdatas gets list of userdata for a server
func (s *API) GetUserdatas(serverID string, metadata bool) (*Userdatas, error) {
	var uri, endpoint string

	endpoint = s.computeAPI
	if metadata {
		uri = "/user_data"
		endpoint = MetadataAPI
	} else {
		uri = fmt.Sprintf("servers/%s/user_data", serverID)
	}

	resp, err := s.GetResponsePaginate(endpoint, uri, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var userdatas Userdatas

	if err = json.Unmarshal(body, &userdatas); err != nil {
		return nil, err
	}
	return &userdatas, nil
}

func (s *Userdata) String() string {
	return string(*s)
}

// GetUserdata gets a specific userdata for a server
func (s *API) GetUserdata(serverID, key string, metadata bool) (*Userdata, error) {
	var uri, endpoint string

	endpoint = s.computeAPI
	if metadata {
		uri = fmt.Sprintf("/user_data/%s", key)
		endpoint = MetadataAPI
	} else {
		uri = fmt.Sprintf("servers/%s/user_data/%s", serverID, key)
	}

	var err error
	resp, err := s.GetResponsePaginate(endpoint, uri, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no such user_data %q (%d)", key, resp.StatusCode)
	}
	var data Userdata
	data, err = ioutil.ReadAll(resp.Body)
	return &data, err
}

// PatchUserdata sets a user data
func (s *API) PatchUserdata(serverID, key string, value []byte, metadata bool) error {
	var resource, endpoint string

	endpoint = s.computeAPI
	if metadata {
		resource = fmt.Sprintf("/user_data/%s", key)
		endpoint = MetadataAPI
	} else {
		resource = fmt.Sprintf("servers/%s/user_data/%s", serverID, key)
	}

	uri := fmt.Sprintf("%s/%s", strings.TrimRight(endpoint, "/"), resource)
	payload := new(bytes.Buffer)
	payload.Write(value)

	req, err := http.NewRequest("PATCH", uri, payload)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("cannot set user_data (%d)", resp.StatusCode)
}

// DeleteUserdata deletes a server user_data
func (s *API) DeleteUserdata(serverID, key string, metadata bool) error {
	var url, endpoint string

	endpoint = s.computeAPI
	if metadata {
		url = fmt.Sprintf("/user_data/%s", key)
		endpoint = MetadataAPI
	} else {
		url = fmt.Sprintf("servers/%s/user_data/%s", serverID, key)
	}

	resp, err := s.DeleteResponse(endpoint, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}
