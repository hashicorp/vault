// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

// Agent encapsulates an API client which talks to Nomad's
// agent endpoints for a specific node.
type Agent struct {
	client *Client

	// Cache static agent info
	nodeName   string
	datacenter string
	region     string
}

// KeyringResponse is a unified key response and can be used for install,
// remove, use, as well as listing key queries.
type KeyringResponse struct {
	Messages map[string]string
	Keys     map[string]int
	NumNodes int
}

// KeyringRequest is request objects for serf key operations.
type KeyringRequest struct {
	Key string
}

// ForceLeaveOpts are used to configure the ForceLeave method.
type ForceLeaveOpts struct {
	// Prune indicates whether to remove a node from the list of members
	Prune bool
}

// Agent returns a new agent which can be used to query
// the agent-specific endpoints.
func (c *Client) Agent() *Agent {
	return &Agent{client: c}
}

// Self is used to query the /v1/agent/self endpoint and
// returns information specific to the running agent.
func (a *Agent) Self() (*AgentSelf, error) {
	var out *AgentSelf

	// Query the self endpoint on the agent
	_, err := a.client.query("/v1/agent/self", &out, nil)
	if err != nil {
		return nil, fmt.Errorf("failed querying self endpoint: %s", err)
	}

	// Populate the cache for faster queries
	a.populateCache(out)

	return out, nil
}

// populateCache is used to insert various pieces of static
// data into the agent handle. This is used during subsequent
// lookups for the same data later on to save the round trip.
func (a *Agent) populateCache(self *AgentSelf) {
	if a.nodeName == "" {
		a.nodeName = self.Member.Name
	}
	if a.datacenter == "" {
		if val, ok := self.Config["Datacenter"]; ok {
			a.datacenter, _ = val.(string)
		}
	}
	if a.region == "" {
		if val, ok := self.Config["Region"]; ok {
			a.region, _ = val.(string)
		}
	}
}

// NodeName is used to query the Nomad agent for its node name.
func (a *Agent) NodeName() (string, error) {
	// Return from cache if we have it
	if a.nodeName != "" {
		return a.nodeName, nil
	}

	// Query the node name
	_, err := a.Self()
	return a.nodeName, err
}

// Datacenter is used to return the name of the datacenter which
// the agent is a member of.
func (a *Agent) Datacenter() (string, error) {
	// Return from cache if we have it
	if a.datacenter != "" {
		return a.datacenter, nil
	}

	// Query the agent for the DC
	_, err := a.Self()
	return a.datacenter, err
}

// Region is used to look up the region the agent is in.
func (a *Agent) Region() (string, error) {
	// Return from cache if we have it
	if a.region != "" {
		return a.region, nil
	}

	// Query the agent for the region
	_, err := a.Self()
	return a.region, err
}

// Join is used to instruct a server node to join another server
// via the gossip protocol. Multiple addresses may be specified.
// We attempt to join all the hosts in the list. Returns the
// number of nodes successfully joined and any error. If one or
// more nodes have a successful result, no error is returned.
func (a *Agent) Join(addrs ...string) (int, error) {
	// Accumulate the addresses
	v := url.Values{}
	for _, addr := range addrs {
		v.Add("address", addr)
	}

	// Send the join request
	var resp joinResponse
	_, err := a.client.put("/v1/agent/join?"+v.Encode(), nil, &resp, nil)
	if err != nil {
		return 0, fmt.Errorf("failed joining: %s", err)
	}
	if resp.Error != "" {
		return 0, fmt.Errorf("failed joining: %s", resp.Error)
	}
	return resp.NumJoined, nil
}

// Members is used to query all of the known server members
func (a *Agent) Members() (*ServerMembers, error) {
	var resp *ServerMembers

	// Query the known members
	_, err := a.client.query("/v1/agent/members", &resp, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Members is used to query all of the known server members
// with the ability to set QueryOptions
func (a *Agent) MembersOpts(opts *QueryOptions) (*ServerMembers, error) {
	var resp *ServerMembers
	_, err := a.client.query("/v1/agent/members", &resp, opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ForceLeave is used to eject an existing node from the cluster.
func (a *Agent) ForceLeave(node string) error {
	v := url.Values{}
	v.Add("node", node)
	_, err := a.client.put("/v1/agent/force-leave?"+v.Encode(), nil, nil, nil)
	return err
}

// ForceLeaveWithOptions is used to eject an existing node from the cluster
// with additional options such as prune.
func (a *Agent) ForceLeaveWithOptions(node string, opts ForceLeaveOpts) error {
	v := url.Values{}
	v.Add("node", node)
	if opts.Prune {
		v.Add("prune", "1")
	}
	_, err := a.client.put("/v1/agent/force-leave?"+v.Encode(), nil, nil, nil)
	return err
}

// Servers is used to query the list of servers on a client node.
func (a *Agent) Servers() ([]string, error) {
	var resp []string
	_, err := a.client.query("/v1/agent/servers", &resp, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetServers is used to update the list of servers on a client node.
func (a *Agent) SetServers(addrs []string) error {
	// Accumulate the addresses
	v := url.Values{}
	for _, addr := range addrs {
		v.Add("address", addr)
	}

	_, err := a.client.put("/v1/agent/servers?"+v.Encode(), nil, nil, nil)
	return err
}

// ListKeys returns the list of installed keys
func (a *Agent) ListKeys() (*KeyringResponse, error) {
	var resp KeyringResponse
	_, err := a.client.query("/v1/agent/keyring/list", &resp, nil)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// InstallKey installs a key in the keyrings of all the serf members
func (a *Agent) InstallKey(key string) (*KeyringResponse, error) {
	args := KeyringRequest{
		Key: key,
	}
	var resp KeyringResponse
	_, err := a.client.put("/v1/agent/keyring/install", &args, &resp, nil)
	return &resp, err
}

// UseKey uses a key from the keyring of serf members
func (a *Agent) UseKey(key string) (*KeyringResponse, error) {
	args := KeyringRequest{
		Key: key,
	}
	var resp KeyringResponse
	_, err := a.client.put("/v1/agent/keyring/use", &args, &resp, nil)
	return &resp, err
}

// RemoveKey removes a particular key from keyrings of serf members
func (a *Agent) RemoveKey(key string) (*KeyringResponse, error) {
	args := KeyringRequest{
		Key: key,
	}
	var resp KeyringResponse
	_, err := a.client.put("/v1/agent/keyring/remove", &args, &resp, nil)
	return &resp, err
}

// Health queries the agent's health
func (a *Agent) Health() (*AgentHealthResponse, error) {
	req, err := a.client.newRequest("GET", "/v1/agent/health")
	if err != nil {
		return nil, err
	}

	var health AgentHealthResponse
	_, resp, err := a.client.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Always try to decode the response as JSON
	err = json.NewDecoder(resp.Body).Decode(&health)
	if err == nil {
		return &health, nil
	}

	// Return custom error when response is not expected JSON format
	return nil, fmt.Errorf("unable to unmarshal response with status %d: %v", resp.StatusCode, err)
}

// Host returns debugging context about the agent's host operating system
func (a *Agent) Host(serverID, nodeID string, q *QueryOptions) (*HostDataResponse, error) {
	if q == nil {
		q = &QueryOptions{}
	}
	if q.Params == nil {
		q.Params = make(map[string]string)
	}

	if serverID != "" {
		q.Params["server_id"] = serverID
	}

	if nodeID != "" {
		q.Params["node_id"] = nodeID
	}

	var resp HostDataResponse
	_, err := a.client.query("/v1/agent/host", &resp, q)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Monitor returns a channel which will receive streaming logs from the agent
// Providing a non-nil stopCh can be used to close the connection and stop log streaming
func (a *Agent) Monitor(stopCh <-chan struct{}, q *QueryOptions) (<-chan *StreamFrame, <-chan error) {
	errCh := make(chan error, 1)
	r, err := a.client.newRequest("GET", "/v1/agent/monitor")
	if err != nil {
		errCh <- err
		return nil, errCh
	}

	r.setQueryOptions(q)
	_, resp, err := requireOK(a.client.doRequest(r)) //nolint:bodyclose
	if err != nil {
		errCh <- err
		return nil, errCh
	}

	frames := make(chan *StreamFrame, 10)
	go func() {
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)

		for {
			select {
			case <-stopCh:
				close(frames)
				return
			default:
			}

			// Decode the next frame
			var frame StreamFrame
			if err := dec.Decode(&frame); err != nil {
				close(frames)
				errCh <- err
				return
			}

			// Discard heartbeat frame
			if frame.IsHeartbeat() {
				continue
			}

			frames <- &frame
		}
	}()

	return frames, errCh
}

// PprofOptions contain a set of parameters for profiling a node or server.
type PprofOptions struct {
	// ServerID is the server ID, name, or special value "leader" to
	// specify the server that a given profile should be run on.
	ServerID string

	// NodeID is the node ID that a given profile should be run on.
	NodeID string

	// Seconds specifies the amount of time a profile should be run for.
	// Seconds only applies for certain runtime profiles like CPU and Trace.
	Seconds int

	// GC determines if a runtime.GC() should be called before a heap
	// profile.
	GC int

	// Debug specifies if the output of a lookup profile should be returned
	// in human readable format instead of binary.
	Debug int
}

// CPUProfile returns a runtime/pprof cpu profile for a given server or node.
// The profile will run for the amount of seconds passed in or default to 1.
// If no serverID or nodeID are provided the current Agents server will be
// used.
//
// The call blocks until the profile finishes, and returns the raw bytes of the
// profile.
func (a *Agent) CPUProfile(opts PprofOptions, q *QueryOptions) ([]byte, error) {
	return a.pprofRequest("profile", opts, q)
}

// Trace returns a runtime/pprof trace for a given server or node.
// The trace will run for the amount of seconds passed in or default to 1.
// If no serverID or nodeID are provided the current Agents server will be
// used.
//
// The call blocks until the profile finishes, and returns the raw bytes of the
// profile.
func (a *Agent) Trace(opts PprofOptions, q *QueryOptions) ([]byte, error) {
	return a.pprofRequest("trace", opts, q)
}

// Lookup returns a runtime/pprof profile using pprof.Lookup to determine
// which profile to run. Accepts a client or server ID but not both simultaneously.
//
// The call blocks until the profile finishes, and returns the raw bytes of the
// profile unless debug is set.
func (a *Agent) Lookup(profile string, opts PprofOptions, q *QueryOptions) ([]byte, error) {
	return a.pprofRequest(profile, opts, q)
}

func (a *Agent) pprofRequest(req string, opts PprofOptions, q *QueryOptions) ([]byte, error) {
	if q == nil {
		q = &QueryOptions{}
	}
	if q.Params == nil {
		q.Params = make(map[string]string)
	}

	q.Params["seconds"] = strconv.Itoa(opts.Seconds)
	q.Params["debug"] = strconv.Itoa(opts.Debug)
	q.Params["gc"] = strconv.Itoa(opts.GC)
	q.Params["node_id"] = opts.NodeID
	q.Params["server_id"] = opts.ServerID

	body, err := a.client.rawQuery(fmt.Sprintf("/v1/agent/pprof/%s", req), q)
	if err != nil {
		return nil, err
	}

	resp, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// joinResponse is used to decode the response we get while
// sending a member join request.
type joinResponse struct {
	NumJoined int    `json:"num_joined"`
	Error     string `json:"error"`
}

type ServerMembers struct {
	ServerName   string
	ServerRegion string
	ServerDC     string
	Members      []*AgentMember
}

type AgentSelf struct {
	Config map[string]interface{}       `json:"config"`
	Member AgentMember                  `json:"member"`
	Stats  map[string]map[string]string `json:"stats"`
}

// AgentMember represents a cluster member known to the agent
type AgentMember struct {
	Name        string
	Addr        string
	Port        uint16
	Tags        map[string]string
	Status      string
	ProtocolMin uint8
	ProtocolMax uint8
	ProtocolCur uint8
	DelegateMin uint8
	DelegateMax uint8
	DelegateCur uint8
}

// AgentMembersNameSort implements sort.Interface for []*AgentMembersNameSort
// based on the Name, DC and Region
type AgentMembersNameSort []*AgentMember

func (a AgentMembersNameSort) Len() int      { return len(a) }
func (a AgentMembersNameSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a AgentMembersNameSort) Less(i, j int) bool {
	if a[i].Tags["region"] != a[j].Tags["region"] {
		return a[i].Tags["region"] < a[j].Tags["region"]
	}

	if a[i].Tags["dc"] != a[j].Tags["dc"] {
		return a[i].Tags["dc"] < a[j].Tags["dc"]
	}

	return a[i].Name < a[j].Name

}

// AgentHealthResponse is the response from the Health endpoint describing an
// agent's health.
type AgentHealthResponse struct {
	Client *AgentHealth `json:"client,omitempty"`
	Server *AgentHealth `json:"server,omitempty"`
}

// AgentHealth describes the Client or Server's health in a Health request.
type AgentHealth struct {
	// Ok is false if the agent is unhealthy
	Ok bool `json:"ok"`

	// Message describes why the agent is unhealthy
	Message string `json:"message"`
}

type HostData struct {
	OS          string
	Network     []map[string]string
	ResolvConf  string
	Hosts       string
	Environment map[string]string
	Disk        map[string]DiskUsage
}

type DiskUsage struct {
	DiskMB int64
	UsedMB int64
}

type HostDataResponse struct {
	AgentID  string
	HostData *HostData `json:",omitempty"`
}

// GetSchedulerWorkerConfig returns the targeted agent's worker pool configuration
func (a *Agent) GetSchedulerWorkerConfig(q *QueryOptions) (*SchedulerWorkerPoolArgs, error) {
	var resp AgentSchedulerWorkerConfigResponse
	_, err := a.client.query("/v1/agent/schedulers/config", &resp, q)
	if err != nil {
		return nil, err
	}

	return &SchedulerWorkerPoolArgs{NumSchedulers: resp.NumSchedulers, EnabledSchedulers: resp.EnabledSchedulers}, nil
}

// SetSchedulerWorkerConfig attempts to update the targeted agent's worker pool configuration
func (a *Agent) SetSchedulerWorkerConfig(args SchedulerWorkerPoolArgs, q *WriteOptions) (*SchedulerWorkerPoolArgs, error) {
	req := AgentSchedulerWorkerConfigRequest(args)
	var resp AgentSchedulerWorkerConfigResponse

	_, err := a.client.put("/v1/agent/schedulers/config", &req, &resp, q)
	if err != nil {
		return nil, err
	}

	return &SchedulerWorkerPoolArgs{NumSchedulers: resp.NumSchedulers, EnabledSchedulers: resp.EnabledSchedulers}, nil
}

type SchedulerWorkerPoolArgs struct {
	NumSchedulers     int
	EnabledSchedulers []string
}

// AgentSchedulerWorkerConfigRequest is used to provide new scheduler worker configuration
// to a specific Nomad server. EnabledSchedulers must contain at least the `_core` scheduler
// to be valid.
type AgentSchedulerWorkerConfigRequest struct {
	NumSchedulers     int      `json:"num_schedulers"`
	EnabledSchedulers []string `json:"enabled_schedulers"`
}

// AgentSchedulerWorkerConfigResponse contains the Nomad server's current running configuration
// as well as the server's id as a convenience. This can be used to provide starting values for
// creating an AgentSchedulerWorkerConfigRequest to make changes to the running configuration.
type AgentSchedulerWorkerConfigResponse struct {
	ServerID          string   `json:"server_id"`
	NumSchedulers     int      `json:"num_schedulers"`
	EnabledSchedulers []string `json:"enabled_schedulers"`
}

// GetSchedulerWorkersInfo returns the current status of all of the scheduler workers on
// a Nomad server.
func (a *Agent) GetSchedulerWorkersInfo(q *QueryOptions) (*AgentSchedulerWorkersInfo, error) {
	var out *AgentSchedulerWorkersInfo

	_, err := a.client.query("/v1/agent/schedulers", &out, q)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// AgentSchedulerWorkersInfo is the response from the scheduler information endpoint containing
// a detailed status of each scheduler worker running on the server.
type AgentSchedulerWorkersInfo struct {
	ServerID   string                     `json:"server_id"`
	Schedulers []AgentSchedulerWorkerInfo `json:"schedulers"`
}

// AgentSchedulerWorkerInfo holds the detailed status information for a single scheduler worker.
type AgentSchedulerWorkerInfo struct {
	ID                string   `json:"id"`
	EnabledSchedulers []string `json:"enabled_schedulers"`
	Started           string   `json:"started"`
	Status            string   `json:"status"`
	WorkloadStatus    string   `json:"workload_status"`
}
