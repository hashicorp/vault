package redis

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/logical"
)

func getBackend(t *testing.T) (*backend, context.Context, logical.Storage, string, func()) {
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "redis",
		ImageTag:      "7",
		ContainerName: "reddis",
		Ports:         []string{"6379/tcp"},
	})
	if err != nil {
		t.Fatalf("could not start docker redis: %s", err)
	}

	ctx := context.Background()
	svc, err := runner.StartService(ctx, func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		client := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", host, port),
		})
		if result := client.Ping(ctx); result.Err() != nil {
			return nil, result.Err()
		}

		return docker.NewServiceHostPort(host, port), nil
	})
	if err != nil {
		t.Fatalf("could not start docker redis: %s", err)
	}

	b := Backend()
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	if err := b.Setup(ctx, config); err != nil {
		t.Fatal(err)
	}

	return b, ctx, config.StorageView, svc.Config.Address(), svc.Cleanup
}

func setConfig(t *testing.T, b *backend, ctx context.Context, s logical.Storage, addr string, rotate bool) {
	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   s,
		Data: map[string]interface{}{
			"address":  addr,
			"username": "default",
			"rotate":   rotate,
		},
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
}

func getBackendAndSetConfig(t *testing.T) (*backend, context.Context, logical.Storage, string, func()) {
	b, ctx, s, addr, stop := getBackend(t)
	setConfig(t, b, ctx, s, addr, false)
	return b, ctx, s, addr, stop
}

func TestBackend_NoConfiguration(t *testing.T) {
	b, ctx, s, _, stop := getBackend(t)
	defer stop()

	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   s,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("Expected error, got: %#v", resp)
	}
	if resp.Error().Error() != "No configuration found" {
		t.Fatalf("Wrong error: %s", resp.Error())
	}
}

func TestBackend_Configuration(t *testing.T) {
	b, ctx, s, addr, stop := getBackend(t)
	defer stop()

	conf := map[string]interface{}{
		"address":  addr,
		"username": "default",
	}

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   s,
		Data:      conf,
	}

	resp, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(conf, resp.Data) {
		t.Fatalf("Expected: %#v\nActual: %#v", conf, resp.Data)
	}

	req.Operation = logical.ReadOperation
	req.Data = nil
	resp, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(conf, resp.Data) {
		t.Fatalf("Expected: %#v\nActual: %#v", conf, resp.Data)
	}
}

func TestBackend_Rotate(t *testing.T) {
	b, ctx, s, addr, stop := getBackendAndSetConfig(t)
	defer stop()

	// We can connect without a password
	if err := tryToConnect(ctx, addr); err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Path:      "config/rotate",
		Operation: logical.UpdateOperation,
		Storage:   s,
	}

	resp, err := b.HandleRequest(ctx, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}

	// Trying to connect without a password should not work anymore
	if err := tryToConnect(ctx, addr); err == nil {
		t.Fatal("an error should have happened")
	}

	// We can rotate the password again because the new one has been saved
	// and the client has been refreshed
	resp, err = b.HandleRequest(ctx, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
}

func TestBackend_RotateNoConfig(t *testing.T) {
	b, ctx, s, _, stop := getBackend(t)
	defer stop()

	resp, err := b.HandleRequest(ctx, &logical.Request{
		Path:      "config/rotate",
		Operation: logical.UpdateOperation,
		Storage:   s,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("Expected error, got: %#v", resp)
	}
	if resp.Error().Error() != "the configuration has not been set" {
		t.Fatalf("Wrong error: %s", resp.Error())
	}
}

func TestBackend_RotateDuringConfiguration(t *testing.T) {
	b, ctx, s, addr, stop := getBackend(t)
	defer stop()

	if err := tryToConnect(ctx, addr); err != nil {
		t.Fatal(err)
	}

	setConfig(t, b, ctx, s, addr, true)

	// We cannot connect anymore because the password has been rotated
	if err := tryToConnect(ctx, addr); err == nil {
		t.Fatal("an error should have happened")
	}

	// Calling rotate works because the client has been refreshed
	resp, err := b.HandleRequest(ctx, &logical.Request{
		Path:      "config/rotate",
		Operation: logical.UpdateOperation,
		Storage:   s,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
}

func tryToConnect(ctx context.Context, addr string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: "default",
	})

	_, err := client.Ping(ctx).Result()
	return err
}
