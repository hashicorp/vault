package cfclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
)

// ProcessListResponse is the json body returned from the API
type ProcessListResponse struct {
	Pagination Pagination `json:"pagination"`
	Processes  []Process  `json:"resources"`
}

// Process represents a running process in a container.
type Process struct {
	GUID        string `json:"guid"`
	Type        string `json:"type"`
	Instances   int    `json:"instances"`
	MemoryInMB  int    `json:"memory_in_mb"`
	DiskInMB    int    `json:"disk_in_mb"`
	Ports       []int  `json:"ports,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	HealthCheck struct {
		Type string `json:"type"`
		Data struct {
			Timeout           int    `json:"timeout"`
			InvocationTimeout int    `json:"invocation_timeout"`
			Endpoint          string `json:"endpoint"`
		} `json:"data"`
	} `json:"health_check"`
	Links struct {
		Self  Link `json:"self"`
		Scale Link `json:"scale"`
		App   Link `json:"app"`
		Space Link `json:"space"`
		Stats Link `json:"stats"`
	} `json:"links"`
}

// ListAllProcesses will call the v3 processes api
func (c *Client) ListAllProcesses() ([]Process, error) {
	return c.ListAllProcessesByQuery(url.Values{})
}

// ListAllProcessesByQuery will call the v3 processes api
func (c *Client) ListAllProcessesByQuery(query url.Values) ([]Process, error) {
	var allProcesses []Process

	urlPath := "/v3/processes"
	for {
		resp, err := c.getProcessPage(urlPath, query)
		if err != nil {
			return nil, err
		}

		if resp.Pagination.TotalResults == 0 {
			return nil, nil
		}

		if allProcesses == nil {
			allProcesses = make([]Process, 0, resp.Pagination.TotalResults)
		}

		allProcesses = append(allProcesses, resp.Processes...)
		if resp.Pagination.Next == nil {
			return allProcesses, nil
		}

		var nextURL string

		if resp.Pagination.Next == nil {
			return allProcesses, nil
		}

		switch resp.Pagination.Next.(type) {
		case string:
			nextURL = resp.Pagination.Next.(string)
		case map[string]interface{}:
			m := resp.Pagination.Next.(map[string]interface{})
			u, ok := m["href"]
			if ok {
				nextURL = u.(string)
			}
		default:
			return nil, fmt.Errorf("Unexpected type [%s] for next url", reflect.TypeOf(resp.Pagination.Next).String())
		}

		if nextURL == "" {
			return allProcesses, nil
		}

		u, err := url.Parse(nextURL)
		if err != nil {
			return nil, err
		}

		urlPath = u.Path
		query, err = url.ParseQuery(u.RawQuery)
		if err != nil {
			return nil, err
		}
	}
}

func (c *Client) getProcessPage(urlPath string, query url.Values) (*ProcessListResponse, error) {
	req := c.NewRequest("GET", fmt.Sprintf("%s?%s", urlPath, query.Encode()))

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	procResp := new(ProcessListResponse)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(procResp)
	if err != nil {
		return nil, err
	}

	return procResp, nil
}
