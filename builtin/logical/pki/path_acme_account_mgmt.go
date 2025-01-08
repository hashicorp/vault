// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAcmeMgmtAccountList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "acme/mgmt/account/keyid/?$",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathAcmeMgmtListAccounts,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationPrefix: operationPrefixPKI,
					OperationVerb:   "list-acme-account-keys",
					Description:     "List all ACME account key identifiers.",
				},
			},
		},

		HelpSynopsis:    "List all ACME account key identifiers.",
		HelpDescription: `Allows an operator to list all ACME account key identifiers.`,
	}
}

func pathAcmeMgmtAccountRead(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "acme/mgmt/account/keyid/" + framework.GenericNameRegex("keyid"),
		Fields: map[string]*framework.FieldSchema{
			"keyid": {
				Type:        framework.TypeString,
				Description: "The key identifier of the account.",
				Required:    true,
			},
			"status": {
				Type:          framework.TypeString,
				Description:   "The status of the account.",
				Required:      true,
				AllowedValues: []interface{}{AccountStatusValid.String(), AccountStatusRevoked.String()},
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathAcmeMgmtReadAccount,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationPrefix: operationPrefixPKI,
					OperationSuffix: "acme-key-id",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathAcmeMgmtUpdateAccount,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationPrefix: operationPrefixPKI,
					OperationSuffix: "acme-key-id",
				},
			},
		},

		HelpSynopsis:    "Fetch the details or update the status of an ACME account by key identifier.",
		HelpDescription: `Allows an operator to retrieve details of an ACME account and to update the account status.`,
	}
}

func (b *backend) pathAcmeMgmtListAccounts(ctx context.Context, r *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, r.Storage)

	accountIds, err := b.GetAcmeState().ListAccountIds(sc)
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(accountIds), nil
}

func (b *backend) pathAcmeMgmtReadAccount(ctx context.Context, r *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyId := d.Get("keyid").(string)
	if len(keyId) == 0 {
		return logical.ErrorResponse("keyid is required"), logical.ErrInvalidRequest
	}

	sc := b.makeStorageContext(ctx, r.Storage)
	as := b.GetAcmeState()

	accountEntry, err := as.LoadAccountWithoutDirEnforcement(sc, keyId)
	if err != nil {
		if errors.Is(err, ErrAccountDoesNotExist) {
			return logical.ErrorResponse("ACME key id %s did not exist", keyId), logical.ErrNotFound
		}
		return nil, fmt.Errorf("failed loading ACME account id %q: %w", keyId, err)
	}

	orders, err := as.LoadAccountOrders(sc, accountEntry.KeyId)
	if err != nil {
		return nil, fmt.Errorf("failed loading orders for account %q: %w", accountEntry.KeyId, err)
	}

	orderData := make([]map[string]interface{}, 0, len(orders))
	for _, order := range orders {
		orderData = append(orderData, acmeOrderToDataMap(order))
	}

	dataMap := acmeAccountToDataMap(accountEntry)
	dataMap["orders"] = orderData
	return &logical.Response{Data: dataMap}, nil
}

func (b *backend) pathAcmeMgmtUpdateAccount(ctx context.Context, r *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyId := d.Get("keyid").(string)
	if len(keyId) == 0 {
		return logical.ErrorResponse("keyid is required"), logical.ErrInvalidRequest
	}

	status, err := convertToAccountStatus(d.Get("status"))
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	if status != AccountStatusValid && status != AccountStatusRevoked {
		return logical.ErrorResponse("invalid status %q", status), logical.ErrInvalidRequest
	}

	sc := b.makeStorageContext(ctx, r.Storage)
	as := b.GetAcmeState()

	accountEntry, err := as.LoadAccountWithoutDirEnforcement(sc, keyId)
	if err != nil {
		if errors.Is(err, ErrAccountDoesNotExist) {
			return logical.ErrorResponse("ACME key id %q did not exist", keyId), logical.ErrNotFound
		}
		return nil, fmt.Errorf("failed loading ACME account id %q: %w", keyId, err)
	}

	if accountEntry.Status != status {
		accountEntry.Status = status

		switch status {
		case AccountStatusRevoked:
			accountEntry.AccountRevokedDate = time.Now()
		case AccountStatusValid:
			accountEntry.AccountRevokedDate = time.Time{}
		}

		if err := as.UpdateAccount(sc, accountEntry); err != nil {
			return nil, fmt.Errorf("failed saving account %q: %w", keyId, err)
		}
	}

	dataMap := acmeAccountToDataMap(accountEntry)
	return &logical.Response{Data: dataMap}, nil
}

func convertToAccountStatus(status any) (ACMEAccountStatus, error) {
	if status == nil {
		return "", fmt.Errorf("status is required")
	}

	statusStr, ok := status.(string)
	if !ok {
		return "", fmt.Errorf("status must be a string")
	}

	switch strings.ToLower(strings.TrimSpace(statusStr)) {
	case AccountStatusValid.String():
		return AccountStatusValid, nil
	case AccountStatusRevoked.String():
		return AccountStatusRevoked, nil
	case AccountStatusDeactivated.String():
		return AccountStatusDeactivated, nil
	default:
		return "", fmt.Errorf("invalid status %q", statusStr)
	}
}

func acmeAccountToDataMap(accountEntry *acmeAccount) map[string]interface{} {
	revokedDate := ""
	if !accountEntry.AccountRevokedDate.IsZero() {
		revokedDate = accountEntry.AccountRevokedDate.Format(time.RFC3339)
	}

	eab := map[string]string{}
	if accountEntry.Eab != nil {
		eab["eab_id"] = accountEntry.Eab.KeyID
		eab["directory"] = accountEntry.Eab.AcmeDirectory
		eab["created_time"] = accountEntry.Eab.CreatedOn.Format(time.RFC3339)
		eab["key_type"] = accountEntry.Eab.KeyType
	}

	return map[string]interface{}{
		"key_id":       accountEntry.KeyId,
		"status":       accountEntry.Status,
		"contacts":     accountEntry.Contact,
		"created_time": accountEntry.AccountCreatedDate.Format(time.RFC3339),
		"revoked_time": revokedDate,
		"directory":    accountEntry.AcmeDirectory,
		"eab":          eab,
	}
}

func acmeOrderToDataMap(order *acmeOrder) map[string]interface{} {
	identifiers := make([]string, 0, len(order.Identifiers))
	for _, identifier := range order.Identifiers {
		identifiers = append(identifiers, identifier.Value)
	}
	var certExpiry string
	if !order.CertificateExpiry.IsZero() {
		certExpiry = order.CertificateExpiry.Format(time.RFC3339)
	}
	return map[string]interface{}{
		"order_id":           order.OrderId,
		"status":             string(order.Status),
		"identifiers":        identifiers,
		"cert_serial_number": strings.ReplaceAll(order.CertificateSerialNumber, "-", ":"),
		"cert_expiry":        certExpiry,
		"order_expiry":       order.Expires.Format(time.RFC3339),
	}
}
