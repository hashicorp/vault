package tencentcloudcos

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func TestTencentCloudCOSBackend(t *testing.T) {
	accessKey := os.Getenv(PROVIDER_SECRET_ID)
	secretKey := os.Getenv(PROVIDER_SECRET_KEY)
	sessionToken := os.Getenv(PROVIDER_SECURITY_TOKEN)
	region := os.Getenv(PROVIDER_REGION)
	appId := os.Getenv("TENCENTCLOUD_COS_APPID")

	if accessKey == "" || secretKey == "" || appId == "" {
		t.SkipNow()
	}

	if region == "" {
		region = "ap-guangzhou"
	}

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-tencentcloud-testacc-%d-%s", randInt, appId)

	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	if err != nil {
		t.Fatalf("unable to create test client: %s", err)
	}

	client := cos.NewClient(
		&cos.BaseURL{BucketURL: u},
		&http.Client{
			Timeout: 60 * time.Second,
			Transport: &cos.AuthorizationTransport{
				SecretID:     accessKey,
				SecretKey:    secretKey,
				SessionToken: sessionToken,
			},
		},
	)

	_, err = client.Bucket.Put(context.Background(), nil)
	if err != nil {
		t.Fatalf("unable to create test bucket: %s", err)
	}

	defer func() {
		// Gotta list all the objects and delete them
		// before being able to delete the bucket
		fs, rsp, err := client.Bucket.Get(context.Background(), &cos.BucketGetOptions{})
		if rsp == nil {
			t.Fatalf("failed to list bucket %v: %v", bucket, fmt.Errorf("no response"))
		}
		defer rsp.Body.Close()

		if rsp.StatusCode == 404 {
			return
		}

		if err != nil {
			t.Fatalf("failed to list bucket %v: %v", bucket, err)
		}

		objects := []cos.Object{}
		for _, v := range fs.Contents {
			objects = append(objects, cos.Object{Key: v.Key})
		}

		if len(objects) > 0 {
			_, _, err = client.Object.DeleteMulti(context.Background(), &cos.ObjectDeleteMultiOptions{Objects: objects})
			if err != nil {
				t.Fatalf("failed to empty bucket %v: %v", bucket, err)
			}
		}

		_, err = client.Bucket.Delete(context.Background())
		if err != nil {
			t.Fatalf("failed to delete bucket %v: %v", bucket, err)
		}
	}()

	logger := logging.NewVaultLogger(log.Debug)

	// This uses the same logic to find the TencentCloud credentials as we did at the beginning of the test
	b, err := NewTencentCloudCOSBackend(
		map[string]string{"bucket": bucket},
		logger,
	)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}
