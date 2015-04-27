package aws

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// RestClient is the underlying client for REST-JSON and REST-XML APIs.
type RestClient struct {
	Context    Context
	Client     *http.Client
	Endpoint   string
	APIVersion string
}

// Whether the byte value can be sent without escaping in AWS URLs
var noEscape [256]bool

// Initialise noEscape
func init() {
	for i := range noEscape {
		// Amazon expects every character except these escaped
		noEscape[i] = (i >= 'A' && i <= 'Z') ||
			(i >= 'a' && i <= 'z') ||
			(i >= '0' && i <= '9') ||
			i == '-' ||
			i == '.' ||
			i == '/' ||
			i == ':' ||
			i == '_' ||
			i == '~'
	}
}

// EscapePath escapes part of a URL path in Amazon style
func EscapePath(path string) string {
	var buf bytes.Buffer
	for i := 0; i < len(path); i++ {
		c := path[i]
		if noEscape[c] {
			buf.WriteByte(c)
		} else {
			buf.WriteByte('%')
			buf.WriteString(strings.ToUpper(strconv.FormatUint(uint64(c), 16)))
		}
	}
	return buf.String()
}

// Do sends an HTTP request and returns an HTTP response, following policy
// (e.g. redirects, cookies, auth) as configured on the client.
func (c *RestClient) Do(req *http.Request) (*http.Response, error) {
	// Set the form for the URL
	req.URL.Opaque = EscapePath(req.URL.Path)
	req.Header.Set("User-Agent", "aws-go")
	if err := c.Context.sign(req); err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if len(bodyBytes) == 0 {
			return nil, APIError{
				StatusCode: resp.StatusCode,
				Message:    resp.Status,
			}
		}
		var restErr restError
		switch resp.Header.Get("Content-Type") {
		case "application/json":
			if err := json.Unmarshal(bodyBytes, &restErr); err != nil {
				return nil, err
			}
			return nil, restErr.Err(resp.StatusCode)
		case "application/xml", "text/xml":
			// AWS XML error documents can have a couple of different formats.
			// Try each before returning a decode error.
			var wrappedErr restErrorResponse
			if err := xml.Unmarshal(bodyBytes, &wrappedErr); err == nil {
				return nil, wrappedErr.Error.Err(resp.StatusCode)
			}
			if err := xml.Unmarshal(bodyBytes, &restErr); err != nil {
				return nil, err
			}
			return nil, restErr.Err(resp.StatusCode)
		default:
			return nil, APIError{
				StatusCode: resp.StatusCode,
				Message:    string(bodyBytes),
			}
		}
	}

	return resp, nil
}

type restErrorResponse struct {
	XMLName xml.Name `xml:"ErrorResponse",json:"-"`
	Error   restError
}

type restError struct {
	XMLName    xml.Name `xml:"Error",json:"-"`
	Code       string
	BucketName string
	Message    string
	RequestID  string
	HostID     string
}

func (e restError) Err(StatusCode int) error {
	return APIError{
		StatusCode: StatusCode,
		Code:       e.Code,
		Message:    e.Message,
		RequestID:  e.RequestID,
		HostID:     e.HostID,
		Specifics: map[string]string{
			"BucketName": e.BucketName,
		},
	}
}
