package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/mapstructure"
)

// ListPluginsInput is used as input to the ListPlugins function.
type ListPluginsInput struct {
	// Type of the plugin. Required.
	Type consts.PluginType `json:"type"`
}

// ListPluginsResponse is the response from the ListPlugins call.
type ListPluginsResponse struct {
	// PluginsByType is the list of plugins by type.
	PluginsByType map[consts.PluginType][]string `json:"types"`

	// Names is the list of names of the plugins.
	//
	// Deprecated: Newer server responses should be returning PluginsByType (json:
	// "types") instead.
	Names []string `json:"names"`
}

// ListPlugins lists all plugins in the catalog and returns their names as a
// list of strings.
func (c *Sys) ListPlugins(i *ListPluginsInput) (*ListPluginsResponse, error) {
	path := ""
	method := ""
	if i.Type == consts.PluginTypeUnknown {
		path = "/v1/sys/plugins/catalog"
		method = "GET"
	} else {
		path = fmt.Sprintf("/v1/sys/plugins/catalog/%s", i.Type)
		method = "LIST"
	}

	req := c.c.NewRequest(method, path)
	if method == "LIST" {
		// Set this for broader compatibility, but we use LIST above to be able
		// to handle the wrapping lookup function
		req.Method = "GET"
		req.Params.Set("list", "true")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, req)
	if err != nil && resp == nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	// We received an Unsupported Operation response from Vault, indicating
	// Vault of an older version that doesn't support the GET method yet;
	// switch it to a LIST.
	if resp.StatusCode == 405 {
		req.Params.Set("list", "true")
		resp, err := c.c.RawRequestWithContext(ctx, req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var result struct {
			Data struct {
				Keys []string `json:"keys"`
			} `json:"data"`
		}
		if err := resp.DecodeJSON(&result); err != nil {
			return nil, err
		}
		return &ListPluginsResponse{Names: result.Data.Keys}, nil
	}

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	result := &ListPluginsResponse{
		PluginsByType: make(map[consts.PluginType][]string),
	}
	if i.Type == consts.PluginTypeUnknown {
		for pluginTypeStr, pluginsRaw := range secret.Data {
			pluginType, err := consts.ParsePluginType(pluginTypeStr)
			if err != nil {
				return nil, err
			}

			pluginsIfc, ok := pluginsRaw.([]interface{})
			if !ok {
				return nil, fmt.Errorf("unable to parse plugins for %q type", pluginTypeStr)
			}

			plugins := make([]string, len(pluginsIfc))
			for i, nameIfc := range pluginsIfc {
				name, ok := nameIfc.(string)
				if !ok {
				}
				plugins[i] = name
			}
			result.PluginsByType[pluginType] = plugins
		}
	} else {
		var respKeys []string
		if err := mapstructure.Decode(secret.Data["keys"], &respKeys); err != nil {
			return nil, err
		}
		result.PluginsByType[i.Type] = respKeys
	}

	return result, nil
}

// GetPluginInput is used as input to the GetPlugin function.
type GetPluginInput struct {
	Name string `json:"-"`

	// Type of the plugin. Required.
	Type consts.PluginType `json:"type"`
}

// GetPluginResponse is the response from the GetPlugin call.
type GetPluginResponse struct {
	Args    []string `json:"args"`
	Builtin bool     `json:"builtin"`
	Command string   `json:"command"`
	Name    string   `json:"name"`
	SHA256  string   `json:"sha256"`
}

// GetPlugin retrieves information about the plugin.
func (c *Sys) GetPlugin(i *GetPluginInput) (*GetPluginResponse, error) {
	path := catalogPathByType(i.Type, i.Name)
	req := c.c.NewRequest(http.MethodGet, path)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data *GetPluginResponse
	}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}
	return result.Data, err
}

// RegisterPluginInput is used as input to the RegisterPlugin function.
type RegisterPluginInput struct {
	// Name is the name of the plugin. Required.
	Name string `json:"-"`

	// Type of the plugin. Required.
	Type consts.PluginType `json:"type"`

	// Args is the list of args to spawn the process with.
	Args []string `json:"args,omitempty"`

	// Command is the command to run.
	Command string `json:"command,omitempty"`

	// SHA256 is the shasum of the plugin.
	SHA256 string `json:"sha256,omitempty"`
}

// RegisterPlugin registers the plugin with the given information.
func (c *Sys) RegisterPlugin(i *RegisterPluginInput) error {
	path := catalogPathByType(i.Type, i.Name)
	req := c.c.NewRequest(http.MethodPut, path)

	if err := req.SetJSONBody(i); err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// DeregisterPluginInput is used as input to the DeregisterPlugin function.
type DeregisterPluginInput struct {
	// Name is the name of the plugin. Required.
	Name string `json:"-"`

	// Type of the plugin. Required.
	Type consts.PluginType `json:"type"`
}

// DeregisterPlugin removes the plugin with the given name from the plugin
// catalog.
func (c *Sys) DeregisterPlugin(i *DeregisterPluginInput) error {
	path := catalogPathByType(i.Type, i.Name)
	req := c.c.NewRequest(http.MethodDelete, path)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// ReloadPluginInput is used as input to the ReloadPlugin function.
type ReloadPluginInput struct {
	// Plugin is the name of the plugin to reload, as registered in the plugin catalog
	Plugin string `json:"plugin"`

	// Mounts is the array of string mount paths of the plugin backends to reload
	Mounts []string `json:"mounts"`

	// Scope is the scope of the plugin reload
	Scope string `json:"scope"`
}

// ReloadPlugin reloads mounted plugin backends, possibly returning
// reloadId for a cluster scoped reload
func (c *Sys) ReloadPlugin(i *ReloadPluginInput) (string, error) {
	path := "/v1/sys/plugins/reload/backend"
	req := c.c.NewRequest(http.MethodPut, path)

	if err := req.SetJSONBody(i); err != nil {
		return "", err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	resp, err := c.c.RawRequestWithContext(ctx, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if i.Scope == "global" {
		// Get the reload id
		secret, parseErr := ParseSecret(resp.Body)
		if parseErr != nil {
			return "", parseErr
		}
		if _, ok := secret.Data["reload_id"]; ok {
			return secret.Data["reload_id"].(string), nil
		}
	}
	return "", err
}

// ReloadStatus is the status of an individual node's plugin reload
type ReloadStatus struct {
	Timestamp time.Time `json:"timestamp" mapstructure:"timestamp"`
	Error     string    `json:"error" mapstructure:"error"`
}

// ReloadStatusResponse is the combined response of all known completed plugin reloads
type ReloadStatusResponse struct {
	ReloadID string                   `mapstructure:"reload_id"`
	Results  map[string]*ReloadStatus `mapstructure:"results"`
}

// ReloadPluginStatusInput is used as input to the ReloadStatusPlugin function.
type ReloadPluginStatusInput struct {
	// ReloadID is the ID of the reload operation
	ReloadID string `json:"reload_id"`
}

// ReloadPluginStatus retrieves the status of a reload operation
func (c *Sys) ReloadPluginStatus(reloadStatusInput *ReloadPluginStatusInput) (*ReloadStatusResponse, error) {
	path := "/v1/sys/plugins/reload/backend/status"
	req := c.c.NewRequest(http.MethodGet, path)
	req.Params.Add("reload_id", reloadStatusInput.ReloadID)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	resp, err := c.c.RawRequestWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp != nil {
		secret, parseErr := ParseSecret(resp.Body)
		if parseErr != nil {
			return nil, err
		}

		var r ReloadStatusResponse
		d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339),
			Result:     &r,
		})
		if err != nil {
			return nil, err
		}
		err = d.Decode(secret.Data)
		if err != nil {
			return nil, err
		}
		return &r, nil
	}
	return nil, nil
}

// catalogPathByType is a helper to construct the proper API path by plugin type
func catalogPathByType(pluginType consts.PluginType, name string) string {
	path := fmt.Sprintf("/v1/sys/plugins/catalog/%s/%s", pluginType, name)

	// Backwards compat, if type is not provided then use old path
	if pluginType == consts.PluginTypeUnknown {
		path = fmt.Sprintf("/v1/sys/plugins/catalog/%s", name)
	}

	return path
}
