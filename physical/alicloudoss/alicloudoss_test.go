// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloudoss

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestAliCloudOSSBackend(t *testing.T) {
	// ex. http://oss-us-east-1.aliyuncs.com
	endpoint := os.Getenv(AlibabaCloudOSSEndpointEnv)
	accessKeyID := os.Getenv(AlibabaCloudAccessKeyEnv)
	accessKeySecret := os.Getenv(AlibabaCloudSecretKeyEnv)

	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" {
		t.SkipNow()
	}

	conn, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		t.Fatalf("unable to create test client: %s", err)
	}

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-alibaba-testacc-%d", randInt)

	err = conn.CreateBucket(bucket)
	if err != nil {
		t.Fatalf("unable to create test bucket: %s", err)
	}

	defer func() {
		// Gotta list all the objects and delete them
		// before being able to delete the bucket
		b, err := conn.Bucket(bucket)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		listResp, err := b.ListObjects()
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		objects := []string{}
		for _, object := range listResp.Objects {
			objects = append(objects, object.Key)
		}

		_, err = b.DeleteObjects(objects)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		err = conn.DeleteBucket(bucket)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	logger := logging.NewVaultLogger(log.Debug)

	// This uses the same logic to find the Alibaba credentials as we did at the beginning of the test
	b, err := NewAliCloudOSSBackend(
		map[string]string{"bucket": bucket},
		logger,
	)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}
