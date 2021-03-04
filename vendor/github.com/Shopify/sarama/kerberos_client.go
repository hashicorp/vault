package sarama

import (
	krb5client "gopkg.in/jcmturner/gokrb5.v7/client"
	krb5config "gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/types"
)

type KerberosGoKrb5Client struct {
	krb5client.Client
}

func (c *KerberosGoKrb5Client) Domain() string {
	return c.Credentials.Domain()
}

func (c *KerberosGoKrb5Client) CName() types.PrincipalName {
	return c.Credentials.CName()
}

/*
*
* Create kerberos client used to obtain TGT and TGS tokens
* used gokrb5 library, which is a pure go kerberos client with
* some GSS-API capabilities, and SPNEGO support. Kafka does not use SPNEGO
* it uses pure Kerberos 5 solution (RFC-4121 and RFC-4120).
*
 */
func NewKerberosClient(config *GSSAPIConfig) (KerberosClient, error) {
	cfg, err := krb5config.Load(config.KerberosConfigPath)
	if err != nil {
		return nil, err
	}
	return createClient(config, cfg)
}

func createClient(config *GSSAPIConfig, cfg *krb5config.Config) (KerberosClient, error) {
	var client *krb5client.Client
	if config.AuthType == KRB5_KEYTAB_AUTH {
		kt, err := keytab.Load(config.KeyTabPath)
		if err != nil {
			return nil, err
		}
		client = krb5client.NewClientWithKeytab(config.Username, config.Realm, kt, cfg)
	} else {
		client = krb5client.NewClientWithPassword(config.Username,
			config.Realm, config.Password, cfg)
	}
	return &KerberosGoKrb5Client{*client}, nil
}
