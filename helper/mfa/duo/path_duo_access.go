package duo

import (
	"fmt"

	"github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/authapi"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathDuoAccess() *framework.Path {
	return &framework.Path{
		Pattern: `duo/access`,
		Fields: map[string]*framework.FieldSchema{
			"skey": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Duo secret key",
			},
			"ikey": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Duo integration key",
			},
			"host": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Duo api host",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation:  pathDuoAccessWrite,
		},

		HelpSynopsis:    pathDuoAccessHelpSyn,
		HelpDescription: pathDuoAccessHelpDesc,
	}
}

func GetDuoAuthClient(req *logical.Request, config *DuoConfig) (*authapi.AuthApi, error) {
	entry, err := req.Storage.Get("duo/access")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf(
			"Duo access credentials haven't been configured. Please configure\n" +
				"them at the 'duo/access' endpoint")
	}
	var access DuoAccess
	if err := entry.DecodeJSON(&access); err != nil {
		return nil, err
	}

	duo_client := duoapi.NewDuoApi(
		access.IKey,
		access.SKey,
		access.Host,
		config.UserAgent,
	)
	duo_auth_client := authapi.NewAuthApi(*duo_client)
	check, err := duo_auth_client.Check()
	if err != nil {
		return nil, err
	}
	if check.StatResult.Stat != "OK" {
		return nil, fmt.Errorf("Could not connect to Duo: %s (%s)", *check.StatResult.Message, *check.StatResult.Message_Detail)
	}
	return duo_auth_client, nil
}

func pathDuoAccessWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("duo/access", DuoAccess{
		SKey: d.Get("skey").(string),
		IKey: d.Get("ikey").(string),
		Host: d.Get("host").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type DuoAccess struct {
	SKey string `json:"skey"`
	IKey string `json:"ikey"`
	Host string `json:"host"`
}

const pathDuoAccessHelpSyn = `
Configure the access keys and host for Duo API connections.
`

const pathDuoAccessHelpDesc = `
To authenticate users with Duo, the backend needs to know what host to connect to
and must authenticate with an integration key and secret key. This endpoint is used
to configure that information.
`
