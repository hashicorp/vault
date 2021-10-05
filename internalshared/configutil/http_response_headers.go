package configutil

import (
	"fmt"
	"net/textproto"
	"strconv"
	"strings"
)

var ValidCustomStatusCodeCollection = []string{
	"default",
	"1xx",
	"2xx",
	"3xx",
	"4xx",
	"5xx",
}

const StrictTransportSecurity = "max-age=31536000; includeSubDomains"

// ParseCustomResponseHeaders takes a raw config values for the
// "custom_response_headers". It makes sure the config entry is passed in
// as a map of status code to a map of header name and header values. It
// verifies the validity of the status codes, and header values. It also
// adds the default headers values.
func ParseCustomResponseHeaders(responseHeaders interface{}) (map[string]map[string]string, error) {
	h := make(map[string]map[string]string)
	// if r is nil, we still should set the default custom headers
	if responseHeaders == nil {
		h["default"] = map[string]string{"Strict-Transport-Security": StrictTransportSecurity}
		return h, nil
	}

	if _, ok := responseHeaders.([]map[string]interface{}); !ok {
		return nil, fmt.Errorf("response headers were not configured correctly. please make sure they're in a slice of maps")
	}

	customResponseHeader := responseHeaders.([]map[string]interface{})

	for _, crh := range customResponseHeader {
		for statusCode, responseHeader := range crh {
			if _, ok := responseHeader.([]map[string]interface{}); !ok {
				return nil, fmt.Errorf("response headers were not configured correctly. please make sure they're in a slice of maps")
			}

			if !IsValidStatusCode(statusCode) {
				return nil, fmt.Errorf("invalid status code found in the server configuration: %v", statusCode)
			}

			headerValList := responseHeader.([]map[string]interface{})
			if len(headerValList) != 1 {
				return nil, fmt.Errorf("invalid number of response headers exist")
			}
			headerValMap := headerValList[0]
			headerVal, err := parseHeaders(headerValMap)
			if err != nil {
				return nil, err
			}

			h[statusCode] = headerVal
		}
	}

	// setting Strict-Transport-Security as a default header
	if h["default"] == nil {
		h["default"] = make(map[string]string)
	}
	if _, ok := h["default"]["Strict-Transport-Security"]; !ok {
		h["default"]["Strict-Transport-Security"] = StrictTransportSecurity
	}

	return h, nil
}

// IsValidStatusCodeCollection checks if the given status code is as expected
func IsValidStatusCodeCollection(sc string) bool {
	for _, v := range ValidCustomStatusCodeCollection {
		if sc == v {
			return true
		}
	}

	return false
}

// IsValidStatusCode checking for status codes outside the boundary
func IsValidStatusCode(sc string) bool {
	if IsValidStatusCodeCollection(sc) {
		return true
	}

	i, err := strconv.Atoi(sc)
	if err != nil {
		return false
	}

	if i >= 600 || i < 100 {
		return false
	}

	return true
}

func parseHeaders(in map[string]interface{}) (map[string]string, error) {
	hvMap := make(map[string]string)
	for k, v := range in {
		// parsing header name
		headerName := textproto.CanonicalMIMEHeaderKey(k)
		// parsing header values
		s, err := parseHeaderValues(v)
		if err != nil {
			return nil, err
		}
		hvMap[headerName] = s
	}
	return hvMap, nil
}

func parseHeaderValues(header interface{}) (string, error) {
	var sl []string
	if _, ok := header.([]interface{}); !ok {
		return "", fmt.Errorf("headers must be given in a list of strings")
	}
	headerValList := header.([]interface{})
	for _, vh := range headerValList {
		if _, ok := vh.(string); !ok {
			return "", fmt.Errorf("found a non-string header value: %v", vh)
		}
		headerVal := strings.TrimSpace(vh.(string))
		if headerVal == "" {
			continue
		}
		sl = append(sl, headerVal)

	}
	s := strings.Join(sl, "; ")

	return s, nil
}
