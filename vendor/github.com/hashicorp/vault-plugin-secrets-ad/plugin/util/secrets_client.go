package util

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/client"
)

func NewSecretsClient(logger hclog.Logger) *SecretsClient {
	return &SecretsClient{adClient: client.NewClient(logger)}
}

// SecretsClient wraps a *activeDirectory.activeDirectoryClient to expose just the common convenience methods needed by the ad secrets backend.
type SecretsClient struct {
	adClient *client.Client
}

func (c *SecretsClient) Get(conf *client.ADConf, serviceAccountName string) (*client.Entry, error) {
	filters := map[*client.Field][]string{
		client.FieldRegistry.UserPrincipalName: {serviceAccountName},
	}

	entries, err := c.adClient.Search(conf, filters)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("unable to find service account named %s in active directory, searches are case sensitive", serviceAccountName)
	}
	if len(entries) > 1 {
		return nil, fmt.Errorf("expected one matching service account, but received %s", entries)
	}
	return entries[0], nil
}

func (c *SecretsClient) GetPasswordLastSet(conf *client.ADConf, serviceAccountName string) (time.Time, error) {
	entry, err := c.Get(conf, serviceAccountName)
	if err != nil {
		return time.Time{}, err
	}

	values, found := entry.Get(client.FieldRegistry.PasswordLastSet)
	if !found {
		return time.Time{}, fmt.Errorf("%+v lacks a PasswordLastSet field", entry)
	}

	if len(values) != 1 {
		return time.Time{}, fmt.Errorf("expected only one value for PasswordLastSet, but received %s", values)
	}

	ticks := values[0]
	if ticks == "0" {
		// password has never been rolled in Active Directory, only created
		return time.Time{}, nil
	}

	t, err := client.ParseTicks(ticks)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (c *SecretsClient) UpdatePassword(conf *client.ADConf, serviceAccountName string, newPassword string) error {
	filters := map[*client.Field][]string{
		client.FieldRegistry.UserPrincipalName: {serviceAccountName},
	}
	return c.adClient.UpdatePassword(conf, filters, newPassword)
}

func (c *SecretsClient) UpdateRootPassword(conf *client.ADConf, bindDN string, newPassword string) error {
	filters := map[*client.Field][]string{
		client.FieldRegistry.DistinguishedName: {bindDN},
	}
	return c.adClient.UpdatePassword(conf, filters, newPassword)
}
