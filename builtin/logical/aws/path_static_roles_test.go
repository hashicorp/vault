// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
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
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("jane-doe"), UserId: aws.String("unique-id")}}),
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
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("ms-impostor"), UserId: aws.String("fake-id")}}),
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
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("jane-doe"), UserId: aws.String("unique-id")}}),
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
			b := Backend(config)
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
	bgCTX := context.Background()

	cases := []struct {
		name string
		// objects to return from mock IAM.
		// You'll need a GetUserOutput (to validate the existence of the user being written,
		// the keys the user has already been assigned,
		// and the new key vault requests.
		opts []awsutil.MockIAMOption // objects to return from the mock IAM
		// the name, username if updating, and rotation_period of the user. This is the inbound request the cod would get.
		data          map[string]interface{}
		expectedError bool
		findUser      bool
		// if data is sent the name "johnny", then we'll match an existing user with rotation period 24 hours.
		isUpdate    bool
		newPriority int64 // update time of new item in queue, skip if isUpdate false. There is a wiggle room of 5 seconds
		// so the deltas between the old and the new update time should be larger than that to ensure the difference
		// can be detected.
	}{
		{
			name: "happy path",
			opts: []awsutil.MockIAMOption{
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("jane-doe"), UserId: aws.String("unique-id")}}),
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
		{
			name: "update existing user, decreased rotation duration",
			opts: []awsutil.MockIAMOption{
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("john-doe"), UserId: aws.String("unique-id")}}),
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{},
					IsTruncated:       aws.Bool(false),
				}),
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String("abcdefghijklmnopqrstuvwxyz"),
						SecretAccessKey: aws.String("zyxwvutsrqponmlkjihgfedcba"),
						UserName:        aws.String("john-doe"),
					},
				}),
			},
			data: map[string]interface{}{
				"name":            "johnny",
				"rotation_period": "19m",
			},
			findUser:    true,
			isUpdate:    true,
			newPriority: time.Now().Add(19 * time.Minute).Unix(),
		},
		{
			name: "update existing user, increased rotation duration",
			opts: []awsutil.MockIAMOption{
				awsutil.WithGetUserOutput(&iam.GetUserOutput{User: &iam.User{UserName: aws.String("john-doe"), UserId: aws.String("unique-id")}}),
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{},
					IsTruncated:       aws.Bool(false),
				}),
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String("abcdefghijklmnopqrstuvwxyz"),
						SecretAccessKey: aws.String("zyxwvutsrqponmlkjihgfedcba"),
						UserName:        aws.String("john-doe"),
					},
				}),
			},
			data: map[string]interface{}{
				"name":            "johnny",
				"rotation_period": "40h",
			},
			findUser:    true,
			isUpdate:    true,
			newPriority: time.Now().Add(40 * time.Hour).Unix(),
		},
	}

	// if a user exists (user doesn't exist is tested in validation)
	// we'll check how many keys the user has - if it's two, we delete one.

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config := logical.TestBackendConfig()
			config.StorageView = &logical.InmemStorage{}

			miam, err := awsutil.NewMockIAM(
				c.opts...,
			)(nil)
			if err != nil {
				t.Fatal(err)
			}

			b := Backend(config)
			b.iamClient = miam
			if err := b.Setup(bgCTX, config); err != nil {
				t.Fatal(err)
			}

			// put a role in storage for update tests
			staticRole := staticRoleEntry{
				Name:           "johnny",
				Username:       "john-doe",
				ID:             "unique-id",
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

			req := &logical.Request{
				Operation: logical.UpdateOperation,
				Storage:   config.StorageView,
				Data:      c.data,
				Path:      "static-roles/" + c.data["name"].(string),
			}

			r, err := b.pathStaticRolesWrite(bgCTX, req, staticRoleFieldData(req.Data))
			if c.expectedError && err == nil {
				t.Fatal(err)
			} else if c.expectedError {
				return // save us some if statements
			}

			if err != nil {
				t.Fatalf("got an error back unexpectedly: %s", err)
			}

			if c.findUser && r == nil {
				t.Fatal("response was nil, but it shouldn't have been")
			}

			role, err := config.StorageView.Get(bgCTX, req.Path)
			if c.findUser && (err != nil || role == nil) {
				t.Fatalf("couldn't find the role we should have stored: %s", err)
			}
			var actualData staticRoleEntry
			err = role.DecodeJSON(&actualData)
			if err != nil {
				t.Fatalf("couldn't convert storage data to role entry: %s", err)
			}

			// construct expected data
			var expectedData staticRoleEntry
			fieldData := staticRoleFieldData(c.data)
			if c.isUpdate {
				// data is johnny + c.data
				expectedData = staticRole
			}

			var actualItem *queue.Item
			if c.isUpdate {
				actualItem, _ = b.credRotationQueue.PopByKey(expectedData.Name)
			}

			if u, ok := fieldData.GetOk("username"); ok {
				expectedData.Username = u.(string)
			}
			if r, ok := fieldData.GetOk("rotation_period"); ok {
				expectedData.RotationPeriod = time.Duration(r.(int)) * time.Second
			}
			if n, ok := fieldData.GetOk("name"); ok {
				expectedData.Name = n.(string)
			}

			// validate fields
			if eu, au := expectedData.Username, actualData.Username; eu != au {
				t.Fatalf("mismatched username, expected %q but got %q", eu, au)
			}
			if er, ar := expectedData.RotationPeriod, actualData.RotationPeriod; er != ar {
				t.Fatalf("mismatched rotation period, expected %q but got %q", er, ar)
			}
			if en, an := expectedData.Name, actualData.Name; en != an {
				t.Fatalf("mismatched role name, expected %q, but got %q", en, an)
			}

			// one-off to avoid importing/casting
			abs := func(x int64) int64 {
				if x < 0 {
					return -x
				}
				return x
			}

			if c.isUpdate {
				if ep, ap := c.newPriority, actualItem.Priority; abs(ep-ap) > 5 { // 5 second wiggle room for how long the test takes
					t.Fatalf("mismatched updated priority, expected %d but got %d", ep, ap)
				}
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

			b := Backend(config)

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

			b := Backend(config)
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
