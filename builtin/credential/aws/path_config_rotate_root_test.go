// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsauth

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestPathConfigRotateRoot(t *testing.T) {
	getIAMClient = func(sess *session.Session) iamiface.IAMAPI {
		return &awsutil.MockIAM{
			CreateAccessKeyOutput: &iam.CreateAccessKeyOutput{
				AccessKey: &iam.AccessKey{
					AccessKeyId:     aws.String("fizz2"),
					SecretAccessKey: aws.String("buzz2"),
				},
			},
			GetUserOutput: &iam.GetUserOutput{
				User: &iam.User{
					UserName: aws.String("ellen"),
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
