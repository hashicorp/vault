package consul

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/sdk/version"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/vault"
)

const redirectAddr = "http://127.0.0.1:8200"

type consulConf map[string]string

// TestConsul_ServiceRegistration tests whether consul ServiceRegistration works
func TestConsul_ServiceRegistration(t *testing.T) {

	// Prepare a docker-based consul instance
	cleanup, addr, token := consul.PrepareTestContainer(t, "1.4.0-rc1")
	defer cleanup()

	// Create a consul client
	cfg := api.DefaultConfig()
	cfg.Address = addr
	cfg.Token = token
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Vault should not yet be registered with Consul
	services, _, err := client.Catalog().Services(nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(services, map[string][]string{
		"consul": {},
	}); diff != nil {
		t.Fatal(diff)
	}

	// Create a ServiceRegistration that points to our consul instance
	logger := logging.NewVaultLogger(log.Trace)
	sd, err := NewServiceRegistration(make(chan struct{}), map[string]string{
		"address": addr,
		"token":   token,
	}, logger, initialState(), redirectAddr)
	if err != nil {
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

	// Vault should soon be registered with Consul in standby mode
	services, _, err = client.Catalog().Services(nil)
	if err != nil {
		t.Fatal(err)
	}
	resultingTags := services["vault"]
	for _, tag := range []string{tagPerfStandby, tagNotActive, tagUninitialized, tagSealed} {
		if !strutil.StrListContains(resultingTags, tag) {
			t.Fatalf("expected %q but received %q", tag, resultingTags)
		}
	}

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

	// Vault should soon be registered with Consul in active mode
	services, _, err = client.Catalog().Services(nil)
	if err != nil {
		t.Fatal(err)
	}
	resultingTags = services["vault"]
	for _, tag := range []string{tagPerfStandby, tagIsActive, tagInitialized, tagUnsealed} {
		if !strutil.StrListContains(resultingTags, tag) {
			t.Fatalf("expected %q but received %q", tag, resultingTags)
		}
	}
}

func TestConsul_ServiceTags(t *testing.T) {

	usersTags := []string{"deadbeef", "cafeefac", "deadc0de", "feedface"}
	negativeState := &sr.State{
		VaultVersion:         "",
		IsInitialized:        false,
		IsSealed:             false,
		IsActive:             false,
		IsPerformanceStandby: false,
	}
	positiveState := &sr.State{
		VaultVersion:         "some-version",
		IsInitialized:        true,
		IsSealed:             true,
		IsActive:             true,
		IsPerformanceStandby: true,
	}

	actual := buildTags(nil, negativeState)
	if !strutil.StrListContains(actual, tagNotPerfStandby) {
		t.Fatalf("expected %q but received %q", tagNotPerfStandby, actual)
	}
	if !strutil.StrListContains(actual, tagNotActive) {
		t.Fatalf("expected %q but received %q", tagNotActive, actual)
	}
	if !strutil.StrListContains(actual, tagUninitialized) {
		t.Fatalf("expected %q but received %q", tagUninitialized, actual)
	}
	if !strutil.StrListContains(actual, tagUnsealed) {
		t.Fatalf("expected %q but received %q", tagUnsealed, actual)
	}

	actual = buildTags(nil, positiveState)
	if !strutil.StrListContains(actual, tagPerfStandby) {
		t.Fatalf("expected %q but received %q", tagPerfStandby, actual)
	}
	if !strutil.StrListContains(actual, tagIsActive) {
		t.Fatalf("expected %q but received %q", tagIsActive, actual)
	}
	if !strutil.StrListContains(actual, tagInitialized) {
		t.Fatalf("expected %q but received %q", tagInitialized, actual)
	}
	if !strutil.StrListContains(actual, tagSealed) {
		t.Fatalf("expected %q but received %q", tagSealed, actual)
	}

	actual = buildTags(usersTags, negativeState)
	if !strutil.StrListContains(actual, tagNotPerfStandby) {
		t.Fatalf("expected %q but received %q", tagNotPerfStandby, actual)
	}
	if !strutil.StrListContains(actual, tagNotActive) {
		t.Fatalf("expected %q but received %q", tagNotActive, actual)
	}
	if !strutil.StrListContains(actual, tagUninitialized) {
		t.Fatalf("expected %q but received %q", tagUninitialized, actual)
	}
	if !strutil.StrListContains(actual, tagUnsealed) {
		t.Fatalf("expected %q but received %q", tagUnsealed, actual)
	}
	for _, tag := range usersTags {
		if !strutil.StrListContains(actual, tag) {
			t.Fatalf("expected %q but received %q", tag, actual)
		}
	}

	actual = buildTags(usersTags, positiveState)
	if !strutil.StrListContains(actual, tagPerfStandby) {
		t.Fatalf("expected %q but received %q", tagPerfStandby, actual)
	}
	if !strutil.StrListContains(actual, tagIsActive) {
		t.Fatalf("expected %q but received %q", tagIsActive, actual)
	}
	if !strutil.StrListContains(actual, tagInitialized) {
		t.Fatalf("expected %q but received %q", tagInitialized, actual)
	}
	if !strutil.StrListContains(actual, tagSealed) {
		t.Fatalf("expected %q but received %q", tagSealed, actual)
	}
	for _, tag := range usersTags {
		if !strutil.StrListContains(actual, tag) {
			t.Fatalf("expected %q but received %q", tag, actual)
		}
	}
}

func TestConsul_newConsulServiceRegistration(t *testing.T) {
	// Prepare a docker-based consul instance
	cleanup, addr, token := consul.PrepareTestContainer(t, "1.4.0-rc1")
	defer cleanup()

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
				"scheme":               "http",
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
			scheme:          "http",
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
			scheme:  "http",

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
		logger := logging.NewVaultLogger(log.Debug)

		test.consulConfig["address"] = addr
		test.consulConfig["token"] = token
		be, err := NewServiceRegistration(make(chan struct{}), test.consulConfig, logger, initialState(), redirectAddr)
		if test.fail {
			if err == nil {
				t.Fatalf(`Expected config "%s" to fail`, test.name)
			} else {
				continue
			}
		} else if !test.fail && err != nil {
			t.Fatalf("Expected config %s to not fail: %v", test.name, err)
		}

		c, ok := be.(*ServiceRegistration)
		if !ok {
			t.Fatalf("Expected ServiceRegistration: %s", test.name)
		}
		c.disableRegistration = true

		if c.disableRegistration == false {
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

		if addr != consulConfigAddress {
			t.Errorf("bad address value: %v != %v", test.address, consulConfigAddress)
		}

		// FIXME(sean@): Unable to test max_parallel
		// if test.max_parallel != cap(c.permitPool) {
		// 	t.Errorf("bad: %v != %v", test.max_parallel, cap(c.permitPool))
		// }
	}
}

func TestConsul_setRedirectAddr(t *testing.T) {
	// Prepare a docker-based consul instance
	cleanup, addr, token := consul.PrepareTestContainer(t, "1.4.0-rc1")
	defer cleanup()
	logger := logging.NewVaultLogger(log.Debug)

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
		be, err := NewServiceRegistration(make(chan struct{}), map[string]string{
			"address": addr,
			"token":   token,
		}, logger, initialState(), test.addr)
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
		c, ok := be.(*ServiceRegistration)
		if !ok {
			t.Fatalf("Expected ServiceRegistration")
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
	// Prepare a docker-based consul instance
	cleanup, addr, token := consul.PrepareTestContainer(t, "1.4.0-rc1")
	defer cleanup()

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
		be, err := NewServiceRegistration(make(chan struct{}), consulConf{
			"service": test.serviceName,
			"address": addr,
			"token":   token,
		}, logger, initialState(), redirectAddr)
		if !test.valid {
			if err == nil {
				t.Fatalf("expected an error initializing for name %q", test.serviceName)
			}
			continue
		}
		if test.valid && err != nil {
			t.Fatalf("expected Consul to initialize: %v", err)
		}

		c, ok := be.(*ServiceRegistration)
		if !ok {
			t.Fatalf("Expected ServiceRegistration")
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

func initialState() *sr.State {
	return &sr.State{
		VaultVersion:         version.GetVersion().VersionNumber(),
		IsInitialized:        false,
		IsSealed:             true,
		IsActive:             false,
		IsPerformanceStandby: true,
	}
}
