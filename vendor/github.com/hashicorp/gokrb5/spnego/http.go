package spnego

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/hashicorp/gokrb5/client"
	"github.com/hashicorp/gokrb5/gssapi"
	"github.com/hashicorp/gokrb5/keytab"
	"github.com/hashicorp/gokrb5/krberror"
	"github.com/hashicorp/gokrb5/service"
	"github.com/hashicorp/gokrb5/types"
	"gopkg.in/jcmturner/goidentity.v3"
)

// Client side functionality //

// Client will negotiate authentication with a server using SPNEGO.
type Client struct {
	*http.Client
	krb5Client *client.Client
	spn        string
	reqs       []*http.Request
}

type redirectErr struct {
	reqTarget *http.Request
}

func (e redirectErr) Error() string {
	return fmt.Sprintf("redirect to %v", e.reqTarget.URL)
}

type teeReadCloser struct {
	io.Reader
	io.Closer
}

// NewClient returns an SPNEGO enabled HTTP client.
func NewClient(krb5Cl *client.Client, httpCl *http.Client, spn string) *Client {
	if httpCl == nil {
		httpCl = http.DefaultClient
	}
	// Add a cookie jar if there isn't one
	if httpCl.Jar == nil {
		httpCl.Jar, _ = cookiejar.New(nil)
	}
	// Add a CheckRedirect function that will execute any functional already defined and then error with a redirectErr
	f := httpCl.CheckRedirect
	httpCl.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if f != nil {
			err := f(req, via)
			if err != nil {
				return err
			}
		}
		return redirectErr{reqTarget: req}
	}
	return &Client{
		Client:     httpCl,
		krb5Client: krb5Cl,
		spn:        spn,
	}
}

// Do is the SPNEGO enabled HTTP client's equivalent of the http.Client's Do method.
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	var body bytes.Buffer
	if req.Body != nil {
		// Use a tee reader to capture any body sent in case we have to replay it again
		teeR := io.TeeReader(req.Body, &body)
		teeRC := teeReadCloser{teeR, req.Body}
		req.Body = teeRC
	}
	resp, err = c.Client.Do(req)
	if err != nil {
		if ue, ok := err.(*url.Error); ok {
			if e, ok := ue.Err.(redirectErr); ok {
				// Picked up a redirect
				e.reqTarget.Header.Del(HTTPHeaderAuthRequest)
				c.reqs = append(c.reqs, e.reqTarget)
				if len(c.reqs) >= 10 {
					return resp, errors.New("stopped after 10 redirects")
				}
				if req.Body != nil {
					// Refresh the body reader so the body can be sent again
					e.reqTarget.Body = ioutil.NopCloser(&body)
				}
				return c.Do(e.reqTarget)
			}
		}
		return resp, err
	}
	if respUnauthorizedNegotiate(resp) {
		err := SetSPNEGOHeader(c.krb5Client, req, c.spn)
		if err != nil {
			return resp, err
		}
		if req.Body != nil {
			// Refresh the body reader so the body can be sent again
			req.Body = ioutil.NopCloser(&body)
		}
		return c.Do(req)
	}
	return resp, err
}

// Get is the SPNEGO enabled HTTP client's equivalent of the http.Client's Get method.
func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post is the SPNEGO enabled HTTP client's equivalent of the http.Client's Post method.
func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

// PostForm is the SPNEGO enabled HTTP client's equivalent of the http.Client's PostForm method.
func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// Head is the SPNEGO enabled HTTP client's equivalent of the http.Client's Head method.
func (c *Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func respUnauthorizedNegotiate(resp *http.Response) bool {
	if resp.StatusCode == http.StatusUnauthorized {
		if resp.Header.Get(HTTPHeaderAuthResponse) == HTTPHeaderAuthResponseValueKey {
			return true
		}
	}
	return false
}

// SetSPNEGOHeader gets the service ticket and sets it as the SPNEGO authorization header on HTTP request object.
// To auto generate the SPN from the request object pass a null string "".
func SetSPNEGOHeader(cl *client.Client, r *http.Request, spn string) error {
	if spn == "" {
		h := strings.TrimSuffix(strings.SplitN(r.URL.Host, ":", 2)[0], ".")
		name, err := net.LookupCNAME(h)
		if err == nil {
			// Underlyng canonical name should be used for SPN
			h = strings.TrimSuffix(name, ".")
		}
		spn = "HTTP/" + h
		r.Host = h
	}
	cl.Log("using SPN %s", spn)
	s := SPNEGOClient(cl, spn)
	err := s.AcquireCred()
	if err != nil {
		return fmt.Errorf("could not acquire client credential: %v", err)
	}
	st, err := s.InitSecContext()
	if err != nil {
		return fmt.Errorf("could not initialize context: %v", err)
	}
	nb, err := st.Marshal()
	if err != nil {
		return krberror.Errorf(err, krberror.EncodingError, "could not marshal SPNEGO")
	}
	hs := "Negotiate " + base64.StdEncoding.EncodeToString(nb)
	r.Header.Set(HTTPHeaderAuthRequest, hs)
	return nil
}

// Service side functionality //

type ctxKey string

const (
	// spnegoNegTokenRespKRBAcceptCompleted - The response on successful authentication always has this header. Capturing as const so we don't have marshaling and encoding overhead.
	spnegoNegTokenRespKRBAcceptCompleted = "Negotiate oRQwEqADCgEAoQsGCSqGSIb3EgECAg=="
	// spnegoNegTokenRespReject - The response on a failed authentication always has this rejection header. Capturing as const so we don't have marshaling and encoding overhead.
	spnegoNegTokenRespReject = "Negotiate oQcwBaADCgEC"
	// spnegoNegTokenRespIncompleteKRB5 - Response token specifying incomplete context and KRB5 as the supported mechtype.
	spnegoNegTokenRespIncompleteKRB5 = "Negotiate oRQwEqADCgEBoQsGCSqGSIb3EgECAg=="
	// CTXKeyAuthenticated is the request context key holding a boolean indicating if the request has been authenticated.
	CTXKeyAuthenticated ctxKey = "github.com/hashicorp/gokrb5/CTXKeyAuthenticated"
	// CTXKeyCredentials is the request context key holding the credentials gopkg.in/jcmturner/goidentity.v2/Identity object.
	CTXKeyCredentials ctxKey = "github.com/hashicorp/gokrb5/CTXKeyCredentials"
	// HTTPHeaderAuthRequest is the header that will hold authn/z information.
	HTTPHeaderAuthRequest = "Authorization"
	// HTTPHeaderAuthResponse is the header that will hold SPNEGO data from the server.
	HTTPHeaderAuthResponse = "WWW-Authenticate"
	// HTTPHeaderAuthResponseValueKey is the key in the auth header for SPNEGO.
	HTTPHeaderAuthResponseValueKey = "Negotiate"
	// UnauthorizedMsg is the message returned in the body when authentication fails.
	UnauthorizedMsg = "Unauthorised.\n"
)

// SPNEGOKRB5Authenticate is a Kerberos SPNEGO authentication HTTP handler wrapper.
func SPNEGOKRB5Authenticate(inner http.Handler, kt *keytab.Keytab, settings ...func(*service.Settings)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the auth header
		s := strings.SplitN(r.Header.Get(HTTPHeaderAuthRequest), " ", 2)
		if len(s) != 2 || s[0] != HTTPHeaderAuthResponseValueKey {
			// No Authorization header set so return 401 with WWW-Authenticate Negotiate header
			w.Header().Set(HTTPHeaderAuthResponse, HTTPHeaderAuthResponseValueKey)
			http.Error(w, UnauthorizedMsg, http.StatusUnauthorized)
			return
		}

		// Set up the SPNEGO GSS-API mechanism
		var spnego *SPNEGO
		h, err := types.GetHostAddress(r.RemoteAddr)
		if err == nil {
			// put in this order so that if the user provides a ClientAddress it will override the one here.
			o := append([]func(*service.Settings){service.ClientAddress(h)}, settings...)
			spnego = SPNEGOService(kt, o...)
		} else {
			spnego = SPNEGOService(kt, settings...)
			spnego.Log("%s - SPNEGO could not parse client address: %v", r.RemoteAddr, err)
		}

		// Decode the header into an SPNEGO context token
		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			spnegoNegotiateKRB5MechType(spnego, w, "%s - SPNEGO error in base64 decoding negotiation header: %v", r.RemoteAddr, err)
			return
		}
		var st SPNEGOToken
		err = st.Unmarshal(b)
		if err != nil {
			spnegoNegotiateKRB5MechType(spnego, w, "%s - SPNEGO error in unmarshaling SPNEGO token: %v", r.RemoteAddr, err)
			return
		}

		// Validate the context token
		authed, ctx, status := spnego.AcceptSecContext(&st)
		if status.Code != gssapi.StatusComplete && status.Code != gssapi.StatusContinueNeeded {
			spnegoResponseReject(spnego, w, "%s - SPNEGO validation error: %v", r.RemoteAddr, status)
			return
		}
		if status.Code == gssapi.StatusContinueNeeded {
			spnegoNegotiateKRB5MechType(spnego, w, "%s - SPNEGO GSS-API continue needed", r.RemoteAddr)
			return
		}
		if authed {
			id := ctx.Value(CTXKeyCredentials).(goidentity.Identity)
			requestCtx := r.Context()
			requestCtx = context.WithValue(requestCtx, CTXKeyCredentials, id)
			requestCtx = context.WithValue(requestCtx, CTXKeyAuthenticated, ctx.Value(CTXKeyAuthenticated))
			spnegoResponseAcceptCompleted(spnego, w, "%s %s@%s - SPNEGO authentication succeeded", r.RemoteAddr, id.UserName(), id.Domain())
			inner.ServeHTTP(w, r.WithContext(requestCtx))
		} else {
			spnegoResponseReject(spnego, w, "%s - SPNEGO Kerberos authentication failed", r.RemoteAddr)
			return
		}
	})
}

func spnegoNegotiateKRB5MechType(s *SPNEGO, w http.ResponseWriter, format string, v ...interface{}) {
	s.Log(format, v...)
	w.Header().Set(HTTPHeaderAuthResponse, spnegoNegTokenRespIncompleteKRB5)
	http.Error(w, UnauthorizedMsg, http.StatusUnauthorized)
}

func spnegoResponseReject(s *SPNEGO, w http.ResponseWriter, format string, v ...interface{}) {
	s.Log(format, v...)
	w.Header().Set(HTTPHeaderAuthResponse, spnegoNegTokenRespReject)
	http.Error(w, UnauthorizedMsg, http.StatusUnauthorized)
}

func spnegoResponseAcceptCompleted(s *SPNEGO, w http.ResponseWriter, format string, v ...interface{}) {
	s.Log(format, v...)
	w.Header().Set(HTTPHeaderAuthResponse, spnegoNegTokenRespKRBAcceptCompleted)
}
