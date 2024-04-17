// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package azurite

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/hashicorp/vault/sdk/helper/docker"
)

type Config struct {
	Endpoint    string
	AccountName string
	AccountKey  string
}

func (c Config) Address() string {
	return c.Endpoint
}

func (c Config) URL() *url.URL {
	return &url.URL{Scheme: "http", Host: c.Endpoint, Path: "/" + accountName}
}

//func (c Config) ConnectionString() string {
//	elems := []string{
//		"DefaultEndpointsProtocol=http",
//		"AccountName=" + accountName,
//		"AccountKey=" + accountKey,
//		"EndpointSuffix=" + c.Endpoint,
//	}
//	return strings.Join(elems, ";")
//}

func (c Config) ContainerURL(ctx context.Context, container string) (*azblob.ContainerURL, error) {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{
		//Log:        pipeline.LogOptions{
		//	Log: func(level pipeline.LogLevel, message string) {
		//		log.Println(message)
		//	},
		//	ShouldLog: func(level pipeline.LogLevel) bool {
		//		return true
		//	},
		//},
	})
	u := *c.URL()
	u.Path += "/" + container
	cu := azblob.NewContainerURL(u, p)
	return &cu, nil
}

var _ docker.ServiceConfig = &Config{}

const (
	accountName = "testaccount"
	accountKey  = "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="
)

func PrepareTestContainer(t *testing.T, version string) (func(), docker.ServiceConfig) {
	if version == "" {
		version = "latest"
	}
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ContainerName: "azurite",
		ImageRepo:     "mcr.microsoft.com/azure-storage/azurite",
		ImageTag:      version,
		Cmd:           []string{"azurite-blob", "--blobHost", "0.0.0.0", "--blobPort", "10000", "-d", "/dev/stderr"},
		Ports:         []string{"10000/tcp"},
		Env:           []string{fmt.Sprintf(`AZURITE_ACCOUNTS=%s:%s`, accountName, accountKey)},
	})
	if err != nil {
		t.Fatalf("Could not start docker Azurite: %s", err)
	}

	svc, err := runner.StartService(context.Background(), connectAzure)
	if err != nil {
		t.Fatalf("Could not start docker Azurite: %s", err)
	}

	return svc.Cleanup, svc.Config
}

func connectAzure(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	cfg := &Config{
		Endpoint:    fmt.Sprintf("%s:%d", host, port),
		AccountName: accountName,
		AccountKey:  accountKey,
	}

	containerURL, err := cfg.ContainerURL(ctx, "testcontainer")
	if err != nil {
		return nil, err
	}
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessContainer)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
