package physical

import (
	"fmt"
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
	const serviceID = "vaultTestService"
	be, err := newConsulBackend(*conf)
	if err != nil {
		t.Fatalf("Expected Consul to initialize: %v", err)
	}

	c, ok := be.(*ConsulBackend)
	if !ok {
		t.Fatalf("Expected ConsulBackend")
	}

	c.service = &api.AgentServiceRegistration{
		ID:                serviceID,
		Name:              c.serviceName,
		Tags:              serviceTags(c.active),
		Port:              8200,
		Address:           testHostIP(),
		EnableTagOverride: false,
	}

	c.sealedCheck = &api.AgentCheckRegistration{
		ID:        c.checkID(),
		Name:      "Vault Sealed Status",
		Notes:     "Vault service is healthy when Vault is in an unsealed status and can become an active Vault server",
		ServiceID: serviceID,
		AgentServiceCheck: api.AgentServiceCheck{
			TTL:    c.checkTimeout.String(),
			Status: api.HealthPassing,
		},
	}

	return c
}

func testConsul_testConsulBackend(t *testing.T) {
	c := testConsulBackend(t)
	if c == nil {
		t.Fatalf("bad")
	}

	if c.active != false {
		t.Fatalf("bad")
	}

	if c.sealed != false {
		t.Fatalf("bad")
	}

	if c.service == nil {
		t.Fatalf("bad")
	}

	if c.sealedCheck == nil {
		t.Fatalf("bad")
	}
}

func TestConsul_newConsulBackend(t *testing.T) {
	tests := []struct {
		Name         string
		Config       map[string]string
		Fail         bool
		checkTimeout time.Duration
		path         string
		service      string
		address      string
		scheme       string
		token        string
		max_parallel int
	}{
		{
			Name:         "Valid default config",
			Config:       map[string]string{},
			checkTimeout: 5 * time.Second,
			path:         "vault",
			service:      "vault",
			address:      "127.0.0.1",
			scheme:       "http",
			token:        "",
			max_parallel: 4,
		},
		{
			Name: "Valid modified config",
			Config: map[string]string{
				"path":          "seaTech/",
				"service":       "astronomy",
				"check_timeout": "6s",
				"address":       "127.0.0.2",
				"scheme":        "https",
				"token":         "deadbeef-cafeefac-deadc0de-feedface",
				"max_parallel":  "4",
			},
			checkTimeout: 6 * time.Second,
			path:         "seaTech/",
			service:      "astronomy",
			address:      "127.0.0.2",
			scheme:       "https",
			token:        "deadbeef-cafeefac-deadc0de-feedface",
			max_parallel: 4,
		},
		{
			Name: "check timeout too short",
			Fail: true,
			Config: map[string]string{
				"check_timeout": "99ms",
			},
		},
	}

	for _, test := range tests {
		be, err := newConsulBackend(test.Config)
		if test.Fail && err == nil {
			t.Fatalf("Expected config %s to fail", test.Name)
		} else if !test.Fail && err != nil {
			t.Fatalf("Expected config %s to not fail: %v", test.Name, err)
		}

		c, ok := be.(*ConsulBackend)
		if !ok {
			t.Fatalf("Expected ConsulBackend")
		}

		if test.checkTimeout != c.checkTimeout {
			t.Errorf("bad: %v != %v", test.checkTimeout, c.checkTimeout)
		}

		if test.path != c.path {
			t.Errorf("bad: %v != %v", test.path, c.path)
		}

		if test.service != c.serviceName {
			t.Errorf("bad: %v != %v", test.service, c.serviceName)
		}

		if test.address != c.consulClientConf.Address {
			t.Errorf("bad: %v != %v", test.address, c.consulClientConf.Address)
		}

		if test.scheme != c.consulClientConf.Scheme {
			t.Errorf("bad: %v != %v", test.scheme, c.consulClientConf.Scheme)
		}

		if test.token != c.consulClientConf.Token {
			t.Errorf("bad: %v != %v", test.token, c.consulClientConf.Token)
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

func TestConsul_UpdateAdvertiseAddr(t *testing.T) {
	tests := []struct {
		addr string
		pass bool
	}{
		{
			addr: "http://127.0.0.1:8200/",
			pass: true,
		},
		{
			addr: "http://127.0.0.1:8200",
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
		if c == nil {
			t.Fatalf("bad")
		}

		err := c.UpdateAdvertiseAddr(test.addr)
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

		if c.advertiseAddr != test.addr {
			t.Fatalf("bad: %v != %v", c.advertiseAddr, test.addr)
		}
	}
}

func TestConsul_AdvertiseActive(t *testing.T) {
	c := testConsulBackend(t)

	if c.active != false {
		t.Fatalf("bad")
	}

	if err := c.AdvertiseActive(true); err != nil {
		t.Fatalf("bad: %v", err)
	}

	if err := c.AdvertiseActive(true); err != nil {
		t.Fatalf("bad: %v", err)
	}

	if err := c.AdvertiseActive(false); err != nil {
		t.Fatalf("bad: %v", err)
	}

	if err := c.AdvertiseActive(false); err != nil {
		t.Fatalf("bad: %v", err)
	}

	if err := c.AdvertiseActive(true); err != nil {
		t.Fatalf("bad: %v", err)
	}
}

func TestConsul_AdvertiseSealed(t *testing.T) {
	c := testConsulBackend(t)

	if c.sealed != false {
		t.Fatalf("bad")
	}

	if err := c.AdvertiseSealed(true); err != nil {
		t.Fatalf("bad: %v", err)
	}
	if c.sealed != true {
		t.Fatalf("bad")
	}

	if err := c.AdvertiseSealed(true); err != nil {
		t.Fatalf("bad: %v", err)
	}
	if c.sealed != true {
		t.Fatalf("bad")
	}

	if err := c.AdvertiseSealed(false); err != nil {
		t.Fatalf("bad: %v", err)
	}
	if c.sealed != false {
		t.Fatalf("bad")
	}

	if err := c.AdvertiseSealed(false); err != nil {
		t.Fatalf("bad: %v", err)
	}
	if c.sealed != false {
		t.Fatalf("bad")
	}

	if err := c.AdvertiseSealed(true); err != nil {
		t.Fatalf("bad: %v", err)
	}
	if c.sealed != true {
		t.Fatalf("bad")
	}
}

func TestConsul_checkID(t *testing.T) {
	c := testConsulBackend(t)
	if c.checkID() != "vault-sealed-check" {
		t.Errorf("bad")
	}
}

func TestConsul_serviceID(t *testing.T) {
	passingTests := []struct {
		advertiseAddr string
		serviceName   string
		expected      string
	}{
		{
			advertiseAddr: "http://127.0.0.1:8200",
			serviceName:   "sea-tech-astronomy",
			expected:      "sea-tech-astronomy:127.0.0.1:8200",
		},
		{
			advertiseAddr: "http://127.0.0.1:8200/",
			serviceName:   "sea-tech-astronomy",
			expected:      "sea-tech-astronomy:127.0.0.1:8200",
		},
		{
			advertiseAddr: "https://127.0.0.1:8200/",
			serviceName:   "sea-tech-astronomy",
			expected:      "sea-tech-astronomy:127.0.0.1:8200",
		},
	}

	for _, test := range passingTests {
		c := testConsulBackendConfig(t, &consulConf{
			"service": test.serviceName,
		})

		if err := c.UpdateAdvertiseAddr(test.advertiseAddr); err != nil {
			t.Fatalf("bad: %v", err)
		}

		serviceID, err := c.serviceID()
		if err != nil {
			t.Fatalf("bad: %v", err)
		}

		if serviceID != test.expected {
			t.Fatalf("bad: %v != %v", serviceID, test.expected)
		}
	}
}

func TestConsulBackend(t *testing.T) {
	addr := os.Getenv("CONSUL_ADDR")
	if addr == "" {
		t.SkipNow()
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

	b, err := NewBackend("consul", map[string]string{
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
	addr := os.Getenv("CONSUL_ADDR")
	if addr == "" {
		t.SkipNow()
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

	b, err := NewBackend("consul", map[string]string{
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
