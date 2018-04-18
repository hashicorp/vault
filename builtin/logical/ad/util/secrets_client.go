package util

import (
	"fmt"

	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/activedirectory"
)

func NewSecretsClient(logger hclog.Logger, adConf *activedirectory.Configuration) *SecretsClient {
	return &SecretsClient{adClient: activedirectory.NewClient(logger, adConf)}
}

// SecretsClient wraps a *activeDirectory.Client to expose just the common convenience methods needed by the ad secrets backend.
type SecretsClient struct {
	adClient *activedirectory.Client
}

func (c *SecretsClient) Get(serviceAccountName string) (*activedirectory.Entry, error) {

	filters := map[*activedirectory.Field][]string{
		activedirectory.FieldRegistry.UserPrincipalName: {serviceAccountName},
	}

	entries, err := c.adClient.Search(filters)
	if err != nil {
		return nil, err
	}

	if len(entries) <= 0 {
		return nil, fmt.Errorf("service account of %s must already exist in active directory, searches are case sensitive", serviceAccountName)
	}
	if len(entries) > 1 {
		return nil, fmt.Errorf("expected one matching service account, but received %s", entries)
	}
	return entries[0], nil
}

func (c *SecretsClient) GetPasswordLastSet(serviceAccountName string) (time.Time, error) {

	entry, err := c.Get(serviceAccountName)
	if err != nil {
		return time.Time{}, err
	}

	values, found := entry.Get(activedirectory.FieldRegistry.PasswordLastSet)
	if !found {
		return time.Time{}, fmt.Errorf("%s lacks a PasswordLastSet field", entry)
	}

	if len(values) != 1 {
		return time.Time{}, fmt.Errorf("expected only one value for PasswordLastSet, but received %s", values)
	}

	ticks := values[0]
	if ticks == "0" {
		// password has never been rolled in Active Directory, only created
		return time.Time{}, nil
	}

	t, err := activedirectory.ParseTime(ticks)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (c *SecretsClient) UpdatePassword(serviceAccountName string, newPassword string) error {
	filters := map[*activedirectory.Field][]string{
		activedirectory.FieldRegistry.UserPrincipalName: {serviceAccountName},
	}
	return c.adClient.UpdatePassword(filters, newPassword)
}
