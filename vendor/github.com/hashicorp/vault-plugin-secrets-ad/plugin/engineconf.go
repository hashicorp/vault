package plugin

import (
	"github.com/hashicorp/vault/helper/ldaputil"
)

type configuration struct {
	PasswordConf *passwordConf
	ADConf       *ldaputil.ConfigEntry
}
