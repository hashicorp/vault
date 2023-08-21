// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"
	"net/http"
)

type PluginRuntimeDetails struct {
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	OCIRuntime   string  `json:"oci_runtime"`
	ParentCGroup string  `json:"parent_cgroup"`
	CPU          float32 `json:"cpu"`
	Memory       uint64  `json:"memory"`
}

// GetPluginRuntimeInput is used as input to the GetPluginRuntime function.
type GetPluginRuntimeInput struct {
	Name string `json:"-"`

	// Type of the plugin runtime. Required.
	Type PluginRuntimeType `json:"type"`
}

// GetPluginRuntimeResponse is the response from the GetPluginRuntime call.
type GetPluginRuntimeResponse struct {
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	OCIRuntime   string  `json:"oci_runtime"`
	ParentCGroup string  `json:"parent_cgroup"`
	CPU          float32 `json:"cpu"`
	Memory       uint64  `json:"memory"`
}

// GetPluginRuntime wraps GetPluginRuntimeWithContext using context.Background.
func (c *Sys) GetPluginRuntime(i *GetPluginRuntimeInput) (*GetPluginRuntimeResponse, error) {
	return c.GetPluginRuntimeWithContext(context.Background(), i)
}

// GetPluginRuntimeWithContext retrieves information about the plugin.
func (c *Sys) GetPluginRuntimeWithContext(ctx context.Context, i *GetPluginRuntimeInput) (*GetPluginRuntimeResponse, error) {
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

	OCIRuntime   string  `json:"oci_runtime,omitempty"`
	ParentCGroup string  `json:"parent_cgroup,omitempty"`
	CPU          float32 `json:"cpu,omitempty"`
	Memory       uint64  `json:"memory,omitempty"`
}

// RegisterPluginRuntime wraps RegisterPluginWithContext using context.Background.
func (c *Sys) RegisterPluginRuntime(i *RegisterPluginRuntimeInput) error {
	return c.RegisterPluginRuntimeWithContext(context.Background(), i)
}

// RegisterPluginRuntimeWithContext registers the plugin with the given information.
func (c *Sys) RegisterPluginRuntimeWithContext(ctx context.Context, i *RegisterPluginRuntimeInput) error {
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

// DeregisterPluginRuntime wraps DeregisterPluginRuntimeWithContext using context.Background.
func (c *Sys) DeregisterPluginRuntime(i *DeregisterPluginRuntimeInput) error {
	return c.DeregisterPluginRuntimeWithContext(context.Background(), i)
}

// DeregisterPluginRuntimeWithContext removes the plugin with the given name from the plugin
// catalog.
func (c *Sys) DeregisterPluginRuntimeWithContext(ctx context.Context, i *DeregisterPluginRuntimeInput) error {
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

// pluginRuntimeCatalogPathByType is a helper to construct the proper API path by plugin type
func pluginRuntimeCatalogPathByType(runtimeType PluginRuntimeType, name string) string {
	return fmt.Sprintf("/v1/sys/plugins/catalog/%s/%s", runtimeType, name)
}
