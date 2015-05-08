package marathon

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"marathon_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The Marathon URL to use for validation",
			},
			"mesos_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The Mesos URL to use for validation",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: pathConfigUpdate,
		},
	}
}

func pathConfigUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config", config{
		MarathonUrl: data.Get("marathon_url").(string),
		MesosUrl:    data.Get("mesos_url").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// Config returns the configuration for this backend.
func (b *backend) Config(s logical.Storage) (*config, error) {
	entry, err := s.Get("config")
	if err != nil {
		return nil, err
	}

	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %s", err)
		}
	}

	return &result, nil
}

type config struct {
	MarathonUrl string `json:"marathon_url"`
	MesosUrl    string `json:"mesos_url"`
}
