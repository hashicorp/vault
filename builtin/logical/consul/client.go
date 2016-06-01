package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
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

	consulConf := api.DefaultNonPooledConfig()
	consulConf.Address = conf.Address
	consulConf.Scheme = conf.Scheme
	consulConf.Token = conf.Token

	client, err := api.NewClient(consulConf)
	return client, nil, err
}
