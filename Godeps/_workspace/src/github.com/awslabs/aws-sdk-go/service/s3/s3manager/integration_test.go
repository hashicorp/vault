// +build integration

package s3manager_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/internal/test/integration"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"github.com/awslabs/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
)

var integBuf12MB = make([]byte, 1024*1024*12)
var integMD512MB = fmt.Sprintf("%x", md5.Sum(integBuf12MB))
var bucketName *string
var _ = integration.Imported

func TestMain(m *testing.M) {
	setup()
	defer teardown() // only called if we panic
	result := m.Run()
	teardown()
	os.Exit(result)
}

func setup() {
	// Create a bucket for testing
	svc := s3.New(nil)
	bucketName = aws.String(
		fmt.Sprintf("aws-sdk-go-integration-%d", time.Now().Unix()))

	for i := 0; i < 10; i++ {
		_, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: bucketName})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// Delete the bucket
func teardown() {
	svc := s3.New(nil)

	objs, _ := svc.ListObjects(&s3.ListObjectsInput{Bucket: bucketName})
	for _, o := range objs.Contents {
		svc.DeleteObject(&s3.DeleteObjectInput{Bucket: bucketName, Key: o.Key})
	}

	uploads, _ := svc.ListMultipartUploads(&s3.ListMultipartUploadsInput{Bucket: bucketName})
	for _, u := range uploads.Uploads {
		svc.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
			Bucket:   bucketName,
			Key:      u.Key,
			UploadID: u.UploadID,
		})
	}

	svc.DeleteBucket(&s3.DeleteBucketInput{Bucket: bucketName})
}

func validate(t *testing.T, key string, md5value string) {
	svc := s3.New(nil)
	resp, err := svc.GetObject(&s3.GetObjectInput{Bucket: bucketName, Key: &key})
	assert.NoError(t, err)
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, md5value, fmt.Sprintf("%x", md5.Sum(b)))
}

func TestUploadConcurrently(t *testing.T) {
	svc := s3.New(nil)
	key := "12mb-1"
	out, err := s3manager.Upload(svc, &s3manager.UploadInput{
		Bucket: bucketName,
		Key:    &key,
		Body:   bytes.NewReader(integBuf12MB),
	}, nil)

	assert.NoError(t, err)
	assert.NotEqual(t, "", out.UploadID)
	assert.Regexp(t, `^https?://.+/`+key+`$`, out.Location)

	validate(t, key, integMD512MB)
}

func TestUploadFailCleanup(t *testing.T) {
	svc := s3.New(nil)

	// Break checksum on 2nd part so it fails
	part := 0
	svc.Handlers.Build.PushBack(func(r *aws.Request) {
		if r.Operation.Name == "UploadPart" {
			if part == 1 {
				r.HTTPRequest.Header.Set("X-Amz-Content-Sha256", "000")
			}
			part++
		}
	})

	key := "12mb-leave"
	u, err := s3manager.Upload(svc, &s3manager.UploadInput{
		Bucket: bucketName,
		Key:    &key,
		Body:   bytes.NewReader(integBuf12MB),
	}, &s3manager.UploadOptions{
		LeavePartsOnError: false,
	})
	assert.Error(t, err)

	_, err = svc.ListParts(&s3.ListPartsInput{
		Bucket: bucketName, Key: &key, UploadID: &u.UploadID})
	assert.Error(t, err)
}
