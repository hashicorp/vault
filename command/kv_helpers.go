package command

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/consts"
)

func kvReadRequest(client *api.Client, path string, params map[string]string) (*api.Secret, error) {
	r := client.NewRequest("GET", "/v1/"+path)
	if r.Headers == nil {
		r.Headers = http.Header{}
	}
	r.Headers.Add(consts.VaultKVCLIClientHeader, "v1")

	for k, v := range params {
		r.Params.Set(k, v)
	}
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func kvListRequest(client *api.Client, path string) (*api.Secret, error) {
	r := client.NewRequest("LIST", "/v1/"+path)
	if r.Headers == nil {
		r.Headers = http.Header{}
	}
	r.Headers.Add(consts.VaultKVCLIClientHeader, "v1")

	// Set this for broader compatibility, but we use LIST above to be able to
	// handle the wrapping lookup function
	r.Method = "GET"
	r.Params.Set("list", "true")
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func kvWriteRequest(client *api.Client, path string, data map[string]interface{}) (*api.Secret, error) {
	r := client.NewRequest("PUT", "/v1/"+path)
	if r.Headers == nil {
		r.Headers = http.Header{}
	}
	r.Headers.Add(consts.VaultKVCLIClientHeader, "v1")
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return api.ParseSecret(resp.Body)
	}

	return nil, nil
}

func kvDeleteRequest(client *api.Client, path string) (*api.Secret, error) {
	r := client.NewRequest("DELETE", "/v1/"+path)
	if r.Headers == nil {
		r.Headers = http.Header{}
	}
	r.Headers.Add(consts.VaultKVCLIClientHeader, "v1")
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return api.ParseSecret(resp.Body)
	}

	return nil, nil
}

func addPrefixToVKVPath(p, apiPrefix string) (string, error) {
	parts := strings.SplitN(p, "/", 2)
	if len(parts) != 2 {
		return "", errors.New("Invalid path")
	}

	return path.Join(parts[0], apiPrefix, parts[1]), nil
}

func getHeaderForMap(header string, data map[string]interface{}) string {
	maxKey := 0
	for k := range data {
		if len(k) > maxKey {
			maxKey = len(k)
		}
	}

	// 4 for the column spaces and 5 for the len("value")
	totalLen := maxKey + 4 + 5

	equalSigns := totalLen - (len(header) + 2)

	// If we have zero or fewer equal signs bump it back up to two on either
	// side of the header.
	if equalSigns <= 0 {
		equalSigns = 4
	}

	// If the number of equal signs is not divisible by two add a sign.
	if equalSigns%2 != 0 {
		equalSigns = equalSigns + 1
	}

	return fmt.Sprintf("%s %s %s", strings.Repeat("=", equalSigns/2), header, strings.Repeat("=", equalSigns/2))
}
