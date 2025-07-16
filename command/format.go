// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/ryanuber/columnize"
)

const (
	// hopeDelim is the delimiter to use when splitting columns. We call it a
	// hopeDelim because we hope that it's never contained in a secret.
	hopeDelim = "♨"
)

type FormatOptions struct {
	Format string
}

func OutputSecret(ui cli.Ui, secret *api.Secret) int {
	return outputWithFormat(ui, secret, secret)
}

func OutputList(ui cli.Ui, data interface{}) int {
	switch data := data.(type) {
	case *api.Secret:
		secret := data
		return outputWithFormat(ui, secret, secret.Data["keys"])
	default:
		return outputWithFormat(ui, nil, data)
	}
}

func OutputData(ui cli.Ui, data interface{}) int {
	return outputWithFormat(ui, nil, data)
}

func outputWithFormat(ui cli.Ui, secret *api.Secret, data interface{}) int {
	format := Format(ui)
	formatter, ok := Formatters[format]
	if !ok {
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}

	if err := formatter.Output(ui, secret, data); err != nil {
		ui.Error(fmt.Sprintf("Could not parse output: %s", err.Error()))
		return 1
	}
	return 0
}

type Formatter interface {
	Output(ui cli.Ui, secret *api.Secret, data interface{}) error
	Format(data interface{}) ([]byte, error)
}

var Formatters = map[string]Formatter{
	"json":   JsonFormatter{},
	"table":  TableFormatter{},
	"yaml":   YamlFormatter{},
	"yml":    YamlFormatter{},
	"pretty": PrettyFormatter{},
	"raw":    RawFormatter{},
}

func Format(ui cli.Ui) string {
	switch ui := ui.(type) {
	case *VaultUI:
		return ui.format
	}

	format := os.Getenv(EnvVaultFormat)
	if format == "" {
		format = "table"
	}

	return format
}

func Detailed(ui cli.Ui) bool {
	switch ui := ui.(type) {
	case *VaultUI:
		return ui.detailed
	}

	return false
}

// An output formatter for json output of an object
type JsonFormatter struct{}

func (j JsonFormatter) Format(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func (j JsonFormatter) Output(ui cli.Ui, secret *api.Secret, data interface{}) error {
	b, err := j.Format(data)
	if err != nil {
		return err
	}

	if secret != nil {
		shouldListWithInfo := Detailed(ui)

		// Show the raw JSON of the LIST call, rather than only the
		// list of keys.
		if shouldListWithInfo {
			b, err = j.Format(secret)
			if err != nil {
				return err
			}
		}
	}

	ui.Output(string(b))
	return nil
}

// An output formatter for raw output of the original request object
type RawFormatter struct{}

func (r RawFormatter) Format(data interface{}) ([]byte, error) {
	byte_data, ok := data.([]byte)
	if !ok {
		return nil, fmt.Errorf("This command does not support the -format=raw option; only `vault read` does.")
	}

	return byte_data, nil
}

func (r RawFormatter) Output(ui cli.Ui, secret *api.Secret, data interface{}) error {
	b, err := r.Format(data)
	if err != nil {
		return err
	}
	ui.Output(string(b))
	return nil
}

// An output formatter for yaml output format of an object
type YamlFormatter struct{}

func (y YamlFormatter) Format(data interface{}) ([]byte, error) {
	return yaml.Marshal(data)
}

func (y YamlFormatter) Output(ui cli.Ui, secret *api.Secret, data interface{}) error {
	b, err := y.Format(data)
	if err == nil {
		ui.Output(strings.TrimSpace(string(b)))
	}
	return err
}

type PrettyFormatter struct{}

func (p PrettyFormatter) Format(data interface{}) ([]byte, error) {
	return nil, nil
}

func (p PrettyFormatter) Output(ui cli.Ui, secret *api.Secret, data interface{}) error {
	switch data.(type) {
	case *api.AutopilotState:
		p.OutputAutopilotState(ui, data)
	default:
		return errors.New("cannot use the pretty formatter for this type")
	}
	return nil
}

func outputStringSlice(buffer *bytes.Buffer, indent string, values []string) {
	for _, val := range values {
		buffer.WriteString(fmt.Sprintf("%s%s\n", indent, val))
	}
}

type mapOutput struct {
	key   string
	value string
}

func formatServer(srv *api.AutopilotServer) string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("   %s\n", srv.ID))
	buffer.WriteString(fmt.Sprintf("      Name:              %s\n", srv.Name))
	buffer.WriteString(fmt.Sprintf("      Address:           %s\n", srv.Address))
	buffer.WriteString(fmt.Sprintf("      Status:            %s\n", srv.Status))
	buffer.WriteString(fmt.Sprintf("      Node Status:       %s\n", srv.NodeStatus))
	buffer.WriteString(fmt.Sprintf("      Healthy:           %t\n", srv.Healthy))
	buffer.WriteString(fmt.Sprintf("      Last Contact:      %s\n", srv.LastContact))
	buffer.WriteString(fmt.Sprintf("      Last Term:         %d\n", srv.LastTerm))
	buffer.WriteString(fmt.Sprintf("      Last Index:        %d\n", srv.LastIndex))
	buffer.WriteString(fmt.Sprintf("      Version:           %s\n", srv.Version))

	if srv.UpgradeVersion != "" {
		buffer.WriteString(fmt.Sprintf("      Upgrade Version:   %s\n", srv.UpgradeVersion))
	}
	if srv.RedundancyZone != "" {
		buffer.WriteString(fmt.Sprintf("      Redundancy Zone:   %s\n", srv.RedundancyZone))
	}
	if srv.NodeType != "" {
		buffer.WriteString(fmt.Sprintf("      Node Type:         %s\n", srv.NodeType))
	}

	return buffer.String()
}

func (p PrettyFormatter) OutputAutopilotState(ui cli.Ui, data interface{}) {
	state := data.(*api.AutopilotState)

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Healthy:                         %t\n", state.Healthy))
	buffer.WriteString(fmt.Sprintf("Failure Tolerance:               %d\n", state.FailureTolerance))
	buffer.WriteString(fmt.Sprintf("Leader:                          %s\n", state.Leader))
	buffer.WriteString("Voters:\n")
	outputStringSlice(&buffer, "   ", state.Voters)

	if len(state.NonVoters) > 0 {
		buffer.WriteString("Non Voters:\n")
		outputStringSlice(&buffer, "   ", state.NonVoters)
	}

	if state.OptimisticFailureTolerance > 0 {
		buffer.WriteString(fmt.Sprintf("Optimistic Failure Tolerance:    %d\n", state.OptimisticFailureTolerance))
	}

	// Servers
	buffer.WriteString("Servers:\n")
	var outputs []mapOutput
	for id, srv := range state.Servers {
		outputs = append(outputs, mapOutput{key: id, value: formatServer(srv)})
	}
	sort.Slice(outputs, func(i, j int) bool {
		return outputs[i].key < outputs[j].key
	})
	for _, output := range outputs {
		buffer.WriteString(output.value)
	}

	// Redundancy Zones
	if len(state.RedundancyZones) > 0 {
		buffer.WriteString("Redundancy Zones:\n")
		zoneList := make([]string, 0, len(state.RedundancyZones))
		for z := range state.RedundancyZones {
			zoneList = append(zoneList, z)
		}
		sort.Strings(zoneList)
		for _, zoneName := range zoneList {
			zone := state.RedundancyZones[zoneName]
			servers := zone.Servers
			voters := zone.Voters
			sort.Strings(servers)
			sort.Strings(voters)
			buffer.WriteString(fmt.Sprintf("   %s\n", zoneName))
			buffer.WriteString(fmt.Sprintf("      Servers: %s\n", strings.Join(servers, ", ")))
			buffer.WriteString(fmt.Sprintf("      Voters: %s\n", strings.Join(voters, ", ")))
			buffer.WriteString(fmt.Sprintf("      Failure Tolerance: %d\n", zone.FailureTolerance))
		}
	}

	// Upgrade Info
	if state.Upgrade != nil {
		buffer.WriteString("Upgrade Info:\n")
		buffer.WriteString(fmt.Sprintf("   Status: %s\n", state.Upgrade.Status))
		buffer.WriteString(fmt.Sprintf("   Target Version: %s\n", state.Upgrade.TargetVersion))
		buffer.WriteString(fmt.Sprintf("   Target Version Voters: %s\n", strings.Join(state.Upgrade.TargetVersionVoters, ", ")))
		buffer.WriteString(fmt.Sprintf("   Target Version Non-Voters: %s\n", strings.Join(state.Upgrade.TargetVersionNonVoters, ", ")))
		buffer.WriteString(fmt.Sprintf("   Other Version Voters: %s\n", strings.Join(state.Upgrade.OtherVersionVoters, ", ")))
		buffer.WriteString(fmt.Sprintf("   Other Version Non-Voters: %s\n", strings.Join(state.Upgrade.OtherVersionNonVoters, ", ")))

		if len(state.Upgrade.RedundancyZones) > 0 {
			buffer.WriteString("   Redundancy Zones:\n")
			for zoneName, zoneVersion := range state.Upgrade.RedundancyZones {
				buffer.WriteString(fmt.Sprintf("      %s\n", zoneName))
				buffer.WriteString(fmt.Sprintf("         Target Version Voters: %s\n", strings.Join(zoneVersion.TargetVersionVoters, ", ")))
				buffer.WriteString(fmt.Sprintf("         Target Version Non-Voters: %s\n", strings.Join(zoneVersion.TargetVersionNonVoters, ", ")))
				buffer.WriteString(fmt.Sprintf("         Other Version Voters: %s\n", strings.Join(zoneVersion.OtherVersionVoters, ", ")))
				buffer.WriteString(fmt.Sprintf("         Other Version Non-Voters: %s\n", strings.Join(zoneVersion.OtherVersionNonVoters, ", ")))
			}
		}
	}

	ui.Output(buffer.String())
}

// An output formatter for table output of an object
type TableFormatter struct{}

// We don't use this due to the TableFormatter introducing a bug when the -field flag is supplied:
// https://github.com/hashicorp/vault/commit/b24cf9a8af2190e96c614205b8cdf06d8c4b6718 .
func (t TableFormatter) Format(data interface{}) ([]byte, error) {
	return nil, nil
}

func (t TableFormatter) Output(ui cli.Ui, secret *api.Secret, data interface{}) error {
	switch data := data.(type) {
	case *api.Secret:
		return t.OutputSecret(ui, secret)
	case []interface{}:
		return t.OutputList(ui, secret, data)
	case []string:
		return t.OutputList(ui, nil, data)
	case map[string]interface{}:
		return t.OutputMap(ui, data)
	case SealStatusOutput:
		return t.OutputSealStatusStruct(ui, nil, data)
	default:
		return errors.New("cannot use the table formatter for this type")
	}
}

func (t TableFormatter) OutputSealStatusStruct(ui cli.Ui, secret *api.Secret, data interface{}) error {
	var status SealStatusOutput = data.(SealStatusOutput)
	var sealPrefix string

	out := []string{}
	out = append(out, "Key | Value")
	out = append(out, fmt.Sprintf("Seal Type | %s", status.Type))
	if status.RecoverySeal {
		sealPrefix = "Recovery "
		out = append(out, fmt.Sprintf("Recovery Seal Type | %s", status.RecoverySealType))
	}
	out = append(out, fmt.Sprintf("Initialized | %t", status.Initialized))
	out = append(out, fmt.Sprintf("Sealed | %t", status.Sealed))
	out = append(out, fmt.Sprintf("Total %sShares | %d", sealPrefix, status.N))
	out = append(out, fmt.Sprintf("Threshold | %d", status.T))

	if status.Sealed {
		out = append(out, fmt.Sprintf("Unseal Progress | %d/%d", status.Progress, status.T))
		out = append(out, fmt.Sprintf("Unseal Nonce | %s", status.Nonce))
	}

	if status.Migration {
		out = append(out, fmt.Sprintf("Seal Migration in Progress | %t", status.Migration))
	}

	out = append(out, fmt.Sprintf("Version | %s", status.Version))
	out = append(out, fmt.Sprintf("Build Date | %s", status.BuildDate))
	out = append(out, fmt.Sprintf("Storage Type | %s", status.StorageType))

	if status.ClusterName != "" && status.ClusterID != "" {
		out = append(out, fmt.Sprintf("Cluster Name | %s", status.ClusterName))
		out = append(out, fmt.Sprintf("Cluster ID | %s", status.ClusterID))
	}

	if status.RemovedFromCluster != nil {
		out = append(out, fmt.Sprintf("Removed From Cluster | %t", *status.RemovedFromCluster))
	}

	// Output if HCP link is configured
	if status.HCPLinkStatus != "" {
		out = append(out, fmt.Sprintf("HCP Link Status | %s", status.HCPLinkStatus))
		out = append(out, fmt.Sprintf("HCP Link Resource ID | %s", status.HCPLinkResourceID))
	}

	// Output if HA is enabled
	out = append(out, fmt.Sprintf("HA Enabled | %t", status.HAEnabled))

	if status.HAEnabled {
		mode := "sealed"
		if !status.Sealed {
			out = append(out, fmt.Sprintf("HA Cluster | %s", status.LeaderClusterAddress))
			mode = "standby"
			showLeaderAddr := false
			if status.IsSelf {
				mode = "active"
			} else {
				if status.LeaderAddress == "" {
					status.LeaderAddress = "<none>"
				}
				showLeaderAddr = true
			}
			out = append(out, fmt.Sprintf("HA Mode | %s", mode))

			if status.IsSelf && !status.ActiveTime.IsZero() {
				out = append(out, fmt.Sprintf("Active Since | %s", status.ActiveTime.Format(time.RFC3339Nano)))
			}
			// This is down here just to keep ordering consistent
			if showLeaderAddr {
				out = append(out, fmt.Sprintf("Active Node Address | %s", status.LeaderAddress))
			}

			if status.PerfStandby {
				out = append(out, fmt.Sprintf("Performance Standby Node | %t", status.PerfStandby))
				out = append(out, fmt.Sprintf("Performance Standby Last Remote WAL | %d", status.PerfStandbyLastRemoteWAL))
			}
		}
	}

	if status.RaftCommittedIndex > 0 {
		out = append(out, fmt.Sprintf("Raft Committed Index | %d", status.RaftCommittedIndex))
	}
	if status.RaftAppliedIndex > 0 {
		out = append(out, fmt.Sprintf("Raft Applied Index | %d", status.RaftAppliedIndex))
	}
	if status.LastWAL != 0 {
		out = append(out, fmt.Sprintf("Last WAL | %d", status.LastWAL))
	}
	if len(status.Warnings) > 0 {
		out = append(out, fmt.Sprintf("Warnings | %v", status.Warnings))
	}

	ui.Output(tableOutput(out, &columnize.Config{
		Delim: "|",
	}))
	return nil
}

func (t TableFormatter) OutputList(ui cli.Ui, secret *api.Secret, data interface{}) error {
	t.printWarnings(ui, secret)

	// Determine if we have additional information from a ListResponseWithInfo endpoint.
	var additionalInfo map[string]interface{}
	if secret != nil {
		shouldListWithInfo := Detailed(ui)
		if additional, ok := secret.Data["key_info"]; shouldListWithInfo && ok && len(additional.(map[string]interface{})) > 0 {
			additionalInfo = additional.(map[string]interface{})
		}
	}

	switch data := data.(type) {
	case []interface{}:
	case []string:
		ui.Output(tableOutput(data, nil))
		return nil
	default:
		return errors.New("error: table formatter cannot output list for this data type")
	}

	list := data.([]interface{})

	if len(list) > 0 {
		keys := make([]string, len(list))
		for i, v := range list {
			typed, ok := v.(string)
			if !ok {
				return fmt.Errorf("%v is not a string", v)
			}
			keys[i] = typed
		}
		sort.Strings(keys)

		// If we have a ListResponseWithInfo endpoint, we'll need to show
		// additional headers. To satisfy the table outputter, we'll need
		// to concat them with the deliminator.
		var headers []string
		header := "Keys"
		if len(additionalInfo) > 0 {
			seenHeaders := make(map[string]bool)
			for key, rawValues := range additionalInfo {
				// Most endpoints use the well-behaved ListResponseWithInfo.
				// However, some use a hand-rolled equivalent, where the
				// returned "keys" doesn't match the key of the "key_info"
				// member (namely, /sys/policies/egp). We seek to exclude
				// headers only visible from "non-visitable" key_info rows,
				// to make table output less confusing. These non-visitable
				// rows will still be visible in the JSON output.
				index := sort.SearchStrings(keys, key)
				if index < len(keys) && keys[index] != key {
					continue
				}

				values := rawValues.(map[string]interface{})
				for key := range values {
					seenHeaders[key] = true
				}
			}

			for key := range seenHeaders {
				headers = append(headers, key)
			}
			sort.Strings(headers)

			header = header + hopeDelim + strings.Join(headers, hopeDelim)
		}

		// Finally, if we have a ListResponseWithInfo, we'll need to update
		// the returned rows to not just have the keys (in the sorted order),
		// but also have the values for each header (in their sorted order).
		rows := keys
		if len(additionalInfo) > 0 && len(headers) > 0 {
			for index, row := range rows {
				formatted := []string{row}
				if rawValues, ok := additionalInfo[row]; ok {
					values := rawValues.(map[string]interface{})
					for _, header := range headers {
						if rawValue, ok := values[header]; ok {
							if looksLikeDuration(header) {
								rawValue = humanDurationInt(rawValue)
							}

							formatted = append(formatted, fmt.Sprintf("%v", rawValue))
						} else {
							// Show a default empty n/a when this field is
							// missing from the additional information.
							formatted = append(formatted, "n/a")
						}
					}
				}

				rows[index] = strings.Join(formatted, hopeDelim)
			}
		}

		// Prepend the header to the formatted rows.
		output := append([]string{header}, rows...)
		ui.Output(tableOutput(output, &columnize.Config{
			Delim: hopeDelim,
		}))
	}

	return nil
}

// printWarnings prints any warnings in the secret.
func (t TableFormatter) printWarnings(ui cli.Ui, secret *api.Secret) {
	if secret != nil && len(secret.Warnings) > 0 {
		ui.Warn("WARNING! The following warnings were returned from Vault:\n")
		for _, warning := range secret.Warnings {
			ui.Warn(wrapAtLengthWithPadding(fmt.Sprintf("* %s", warning), 2))
			ui.Warn("")
		}
	}
}

func (t TableFormatter) OutputSecret(ui cli.Ui, secret *api.Secret) error {
	if secret == nil {
		return nil
	}

	t.printWarnings(ui, secret)

	out := make([]string, 0, 8)
	if secret.LeaseDuration > 0 {
		if secret.LeaseID != "" {
			out = append(out, fmt.Sprintf("lease_id %s %s", hopeDelim, secret.LeaseID))
			out = append(out, fmt.Sprintf("lease_duration %s %v", hopeDelim, humanDurationInt(secret.LeaseDuration)))
			out = append(out, fmt.Sprintf("lease_renewable %s %t", hopeDelim, secret.Renewable))
		} else {
			// This is probably the generic secret backend which has leases, but we
			// print them as refresh_interval to reduce confusion.
			out = append(out, fmt.Sprintf("refresh_interval %s %v", hopeDelim, humanDurationInt(secret.LeaseDuration)))
		}
	}

	if secret.Auth != nil {
		if secret.Auth.MFARequirement != nil {
			out = append(out, fmt.Sprintf("mfa_request_id %s %s", hopeDelim, secret.Auth.MFARequirement.MFARequestID))

			for k, constraintSet := range secret.Auth.MFARequirement.MFAConstraints {
				for _, constraint := range constraintSet.Any {
					out = append(out, fmt.Sprintf("mfa_constraint_%s_%s_id %s %s", k, constraint.Type, hopeDelim, constraint.ID))
					out = append(out, fmt.Sprintf("mfa_constraint_%s_%s_uses_passcode %s %t", k, constraint.Type, hopeDelim, constraint.UsesPasscode))
					if constraint.Name != "" {
						out = append(out, fmt.Sprintf("mfa_constraint_%s_%s_name %s %s", k, constraint.Type, hopeDelim, constraint.Name))
					}
				}
			}
		} else { // Token information only makes sense if no further MFA requirement (i.e. if we actually have a token)
			out = append(out, fmt.Sprintf("token %s %s", hopeDelim, secret.Auth.ClientToken))
			out = append(out, fmt.Sprintf("token_accessor %s %s", hopeDelim, secret.Auth.Accessor))
			// If the lease duration is 0, it's likely a root token, so output the
			// duration as "infinity" to clear things up.
			if secret.Auth.LeaseDuration == 0 {
				out = append(out, fmt.Sprintf("token_duration %s %s", hopeDelim, "∞"))
			} else {
				out = append(out, fmt.Sprintf("token_duration %s %v", hopeDelim, humanDurationInt(secret.Auth.LeaseDuration)))
			}
			out = append(out, fmt.Sprintf("token_renewable %s %t", hopeDelim, secret.Auth.Renewable))
			out = append(out, fmt.Sprintf("token_policies %s %q", hopeDelim, secret.Auth.TokenPolicies))
			out = append(out, fmt.Sprintf("identity_policies %s %q", hopeDelim, secret.Auth.IdentityPolicies))
			out = append(out, fmt.Sprintf("policies %s %q", hopeDelim, secret.Auth.Policies))
			for k, v := range secret.Auth.Metadata {
				out = append(out, fmt.Sprintf("token_meta_%s %s %v", k, hopeDelim, v))
			}
		}
	}

	if secret.WrapInfo != nil {
		out = append(out, fmt.Sprintf("wrapping_token: %s %s", hopeDelim, secret.WrapInfo.Token))
		out = append(out, fmt.Sprintf("wrapping_accessor: %s %s", hopeDelim, secret.WrapInfo.Accessor))
		out = append(out, fmt.Sprintf("wrapping_token_ttl: %s %v", hopeDelim, humanDurationInt(secret.WrapInfo.TTL)))
		out = append(out, fmt.Sprintf("wrapping_token_creation_time: %s %s", hopeDelim, secret.WrapInfo.CreationTime.String()))
		out = append(out, fmt.Sprintf("wrapping_token_creation_path: %s %s", hopeDelim, secret.WrapInfo.CreationPath))
		if secret.WrapInfo.WrappedAccessor != "" {
			out = append(out, fmt.Sprintf("wrapped_accessor: %s %s", hopeDelim, secret.WrapInfo.WrappedAccessor))
		}
	}

	if len(secret.Data) > 0 {
		keys := make([]string, 0, len(secret.Data))
		for k := range secret.Data {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := secret.Data[k]

			// If the field "looks" like a TTL, print it as a time duration instead.
			if looksLikeDuration(k) {
				v = humanDurationInt(v)
			}

			out = append(out, fmt.Sprintf("%s %s %v", k, hopeDelim, v))
		}
	}

	// If we got this far and still don't have any data, there's nothing to print,
	// sorry.
	if len(out) == 0 {
		return nil
	}

	// Prepend the header
	out = append([]string{"Key" + hopeDelim + "Value"}, out...)

	ui.Output(tableOutput(out, &columnize.Config{
		Delim: hopeDelim,
	}))
	return nil
}

func (t TableFormatter) OutputMap(ui cli.Ui, data map[string]interface{}) error {
	out := make([]string, 0, len(data)+1)
	if len(data) > 0 {
		keys := make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := data[k]

			// If the field "looks" like a TTL, print it as a time duration instead.
			if looksLikeDuration(k) {
				v = humanDurationInt(v)
			}

			out = append(out, fmt.Sprintf("%s %s %v", k, hopeDelim, v))
		}
	}

	// If we got this far and still don't have any data, there's nothing to print,
	// sorry.
	if len(out) == 0 {
		return nil
	}

	// Prepend the header
	out = append([]string{"Key" + hopeDelim + "Value"}, out...)

	ui.Output(tableOutput(out, &columnize.Config{
		Delim: hopeDelim,
	}))
	return nil
}

// OutputSealStatus will print *api.SealStatusResponse in the CLI according to the format provided
func OutputSealStatus(ui cli.Ui, client *api.Client, status *api.SealStatusResponse) int {
	sealStatusOutput := SealStatusOutput{SealStatusResponse: *status}

	// Mask the 'Vault is sealed' error, since this means HA is enabled, but that
	// we cannot query for the leader since we are sealed.
	leaderStatus, err := client.Sys().Leader()
	if err != nil && strings.Contains(err.Error(), "Vault is sealed") {
		leaderStatus = &api.LeaderResponse{HAEnabled: true}
		err = nil
	}
	if err != nil {
		ui.Error(fmt.Sprintf("Error checking leader status: %s", err))
		return 1
	}

	// copy leaderStatus fields into sealStatusOutput for display later
	sealStatusOutput.HAEnabled = leaderStatus.HAEnabled
	sealStatusOutput.IsSelf = leaderStatus.IsSelf
	sealStatusOutput.ActiveTime = leaderStatus.ActiveTime
	sealStatusOutput.LeaderAddress = leaderStatus.LeaderAddress
	sealStatusOutput.LeaderClusterAddress = leaderStatus.LeaderClusterAddress
	sealStatusOutput.PerfStandby = leaderStatus.PerfStandby
	sealStatusOutput.PerfStandbyLastRemoteWAL = leaderStatus.PerfStandbyLastRemoteWAL
	sealStatusOutput.LastWAL = leaderStatus.LastWAL
	sealStatusOutput.RaftCommittedIndex = leaderStatus.RaftCommittedIndex
	sealStatusOutput.RaftAppliedIndex = leaderStatus.RaftAppliedIndex
	OutputData(ui, sealStatusOutput)
	return 0
}

// looksLikeDuration checks if the given key "k" looks like a duration value.
// This is used to pretty-format duration values in responses, especially from
// plugins.
func looksLikeDuration(k string) bool {
	return k == "period" || strings.HasSuffix(k, "_period") ||
		k == "ttl" || strings.HasSuffix(k, "_ttl") ||
		k == "duration" || strings.HasSuffix(k, "_duration") ||
		k == "lease_max" || k == "ttl_max"
}

// This struct is responsible for capturing all the fields to be output by a
// vault status command, including fields that do not come from the status API.
// Currently we are adding the fields from api.LeaderResponse
type SealStatusOutput struct {
	api.SealStatusResponse
	HAEnabled                bool      `json:"ha_enabled"`
	IsSelf                   bool      `json:"is_self,omitempty"`
	ActiveTime               time.Time `json:"active_time,omitempty"`
	LeaderAddress            string    `json:"leader_address,omitempty"`
	LeaderClusterAddress     string    `json:"leader_cluster_address,omitempty"`
	PerfStandby              bool      `json:"performance_standby,omitempty"`
	PerfStandbyLastRemoteWAL uint64    `json:"performance_standby_last_remote_wal,omitempty"`
	LastWAL                  uint64    `json:"last_wal,omitempty"`
	RaftCommittedIndex       uint64    `json:"raft_committed_index,omitempty"`
	RaftAppliedIndex         uint64    `json:"raft_applied_index,omitempty"`
}
