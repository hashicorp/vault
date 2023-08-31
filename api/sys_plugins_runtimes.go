// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// GetPluginRuntimeInput is used as input to the GetPluginRuntime function.
type GetPluginRuntimeInput struct {
	Name string `json:"-"`

	// Type of the plugin runtime. Required.
	Type PluginRuntimeType `json:"type"`
}

// GetPluginRuntimeResponse is the response from the GetPluginRuntime call.
type GetPluginRuntimeResponse struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	OCIRuntime   string `json:"oci_runtime"`
	CgroupParent string `json:"cgroup_parent"`
	CPU          int64  `json:"cpu"`
	Memory       int64  `json:"memory"`
}

// GetPluginRuntime retrieves information about the plugin.
func (c *Sys) GetPluginRuntime(ctx context.Context, i *GetPluginRuntimeInput) (*GetPluginRuntimeResponse, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	path := pluginRuntimeCatalogPathByType(i.Type, i.Name)
	req := c.c.NewRequest(http.MethodGet, path)

	resp, err := c.c.rawRequestWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data *GetPluginRuntimeResponse
	}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}
	return result.Data, err
}

// RegisterPluginRuntimeInput is used as input to the RegisterPluginRuntime function.
type RegisterPluginRuntimeInput struct {
	// Name is the name of the plugin. Required.
	Name string `json:"-"`

	// Type of the plugin. Required.
	Type PluginRuntimeType `json:"type"`

	OCIRuntime   string `json:"oci_runtime,omitempty"`
	CgroupParent string `json:"cgroup_parent,omitempty"`
	CPU          int64  `json:"cpu,omitempty"`
	Memory       int64  `json:"memory,omitempty"`
}

// RegisterPluginRuntime registers the plugin with the given information.
func (c *Sys) RegisterPluginRuntime(ctx context.Context, i *RegisterPluginRuntimeInput) error {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	path := pluginRuntimeCatalogPathByType(i.Type, i.Name)
	req := c.c.NewRequest(http.MethodPut, path)

	if err := req.SetJSONBody(i); err != nil {
		return err
	}

	resp, err := c.c.rawRequestWithContext(ctx, req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// DeregisterPluginRuntimeInput is used as input to the DeregisterPluginRuntime function.
type DeregisterPluginRuntimeInput struct {
	// Name is the name of the plugin runtime. Required.
	Name string `json:"-"`

	// Type of the plugin. Required.
	Type PluginRuntimeType `json:"type"`
}

// DeregisterPluginRuntime removes the plugin with the given name from the plugin
// catalog.
func (c *Sys) DeregisterPluginRuntime(ctx context.Context, i *DeregisterPluginRuntimeInput) error {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	path := pluginRuntimeCatalogPathByType(i.Type, i.Name)
	req := c.c.NewRequest(http.MethodDelete, path)
	resp, err := c.c.rawRequestWithContext(ctx, req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

type PluginRuntimeDetails struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	OCIRuntime   string `json:"oci_runtime"`
	CgroupParent string `json:"cgroup_parent"`
	CPU          int64  `json:"cpu"`
	Memory       int64  `json:"memory"`
}

// ListPluginRuntimesInput is used as input to the ListPluginRuntimes function.
type ListPluginRuntimesInput struct {
	// Type of the plugin. Required.
	Type PluginRuntimeType `json:"type"`
}

// ListPluginRuntimesResponse is the response from the ListPluginRuntimes call.
type ListPluginRuntimesResponse struct {
	// RuntimesByType is the list of plugin runtimes by type.
	Runtimes []PluginRuntimeDetails `json:"runtimes"`
}

// ListPluginRuntimes lists all plugin runtimes in the catalog and returns their names as a
// list of strings.
func (c *Sys) ListPluginRuntimes(ctx context.Context, i *ListPluginRuntimesInput) (*ListPluginRuntimesResponse, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	if i != nil && i.Type == PluginRuntimeTypeUnsupported {
		return nil, fmt.Errorf("%q is not a supported runtime type", i.Type.String())
	}

	resp, err := c.c.rawRequestWithContext(ctx, c.c.NewRequest(http.MethodGet, "/v1/sys/plugins/runtimes/catalog"))
	if err != nil && resp == nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}
	if _, ok := secret.Data["runtimes"]; !ok {
		return nil, fmt.Errorf("data from server response does not contain runtimes")
	}

	runtimesRaw, ok := secret.Data["runtimes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to parse runtimes")
	}

	result := &ListPluginRuntimesResponse{
		Runtimes: []PluginRuntimeDetails{},
	}
	var runtimes []PluginRuntimeDetails
	for _, runtimeRaw := range runtimesRaw {
		var runtime PluginRuntimeDetails
		if err = mapstructure.Decode(runtimeRaw, &runtime); err != nil {
			return nil, err
		}
		runtimes = append(runtimes, runtime)
	}

	// return all runtimes in the catalog
	if i == nil {
		result.Runtimes = runtimes
		return result, nil
	}

	switch i.Type {
	default:
		for _, runtime := range runtimes {
			if runtime.Type == i.Type.String() {
				result.Runtimes = append(result.Runtimes, runtime)
			}
		}
	}
	return result, nil
}

// pluginRuntimeCatalogPathByType is a helper to construct the proper API path by plugin type
func pluginRuntimeCatalogPathByType(runtimeType PluginRuntimeType, name string) string {
	return fmt.Sprintf("/v1/sys/plugins/catalog/%s/%s", runtimeType, name)
}
