package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const (
	ErrOutputPolicyRequest = "output a policy, please"
)

var LastOutputPolicyError *OutputPolicyError

type OutputPolicyError struct {
	*retryablehttp.Request
	vaultServerAddress string
	finalHCLString     string
}

func (d *OutputPolicyError) Error() string {
	if d.finalHCLString == "" {
		p, err := d.buildSamplePolicy()
		if err != nil {
			return err.Error()
		}
		d.finalHCLString = p
	}

	return ErrOutputPolicyRequest
}

func (d *OutputPolicyError) HCLString() (string, error) {
	if d.finalHCLString == "" {
		p, err := d.buildSamplePolicy()
		if err != nil {
			return "", err
		}
		d.finalHCLString = p
	}
	return d.finalHCLString, nil
}

// Builds a sample policy document from the request
func (d *OutputPolicyError) buildSamplePolicy() (string, error) {
	var capabilities []string
	switch d.Request.Method {
	case http.MethodGet, "":
		capabilities = append(capabilities, "read")
	case http.MethodPost, http.MethodPut:
		capabilities = append(capabilities, "create")
		capabilities = append(capabilities, "update")
	case http.MethodPatch:
		capabilities = append(capabilities, "patch")
	case http.MethodDelete:
		capabilities = append(capabilities, "delete")
	case "LIST":
		capabilities = append(capabilities, "list")
	}

	// sanitize, then trim the Vault address and v1 from the front of the path
	url, err := url.PathUnescape(d.Request.URL.String())
	if err != nil {
		return "", fmt.Errorf("failed to unescape request URL characters: %v", err)
	}
	apiAddrPrefix := fmt.Sprintf("%sv1/", d.vaultServerAddress)
	path := strings.TrimLeft(url, apiAddrPrefix)

	// determine whether to add sudo capability
	if isSudoPath(path) {
		capabilities = append(capabilities, "sudo")
	}

	capStr := strings.Join(capabilities, `", "`)
	return fmt.Sprintf(
		`path "%s" {
  capabilities = ["%s"]
}`, path, capStr), nil
}

// Determine whether the given path requires the sudo capability
func isSudoPath(path string) bool {
	// Return early if the path is any of the non-templated sudo paths.
	if _, ok := sudoPaths[path]; ok {
		return true
	}

	// Some sudo paths have templated fields in them.
	// (e.g. sys/revoke-prefix/{prefix})
	// The values in the sudoPaths map are actually regular expressions,
	// so we can check if our path matches against them.
	for _, sudoPathRegexp := range sudoPaths {
		match := sudoPathRegexp.Match([]byte(fmt.Sprintf("/%s", path))) // the OpenAPI response has a / in front of each path
		if match {
			return true
		}
	}

	return false
}
