package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/vault/sdk/logical"
)

type mockSTSClient struct {
	stsiface.STSAPI

	AssumeRoleInput  *sts.AssumeRoleInput
	AssumeRoleOutput *sts.AssumeRoleOutput
	AssumeRoleError  error
}

func (m *mockSTSClient) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	m.AssumeRoleInput = input
	if m.AssumeRoleError != nil {
		return nil, m.AssumeRoleError
	}

	return m.AssumeRoleOutput, nil
}

func (m *mockSTSClient) AssumeRoleWithContext(ctx aws.Context, input *sts.AssumeRoleInput, opts ...request.Option) (*sts.AssumeRoleOutput, error) {
	m.AssumeRoleInput = input
	if m.AssumeRoleError != nil {
		return nil, m.AssumeRoleError
	}

	return m.AssumeRoleOutput, nil
}

func TestBackend_PathCredsRead(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &logical.StaticSystemView{
		EntityVal: &logical.Entity{
			ID:   "test-id",
			Name: "test",
			Aliases: []*logical.Alias{
				{
					MountAccessor: "auth_jwt_1234",
					Metadata: map[string]string{
						"group": "dev",
					},
				},
			},
		},
	}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	b.stsClient = &mockSTSClient{
		AssumeRoleOutput: &sts.AssumeRoleOutput{
			AssumedRoleUser: &sts.AssumedRoleUser{
				Arn: aws.String("arn:aws:iam::123456789012:role/test"),
			},
			Credentials: &sts.Credentials{
				Expiration:      aws.Time(time.Now()),
				AccessKeyId:     aws.String("123456"),
				SecretAccessKey: aws.String("123456"),
				SessionToken:    aws.String("123456"),
			},
		},
	}

	roleData := map[string]interface{}{
		"role_arns":                         []string{"arn:aws:iam::123456789012:role/test"},
		"credential_type":                   assumedRoleCred,
		"enable_policy_document_templating": true,
		"policy_document": `
{
	"Version": "2012-10-07",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": "ec2:Describe*",
			"Resource": "*",
            "Condition": {
				"StringEquals": {
					"ec2:ResourceTag/Entity": "{{identity.entity.name}}",
					"ec2:ResourceTag/Group": "{{identity.entity.aliases.auth_jwt_1234.metadata.group}}"
				}
            }
		}
	]
}`,
		"default_sts_ttl": 3600,
		"max_sts_ttl":     3600,
	}

	updateRoleResp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test",
		Storage:   config.StorageView,
		Data:      roleData,
	})
	if err != nil || (updateRoleResp != nil && updateRoleResp.IsError()) {
		t.Fatalf("bad: role creation failed. updateRoleResp:%#v\nerr:%v", updateRoleResp, err)
	}

	readResp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/test",
		Storage:   config.StorageView,
		EntityID:  "test-id",
	})
	if err != nil || (readResp != nil && readResp.IsError()) {
		t.Fatalf("bad: reading creds failed. readResp:%#v\n err:%v", readResp, err)
	}

	expectedPolicy, err := compactJSON(`{
		"Version": "2012-10-07",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": "ec2:Describe*",
				"Resource": "*",
				"Condition": {
					"StringEquals": {
						"ec2:ResourceTag/Entity": "test",
						"ec2:ResourceTag/Group": "dev"
					}
				}
			}
		]
	}`)
	assert.NoError(t, err)

	assert.Equal(t, expectedPolicy, *b.stsClient.(*mockSTSClient).AssumeRoleInput.Policy)
}
