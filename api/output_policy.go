package api

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const (
	ErrOutputPolicyRequest = "output a policy, please"
)

var LastOutputPolicyError *OutputPolicyError

var sudoPathsWithRegexps map[string]*regexp.Regexp

func init() {
	setSudoPathRegexps()
}

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
	needsSudo, err := isSudoPath(path)
	if err != nil {
		return "", fmt.Errorf("unable to determine if path needs sudo capability: %v", err)
	}
	if needsSudo {
		capabilities = append(capabilities, "sudo")
	}

	capStr := strings.Join(capabilities, `", "`)
	return fmt.Sprintf(
		`path "%s" {
  capabilities = ["%s"]
}`, path, capStr), nil
}

// Determine whether the given path requires the sudo capability
func isSudoPath(path string) (bool, error) {
	sudoPaths := GetSudoPaths()

	// Return early if the path is any of the non-templated sudo paths.
	if _, ok := sudoPaths[path]; ok {
		return true, nil
	}

	// Some sudo paths have templated fields in them.
	// (e.g. sys/revoke-prefix/{prefix})
	// The values in the sudoPaths map are actually regular expressions,
	// so we can check if our path matches against them.
	for _, sudoPathRegexp := range sudoPaths {
		match := sudoPathRegexp.Match([]byte(fmt.Sprintf("/%s", path))) // the OpenAPI response has a / in front of each path
		if match {
			return true, nil
		}
	}

	return false, nil
}

func GetSudoPaths() map[string]*regexp.Regexp {
	return sudoPathsWithRegexps
}

func setSudoPathRegexps() {
	sudoPathsWithRegexps = map[string]*regexp.Regexp{
		"/auth/token/accessors/":                        regexp.MustCompile("^/auth/token/accessors/$"),
		"/pki/root":                                     regexp.MustCompile("^/pki/root$"),
		"/pki/root/sign-self-issued":                    regexp.MustCompile("^/pki/root/sign-self-issued$"),
		"/sys/audit":                                    regexp.MustCompile("^/sys/audit$"),
		"/sys/audit/{path}":                             regexp.MustCompile("^/sys/audit/.+$"),
		"/sys/auth/{path}":                              regexp.MustCompile("^/sys/auth/.+$"),
		"/sys/auth/{path}/tune":                         regexp.MustCompile("^/sys/auth/.+/tune$"),
		"/sys/config/auditing/request-headers":          regexp.MustCompile("^/sys/config/auditing/request-headers$"),
		"/sys/config/auditing/request-headers/{header}": regexp.MustCompile("^/sys/config/auditing/request-headers/.+$"),
		"/sys/config/cors":                              regexp.MustCompile("^/sys/config/cors$"),
		"/sys/config/ui/headers/":                       regexp.MustCompile("^/sys/config/ui/headers/$"),
		"/sys/config/ui/headers/{header}":               regexp.MustCompile("^/sys/config/ui/headers/.+$"),
		"/sys/leases":                                   regexp.MustCompile("^/sys/leases$"),
		"/sys/leases/lookup/":                           regexp.MustCompile("^/sys/leases/lookup/$"),
		"/sys/leases/lookup/{prefix}":                   regexp.MustCompile("^/sys/leases/lookup/.+$"),
		"/sys/leases/revoke-force/{prefix}":             regexp.MustCompile("^/sys/leases/revoke-force/.+$"),
		"/sys/leases/revoke-prefix/{prefix}":            regexp.MustCompile("^/sys/leases/revoke-prefix/.+$"),
		"/sys/plugins/catalog/{name}":                   regexp.MustCompile("^/sys/plugins/catalog/.+$"),
		"/sys/plugins/catalog/{type}":                   regexp.MustCompile("^/sys/plugins/catalog/.+$"),
		"/sys/plugins/catalog/{type}/{name}":            regexp.MustCompile("^/sys/plugins/catalog/.+/.+$"),
		"/sys/raw":                                      regexp.MustCompile("^/sys/raw$"),
		"/sys/raw/{path}":                               regexp.MustCompile("^/sys/raw/.+$"),
		"/sys/remount":                                  regexp.MustCompile("^/sys/remount$"),
		"/sys/revoke-force/{prefix}":                    regexp.MustCompile("^/sys/revoke-force/.+$"),
		"/sys/revoke-prefix/{prefix}":                   regexp.MustCompile("^/sys/revoke-prefix/.+$"),
		"/sys/rotate":                                   regexp.MustCompile("^/sys/rotate$"),
	}
}
