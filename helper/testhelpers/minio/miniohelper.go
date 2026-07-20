// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package minio

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	svc, err := runner.StartService(t.Context(), connectMinio)
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

	_, err = s3conn.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	return docker.NewServiceURL(u), nil
}

func (c *Config) Conn() (*s3.Client, error) {
	// Static credentials only: MinIO always uses the hardcoded accessKeyID/secretKey.
	// A full chain (env, shared file, IMDS) is omitted to prevent tests accidentally
	// resolving real AWS credentials and hitting live S3 instead of the local container.
	// Empty shared config/credentials file lists prevent LoadDefaultConfig from reading
	// any host AWS profile, making the helper hermetic even when AWS_PROFILE or
	// AWS_DEFAULT_PROFILE is set in the caller's environment.
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(c.Region),
		config.WithSharedConfigFiles([]string{}),
		config.WithSharedCredentialsFiles([]string{}),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.AccessKeyID,
			c.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://" + c.Endpoint)
		o.UsePathStyle = true
	})
	return client, nil
}
