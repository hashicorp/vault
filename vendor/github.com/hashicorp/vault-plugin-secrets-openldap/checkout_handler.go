// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openldap

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	checkoutStoragePrefix = "checkout/"
	passwordStoragePrefix = "password/"
)

var (
	// errCheckedOut is returned when a check-out request is received
	// for a service account that's already checked out.
	errCheckedOut = errors.New("checked out")

	// errNotFound is used when a requested item doesn't exist.
	errNotFound = errors.New("not found")
)

// CheckOut provides information for a service account that is currently
// checked out.
type CheckOut struct {
	IsAvailable         bool   `json:"is_available"`
	BorrowerEntityID    string `json:"borrower_entity_id"`
	BorrowerClientToken string `json:"borrower_client_token"`
}

// CheckOut attempts to check out a service account. If the account is unavailable, it returns
// errCheckedOut. If the service account isn't managed by this plugin, it returns
// errNotFound.
func (b *backend) CheckOut(ctx context.Context, storage logical.Storage, serviceAccountName string, checkOut *CheckOut) error {
	if ctx == nil {
		return errors.New("ctx must be provided")
	}
	if storage == nil {
		return errors.New("storage must be provided")
	}
	if serviceAccountName == "" {
		return errors.New("service account name must be provided")
	}
	if checkOut == nil {
		return errors.New("check-out must be provided")
	}

	// Check if the service account is currently checked out.
	currentEntry, err := storage.Get(ctx, checkoutStoragePrefix+serviceAccountName)
	if err != nil {
		return err
	}
	if currentEntry == nil {
		return errNotFound
	}
	currentCheckOut := &CheckOut{}
	if err := currentEntry.DecodeJSON(currentCheckOut); err != nil {
		return err
	}
	if !currentCheckOut.IsAvailable {
		return errCheckedOut
	}

	// Since it's not, store the new check-out.
	entry, err := logical.StorageEntryJSON(checkoutStoragePrefix+serviceAccountName, checkOut)
	if err != nil {
		return err
	}
	return storage.Put(ctx, entry)
}

// CheckIn attempts to check in a service account. If an error occurs, the account remains checked out
// and can either be retried by the caller, or eventually may be checked in if it has a ttl
// that ends.
func (b *backend) CheckIn(ctx context.Context, storage logical.Storage, serviceAccountName string) error {
	if ctx == nil {
		return errors.New("ctx must be provided")
	}
	if storage == nil {
		return errors.New("storage must be provided")
	}
	if serviceAccountName == "" {
		return errors.New("service account name must be provided")
	}

	// On check-ins, a new LDAP password is generated, updated in LDAP, and stored.
	config, err := readConfig(ctx, storage)
	if err != nil {
		return err
	}
	if config == nil {
		return errors.New("the config is currently unset")
	}

	newPassword, err := b.GeneratePassword(ctx, config)
	if err != nil {
		return err
	}
	if err := b.client.UpdateUserPassword(config.LDAP, serviceAccountName, newPassword); err != nil {
		return err
	}
	pwdEntry, err := logical.StorageEntryJSON(passwordStoragePrefix+serviceAccountName, newPassword)
	if err != nil {
		return err
	}
	if err := storage.Put(ctx, pwdEntry); err != nil {
		return err
	}

	// That ends the password-handling leg of our journey, now let's deal with the stored check-out itself.
	// Store a check-out status indicating it's available.
	checkOut := &CheckOut{
		IsAvailable: true,
	}
	entry, err := logical.StorageEntryJSON(checkoutStoragePrefix+serviceAccountName, checkOut)
	if err != nil {
		return err
	}
	return storage.Put(ctx, entry)
}

// LoadCheckOut returns either:
//   - A *CheckOut and nil error if the serviceAccountName is currently managed by this engine.
//   - A nil *Checkout and errNotFound if the serviceAccountName is not currently managed by this engine.
func (b *backend) LoadCheckOut(ctx context.Context, storage logical.Storage, serviceAccountName string) (*CheckOut, error) {
	if ctx == nil {
		return nil, errors.New("ctx must be provided")
	}
	if storage == nil {
		return nil, errors.New("storage must be provided")
	}
	if serviceAccountName == "" {
		return nil, errors.New("service account name must be provided")
	}

	entry, err := storage.Get(ctx, checkoutStoragePrefix+serviceAccountName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errNotFound
	}
	checkOut := &CheckOut{}
	if err := entry.DecodeJSON(checkOut); err != nil {
		return nil, err
	}
	return checkOut, nil
}

// DeleteCheckout cleans up anything we were tracking from the service account that
// we will no longer need.
func (b *backend) DeleteCheckout(ctx context.Context, storage logical.Storage, serviceAccountName string) error {
	if ctx == nil {
		return errors.New("ctx must be provided")
	}
	if storage == nil {
		return errors.New("storage must be provided")
	}
	if serviceAccountName == "" {
		return errors.New("service account name must be provided")
	}

	if err := storage.Delete(ctx, passwordStoragePrefix+serviceAccountName); err != nil {
		return err
	}
	return storage.Delete(ctx, checkoutStoragePrefix+serviceAccountName)
}

// retrievePassword is a utility function for grabbing a service account's password from storage.
// retrievePassword will return:
//   - "password", nil if it was successfully able to retrieve the password.
//   - errNotFound if there's no password presently.
//   - Some other error if it was unable to complete successfully.
func retrievePassword(ctx context.Context, storage logical.Storage, serviceAccountName string) (string, error) {
	entry, err := storage.Get(ctx, passwordStoragePrefix+serviceAccountName)
	if err != nil {
		return "", err
	}
	if entry == nil {
		return "", errNotFound
	}
	password := ""
	if err := entry.DecodeJSON(&password); err != nil {
		return "", err
	}
	return password, nil
}
