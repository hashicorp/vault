package aws

import (
	"context"
	"testing"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type extendedMockIAM struct {
	iamiface.IAMAPI

	ListAccessKeysOutput *iam.ListAccessKeysOutput
	ListAccessKeysError  error
}

// Hang a new mock function off of a temporary test struct until the mock dependency we're using gets updated
func (e *extendedMockIAM) ListAccessKeys(*iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	if e.ListAccessKeysError != nil {
		return nil, e.ListAccessKeysError
	}

	return e.ListAccessKeysOutput, nil
}

func TestStaticRolesValidation(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	mockIAM, err := awsutil.NewMockIAM(
		awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("jane-doe")}}),
		awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
			AccessKey: &iam.AccessKey{
				AccessKeyId:     aws.String("abcdefghijklmnopqrstuvwxyz"),
				SecretAccessKey: aws.String("zyxwvutsrqponmlkjihgfedcba"),
				UserName:        aws.String("jane-doe"),
			},
		}),
	)(nil)

	goodUser := &extendedMockIAM{
		IAMAPI: mockIAM,
		ListAccessKeysOutput: &iam.ListAccessKeysOutput{
			AccessKeyMetadata: []*iam.AccessKeyMetadata{},
			IsTruncated:       aws.Bool(false),
		},
	}

	//_ = &extendedMockIAM{
	//	MockIAM: awsutil.MockIAM{
	//		GetUserError: errors.New("oh no"),
	//	},
	//}

	b := Backend()
	b.iamClient = goodUser
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
		"name":            "test",
		"username":        "jane-doe",
		"rotation_period": 24601,
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Data:      roleData,
		Path:      "static-roles/test",
	}

	// everything good
	//err := b.validateIAMUserExists(context.Background(), roleReq, "jane-doe")
	//if err != nil {
	//	t.Fatalf("couldn't validate user: %s", err)
	//}
	resp, err := b.pathStaticRolesWrite(context.Background(), roleReq, staticRoleFieldData(roleData))
	if err != nil {
		t.Fatalf("couldn't validate an expected good request: %s", err)
	}
	if resp == nil {
		t.Fatalf("didn't get a response from an expected good request")
	}
	//// bad user
	//b.iamClient = badUser
	//err = b.validateIAMUserExists(context.Background(), roleReq, "jane-doe")
	//if err == nil {
	//	t.Fatalf("expected an IAM get user error but didn't get one")
	//}
	//
	//// bad duration
	//err = b.validateRotationPeriod(time.Duration(0))
	//if err == nil {
	//	t.Fatalf("expected duration to be invalid but it was accepted")
	//}
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
