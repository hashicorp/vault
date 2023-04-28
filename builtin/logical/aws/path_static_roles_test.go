package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestStaticRolesValidation(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	cases := []struct {
		name        string
		opts        []awsutil.MockIAMOption
		requestData map[string]interface{}
		isError     bool
	}{
		{
			name: "all good",
			opts: []awsutil.MockIAMOption{
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("jane-doe")}}),
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String("abcdefghijklmnopqrstuvwxyz"),
						SecretAccessKey: aws.String("zyxwvutsrqponmlkjihgfedcba"),
						UserName:        aws.String("jane-doe"),
					},
				}),
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{},
					IsTruncated:       aws.Bool(false),
				}),
			},
			requestData: map[string]interface{}{
				"name":            "test",
				"username":        "jane-doe",
				"rotation_period": "1d",
			},
		},
		{
			name: "bad user",
			opts: []awsutil.MockIAMOption{
				awsutil.WithGetUserError(errors.New("oh no")),
			},
			requestData: map[string]interface{}{
				"name":            "test",
				"username":        "jane-doe",
				"rotation_period": "24h",
			},
			isError: true,
		},
		{
			name: "user mismatch",
			opts: []awsutil.MockIAMOption{
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("ms-impostor")}}),
			},
			requestData: map[string]interface{}{
				"name":            "test",
				"username":        "jane-doe",
				"rotation_period": "1d2h",
			},
			isError: true,
		},
		{
			name: "bad rotation period",
			opts: []awsutil.MockIAMOption{
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("jane-doe")}}),
			},
			requestData: map[string]interface{}{
				"name":            "test",
				"username":        "jane-doe",
				"rotation_period": "45s",
			},
			isError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := Backend()
			miam, err := awsutil.NewMockIAM(c.opts...)(nil)
			if err != nil {
				t.Fatal(err)
			}
			b.iamClient = miam
			if err := b.Setup(context.Background(), config); err != nil {
				t.Fatal(err)
			}
			req := &logical.Request{
				Operation: logical.UpdateOperation,
				Storage:   config.StorageView,
				Data:      c.requestData,
				Path:      "static-roles/test",
			}
			_, err = b.pathStaticRolesWrite(context.Background(), req, staticRoleFieldData(req.Data))
			if c.isError && err == nil {
				t.Fatal("expected an error but didn't get one")
			} else if !c.isError && err != nil {
				t.Fatalf("got an unexpected error: %s", err)
			}
		})
	}
}

func staticRoleFieldData(data map[string]interface{}) *framework.FieldData {
	schema := map[string]*framework.FieldSchema{
		paramRoleName: {
			Type:        framework.TypeString,
			Description: descRoleName,
		},
		paramUsername: {
			Type:        framework.TypeString,
			Description: descUsername,
		},
		paramRotationPeriod: {
			Type:        framework.TypeDurationSecond,
			Description: descRotationPeriod,
		},
	}

	return &framework.FieldData{
		Raw:    data,
		Schema: schema,
	}
}
