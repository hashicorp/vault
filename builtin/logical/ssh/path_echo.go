package ssh

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathEcho(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "echo",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathEchoRead,
		},
		HelpSynopsis:    pathEchoHelpSyn,
		HelpDescription: pathEchoHelpDesc,
	}
}

func (b *backend) pathEchoRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{
		Data: map[string]interface{}{
			"echo": "vault-echo",
		},
	}, nil
}

const pathEchoHelpSyn = `
Responds with a echo message.
`

const pathEchoHelpDesc = `
This path will be used by the vault agent running in the
target machine to check if the agent installation is proper.
`
