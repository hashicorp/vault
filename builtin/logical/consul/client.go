package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/logical"
)

func client(s logical.Storage) (*api.Client, error) {
	entry, err := s.Get("config/access")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf(
			"root credentials haven't been configured. Please configure\n" +
				"them at the '/root' endpoint")
	}

	var conf accessConfig
	if err := entry.DecodeJSON(&conf); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %s", err)
	}

	consulConf := api.DefaultNonPooledConfig()
	consulConf.Address = conf.Address
	consulConf.Scheme = conf.Scheme
	consulConf.Token = conf.Token

	return api.NewClient(consulConf)
}
