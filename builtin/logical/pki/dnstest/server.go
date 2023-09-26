// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dnstest

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/stretchr/testify/require"
)

type TestServer struct {
	t   *testing.T
	ctx context.Context
	log hclog.Logger

	runner  *docker.Runner
	network string
	startup *docker.Service

	lock       sync.Mutex
	serial     int
	forwarders []string
	domains    []string
	records    map[string]map[string][]string // domain -> record -> value(s).

	cleanup func()
}

func SetupResolver(t *testing.T, domain string) *TestServer {
	return SetupResolverOnNetwork(t, domain, "")
}

func SetupResolverOnNetwork(t *testing.T, domain string, network string) *TestServer {
	var ts TestServer
	ts.t = t
	ts.ctx = context.Background()
	ts.domains = []string{domain}
	ts.records = map[string]map[string][]string{}
	ts.network = network
	ts.log = hclog.L()

	ts.setupRunner(domain, network)
	ts.startContainer(network)
	ts.PushConfig()

	return &ts
}

func (ts *TestServer) setupRunner(domain string, network string) {
	var err error
	ts.runner, err = docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "ubuntu/bind9",
		ImageTag:      "latest",
		ContainerName: "bind9-dns-" + strings.ReplaceAll(domain, ".", "-"),
		NetworkName:   network,
		Ports:         []string{"53/udp"},
		// DNS container logging was disabled to reduce content within CI logs.
		//LogConsumer: func(s string) {
		//	ts.log.Info(s)
		//},
	})
	require.NoError(ts.t, err)
}

func (ts *TestServer) startContainer(network string) {
	connUpFunc := func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		// Perform a simple connection to this resolver, even though the
		// default configuration doesn't do anything useful.
		peer, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return nil, fmt.Errorf("failed to resolve peer: %v / %v: %w", host, port, err)
		}

		conn, err := net.DialUDP("udp", nil, peer)
		if err != nil {
			return nil, fmt.Errorf("failed to dial peer: %v / %v / %v: %w", host, port, peer, err)
		}
		defer conn.Close()

		_, err = conn.Write([]byte("garbage-in"))
		if err != nil {
			return nil, fmt.Errorf("failed to write to peer: %v / %v / %v: %w", host, port, peer, err)
		}

		// Connection worked.
		return docker.NewServiceHostPort(host, port), nil
	}

	result, _, err := ts.runner.StartNewService(ts.ctx, true, true, connUpFunc)
	require.NoError(ts.t, err, "failed to start dns resolver for "+ts.domains[0])
	ts.startup = result

	if ts.startup.StartResult.RealIP == "" {
		mapping, err := ts.runner.GetNetworkAndAddresses(ts.startup.Container.ID)
		require.NoError(ts.t, err, "failed to fetch network addresses to correct missing real IP address")
		if len(network) == 0 {
			require.Equal(ts.t, 1, len(mapping), "expected exactly one network address")
			for network = range mapping {
				// Because mapping is a map of network name->ip, we need
				// to use the above range's assignment to get the name,
				// as there is no other way of getting the keys of a map.
			}
		}
		require.Contains(ts.t, mapping, network, "expected network to be part of the mapping")
		ts.startup.StartResult.RealIP = mapping[network]
	}

	ts.log.Info(fmt.Sprintf("[dnsserv] Addresses of DNS resolver: local=%v / container=%v", ts.GetLocalAddr(), ts.GetRemoteAddr()))
}

func (ts *TestServer) buildNamedConf() string {
	forwarders := "\n"
	if len(ts.forwarders) > 0 {
		forwarders = "\tforwarders {\n"
		for _, forwarder := range ts.forwarders {
			forwarders += "\t\t" + forwarder + ";\n"
		}
		forwarders += "\t};\n"
	}

	zones := "\n"
	for _, domain := range ts.domains {
		zones += fmt.Sprintf("zone \"%s\" {\n", domain)
		zones += "\ttype primary;\n"
		zones += fmt.Sprintf("\tfile \"%s.zone\";\n", domain)
		zones += "\tallow-update {\n\t\tnone;\n\t};\n"
		zones += "\tnotify no;\n"
		zones += "};\n\n"
	}

	// Reverse lookups are not handles as they're not presently necessary.

	cfg := `options {
	directory "/var/cache/bind";

	dnssec-validation no;

	` + forwarders + `
};

` + zones

	return cfg
}

func (ts *TestServer) buildZoneFile(target string) string {
	// One second TTL by default to allow quick refreshes.
	zone := "$TTL 1;\n"

	ts.serial += 1
	zone += fmt.Sprintf("@\tIN\tSOA\tns.%v.\troot.%v.\t(\n", target, target)
	zone += fmt.Sprintf("\t\t\t%d;\n\t\t\t1;\n\t\t\t1;\n\t\t\t2;\n\t\t\t1;\n\t\t\t)\n\n", ts.serial)
	zone += fmt.Sprintf("@\tIN\tNS\tns%d.%v.\n", ts.serial, target)
	zone += fmt.Sprintf("ns%d.%v.\tIN\tA\t%v\n", ts.serial, target, "127.0.0.1")

	for domain, records := range ts.records {
		if !strings.HasSuffix(domain, target) {
			continue
		}

		for recordType, values := range records {
			for _, value := range values {
				zone += fmt.Sprintf("%s.\tIN\t%s\t%s\n", domain, recordType, value)
			}
		}
	}

	return zone
}

func (ts *TestServer) pushNamedConf() {
	contents := docker.NewBuildContext()
	cfgPath := "/etc/bind/named.conf.options"
	namedCfg := ts.buildNamedConf()
	contents[cfgPath] = docker.PathContentsFromString(namedCfg)
	contents[cfgPath].SetOwners(0, 142) // root, bind

	ts.log.Info(fmt.Sprintf("Generated bind9 config (%s):\n%v\n", cfgPath, namedCfg))

	err := ts.runner.CopyTo(ts.startup.Container.ID, "/", contents)
	require.NoError(ts.t, err, "failed pushing updated named.conf.options to container")
}

func (ts *TestServer) pushZoneFiles() {
	contents := docker.NewBuildContext()

	for _, domain := range ts.domains {
		path := "/var/cache/bind/" + domain + ".zone"
		zoneFile := ts.buildZoneFile(domain)
		contents[path] = docker.PathContentsFromString(zoneFile)
		contents[path].SetOwners(0, 142) // root, bind

		ts.log.Info(fmt.Sprintf("Generated bind9 zone file for %v (%s):\n%v\n", domain, path, zoneFile))
	}

	err := ts.runner.CopyTo(ts.startup.Container.ID, "/", contents)
	require.NoError(ts.t, err, "failed pushing updated named.conf.options to container")
}

func (ts *TestServer) PushConfig() {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	_, _, _, err := ts.runner.RunCmdWithOutput(ts.ctx, ts.startup.Container.ID, []string{"rndc", "freeze"})
	require.NoError(ts.t, err, "failed to freeze DNS config")

	// There's two cases here:
	//
	// 1. We've added a new top-level domain name. Here, we want to make
	//    sure the new zone file is pushed before we push the reference
	//    to it.
	// 2. We've just added a new. Here, the order doesn't matter, but
	//    mostly likely the second push will be a no-op.
	ts.pushZoneFiles()
	ts.pushNamedConf()

	_, _, _, err = ts.runner.RunCmdWithOutput(ts.ctx, ts.startup.Container.ID, []string{"rndc", "thaw"})
	require.NoError(ts.t, err, "failed to thaw DNS config")

	// Wait until our config has taken.
	corehelpers.RetryUntil(ts.t, 15*time.Second, func() error {
		// bind reloads based on file mtime, touch files before starting
		// to make sure it has been updated more recently than when the
		// last update was written. Then issue a new SIGHUP.
		for _, domain := range ts.domains {
			path := "/var/cache/bind/" + domain + ".zone"
			touchCmd := []string{"touch", path}

			_, _, _, err := ts.runner.RunCmdWithOutput(ts.ctx, ts.startup.Container.ID, touchCmd)
			if err != nil {
				return fmt.Errorf("failed to update zone mtime: %w", err)
			}
		}
		ts.runner.DockerAPI.ContainerKill(ts.ctx, ts.startup.Container.ID, "SIGHUP")

		// Connect to our bind resolver.
		resolver := &net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 10 * time.Second,
				}
				return d.DialContext(ctx, network, ts.GetLocalAddr())
			},
		}

		// last domain has the given serial number, which also appears in the
		// NS record so we can fetch it via Go.
		lastDomain := ts.domains[len(ts.domains)-1]
		records, err := resolver.LookupNS(ts.ctx, lastDomain)
		if err != nil {
			return fmt.Errorf("failed to lookup NS record for %v: %w", lastDomain, err)
		}

		if len(records) != 1 {
			return fmt.Errorf("expected only 1 NS record for %v, got %v/%v", lastDomain, len(records), records)
		}

		expectedNS := fmt.Sprintf("ns%d.%v.", ts.serial, lastDomain)
		if records[0].Host != expectedNS {
			return fmt.Errorf("expected to find NS %v, got %v indicating reload hadn't completed", expectedNS, records[0])
		}

		return nil
	})
}

func (ts *TestServer) GetLocalAddr() string {
	return ts.startup.Config.Address()
}

func (ts *TestServer) GetRemoteAddr() string {
	return fmt.Sprintf("%s:%d", ts.startup.StartResult.RealIP, 53)
}

func (ts *TestServer) AddDomain(domain string) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	for _, existing := range ts.domains {
		if existing == domain {
			return
		}
	}

	ts.domains = append(ts.domains, domain)
}

func (ts *TestServer) AddRecord(domain string, record string, value string) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	foundDomain := false
	for _, existing := range ts.domains {
		if strings.HasSuffix(domain, existing) {
			foundDomain = true
			break
		}
	}
	if !foundDomain {
		ts.t.Fatalf("cannot add record %v/%v :: [%v] -- no domain zone matching (%v)", record, domain, value, ts.domains)
	}

	value = strings.TrimSpace(value)
	if _, present := ts.records[domain]; !present {
		ts.records[domain] = map[string][]string{}
	}

	if values, present := ts.records[domain][record]; present {
		for _, candidate := range values {
			if candidate == value {
				// Already present; skip adding.
				return
			}
		}
	}

	ts.records[domain][record] = append(ts.records[domain][record], value)
}

func (ts *TestServer) RemoveRecord(domain string, record string, value string) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	foundDomain := false
	for _, existing := range ts.domains {
		if strings.HasSuffix(domain, existing) {
			foundDomain = true
			break
		}
	}
	if !foundDomain {
		// Not found.
		return
	}

	value = strings.TrimSpace(value)
	if _, present := ts.records[domain]; !present {
		// Not found.
		return
	}

	var remaining []string
	if values, present := ts.records[domain][record]; present {
		for _, candidate := range values {
			if candidate != value {
				remaining = append(remaining, candidate)
			}
		}
	}

	ts.records[domain][record] = remaining
}

func (ts *TestServer) RemoveRecordsOfTypeForDomain(domain string, record string) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	foundDomain := false
	for _, existing := range ts.domains {
		if strings.HasSuffix(domain, existing) {
			foundDomain = true
			break
		}
	}
	if !foundDomain {
		// Not found.
		return
	}

	if _, present := ts.records[domain]; !present {
		// Not found.
		return
	}

	delete(ts.records[domain], record)
}

func (ts *TestServer) RemoveRecordsForDomain(domain string) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	foundDomain := false
	for _, existing := range ts.domains {
		if strings.HasSuffix(domain, existing) {
			foundDomain = true
			break
		}
	}
	if !foundDomain {
		// Not found.
		return
	}

	if _, present := ts.records[domain]; !present {
		// Not found.
		return
	}

	ts.records[domain] = map[string][]string{}
}

func (ts *TestServer) RemoveAllRecords() {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	ts.records = map[string]map[string][]string{}
}

func (ts *TestServer) Cleanup() {
	if ts.cleanup != nil {
		ts.cleanup()
	}
	if ts.startup != nil && ts.startup.Cleanup != nil {
		ts.startup.Cleanup()
	}
}
