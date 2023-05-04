package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

func TestRotation(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b := Backend()
	err := b.initQueue(context.Background(), nil)
	if err != nil {
		t.Fatalf("couldn't initialize queue: %s", err)
	}

	miam, err := awsutil.NewMockIAM(
		//awsutil.WithGetUserOutput(),
		awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
			AccessKeyMetadata: []*iam.AccessKeyMetadata{
				{},
			},
		}),
		awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
			AccessKey: &iam.AccessKey{
				AccessKeyId:     aws.String("abcdefghijklmnopqrstuvwxyz"),
				SecretAccessKey: aws.String("bigsecret"),
			},
		}),
	)(nil)
	if err != nil {
		t.Fatalf("couldn't initialze mock IAM handler: %s", err)
	}
	b.iamClient = miam

	staticConfig := staticRoleConfig{
		Name:           "test",
		Username:       "jane-doe",
		RotationPeriod: 2 * time.Second,
	}
	err = b.createCredential(context.Background(), config.StorageView, staticConfig)
	if err != nil {
		t.Fatalf("couldn't initialize queue: %s", err)
	}

	item := &queue.Item{
		Key:      "test",
		Value:    staticConfig,
		Priority: time.Now().Add(staticConfig.RotationPeriod).Unix(),
	}
	err = b.credRotationQueue.Push(item)
	if err != nil {
		t.Fatalf("couldn't push item onto queue: %s", err)
	}

	time.Sleep(5 * time.Second)
	// update aws responses
	miam, err = awsutil.NewMockIAM(
		// old key
		awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
			AccessKeyMetadata: []*iam.AccessKeyMetadata{
				{
					AccessKeyId: aws.String("abcdefghijklmnopqrstuvwxyz"),
				},
			},
		}),
		// new key
		awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
			AccessKey: &iam.AccessKey{
				AccessKeyId:     aws.String("zyxwvutsrqponmlkjihgfedcba"),
				SecretAccessKey: aws.String("biggersecret"),
			},
		}),
	)(nil)
	if err != nil {
		t.Fatalf("couldn't initialze mock IAM handler: %s", err)
	}
	b.iamClient = miam

	req := &logical.Request{
		Storage: config.StorageView,
	}
	err = b.rotateExpiredStaticCreds(context.Background(), req)
	if err != nil {
		t.Fatalf("got an error rotating credentials: %s", err)
	}

	entry, err := config.StorageView.Get(context.Background(), formatCredsStoragePath("test"))
	if err != nil {
		t.Fatalf("got an error retrieving credentials")
	}
	var out awsCredentials
	entry.DecodeJSON(&out)
	if out.SecretAccessKey != "biggersecret" {
		t.Fatal("mismatched secret")
	}

}
