// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/protobuf/field_mask"

	multierror "github.com/hashicorp/go-multierror"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func (b *backend) pathKeys() *framework.Path {
	return &framework.Path{
		Pattern: "keys/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationVerb:   "list",
			OperationSuffix: "keys",
		},

		HelpSynopsis:    "List named keys",
		HelpDescription: "List the named keys available for use.",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: withFieldValidator(b.pathKeysList),
		},
	}
}

func (b *backend) pathKeysCRUD() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("key"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloudKMS,
			OperationSuffix: "key",
		},

		HelpSynopsis: "Interact with crypto keys in Vault and Google Cloud KMS",
		HelpDescription: `
This endpoint is used for the CRUD operations for keys in Vault.

To create a new key or to update an existing key, perform a write operation with
the name of the key and the configured parameters below. Vault will also create
or modify the underlying Google Cloud KMS crypto key and store a reference to
it.

    $ vault write gcpkms/keys/my-key \
        key_ring="projects/my-project/locations/global/keyRings/vault" \
        rotation_period="72h" \
        labels="test=true"

To read data about a Google Cloud KMS crypto key, including the key status and
current primary key version, read from the path:

    $ vault read gcpkms/keys/my-key

To delete a key from both Vault and Google Cloud KMS, perform a delete operation
on the name of the key. This will disable automatic rotation of the key in
Google Cloud KMS, disable all crypto key versions for this crypto key in Google
Cloud KMS, and delete Vault's reference to the crypto key.

    $ vault delete gcpkms/keys/my-key

For more information about any of the options, please see the parameter
documentation below. `,

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the key in Vault.
`,
			},

			"algorithm": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Algorithm to use for encryption, decryption, or signing. The value depends on
the key purpose. The value cannot be changed after creation.

For a key purpose of "encrypt_decrypt", the valid values are:

	- symmetric_encryption (default)

For a key purpose of "asymmetric_sign", valid values are:

	- rsa_sign_pss_2048_sha256
	- rsa_sign_pss_3072_sha256
	- rsa_sign_pss_4096_sha256
	- rsa_sign_pkcs1_2048_sha256
	- rsa_sign_pkcs1_3072_sha256
	- rsa_sign_pkcs1_4096_sha256
	- ec_sign_p256_sha256
	- ec_sign_p384_sha384

For a key purpose of "asymmetric_decrypt", valid values are:

	- rsa_decrypt_oaep_2048_sha256
	- rsa_decrypt_oaep_3072_sha256
	- rsa_decrypt_oaep_4096_sha256
`,
			},

			"key_ring": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Full Google Cloud resource ID of the key ring with the project and location
(e.g. projects/my-project/locations/global/keyRings/my-keyring). If the given
key ring does not exist, Vault will try to create it during a create operation.
`,
			},

			"crypto_key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Name of the crypto key to use. If the given crypto key does not exist, Vault
will try to create it. This defaults to the name of the key given to Vault as
the parameter if unspecified.
`,
			},

			"protection_level": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Level of protection to use for the key management. Valid values are "software"
and "hsm". The default value is "software". The value cannot be changed after
creation.
`,
			},

			"purpose": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Purpose of the key. Valid options are "asymmetric_decrypt", "asymmetric_sign",
and "encrypt_decrypt". The default value is "encrypt_decrypt". The value cannot
be changed after creation.
`,
			},

			"rotation_period": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `
Amount of time between crypto key version rotations. This is specified as a time
duration value like 72h (72 hours). The smallest possible value is 24h. This
value only applies to keys with a purpose of "encrypt_decrypt".
`,
			},

			"labels": &framework.FieldSchema{
				Type: framework.TypeKVPairs,
				Description: `
Arbitrary key=value label to apply to the crypto key. To specify multiple
labels, specify this argument multiple times (e.g. labels="a=b" labels="c=d").
`,
			},
		},

		ExistenceCheck: b.pathKeysExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   withFieldValidator(b.pathKeysRead),
			logical.CreateOperation: withFieldValidator(b.pathKeysWrite),
			logical.UpdateOperation: withFieldValidator(b.pathKeysWrite),
			logical.DeleteOperation: withFieldValidator(b.pathKeysDelete),
		},
	}
}

// pathKeysExistenceCheck is used to check if a given key exists.
func (b *backend) pathKeysExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	key := d.Get("key").(string)
	if k, err := b.Key(ctx, req.Storage, key); err != nil || k == nil {
		return false, nil
	}
	return true, nil
}

// pathKeysRead corresponds to GET gcpkms/keys/:name and is used to show
// information about the key.
func (b *backend) pathKeysRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	key := d.Get("key").(string)

	k, err := b.Key(ctx, req.Storage, key)
	if err != nil {
		if err == ErrKeyNotFound {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}

	kmsClient, closer, err := b.KMSClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	cryptoKey, err := kmsClient.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{
		Name: k.CryptoKeyID,
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to read crypto key: {{err}}", err)
	}

	data := map[string]interface{}{
		"id":      cryptoKey.Name,
		"purpose": purposeToString(cryptoKey.Purpose),
	}

	if len(cryptoKey.Labels) > 0 {
		data["labels"] = cryptoKey.Labels
	}
	if cryptoKey.NextRotationTime != nil {
		data["next_rotation_time_seconds"] = cryptoKey.NextRotationTime.Seconds
	}
	if cryptoKey.RotationSchedule != nil {
		if t, ok := cryptoKey.RotationSchedule.(*kmspb.CryptoKey_RotationPeriod); ok && t.RotationPeriod != nil {
			data["rotation_schedule_seconds"] = t.RotationPeriod.Seconds
		}
	}
	if cryptoKey.Primary != nil {
		data["primary_version"] = path.Base(cryptoKey.Primary.Name)
		data["state"] = strings.ToLower(cryptoKey.Primary.State.String())
	}
	if vt := cryptoKey.VersionTemplate; vt != nil {
		data["protection_level"] = protectionLevelToString(vt.ProtectionLevel)
		data["algorithm"] = algorithmToString(vt.Algorithm)
	}

	return &logical.Response{
		Data: data,
	}, nil
}

// pathKeysList corresponds to LIST gcpkms/keys and is used to list all keys
// in the system.
func (b *backend) pathKeysList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keys, err := b.Keys(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(keys), nil
}

// pathKeysWrite corresponds to PUT/POST gcpkms/keys/create/:key and creates a
// new GCP KMS key and registers it for use in Vault.
func (b *backend) pathKeysWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	kmsClient, closer, err := b.KMSClient(req.Storage)
	if err != nil {
		return nil, err
	}
	defer closer()

	key := d.Get("key").(string)
	keyRing := d.Get("key_ring").(string)
	labels := d.Get("labels").(map[string]string)

	// Default crypto key name to the key name if unspecified
	cryptoKey := d.Get("crypto_key").(string)
	if cryptoKey == "" {
		cryptoKey = key
	}

	// Base key
	ck := &kmspb.CryptoKey{
		Labels:          labels,
		VersionTemplate: new(kmspb.CryptoKeyVersionTemplate),
	}

	// Set purpose if given
	if v, ok := d.GetOk("purpose"); ok {
		if req.Operation == logical.UpdateOperation {
			return nil, errImmutable("purpose")
		}

		purpose, ok := keyPurposes[strings.ToLower(v.(string))]
		if !ok {
			return nil, logical.CodedError(400, fmt.Sprintf(
				"unknown purpose %q, valid purposes are %q", v, keyPurposeNames()))
		}
		ck.Purpose = purpose
	} else {
		ck.Purpose = kmspb.CryptoKey_ENCRYPT_DECRYPT
	}

	// Set algorithm if given
	if v, ok := d.GetOk("algorithm"); ok {
		algorithm, ok := keyAlgorithms[strings.ToLower(v.(string))]
		if !ok {
			return nil, logical.CodedError(400, fmt.Sprintf(
				"unknown algorithm %q, valid algorithms are %q", v, keyAlgorithmNames()))
		}
		ck.VersionTemplate.Algorithm = algorithm
	} else {
		if ck.Purpose == kmspb.CryptoKey_ENCRYPT_DECRYPT {
			ck.VersionTemplate.Algorithm = kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION
		} else {
			return nil, errMissingFields("algorithm")
		}
	}

	// Set the protection level
	if v, ok := d.GetOk("protection_level"); ok {
		if req.Operation == logical.UpdateOperation {
			return nil, errImmutable("protection level")
		}

		protectionLevel, ok := keyProtectionLevels[strings.ToLower(v.(string))]
		if !ok {
			return nil, logical.CodedError(400, fmt.Sprintf(
				"unknown protection level %q, valid protection levels are %q", v, keyProtectionLevelNames()))
		}
		ck.VersionTemplate.ProtectionLevel = protectionLevel
	} else {
		ck.VersionTemplate.ProtectionLevel = kmspb.ProtectionLevel_SOFTWARE
	}

	// Set rotation period
	if v, ok := d.GetOk("rotation_period"); ok {
		t := int64(v.(int))

		ck.RotationSchedule = &kmspb.CryptoKey_RotationPeriod{
			RotationPeriod: &duration.Duration{
				Seconds: int64(t),
			},
		}

		ck.NextRotationTime = &timestamp.Timestamp{
			Seconds: time.Now().UTC().Add(time.Duration(t) * time.Second).Unix(),
		}
	}

	// Check if the key ring exists
	kr, err := kmsClient.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{
		Name: keyRing,
	})
	if err != nil {
		if terr, ok := grpcstatus.FromError(err); ok && terr.Code() == grpccodes.NotFound {
			// Key ring does not exist, try to create it
			kr, err = kmsClient.CreateKeyRing(ctx, &kmspb.CreateKeyRingRequest{
				Parent:    path.Dir(path.Dir(keyRing)),
				KeyRingId: path.Base(keyRing),
			})
			if err != nil {
				return nil, errwrap.Wrapf("failed to create key ring: {{err}}", err)
			}
		} else {
			return nil, errwrap.Wrapf("failed to check if key ring exists: {{err}}", err)
		}
	}

	resp, err := kmsClient.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      kr.Name,
		CryptoKeyId: cryptoKey,
		CryptoKey:   ck,
	})
	if err != nil {
		if terr, ok := grpcstatus.FromError(err); ok && terr.Code() == grpccodes.AlreadyExists {
			if req.Operation != logical.UpdateOperation {
				resp := logical.ErrorResponse(
					"cannot update a key that is not already registered - register the " +
						"key first using the /keys/register endpoint, and then update any " +
						"configuration fields.")
				return resp, logical.ErrPermissionDenied
			}

			var paths []string
			ck.Name = fmt.Sprintf("%s/cryptoKeys/%s", kr.Name, cryptoKey)

			if ck.Labels != nil {
				paths = append(paths, "labels")
			}

			if ck.RotationSchedule != nil {
				paths = append(paths, "rotation_period")
			}

			if ck.NextRotationTime != nil {
				paths = append(paths, "next_rotation_time")
			}

			resp, err = kmsClient.UpdateCryptoKey(ctx, &kmspb.UpdateCryptoKeyRequest{
				CryptoKey: ck,
				UpdateMask: &field_mask.FieldMask{
					Paths: paths,
				},
			})
			if err != nil {
				return nil, errwrap.Wrapf("failed to update crypto key: {{err}}", err)
			}
		} else {
			return nil, errwrap.Wrapf("failed to create crypto key: {{err}}", err)
		}
	}

	// Save it
	entry, err := logical.StorageEntryJSON("keys/"+key, &Key{
		Name:        key,
		CryptoKeyID: resp.Name,
	})
	if err != nil {
		return nil, errwrap.Wrapf("failed to create storage entry: {{err}}", err)
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, errwrap.Wrapf("failed to write to storage: {{err}}", err)
	}

	return nil, nil
}

// pathKeysDelete corresponds to PUT/POST gcpkms/keys/delete/:key and deletes an
// existing GCP KMS key and deregisters it from Vault.
func (b *backend) pathKeysDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

	// Disable automatic key rotation
	if _, err := kmsClient.UpdateCryptoKey(ctx, &kmspb.UpdateCryptoKeyRequest{
		CryptoKey: &kmspb.CryptoKey{
			Name:             k.CryptoKeyID,
			NextRotationTime: nil,
			RotationSchedule: nil,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"next_rotation_time", "rotation_period"},
		},
	}); err != nil {
		return nil, errwrap.Wrapf("failed to disable rotation on crypto key: {{err}}", err)
	}

	// Collect the list of all key versions
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

		if resp.State != kmspb.CryptoKeyVersion_DESTROYED &&
			resp.State != kmspb.CryptoKeyVersion_DESTROY_SCHEDULED {
			ckvs = append(ckvs, resp.Name)
		}
	}

	// Iterate over each key version and schedule deletion
	var mu sync.Mutex
	var errs *multierror.Error
	wp := workerpool.New(25)
	for _, ckv := range ckvs {
		ckv := ckv

		wp.Submit(func() {
			if err := retryFib(func() error {
				_, err := kmsClient.DestroyCryptoKeyVersion(ctx, &kmspb.DestroyCryptoKeyVersionRequest{
					Name: ckv,
				})
				return err
			}); err != nil {
				mu.Lock()
				errs = multierror.Append(errs, errwrap.Wrapf(fmt.Sprintf("failed to destroy crypto key version %s: {{err}}", ckv), err))
				mu.Unlock()
			}
		})
	}

	wp.StopWait()

	// Return errors if any happened
	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	// Delete the key from our storage
	if err := req.Storage.Delete(ctx, "keys/"+key); err != nil {
		return nil, errwrap.Wrapf("failed to delete from storage: {{err}}", err)
	}
	return nil, nil
}

// keyPurposes is the list of purposes to key types
var keyPurposes = map[string]kmspb.CryptoKey_CryptoKeyPurpose{
	"asymmetric_decrypt": kmspb.CryptoKey_ASYMMETRIC_DECRYPT,
	"asymmetric_sign":    kmspb.CryptoKey_ASYMMETRIC_SIGN,
	"encrypt_decrypt":    kmspb.CryptoKey_ENCRYPT_DECRYPT,
	"unspecified":        kmspb.CryptoKey_CRYPTO_KEY_PURPOSE_UNSPECIFIED,
}

// keyPurposeNames returns the list of key purposes.
func keyPurposeNames() []string {
	list := make([]string, 0, len(keyPurposes))
	for k := range keyPurposes {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// purposeToString accepts a kmspb and maps that to the user readable purpose.
// Instead of maintaining two maps, this iterates over the purposes map because
// N will always be ridiculously small.
func purposeToString(p kmspb.CryptoKey_CryptoKeyPurpose) string {
	for k, v := range keyPurposes {
		if p == v {
			return k
		}
	}
	return "unspecified"
}

// keyAlgorithms is the list of key algorithms.
var keyAlgorithms = map[string]kmspb.CryptoKeyVersion_CryptoKeyVersionAlgorithm{
	"symmetric_encryption":         kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
	"rsa_sign_pss_2048_sha256":     kmspb.CryptoKeyVersion_RSA_SIGN_PSS_2048_SHA256,
	"rsa_sign_pss_3072_sha256":     kmspb.CryptoKeyVersion_RSA_SIGN_PSS_3072_SHA256,
	"rsa_sign_pss_4096_sha256":     kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA256,
	"rsa_sign_pkcs1_2048_sha256":   kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256,
	"rsa_sign_pkcs1_3072_sha256":   kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_3072_SHA256,
	"rsa_sign_pkcs1_4096_sha256":   kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA256,
	"rsa_decrypt_oaep_2048_sha256": kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_2048_SHA256,
	"rsa_decrypt_oaep_3072_sha256": kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_3072_SHA256,
	"rsa_decrypt_oaep_4096_sha256": kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_4096_SHA256,
	"ec_sign_p256_sha256":          kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256,
	"ec_sign_p384_sha384":          kmspb.CryptoKeyVersion_EC_SIGN_P384_SHA384,
}

// keyAlgorithmNames returns the list of key algorithms.
func keyAlgorithmNames() []string {
	list := make([]string, 0, len(keyAlgorithms))
	for k := range keyAlgorithms {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// algorithmToString accepts a kmspb and maps that to the user readable algorithm.
// Instead of maintaining two maps, this iterates over the algorithms map because
// N will always be ridiculously small.
func algorithmToString(p kmspb.CryptoKeyVersion_CryptoKeyVersionAlgorithm) string {
	for k, v := range keyAlgorithms {
		if p == v {
			return k
		}
	}
	return "unspecified"
}

// keyProtectionLevels is the list of key protection levels.
var keyProtectionLevels = map[string]kmspb.ProtectionLevel{
	"hsm":      kmspb.ProtectionLevel_HSM,
	"software": kmspb.ProtectionLevel_SOFTWARE,
}

// keyProtectionLevelNames returns the list of key protection levels.
func keyProtectionLevelNames() []string {
	list := make([]string, 0, len(keyProtectionLevels))
	for k := range keyProtectionLevels {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// protectionLevelToString accepts a kmspb and maps that to the user readable algorithm.
// Instead of maintaining two maps, this iterates over the algorithms map because
// N will always be ridiculously small.
func protectionLevelToString(p kmspb.ProtectionLevel) string {
	for k, v := range keyProtectionLevels {
		if p == v {
			return k
		}
	}
	return "unknown"
}

// errImmutable is a logical coded error that is returned when the user tries to
// modfiy an immutable field.
func errImmutable(s string) error {
	return logical.CodedError(400, fmt.Sprintf("cannot change %s after key creation", s))
}
