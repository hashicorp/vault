// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
)

// CheckBundleConfig configuration specific to check type
type CheckBundleConfig struct {
	AsyncMetrics  bool   `json:"async_metrics"`
	Secret        string `json:"secret"`
	SubmissionURL string `json:"submission_url"`
	ReverseSecret string `json:"reverse:secret_key"`
	HTTPVersion   string `json:"http_version,omitempty"`
	Method        string `json:"method,omitempty"`
	Payload       string `json:"payload,omitempty"`
	Port          string `json:"port,omitempty"`
	ReadLimit     string `json:"read_limit,omitempty"`
	URL           string `json:"url,omitempty"`
}

// CheckBundleMetric individual metric configuration
type CheckBundleMetric struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Units  string `json:"units"`
	Status string `json:"status"`
}

// CheckBundle definition
type CheckBundle struct {
	CheckUUIDs         []string            `json:"_check_uuids,omitempty"`
	Checks             []string            `json:"_checks,omitempty"`
	Cid                string              `json:"_cid,omitempty"`
	Created            int                 `json:"_created,omitempty"`
	LastModified       int                 `json:"_last_modified,omitempty"`
	LastModifedBy      string              `json:"_last_modifed_by,omitempty"`
	ReverseConnectURLs []string            `json:"_reverse_connection_urls"`
	Brokers            []string            `json:"brokers"`
	Config             CheckBundleConfig   `json:"config"`
	DisplayName        string              `json:"display_name"`
	Metrics            []CheckBundleMetric `json:"metrics"`
	MetricLimit        int                 `json:"metric_limit"`
	Notes              string              `json:"notes"`
	Period             int                 `json:"period"`
	Status             string              `json:"status"`
	Tags               []string            `json:"tags"`
	Target             string              `json:"target"`
	Timeout            int                 `json:"timeout"`
	Type               string              `json:"type"`
}

// FetchCheckBundleByID fetch a check bundle configuration by id
func (a *API) FetchCheckBundleByID(id IDType) (*CheckBundle, error) {
	cid := CIDType(fmt.Sprintf("/check_bundle/%d", id))
	return a.FetchCheckBundleByCID(cid)
}

// FetchCheckBundleByCID fetch a check bundle configuration by id
func (a *API) FetchCheckBundleByCID(cid CIDType) (*CheckBundle, error) {
	result, err := a.Get(string(cid))
	if err != nil {
		return nil, err
	}

	checkBundle := &CheckBundle{}
	if err := json.Unmarshal(result, checkBundle); err != nil {
		return nil, err
	}

	return checkBundle, nil
}

// CheckBundleSearch returns list of check bundles matching a search query
//    - a search query not a filter (see: https://login.circonus.com/resources/api#searching)
func (a *API) CheckBundleSearch(searchCriteria SearchQueryType) ([]CheckBundle, error) {
	apiPath := fmt.Sprintf("/check_bundle?search=%s", searchCriteria)

	response, err := a.Get(apiPath)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] API call error %+v", err)
	}

	var results []CheckBundle
	if err := json.Unmarshal(response, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// CreateCheckBundle create a new check bundle (check)
func (a *API) CreateCheckBundle(config CheckBundle) (*CheckBundle, error) {
	cfgJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	response, err := a.Post("/check_bundle", cfgJSON)
	if err != nil {
		return nil, err
	}

	checkBundle := &CheckBundle{}
	if err := json.Unmarshal(response, checkBundle); err != nil {
		return nil, err
	}

	return checkBundle, nil
}

// UpdateCheckBundle updates a check bundle configuration
func (a *API) UpdateCheckBundle(config *CheckBundle) (*CheckBundle, error) {
	if a.Debug {
		a.Log.Printf("[DEBUG] Updating check bundle with new metrics.")
	}

	cfgJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	response, err := a.Put(config.Cid, cfgJSON)
	if err != nil {
		return nil, err
	}

	checkBundle := &CheckBundle{}
	if err := json.Unmarshal(response, checkBundle); err != nil {
		return nil, err
	}

	return checkBundle, nil
}
