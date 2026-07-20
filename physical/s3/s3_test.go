// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package s3

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestDefaultS3Backend(t *testing.T) {
	DoS3BackendTest(t, "")
}

func TestS3BackendSseKms(t *testing.T) {
	DoS3BackendTest(t, "alias/aws/s3")
}

func DoS3BackendTest(t *testing.T, kmsKeyId string) {
	if enabled := os.Getenv("VAULT_ACC"); enabled == "" {
		t.Skip()
	}

	if !hasAWSCredentials() {
		t.Skip("Skipping because AWS credentials could not be resolved. See https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-gosdk.html#specifying-credentials for information on how to set up AWS credentials.")
	}

	logger := logging.NewVaultLogger(log.Debug)

	// If the variable is empty or doesn't exist, the default AWS endpoints will be used.
	// Must be a full URL including scheme (e.g. http://127.0.0.1:9000); a bare host:port
	// is not valid because o.BaseEndpoint expects a full URL.
	endpoint := os.Getenv("AWS_S3_ENDPOINT")

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if region == "" {
		region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(t.Context(), config.WithRegion(region))
	if err != nil {
		t.Fatal(err)
	}

	s3conn := s3v2.NewFromConfig(cfg, func(o *s3v2.Options) {
		if endpoint != "" {
			o.BaseEndpoint = awsv2.String(endpoint)
		}
	})

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-s3-testacc-%d", randInt)

	_, err = s3conn.CreateBucket(t.Context(), &s3v2.CreateBucketInput{
		Bucket: awsv2.String(bucket),
		CreateBucketConfiguration: func() *types.CreateBucketConfiguration {
			// AWS S3 requires a LocationConstraint for all regions except us-east-1.
			// Skip for custom endpoints (e.g. MinIO) which don't enforce this.
			if endpoint == "" && region != "us-east-1" {
				return &types.CreateBucketConfiguration{
					LocationConstraint: types.BucketLocationConstraint(region),
				}
			}
			return nil
		}(),
	})
	if err != nil {
		t.Fatalf("unable to create test bucket: %s", err)
	}

	defer func() {
		// List all objects and delete them before deleting the bucket.
		listResp, err := s3conn.ListObjectsV2(context.Background(), &s3v2.ListObjectsV2Input{
			Bucket: awsv2.String(bucket),
		})
		if err == nil && len(listResp.Contents) > 0 {
			objects := &types.Delete{}
			for _, key := range listResp.Contents {
				objects.Objects = append(objects.Objects, types.ObjectIdentifier{Key: key.Key})
			}

			if _, err := s3conn.DeleteObjects(context.Background(), &s3v2.DeleteObjectsInput{
				Bucket: awsv2.String(bucket),
				Delete: objects,
			}); err != nil {
				t.Logf("cleanup: failed to delete objects from bucket %s: %s", bucket, err)
			}
		}

		_, err = s3conn.DeleteBucket(context.Background(), &s3v2.DeleteBucketInput{Bucket: awsv2.String(bucket)})
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	// This uses the same logic to find the AWS credentials as we did at the beginning of the test
	b, err := NewS3Backend(map[string]string{
		"bucket":     bucket,
		"kms_key_id": kmsKeyId,
		"path":       "test/vault",
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

// TestNewS3Backend_InvalidConfig verifies that NewS3Backend returns a
// descriptive error for malformed config values that are validated before
// any network call is made. These cases require no AWS credentials or
// external infrastructure.
//
// Note: the "max_parallel" error path is not covered here because that
// validation occurs after the bucket-accessibility check (ListObjectsV2),
// which requires live S3 or MinIO connectivity.
func TestNewS3Backend_InvalidConfig(t *testing.T) {
	// No t.Parallel: t.Setenv cannot be used after t.Parallel, nor in subtests
	// of a parallel parent. The tests are fast (no network) so sequential is fine.

	// Ensure AWS_S3_BUCKET does not shadow the missing-bucket conf key. NewS3Backend
	// prefers the env var over conf["bucket"], so a set env var would bypass that check.
	t.Setenv("AWS_S3_BUCKET", "")

	logger := logging.NewVaultLogger(log.Error)

	cases := []struct {
		name          string
		conf          map[string]string
		wantErrSubstr string
	}{
		{
			name:          "missing bucket",
			conf:          map[string]string{},
			wantErrSubstr: "'bucket' must be set",
		},
		{
			name: "invalid s3_force_path_style",
			conf: map[string]string{
				"bucket":              "any-bucket",
				"s3_force_path_style": "not-a-bool",
			},
			wantErrSubstr: "invalid boolean set for s3_force_path_style",
		},
		{
			name: "invalid disable_ssl",
			conf: map[string]string{
				"bucket":      "any-bucket",
				"disable_ssl": "not-a-bool",
			},
			wantErrSubstr: "invalid boolean set for disable_ssl",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewS3Backend(tc.conf, logger)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantErrSubstr)
			}
			if !strings.Contains(err.Error(), tc.wantErrSubstr) {
				t.Fatalf("expected error containing %q, got: %v", tc.wantErrSubstr, err)
			}
		})
	}
}

// TestResolveS3Region verifies the documented priority order for region resolution
// (configuration/storage/s3.mdx): AWS_REGION → AWS_DEFAULT_REGION → conf key → us-east-1.
// Tests are sequential (no t.Parallel) because t.Setenv modifies the global environment.
func TestResolveS3Region(t *testing.T) {
	cases := []struct {
		name       string
		envRegion  string
		envDefault string
		confRegion string
		want       string
	}{
		{
			name:      "AWS_REGION takes highest priority",
			envRegion: "eu-west-1", envDefault: "ap-east-1", confRegion: "sa-east-1",
			want: "eu-west-1",
		},
		{
			name:       "AWS_DEFAULT_REGION used when AWS_REGION absent",
			envDefault: "ap-east-1", confRegion: "sa-east-1",
			want: "ap-east-1",
		},
		{
			name:       "conf key used when env vars absent",
			confRegion: "sa-east-1",
			want:       "sa-east-1",
		},
		{
			name: "defaults to us-east-1 when all sources absent",
			want: "us-east-1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("AWS_REGION", tc.envRegion)
			t.Setenv("AWS_DEFAULT_REGION", tc.envDefault)

			conf := map[string]string{}
			if tc.confRegion != "" {
				conf["region"] = tc.confRegion
			}

			got := resolveS3Region(conf)
			if got != tc.want {
				t.Errorf("resolveS3Region() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestBuildS3CredentialChain_StaticCredsWithDefaultAWSProfile verifies that
// buildS3CredentialChain succeeds when explicit static keys are provided even
// when ~/.aws/config contains a [default] profile. This tests the
// WithSharedCredentials(false) fix: awsutil defaults to withSharedCredentials=true
// which injects an empty WithSharedCredentialsFiles("") override that causes
// LoadDefaultConfig to fail with "failed to get shared config profile, default".
func TestBuildS3CredentialChain_StaticCredsWithDefaultAWSProfile(t *testing.T) {
	configFile := t.TempDir() + "/config"
	if err := os.WriteFile(configFile, []byte("[default]\nregion = us-west-2\n"), 0o600); err != nil {
		t.Fatalf("failed to write temp AWS config file: %v", err)
	}
	t.Setenv("AWS_CONFIG_FILE", configFile)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", t.TempDir()+"/nonexistent")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_DEFAULT_PROFILE", "")

	logger := logging.NewVaultLogger(log.Error)
	cfg, err := buildS3CredentialChain("statickey", "staticsecret", "", "us-east-1", logger, nil)
	if err != nil {
		t.Fatalf("buildS3CredentialChain with static creds and [default] profile should not fail: %v", err)
	}
	creds, err := cfg.Credentials.Retrieve(t.Context())
	if err != nil {
		t.Fatalf("Credentials.Retrieve() should not fail: %v", err)
	}
	if creds.AccessKeyID != "statickey" {
		t.Errorf("expected AccessKeyID %q, got %q", "statickey", creds.AccessKeyID)
	}
}

// TestBuildS3CredentialChain_StaticCredsTakePrecedenceOverCredentialsFile verifies
// that when both explicit static keys and AWS_SHARED_CREDENTIALS_FILE are present,
// the static keys win. This is the exact scenario from the s3.mdx credential chain:
// "static config" takes priority over "AWS credential files".
func TestBuildS3CredentialChain_StaticCredsTakePrecedenceOverCredentialsFile(t *testing.T) {
	credsFile := t.TempDir() + "/credentials"
	content := "[default]\naws_access_key_id = fileaccesskey\naws_secret_access_key = filesecretkey\n"
	if err := os.WriteFile(credsFile, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp credentials file: %v", err)
	}
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", credsFile)
	t.Setenv("AWS_CONFIG_FILE", t.TempDir()+"/nonexistent")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_DEFAULT_PROFILE", "")

	logger := logging.NewVaultLogger(log.Error)
	cfg, err := buildS3CredentialChain("statickey", "staticsecret", "", "us-east-1", logger, nil)
	if err != nil {
		t.Fatalf("buildS3CredentialChain with static creds should not fail: %v", err)
	}
	creds, err := cfg.Credentials.Retrieve(t.Context())
	if err != nil {
		t.Fatalf("Credentials.Retrieve() should not fail: %v", err)
	}
	if creds.AccessKeyID != "statickey" {
		t.Errorf("expected static AccessKeyID %q, got %q — file credentials must not override static", "statickey", creds.AccessKeyID)
	}
}

// TestBuildS3CredentialChain_PartialCredentials_Fail verifies that providing
// only one of access_key / secret_key is rejected with awsutil.ErrBadStaticCreds.
// This is validated before any network call.
func TestBuildS3CredentialChain_PartialCredentials_Fail(t *testing.T) {
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")

	logger := logging.NewVaultLogger(log.Error)

	// access_key without secret_key
	_, err := buildS3CredentialChain("accesskey", "", "", "us-east-1", logger, nil)
	if err == nil {
		t.Fatal("expected error when only access_key is set")
	}
	if !strings.Contains(err.Error(), "static AWS client credentials") {
		t.Errorf("unexpected error (want ErrBadStaticCreds): %v", err)
	}
}

// TestBuildS3CredentialChain_NoCredentials_RetrieveFails verifies that when all
// credential sources are removed, Retrieve() returns an error rather than silently
// returning empty credentials.
func TestBuildS3CredentialChain_NoCredentials_RetrieveFails(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", t.TempDir()+"/nonexistent")
	t.Setenv("AWS_CONFIG_FILE", t.TempDir()+"/nonexistent")
	t.Setenv("AWS_ROLE_ARN", "")
	t.Setenv("AWS_WEB_IDENTITY_TOKEN_FILE", "")
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_DEFAULT_PROFILE", "")
	t.Setenv("AWS_CONTAINER_CREDENTIALS_FULL_URI", "")
	t.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "")

	logger := logging.NewVaultLogger(log.Error)
	cfg, err := buildS3CredentialChain("", "", "", "us-east-1", logger, nil)
	if err != nil {
		t.Fatalf("buildS3CredentialChain with no static keys should not fail at build time: %v", err)
	}
	_, err = cfg.Credentials.Retrieve(t.Context())
	if err == nil {
		t.Fatal("Credentials.Retrieve() should fail when no credential sources are available")
	}
}

func hasAWSCredentials() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return false
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return false
	}

	return creds.HasKeys()
}
