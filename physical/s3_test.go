package physical

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestS3Backend(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.SkipNow()
	}

	creds, err := credentials.NewEnvCredentials().Get()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// If the variable is empty or doesn't exist, the default
	// AWS endpoints will be used
	endpoint := os.Getenv("AWS_S3_ENDPOINT")

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	s3conn := s3.New(session.New(&aws.Config{
		Credentials: credentials.NewEnvCredentials(),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}))

	var randInt = rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-s3-testacc-%d", randInt)

	_, err = s3conn.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Fatalf("unable to create test bucket: %s", err)
	}

	defer func() {
		// Gotta list all the objects and delete them
		// before being able to delete the bucket
		listResp, _ := s3conn.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucket),
		})

		objects := &s3.Delete{}
		for _, key := range listResp.Contents {
			oi := &s3.ObjectIdentifier{Key: key.Key}
			objects.Objects = append(objects.Objects, oi)
		}

		s3conn.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: objects,
		})

		_, err := s3conn.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String(bucket)})
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	b, err := NewBackend("s3", logger, map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"bucket":        bucket,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}
