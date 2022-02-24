package api

import (
	"fmt"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const (
	ErrOutputPolicyRequest = "output a policy, please"
)

var LastOutputPolicyError *OutputPolicyError

type OutputPolicyError struct {
	*retryablehttp.Request
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

func (d *OutputPolicyError) parseRequest() {
	body, err := d.Request.BodyBytes()
	if err != nil {
		d.parsingError = err
		return
	}

	// Build policy document
	d.parsedHCLString = "path "

	path := d.Request.URL.String()

	if d.Request.Method == "GET" {
		d.parsedHCLString
	}
	// if d.TLSSkipVerify {
	// 	d.parsedCurlString += "--insecure "
	// }
	// if d.Request.Method != "GET" {
	// 	d.parsedCurlString = fmt.Sprintf("%s-X %s ", d.parsedCurlString, d.Request.Method)
	// }
	// if d.ClientCACert != "" {
	// 	clientCACert := strings.Replace(d.ClientCACert, "'", "'\"'\"'", -1)
	// 	d.parsedCurlString = fmt.Sprintf("%s--cacert '%s' ", d.parsedCurlString, clientCACert)
	// }
	// if d.ClientCAPath != "" {
	// 	clientCAPath := strings.Replace(d.ClientCAPath, "'", "'\"'\"'", -1)
	// 	d.parsedCurlString = fmt.Sprintf("%s--capath '%s' ", d.parsedCurlString, clientCAPath)
	// }
	// if d.ClientCert != "" {
	// 	clientCert := strings.Replace(d.ClientCert, "'", "'\"'\"'", -1)
	// 	d.parsedCurlString = fmt.Sprintf("%s--cert '%s' ", d.parsedCurlString, clientCert)
	// }
	// if d.ClientKey != "" {
	// 	clientKey := strings.Replace(d.ClientKey, "'", "'\"'\"'", -1)
	// 	d.parsedCurlString = fmt.Sprintf("%s--key '%s' ", d.parsedCurlString, clientKey)
	// }
	// for k, v := range d.Request.Header {
	// 	for _, h := range v {
	// 		if strings.ToLower(k) == "x-vault-token" {
	// 			h = `$(vault print token)`
	// 		}
	// 		d.parsedCurlString = fmt.Sprintf("%s-H \"%s: %s\" ", d.parsedCurlString, k, h)
	// 	}
	// }

	// if len(body) > 0 {
	// 	// We need to escape single quotes since that's what we're using to
	// 	// quote the body
	// 	escapedBody := strings.Replace(string(body), "'", "'\"'\"'", -1)
	// 	d.parsedCurlString = fmt.Sprintf("%s-d '%s' ", d.parsedCurlString, escapedBody)
	// }

	d.parsedHCLString = fmt.Sprintf("%s%s", d.parsedHCLString, d.Request.URL.String())
}

func (d *OutputPolicyError) HCLString() string {
	if d.parsedHCLString == "" {
		d.parseRequest()
	}
	return d.parsedHCLString
}
