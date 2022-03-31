package pki

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

func updateDefaultKeyId(ctx context.Context, s logical.Storage, id keyId) error {
	config, err := getKeysConfig(ctx, s)
	if err != nil {
		return err
	}

	if config.DefaultKeyId != id {
		return setKeysConfig(ctx, s, &keyConfig{
			DefaultKeyId: id,
		})
	}

	return nil
}

func updateDefaultIssuerId(ctx context.Context, s logical.Storage, id issuerId) error {
	config, err := getIssuersConfig(ctx, s)
	if err != nil {
		return err
	}

	if config.DefaultIssuerId != id {
		return setIssuersConfig(ctx, s, &issuerConfig{
			DefaultIssuerId: id,
		})
	}

	return nil
}
