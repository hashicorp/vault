// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	"github.com/posener/complete"
	"github.com/ryanuber/columnize"
)

const (
	healthWALLagWarning      = 100
	healthWALLagCritical     = 1000
	healthHeartbeatWarning   = 30 * time.Second
	healthHeartbeatCritical  = 2 * time.Minute
	healthTLSExpiryWarning   = 30 * 24 * time.Hour
	healthTLSExpiryCritical  = 7 * 24 * time.Hour
	healthMaximumDiscovered  = 64
	healthLevelHealthy       = "healthy"
	healthLevelWarning       = "warning"
	healthLevelCritical      = "critical"
	healthLevelNotApplicable = "not_applicable"
)

var (
	_ cli.Command             = (*HealthCommand)(nil)
	_ cli.CommandAutocomplete = (*HealthCommand)(nil)
)

type HealthCommand struct {
	*BaseCommand

	now func() time.Time
}

type healthReport struct {
	Address string        `json:"address"`
	Overall string        `json:"overall"`
	Nodes   []healthNode  `json:"nodes"`
	Errors  []healthError `json:"errors,omitempty"`
}

type healthNode struct {
	Name        string                  `json:"name"`
	Address     string                  `json:"address"`
	Version     string                  `json:"version,omitempty"`
	Initialized bool                    `json:"initialized"`
	Sealed      bool                    `json:"sealed"`
	Standby     bool                    `json:"standby"`
	TLS         healthTLS               `json:"tls"`
	Replication []healthReplicationMode `json:"replication,omitempty"`
}

type healthTLS struct {
	Enabled       bool       `json:"enabled"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	ExpiresInDays *int       `json:"expires_in_days,omitempty"`
}

type healthReplicationMode struct {
	Mode        string   `json:"mode"`
	Role        string   `json:"role"`
	State       string   `json:"state,omitempty"`
	ClusterID   string   `json:"cluster_id,omitempty"`
	WAL         *uint64  `json:"wal,omitempty"`
	RemoteWAL   *uint64  `json:"remote_wal,omitempty"`
	WALLag      *int64   `json:"wal_lag,omitempty"`
	Health      string   `json:"health"`
	Diagnostics []string `json:"diagnostics,omitempty"`
}

type healthError struct {
	Address string `json:"address"`
	Check   string `json:"check"`
	Error   string `json:"error"`
}

type healthProbe struct {
	address     string
	health      *api.HealthResponse
	tls         healthTLS
	replication *healthReplicationStatus
	healthErr   error
	replErr     error
}

type healthReplicationStatus struct {
	DR          healthReplicationStatusMode `mapstructure:"dr"`
	Performance healthReplicationStatusMode `mapstructure:"performance"`
}

type healthReplicationStatusMode struct {
	LastDRWAL           uint64                  `mapstructure:"last_dr_wal"`
	ClusterID           string                  `mapstructure:"cluster_id"`
	LastWAL             uint64                  `mapstructure:"last_wal"`
	CorruptedMerkleTree bool                    `mapstructure:"corrupted_merkle_tree"`
	Mode                string                  `mapstructure:"mode"`
	LastPerformanceWAL  uint64                  `mapstructure:"last_performance_wal"`
	State               string                  `mapstructure:"state"`
	ConnectionState     string                  `mapstructure:"connection_state"`
	LastRemoteWAL       uint64                  `mapstructure:"last_remote_wal"`
	Primaries           []healthReplicationPeer `mapstructure:"primaries"`
	Secondaries         []healthReplicationPeer `mapstructure:"secondaries"`
}

type healthReplicationPeer struct {
	APIAddress       string `mapstructure:"api_address"`
	ConnectionStatus string `mapstructure:"connection_status"`
	LastHeartbeat    string `mapstructure:"last_heartbeat"`
}

func (c *HealthCommand) Synopsis() string {
	return "Print health and replication status"
}

func (c *HealthCommand) Help() string {
	helpText := `
Usage: vault health [options]

  Prints health and replication status for the configured Vault address. The
  command discovers other replication participants from the status returned by
  Vault and checks every API address it can reach. No topology configuration is
  required.

  The exit code reflects the overall health:

      - 0 - healthy, or replication is not configured
      - 1 - warning or command error
      - 2 - critical

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *HealthCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *HealthCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *HealthCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *HealthCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	report := c.buildHealthReport(client)
	if c.flagFormat == "table" {
		c.outputHealthReport(report)
	} else if code := OutputData(c.UI, report); code != 0 {
		return code
	}

	switch report.Overall {
	case healthLevelCritical:
		return 2
	case healthLevelWarning:
		return 1
	default:
		return 0
	}
}

func (c *HealthCommand) buildHealthReport(client *api.Client) healthReport {
	startAddress := normalizeHealthAddress(client.Address())
	probes, discoveryErr := c.discoverHealthProbes(client, startAddress)
	report := healthReport{
		Address: startAddress,
		Overall: healthLevelNotApplicable,
		Nodes:   make([]healthNode, 0, len(probes)),
	}
	if discoveryErr != nil {
		report.Errors = append(report.Errors, healthError{
			Address: startAddress,
			Check:   "discovery",
			Error:   discoveryErr.Error(),
		})
		report.Overall = healthLevelWarning
	}

	for _, probe := range probes {
		node := healthNode{
			Name:    healthNodeName(probe.address),
			Address: probe.address,
			TLS:     probe.tls,
		}
		if probe.healthErr != nil {
			report.Errors = append(report.Errors, healthError{
				Address: probe.address,
				Check:   "health",
				Error:   probe.healthErr.Error(),
			})
			report.Overall = worseHealthLevel(report.Overall, healthLevelCritical)
		} else {
			node.Version = probe.health.Version
			node.Initialized = probe.health.Initialized
			node.Sealed = probe.health.Sealed
			node.Standby = probe.health.Standby
			if !node.Initialized || node.Sealed {
				report.Overall = worseHealthLevel(report.Overall, healthLevelCritical)
			}
		}

		if probe.replErr != nil {
			report.Errors = append(report.Errors, healthError{
				Address: probe.address,
				Check:   "replication",
				Error:   probe.replErr.Error(),
			})
			report.Overall = worseHealthLevel(report.Overall, healthLevelWarning)
		} else {
			node.Replication = appendReplicationModes(node.Replication, probe.replication)
		}
		report.Nodes = append(report.Nodes, node)
	}

	sort.Slice(report.Nodes, func(i, j int) bool {
		return report.Nodes[i].Address < report.Nodes[j].Address
	})
	sort.Slice(report.Errors, func(i, j int) bool {
		if report.Errors[i].Address == report.Errors[j].Address {
			return report.Errors[i].Check < report.Errors[j].Check
		}
		return report.Errors[i].Address < report.Errors[j].Address
	})

	c.evaluateReplicationHealth(&report, probes)
	return report
}

func (c *HealthCommand) discoverHealthProbes(client *api.Client, startAddress string) ([]healthProbe, error) {
	queue := []string{startAddress}
	seen := make(map[string]struct{})
	probes := make([]healthProbe, 0)

	for len(queue) > 0 {
		address := queue[0]
		queue = queue[1:]
		if _, ok := seen[address]; ok {
			continue
		}
		if len(seen) == healthMaximumDiscovered {
			return probes, fmt.Errorf("stopped after discovering %d Vault addresses", healthMaximumDiscovered)
		}
		seen[address] = struct{}{}

		probeClient, err := cloneHealthClient(client, address)
		if err != nil {
			probes = append(probes, healthProbe{
				address:   address,
				healthErr: err,
				replErr:   err,
			})
			continue
		}

		probe := healthProbe{address: address}
		probe.health, probe.tls, probe.healthErr = readHealthAndTLS(probeClient)
		probe.replication, probe.replErr = readHealthReplicationStatus(probeClient)
		probes = append(probes, probe)

		if probe.replErr == nil {
			for _, peerAddress := range replicationPeerAddresses(probe.replication) {
				if _, ok := seen[peerAddress]; !ok {
					queue = append(queue, peerAddress)
				}
			}
		}
	}

	return probes, nil
}

func readHealthAndTLS(client *api.Client) (*api.HealthResponse, healthTLS, error) {
	address, err := url.Parse(client.Address())
	if err != nil {
		return nil, healthTLS{}, fmt.Errorf("failed to parse Vault address: %w", err)
	}
	tlsInfo := healthTLS{Enabled: address.Scheme == "https"}

	request := client.NewRequest(http.MethodGet, "/v1/sys/health")
	for name, value := range map[string]string{
		"uninitcode":             "299",
		"sealedcode":             "299",
		"standbycode":            "299",
		"drsecondarycode":        "299",
		"performancestandbycode": "299",
		"removedcode":            "299",
		"haunhealthycode":        "299",
	} {
		request.Params.Set(name, value)
	}

	response, err := client.RawRequest(request)
	if err != nil {
		return nil, tlsInfo, err
	}
	defer response.Body.Close()

	if tlsInfo.Enabled {
		if response.TLS == nil || len(response.TLS.PeerCertificates) == 0 {
			return nil, tlsInfo, fmt.Errorf("HTTPS response did not include a peer certificate")
		}
		expiresAt := response.TLS.PeerCertificates[0].NotAfter
		tlsInfo.ExpiresAt = &expiresAt
	}

	var health api.HealthResponse
	if err := response.DecodeJSON(&health); err != nil {
		return nil, tlsInfo, fmt.Errorf("failed to decode health response: %w", err)
	}
	return &health, tlsInfo, nil
}

func readHealthReplicationStatus(client *api.Client) (*healthReplicationStatus, error) {
	secret, err := client.Logical().Read("sys/replication/status")
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("empty replication status response")
	}

	var status healthReplicationStatus
	if err := mapstructure.WeakDecode(secret.Data, &status); err != nil {
		return nil, fmt.Errorf("failed to decode replication status: %w", err)
	}
	return &status, nil
}

func cloneHealthClient(client *api.Client, address string) (*api.Client, error) {
	cloned, err := client.CloneWithHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed to clone Vault client: %w", err)
	}
	cloned.SetToken(client.Token())
	if err := cloned.SetAddress(address); err != nil {
		return nil, fmt.Errorf("failed to use discovered address %q: %w", address, err)
	}
	return cloned, nil
}

func replicationPeerAddresses(status *healthReplicationStatus) []string {
	if status == nil {
		return nil
	}

	addresses := make(map[string]struct{})
	for _, replication := range []*healthReplicationStatusMode{
		&status.DR,
		&status.Performance,
	} {
		for _, peer := range append(replication.Primaries, replication.Secondaries...) {
			if address := normalizeHealthAddress(peer.APIAddress); address != "" {
				addresses[address] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(addresses))
	for address := range addresses {
		result = append(result, address)
	}
	sort.Strings(result)
	return result
}

func appendReplicationModes(modes []healthReplicationMode, status *healthReplicationStatus) []healthReplicationMode {
	if status == nil {
		return modes
	}

	for _, item := range []struct {
		name   string
		status healthReplicationStatusMode
	}{
		{name: "dr", status: status.DR},
		{name: "performance", status: status.Performance},
	} {
		if item.status.Mode == "" || item.status.Mode == "disabled" {
			continue
		}

		mode := healthReplicationMode{
			Mode:      item.name,
			Role:      item.status.Mode,
			State:     item.status.State,
			ClusterID: item.status.ClusterID,
			Health:    healthLevelHealthy,
		}
		if item.status.Mode == "primary" {
			wal := item.status.LastWAL
			if item.name == "dr" && item.status.LastDRWAL != 0 {
				wal = item.status.LastDRWAL
			}
			if item.name == "performance" && item.status.LastPerformanceWAL != 0 {
				wal = item.status.LastPerformanceWAL
			}
			mode.WAL = &wal
		} else {
			remoteWAL := item.status.LastRemoteWAL
			mode.RemoteWAL = &remoteWAL
		}
		modes = append(modes, mode)
	}

	return modes
}

func (c *HealthCommand) evaluateReplicationHealth(report *healthReport, probes []healthProbe) {
	now := time.Now()
	if c.now != nil {
		now = c.now()
	}

	for nodeIndex := range report.Nodes {
		node := &report.Nodes[nodeIndex]
		if node.TLS.ExpiresAt != nil {
			remaining := node.TLS.ExpiresAt.Sub(now)
			days := int(math.Floor(remaining.Hours() / 24))
			node.TLS.ExpiresInDays = &days
			switch {
			case remaining <= healthTLSExpiryCritical:
				report.Overall = worseHealthLevel(report.Overall, healthLevelCritical)
			case remaining <= healthTLSExpiryWarning:
				report.Overall = worseHealthLevel(report.Overall, healthLevelWarning)
			}
		}

		for modeIndex := range node.Replication {
			mode := &node.Replication[modeIndex]

			switch {
			case mode.Role == "primary" && mode.State != "running":
				mode.Health = healthLevelWarning
				mode.Diagnostics = append(mode.Diagnostics, "state="+mode.State)
			case mode.Role == "secondary" && mode.State != "stream-wals":
				mode.Health = healthLevelCritical
				mode.Diagnostics = append(mode.Diagnostics, "state="+mode.State)
			}

			primary := findPrimaryMode(report.Nodes, mode.Mode, mode.ClusterID)
			if mode.Role == "secondary" && mode.RemoteWAL != nil {
				if *mode.RemoteWAL == 0 {
					mode.Health = healthLevelCritical
					mode.Diagnostics = append(mode.Diagnostics, "never streamed (remote_wal=0)")
				} else if primary != nil && primary.WAL != nil {
					lag := int64(*primary.WAL) - int64(*mode.RemoteWAL)
					mode.WALLag = &lag
					absoluteLag := math.Abs(float64(lag))
					switch {
					case absoluteLag >= healthWALLagCritical:
						mode.Health = healthLevelCritical
						mode.Diagnostics = append(mode.Diagnostics, fmt.Sprintf("WAL lag %d", lag))
					case absoluteLag >= healthWALLagWarning:
						mode.Health = worseHealthLevel(mode.Health, healthLevelWarning)
						mode.Diagnostics = append(mode.Diagnostics, fmt.Sprintf("WAL lag %d", lag))
					}
				}
			}

			c.evaluatePeerHealth(mode, node.Address, probes, now)
			report.Overall = worseHealthLevel(report.Overall, mode.Health)
		}
	}
}

func (c *HealthCommand) evaluatePeerHealth(mode *healthReplicationMode, address string, probes []healthProbe, now time.Time) {
	replication := findReplicationStatus(address, mode.Mode, probes)
	if replication == nil {
		return
	}

	if replication.CorruptedMerkleTree {
		mode.Health = healthLevelCritical
		mode.Diagnostics = append(mode.Diagnostics, "merkle tree corrupted")
	}
	if mode.Role == "secondary" && replication.ConnectionState != "" && replication.ConnectionState != "ready" {
		mode.Health = healthLevelCritical
		mode.Diagnostics = append(mode.Diagnostics, "connection="+replication.ConnectionState)
	}

	for _, peer := range append(replication.Primaries, replication.Secondaries...) {
		name := healthNodeName(peer.APIAddress)
		if peer.ConnectionStatus != "" && peer.ConnectionStatus != "connected" && peer.ConnectionStatus != "ready" {
			mode.Health = healthLevelCritical
			mode.Diagnostics = append(mode.Diagnostics, fmt.Sprintf("%s %s", name, peer.ConnectionStatus))
		}
		if heartbeat, err := time.Parse(time.RFC3339Nano, peer.LastHeartbeat); err == nil {
			age := now.Sub(heartbeat)
			switch {
			case age >= healthHeartbeatCritical:
				mode.Health = healthLevelCritical
				mode.Diagnostics = append(mode.Diagnostics, fmt.Sprintf("%s heartbeat %s ago", name, age.Round(time.Second)))
			case age >= healthHeartbeatWarning:
				mode.Health = worseHealthLevel(mode.Health, healthLevelWarning)
				mode.Diagnostics = append(mode.Diagnostics, fmt.Sprintf("%s heartbeat %s ago", name, age.Round(time.Second)))
			}
		}
	}
}

func findHealthProbe(address string, probes []healthProbe) *healthProbe {
	for i := range probes {
		if probes[i].address == address {
			return &probes[i]
		}
	}
	return nil
}

func findReplicationStatus(address, mode string, probes []healthProbe) *healthReplicationStatusMode {
	probe := findHealthProbe(address, probes)
	if probe == nil || probe.replication == nil {
		return nil
	}
	if mode == "dr" {
		return &probe.replication.DR
	}
	return &probe.replication.Performance
}

func findPrimaryMode(nodes []healthNode, modeName, clusterID string) *healthReplicationMode {
	for nodeIndex := range nodes {
		for modeIndex := range nodes[nodeIndex].Replication {
			mode := &nodes[nodeIndex].Replication[modeIndex]
			if mode.Mode == modeName && mode.Role == "primary" && mode.ClusterID == clusterID {
				return mode
			}
		}
	}
	return nil
}

func worseHealthLevel(current, candidate string) string {
	rank := map[string]int{
		healthLevelNotApplicable: 0,
		healthLevelHealthy:       1,
		healthLevelWarning:       2,
		healthLevelCritical:      3,
	}
	if rank[candidate] > rank[current] {
		return candidate
	}
	return current
}

func normalizeHealthAddress(address string) string {
	return strings.TrimRight(strings.TrimSpace(address), "/")
}

func healthNodeName(address string) string {
	parsed, err := url.Parse(address)
	if err != nil || parsed.Hostname() == "" {
		return address
	}
	host := parsed.Hostname()
	if dot := strings.IndexByte(host, '.'); dot >= 0 {
		return host[:dot]
	}
	return host
}

func (c *HealthCommand) outputHealthReport(report healthReport) {
	nodes := []string{"Node | Address | Version | Initialized | Sealed | HA State | HTTPS | Cert Expires (days)"}
	for _, node := range report.Nodes {
		haState := "active"
		if node.Standby {
			haState = "standby"
		}
		nodes = append(nodes, fmt.Sprintf(
			"%s | %s | %s | %t | %t | %s | %t | %s",
			node.Name,
			node.Address,
			valueOrNA(node.Version),
			node.Initialized,
			node.Sealed,
			haState,
			node.TLS.Enabled,
			formatHealthIntPointer(node.TLS.ExpiresInDays),
		))
	}
	c.UI.Output(tableOutput(nodes, columnize.DefaultConfig()))

	replication := []string{"Mode | Node | Role | State | WAL | Remote WAL | Lag | Health | Diagnostics"}
	for _, node := range report.Nodes {
		for _, mode := range node.Replication {
			replication = append(replication, fmt.Sprintf(
				"%s | %s | %s | %s | %s | %s | %s | %s | %s",
				mode.Mode,
				node.Name,
				mode.Role,
				valueOrNA(mode.State),
				formatHealthUint(mode.WAL),
				formatHealthUint(mode.RemoteWAL),
				formatHealthInt(mode.WALLag),
				mode.Health,
				valueOrNA(strings.Join(mode.Diagnostics, "; ")),
			))
		}
	}
	c.UI.Output("")
	if len(replication) == 1 {
		c.UI.Output("Replication is not configured.")
	} else {
		c.UI.Output(tableOutput(replication, columnize.DefaultConfig()))
	}

	if len(report.Errors) > 0 {
		errors := []string{"Address | Check | Error"}
		for _, healthErr := range report.Errors {
			errors = append(errors, fmt.Sprintf(
				"%s | %s | %s",
				healthErr.Address,
				healthErr.Check,
				healthErr.Error,
			))
		}
		c.UI.Output("")
		c.UI.Output(tableOutput(errors, columnize.DefaultConfig()))
	}

	c.UI.Output("")
	c.UI.Output("Overall: " + strings.ToUpper(report.Overall))
}

func formatHealthUint(value *uint64) string {
	if value == nil {
		return "n/a"
	}
	return strconv.FormatUint(*value, 10)
}

func formatHealthInt(value *int64) string {
	if value == nil {
		return "n/a"
	}
	return fmt.Sprintf("%+d", *value)
}

func formatHealthIntPointer(value *int) string {
	if value == nil {
		return "n/a"
	}
	return strconv.Itoa(*value)
}

func valueOrNA(value string) string {
	if value == "" {
		return "n/a"
	}
	return value
}
