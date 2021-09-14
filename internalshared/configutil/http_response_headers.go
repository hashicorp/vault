package configutil

import (
    "fmt"
    "net/textproto"
    "strconv"
    "strings"
)

var DefaultHeaderNames = []string {
    "Content-Security-Policy",
    "X-XSS-Protection",
    "X-Frame-Options",
    "X-Content-Type-Options",
    "Strict-Transport-Security",
    "Content-Type",
}

var ValidCustomStatusCodeCollection = []string {
    "default",
    "1xx",
    "2xx",
    "3xx",
    "4xx",
    "5xx",
}

const (
    contentSecurityPolicy = "default-src 'none';   connect-src 'self';  img-src 'self' data:; script-src 'self'; style-src 'unsafe-inline' 'self'; form-action  'none'; frame-ancestors 'none'; font-src 'self'"
    xXssProtection = "1; mode=block"
    xFrameOptions = "Deny"
    xContentTypeOptions = "nosniff"
    strictTransportSecurity = "max-age=31536000; includeSubDomains"
    contentType = "application/json"
)

func GetDefaultHeaderValue(h string) string {
    switch h {
    case "Content-Security-Policy":
        return contentSecurityPolicy
    case "X-XSS-Protection":
        return xXssProtection
    case "X-Frame-Options":
        return xFrameOptions
    case "X-Content-Type-Options":
        return xContentTypeOptions
    case "Strict-Transport-Security":
        return strictTransportSecurity
    case "Content-Type":
        return contentType
    default:
        return ""
    }
}

func setDefaultResponseHeaders(c map[string]string) map[string]string {
    defaults := make(map[string]string)
    // adding all parsed default headers
    for k, v := range c {
        defaults[k] = v
    }

    // setting all default headers that are not included in the config
    // file under the "default" category
    for _, hn := range DefaultHeaderNames {
        if _, ok := c[hn]; ok {
            continue
        }
        hv := GetDefaultHeaderValue(hn)
        if hv != "" {
            defaults[hn] = hv
        }
    }

    return defaults
}

func ParseCustomResponseHeaders(r interface{}) (map[string]map[string]string, error) {
    if _, ok := r.([]map[string]interface{}); !ok {
        return nil, fmt.Errorf("response headers were not configured correctly. please make sure they're in a map")
    }

    customResponseHeader := r.([]map[string]interface{})
    h := make(map[string]map[string]string)

    for _, crh := range customResponseHeader {
        for statusCode, responseHeader := range crh {
            if _, ok := responseHeader.([]map[string]interface{}); !ok {
                return nil, fmt.Errorf("response headers were not configured correctly. please make sure they're in a map")
            }

            if !IsValidStatusCode(statusCode) {
                return nil, fmt.Errorf("invalid status code found in the config file: %v", statusCode)
            }

            hvl := responseHeader.([]map[string]interface{})
            if len(hvl) != 1 {
                return nil, fmt.Errorf("invalid number of response headers exist")
            }
            hvm := hvl[0]
            hv, err := parseHeaders(hvm)
            if err != nil {
                return nil, err
            }

            h[statusCode] = hv
        }
    }

    // setting default custom headers
    de := h["default"]
    h["default"] = setDefaultResponseHeaders(de)

	return h, nil
}

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
        hn := textproto.CanonicalMIMEHeaderKey(k)
        // parsing header values
        s, err := parseHeaderValues(v)
        if err != nil {
            return nil, err
        }
        hvMap[hn] = s
    }
    return hvMap, nil
}

func parseHeaderValues(h interface{}) (string, error) {
    var sl []string
    if _, ok := h.([]interface{}); !ok {
        return "", fmt.Errorf("headers must be given in a list of strings")
    }
    vli := h.([]interface{})
    for _, vh := range vli {
        if vh.(string) == "" {
           continue
        }
        sl = append(sl, vh.(string))
    }
    s := strings.Join(sl, "; ")

    return s, nil
}