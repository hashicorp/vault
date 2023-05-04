package aws

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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

func TestStaticRolesWrite(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	cases := []struct {
		name          string
		opts          []awsutil.MockIAMOption
		data          map[string]interface{}
		expectedError bool
		findUser      bool
	}{
		{
			name: "happy path",
			opts: []awsutil.MockIAMOption{
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("jane-doe")}}),
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{},
					IsTruncated:       aws.Bool(false),
				}),
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String("abcdefghijklmnopqrstuvwxyz"),
						SecretAccessKey: aws.String("zyxwvutsrqponmlkjihgfedcba"),
						UserName:        aws.String("jane-doe"),
					},
				}),
			},
			data: map[string]interface{}{
				"name":            "test",
				"username":        "jane-doe",
				"rotation_period": "1d",
			},
			// writes role, writes cred
			findUser: true,
		},
		{
			name: "no aws user",
			opts: []awsutil.MockIAMOption{
				awsutil.WithGetUserError(errors.New("no such user, etc etc")),
			},
			data: map[string]interface{}{
				"name":            "test",
				"username":        "a-nony-mous",
				"rotation_period": "15s",
			},
			expectedError: true,
		},
	}

	// if a user exists (user doesn't exist is tested in validation)
	// we'll check how many keys the user has - if it's two, we delete one.

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			miam, err := awsutil.NewMockIAM(
				c.opts...,
			)(nil)
			if err != nil {
				t.Fatal(err)
			}

			b := Backend()
			b.iamClient = miam
			if err := b.Setup(context.Background(), config); err != nil {
				t.Fatal(err)
			}

			req := &logical.Request{
				Operation: logical.UpdateOperation,
				Storage:   config.StorageView,
				Data:      c.data,
				Path:      "static-roles/test",
			}

			r, err := b.pathStaticRolesWrite(context.Background(), req, staticRoleFieldData(req.Data))
			if c.expectedError && err == nil {
				t.Fatal(err)
			}
			if c.findUser && r == nil {
				t.Fatal("response was nil, but it shouldn't have been")
			}

			role, err := config.StorageView.Get(context.Background(), req.Path)
			if c.findUser && (err != nil || role == nil) {
				t.Fatalf("couldn't find the role we should have stored: %s", err)
			}
		})
	}
}

func TestStaticRoleRead(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	staticRole := staticRoleConfig{
		Name:           "test",
		Username:       "jane-doe",
		RotationPeriod: 24 * time.Hour,
	}
	entry, err := logical.StorageEntryJSON(formatRoleStoragePath(staticRole.Name), staticRole)
	if err != nil {
		t.Fatal(err)
	}
	err = config.StorageView.Put(context.Background(), entry)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"name": "test",
		},
		Path: "static-roles/test",
	}

	b := Backend()

	r, err := b.pathStaticRolesRead(context.Background(), req, staticRoleFieldData(req.Data))
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("response was nil, but it shouldn't have been")
	}
}

func TestStaticRoleDelete(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	staticRole := staticRoleConfig{
		Name:           "test",
		Username:       "jane-doe",
		RotationPeriod: 24 * time.Hour,
	}
	entry, err := logical.StorageEntryJSON(formatRoleStoragePath(staticRole.Name), staticRole)
	if err != nil {
		t.Fatal(err)
	}
	err = config.StorageView.Put(context.Background(), entry)
	if err != nil {
		t.Fatal(err)
	}

	l, err := config.StorageView.List(context.Background(), "")
	if err != nil || len(l) != 1 {
		t.Fatalf("couldn't add an entry to storage during test setup: %s", err)
	}
	fmt.Println(len(l))

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   config.StorageView,
		Data: map[string]interface{}{
			"name": "test",
		},
		Path: "static-roles/test",
	}

	b := Backend()

	r, err := b.pathStaticRolesDelete(context.Background(), req, staticRoleFieldData(req.Data))
	if err != nil {
		t.Fatal(err)
	}
	if r != nil {
		t.Fatal("response wasn't nil, but it should have been")
	}

	l, err = config.StorageView.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 0 {
		t.Fatal("size of role storage is non zero after delete")
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
