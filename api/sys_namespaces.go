package api

import (
	"fmt"
	"net/http"
)

// ListNamespacesResponse is the response from the ListNamespaces call.
type ListNamespacesResponse struct {
	// NamespacePaths is the list of child namespace paths
	NamespacePaths []string `json:"namespace_paths"`
}

type GetNamespaceResponse struct {
	Path string `json:"path"`
}

// ListNamespaces lists any existing namespace relative to the namespace
// provided in the client's namespace header.
func (c *Sys) ListNamespaces() (*ListNamespacesResponse, error) {
	r := c.c.NewRequest("LIST", "/v1/sys/namespaces")

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}

	return &ListNamespacesResponse{NamespacePaths: result.Data.Keys}, nil
}

// GetNamespace returns namespace information
func (c *Sys) GetNamespace(path string) (*GetNamespaceResponse, error) {
	r := c.c.NewRequest("GET", fmt.Sprintf("/v1/sys/namespaces/%s", path))
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	ret := &GetNamespaceResponse{}
	result := map[string]interface{}{
		"data": map[string]interface{}{},
	}
	if err := resp.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if data, ok := result["data"]; ok {
		if pathOk, ok := data.(map[string]interface{})["path"]; ok {
			if pathRaw, ok := pathOk.(string); ok {
				ret.Path = pathRaw
			}
		}
	}

	return ret, nil
}

// CreateNamespace creates a new namespace relative to the namespace provided
// in the client's namespace header.
func (c *Sys) CreateNamespace(path string) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/sys/namespaces/%s", path))
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteNamespace delete an existing namespace relative to the namespace
// provided in the client's namespace header.
func (c *Sys) DeleteNamespace(path string) error {
	r := c.c.NewRequest("DELETE", fmt.Sprintf("/v1/sys/namespaces/%s", path))
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
