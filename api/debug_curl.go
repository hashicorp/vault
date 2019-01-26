package api

import (
	"fmt"
	"strings"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const (
	ErrDebugCurl = "output cURL, please:"
)

var (
	LastDebugCurlError *DebugCurlError
)

type DebugCurlError struct {
	*retryablehttp.Request
	parsingError error
	parsedString string
}

func (d *DebugCurlError) Error() string {
	if d.parsedString == "" {
		d.parseRequest()
		if d.parsingError != nil {
			return d.parsingError.Error()
		}
	}

	return ErrDebugCurl
}

func (d *DebugCurlError) parseRequest() {
	var err error
	d.parsedString = "curl "
	d.parsedString = fmt.Sprintf("%s-X %s ", d.parsedString, d.Request.Method)
	for k, v := range d.Request.Header {
		for _, h := range v {
			d.parsedString = fmt.Sprintf("%s-H \"%s: %s\" ", d.parsedString, k, h)
		}
	}

	body, err := d.Request.BodyBytes()
	if err != nil {
		d.parsingError = err
		return
	}
	if len(body) > 0 {
		// We need to escape single quotes since that's what we're using to
		// quote the body
		escapedBody := strings.Replace(string(body), "'", "'\"'\"'", -1)
		d.parsedString = fmt.Sprintf("%s-d '%s' ", d.parsedString, escapedBody)
	}

	d.parsedString = fmt.Sprintf("%s%s", d.parsedString, d.Request.URL.String())
}

func (d *DebugCurlError) CurlString() string {
	if d.parsedString == "" {
		d.parseRequest()
	}
	return d.parsedString
}
