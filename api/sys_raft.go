package api

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/mapstructure"
)

var ErrIncompleteSnapshot = errors.New("incomplete snapshot, unable to read SHA256SUMS.sealed file")

// RaftJoinResponse represents the response of the raft join API
type RaftJoinResponse struct {
	Joined bool `json:"joined"`
}

// RaftJoinRequest represents the parameters consumed by the raft join API
type RaftJoinRequest struct {
	AutoJoin         string `json:"auto_join"`
	AutoJoinScheme   string `json:"auto_join_scheme"`
	AutoJoinPort     uint   `json:"auto_join_port"`
	LeaderAPIAddr    string `json:"leader_api_addr"`
	LeaderCACert     string `json:"leader_ca_cert"`
	LeaderClientCert string `json:"leader_client_cert"`
	LeaderClientKey  string `json:"leader_client_key"`
	Retry            bool   `json:"retry"`
	NonVoter         bool   `json:"non_voter"`
}

// AutopilotConfig is used for querying/setting the Autopilot configuration.
type AutopilotConfig struct {
	CleanupDeadServers             bool          `json:"cleanup_dead_servers" mapstructure:"cleanup_dead_servers"`
	LastContactThreshold           time.Duration `json:"last_contact_threshold" mapstructure:"-"`
	DeadServerLastContactThreshold time.Duration `json:"dead_server_last_contact_threshold" mapstructure:"-"`
	MaxTrailingLogs                uint64        `json:"max_trailing_logs" mapstructure:"max_trailing_logs"`
	MinQuorum                      uint          `json:"min_quorum" mapstructure:"min_quorum"`
	ServerStabilizationTime        time.Duration `json:"server_stabilization_time" mapstructure:"-"`
}

// MarshalJSON makes the autopilot config fields JSON compatible
func (ac *AutopilotConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"cleanup_dead_servers":               ac.CleanupDeadServers,
		"last_contact_threshold":             ac.LastContactThreshold.String(),
		"dead_server_last_contact_threshold": ac.DeadServerLastContactThreshold.String(),
		"max_trailing_logs":                  ac.MaxTrailingLogs,
		"min_quorum":                         ac.MinQuorum,
		"server_stabilization_time":          ac.ServerStabilizationTime.String(),
	})
}

// UnmarshalJSON parses the autopilot config JSON blob
func (ac *AutopilotConfig) UnmarshalJSON(b []byte) error {
	var data interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	conf := data.(map[string]interface{})
	if err = mapstructure.WeakDecode(conf, ac); err != nil {
		return err
	}
	if ac.LastContactThreshold, err = parseutil.ParseDurationSecond(conf["last_contact_threshold"]); err != nil {
		return err
	}
	if ac.DeadServerLastContactThreshold, err = parseutil.ParseDurationSecond(conf["dead_server_last_contact_threshold"]); err != nil {
		return err
	}
	if ac.ServerStabilizationTime, err = parseutil.ParseDurationSecond(conf["server_stabilization_time"]); err != nil {
		return err
	}
	return nil
}

// AutopilotState represents the response of the raft autopilot state API
type AutopilotState struct {
	Healthy          bool                        `mapstructure:"healthy"`
	FailureTolerance int                         `mapstructure:"failure_tolerance"`
	Servers          map[string]*AutopilotServer `mapstructure:"servers"`
	Leader           string                      `mapstructure:"leader"`
	Voters           []string                    `mapstructure:"voters"`
	NonVoters        []string                    `mapstructure:"non_voters"`
}

// AutopilotServer represents the server blocks in the response of the raft
// autopilot state API.
type AutopilotServer struct {
	ID          string            `mapstructure:"id"`
	Name        string            `mapstructure:"name"`
	Address     string            `mapstructure:"address"`
	NodeStatus  string            `mapstructure:"node_status"`
	LastContact string            `mapstructure:"last_contact"`
	LastTerm    uint64            `mapstructure:"last_term"`
	LastIndex   uint64            `mapstructure:"last_index"`
	Healthy     bool              `mapstructure:"healthy"`
	StableSince string            `mapstructure:"stable_since"`
	Status      string            `mapstructure:"status"`
	Meta        map[string]string `mapstructure:"meta"`
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

	// Avoiding the use of RawRequestWithContext which reads the response body
	// to determine if the body contains error message.
	var result *Response
	resp, err := c.c.config.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp == nil {
		return nil
	}

	// Check for a redirect, only allowing for a single redirect
	if resp.StatusCode == 301 || resp.StatusCode == 302 || resp.StatusCode == 307 {
		// Parse the updated location
		respLoc, err := resp.Location()
		if err != nil {
			return err
		}

		// Ensure a protocol downgrade doesn't happen
		if req.URL.Scheme == "https" && respLoc.Scheme != "https" {
			return fmt.Errorf("redirect would cause protocol downgrade")
		}

		// Update the request
		req.URL = respLoc

		// Retry the request
		resp, err = c.c.config.HttpClient.Do(req)
		if err != nil {
			return err
		}
	}

	result = &Response{Response: resp}
	if err := result.Error(); err != nil {
		return err
	}

	// Make sure that the last file in the archive, SHA256SUMS.sealed, is present
	// and non-empty.  This is to catch cases where the snapshot failed midstream,
	// e.g. due to a problem with the seal that prevented encryption of that file.
	var wg sync.WaitGroup
	wg.Add(1)
	var verified bool

	rPipe, wPipe := io.Pipe()
	dup := io.TeeReader(resp.Body, wPipe)
	go func() {
		defer func() {
			io.Copy(ioutil.Discard, rPipe)
			rPipe.Close()
			wg.Done()
		}()

		uncompressed, err := gzip.NewReader(rPipe)
		if err != nil {
			return
		}

		t := tar.NewReader(uncompressed)
		var h *tar.Header
		for {
			h, err = t.Next()
			if err != nil {
				return
			}
			if h.Name != "SHA256SUMS.sealed" {
				continue
			}
			var b []byte
			b, err = ioutil.ReadAll(t)
			if err != nil || len(b) == 0 {
				return
			}
			verified = true
			return
		}
	}()

	// Copy bytes from dup to snapWriter.  This will have a side effect that
	// everything read from dup will be written to wPipe.
	_, err = io.Copy(snapWriter, dup)
	wPipe.Close()
	if err != nil {
		rPipe.CloseWithError(err)
		return err
	}
	wg.Wait()

	if !verified {
		return ErrIncompleteSnapshot
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

// RaftAutopilotState returns the state of the raft cluster as seen by autopilot.
func (c *Sys) RaftAutopilotState() (*AutopilotState, error) {
	r := c.c.NewRequest("GET", "/v1/sys/storage/raft/autopilot/state")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return nil, nil
		}
	}
	if err != nil {
		return nil, err
	}

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result AutopilotState
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

// RaftAutopilotConfiguration fetches the autopilot config.
func (c *Sys) RaftAutopilotConfiguration() (*AutopilotConfig, error) {
	r := c.c.NewRequest("GET", "/v1/sys/storage/raft/autopilot/configuration")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return nil, nil
		}
	}
	if err != nil {
		return nil, err
	}

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result AutopilotConfig
	if err = mapstructure.Decode(secret.Data, &result); err != nil {
		return nil, err
	}
	if result.LastContactThreshold, err = parseutil.ParseDurationSecond(secret.Data["last_contact_threshold"]); err != nil {
		return nil, err
	}
	if result.DeadServerLastContactThreshold, err = parseutil.ParseDurationSecond(secret.Data["dead_server_last_contact_threshold"]); err != nil {
		return nil, err
	}
	if result.ServerStabilizationTime, err = parseutil.ParseDurationSecond(secret.Data["server_stabilization_time"]); err != nil {
		return nil, err
	}

	return &result, err
}

// PutRaftAutopilotConfiguration allows modifying the raft autopilot configuration
func (c *Sys) PutRaftAutopilotConfiguration(opts *AutopilotConfig) error {
	r := c.c.NewRequest("POST", "/v1/sys/storage/raft/autopilot/configuration")

	if err := r.SetJSONBody(opts); err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
