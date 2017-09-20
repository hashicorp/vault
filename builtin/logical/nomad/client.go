package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/vault/logical"
)

func client(s logical.Storage) (*api.Client, error, error) {
	conf, userErr, intErr := readConfigAccess(s)
	if intErr != nil {
		return nil, nil, intErr
	}
	if userErr != nil {
		return nil, userErr, nil
	}
	if conf == nil {
		return nil, nil, fmt.Errorf("no error received but no configuration found")
	}

	nomadConf := api.DefaultConfig()
	nomadConf.Address = conf.Address
	nomadConf.SecretID = conf.Token

	client, err := api.NewClient(nomadConf)
	return client, nil, err
}
