package duo

import (
	"context"
	"fmt"
	"net/url"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/authapi"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type AuthClient interface {
	Preauth(options ...func(*url.Values)) (*authapi.PreauthResult, error)
	Auth(factor string, options ...func(*url.Values)) (*authapi.AuthResult, error)
}

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
			logical.UpdateOperation: pathDuoAccessWrite,
		},

		HelpSynopsis:    pathDuoAccessHelpSyn,
		HelpDescription: pathDuoAccessHelpDesc,
	}
}

func GetDuoAuthClient(ctx context.Context, req *logical.Request, config *DuoConfig) (AuthClient, error) {
	entry, err := req.Storage.Get(ctx, "duo/access")
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

	duoClient := duoapi.NewDuoApi(
		access.IKey,
		access.SKey,
		access.Host,
		config.UserAgent,
	)
	duoAuthClient := authapi.NewAuthApi(*duoClient)
	check, err := duoAuthClient.Check()
	if err != nil {
		return nil, err
	}
	if check == nil {
		return nil, fmt.Errorf("could not connect to Duo; got nil result back from API check call")
	}
	var msg, detail string
	if check.StatResult.Message != nil {
		msg = *check.StatResult.Message
	}
	if check.StatResult.Message_Detail != nil {
		detail = *check.StatResult.Message_Detail
	}
	if check.StatResult.Stat != "OK" {
		return nil, fmt.Errorf("could not connect to Duo: %q (%q)", msg, detail)
	}
	return duoAuthClient, nil
}

func pathDuoAccessWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("duo/access", DuoAccess{
		SKey: d.Get("skey").(string),
		IKey: d.Get("ikey").(string),
		Host: d.Get("host").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
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
