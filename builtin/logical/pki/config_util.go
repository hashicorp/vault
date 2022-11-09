package pki

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/sdk/logical"
)

func isDefaultKeySet(ctx context.Context, s logical.Storage) (bool, error) {
	config, err := getKeysConfig(ctx, s)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultKeyId.String()) != "", nil
}

func isDefaultIssuerSet(ctx context.Context, s logical.Storage) (bool, error) {
	config, err := getIssuersConfig(ctx, s)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultIssuerId.String()) != "", nil
}

func updateDefaultKeyId(ctx context.Context, s logical.Storage, id keyID) error {
	config, err := getKeysConfig(ctx, s)
	if err != nil {
		return err
	}

	if config.DefaultKeyId != id {
		return setKeysConfig(ctx, s, &keyConfigEntry{
			DefaultKeyId: id,
		})
	}

	return nil
}

func updateDefaultIssuerId(ctx context.Context, s logical.Storage, id issuerID) error {
	config, err := getIssuersConfig(ctx, s)
	if err != nil {
		return err
	}

	if config.DefaultIssuerId != id {
		config.DefaultIssuerId = id
		return setIssuersConfig(ctx, s, config)
	}

	return nil
}
