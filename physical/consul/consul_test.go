package consul

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestConsul_newConsulBackend(t *testing.T) {
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
	}

	for _, test := range tests {
		logger := logging.NewVaultLogger(log.Debug)

		be, err := NewConsulBackend(test.consulConfig, logger)
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

		if test.path != c.path {
			t.Errorf("bad: %s %v != %v", test.name, test.path, c.path)
		}

		if test.consistencyMode != c.consistencyMode {
			t.Errorf("bad consistency_mode value: %v != %v", test.consistencyMode, c.consistencyMode)
		}

		// The configuration stored in the Consul "client" object is not exported, so
		// we either have to skip validating it, or add a method to export it, or use reflection.
		consulConfig := reflect.Indirect(reflect.ValueOf(c.client)).FieldByName("config")
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
	}
}

func TestConsulBackend(t *testing.T) {
	consulToken := os.Getenv("CONSUL_HTTP_TOKEN")
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	if addr == "" {
		cleanup, connURL, token := consul.PrepareTestContainer(t, "1.4.4")
		defer cleanup()
		addr, consulToken = connURL, token
	}

	conf := api.DefaultConfig()
	conf.Address = addr
	conf.Token = consulToken
	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewConsulBackend(map[string]string{
		"address":      conf.Address,
		"path":         randPath,
		"max_parallel": "256",
		"token":        conf.Token,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestConsul_TooLarge(t *testing.T) {
	consulToken := os.Getenv("CONSUL_HTTP_TOKEN")
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	if addr == "" {
		cleanup, connURL, token := consul.PrepareTestContainer(t, "1.4.4")
		defer cleanup()
		addr, consulToken = connURL, token
	}

	conf := api.DefaultConfig()
	conf.Address = addr
	conf.Token = consulToken
	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewConsulBackend(map[string]string{
		"address":      conf.Address,
		"path":         randPath,
		"max_parallel": "256",
		"token":        conf.Token,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	zeros := make([]byte, 600000, 600000)
	n, err := rand.Read(zeros)
	if n != 600000 {
		t.Fatalf("expected 500k zeros, read %d", n)
	}
	if err != nil {
		t.Fatal(err)
	}

	err = b.Put(context.Background(), &physical.Entry{
		Key:   "foo",
		Value: zeros,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), physical.ErrValueTooLarge) {
		t.Fatalf("expected value too large error, got %v", err)
	}

	err = b.(physical.Transactional).Transaction(context.Background(), []*physical.TxnEntry{
		{
			Operation: physical.PutOperation,
			Entry: &physical.Entry{
				Key:   "foo",
				Value: zeros,
			},
		},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), physical.ErrValueTooLarge) {
		t.Fatalf("expected value too large error, got %v", err)
	}
}

func TestConsulHABackend(t *testing.T) {
	consulToken := os.Getenv("CONSUL_HTTP_TOKEN")
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	if addr == "" {
		cleanup, connURL, token := consul.PrepareTestContainer(t, "1.4.4")
		defer cleanup()
		addr, consulToken = connURL, token
	}

	conf := api.DefaultConfig()
	conf.Address = addr
	conf.Token = consulToken
	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())
	defer func() {
		client.KV().DeleteTree(randPath, nil)
	}()

	logger := logging.NewVaultLogger(log.Debug)
	config := map[string]string{
		"address":      conf.Address,
		"path":         randPath,
		"max_parallel": "-1",
		"token":        conf.Token,
	}

	b, err := NewConsulBackend(config, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	b2, err := NewConsulBackend(config, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseHABackend(t, b.(physical.HABackend), b2.(physical.HABackend))

	detect, ok := b.(physical.RedirectDetect)
	if !ok {
		t.Fatalf("consul does not implement RedirectDetect")
	}
	host, err := detect.DetectHostAddr()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if host == "" {
		t.Fatalf("bad addr: %v", host)
	}
}
