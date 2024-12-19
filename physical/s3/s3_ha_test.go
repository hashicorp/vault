// s3_test.go

package s3

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestHABackend(t *testing.T) {
	if enabled := os.Getenv("VAULT_ACC"); enabled == "" {
		t.Skip()
	}

	if !hasAWSCredentials(t) {
		t.Skip("Skipping because AWS credentials could not be resolved. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials for information on how to set up AWS credentials.")
	}

	ctx := context.TODO()

	// If the variable is empty or doesn't exist, the default
	// AWS endpoints will be used
	endpoint := os.Getenv("AWS_S3_ENDPOINT")
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if endpoint != "" {
					return aws.Endpoint{URL: endpoint}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			})),
	)
	if err != nil {
		t.Fatal(err)
	}
	s3conn := s3.NewFromConfig(cfg)

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-s3-testacc-%d", randInt)

	_, err = s3conn.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Fatalf("unable to create test bucket: %s", err)
	}

	defer func() {
		listResp, _ := s3conn.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket: aws.String(bucket),
		})

		var objIds []types.ObjectIdentifier
		for _, obj := range listResp.Contents {
			objIds = append(objIds, types.ObjectIdentifier{Key: obj.Key})
		}

		if len(objIds) > 0 {
			_, _ = s3conn.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: &types.Delete{Objects: objIds},
			})
		}

		_, err := s3conn.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	// Configure the backend
	config := map[string]string{
		"bucket":     bucket,
		"region":     region,
		"ha_enabled": "true",
	}

	logger := logging.NewVaultLogger(log.Debug)
	// Create first backend instance
	b1, err := NewS3Backend(config, logger)
	if err != nil {
		t.Fatal(err)
	}

	// Create second backend instance
	b2, err := NewS3Backend(config, logger)
	if err != nil {
		t.Fatal(err)
	}

	// Run the HA tests
	physical.ExerciseBackend(t, b1)
	physical.ExerciseBackend_ListPrefix(t, b1)
	physical.ExerciseHABackend(t, b1.(physical.HABackend), b2.(physical.HABackend))
}
