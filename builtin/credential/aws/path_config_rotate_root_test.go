// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"testing"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	awsutil "github.com/hashicorp/go-secure-stdlib/awsutil/v2"
	"github.com/hashicorp/vault/sdk/logical"
)

type mockIAMClient = awsutil.MockIAM

func TestPathConfigRotateRoot(t *testing.T) {
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_DEFAULT_REGION", "us-east-1")

	getIAMClient = func(cfg *awsv2.Config) awsutil.IAMClient {
		return &mockIAMClient{
			CreateAccessKeyOutput: &iam.CreateAccessKeyOutput{
				AccessKey: &iamtypes.AccessKey{
					AccessKeyId:     awsv2.String("fizz2"),
					SecretAccessKey: awsv2.String("buzz2"),
				},
			},
			GetUserOutput: &iam.GetUserOutput{
				User: &iamtypes.User{
					UserName: awsv2.String("ellen"),
				},
			},
		}
	}

	ctx := context.Background()
	config := logical.TestBackendConfig()
	logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	clientConf := &clientConfig{
		AccessKey: "fizz1",
		SecretKey: "buzz1",
	}
	entry, err := logical.StorageEntryJSON("config/client", clientConf)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(ctx, entry); err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/rotate-root",
		Storage:   storage,
	}
	resp, err := b.HandleRequest(ctx, req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%v", resp, err)
	}
	if resp == nil {
		t.Fatal("expected nil response to represent a 204")
	}
	if resp.Data == nil {
		t.Fatal("expected resp.Data")
	}
	if resp.Data["access_key"].(string) != "fizz2" {
		t.Fatalf("expected new access key buzz2 but received %s", resp.Data["access_key"])
	}
	newClientConf, err := b.nonLockedClientConfigEntry(ctx, req.Storage)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["access_key"].(string) != newClientConf.AccessKey {
		t.Fatalf("expected new access key buzz2 to be saved to storage but receieved %s", clientConf.AccessKey)
	}
}
