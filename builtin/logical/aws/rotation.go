// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

// rotateExpiredStaticCreds will pop expired credentials (credentials whose priority
// represents a time before the present), rotate the associated credential, and push
// them back onto the queue with the new priority.
func (b *backend) rotateExpiredStaticCreds(ctx context.Context, req *logical.Request) error {
	var errs *multierror.Error

	for {
		keepGoing, err := b.rotateCredential(ctx, req.Storage)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if !keepGoing {
			if errs.ErrorOrNil() != nil {
				return fmt.Errorf("error(s) occurred while rotating expired static credentials: %w", errs)
			} else {
				return nil
			}
		}
	}
}

// rotateCredential pops an element from the priority queue, and if it is expired, rotate and re-push.
// If a cred was ready for rotation, return true, otherwise return false.
func (b *backend) rotateCredential(ctx context.Context, storage logical.Storage) (wasReady bool, err error) {
	// If queue is empty or first item does not need a rotation (priority is next rotation timestamp) there is nothing to do
	item, err := b.credRotationQueue.Pop()
	if err != nil {
		// the queue is just empty, which is fine.
		if errors.Is(err, queue.ErrEmpty) {
			return false, nil
		}
		return false, fmt.Errorf("failed to pop from queue for role %q: %w", item.Key, err)
	}
	if item.Priority > time.Now().Unix() {
		// no rotation required
		// push the item back into priority queue
		err = b.credRotationQueue.Push(item)
		if err != nil {
			return false, fmt.Errorf("failed to add item into the rotation queue for role %q: %w", item.Key, err)
		}
		return false, nil
	}

	b.Logger().Debug("rotating credential", "role", item.Key)
	cfg := item.Value.(staticRoleEntry)

	creds, err := b.createCredential(ctx, storage, cfg, true)
	if err != nil {
		b.Logger().Error("failed to create credential, re-queueing", "error", err)
		// put it back in the queue with a backoff
		item.Priority = time.Now().Add(10 * time.Second).Unix()
		innerErr := b.credRotationQueue.Push(item)
		if innerErr != nil {
			return true, fmt.Errorf("failed to add item into the rotation queue for role %q(%w), while attempting to recover from failure to create credential: %w", cfg.Name, innerErr, err)
		}
		// there was one that "should have" rotated, so we want to keep looking further down the queue
		return true, err
	}

	// set new priority and re-queue
	item.Priority = creds.priority(cfg)
	err = b.credRotationQueue.Push(item)
	if err != nil {
		return true, fmt.Errorf("failed to add item into the rotation queue for role %q: %w", cfg.Name, err)
	}

	return true, nil
}

// createCredential will create a new iam credential, deleting the oldest one if necessary.
func (b *backend) createCredential(ctx context.Context, storage logical.Storage, cfg staticRoleEntry, shouldLockStorage bool) (*awsCredentials, error) {
	// Always create a fresh client
	iamClient, err := b.getNonCachedIAMClient(ctx, storage, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get IAM client for role %q: %w", cfg.Name, err)
	}

	// IAM users can have a most 2 sets of keys at a time.
	// (https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_iam-quotas.html)
	// Ideally we would get this value through an api check, but I'm not sure one exists.
	const maxAllowedKeys = 2

	err = b.validateIAMUserExists(ctx, storage, &cfg, false)
	if err != nil {
		return nil, fmt.Errorf("iam user didn't exist, or username/userid didn't match: %w", err)
	}

	accessKeys, err := iamClient.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(cfg.Username),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list existing access keys for IAM user %q: %w", cfg.Username, err)
	}

	// If we have the maximum number of keys, we have to delete one to make another (so we can get the credentials).
	// We'll delete the oldest one.
	//
	// Since this check relies on a pre-coded maximum, it's a bit fragile. If the number goes up, we risk deleting
	// a key when we didn't need to. If this number goes down, we'll start throwing errors because we think we're
	// allowed to create a key and aren't. In either case, adjusting the constant should be sufficient to fix things.
	if len(accessKeys.AccessKeyMetadata) >= maxAllowedKeys {
		oldestKey := accessKeys.AccessKeyMetadata[0]

		for i := 1; i < len(accessKeys.AccessKeyMetadata); i++ {
			if accessKeys.AccessKeyMetadata[i].CreateDate.Before(*oldestKey.CreateDate) {
				oldestKey = accessKeys.AccessKeyMetadata[i]
			}
		}

		_, err := iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: oldestKey.AccessKeyId,
			UserName:    oldestKey.UserName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to delete oldest access keys for user %q: %w", cfg.Username, err)
		}
	}

	// Create new set of keys
	out, err := iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(cfg.Username),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create new access keys for user %q: %w", cfg.Username, err)
	}
	expiration := time.Now().UTC().Add(cfg.RotationPeriod)

	creds := &awsCredentials{
		AccessKeyID:     *out.AccessKey.AccessKeyId,
		SecretAccessKey: *out.AccessKey.SecretAccessKey,
		Expiration:      &expiration,
	}
	// Persist new keys
	entry, err := logical.StorageEntryJSON(formatCredsStoragePath(cfg.Name), creds)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object to JSON: %w", err)
	}
	if shouldLockStorage {
		b.roleMutex.Lock()
		defer b.roleMutex.Unlock()
	}
	err = storage.Put(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("failed to save object in storage: %w", err)
	}

	return creds, nil
}

// delete credential will remove the credential associated with the role from storage.
func (b *backend) deleteCredential(ctx context.Context, storage logical.Storage, cfg staticRoleEntry, shouldLockStorage bool) error {
	// synchronize storage access if we didn't in the caller.
	if shouldLockStorage {
		b.roleMutex.Lock()
		defer b.roleMutex.Unlock()
	}

	key, err := storage.Get(ctx, formatCredsStoragePath(cfg.Name))
	if err != nil {
		return fmt.Errorf("couldn't find key in storage: %w", err)
	}
	// no entry, so i guess we deleted it already
	if key == nil {
		return nil
	}
	var creds awsCredentials
	err = key.DecodeJSON(&creds)
	if err != nil {
		return fmt.Errorf("couldn't decode storage entry to a valid credential: %w", err)
	}

	err = storage.Delete(ctx, formatCredsStoragePath(cfg.Name))
	if err != nil {
		return fmt.Errorf("couldn't delete from storage: %w", err)
	}

	iamClient, err := b.nonCachedClientIAM(ctx, storage, b.Logger(), &cfg)
	if err != nil {
		return fmt.Errorf("failed to get IAM client for role %q while deleting: %w", cfg.Name, err)
	}

	// because we have the information, this is the one we created, so it's safe for us to delete.
	_, err = iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(creds.AccessKeyID),
		UserName:    aws.String(cfg.Username),
	})
	if err != nil {
		return fmt.Errorf("couldn't delete from IAM: %w", err)
	}

	return nil
}
