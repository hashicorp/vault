package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/pkg/errors"
)

type StacksResponse struct {
	Count     int              `json:"total_results"`
	Pages     int              `json:"total_pages"`
	NextUrl   string           `json:"next_url"`
	Resources []StacksResource `json:"resources"`
}

type StacksResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Stack `json:"entity"`
}

type Stack struct {
	Guid        string `json:"guid"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Description string `json:"description"`
	c           *Client
}

func (c *Client) ListStacksByQuery(query url.Values) ([]Stack, error) {
	var stacks []Stack
	requestURL := "/v2/stacks?" + query.Encode()
	for {
		stacksResp, err := c.getStacksResponse(requestURL)
		if err != nil {
			return []Stack{}, err
		}
		for _, stack := range stacksResp.Resources {
			stack.Entity.Guid = stack.Meta.Guid
			stack.Entity.CreatedAt = stack.Meta.CreatedAt
			stack.Entity.UpdatedAt = stack.Meta.UpdatedAt
			stack.Entity.c = c
			stacks = append(stacks, stack.Entity)
		}
		requestURL = stacksResp.NextUrl
		if requestURL == "" {
			break
		}
	}
	return stacks, nil
}

func (c *Client) ListStacks() ([]Stack, error) {
	return c.ListStacksByQuery(nil)
}

func (c *Client) GetStackByGuid(stackGUID string) (Stack, error) {
	var stacksRes StacksResource
	requestURL := fmt.Sprintf("/v2/stacks/%s", stackGUID)
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Stack{}, errors.Wrap(err, "Error requesting stack info")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return Stack{}, errors.Wrap(err, "Error reading stack body")
	}
	err = json.Unmarshal(resBody, &stacksRes)
	if err != nil {
		return Stack{}, errors.Wrap(err, "Error unmarshalling stack")
	}

	stacksRes.Entity.Guid = stacksRes.Meta.Guid
	stacksRes.Entity.CreatedAt = stacksRes.Meta.CreatedAt
	stacksRes.Entity.UpdatedAt = stacksRes.Meta.UpdatedAt
	stacksRes.Entity.c = c

	return stacksRes.Entity, nil
}

func (c *Client) getStacksResponse(requestURL string) (StacksResponse, error) {
	var stacksResp StacksResponse
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return StacksResponse{}, errors.Wrap(err, "Error requesting stacks")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return StacksResponse{}, errors.Wrap(err, "Error reading stacks body")
	}
	err = json.Unmarshal(resBody, &stacksResp)
	if err != nil {
		return StacksResponse{}, errors.Wrap(err, "Error unmarshalling stacks")
	}
	return stacksResp, nil
}
