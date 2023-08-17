// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package minio

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/vault/sdk/helper/docker"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

const (
	accessKeyID = "min-access-key"
	secretKey   = "min-secret-key"
)

func PrepareTestContainer(t *testing.T, version string) (func(), *Config) {
	if version == "" {
		version = "latest"
	}
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ContainerName: "minio",
		ImageRepo:     "docker.mirror.hashicorp.services/minio/minio",
		ImageTag:      version,
		Env: []string{
			"MINIO_ACCESS_KEY=" + accessKeyID,
			"MINIO_SECRET_KEY=" + secretKey,
		},
		Cmd:   []string{"server", "/data"},
		Ports: []string{"9000/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start docker Minio: %s", err)
	}

	svc, err := runner.StartService(context.Background(), connectMinio)
	if err != nil {
		t.Fatalf("Could not start docker Minio: %s", err)
	}

	return svc.Cleanup, &Config{
		Endpoint:        svc.Config.URL().Host,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretKey,
		Region:          "us-east-1",
	}
}

func connectMinio(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme: "s3",
		Host:   fmt.Sprintf("%s:%d", host, port),
	}

	c := &Config{
		Endpoint:        u.Host,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretKey,
		Region:          "us-east-1",
	}
	s3conn, err := c.Conn()
	if err != nil {
		return nil, err
	}

	_, err = s3conn.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	return docker.NewServiceURL(u), nil
}

func (c *Config) Conn() (*s3.S3, error) {
	cfg := &aws.Config{
		DisableSSL:       aws.Bool(true),
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String(c.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.StaticProvider{
					Value: credentials.Value{
						AccessKeyID:     accessKeyID,
						SecretAccessKey: secretKey,
					},
				},
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
				defaults.RemoteCredProvider(*(defaults.Config()), defaults.Handlers()),
			}),
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}
