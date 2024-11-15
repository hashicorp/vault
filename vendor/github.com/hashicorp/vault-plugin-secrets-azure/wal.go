// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresecrets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	walAppKey            = "appCreate"
	walRotateRootCreds   = "rotateRootCreds"
	walAppRoleAssignment = "appRoleAssign"
)

// Eventually expire the WAL if for some reason the rollback operation consistently fails
var maxWALAge = 24 * time.Hour

func (b *azureSecretBackend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	switch kind {
	case walAppKey:
		return b.rollbackAppWAL(ctx, req, data)
	case walRotateRootCreds:
		return b.rollbackRootWAL(ctx, req, data)
	case walAppRoleAssignment:
		return b.rollbackRoleAssignWAL(ctx, req, data)
	default:
		return fmt.Errorf("unknown rollback type %q", kind)
	}
}

type walApp struct {
	AppID      string
	AppObjID   string
	Expiration time.Time
}

func (b *azureSecretBackend) rollbackAppWAL(ctx context.Context, req *logical.Request, data interface{}) error {
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
	if err := client.deleteApp(ctx, entry.AppObjID, true); err != nil {
		// Prevent noisy logs for non-existent or deleted out-of-band errors
		respErr := new(azcore.ResponseError)
		if errors.As(err, &respErr) && (respErr.StatusCode == http.StatusNoContent || respErr.StatusCode == http.StatusNotFound) {
			b.Logger().Trace("app already deleted or does not exist", "err", err.Error())
		} else {
			b.Logger().Warn("rollback error deleting App", "err", err)
		}

		if time.Now().After(entry.Expiration) {
			b.Logger().Warn("app WAL expired prior to rollback; resources may still exist")
			return nil
		}
		return err
	}

	return nil
}

type walRotateRoot struct{}

func (b *azureSecretBackend) rollbackRootWAL(ctx context.Context, req *logical.Request, data interface{}) error {
	b.Logger().Debug("rolling back config")
	config, err := b.getConfig(ctx, req.Storage)
	if err != nil {
		return err
	}

	config.NewClientSecret = ""
	config.NewClientSecretCreated = time.Time{}
	config.NewClientSecretExpirationDate = time.Time{}
	config.NewClientSecretKeyID = ""

	err = b.saveConfig(ctx, config, req.Storage)
	if err != nil {
		return err
	}

	b.updatePassword = false

	return nil
}

type walAppRoleAssign struct {
	SpID          string
	AssignmentIDs []string
	AzureRoles    []*AzureRole
	Expiration    time.Time
}

func (b *azureSecretBackend) rollbackRoleAssignWAL(ctx context.Context, req *logical.Request, data interface{}) error {
	// Decode the WAL data
	var entry walAppRoleAssign
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

	b.Logger().Debug("rolling back role assignments for service principal", "ID", entry.SpID)

	// Return if there aren't any roles to unassign
	if entry.AzureRoles == nil {
		b.Logger().Error("no azure roles associated with role")
		return nil
	}

	// Assemble all App Role Assignment IDs
	var roleAssignments []string
	for i, assignmentID := range entry.AssignmentIDs {
		if entry.AzureRoles[i] == nil {
			return fmt.Errorf("azure role was nil")
		}
		roleAssignments = append(roleAssignments, fmt.Sprintf("%s/providers/Microsoft.Authorization/roleAssignments/%s",
			entry.AzureRoles[i].Scope,
			assignmentID))
	}

	// Check any errors to filter out expected responses. Azure will return
	// a 204 when trying to delete a role assignment that has already been
	// deleted, or does not exist. We may hit this case during rollback.
	if err := client.unassignRoles(ctx, roleAssignments); err != nil {
		for _, e := range err.(*multierror.Error).Errors {
			switch {
			case strings.Contains(e.Error(), "StatusCode=204"):
				b.Logger().Trace("role assignment already deleted or does not exist", "err", e.Error())
			default:
				return fmt.Errorf("rollback error unassinging role: %w", e)
			}
		}
		if time.Now().After(entry.Expiration) {
			b.Logger().Warn("role assignment WAL expired prior to rollback; resources may still exist")
			return nil
		}
	}
	return nil
}
