package openldap

import (
	"fmt"

	"github.com/hashicorp/vault-plugin-secrets-openldap/client"
)

type ldapClient interface {
	Get(conf *client.Config, dn string) (*client.Entry, error)
	UpdatePassword(conf *client.Config, dn string, newPassword string) error
	UpdateRootPassword(conf *client.Config, newPassword string) error
}

func NewClient() *Client {
	return &Client{
		ldap: client.New(),
	}
}

type Client struct {
	ldap client.Client
}

func (c *Client) Get(conf *client.Config, dn string) (*client.Entry, error) {
	filters := map[*client.Field][]string{
		client.FieldRegistry.ObjectClass: {"*"},
	}

	entries, err := c.ldap.Search(conf, dn, filters)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("unable to find entry %s in openldap", dn)
	}
	if len(entries) > 1 {
		return nil, fmt.Errorf("expected one matching entry, but received %d", len(entries))
	}
	return entries[0], nil
}

func (c *Client) UpdatePassword(conf *client.Config, dn string, newPassword string) error {
	filters := map[*client.Field][]string{client.FieldRegistry.ObjectClass: {"*"}}

	newValues, err := client.GetSchemaFieldRegistry(conf.Schema, newPassword)
	if err != nil {
		return fmt.Errorf("error updating password: %s", err)
	}

	return c.ldap.UpdatePassword(conf, dn, newValues, filters)
}

func (c *Client) UpdateRootPassword(conf *client.Config, newPassword string) error {
	filters := map[*client.Field][]string{client.FieldRegistry.ObjectClass: {"*"}}
	newValues := map[*client.Field][]string{client.FieldRegistry.UserPassword: {newPassword}}

	return c.ldap.UpdatePassword(conf, conf.BindDN, newValues, filters)
}
