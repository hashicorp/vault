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

// TestRotation verifies that the rotation code and priority queue correctly selects and rotates credentials
// for static secrets.
func TestRotation(t *testing.T) {
	bgCTX := context.Background()

	type credToInsert struct {
		config staticRoleEntry // role configuration from a normal createRole request
		age    time.Duration   // how old the cred should be - if this is longer than the config.RotationPeriod,
		// the cred is 'pre-expired'

		changed bool // whether we expect the cred to change - this is technically redundant to a comparison between
		// rotationPeriod and age.
	}

	cases := []struct {
		name  string
		creds []credToInsert
	}{
		{
			name: "refresh one",
			creds: []credToInsert{
				{
					config: staticRoleEntry{
						Name:           "test",
						Username:       "jane-doe",
						RotationPeriod: 2 * time.Second,
					},
					age:     5 * time.Second,
					changed: true,
				},
			},
		},
		{
			name: "refresh none",
			creds: []credToInsert{
				{
					config: staticRoleEntry{
						Name:           "test",
						Username:       "jane-doe",
						RotationPeriod: 1 * time.Minute,
					},
					age:     5 * time.Second,
					changed: false,
				},
			},
		},
		{
			name: "refresh one of two",
			creds: []credToInsert{
				{
					config: staticRoleEntry{
						Name:           "test",
						Username:       "jane-doe",
						RotationPeriod: 1 * time.Minute,
					},
					age:     5 * time.Second,
					changed: false,
				},
				{
					config: staticRoleEntry{
						Name:           "toast",
						Username:       "john-doe",
						RotationPeriod: 1 * time.Second,
					},
					age:     5 * time.Second,
					changed: true,
				},
			},
		},
		{
			name:  "no creds to rotate",
			creds: []credToInsert{},
		},
	}

	ak := "long-access-key-id"
	oldSecret := "abcdefghijklmnopqrstuvwxyz"
	newSecret := "zyxwvutsrqponmlkjihgfedcba"

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config := logical.TestBackendConfig()
			config.StorageView = &logical.InmemStorage{}

			b := Backend()

			// this means the creds will be the same for every user, but that's okay
			// since what we care about is whether they changed on a single-user basis.
			miam, err := awsutil.NewMockIAM(
				// blank list for existing user
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{
						{},
					},
				}),
				// initial key to store
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String(ak),
						SecretAccessKey: aws.String(oldSecret),
					},
				}),
			)(nil)
			if err != nil {
				t.Fatalf("couldn't initialze mock IAM handler: %s", err)
			}

			// insert all our creds
			for i, cred := range c.creds {
				b.iamClient = miam

				err = b.createCredential(bgCTX, config.StorageView, cred.config)
				if err != nil {
					t.Fatalf("couldn't insert credential %d: %s", i, err)
				}

				item := &queue.Item{
					Key:      cred.config.Name,
					Value:    cred.config,
					Priority: time.Now().Add(-1 * cred.age).Add(cred.config.RotationPeriod).Unix(),
				}
				err = b.credRotationQueue.Push(item)
				if err != nil {
					t.Fatalf("couldn't push item onto queue: %s", err)
				}
			}

			// update aws responses, same argument for why it's okay every cred will be the same
			miam, err = awsutil.NewMockIAM(
				// old key
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{
						{
							AccessKeyId: aws.String(ak),
						},
					},
				}),
				// new key
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String(ak),
						SecretAccessKey: aws.String(newSecret),
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
			err = b.rotateExpiredStaticCreds(bgCTX, req)
			if err != nil {
				t.Fatalf("got an error rotating credentials: %s", err)
			}

			// check our credentials
			for i, cred := range c.creds {
				entry, err := config.StorageView.Get(bgCTX, formatCredsStoragePath(cred.config.Name))
				if err != nil {
					t.Fatalf("got an error retrieving credentials %d", i)
				}
				var out awsCredentials
				err = entry.DecodeJSON(&out)
				if err != nil {
					t.Fatalf("could not unmarshal storage view entry for cred %d to an aws credential: %s", i, err)
				}

				if cred.changed && out.SecretAccessKey != newSecret {
					t.Fatalf("expected the key for cred %d to have changed, but it hasn't", i)
				} else if !cred.changed && out.SecretAccessKey != oldSecret {
					t.Fatalf("expected the key for cred %d to have stayed the same, but it changed", i)
				}
			}
		})
	}
}
