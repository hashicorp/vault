// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/iterator"

	multierror "github.com/hashicorp/go-multierror"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func (b *backend) pathKeysTrim() *framework.Path {
	return &framework.Path{
		Pattern: "keys/trim/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "trim",
		},

		HelpSynopsis: "Delete old crypto key versions from Google Cloud KMS",
		HelpDescription: `
This endpoint deletes old crypto key versions from Google Cloud KMS that are
older than the key's min_version. If min_version is unset, no keys are deleted.

To trim a collection of keys, first decide on the minimum version which you want
to allow:

    $ vault write gcpkms/keys/config/my-key \
        min_version=42

Then execute the trim call:

    $ vault write -f gcpkms/keys/trim/my-key

This will delete all crypto key versions from Google Cloud KMS which are older
than the specified version (version 42 in this example). Note that this will
make it impossible to decrypt data previously encrypted with these older keys
through conventional methods.
`,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key in Vault.
`,
			},
		},

		ExistenceCheck: b.pathKeysExistenceCheck,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathKeysTrimWrite),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "key-versions",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathKeysTrimWrite),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "key-versions",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: withFieldValidator(b.pathKeysTrimWrite),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "key-versions2",
				},
			},
		},
	}
}

// pathKeysTrimWrite corresponds to PUT/POST/DELETE gcpkms/keys/trim/:key and
// deletes all crypto key versions from Google Cloud KMS which are older than
// the key's min_version.
func (b *backend) pathKeysTrimWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	kmsClient, closer, err := b.KMSClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	key := d.Get("key").(string)

	k, err := b.Key(ctx, req.Storage, key)
	if err != nil {
		if err == ErrKeyNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	// If a min version was not set, there's no point in iterating
	if k.MinVersion < 1 {
		return nil, nil
	}

	// Collect the list of all key versions
	var errs *multierror.Error
	var ckvs []string
	it := kmsClient.ListCryptoKeyVersions(ctx, &kmspb.ListCryptoKeyVersionsRequest{
		Parent: k.CryptoKeyID,
	})
	for {
		resp, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, errwrap.Wrapf("failed to list crypto key versions: {{err}}", err)
		}

		if resp.State == kmspb.CryptoKeyVersion_DESTROYED ||
			resp.State == kmspb.CryptoKeyVersion_DESTROY_SCHEDULED {
			continue
		}

		parts := strings.Split(resp.Name, "/")
		if len(parts) < 1 {
			errs = multierror.Append(errs,
				fmt.Errorf("failed to delete crypto key version %s: malformed name", resp.Name))
		}

		v, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			errs = multierror.Append(errs,
				fmt.Errorf("failed to delete crypto key version %s: not an integer version", resp.Name))
		}

		if v < k.MinVersion {
			ckvs = append(ckvs, resp.Name)
		}
	}

	// Iterate over each key version and schedule deletion
	var mu sync.Mutex
	wp := workerpool.New(25)
	for _, ckv := range ckvs {
		ckv := ckv

		wp.Submit(func() {
			if _, err := kmsClient.DestroyCryptoKeyVersion(ctx, &kmspb.DestroyCryptoKeyVersionRequest{
				Name: ckv,
			}); err != nil {
				mu.Lock()
				errs = multierror.Append(errs, errwrap.Wrapf(fmt.Sprintf(
					"failed to delete crypto key version %s: {{err}}", ckv), err))
				mu.Unlock()
			}
		})
	}

	wp.StopWait()

	// Return an error if any failed
	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	return nil, nil
}
