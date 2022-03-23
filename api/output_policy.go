package api

import (
	"fmt"
	"strings"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const (
	ErrOutputPolicyRequest = "output policy request"
)

var LastOutputPolicyError *OutputPolicyError

type OutputPolicyError struct {
	*retryablehttp.Request
	VaultAddress    string
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
	apiAddrPrefix := fmt.Sprintf("%sv1/", d.VaultAddress)
	path := strings.Trim(url, apiAddrPrefix)

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
