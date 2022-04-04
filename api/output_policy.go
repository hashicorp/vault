package api

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

const (
	ErrOutputPolicyRequest = "output policy request"
)

var LastOutputPolicyError *OutputPolicyError

type OutputPolicyError struct {
	*retryablehttp.Request
	VaultClient     *Client
	parsingError    error
	parsedHCLString string
}

func (d *OutputPolicyError) Error() string {
	if d.parsedHCLString == "" {
		d.parseRequest()
		if d.parsingError != nil {
			return d.parsingError.Error()
		}
	}

	return ErrOutputPolicyRequest
}

// Builds a sample policy document from the request
func (d *OutputPolicyError) parseRequest() {
	capabilities := []string{}
	switch d.Request.Method {
	case "GET":
		capabilities = append(capabilities, "read")
	case "LIST":
		capabilities = append(capabilities, "list")
	case "POST", "PUT":
		capabilities = append(capabilities, "create")
		capabilities = append(capabilities, "update")
	case "PATCH":
		capabilities = append(capabilities, "patch")
	case "DELETE":
		capabilities = append(capabilities, "delete")
	}

	// trim the Vault address and v1 from the front of the path
	url := d.Request.URL.String()
	apiAddrPrefix := fmt.Sprintf("%sv1/", d.VaultClient.config.Address)
	path := strings.Trim(url, apiAddrPrefix)

	// determine whether to add sudo capability
	needsSudo, err := isSudoPath(d.VaultClient, path)
	if err != nil {
		d.parsingError = err
		return
	}
	if needsSudo {
		capabilities = append(capabilities, "sudo")
	}

	capStr := strings.Join(capabilities, `", "`)
	d.parsedHCLString = fmt.Sprintf(
		`path "%s" {
  capabilities = ["%s"]
}`, path, capStr)
}

func (d *OutputPolicyError) HCLString() string {
	if d.parsedHCLString == "" {
		d.parseRequest()
	}
	return d.parsedHCLString
}

// Determine whether the given path requires the sudo capability
func isSudoPath(client *Client, path string) (bool, error) {
	sudoPaths, err := getSudoPaths(client)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve list of paths that require sudo capability: %v", err)
	}
	if sudoPaths == nil || len(sudoPaths) < 1 {
		// OpenAPI spec did not return any paths that require sudo for
		// some reason, but the user probably still shouldn't see an error.
		return false, nil
	}

	// Return early if the path is clearly one of the sudo paths.
	if _, ok := sudoPaths[path]; ok {
		return true, nil
	}

	// Some sudo paths have templated fields in them.
	// (e.g. sys/revoke-prefix/{prefix})
	// The keys in the sudoPaths map are actually regular expressions,
	// so we can check if our path matches against them.
	for sudoPath := range sudoPaths {
		r, err := regexp.Compile(fmt.Sprintf("^%s$", sudoPath))
		if err != nil {
			continue
		}

		match := r.Match([]byte(fmt.Sprintf("/%s", path))) // the OpenAPI response has a / in front of each path
		if match {
			return true, nil
		}
	}

	return false, nil
}

func getSudoPaths(client *Client) (map[string]bool, error) {
	// We don't want to use a wrapping call or any special flags here
	// so save any custom value and restore after.
	currentWrappingLookupFunc := client.CurrentWrappingLookupFunc()
	client.SetWrappingLookupFunc(nil)
	defer client.SetWrappingLookupFunc(currentWrappingLookupFunc)
	currentOutputCurlString := client.OutputCurlString()
	client.SetOutputCurlString(false)
	defer client.SetOutputCurlString(currentOutputCurlString)
	currentOutputPolicy := client.OutputPolicy()
	client.SetOutputPolicy(false)
	defer client.SetOutputPolicy(currentOutputPolicy)

	r := client.NewRequest("GET", "/v1/sys/internal/specs/openapi")
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve sudo endpoints: %v", err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body from OpenAPI: %v", err)
	}
	if buf.Len() == 0 {
		return nil, fmt.Errorf("OpenAPI response had no content")
	}

	oasInfo := make(map[string]interface{})
	if err := jsonutil.DecodeJSONFromReader(&buf, &oasInfo); err != nil {
		return nil, fmt.Errorf("unable to decode JSON from OpenAPI response: %v", err)
	}

	paths, ok := oasInfo["paths"]
	if !ok {
		return nil, fmt.Errorf("OpenAPI response did not include paths")
	}

	pathsMap, ok := paths.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("OpenAPI response did not return valid paths")
	}

	sudoPaths := make(map[string]bool) // this could be a slice, but we're just making it a map so we can do quick lookups for the paths that don't have any templating
	for pathName, pathInfo := range pathsMap {
		pathInfoMap, ok := pathInfo.(map[string]interface{})
		if !ok {
			continue
		}

		if sudo, ok := pathInfoMap["x-vault-sudo"]; ok {
			if sudo == true {
				// Since many paths have templated fields like {name},
				// our list of sudo paths will actually be a list of
				// regular expressions that we can match against.
				pathRegex := buildPathRegexp(pathName)
				sudoPaths[pathRegex] = true
			}
		}
	}

	return sudoPaths, nil
}

// Replaces any template fields in a path with the characters ".+" so that
// we can later allow any characters to match those fields.
func buildPathRegexp(pathName string) string {
	templateFields := []string{"{path}", "{header}", "{prefix}", "{name}", "{type}"}
	pathWithRegexPatterns := pathName
	for _, field := range templateFields {
		r, _ := regexp.Compile(field)
		pathWithRegexPatterns = r.ReplaceAllString(pathWithRegexPatterns, ".+")
	}

	return pathWithRegexPatterns
}
