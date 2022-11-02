package fakegcsserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// In principle we don't need docker for fake-gcs-server, we could run it in
// memory instead.  However I had an error trying to use it:
//  go: finding module for package google.golang.org/grpc/naming
//  github.com/hashicorp/vault/vault imports
//        google.golang.org/grpc/naming: module google.golang.org/grpc@latest found (v1.32.0), but does not contain package google.golang.org/grpc/naming
// so it seemed easiest to go this route.  Vault already has too many deps anyway.

func PrepareTestContainer(t *testing.T, version string) (func(), docker.ServiceConfig) {
	if version == "" {
		version = "latest"
	}
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ContainerName: "fake-gcs-server",
		ImageRepo:     "docker.mirror.hashicorp.services/fsouza/fake-gcs-server",
		ImageTag:      version,
		Cmd:           []string{"-scheme", "http", "-public-host", "storage.gcs.127.0.0.1.nip.io:4443"},
		Ports:         []string{"4443/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start docker fake-gcs-server: %s", err)
	}

	svc, err := runner.StartService(context.Background(), connectGCS)
	if err != nil {
		t.Fatalf("Could not start docker fake-gcs-server: %s", err)
	}

	return svc.Cleanup, svc.Config
}

func connectGCS(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   "storage/v1/b",
	}
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	httpClient := &http.Client{Transport: transCfg}
	client, err := storage.NewClient(context.TODO(), option.WithEndpoint(u.String()), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	it := client.Buckets(ctx, "test")
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return docker.NewServiceURL(u), nil
}
