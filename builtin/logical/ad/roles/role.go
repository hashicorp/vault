package roles

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func newRole(logger hclog.Logger, ctx context.Context, storage logical.Storage, configReader config.Reader, name string, fieldData *framework.FieldData) (*Role, error) {

	engineConf, err := configReader.Config(ctx, storage)
	if err != nil {
		return nil, err
	}
	if engineConf == nil {
		return nil, errors.New("config must be set to create a role")
	}

	adClient := activedirectory.NewClient(logger, engineConf.ADConf)

	serviceAccountName, err := getServiceAccountName(fieldData)
	if err != nil {
		return nil, err
	}

	if err := verifyAccountExists(adClient, serviceAccountName); err != nil {
		return nil, err
	}

	ttl, err := getTTL(engineConf.PasswordConf, fieldData)
	if err != nil {
		return nil, err
	}

	return &Role{
		Name:               name,
		ServiceAccountName: serviceAccountName,
		TTL:                ttl,
	}, nil
}

type Role struct {
	Name               string     `json:"name"`
	ServiceAccountName string     `json:"service_account_name"`
	TTL                int        `json:"ttl"`
	LastVaultRotation  *time.Time `json:"last_vault_rotation,omitempty"`
	PasswordLastSet    *time.Time `json:"password_last_set,omitempty"`
}

func (r *Role) Map() map[string]interface{} {
	return map[string]interface{}{
		"name":                 r.Name,
		"service_account_name": r.ServiceAccountName,
		"ttl": r.TTL,
		"last_vault_rotation": r.LastVaultRotation,
		"password_last_set":   r.PasswordLastSet,
	}
}

func getServiceAccountName(fieldData *framework.FieldData) (string, error) {
	serviceAccountName := fieldData.Get("service_account_name")
	if serviceAccountName == "" {
		return "", errors.New("\"service_account_name\" is required")
	}
	return "", nil
}

func verifyAccountExists(adClient *activedirectory.Client, serviceAccountName string) error {

	filters := map[*activedirectory.Field][]string{
		activedirectory.FieldRegistry.UserPrincipalName: {serviceAccountName},
	}

	entries, err := adClient.Search(filters)
	if err != nil {
		return err
	}

	if len(entries) <= 0 {
		return fmt.Errorf("service account of %s must already exist in active directory, searches are case sensitive", serviceAccountName)
	}
	if len(entries) > 1 {
		return fmt.Errorf("expected one matching service account, but received %s", entries)
	}
	return nil
}

func getTTL(passwordConf *config.PasswordConf, fieldData *framework.FieldData) (int, error) {

	ttl := fieldData.Get("ttl").(int)
	if ttl == unsetTTL {
		ttl = passwordConf.TTL
	}

	if ttl > passwordConf.MaxTTL {
		return 0, fmt.Errorf("requested ttl of %d seconds is over the max ttl of %d seconds", ttl, passwordConf.MaxTTL)
	}

	return ttl, nil
}
