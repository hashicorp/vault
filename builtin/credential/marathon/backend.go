package marathon

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"os"

	marathon "github.com/gambol99/go-marathon"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		Paths: append([]*framework.Path{
			pathConfig(),
			pathLogin(&b),
		}),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config",
			},
			Unauthenticated: []string{
				"login",
			},
		},

		AuthRenew: b.pathLoginRenew,
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}

// Client returns the Marathon client to communicate to Marathon
func (b *backend) Client(marathonUrl string) (marathon.Marathon, error) {
	config := marathon.NewDefaultConfig()
	config.URL = marathonUrl
	config.LogOutput = os.Stdout

	if client, err := marathon.NewClient(config); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

const backendHelp = `
The Marathon credential provider allows task authentication via Marathon.

Tasks provide a marathon_app_id, marathon_app_version and mesos_task_id
and the credential provider can authenticate the task with the Marathon
API.
`
