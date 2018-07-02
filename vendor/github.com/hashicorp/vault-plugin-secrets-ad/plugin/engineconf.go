package plugin

import (
	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/client"
)

type configuration struct {
	PasswordConf *passwordConf
	ADConf       *client.ADConf
}
