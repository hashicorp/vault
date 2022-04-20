package api

import (
	"fmt"
	"net/http"
	"strings"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const (
	ErrOutputStringRequest = "output a string, please"
)

var LastOutputStringError *OutputStringError

type OutputStringError struct {
	*retryablehttp.Request
	TLSSkipVerify              bool
	ClientCACert, ClientCAPath string
	ClientCert, ClientKey      string
	parsingError               error
	finalCurlString            string
}

func (d *OutputStringError) Error() string {
	if d.finalCurlString == "" {
		d.buildCurlString()
		if d.parsingError != nil {
			return d.parsingError.Error()
		}
	}

	return ErrOutputStringRequest
}

func (d *OutputStringError) buildCurlString() {
	body, err := d.Request.BodyBytes()
	if err != nil {
		d.parsingError = err
		return
	}

	// Build cURL string
	d.finalCurlString = "curl "
	if d.TLSSkipVerify {
		d.finalCurlString += "--insecure "
	}
	if d.Request.Method != http.MethodGet {
		d.finalCurlString = fmt.Sprintf("%s-X %s ", d.finalCurlString, d.Request.Method)
	}
	if d.ClientCACert != "" {
		clientCACert := strings.Replace(d.ClientCACert, "'", "'\"'\"'", -1)
		d.finalCurlString = fmt.Sprintf("%s--cacert '%s' ", d.finalCurlString, clientCACert)
	}
	if d.ClientCAPath != "" {
		clientCAPath := strings.Replace(d.ClientCAPath, "'", "'\"'\"'", -1)
		d.finalCurlString = fmt.Sprintf("%s--capath '%s' ", d.finalCurlString, clientCAPath)
	}
	if d.ClientCert != "" {
		clientCert := strings.Replace(d.ClientCert, "'", "'\"'\"'", -1)
		d.finalCurlString = fmt.Sprintf("%s--cert '%s' ", d.finalCurlString, clientCert)
	}
	if d.ClientKey != "" {
		clientKey := strings.Replace(d.ClientKey, "'", "'\"'\"'", -1)
		d.finalCurlString = fmt.Sprintf("%s--key '%s' ", d.finalCurlString, clientKey)
	}
	for k, v := range d.Request.Header {
		for _, h := range v {
			if strings.ToLower(k) == "x-vault-token" {
				h = `$(vault print token)`
			}
			d.finalCurlString = fmt.Sprintf("%s-H \"%s: %s\" ", d.finalCurlString, k, h)
		}
	}

	if len(body) > 0 {
		// We need to escape single quotes since that's what we're using to
		// quote the body
		escapedBody := strings.Replace(string(body), "'", "'\"'\"'", -1)
		d.finalCurlString = fmt.Sprintf("%s-d '%s' ", d.finalCurlString, escapedBody)
	}

	d.finalCurlString = fmt.Sprintf("%s%s", d.finalCurlString, d.Request.URL.String())
}

func (d *OutputStringError) CurlString() string {
	if d.finalCurlString == "" {
		d.buildCurlString()
	}
	return d.finalCurlString
}
