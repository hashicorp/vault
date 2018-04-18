package roles

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/logical/framework"
)

func newRole(adClient *activedirectory.Client, passwordConf *config.PasswordConf, name string, fieldData *framework.FieldData) (*Role, error) {

	serviceAccountName, err := getServiceAccountName(fieldData)
	if err != nil {
		return nil, err
	}

	// verify service account exists
	_, err = getServiceAccountByName(adClient, serviceAccountName)
	if err != nil {
		return nil, err
	}

	ttl, err := getTTL(passwordConf, fieldData)
	if err != nil {
		return nil, err
	}

	return &Role{
		Name:               name,
		ServiceAccountName: serviceAccountName,
		TTL:                ttl,
	}, nil
}

func getServiceAccountName(fieldData *framework.FieldData) (string, error) {
	serviceAccountName := fieldData.Get("service_account_name").(string)
	if serviceAccountName == "" {
		return "", errors.New("\"service_account_name\" is required")
	}
	return serviceAccountName, nil
}

func getTTL(passwordConf *config.PasswordConf, fieldData *framework.FieldData) (int, error) {

	ttl := fieldData.Get("ttl").(int)
	if ttl == unsetTTL {
		ttl = passwordConf.TTL
	}

	if ttl > passwordConf.MaxTTL {
		return 0, fmt.Errorf("requested ttl of %d seconds is over the max ttl of %d seconds", ttl, passwordConf.MaxTTL)
	}

	if ttl <= 0 {
		return 0, fmt.Errorf("negative ttls are not allowed as they could side-step the preset max ttl")
	}

	return ttl, nil
}

type Role struct {
	Name               string    `json:"name"`
	ServiceAccountName string    `json:"service_account_name"`
	TTL                int       `json:"ttl"`
	LastVaultRotation  time.Time `json:"last_vault_rotation"`
	PasswordLastSet    time.Time `json:"password_last_set"`
}

func (r *Role) Map() map[string]interface{} {
	m := map[string]interface{}{
		"name":                 r.Name,
		"service_account_name": r.ServiceAccountName,
		"ttl": r.TTL,
	}

	var unset time.Time

	if r.LastVaultRotation == unset {
		m["last_vault_rotation"] = nil
	} else {
		m["last_vault_rotation"] = r.LastVaultRotation
	}

	if r.PasswordLastSet == unset {
		m["password_last_set"] = nil
	} else {
		m["password_last_set"] = r.PasswordLastSet
	}

	return m
}
