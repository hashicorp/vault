package aws

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestStaticRolesValidation(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	goodUser := &awsutil.MockIAM{
		GetUserOutput: &iam.GetUserOutput{
			User: &iam.User{
				UserName: aws.String("jane-doe"),
			},
		},
	}

	badUser := &awsutil.MockIAM{
		GetUserError: errors.New("oh no"),
	}

	b := Backend()
	b.iamClient = goodUser
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
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
	err := b.validateIAMUserExists(context.Background(), roleReq, "jane-doe")
	if err != nil {
		t.Fatalf("couldn't validate user: %s", err)
	}

	// bad user
	b.iamClient = badUser
	err = b.validateIAMUserExists(context.Background(), roleReq, "jane-doe")
	if err == nil {
		t.Fatalf("expected an IAM get user error but didn't get one")
	}

	// bad duration
	err = b.validateRotationPeriod(time.Duration(0))
	if err == nil {
		t.Fatalf("expected duration to be invalid but it was accepted")
	}
}
