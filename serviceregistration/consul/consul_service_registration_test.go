// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

type consulConf map[string]string

func testConsulServiceRegistration(t *testing.T) *serviceRegistration {
	return testConsulServiceRegistrationConfig(t, &consulConf{})
}

func testConsulServiceRegistrationConfig(t *testing.T, conf *consulConf) *serviceRegistration {
	logger := logging.NewVaultLogger(log.Debug)

	shutdownCh := make(chan struct{})
	defer func() {
		close(shutdownCh)
	}()
	be, err := NewServiceRegistration(*conf, logger, sr.State{})
	if err != nil {
		t.Fatalf("Expected Consul to initialize: %v", err)
	}
	if err := be.Run(shutdownCh, &sync.WaitGroup{}, ""); err != nil {
		t.Fatal(err)
	}

	c, ok := be.(*serviceRegistration)
	if !ok {
		t.Fatalf("Expected serviceRegistration")
	}

	return c
}

// TestConsul_ServiceRegistration tests whether consul ServiceRegistration works
func TestConsul_ServiceRegistration(t *testing.T) {
	// Prepare a docker-based consul instance
	cleanup, config := consul.PrepareTestContainer(t, "", false, true)
	defer cleanup()

	// Create a consul client
	client, err := api.NewClient(config.APIConfig())
	if err != nil {
		t.Fatal(err)
	}

	// update the agent's ACL token so that we can successfully deregister the
	// service later in the test
	_, err = client.Agent().UpdateAgentACLToken(config.Token, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Agent().UpdateDefaultACLToken(config.Token, nil)
	if err != nil {
		t.Fatal(err)
	}

	// waitForServices waits for the services in the Consul catalog to
	// reach an expected value, returning the delta if that doesn't happen in time.
	waitForServices := func(t *testing.T, expected map[string][]string) map[string][]string {
		t.Helper()
		// Wait for up to 10 seconds
		var services map[string][]string
		var err error
		for i := 0; i < 10; i++ {
			services, _, err = client.Catalog().Services(nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff := deep.Equal(services, expected); diff == nil {
				return services
			}
			time.Sleep(time.Second)
		}
		t.Fatalf("Catalog Services never reached: got: %v, expected state: %v", services, expected)
		return nil
	}

	shutdownCh := make(chan struct{})
	defer func() {
		close(shutdownCh)
	}()
	const redirectAddr = "http://127.0.0.1:8200"

	// Create a ServiceRegistration that points to our consul instance
	logger := logging.NewVaultLogger(log.Trace)
	srConfig := map[string]string{
		"address": config.Address(),
		"token":   config.Token,
		// decrease reconcile timeout to make test run faster
		"reconcile_timeout": "1s",
	}
	sd, err := NewServiceRegistration(srConfig, logger, sr.State{})
	if err != nil {
		t.Fatal(err)
	}
	if err := sd.Run(shutdownCh, &sync.WaitGroup{}, redirectAddr); err != nil {
		t.Fatal(err)
	}

	// Create the core
	inm, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	core, err := vault.NewCore(&vault.CoreConfig{
		ServiceRegistration: sd,
		Physical:            inm,
		HAPhysical:          inmha.(physical.HABackend),
		RedirectAddr:        redirectAddr,
		DisableMlock:        true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer core.Shutdown()

	waitForServices(t, map[string][]string{
		"consul": {},
		"vault":  {"standby"},
	})

	// Initialize and unseal the core
	keys, _ := vault.TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := vault.TestCoreUnseal(core, vault.TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}
	if core.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Wait for the core to become active
	vault.TestWaitActive(t, core)

	waitForServices(t, map[string][]string{
		"consul": {},
		"vault":  {"active", "initialized"},
	})

	// change the token and trigger reload
	if sd.(*serviceRegistration).config.Token == "" {
		t.Fatal("expected service registration token to not be '' before configuration reload")
	}

	srConfigWithoutToken := make(map[string]string)
	for k, v := range srConfig {
		srConfigWithoutToken[k] = v
	}
	srConfigWithoutToken["token"] = ""
	err = sd.NotifyConfigurationReload(&srConfigWithoutToken)
	if err != nil {
		t.Fatal(err)
	}

	if sd.(*serviceRegistration).config.Token != "" {
		t.Fatal("expected service registration token to be '' after configuration reload")
	}

	// reconfigure the configuration back to its original state and verify vault is registered
	err = sd.NotifyConfigurationReload(&srConfig)
	if err != nil {
		t.Fatal(err)
	}

	waitForServices(t, map[string][]string{
		"consul": {},
		"vault":  {"active", "initialized"},
	})

	// send 'nil' configuration to verify that the service is deregistered
	err = sd.NotifyConfigurationReload(nil)
	if err != nil {
		t.Fatal(err)
	}

	waitForServices(t, map[string][]string{
		"consul": {},
	})

	// reconfigure the configuration back to its original state and verify vault
	// is re-registered
	err = sd.NotifyConfigurationReload(&srConfig)
	if err != nil {
		t.Fatal(err)
	}

	waitForServices(t, map[string][]string{
		"consul": {},
		"vault":  {"active", "initialized"},
	})
}

func TestConsul_ServiceAddress(t *testing.T) {
	tests := []struct {
		consulConfig   map[string]string
		serviceAddrNil bool
	}{
		{
			consulConfig: map[string]string{
				"service_address": "",
			},
		},
		{
			consulConfig: map[string]string{
				"service_address": "vault.example.com",
			},
		},
		{
			serviceAddrNil: true,
		},
	}

	for _, test := range tests {
		shutdownCh := make(chan struct{})
		logger := logging.NewVaultLogger(log.Debug)

		be, err := NewServiceRegistration(test.consulConfig, logger, sr.State{})
		if err != nil {
			t.Fatalf("expected Consul to initialize: %v", err)
		}
		if err := be.Run(shutdownCh, &sync.WaitGroup{}, ""); err != nil {
			t.Fatal(err)
		}

		c, ok := be.(*serviceRegistration)
		if !ok {
			t.Fatalf("Expected ConsulServiceRegistration")
		}

		if test.serviceAddrNil {
			if c.serviceAddress != nil {
				t.Fatalf("expected service address to be nil")
			}
		} else {
			if c.serviceAddress == nil {
				t.Fatalf("did not expect service address to be nil")
			}
		}
		close(shutdownCh)
	}
}

func TestConsul_newConsulServiceRegistration(t *testing.T) {
	tests := []struct {
		name            string
		consulConfig    map[string]string
		fail            bool
		redirectAddr    string
		checkTimeout    time.Duration
		path            string
		service         string
		address         string
		scheme          string
		token           string
		max_parallel    int
		disableReg      bool
		consistencyMode string
	}{
		{
			name:            "Valid default config",
			consulConfig:    map[string]string{},
			checkTimeout:    5 * time.Second,
			redirectAddr:    "http://127.0.0.1:8200",
			path:            "vault/",
			service:         "vault",
			address:         "127.0.0.1:8500",
			scheme:          "http",
			token:           "",
			max_parallel:    4,
			disableReg:      false,
			consistencyMode: "default",
		},
		{
			name: "Valid modified config",
			consulConfig: map[string]string{
				"path":                 "seaTech/",
				"service":              "astronomy",
				"redirect_addr":        "http://127.0.0.2:8200",
				"check_timeout":        "6s",
				"address":              "127.0.0.2",
				"scheme":               "https",
				"token":                "deadbeef-cafeefac-deadc0de-feedface",
				"max_parallel":         "4",
				"disable_registration": "false",
				"consistency_mode":     "strong",
			},
			checkTimeout:    6 * time.Second,
			path:            "seaTech/",
			service:         "astronomy",
			redirectAddr:    "http://127.0.0.2:8200",
			address:         "127.0.0.2",
			scheme:          "https",
			token:           "deadbeef-cafeefac-deadc0de-feedface",
			max_parallel:    4,
			consistencyMode: "strong",
		},
		{
			name: "Unix socket",
			consulConfig: map[string]string{
				"address": "unix:///tmp/.consul.http.sock",
			},
			address: "/tmp/.consul.http.sock",
			scheme:  "http", // Default, not overridden?

			// Defaults
			checkTimeout:    5 * time.Second,
			redirectAddr:    "http://127.0.0.1:8200",
			path:            "vault/",
			service:         "vault",
			token:           "",
			max_parallel:    4,
			disableReg:      false,
			consistencyMode: "default",
		},
		{
			name: "Scheme in address",
			consulConfig: map[string]string{
				"address": "https://127.0.0.2:5000",
			},
			address: "127.0.0.2:5000",
			scheme:  "https",

			// Defaults
			checkTimeout:    5 * time.Second,
			redirectAddr:    "http://127.0.0.1:8200",
			path:            "vault/",
			service:         "vault",
			token:           "",
			max_parallel:    4,
			disableReg:      false,
			consistencyMode: "default",
		},
		{
			name: "check timeout too short",
			fail: true,
			consulConfig: map[string]string{
				"check_timeout": "99ms",
			},
		},
	}

	for _, test := range tests {
		shutdownCh := make(chan struct{})
		logger := logging.NewVaultLogger(log.Debug)

		be, err := NewServiceRegistration(test.consulConfig, logger, sr.State{})
		if test.fail {
			if err == nil {
				t.Fatalf(`Expected config "%s" to fail`, test.name)
			} else {
				continue
			}
		} else if !test.fail && err != nil {
			t.Fatalf("Expected config %s to not fail: %v", test.name, err)
		}
		if err := be.Run(shutdownCh, &sync.WaitGroup{}, ""); err != nil {
			t.Fatal(err)
		}

		c, ok := be.(*serviceRegistration)
		if !ok {
			t.Fatalf("Expected ConsulServiceRegistration: %s", test.name)
		}
		c.disableRegistration = true

		if !c.disableRegistration {
			addr := os.Getenv("CONSUL_HTTP_ADDR")
			if addr == "" {
				continue
			}
		}

		if test.checkTimeout != c.checkTimeout {
			t.Errorf("bad: %v != %v", test.checkTimeout, c.checkTimeout)
		}

		if test.service != c.serviceName {
			t.Errorf("bad: %v != %v", test.service, c.serviceName)
		}

		// The configuration stored in the Consul "client" object is not exported, so
		// we either have to skip validating it, or add a method to export it, or use reflection.
		consulConfig := reflect.Indirect(reflect.ValueOf(c.Client)).FieldByName("config")
		consulConfigScheme := consulConfig.FieldByName("Scheme").String()
		consulConfigAddress := consulConfig.FieldByName("Address").String()

		if test.scheme != consulConfigScheme {
			t.Errorf("bad scheme value: %v != %v", test.scheme, consulConfigScheme)
		}

		if test.address != consulConfigAddress {
			t.Errorf("bad address value: %v != %v", test.address, consulConfigAddress)
		}

		// FIXME(sean@): Unable to test max_parallel
		// if test.max_parallel != cap(c.permitPool) {
		// 	t.Errorf("bad: %v != %v", test.max_parallel, cap(c.permitPool))
		// }
		close(shutdownCh)
	}
}

func TestConsul_serviceTags(t *testing.T) {
	tests := []struct {
		active      bool
		perfStandby bool
		initialized bool
		tags        []string
	}{
		{
			active:      true,
			perfStandby: false,
			initialized: false,
			tags:        []string{"active"},
		},
		{
			active:      false,
			perfStandby: false,
			initialized: false,
			tags:        []string{"standby"},
		},
		{
			active:      false,
			perfStandby: true,
			initialized: false,
			tags:        []string{"performance-standby"},
		},
		{
			active:      true,
			perfStandby: true,
			initialized: false,
			tags:        []string{"performance-standby"},
		},
		{
			active:      true,
			perfStandby: false,
			initialized: true,
			tags:        []string{"active", "initialized"},
		},
		{
			active:      false,
			perfStandby: false,
			initialized: true,
			tags:        []string{"standby", "initialized"},
		},
		{
			active:      false,
			perfStandby: true,
			initialized: true,
			tags:        []string{"performance-standby", "initialized"},
		},
		{
			active:      true,
			perfStandby: true,
			initialized: true,
			tags:        []string{"performance-standby", "initialized"},
		},
	}

	c := testConsulServiceRegistration(t)

	for _, test := range tests {
		tags := c.fetchServiceTags(test.active, test.perfStandby, test.initialized)
		if !reflect.DeepEqual(tags[:], test.tags[:]) {
			t.Errorf("Bad %v: %v %v", test.active, tags, test.tags)
		}
	}
}

func TestConsul_setRedirectAddr(t *testing.T) {
	tests := []struct {
		addr string
		host string
		port int64
		pass bool
	}{
		{
			addr: "http://127.0.0.1:8200/",
			host: "127.0.0.1",
			port: 8200,
			pass: true,
		},
		{
			addr: "http://127.0.0.1:8200",
			host: "127.0.0.1",
			port: 8200,
			pass: true,
		},
		{
			addr: "https://127.0.0.1:8200",
			host: "127.0.0.1",
			port: 8200,
			pass: true,
		},
		{
			addr: "unix:///tmp/.vault.addr.sock",
			host: "/tmp/.vault.addr.sock",
			port: -1,
			pass: true,
		},
		{
			addr: "127.0.0.1:8200",
			pass: false,
		},
		{
			addr: "127.0.0.1",
			pass: false,
		},
	}
	for _, test := range tests {
		c := testConsulServiceRegistration(t)
		err := c.setRedirectAddr(test.addr)
		if test.pass {
			if err != nil {
				t.Fatalf("bad: %v", err)
			}
		} else {
			if err == nil {
				t.Fatalf("bad, expected fail")
			} else {
				continue
			}
		}

		if c.redirectHost != test.host {
			t.Fatalf("bad: %v != %v", c.redirectHost, test.host)
		}

		if c.redirectPort != test.port {
			t.Fatalf("bad: %v != %v", c.redirectPort, test.port)
		}
	}
}

func TestConsul_serviceID(t *testing.T) {
	tests := []struct {
		name         string
		redirectAddr string
		serviceName  string
		expected     string
		valid        bool
	}{
		{
			name:         "valid host w/o slash",
			redirectAddr: "http://127.0.0.1:8200",
			serviceName:  "sea-tech-astronomy",
			expected:     "sea-tech-astronomy:127.0.0.1:8200",
			valid:        true,
		},
		{
			name:         "valid host w/ slash",
			redirectAddr: "http://127.0.0.1:8200/",
			serviceName:  "sea-tech-astronomy",
			expected:     "sea-tech-astronomy:127.0.0.1:8200",
			valid:        true,
		},
		{
			name:         "valid https host w/ slash",
			redirectAddr: "https://127.0.0.1:8200/",
			serviceName:  "sea-tech-astronomy",
			expected:     "sea-tech-astronomy:127.0.0.1:8200",
			valid:        true,
		},
		{
			name:         "invalid host name",
			redirectAddr: "https://127.0.0.1:8200/",
			serviceName:  "sea_tech_astronomy",
			expected:     "",
			valid:        false,
		},
	}

	logger := logging.NewVaultLogger(log.Debug)

	for _, test := range tests {
		shutdownCh := make(chan struct{})
		be, err := NewServiceRegistration(consulConf{
			"service": test.serviceName,
		}, logger, sr.State{})
		if !test.valid {
			if err == nil {
				t.Fatalf("expected an error initializing for name %q", test.serviceName)
			}
			continue
		}
		if test.valid && err != nil {
			t.Fatalf("expected Consul to initialize: %v", err)
		}
		if err := be.Run(shutdownCh, &sync.WaitGroup{}, ""); err != nil {
			t.Fatal(err)
		}

		c, ok := be.(*serviceRegistration)
		if !ok {
			t.Fatalf("Expected serviceRegistration")
		}

		if err := c.setRedirectAddr(test.redirectAddr); err != nil {
			t.Fatalf("bad: %s %v", test.name, err)
		}

		serviceID := c.serviceID()
		if serviceID != test.expected {
			t.Fatalf("bad: %v != %v", serviceID, test.expected)
		}
	}
}

// TestConsul_NewServiceRegistration_serviceTags ensures that we do not modify
// the case of any 'service_tags' set by the config.
// We do expect tags to be sorted in lexicographic order (A-Z).
func TestConsul_NewServiceRegistration_serviceTags(t *testing.T) {
	tests := map[string]struct {
		Tags         string
		ExpectedTags []string
	}{
		"lowercase": {
			Tags:         "foo,bar",
			ExpectedTags: []string{"bar", "foo"},
		},
		"uppercase": {
			Tags:         "FOO,BAR",
			ExpectedTags: []string{"BAR", "FOO"},
		},
		"PascalCase": {
			Tags:         "FooBar, Feedface",
			ExpectedTags: []string{"Feedface", "FooBar"},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg := map[string]string{"service_tags": tc.Tags}
			logger := logging.NewVaultLogger(log.Trace)
			be, err := NewServiceRegistration(cfg, logger, sr.State{})
			require.NoError(t, err)
			require.NotNil(t, be)
			c, ok := be.(*serviceRegistration)
			require.True(t, ok)
			require.NotNil(t, c)
			require.Equal(t, tc.ExpectedTags, c.serviceTags)
		})
	}
}
