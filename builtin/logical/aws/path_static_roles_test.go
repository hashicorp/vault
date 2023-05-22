package aws

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/queue"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestStaticRolesValidation verifies that valid requests pass validation and that invalid requests fail validation.
// This includes the user already existing in IAM roles, and the rotation period being sufficiently long.
func TestStaticRolesValidation(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	bgCTX := context.Background() // for brevity

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
			if err := b.Setup(bgCTX, config); err != nil {
				t.Fatal(err)
			}
			req := &logical.Request{
				Operation: logical.UpdateOperation,
				Storage:   config.StorageView,
				Data:      c.requestData,
				Path:      "static-roles/test",
			}
			_, err = b.pathStaticRolesWrite(bgCTX, req, staticRoleFieldData(req.Data))
			if c.isError && err == nil {
				t.Fatal("expected an error but didn't get one")
			} else if !c.isError && err != nil {
				t.Fatalf("got an unexpected error: %s", err)
			}
		})
	}
}

// TestStaticRolesWrite validates that we can write a new entry for a new static role, and that we correctly
// do not write if the request is invalid in some way.
func TestStaticRolesWrite(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	bgCTX := context.Background()

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
			if err := b.Setup(bgCTX, config); err != nil {
				t.Fatal(err)
			}

			req := &logical.Request{
				Operation: logical.UpdateOperation,
				Storage:   config.StorageView,
				Data:      c.data,
				Path:      "static-roles/test",
			}

			r, err := b.pathStaticRolesWrite(bgCTX, req, staticRoleFieldData(req.Data))
			if c.expectedError && err == nil {
				t.Fatal(err)
			}
			if c.findUser && r == nil {
				t.Fatal("response was nil, but it shouldn't have been")
			}

			role, err := config.StorageView.Get(bgCTX, req.Path)
			if c.findUser && (err != nil || role == nil) {
				t.Fatalf("couldn't find the role we should have stored: %s", err)
			}
		})
	}
}

// TestStaticRoleRead validates that we can read a configured role and correctly do not read anything if we
// request something that doesn't exist.
func TestStaticRoleRead(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	bgCTX := context.Background()

	// test cases are run against an inmem storage holding a role called "test" attached to an IAM user called "jane-doe"
	cases := []struct {
		name     string
		roleName string
		found    bool
	}{
		{
			name:     "role name exists",
			roleName: "test",
			found:    true,
		},
		{
			name:     "role name not found",
			roleName: "toast",
			found:    false, // implied, but set for clarity
		},
	}

	staticRole := staticRoleEntry{
		Name:           "test",
		Username:       "jane-doe",
		RotationPeriod: 24 * time.Hour,
	}
	entry, err := logical.StorageEntryJSON(formatRoleStoragePath(staticRole.Name), staticRole)
	if err != nil {
		t.Fatal(err)
	}
	err = config.StorageView.Put(bgCTX, entry)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &logical.Request{
				Operation: logical.ReadOperation,
				Storage:   config.StorageView,
				Data: map[string]interface{}{
					"name": c.roleName,
				},
				Path: formatRoleStoragePath(c.roleName),
			}

			b := Backend()

			r, err := b.pathStaticRolesRead(bgCTX, req, staticRoleFieldData(req.Data))
			if err != nil {
				t.Fatal(err)
			}
			if c.found {
				if r == nil {
					t.Fatal("response was nil, but it shouldn't have been")
				}
			} else {
				if r != nil {
					t.Fatal("response should have been nil on a non-existent role")
				}
			}
		})
	}
}

// TestStaticRoleDelete validates that we correctly remove a role on a delete request, and that we correctly do not
// remove anything if a role does not exist with that name.
func TestStaticRoleDelete(t *testing.T) {
	bgCTX := context.Background()

	// test cases are run against an inmem storage holding a role called "test" attached to an IAM user called "jane-doe"
	cases := []struct {
		name  string
		role  string
		found bool
	}{
		{
			name:  "role found",
			role:  "test",
			found: true,
		},
		{
			name:  "role not found",
			role:  "tossed",
			found: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config := logical.TestBackendConfig()
			config.StorageView = &logical.InmemStorage{}

			// fake an IAM
			var iamfunc awsutil.IAMAPIFunc
			if !c.found {
				iamfunc = awsutil.NewMockIAM(awsutil.WithDeleteAccessKeyError(errors.New("shouldn't have called delete")))
			} else {
				iamfunc = awsutil.NewMockIAM()
			}
			miam, err := iamfunc(nil)
			if err != nil {
				t.Fatalf("couldn't initialize mockiam: %s", err)
			}

			b := Backend()
			b.iamClient = miam

			// put in storage
			staticRole := staticRoleEntry{
				Name:           "test",
				Username:       "jane-doe",
				RotationPeriod: 24 * time.Hour,
			}
			entry, err := logical.StorageEntryJSON(formatRoleStoragePath(staticRole.Name), staticRole)
			if err != nil {
				t.Fatal(err)
			}
			err = config.StorageView.Put(bgCTX, entry)
			if err != nil {
				t.Fatal(err)
			}

			l, err := config.StorageView.List(bgCTX, "")
			if err != nil || len(l) != 1 {
				t.Fatalf("couldn't add an entry to storage during test setup: %s", err)
			}

			// put in queue
			err = b.credRotationQueue.Push(&queue.Item{
				Key:      staticRole.Name,
				Value:    staticRole,
				Priority: time.Now().Add(90 * time.Hour).Unix(),
			})
			if err != nil {
				t.Fatalf("couldn't add items to pq")
			}

			req := &logical.Request{
				Operation: logical.ReadOperation,
				Storage:   config.StorageView,
				Data: map[string]interface{}{
					"name": c.role,
				},
				Path: formatRoleStoragePath(c.role),
			}

			r, err := b.pathStaticRolesDelete(bgCTX, req, staticRoleFieldData(req.Data))
			if err != nil {
				t.Fatal(err)
			}
			if r != nil {
				t.Fatal("response wasn't nil, but it should have been")
			}

			l, err = config.StorageView.List(bgCTX, "")
			if err != nil {
				t.Fatal(err)
			}
			if c.found && len(l) != 0 {
				t.Fatal("size of role storage is non zero after delete")
			} else if !c.found && len(l) != 1 {
				t.Fatal("size of role storage changed after what should have been no deletion")
			}

			if c.found && b.credRotationQueue.Len() != 0 {
				t.Fatal("size of queue is non-zero after delete")
			} else if !c.found && b.credRotationQueue.Len() != 1 {
				t.Fatal("size of queue changed after what should have been no deletion")
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
