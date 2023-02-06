package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

func (b *backend) initQueue(_ context.Context, _ *logical.InitializationRequest) error {
	return nil
}

func (b *backend) rotateExpiredStaticCreds(ctx context.Context, req *logical.Request) error {
	var errs *multierror.Error

	for {
		keepGoing, err := b.rotateCredential(ctx, req.Storage)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if !keepGoing {
			if errs != nil && len(errs.Errors) != 0 {
				return fmt.Errorf("error(s) occurred while rotating expired static credentials: %w", errs)
			} else {
				return nil
			}
		}
	}
}

func (b *backend) rotateCredential(ctx context.Context, storage logical.Storage) (bool, error) {
	// If queue is empty or first item does not need a rotation (priority is next rotation timestamp) there is nothing to do
	item, err := b.credRotationQueue.Pop()
	if err != nil && err == queue.ErrEmpty {
		return false, nil
	}
	if item.Priority > time.Now().Unix() {
		return false, nil
	}

	cfg := item.Value.(staticRoleConfig)

	item.Priority = time.Now().Add(cfg.RotationPeriod).Unix()
	err = b.credRotationQueue.Push(item)
	if err != nil {
		return false, fmt.Errorf("failed to add item into the rotation queue for role '%q': %w", cfg.Name, err)
	}

	err = b.createCredential(ctx, storage, cfg)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (b *backend) createCredential(ctx context.Context, storage logical.Storage, cfg staticRoleConfig) error {
	iamClient, err := b.clientIAM(ctx, storage)
	if err != nil {
		return fmt.Errorf("unable to get the AWS IAM client: %w", err)
	}

	// IAM users can have at most 2 set of keys at a time. If 2 sets already exists we must delete the oldest one to be
	// able to create a new set of credentials.
	accessKeys, err := iamClient.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(cfg.Username),
	})
	if err != nil {
		return fmt.Errorf("unable to list existing access keys for IAM user '%s': %w", cfg.Username, err)
	}
	if len(accessKeys.AccessKeyMetadata) >= 2 {
		oldestKey := accessKeys.AccessKeyMetadata[0]
		if accessKeys.AccessKeyMetadata[1].CreateDate.Before(*accessKeys.AccessKeyMetadata[0].CreateDate) {
			oldestKey = accessKeys.AccessKeyMetadata[1]
		}

		_, err := iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: oldestKey.AccessKeyId,
			UserName:    oldestKey.UserName,
		})
		if err != nil {
			return fmt.Errorf("unable to delete oldest access keys for user '%s': %w", cfg.Username, err)
		}
	}

	// Create new set of keys
	out, err := iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(cfg.Username),
	})
	if err != nil {
		return fmt.Errorf("unable to create new access keys for user '%s': %w", cfg.Username, err)
	}

	// Persist new keys
	entry, err := logical.StorageEntryJSON(formatCredsStoragePath(cfg.Name), &awsCredentials{
		AccessKeyID:     *out.AccessKey.AccessKeyId,
		SecretAccessKey: *out.AccessKey.SecretAccessKey,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal object to JSON: %w", err)
	}
	err = storage.Put(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to save object in storage: %w", err)
	}

	return nil
}

func (b *backend) deleteCredential(ctx context.Context, storage logical.Storage, name string) error {
	return storage.Delete(ctx, formatCredsStoragePath(name))
}
