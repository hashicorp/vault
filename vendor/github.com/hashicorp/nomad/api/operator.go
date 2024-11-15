// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Operator can be used to perform low-level operator tasks for Nomad.
type Operator struct {
	c *Client
}

// Operator returns a handle to the operator endpoints.
func (c *Client) Operator() *Operator {
	return &Operator{c}
}

// RaftServer has information about a server in the Raft configuration.
type RaftServer struct {
	// ID is the unique ID for the server. These are currently the same
	// as the address, but they will be changed to a real GUID in a future
	// release of Nomad.
	ID string

	// Node is the node name of the server, as known by Nomad, or this
	// will be set to "(unknown)" otherwise.
	Node string

	// Address is the IP:port of the server, used for Raft communications.
	Address string

	// Leader is true if this server is the current cluster leader.
	Leader bool

	// Voter is true if this server has a vote in the cluster. This might
	// be false if the server is staging and still coming online, or if
	// it's a non-voting server, which will be added in a future release of
	// Nomad.
	Voter bool

	// RaftProtocol is the version of the Raft protocol spoken by this server.
	RaftProtocol string
}

// RaftConfiguration is returned when querying for the current Raft configuration.
type RaftConfiguration struct {
	// Servers has the list of servers in the Raft configuration.
	Servers []*RaftServer

	// Index has the Raft index of this configuration.
	Index uint64
}

// RaftGetConfiguration is used to query the current Raft peer set.
func (op *Operator) RaftGetConfiguration(q *QueryOptions) (*RaftConfiguration, error) {
	r, err := op.c.newRequest("GET", "/v1/operator/raft/configuration")
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	_, resp, err := requireOK(op.c.doRequest(r)) //nolint:bodyclose
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out RaftConfiguration
	if err := decodeBody(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RaftRemovePeerByAddress is used to kick a stale peer (one that it in the Raft
// quorum but no longer known to Serf or the catalog) by address in the form of
// "IP:port".
func (op *Operator) RaftRemovePeerByAddress(address string, q *WriteOptions) error {
	r, err := op.c.newRequest("DELETE", "/v1/operator/raft/peer")
	if err != nil {
		return err
	}
	r.setWriteOptions(q)

	r.params.Set("address", address)

	_, resp, err := requireOK(op.c.doRequest(r)) //nolint:bodyclose
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

// RaftRemovePeerByID is used to kick a stale peer (one that is in the Raft
// quorum but no longer known to Serf or the catalog) by ID.
func (op *Operator) RaftRemovePeerByID(id string, q *WriteOptions) error {
	r, err := op.c.newRequest("DELETE", "/v1/operator/raft/peer")
	if err != nil {
		return err
	}
	r.setWriteOptions(q)

	r.params.Set("id", id)

	_, resp, err := requireOK(op.c.doRequest(r)) //nolint:bodyclose
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

// RaftTransferLeadershipByAddress is used to transfer leadership to a
// different peer using its address in the form of "IP:port".
func (op *Operator) RaftTransferLeadershipByAddress(address string, q *WriteOptions) error {
	r, err := op.c.newRequest("PUT", "/v1/operator/raft/transfer-leadership")
	if err != nil {
		return err
	}
	r.setWriteOptions(q)

	r.params.Set("address", address)

	_, resp, err := requireOK(op.c.doRequest(r)) //nolint:bodyclose
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

// RaftTransferLeadershipByID is used to transfer leadership to a
// different peer using its Raft ID.
func (op *Operator) RaftTransferLeadershipByID(id string, q *WriteOptions) error {
	r, err := op.c.newRequest("PUT", "/v1/operator/raft/transfer-leadership")
	if err != nil {
		return err
	}
	r.setWriteOptions(q)

	r.params.Set("id", id)

	_, resp, err := requireOK(op.c.doRequest(r)) //nolint:bodyclose
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

// SchedulerConfiguration is the config for controlling scheduler behavior
type SchedulerConfiguration struct {
	// SchedulerAlgorithm lets you select between available scheduling algorithms.
	SchedulerAlgorithm SchedulerAlgorithm

	// PreemptionConfig specifies whether to enable eviction of lower
	// priority jobs to place higher priority jobs.
	PreemptionConfig PreemptionConfig

	// MemoryOversubscriptionEnabled specifies whether memory oversubscription is enabled
	MemoryOversubscriptionEnabled bool

	// RejectJobRegistration disables new job registrations except with a
	// management ACL token
	RejectJobRegistration bool

	// PauseEvalBroker stops the leader evaluation broker process from running
	// until the configuration is updated and written to the Nomad servers.
	PauseEvalBroker bool

	// CreateIndex/ModifyIndex store the create/modify indexes of this configuration.
	CreateIndex uint64
	ModifyIndex uint64
}

// SchedulerConfigurationResponse is the response object that wraps SchedulerConfiguration
type SchedulerConfigurationResponse struct {
	// SchedulerConfig contains scheduler config options
	SchedulerConfig *SchedulerConfiguration

	QueryMeta
}

// SchedulerSetConfigurationResponse is the response object used
// when updating scheduler configuration
type SchedulerSetConfigurationResponse struct {
	// Updated returns whether the config was actually updated
	// Only set when the request uses CAS
	Updated bool

	WriteMeta
}

// SchedulerAlgorithm is an enum string that encapsulates the valid options for a
// SchedulerConfiguration block's SchedulerAlgorithm. These modes will allow the
// scheduler to be user-selectable.
type SchedulerAlgorithm string

const (
	SchedulerAlgorithmBinpack SchedulerAlgorithm = "binpack"
	SchedulerAlgorithmSpread  SchedulerAlgorithm = "spread"
)

// PreemptionConfig specifies whether preemption is enabled based on scheduler type
type PreemptionConfig struct {
	SystemSchedulerEnabled   bool
	SysBatchSchedulerEnabled bool
	BatchSchedulerEnabled    bool
	ServiceSchedulerEnabled  bool
}

// SchedulerGetConfiguration is used to query the current Scheduler configuration.
func (op *Operator) SchedulerGetConfiguration(q *QueryOptions) (*SchedulerConfigurationResponse, *QueryMeta, error) {
	var resp SchedulerConfigurationResponse
	qm, err := op.c.query("/v1/operator/scheduler/configuration", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// SchedulerSetConfiguration is used to set the current Scheduler configuration.
func (op *Operator) SchedulerSetConfiguration(conf *SchedulerConfiguration, q *WriteOptions) (*SchedulerSetConfigurationResponse, *WriteMeta, error) {
	var out SchedulerSetConfigurationResponse
	wm, err := op.c.put("/v1/operator/scheduler/configuration", conf, &out, q)
	if err != nil {
		return nil, nil, err
	}
	return &out, wm, nil
}

// SchedulerCASConfiguration is used to perform a Check-And-Set update on the
// Scheduler configuration. The ModifyIndex value will be respected. Returns
// true on success or false on failures.
func (op *Operator) SchedulerCASConfiguration(conf *SchedulerConfiguration, q *WriteOptions) (*SchedulerSetConfigurationResponse, *WriteMeta, error) {
	var out SchedulerSetConfigurationResponse
	wm, err := op.c.put("/v1/operator/scheduler/configuration?cas="+strconv.FormatUint(conf.ModifyIndex, 10), conf, &out, q)
	if err != nil {
		return nil, nil, err
	}

	return &out, wm, nil
}

// Snapshot is used to capture a snapshot state of a running cluster.
// The returned reader that must be consumed fully
func (op *Operator) Snapshot(q *QueryOptions) (io.ReadCloser, error) {
	r, err := op.c.newRequest("GET", "/v1/operator/snapshot")
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	_, resp, err := requireOK(op.c.doRequest(r)) //nolint:bodyclose
	if err != nil {
		return nil, err
	}

	digest := resp.Header.Get("Digest")

	cr, err := newChecksumValidatingReader(resp.Body, digest)
	if err != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		return nil, err
	}

	return cr, nil
}

// SnapshotRestore is used to restore a running nomad cluster to an original
// state.
func (op *Operator) SnapshotRestore(in io.Reader, q *WriteOptions) (*WriteMeta, error) {
	wm, err := op.c.put("/v1/operator/snapshot", in, nil, q)
	if err != nil {
		return nil, err
	}

	return wm, nil
}

type License struct {
	// The unique identifier of the license
	LicenseID string

	// The customer ID associated with the license
	CustomerID string

	// If set, an identifier that should be used to lock the license to a
	// particular site, cluster, etc.
	InstallationID string

	// The time at which the license was issued
	IssueTime time.Time

	// The time at which the license starts being valid
	StartTime time.Time

	// The time after which the license expires
	ExpirationTime time.Time

	// The time at which the license ceases to function and can
	// no longer be used in any capacity
	TerminationTime time.Time

	// The product the license is valid for
	Product string

	// License Specific Flags
	Flags map[string]interface{}

	// Modules is a list of the licensed enterprise modules
	Modules []string

	// List of features enabled by the license
	Features []string
}

type LicenseReply struct {
	License        *License
	ConfigOutdated bool
	QueryMeta
}

type ApplyLicenseOptions struct {
	Force bool
}

func (op *Operator) LicensePut(license string, q *WriteOptions) (*WriteMeta, error) {
	return op.ApplyLicense(license, nil, q)
}

func (op *Operator) ApplyLicense(license string, opts *ApplyLicenseOptions, q *WriteOptions) (*WriteMeta, error) {
	r, err := op.c.newRequest("PUT", "/v1/operator/license")
	if err != nil {
		return nil, err
	}

	if opts != nil && opts.Force {
		r.params.Add("force", "true")
	}

	r.setWriteOptions(q)
	r.body = strings.NewReader(license)

	rtt, resp, err := requireOK(op.c.doRequest(r)) //nolint:bodyclose
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	wm := &WriteMeta{RequestTime: rtt}
	parseWriteMeta(resp, wm)

	return wm, nil
}

func (op *Operator) LicenseGet(q *QueryOptions) (*LicenseReply, *QueryMeta, error) {
	req, err := op.c.newRequest("GET", "/v1/operator/license")
	if err != nil {
		return nil, nil, err
	}
	req.setQueryOptions(q)

	var reply LicenseReply
	rtt, resp, err := op.c.doRequest(req) //nolint:bodyclose
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil, errors.New("Nomad Enterprise only endpoint")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, newUnexpectedResponseError(
			fromHTTPResponse(resp),
			withExpectedStatuses([]int{http.StatusOK, http.StatusNoContent}),
		)
	}

	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return nil, nil, err
	}

	qm := &QueryMeta{}
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	return &reply, qm, nil
}

type LeadershipTransferResponse struct {
	From RaftServer
	To   RaftServer
	Noop bool
	Err  error

	WriteMeta
}

// VaultWorkloadIdentityUpgradeCheck is the result of verifying if the cluster
// is ready to switch to workload identities for Vault.
type VaultWorkloadIdentityUpgradeCheck struct {
	// JobsWithoutVaultIdentity is the list of jobs that have a `vault` block
	// but do not have an `identity` for Vault.
	JobsWithoutVaultIdentity []*JobListStub

	// OutdatedNodes is the list of nodes running a version of Nomad that does
	// not support workload identities for Vault.
	OutdatedNodes []*NodeListStub

	// VaultTokens is the list of Vault ACL token accessors that Nomad created
	// and will no longer manage after the cluster is migrated to workload
	// identities.
	VaultTokens []*VaultAccessor
}

// Ready returns true if the cluster is ready to migrate to workload identities
// with Vault.
func (v *VaultWorkloadIdentityUpgradeCheck) Ready() bool {
	return v != nil &&
		len(v.VaultTokens) == 0 &&
		len(v.OutdatedNodes) == 0 &&
		len(v.JobsWithoutVaultIdentity) == 0
}

// VaultAccessor is a Vault ACL token created by Nomad for a task to access
// Vault using the legacy authentication flow.
type VaultAccessor struct {
	// AllocID is the ID of the allocation that requested this token.
	AllocID string

	// Task is the name of the task that requested this token.
	Task string

	// NodeID is the ID of the node running the allocation that requested this
	// token.
	NodeID string

	// Accessor is the Vault ACL token accessor ID.
	Accessor string

	// CreationTTL is the TTL set when the token was created.
	CreationTTL int

	// CreateIndex is the Raft index when the token was created.
	CreateIndex uint64
}

// UpgradeCheckVaultWorkloadIdentity retrieves the cluster status for migrating
// to workload identities with Vault.
func (op *Operator) UpgradeCheckVaultWorkloadIdentity(q *QueryOptions) (*VaultWorkloadIdentityUpgradeCheck, *QueryMeta, error) {
	var resp VaultWorkloadIdentityUpgradeCheck
	qm, err := op.c.query("/v1/operator/upgrade-check/vault-workload-identity", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}
