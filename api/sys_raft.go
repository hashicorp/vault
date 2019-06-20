package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"golang.org/x/net/http2"
)

// RaftJoinResponse represents the response of the raft join API
type RaftJoinResponse struct {
	Joined bool `json:"joined"`
}

// RaftJoinRequest represents the parameters consumed by the raft join API
type RaftJoinRequest struct {
	LeaderAddr string `json:"leader_api_addr"`
	CACert     string `json:"ca_cert":`
	Retry      bool   `json:"retry"`
}

// RaftJoin adds the node from which this call is invoked from to the raft
// cluster represented by the leader address in the parameter.
func (c *Sys) RaftJoin(opts *RaftJoinRequest) (*RaftJoinResponse, error) {
	r := c.c.NewRequest("POST", "/v1/sys/storage/raft/join")

	if err := r.SetJSONBody(opts); err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RaftJoinResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// RaftSnapshot invokes the API that takes the snapshot of the raft cluster and
// writes it to the supplied io.Writer.
func (c *Sys) RaftSnapshot(snapWriter io.Writer) error {
	r := c.c.NewRequest("GET", "/v1/sys/storage/raft/snapshot")
	r.URL.RawQuery = r.Params.Encode()

	req, err := http.NewRequest(http.MethodGet, r.URL.RequestURI(), nil)
	if err != nil {
		return err
	}

	req.URL.User = r.URL.User
	req.URL.Scheme = r.URL.Scheme
	req.URL.Host = r.URL.Host
	req.Host = r.URL.Host

	if r.Headers != nil {
		for header, vals := range r.Headers {
			for _, val := range vals {
				req.Header.Add(header, val)
			}
		}
	}

	if len(r.ClientToken) != 0 {
		req.Header.Set(consts.AuthHeaderName, r.ClientToken)
	}

	if len(r.WrapTTL) != 0 {
		req.Header.Set("X-Vault-Wrap-TTL", r.WrapTTL)
	}

	if len(r.MFAHeaderVals) != 0 {
		for _, mfaHeaderVal := range r.MFAHeaderVals {
			req.Header.Add("X-Vault-MFA", mfaHeaderVal)
		}
	}

	if r.PolicyOverride {
		req.Header.Set("X-Vault-Policy-Override", "true")
	}

	config := &Config{
		Address: "https://127.0.0.1:8200",
		// Avoiding the use of retryablehttp which reads the response body. The
		// response body of the snapshot API is expected to be huge and it is
		// not desired that the response is held in memory.
		HttpClient: &http.Client{
			Transport: cleanhttp.DefaultPooledTransport(),
			// Giving more time to retrieve all the data
			Timeout: time.Minute * 10,
		},
	}

	transport := config.HttpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		return err
	}
	if err := config.ReadEnvironment(); err != nil {
		return err
	}

	// Avoiding the use of RawRequestWithContext which reads the response body
	// to determine if the body contains error message.
	resp, err := config.HttpClient.Do(req)

	// By not parsing the body and checking only status code here, the error
	// message will land in the snapshot file, but not without the API erroring
	// out.
	if !(resp.StatusCode >= 200 && resp.StatusCode < 400) && resp.StatusCode != 429 {
		return fmt.Errorf("http error code: %d", resp.StatusCode)
	}
	if err != nil {
		return err
	}

	_, err = io.Copy(snapWriter, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// RaftSnapshotRestore reads the snapshot from the io.Reader and installs that
// snapshot, returning the cluster to the state defined by it.
func (c *Sys) RaftSnapshotRestore(snapReader io.Reader, force bool) error {
	path := "/v1/sys/storage/raft/snapshot"
	if force {
		path = "/v1/sys/storage/raft/snapshot-force"
	}
	r := c.c.NewRequest("POST", path)

	r.Body = snapReader

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
