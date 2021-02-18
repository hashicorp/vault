package rabbithole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	// URI of a RabbitMQ node to use, not including the path, e.g. http://127.0.0.1:15672.
	Endpoint string
	// Username to use. This RabbitMQ user must have the "management" tag.
	Username string
	// Password to use.
	Password  string
	host      string
	transport *http.Transport
	timeout   time.Duration
}

func NewClient(uri string, username string, password string) (me *Client, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	me = &Client{
		Endpoint: uri,
		host:     u.Host,
		Username: username,
		Password: password,
	}

	return me, nil
}

// Creates a client with a transport; it is up to the developer to make that layer secure.
func NewTLSClient(uri string, username string, password string, transport *http.Transport) (me *Client, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	me = &Client{
		Endpoint:  uri,
		host:      u.Host,
		Username:  username,
		Password:  password,
		transport: transport,
	}

	return me, nil
}

//SetTransport changes the Transport Layer that the Client will use.
func (c *Client) SetTransport(transport *http.Transport) {
	c.transport = transport
}

// SetTimeout changes the HTTP timeout that the Client will use.
// By default there is no timeout.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func newGETRequest(client *Client, path string) (*http.Request, error) {
	s := client.Endpoint + "/api/" + path
	req, err := http.NewRequest("GET", s, nil)

	req.Close = true
	req.SetBasicAuth(client.Username, client.Password)

	return req, err
}

func newGETRequestWithParameters(client *Client, path string, qs url.Values) (*http.Request, error) {
	s := client.Endpoint + "/api/" + path + "?" + qs.Encode()

	req, err := http.NewRequest("GET", s, nil)
	req.Close = true
	req.SetBasicAuth(client.Username, client.Password)

	return req, err
}

func newRequestWithBody(client *Client, method string, path string, body []byte) (*http.Request, error) {
	s := client.Endpoint + "/api/" + path

	req, err := http.NewRequest(method, s, bytes.NewReader(body))

	req.Close = true
	req.SetBasicAuth(client.Username, client.Password)

	req.Header.Add("Content-Type", "application/json")

	return req, err
}

func executeRequest(client *Client, req *http.Request) (res *http.Response, err error) {
	httpc := &http.Client{
		Timeout: client.timeout,
	}
	if client.transport != nil {
		httpc.Transport = client.transport
	}
	return httpc.Do(req)
}

func executeAndParseRequest(client *Client, req *http.Request, rec interface{}) (err error) {
	res, err := executeRequest(client, req)
	if err != nil {
		return err
	}
	defer res.Body.Close() // always close body

	if res.StatusCode >= http.StatusBadRequest {
		rme := ErrorResponse{}
		err = json.NewDecoder(res.Body).Decode(&rme)
		if err != nil {
			return fmt.Errorf("Error %d from RabbitMQ: %s", res.StatusCode, err)
		}
		rme.StatusCode = res.StatusCode
		return rme
	}

	err = json.NewDecoder(res.Body).Decode(&rec)
	if err != nil {
		return err
	}

	return nil
}

// This is an ugly hack: we copy relevant bits from
// https://github.com/golang/go/blob/7e2bf952a905f16a17099970392ea17545cdd193/src/net/url/url.go
// because up to Go 1.8 there is no built-in method
// (and url.QueryEscape isn't suitable since it encodes
// spaces as + and not %20).
//
// See https://github.com/golang/go/issues/13737,
//     https://github.com/golang/go/commit/7e2bf952a905f16a17099970392ea17545cdd193

// PathEscape escapes the string so it can be safely placed
// inside a URL path segment.
func PathEscape(s string) string {
	return escape(s, encodePathSegment)
}

type encoding int

const (
	encodePath encoding = 1 + iota
	encodePathSegment
	encodeHost
	encodeZone
	encodeUserPassword
	encodeQueryComponent
	encodeFragment
)

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent:
			t[j] = '+'
			j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
//
// Please be informed that for now shouldEscape does not check all
// reserved characters correctly. See golang.org/issue/5684.
func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == encodeHost || mode == encodeZone {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't use %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case encodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last three as well. That leaves only ? to escape.
			return c == '?'

		case encodePathSegment: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments.
			return c == '/' || c == ';' || c == ',' || c == '?'

		case encodeUserPassword: // §3.2.1
			// The RFC allows ';', ':', '&', '=', '+', '$', and ',' in
			// userinfo, so we must escape only '@', '/', and '?'.
			// The parsing of userinfo treats ':' as special so we must escape
			// that too.
			return c == '@' || c == '/' || c == '?' || c == ':'

		case encodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	// Everything else must be escaped.
	return true
}
