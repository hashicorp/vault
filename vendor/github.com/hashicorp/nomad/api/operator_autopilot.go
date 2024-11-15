// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"encoding/json"
	"strconv"
	"time"
)

// AutopilotConfiguration is used for querying/setting the Autopilot configuration.
// Autopilot helps manage operator tasks related to Nomad servers like removing
// failed servers from the Raft quorum.
type AutopilotConfiguration struct {
	// CleanupDeadServers controls whether to remove dead servers from the Raft
	// peer list when a new server joins
	CleanupDeadServers bool

	// LastContactThreshold is the limit on the amount of time a server can go
	// without leader contact before being considered unhealthy.
	LastContactThreshold time.Duration

	// MaxTrailingLogs is the amount of entries in the Raft Log that a server can
	// be behind before being considered unhealthy.
	MaxTrailingLogs uint64

	// MinQuorum sets the minimum number of servers allowed in a cluster before
	// autopilot can prune dead servers.
	MinQuorum uint

	// ServerStabilizationTime is the minimum amount of time a server must be
	// in a stable, healthy state before it can be added to the cluster. Only
	// applicable with Raft protocol version 3 or higher.
	ServerStabilizationTime time.Duration

	// (Enterprise-only) EnableRedundancyZones specifies whether to enable redundancy zones.
	EnableRedundancyZones bool

	// (Enterprise-only) DisableUpgradeMigration will disable Autopilot's upgrade migration
	// strategy of waiting until enough newer-versioned servers have been added to the
	// cluster before promoting them to voters.
	DisableUpgradeMigration bool

	// (Enterprise-only) EnableCustomUpgrades specifies whether to enable using custom
	// upgrade versions when performing migrations.
	EnableCustomUpgrades bool

	// CreateIndex holds the index corresponding the creation of this configuration.
	// This is a read-only field.
	CreateIndex uint64

	// ModifyIndex will be set to the index of the last update when retrieving the
	// Autopilot configuration. Resubmitting a configuration with
	// AutopilotCASConfiguration will perform a check-and-set operation which ensures
	// there hasn't been a subsequent update since the configuration was retrieved.
	ModifyIndex uint64
}

func (u *AutopilotConfiguration) MarshalJSON() ([]byte, error) {
	type Alias AutopilotConfiguration
	return json.Marshal(&struct {
		LastContactThreshold    string
		ServerStabilizationTime string
		*Alias
	}{
		LastContactThreshold:    u.LastContactThreshold.String(),
		ServerStabilizationTime: u.ServerStabilizationTime.String(),
		Alias:                   (*Alias)(u),
	})
}

func (u *AutopilotConfiguration) UnmarshalJSON(data []byte) error {
	type Alias AutopilotConfiguration
	aux := &struct {
		LastContactThreshold    string
		ServerStabilizationTime string
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	if aux.LastContactThreshold != "" {
		if u.LastContactThreshold, err = time.ParseDuration(aux.LastContactThreshold); err != nil {
			return err
		}
	}
	if aux.ServerStabilizationTime != "" {
		if u.ServerStabilizationTime, err = time.ParseDuration(aux.ServerStabilizationTime); err != nil {
			return err
		}
	}
	return nil
}

// ServerHealth is the health (from the leader's point of view) of a server.
type ServerHealth struct {
	// ID is the raft ID of the server.
	ID string

	// Name is the node name of the server.
	Name string

	// Address is the address of the server.
	Address string

	// The status of the SerfHealth check for the server.
	SerfStatus string

	// Version is the Nomad version of the server.
	Version string

	// Leader is whether this server is currently the leader.
	Leader bool

	// LastContact is the time since this node's last contact with the leader.
	LastContact time.Duration

	// LastTerm is the highest leader term this server has a record of in its Raft log.
	LastTerm uint64

	// LastIndex is the last log index this server has a record of in its Raft log.
	LastIndex uint64

	// Healthy is whether or not the server is healthy according to the current
	// Autopilot config.
	Healthy bool

	// Voter is whether this is a voting server.
	Voter bool

	// StableSince is the last time this server's Healthy value changed.
	StableSince time.Time
}

func (u *ServerHealth) MarshalJSON() ([]byte, error) {
	type Alias ServerHealth
	return json.Marshal(&struct {
		LastContact string
		*Alias
	}{
		LastContact: u.LastContact.String(),
		Alias:       (*Alias)(u),
	})
}

func (u *ServerHealth) UnmarshalJSON(data []byte) error {
	type Alias ServerHealth
	aux := &struct {
		LastContact string
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	if aux.LastContact != "" {
		if u.LastContact, err = time.ParseDuration(aux.LastContact); err != nil {
			return err
		}
	}
	return nil
}

// OperatorHealthReply is a representation of the overall health of the cluster
type OperatorHealthReply struct {
	// Healthy is true if all the servers in the cluster are healthy.
	Healthy bool

	// FailureTolerance is the number of healthy servers that could be lost without
	// an outage occurring.
	FailureTolerance int

	// Servers holds the health of each server.
	Servers []ServerHealth

	// The ID of the current leader.
	Leader string

	// List of servers that are voters in the Raft configuration.
	Voters []string

	// ReadReplicas holds the list of servers that are
	// read replicas in the Raft configuration. (Enterprise only)
	ReadReplicas []string `json:",omitempty"`

	// RedundancyZones holds the list of servers in each redundancy zone.
	// (Enterprise only)
	RedundancyZones map[string]AutopilotZone `json:",omitempty"`

	// Upgrade holds the current upgrade status.
	Upgrade *AutopilotUpgrade `json:",omitempty"`

	// The number of servers that could be lost without an outage
	// occurring if all the voters don't fail at once.  (Enterprise only)
	OptimisticFailureTolerance int `json:",omitempty"`
}

// AutopilotZone holds the list of servers in a redundancy zone.  (Enterprise only)
type AutopilotZone struct {
	// Servers holds the list of servers in the redundancy zone.
	Servers []string

	// Voters holds the list of servers that are voters in the redundancy zone.
	Voters []string

	// FailureTolerance is the number of servers that could be lost without an
	// outage occurring.
	FailureTolerance int
}

// AutopilotUpgrade holds the current upgrade status.  (Enterprise only)
type AutopilotUpgrade struct {
	// Status of the upgrade.
	Status string

	// TargetVersion is the version that the cluster is upgrading to.
	TargetVersion string

	// TargetVersionVoters holds the list of servers that are voters in the Raft
	// configuration of the TargetVersion.
	TargetVersionVoters []string

	// TargetVersionNonVoters holds the list of servers that are non-voters in
	// the Raft configuration of the TargetVersion.
	TargetVersionNonVoters []string

	// TargetVersionReadReplicas holds the list of servers that are read
	// replicas in the Raft configuration of the TargetVersion.
	TargetVersionReadReplicas []string

	// OtherVersionVoters holds the list of servers that are voters in the Raft
	// configuration of a version other than the TargetVersion.
	OtherVersionVoters []string

	// OtherVersionNonVoters holds the list of servers that are non-voters in
	// the Raft configuration of a version other than the TargetVersion.
	OtherVersionNonVoters []string

	// OtherVersionReadReplicas holds the list of servers that are read replicas
	// in the Raft configuration of a version other than the TargetVersion.
	OtherVersionReadReplicas []string

	// RedundancyZones holds the list of servers in each redundancy zone for the
	// TargetVersion.
	RedundancyZones map[string]AutopilotZoneUpgradeVersions
}

// AutopilotZoneUpgradeVersions holds the list of servers
// in a redundancy zone for a specific version.  (Enterprise only)
type AutopilotZoneUpgradeVersions struct {
	TargetVersionVoters    []string
	TargetVersionNonVoters []string
	OtherVersionVoters     []string
	OtherVersionNonVoters  []string
}

// AutopilotGetConfiguration is used to query the current Autopilot configuration.
func (op *Operator) AutopilotGetConfiguration(q *QueryOptions) (*AutopilotConfiguration, *QueryMeta, error) {
	var resp AutopilotConfiguration
	qm, err := op.c.query("/v1/operator/autopilot/configuration", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// AutopilotSetConfiguration is used to set the current Autopilot configuration.
func (op *Operator) AutopilotSetConfiguration(conf *AutopilotConfiguration, q *WriteOptions) (*WriteMeta, error) {
	var out bool
	wm, err := op.c.put("/v1/operator/autopilot/configuration", conf, &out, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// AutopilotCASConfiguration is used to perform a Check-And-Set update on the
// Autopilot configuration. The ModifyIndex value will be respected. Returns
// true on success or false on failures.
func (op *Operator) AutopilotCASConfiguration(conf *AutopilotConfiguration, q *WriteOptions) (bool, *WriteMeta, error) {
	var out bool
	wm, err := op.c.put("/v1/operator/autopilot/configuration?cas="+strconv.FormatUint(conf.ModifyIndex, 10), conf, &out, q)
	if err != nil {
		return false, nil, err
	}

	return out, wm, nil
}

// AutopilotServerHealth is used to query Autopilot's top-level view of the health
// of each Nomad server.
func (op *Operator) AutopilotServerHealth(q *QueryOptions) (*OperatorHealthReply, *QueryMeta, error) {
	var out OperatorHealthReply
	qm, err := op.c.query("/v1/operator/autopilot/health", &out, q)
	if err != nil {
		return nil, nil, err
	}
	return &out, qm, nil
}
