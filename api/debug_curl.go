package api

import (
	"net/http"
)

const (
	ErrDebugCurl = "output cURL, please:"
)

var (
	LastDebugCurlError *DebugCurlError
)

type DebugCurlError struct {
	*http.Request
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
	d.parsedString = "<goes here>"
}

func (d *DebugCurlError) CurlString() string {
	if d.parsedString == "" {
		d.parseRequest()
	}
	return d.parsedString
}
