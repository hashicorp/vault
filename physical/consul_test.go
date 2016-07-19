package physical

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)

type consulConf map[string]string

var (
	addrCount int = 0
)

func testHostIP() string {
	a := addrCount
	addrCount++
	return fmt.Sprintf("127.0.0.%d", a)
}

func testConsulBackend(t *testing.T) *ConsulBackend {
	return testConsulBackendConfig(t, &consulConf{})
}

func testConsulBackendConfig(t *testing.T, conf *consulConf) *ConsulBackend {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	be, err := newConsulBackend(*conf, logger)
	if err != nil {
		t.Fatalf("Expected Consul to initialize: %v", err)
	}

	c, ok := be.(*ConsulBackend)
	if !ok {
		t.Fatalf("Expected ConsulBackend")
	}

	return c
}

func testConsul_testConsulBackend(t *testing.T) {
	c := testConsulBackend(t)
	if c == nil {
		t.Fatalf("bad")
	}
}

func testActiveFunc(activePct float64) activeFunction {
	return func() bool {
		var active bool
		standbyProb := rand.Float64()
		if standbyProb > activePct {
			active = true
		}
		return active
	}
}

func testSealedFunc(sealedPct float64) sealedFunction {
	return func() bool {
		var sealed bool
		unsealedProb := rand.Float64()
		if unsealedProb > sealedPct {
			sealed = true
		}
		return sealed
	}
}

func TestConsul_newConsulBackend(t *testing.T) {
	tests := []struct {
		name          string
		consulConfig  map[string]string
		fail          bool
		advertiseAddr string
		checkTimeout  time.Duration
		path          string
		service       string
		address       string
		scheme        string
		token         string
		max_parallel  int
		disableReg    bool
	}{
		{
			name:          "Valid default config",
			consulConfig:  map[string]string{},
			checkTimeout:  5 * time.Second,
			advertiseAddr: "http://127.0.0.1:8200",
			path:          "vault/",
			service:       "vault",
			address:       "127.0.0.1:8500",
			scheme:        "http",
			token:         "",
			max_parallel:  4,
			disableReg:    false,
		},
		{
			name: "Valid modified config",
			consulConfig: map[string]string{
				"path":                 "seaTech/",
				"service":              "astronomy",
				"advertiseAddr":        "http://127.0.0.2:8200",
				"check_timeout":        "6s",
				"address":              "127.0.0.2",
				"scheme":               "https",
				"token":                "deadbeef-cafeefac-deadc0de-feedface",
				"max_parallel":         "4",
				"disable_registration": "false",
			},
			checkTimeout:  6 * time.Second,
			path:          "seaTech/",
			service:       "astronomy",
			advertiseAddr: "http://127.0.0.2:8200",
			address:       "127.0.0.2",
			scheme:        "https",
			token:         "deadbeef-cafeefac-deadc0de-feedface",
			max_parallel:  4,
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
		logger := log.New(os.Stderr, "", log.LstdFlags)
		be, err := newConsulBackend(test.consulConfig, logger)
		if test.fail {
			if err == nil {
				t.Fatalf(`Expected config "%s" to fail`, test.name)
			} else {
				continue
			}
		} else if !test.fail && err != nil {
			t.Fatalf("Expected config %s to not fail: %v", test.name, err)
		}

		c, ok := be.(*ConsulBackend)
		if !ok {
			t.Fatalf("Expected ConsulBackend: %s", test.name)
		}
		c.disableRegistration = true

		if c.disableRegistration == false {
			addr := os.Getenv("CONSUL_HTTP_ADDR")
			if addr == "" {
				continue
			}
		}

		var shutdownCh ShutdownChannel
		if err := c.RunServiceDiscovery(shutdownCh, test.advertiseAddr, testActiveFunc(0.5), testSealedFunc(0.5)); err != nil {
			t.Fatalf("bad: %v", err)
		}

		if test.checkTimeout != c.checkTimeout {
			t.Errorf("bad: %v != %v", test.checkTimeout, c.checkTimeout)
		}

		if test.path != c.path {
			t.Errorf("bad: %s %v != %v", test.name, test.path, c.path)
		}

		if test.service != c.serviceName {
			t.Errorf("bad: %v != %v", test.service, c.serviceName)
		}

		// FIXME(sean@): Unable to test max_parallel
		// if test.max_parallel != cap(c.permitPool) {
		// 	t.Errorf("bad: %v != %v", test.max_parallel, cap(c.permitPool))
		// }
	}
}

func TestConsul_serviceTags(t *testing.T) {
	tests := []struct {
		active bool
		tags   []string
	}{
		{
			active: true,
			tags:   []string{"active"},
		},
		{
			active: false,
			tags:   []string{"standby"},
		},
	}

	for _, test := range tests {
		tags := serviceTags(test.active)
		if !reflect.DeepEqual(tags[:], test.tags[:]) {
			t.Errorf("Bad %v: %v %v", test.active, tags, test.tags)
		}
	}
}

func TestConsul_setAdvertiseAddr(t *testing.T) {
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
		c := testConsulBackend(t)
		err := c.setAdvertiseAddr(test.addr)
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

		if c.advertiseHost != test.host {
			t.Fatalf("bad: %v != %v", c.advertiseHost, test.host)
		}

		if c.advertisePort != test.port {
			t.Fatalf("bad: %v != %v", c.advertisePort, test.port)
		}
	}
}

func TestConsul_NotifyActiveStateChange(t *testing.T) {
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	if addr == "" {
		t.Skipf("No consul process running, skipping test")
	}

	c := testConsulBackend(t)

	if err := c.NotifyActiveStateChange(); err != nil {
		t.Fatalf("bad: %v", err)
	}
}

func TestConsul_NotifySealedStateChange(t *testing.T) {
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	if addr == "" {
		t.Skipf("No consul process running, skipping test")
	}

	c := testConsulBackend(t)

	if err := c.NotifySealedStateChange(); err != nil {
		t.Fatalf("bad: %v", err)
	}
}

func TestConsul_serviceID(t *testing.T) {
	passingTests := []struct {
		name          string
		advertiseAddr string
		serviceName   string
		expected      string
	}{
		{
			name:          "valid host w/o slash",
			advertiseAddr: "http://127.0.0.1:8200",
			serviceName:   "sea-tech-astronomy",
			expected:      "sea-tech-astronomy:127.0.0.1:8200",
		},
		{
			name:          "valid host w/ slash",
			advertiseAddr: "http://127.0.0.1:8200/",
			serviceName:   "sea-tech-astronomy",
			expected:      "sea-tech-astronomy:127.0.0.1:8200",
		},
		{
			name:          "valid https host w/ slash",
			advertiseAddr: "https://127.0.0.1:8200/",
			serviceName:   "sea-tech-astronomy",
			expected:      "sea-tech-astronomy:127.0.0.1:8200",
		},
	}

	for _, test := range passingTests {
		c := testConsulBackendConfig(t, &consulConf{
			"service": test.serviceName,
		})

		if err := c.setAdvertiseAddr(test.advertiseAddr); err != nil {
			t.Fatalf("bad: %s %v", test.name, err)
		}

		serviceID := c.serviceID()
		if serviceID != test.expected {
			t.Fatalf("bad: %v != %v", serviceID, test.expected)
		}
	}
}

func TestConsulBackend(t *testing.T) {
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	if addr == "" {
		t.Skipf("No consul process running, skipping test")
	}

	conf := api.DefaultConfig()
	conf.Address = addr
	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	b, err := NewBackend("consul", logger, map[string]string{
		"address":      addr,
		"path":         randPath,
		"max_parallel": "256",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}

func TestConsulHABackend(t *testing.T) {
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	if addr == "" {
		t.Skipf("No consul process running, skipping test")
	}

	conf := api.DefaultConfig()
	conf.Address = addr
	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	b, err := NewBackend("consul", logger, map[string]string{
		"address":      addr,
		"path":         randPath,
		"max_parallel": "-1",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	ha, ok := b.(HABackend)
	if !ok {
		t.Fatalf("consul does not implement HABackend")
	}
	testHABackend(t, ha, ha)

	detect, ok := b.(AdvertiseDetect)
	if !ok {
		t.Fatalf("consul does not implement AdvertiseDetect")
	}
	host, err := detect.DetectHostAddr()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if host == "" {
		t.Fatalf("bad addr: %v", host)
	}
}
