package azureservicebus

import (
	"fmt"
	"net/url"
	"time"

	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "token/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathTokenRead,
		},

		HelpSynopsis:    pathTokenHelpSyn,
		HelpDescription: pathTokenHelpDesc,
	}
}

func (b *backend) pathTokenRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Get the role
	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	ttl := role.TTL
	// Determine if we have a lease configuration
	if ttl == 0 {
		leaseConfig, err := b.LeaseConfig(req.Storage)
		if err != nil {
			return nil, err
		}
		if leaseConfig == nil {
			leaseConfig = &configLease{}
		}
		ttl = leaseConfig.TTL
	}
	resourceConfig, err := b.ResourceConfig(req.Storage)
	if err != nil {
		return nil, err
	}

	//Encode the SAS token. Reference: https://azure.microsoft.com/en-us/documentation/articles/service-bus-shared-access-signature-authentication/
	uri := strings.ToLower(url.QueryEscape(resourceConfig.ResourceURI))
	expirytime := time.Now().Add(ttl).Unix()
	signstring := fmt.Sprintf("%v\n%v", uri, expirytime)
	signature := url.QueryEscape(ComputeHmac256(signstring, role.SASPolicyKey))
	token := fmt.Sprintf("SharedAccessSignature sr=%v&sig=%v&se=%v&skn=%v", uri, signature, expirytime, role.SASPolicyName)

	// Return the secret. Nothing need to be saved in the secret itself
	resp := b.Secret(SecretTokenType).Response(map[string]interface{}{
		"policy_name": name,
		"token":       token,
	}, map[string]interface{}{})
	resp.Secret.TTL = ttl
	return resp, nil
}

const pathTokenHelpSyn = `
Request a SAS token for a certain role.
`

const pathTokenHelpDesc = `
This path generates a SAS token for a certain role. The
token is generated on demand and will automatically 
expire when the lease is up.
`
