package azuresecrets

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const walAppKey = "appCreate"

// Eventually expire the WAL if for some reason the rollback operation consistently fails
var maxWALAge = 24 * time.Hour

type walApp struct {
	AppID      string
	AppObjID   string
	Expiration time.Time
}

func (b *azureSecretBackend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	if kind != walAppKey {
		return fmt.Errorf("unknown rollback type %q", kind)
	}
	// Decode the WAL data
	var entry walApp
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339),
		Result:     &entry,
	})
	if err != nil {
		return err
	}
	err = d.Decode(data)
	if err != nil {
		return err
	}

	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return err
	}

	b.Logger().Debug("rolling back SP", "appID", entry.AppID, "appObjID", entry.AppObjID)

	// Attempt to delete the App. deleteApp doesn't return an error if the app isn't
	// found, so no special handling is needed for that case. If we don't succeed within
	// maxWALAge (e.g. client creds have changed and the delete will never succeed),
	// unconditionally remove the WAL.
	if err := client.deleteApp(ctx, entry.AppObjID); err != nil {
		b.Logger().Warn("rollback error deleting App", "err", err)

		if time.Now().After(entry.Expiration) {
			return nil
		}
		return err
	}

	return nil
}
